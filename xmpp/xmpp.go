// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"reflect"
	"sort"
	"strconv"
	"strings"
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

	lock          sync.Mutex
	inflights     map[Cookie]inflight
	customStorage map[xml.Name]reflect.Type
}

func (conn *Conn) Close() error {
	return conn.config.Conn.Close()
}

// inflight contains the details of a pending request to which we are awaiting
// a reply.
type inflight struct {
	// replyChan is the channel to which we'll send the reply.
	replyChan chan<- Stanza
	// to is the address to which we sent the request.
	to string
}

// Stanza represents a message from the XMPP server.
type Stanza struct {
	Name  xml.Name
	Value interface{}
}

// Cookie is used to give a unique identifier to each request.
type Cookie uint64

func (c *Conn) getCookie() Cookie {
	var buf [8]byte
	if _, err := rand.Reader.Read(buf[:]); err != nil {
		panic("Failed to read random bytes: " + err.Error())
	}
	return Cookie(binary.LittleEndian.Uint64(buf[:]))
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
			if cookieValue, err = strconv.ParseUint(iq.Id, 16, 64); err != nil {
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

	panic("unreachable")
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

// RequestRoster requests the user's roster from the server. It returns a
// channel on which the reply can be read when received and a Cookie that can
// be used to cancel the request.
func (c *Conn) RequestRoster() (<-chan Stanza, Cookie, error) {
	cookie := c.getCookie()
	if _, err := fmt.Fprintf(c.out, "<iq type='get' id='%x'><query xmlns='jabber:iq:roster'/></iq>", cookie); err != nil {
		return nil, 0, err
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	ch := make(chan Stanza, 1)
	c.inflights[cookie] = inflight{ch, ""}
	return ch, cookie, nil
}

type rosterEntries []RosterEntry

func (entries rosterEntries) Len() int {
	return len(entries)
}

func (entries rosterEntries) Less(i, j int) bool {
	return entries[i].Jid < entries[j].Jid
}

func (entries rosterEntries) Swap(i, j int) {
	entries[i], entries[j] = entries[j], entries[i]
}

// ParseRoster extracts roster information from the given Stanza.
func ParseRoster(reply Stanza) ([]RosterEntry, error) {
	iq, ok := reply.Value.(*ClientIQ)
	if !ok {
		return nil, errors.New("xmpp: roster request resulted in tag of type " + reply.Name.Local)
	}

	var roster Roster
	if err := xml.NewDecoder(bytes.NewBuffer(iq.Query)).Decode(&roster); err != nil {
		return nil, err
	}
	sort.Sort(rosterEntries(roster.Item))
	return roster.Item, nil
}

// SendIQ sends an info/query message to the given user. It returns a channel
// on which the reply can be read when received and a Cookie that can be used
// to cancel the request.
func (c *Conn) SendIQ(to, typ string, value interface{}) (reply chan Stanza, cookie Cookie, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	cookie = c.getCookie()
	reply = make(chan Stanza, 1)

	toAttr := ""
	if len(to) > 0 {
		toAttr = "to='" + xmlEscape(to) + "'"
	}
	if _, err = fmt.Fprintf(c.out, "<iq %s from='%s' type='%s' id='%x'>", toAttr, xmlEscape(c.jid), xmlEscape(typ), cookie); err != nil {
		return
	}
	if _, ok := value.(EmptyReply); !ok {
		if err = xml.NewEncoder(c.out).Encode(value); err != nil {
			return
		}
	}
	if _, err = fmt.Fprintf(c.out, "</iq>"); err != nil {
		return
	}

	c.inflights[cookie] = inflight{reply, to}
	return
}

// SendIQReply sends a reply to an IQ query.
func (c *Conn) SendIQReply(to, typ, id string, value interface{}) error {
	if _, err := fmt.Fprintf(c.out, "<iq to='%s' from='%s' type='%s' id='%s'>", xmlEscape(to), xmlEscape(c.jid), xmlEscape(typ), xmlEscape(id)); err != nil {
		return err
	}
	if _, ok := value.(EmptyReply); !ok {
		if err := xml.NewEncoder(c.out).Encode(value); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(c.out, "</iq>")
	return err
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

// SendPresence sends a presence stanza. If id is empty, a unique id is
// generated.
func (c *Conn) SendPresence(to, typ, id string) error {
	if len(id) == 0 {
		id = strconv.FormatUint(uint64(c.getCookie()), 10)
	}
	_, err := fmt.Fprintf(c.out, "<presence id='%s' to='%s' type='%s'/>", xmlEscape(id), xmlEscape(to), xmlEscape(typ))
	return err
}

func (c *Conn) SignalPresence(state string) error {
	_, err := fmt.Fprintf(c.out, "<presence><show>%s</show></presence>", xmlEscape(state))
	return err
}

func (c *Conn) SendStanza(s interface{}) error {
	return xml.NewEncoder(c.out).Encode(s)
}

func (c *Conn) SetCustomStorage(space, local string, s interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.customStorage == nil {
		c.customStorage = make(map[xml.Name]reflect.Type)
	}
	key := xml.Name{Space: space, Local: local}
	if s == nil {
		delete(c.customStorage, key)
	} else {
		c.customStorage[key] = reflect.TypeOf(s)
	}
}

// rfc3920 section 5.2
func (c *Conn) getFeatures(domain string) (features streamFeatures, err error) {
	if _, err = fmt.Fprintf(c.out, "<?xml version='1.0'?><stream:stream to='%s' xmlns='%s' xmlns:stream='%s' version='1.0'>\n", xmlEscape(domain), NsClient, NsStream); err != nil {
		return
	}

	se, err := nextStart(c.in)
	if err != nil {
		return
	}
	if se.Name.Space != NsStream || se.Name.Local != "stream" {
		err = errors.New("xmpp: expected <stream> but got <" + se.Name.Local + "> in " + se.Name.Space)
		return
	}

	// Now we're in the stream and can use Unmarshal.
	// Next message should be <features> to tell us authentication options.
	// See section 4.6 in RFC 3920.
	if err = c.in.DecodeElement(&features, nil); err != nil {
		err = errors.New("unmarshal <features>: " + err.Error())
		return
	}

	return
}

func (c *Conn) authenticate(features streamFeatures, user, password string) (err error) {
	havePlain := false
	for _, m := range features.Mechanisms.Mechanism {
		if m == "PLAIN" {
			havePlain = true
			break
		}
	}
	if !havePlain {
		return errors.New("xmpp: PLAIN authentication is not an option")
	}

	// Plain authentication: send base64-encoded \x00 user \x00 password.
	raw := "\x00" + user + "\x00" + password
	enc := make([]byte, base64.StdEncoding.EncodedLen(len(raw)))
	base64.StdEncoding.Encode(enc, []byte(raw))
	fmt.Fprintf(c.rawOut, "<auth xmlns='%s' mechanism='PLAIN'>%s</auth>\n", NsSASL, enc)

	// Next message should be either success or failure.
	name, val, err := next(c)
	switch v := val.(type) {
	case *saslSuccess:
	case *saslFailure:
		// v.Any is type of sub-element in failure,
		// which gives a description of what failed.
		return errors.New("xmpp: authentication failure: " + v.Any.Local)
	default:
		return errors.New("expected <success> or <failure>, got <" + name.Local + "> in " + name.Space)
	}

	return nil
}

func certName(cert *x509.Certificate) string {
	name := cert.Subject
	ret := ""

	for _, org := range name.Organization {
		ret += "O=" + org + "/"
	}
	for _, ou := range name.OrganizationalUnit {
		ret += "OU=" + ou + "/"
	}
	if len(name.CommonName) > 0 {
		ret += "CN=" + name.CommonName + "/"
	}
	return ret
}

// Resolve performs a DNS SRV lookup for the XMPP server that serves the given
// domain.
func Resolve(domain string) (host string, port uint16, err error) {
	_, addrs, err := net.LookupSRV("xmpp-client", "tcp", domain)
	if err != nil {
		return "", 0, err
	}
	if len(addrs) == 0 {
		return "", 0, errors.New("xmpp: no SRV records found for " + domain)
	}

	return addrs[0].Target, addrs[0].Port, nil
}

// Config contains options for an XMPP connection.
type Config struct {
	// Conn is the connection to the server, if non-nill.
	Conn net.Conn
	// InLog is an optional Writer which receives the raw contents of the
	// XML from the server.
	InLog io.Writer
	// OutLog is an optional Writer which receives the raw XML sent to the
	// server.
	OutLog io.Writer
	// Log is an optional Writer which receives human readable log messages
	// during the connection.
	Log io.Writer
	// CreateCallback, if not nil, causes a new account to be created on
	// the server. The callback is needed in order to be able to handle
	// XMPP forms.
	CreateCallback FormCallback
	// TrustedAddress, if true, means that the address passed to Dial is
	// trusted and that certificates for that name should be accepted.
	TrustedAddress bool
	// Archive determines whether we disable archiving for messages. If
	// false, XML is sent with each message to disable recording on the
	// server.
	Archive bool
	// ServerCertificateSHA256 contains the SHA-256 hash of the server's
	// leaf certificate, or may be empty to use normal X.509 verification.
	// If this is specified then normal X.509 verification is disabled.
	ServerCertificateSHA256 []byte
	// SkipTLS, if true, causes the TLS handshake to be skipped.
	// WARNING: this should only be used if Conn is already secure.
	SkipTLS bool
	// TLSConfig contains the configuration to be used by the TLS
	// handshake. If nil, sensible defaults will be used.
	TLSConfig *tls.Config
}

var tlsVersionStrings = map[uint16]string{
	tls.VersionSSL30: "SSL 3.0",
	tls.VersionTLS10: "TLS 1.0",
	tls.VersionTLS11: "TLS 1.1",
	tls.VersionTLS12: "TLS 1.2",
}

var tlsCipherSuiteNames = map[uint16]string{
	0x0005: "TLS_RSA_WITH_RC4_128_SHA",
	0x000a: "TLS_RSA_WITH_3DES_EDE_CBC_SHA",
	0x002f: "TLS_RSA_WITH_AES_128_CBC_SHA",
	0x0035: "TLS_RSA_WITH_AES_256_CBC_SHA",
	0xc007: "TLS_ECDHE_ECDSA_WITH_RC4_128_SHA",
	0xc009: "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA",
	0xc00a: "TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA",
	0xc011: "TLS_ECDHE_RSA_WITH_RC4_128_SHA",
	0xc012: "TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA",
	0xc013: "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA",
	0xc014: "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA",
	0xc02f: "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
	0xc02b: "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
}

func printTLSDetails(w io.Writer, tlsState tls.ConnectionState) {
	version, ok := tlsVersionStrings[tlsState.Version]
	if !ok {
		version = "unknown"
	}

	cipherSuite, ok := tlsCipherSuiteNames[tlsState.CipherSuite]
	if !ok {
		cipherSuite = "unknown"
	}

	fmt.Fprintf(w, "  SSL/TLS version: %s\n", version)
	fmt.Fprintf(w, "  Cipher suite: %s\n", cipherSuite)
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
		if err != nil {
			return nil, err
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
	if &iq.Bind == nil {
		return nil, errors.New("<iq> result missing <bind>")
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

var xmlSpecial = map[byte]string{
	'<':  "&lt;",
	'>':  "&gt;",
	'"':  "&quot;",
	'\'': "&apos;",
	'&':  "&amp;",
}

func xmlEscape(s string) string {
	var b bytes.Buffer
	for i := 0; i < len(s); i++ {
		c := s[i]
		if s, ok := xmlSpecial[c]; ok {
			b.WriteString(s)
		} else {
			b.WriteByte(c)
		}
	}
	return b.String()
}

// Scan XML token stream to find next StartElement.
func nextStart(p *xml.Decoder) (elem xml.StartElement, err error) {
	for {
		var t xml.Token
		t, err = p.Token()
		if err != nil {
			return
		}
		switch t := t.(type) {
		case xml.StartElement:
			elem = t
			return
		}
	}
	panic("unreachable")
}

// RFC 3920  C.1  Streams name space

type streamFeatures struct {
	XMLName    xml.Name `xml:"http://etherx.jabber.org/streams features"`
	StartTLS   tlsStartTLS
	Mechanisms saslMechanisms
	Bind       bindBind
	// This is a hack for now to get around the fact that the new encoding/xml
	// doesn't unmarshal to XMLName elements.
	Session *string `xml:"session"`
}

type StreamError struct {
	XMLName xml.Name `xml:"http://etherx.jabber.org/streams error"`
	Any     xml.Name `xml:",any"`
	Text    string   `xml:"text"`
}

// RFC 3920  C.3  TLS name space

type tlsStartTLS struct {
	XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls starttls"`
	Required xml.Name `xml:"required"`
}

type tlsProceed struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls proceed"`
}

type tlsFailure struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls failure"`
}

// RFC 3920  C.4  SASL name space

type saslMechanisms struct {
	XMLName   xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl mechanisms"`
	Mechanism []string `xml:"mechanism"`
}

type saslAuth struct {
	XMLName   xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl auth"`
	Mechanism string   `xml:"mechanism,attr"`
}

type saslChallenge string

type saslResponse string

type saslAbort struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl abort"`
}

type saslSuccess struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl success"`
}

type saslFailure struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl failure"`
	Any     xml.Name `xml:",any"`
}

// RFC 3920  C.5  Resource binding name space

type bindBind struct {
	XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-bind bind"`
	Resource string   `xml:"resource"`
	Jid      string   `xml:"jid"`
}

// XEP-0203: Delayed Delivery of <message/> and <presence/> stanzas.
type Delay struct {
	XMLName xml.Name `xml:"urn:xmpp:delay delay"`
	From    string   `xml:"from,attr,omitempty"`
	Stamp   string   `xml:"stamp,attr"`

	Body string `xml:",chardata"`
}

// RFC 3921  B.1  jabber:client
type ClientMessage struct {
	XMLName xml.Name `xml:"jabber:client message"`
	From    string   `xml:"from,attr"`
	Id      string   `xml:"id,attr"`
	To      string   `xml:"to,attr"`
	Type    string   `xml:"type,attr"` // chat, error, groupchat, headline, or normal

	// These should technically be []clientText,
	// but string is much more convenient.
	Subject string `xml:"subject"`
	Body    string `xml:"body"`
	Thread  string `xml:"thread"`
	Delay   *Delay `xml:"delay,omitempty"`
}

type ClientText struct {
	Lang string `xml:"lang,attr"`
	Body string `xml:",chardata"`
}

type ClientPresence struct {
	XMLName xml.Name `xml:"jabber:client presence"`
	From    string   `xml:"from,attr,omitempty"`
	Id      string   `xml:"id,attr,omitempty"`
	To      string   `xml:"to,attr,omitempty"`
	Type    string   `xml:"type,attr,omitempty"` // error, probe, subscribe, subscribed, unavailable, unsubscribe, unsubscribed
	Lang    string   `xml:"lang,attr,omitempty"`

	Show     string       `xml:"show,omitempty"`   // away, chat, dnd, xa
	Status   string       `xml:"status,omitempty"` // sb []clientText
	Priority string       `xml:"priority,omitempty"`
	Caps     *ClientCaps  `xml:"c"`
	Error    *ClientError `xml:"error"`
	Delay    Delay        `xml:"delay"`
}

type ClientCaps struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/caps c"`
	Ext     string   `xml:"ext,attr"`
	Hash    string   `xml:"hash,attr"`
	Node    string   `xml:"node,attr"`
	Ver     string   `xml:"ver,attr"`
}

type ClientIQ struct { // info/query
	XMLName xml.Name    `xml:"jabber:client iq"`
	From    string      `xml:"from,attr"`
	Id      string      `xml:"id,attr"`
	To      string      `xml:"to,attr"`
	Type    string      `xml:"type,attr"` // error, get, result, set
	Error   ClientError `xml:"error"`
	Bind    bindBind    `xml:"bind"`
	Query   []byte      `xml:",innerxml"`
}

type ClientError struct {
	XMLName xml.Name `xml:"jabber:client error"`
	Code    string   `xml:"code,attr"`
	Type    string   `xml:"type,attr"`
	Any     xml.Name `xml:",any"`
	Text    string   `xml:"text"`
}

type Roster struct {
	XMLName xml.Name      `xml:"jabber:iq:roster query"`
	Item    []RosterEntry `xml:"item"`
}

type RosterEntry struct {
	Jid          string   `xml:"jid,attr"`
	Subscription string   `xml:"subscription,attr"`
	Name         string   `xml:"name,attr"`
	Group        []string `xml:"group"`
}

type RegisterQuery struct {
	XMLName  xml.Name  `xml:"jabber:iq:register query"`
	Username *xml.Name `xml:"username"`
	Password *xml.Name `xml:"password"`
	Form     Form      `xml:"x"`
	Datas    []bobData `xml:"data"`
}

// bobData is a data element from http://xmpp.org/extensions/xep-0231.html.
type bobData struct {
	XMLName  xml.Name `xml:"urn:xmpp:bob data"`
	CID      string   `xml:"cid,attr"`
	MIMEType string   `xml:"type,attr"`
	Base64   string   `xml:",chardata"`
}

// Scan XML token stream for next element and save into val.
// If val == nil, allocate new element based on proto map.
// Either way, return val.
func next(c *Conn) (xml.Name, interface{}, error) {
	// Read start element to find out what type we want.
	se, err := nextStart(c.in)
	if err != nil {
		return xml.Name{}, nil, err
	}

	// Put it in an interface and allocate one.
	var nv interface{}
	c.lock.Lock()
	defer c.lock.Unlock()
	if t, e := c.customStorage[se.Name]; e {
		nv = reflect.New(t).Interface()
	} else if t, e := defaultStorage[se.Name]; e {
		nv = reflect.New(t).Interface()
	} else {
		return xml.Name{}, nil, errors.New("unexpected XMPP message " +
			se.Name.Space + " <" + se.Name.Local + "/>")
	}

	// Unmarshal into that storage.
	if err = c.in.DecodeElement(nv, &se); err != nil {
		return xml.Name{}, nil, err
	}
	return se.Name, nv, err
}

var defaultStorage = map[xml.Name]reflect.Type{
	xml.Name{Space: NsStream, Local: "features"}: reflect.TypeOf(streamFeatures{}),
	xml.Name{Space: NsStream, Local: "error"}:    reflect.TypeOf(StreamError{}),
	xml.Name{Space: NsTLS, Local: "starttls"}:    reflect.TypeOf(tlsStartTLS{}),
	xml.Name{Space: NsTLS, Local: "proceed"}:     reflect.TypeOf(tlsProceed{}),
	xml.Name{Space: NsTLS, Local: "failure"}:     reflect.TypeOf(tlsFailure{}),
	xml.Name{Space: NsSASL, Local: "mechanisms"}: reflect.TypeOf(saslMechanisms{}),
	xml.Name{Space: NsSASL, Local: "challenge"}:  reflect.TypeOf(""),
	xml.Name{Space: NsSASL, Local: "response"}:   reflect.TypeOf(""),
	xml.Name{Space: NsSASL, Local: "abort"}:      reflect.TypeOf(saslAbort{}),
	xml.Name{Space: NsSASL, Local: "success"}:    reflect.TypeOf(saslSuccess{}),
	xml.Name{Space: NsSASL, Local: "failure"}:    reflect.TypeOf(saslFailure{}),
	xml.Name{Space: NsBind, Local: "bind"}:       reflect.TypeOf(bindBind{}),
	xml.Name{Space: NsClient, Local: "message"}:  reflect.TypeOf(ClientMessage{}),
	xml.Name{Space: NsClient, Local: "presence"}: reflect.TypeOf(ClientPresence{}),
	xml.Name{Space: NsClient, Local: "iq"}:       reflect.TypeOf(ClientIQ{}),
	xml.Name{Space: NsClient, Local: "error"}:    reflect.TypeOf(ClientError{}),
}

type DiscoveryReply struct {
	XMLName    xml.Name            `xml:"http://jabber.org/protocol/disco#info query"`
	Node       string              `xml:"node"`
	Identities []DiscoveryIdentity `xml:"identity"`
	Features   []DiscoveryFeature  `xml:"feature"`
	Forms      []Form              `xml:"jabber:x:data x"`
}

type DiscoveryIdentity struct {
	XMLName  xml.Name `xml:"http://jabber.org/protocol/disco#info identity"`
	Lang     string   `xml:"lang,attr,omitempty"`
	Category string   `xml:"category,attr"`
	Type     string   `xml:"type,attr"`
	Name     string   `xml:"name,attr"`
}

type DiscoveryFeature struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/disco#info feature"`
	Var     string   `xml:"var,attr"`
}

type Form struct {
	XMLName      xml.Name    `xml:"jabber:x:data x"`
	Type         string      `xml:"type,attr"`
	Title        string      `xml:"title,omitempty"`
	Instructions string      `xml:"instructions,omitempty"`
	Fields       []formField `xml:"field"`
}

type formField struct {
	XMLName  xml.Name           `xml:"field"`
	Desc     string             `xml:"desc,omitempty"`
	Var      string             `xml:"var,attr"`
	Type     string             `xml:"type,attr,omitempty"`
	Label    string             `xml:"label,attr,omitempty"`
	Required *formFieldRequired `xml:"required"`
	Values   []string           `xml:"value"`
	Options  []formFieldOption  `xml:"option"`
	Media    []formFieldMedia   `xml:"media"`
}

type formFieldMedia struct {
	XMLName xml.Name   `xml:"urn:xmpp:media-element media"`
	URIs    []mediaURI `xml:"uri"`
}

type mediaURI struct {
	XMLName  xml.Name `xml:"urn:xmpp:media-element uri"`
	MIMEType string   `xml:"type,attr,omitempty"`
	URI      string   `xml:",chardata"`
}

type formFieldRequired struct {
	XMLName xml.Name `xml:"required"`
}

type formFieldOption struct {
	Label string `xml:"var,attr,omitempty"`
	Value string `xml:"value"`
}

// FormField is the type of a generic form field. One should type cast to a
// specific type of field before processing.
type FormField struct {
	// Label is a human readable label for this field.
	Label string
	// Type is the XMPP-internal type of this field. One should type cast
	// rather than inspect this.
	Type string
	// Name gives the internal name of the field.
	Name     string
	Required bool
	// Media contains one of more items of media associated with this
	// field and, for each item, one or more representations of it.
	Media [][]Media
}

type Media struct {
	MIMEType string
	// URI contains a URI to the data. It may be empty if Data is not.
	URI string
	// Data contains the raw data itself. It may be empty if URI is not.
	Data []byte
}

// FixedFormField is used to indicate a section heading. It's for the form to
// send data to the user rather than the other way around.
type FixedFormField struct {
	FormField

	Text string
}

// BooleanFormField is for a yes/no answer. The Result member should be set to
// the user's answer.
type BooleanFormField struct {
	FormField

	Result bool
}

// TextFormField is for the entry of a single textual item. The Result member
// should be set to the data entered.
type TextFormField struct {
	FormField

	Default string
	Result  string
	// Private is true if this is a password or other sensitive entry.
	Private bool
}

// MultiTextFormField is for the entry of a several textual items. The Results
// member should be set to the data entered.
type MultiTextFormField struct {
	FormField

	Defaults []string
	Results  []string
}

// SelectionFormField asks the user to pick a single element from a set of
// choices. The Result member should be set to an index of the Values array.
type SelectionFormField struct {
	FormField

	Values []string
	Ids    []string
	Result int
}

// MultiSelectionFormField asks the user to pick a subset of possible choices.
// The Result member should be set to a series of indexes of the Results array.
type MultiSelectionFormField struct {
	FormField

	Values  []string
	Ids     []string
	Results []int
}

// FormCallback is the type of a function called to process a form. The
// argument is a list of pointers to FormField types. The function should type
// cast the elements, prompt the user and fill in the result field in each
// struct.
type FormCallback func(title, instructions string, fields []interface{}) error

// processForm calls the callback with the given XMPP form and returns the
// result form. The datas argument contains any additional XEP-0231 blobs that
// might contain media for the questions in the form.
func processForm(form *Form, datas []bobData, callback FormCallback) (*Form, error) {
	var fields []interface{}

	for _, field := range form.Fields {
		base := FormField{
			Label:    field.Label,
			Type:     field.Type,
			Name:     field.Var,
			Required: field.Required != nil,
		}

		for _, media := range field.Media {
			var options []Media
			for _, uri := range media.URIs {
				media := Media{
					MIMEType: uri.MIMEType,
					URI:      uri.URI,
				}
				if strings.HasPrefix(media.URI, "cid:") {
					// cid URIs are references to data
					// blobs that, hopefully, were sent
					// along with the form.
					cid := media.URI[4:]
					media.URI = ""

					for _, data := range datas {
						if data.CID == cid {
							var err error
							if media.Data, err = base64.StdEncoding.DecodeString(data.Base64); err != nil {
								media.Data = nil
							}
						}
					}
				}
				if len(media.URI) > 0 || len(media.Data) > 0 {
					options = append(options, media)
				}
			}

			base.Media = append(base.Media, options)
		}

		switch field.Type {
		case "fixed":
			if len(field.Values) < 1 {
				continue
			}
			f := &FixedFormField{
				FormField: base,
				Text:      field.Values[0],
			}
			fields = append(fields, f)
		case "boolean":
			f := &BooleanFormField{
				FormField: base,
			}
			fields = append(fields, f)
		case "jid-multi", "text-multi":
			f := &MultiTextFormField{
				FormField: base,
				Defaults:  field.Values,
			}
			fields = append(fields, f)
		case "list-single":
			f := &SelectionFormField{
				FormField: base,
			}
			for _, opt := range field.Options {
				f.Ids = append(f.Ids, opt.Value)
				f.Values = append(f.Values, opt.Label)
			}
			fields = append(fields, f)
		case "list-multi":
			f := &MultiSelectionFormField{
				FormField: base,
			}
			for _, opt := range field.Options {
				f.Ids = append(f.Ids, opt.Value)
				f.Values = append(f.Values, opt.Label)
			}
			fields = append(fields, f)
		case "hidden":
			continue
		default:
			f := &TextFormField{
				FormField: base,
				Private:   field.Type == "text-private",
			}
			if len(field.Values) > 0 {
				f.Default = field.Values[0]
			}
			fields = append(fields, f)
		}
	}

	if err := callback(form.Title, form.Instructions, fields); err != nil {
		return nil, err
	}

	result := &Form{
		Type: "submit",
	}

	// Copy the hidden fields across.
	for _, field := range form.Fields {
		if field.Type != "hidden" {
			continue
		}
		result.Fields = append(result.Fields, formField{
			Var:    field.Var,
			Values: field.Values,
		})
	}

	for _, field := range fields {
		switch field := field.(type) {
		case *BooleanFormField:
			value := "false"
			if field.Result {
				value = "true"
			}
			result.Fields = append(result.Fields, formField{
				Var:    field.Name,
				Values: []string{value},
			})
		case *TextFormField:
			result.Fields = append(result.Fields, formField{
				Var:    field.Name,
				Values: []string{field.Result},
			})
		case *MultiTextFormField:
			result.Fields = append(result.Fields, formField{
				Var:    field.Name,
				Values: field.Results,
			})
		case *SelectionFormField:
			result.Fields = append(result.Fields, formField{
				Var:    field.Name,
				Values: []string{field.Ids[field.Result]},
			})
		case *MultiSelectionFormField:
			var values []string
			for _, selected := range field.Results {
				values = append(values, field.Ids[selected])
			}

			result.Fields = append(result.Fields, formField{
				Var:    field.Name,
				Values: values,
			})
		case *FixedFormField:
			continue
		default:
			panic(fmt.Sprintf("unknown field type in result from callback: %T", field))
		}
	}

	return result, nil
}

// VerificationString returns a SHA-1 verification string as defined in XEP-0115.
// See http://xmpp.org/extensions/xep-0115.html#ver
func (r *DiscoveryReply) VerificationString() (string, error) {
	h := sha1.New()

	seen := make(map[string]bool)
	identitySorter := &xep0115Sorter{}
	for i := range r.Identities {
		identitySorter.add(&r.Identities[i])
	}
	sort.Sort(identitySorter)
	for _, id := range identitySorter.s {
		id := id.(*DiscoveryIdentity)
		c := id.Category + "/" + id.Type + "/" + id.Lang + "/" + id.Name + "<"
		if seen[c] {
			return "", errors.New("duplicate discovery identity")
		}
		seen[c] = true
		io.WriteString(h, c)
	}

	seen = make(map[string]bool)
	featureSorter := &xep0115Sorter{}
	for i := range r.Features {
		featureSorter.add(&r.Features[i])
	}
	sort.Sort(featureSorter)
	for _, f := range featureSorter.s {
		f := f.(*DiscoveryFeature)
		if seen[f.Var] {
			return "", errors.New("duplicate discovery feature")
		}
		seen[f.Var] = true
		io.WriteString(h, f.Var+"<")
	}

	seen = make(map[string]bool)
	for _, f := range r.Forms {
		if len(f.Fields) == 0 {
			continue
		}
		fieldSorter := &xep0115Sorter{}
		for i := range f.Fields {
			fieldSorter.add(&f.Fields[i])
		}
		sort.Sort(fieldSorter)
		formTypeField := fieldSorter.s[0].(*formField)
		if formTypeField.Var != "FORM_TYPE" {
			continue
		}
		if seen[formTypeField.Type] {
			return "", errors.New("multiple forms of the same type")
		}
		seen[formTypeField.Type] = true
		if len(formTypeField.Values) != 1 {
			return "", errors.New("form does not have a single FORM_TYPE value")
		}
		if formTypeField.Type != "hidden" {
			continue
		}
		io.WriteString(h, formTypeField.Values[0]+"<")
		for _, field := range fieldSorter.s[1:] {
			field := field.(*formField)
			io.WriteString(h, field.Var+"<")
			values := append([]string{}, field.Values...)
			sort.Strings(values)
			for _, v := range values {
				io.WriteString(h, v+"<")
			}
		}
	}

	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

type xep0115Less interface {
	xep0115Less(interface{}) bool
}

type xep0115Sorter struct{ s []xep0115Less }

func (s *xep0115Sorter) add(c xep0115Less)  { s.s = append(s.s, c) }
func (s *xep0115Sorter) Len() int           { return len(s.s) }
func (s *xep0115Sorter) Swap(i, j int)      { s.s[i], s.s[j] = s.s[j], s.s[i] }
func (s *xep0115Sorter) Less(i, j int) bool { return s.s[i].xep0115Less(s.s[j]) }

func (a *DiscoveryIdentity) xep0115Less(other interface{}) bool {
	b := other.(*DiscoveryIdentity)
	if a.Category != b.Category {
		return a.Category < b.Category
	}
	if a.Type != b.Type {
		return a.Type < b.Type
	}
	return a.Lang < b.Lang
}

func (a *DiscoveryFeature) xep0115Less(other interface{}) bool {
	b := other.(*DiscoveryFeature)
	return a.Var < b.Var
}

func (a *formField) xep0115Less(other interface{}) bool {
	b := other.(*formField)
	if a.Var == "FORM_TYPE" {
		return true
	} else if b.Var == "FORM_TYPE" {
		return false
	}
	return a.Var < b.Var
}

type VersionQuery struct {
	XMLName xml.Name `xml:"jabber:iq:version query"`
}

type VersionReply struct {
	XMLName xml.Name `xml:"jabber:iq:version query"`
	Name    string   `xml:"name"`
	Version string   `xml:"version"`
	OS      string   `xml:"os"`
}

// ErrorReply reflects an XMPP error stanza. See
// http://xmpp.org/rfcs/rfc6120.html#stanzas-error-syntax
type ErrorReply struct {
	XMLName xml.Name    `xml:"error"`
	Type    string      `xml:"type,attr"`
	Error   interface{} `xml:"error"`
}

// ErrorBadRequest reflects a bad-request stanza. See
// http://xmpp.org/rfcs/rfc6120.html#stanzas-error-conditions-bad-request
type ErrorBadRequest struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-stanzas bad-request"`
}

// RosterRequest is used to request that the server update the user's roster.
// See RFC 6121, section 2.3.
type RosterRequest struct {
	XMLName xml.Name          `xml:"jabber:iq:roster query"`
	Item    RosterRequestItem `xml:"item"`
}

type RosterRequestItem struct {
	Jid          string   `xml:"jid,attr"`
	Subscription string   `xml:"subscription,attr"`
	Name         string   `xml:"name,attr"`
	Group        []string `xml:"group"`
}

// An EmptyReply results in in no XML.
type EmptyReply struct {
}
