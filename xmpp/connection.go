// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"reflect"
	"strconv"
	"sync"
)

// Conn represents a connection to an XMPP server.
type Conn struct {
	config *Config

	out     io.Writer
	rawOut  io.Writer // doesn't log. Used for <auth>
	in      *xml.Decoder
	jid     string
	archive bool
	Rand    io.Reader

	lock          sync.Mutex
	inflights     map[Cookie]inflight
	customStorage map[xml.Name]reflect.Type
}

// NewConn creates a new connection
func NewConn(in *xml.Decoder, out io.Writer, jid string) *Conn {
	return &Conn{
		in:  in,
		out: out,
		jid: jid,
	}
}

// Close closes the underlying connection
func (c *Conn) Close() error {
	return c.config.Conn.Close()
}

// Next reads stanzas from the server. If the stanza is a reply, it dispatches
// it to the correct channel and reads the next message. Otherwise it returns
// the stanza for processing.
func (c *Conn) Next() (stanza Stanza, err error) {
	for {
		if stanza.Name, stanza.Value, err = next(c); err != nil {
			return
		}

		if iq, ok := stanza.Value.(*ClientIQ); ok && (iq.Type == "result" || iq.Type == "error") {
			var cookieValue uint64
			if cookieValue, err = strconv.ParseUint(iq.ID, 16, 64); err != nil {
				err = errors.New("xmpp: failed to parse id from iq: " + err.Error())
				return
			}
			cookie := Cookie(cookieValue)

			c.lock.Lock()
			inflight, ok := c.inflights[cookie]
			c.lock.Unlock()

			if !ok {
				continue
			}
			if len(inflight.to) > 0 {
				// The reply must come from the address to
				// which we sent the request.
				if inflight.to != iq.From {
					continue
				}
			} else {
				// If there was no destination on the request
				// then the matching is more complex because
				// servers differ in how they construct the
				// reply.
				if len(iq.From) > 0 && iq.From != c.jid && iq.From != RemoveResourceFromJid(c.jid) && iq.From != domainFromJid(c.jid) {
					continue
				}
			}

			c.lock.Lock()
			delete(c.inflights, cookie)
			c.lock.Unlock()

			inflight.replyChan <- stanza
			continue
		}

		return
	}
}

// Cancel cancels and outstanding request. The request's channel is closed.
func (c *Conn) Cancel(cookie Cookie) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	inflight, ok := c.inflights[cookie]
	if !ok {
		return false
	}

	delete(c.inflights, cookie)
	close(inflight.replyChan)
	return true
}

// Dial creates a new connection to an XMPP server, authenticates as the
// given user.
func Dial(address, user, domain, password string, config *Config) (c *Conn, err error) {
	c = new(Conn)
	c.config = config
	c.inflights = make(map[Cookie]inflight)
	c.archive = config.Archive

	log := ioutil.Discard
	if config != nil && config.Log != nil {
		log = config.Log
	}

	var conn net.Conn
	if config != nil && config.Conn != nil {
		conn = config.Conn
	} else {
		io.WriteString(log, "Making TCP connection to "+address+"\n")

		if conn, err = net.Dial("tcp", address); err != nil {
			return nil, err
		}
	}

	c.in, c.out = makeInOut(conn, config)

	features, err := c.getFeatures(domain)
	if err != nil {
		return nil, err
	}

	if !config.SkipTLS {
		if features.StartTLS.XMLName.Local == "" {
			return nil, errors.New("xmpp: server doesn't support TLS")
		}

		fmt.Fprintf(c.out, "<starttls xmlns='%s'/>", NsTLS)

		proceed, err := nextStart(c.in)
		if err != nil {
			return nil, err
		}
		if proceed.Name.Space != NsTLS || proceed.Name.Local != "proceed" {
			return nil, errors.New("xmpp: expected <proceed> after <starttls> but got <" + proceed.Name.Local + "> in " + proceed.Name.Space)
		}

		io.WriteString(log, "Starting TLS handshake\n")

		haveCertHash := len(config.ServerCertificateSHA256) != 0

		var tlsConfig tls.Config
		if config.TLSConfig != nil {
			tlsConfig = *config.TLSConfig
		}
		tlsConfig.ServerName = domain
		tlsConfig.InsecureSkipVerify = true

		tlsConn := tls.Client(conn, &tlsConfig)
		if err := tlsConn.Handshake(); err != nil {
			return nil, err
		}

		tlsState := tlsConn.ConnectionState()
		printTLSDetails(log, tlsState)

		if haveCertHash {
			h := sha256.New()
			h.Write(tlsState.PeerCertificates[0].Raw)
			if digest := h.Sum(nil); !bytes.Equal(digest, config.ServerCertificateSHA256) {
				return nil, fmt.Errorf("xmpp: server certificate does not match expected hash (got: %x, want: %x)", digest, config.ServerCertificateSHA256)
			}
		} else {
			if len(tlsState.PeerCertificates) == 0 {
				return nil, errors.New("xmpp: server has no certificates")
			}

			opts := x509.VerifyOptions{
				Intermediates: x509.NewCertPool(),
			}
			for _, cert := range tlsState.PeerCertificates[1:] {
				opts.Intermediates.AddCert(cert)
			}
			verifiedChains, err := tlsState.PeerCertificates[0].Verify(opts)
			if err != nil {
				return nil, errors.New("xmpp: failed to verify TLS certificate: " + err.Error())
			}

			for i, cert := range verifiedChains[0] {
				fmt.Fprintf(log, "  certificate %d: %s\n", i, certName(cert))
			}
			leafCert := verifiedChains[0][0]

			if err := leafCert.VerifyHostname(domain); err != nil {
				if config.TrustedAddress {
					fmt.Fprintf(log, "Certificate fails to verify against domain in username: %s\n", err)
					host, _, err := net.SplitHostPort(address)
					if err != nil {
						return nil, errors.New("xmpp: failed to split address when checking whether TLS certificate is valid: " + err.Error())
					}
					if err = leafCert.VerifyHostname(host); err != nil {
						return nil, errors.New("xmpp: failed to match TLS certificate to address after failing to match to username: " + err.Error())
					}
					fmt.Fprintf(log, "Certificate matches against trusted server hostname: %s\n", host)
				} else {
					return nil, errors.New("xmpp: failed to match TLS certificate to name: " + err.Error())
				}
			}
		}

		c.in, c.out = makeInOut(tlsConn, config)
		c.rawOut = tlsConn

		if features, err = c.getFeatures(domain); err != nil {
			return nil, err
		}
	} else {
		c.rawOut = conn
	}

	if config != nil && config.CreateCallback != nil {
		io.WriteString(log, "Attempting to create account\n")
		fmt.Fprintf(c.out, "<iq type='get' id='create_1'><query xmlns='jabber:iq:register'/></iq>")
		var iq ClientIQ
		if err = c.in.DecodeElement(&iq, nil); err != nil {
			return nil, errors.New("unmarshal <iq>: " + err.Error())
		}
		if iq.Type != "result" {
			return nil, errors.New("xmpp: account creation failed")
		}
		var register RegisterQuery
		if err := xml.NewDecoder(bytes.NewBuffer(iq.Query)).Decode(&register); err != nil {
			return nil, err
		}

		if len(register.Form.Type) > 0 {
			reply, err := processForm(&register.Form, register.Datas, config.CreateCallback)
			fmt.Fprintf(c.rawOut, "<iq type='set' id='create_2'><query xmlns='jabber:iq:register'>")
			if err = xml.NewEncoder(c.rawOut).Encode(reply); err != nil {
				return nil, err
			}
			fmt.Fprintf(c.rawOut, "</query></iq>")
		} else if register.Username != nil && register.Password != nil {
			// Try the old-style registration.
			fmt.Fprintf(c.rawOut, "<iq type='set' id='create_2'><query xmlns='jabber:iq:register'><username>%s</username><password>%s</password></query></iq>", user, password)
		}
		var iq2 ClientIQ
		if err = c.in.DecodeElement(&iq2, nil); err != nil {
			return nil, errors.New("unmarshal <iq>: " + err.Error())
		}
		if iq2.Type == "error" {
			return nil, errors.New("xmpp: account creation failed")
		}
	}

	io.WriteString(log, "Authenticating as "+user+"\n")
	if err := c.authenticate(features, user, password); err != nil {
		return nil, err
	}
	io.WriteString(log, "Authentication successful\n")

	if features, err = c.getFeatures(domain); err != nil {
		return nil, err
	}

	// Send IQ message asking to bind to the local user name.
	fmt.Fprintf(c.out, "<iq type='set' id='bind_1'><bind xmlns='%s'/></iq>", NsBind)
	var iq ClientIQ
	if err = c.in.DecodeElement(&iq, nil); err != nil {
		return nil, errors.New("unmarshal <iq>: " + err.Error())
	}
	c.jid = iq.Bind.Jid // our local id

	if features.Session != nil {
		// The server needs a session to be established. See RFC 3921,
		// section 3.
		fmt.Fprintf(c.out, "<iq to='%s' type='set' id='sess_1'><session xmlns='%s'/></iq>", domain, NsSession)
		if err = c.in.DecodeElement(&iq, nil); err != nil {
			return nil, errors.New("xmpp: unmarshal <iq>: " + err.Error())
		}
		if iq.Type != "result" {
			return nil, errors.New("xmpp: session establishment failed")
		}
	}
	return c, nil
}

// Send sends an IM message to the given user.
func (c *Conn) Send(to, msg string) error {
	archive := ""
	if !c.archive {
		// The first part of archive is from google:
		// See https://developers.google.com/talk/jep_extensions/otr
		// The second part of the stanza is from XEP-0136
		// http://xmpp.org/extensions/xep-0136.html#pref-syntax-item-otr
		// http://xmpp.org/extensions/xep-0136.html#otr-nego
		archive = "<nos:x xmlns:nos='google:nosave' value='enabled'/><arc:record xmlns:arc='http://jabber.org/protocol/archive' otr='require'/>"
	}
	_, err := fmt.Fprintf(c.out, "<message to='%s' from='%s' type='chat'><body>%s</body>%s</message>", xmlEscape(to), xmlEscape(c.jid), xmlEscape(msg), archive)
	return err
}

func makeInOut(conn io.ReadWriter, config *Config) (in *xml.Decoder, out io.Writer) {
	if config != nil && config.InLog != nil {
		in = xml.NewDecoder(io.TeeReader(conn, config.InLog))
	} else {
		in = xml.NewDecoder(conn)
	}

	if config != nil && config.OutLog != nil {
		out = io.MultiWriter(conn, config.OutLog)
	} else {
		out = conn
	}

	return
}
