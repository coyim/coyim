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
	"log"
	"net"
	"strings"

	"golang.org/x/net/proxy"
)

// A Dialer connects and authenticates to an XMPP server
type Dialer struct {
	// JID represents the user's "bare JID" as specified in RFC 6120
	JID string

	// Password used to authenticate to the server
	Password string

	// ServerAddress associates a particular FQDN with the origin domain specified by the JID.
	ServerAddress string

	// Proxy configures a proxy used to connect to the server
	Proxy proxy.Dialer

	// Config configures the XMPP protocol
	Config Config
}

func (d *Dialer) hardcodedServer() bool {
	return d.ServerAddress != ""
}

func (d *Dialer) getHardcodedDomain() string {
	h, _, err := net.SplitHostPort(d.ServerAddress)
	if err != nil {
		//TODO: error
		return ""
	}

	return h
}

func (d *Dialer) getJIDLocalpart() string {
	parts := strings.SplitN(d.JID, "@", 2)
	return parts[0]
}

func (d *Dialer) getJIDDomainpart() string {
	//TODO: remove any existing resourcepart although our doc says it is a bare JID (without resourcepart) but it would be nice
	parts := strings.SplitN(d.JID, "@", 2)
	return parts[1]
}

func (d *Dialer) connect(addr string, conn net.Conn) (*Conn, error) {
	config := d.Config

	//JID domainpart is separated from localpart because it is used as "origin domain" for the TLS cert
	return dial(addr,
		d.getJIDLocalpart(),
		d.getJIDDomainpart(),
		d.Password,
		&config,
		conn,
	)
}

// Dial creates a new connection to an XMPP server with the given proxy
// and authenticates as the given user.
func (d *Dialer) Dial() (*Conn, error) {
	if d.Proxy == nil {
		d.Proxy = proxy.Direct
	}

	//RFC 6120, Section 3.2.3
	//See: https://xmpp.org/rfcs/rfc6120.html#tcp-resolution-srvnot
	if d.hardcodedServer() {
		addr := d.getHardcodedDomain()
		conn, err := connectWithProxy(addr, d.Proxy)
		if err != nil {
			return nil, err
		}

		return d.connect(addr, conn)
	}

	addr := d.getJIDDomainpart()
	xmppAddrs, err := ResolveProxy(d.Proxy, addr)
	if err != nil {
		return nil, err
	}

	//RFC 6120, Section 3.2.1, item 9
	//If the SRV has no response, we fallback to use
	//the domain at default port
	if len(xmppAddrs) == 0 {
		//TODO: in this case, a failure to connect might be recovered using HTTP binding
		//See: RFC 6120, Section 3.2.2
		xmppAddrs = []string{net.JoinHostPort(addr, "5222")}
	}

	conn, addr, err := connectToFirstAvailable(xmppAddrs, d.Proxy)
	if err != nil {
		return nil, err
	}

	return d.connect(addr, conn)
}

func connectToFirstAvailable(xmppAddrs []string, dialer proxy.Dialer) (net.Conn, string, error) {
	if dialer == nil {
		dialer = proxy.Direct
	}

	for _, addr := range xmppAddrs {
		conn, err := connectWithProxy(addr, dialer)
		if err == nil {
			return conn, addr, nil
		}
	}

	return nil, "", errors.New("Failed to connect to XMPP server: exhausted list of XMPP SRV for server")
}

func connectWithProxy(addr string, dialer proxy.Dialer) (conn net.Conn, err error) {
	log.Printf("Connecting to %s\n", addr)

	//TODO: It is not clear to me if this follows
	//RFC 6120, Section 3.2.1, item 6
	//See: https://xmpp.org/rfcs/rfc6120.html#tcp-resolution
	conn, err = dialer.Dial("tcp", addr)
	if err != nil {
		log.Printf("Failed to connect to %s: %s\n", addr, err)
		return
	}

	return
}

func dial(address, user, domain, password string, config *Config, conn net.Conn) (c *Conn, err error) {
	c = new(Conn)
	c.config = config
	c.inflights = make(map[Cookie]inflight)
	c.archive = config.Archive

	c.in, c.out = makeInOut(conn, config)
	c.rawOut = conn

	features, err := c.getFeatures(domain)
	if err != nil {
		return nil, err
	}

	if !config.SkipTLS {
		if features.StartTLS.XMLName.Local == "" {
			return nil, errors.New("xmpp: server doesn't support TLS")
		}

		if err := c.startTLS(address, domain, conn); err != nil {
			return nil, err
		}

		if features, err = c.getFeatures(domain); err != nil {
			return nil, err
		}
	}

	if err := createAccount(user, password, config, c); err != nil {
		return nil, err
	}

	if err := authenticate(features, user, password, config, c); err != nil {
		return nil, err
	}

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

func (c *Conn) startTLS(address, domain string, conn net.Conn) error {
	fmt.Fprintf(c.out, "<starttls xmlns='%s'/>", NsTLS)

	proceed, err := nextStart(c.in)
	if err != nil {
		return err
	}

	if proceed.Name.Space != NsTLS || proceed.Name.Local != "proceed" {
		return errors.New("xmpp: expected <proceed> after <starttls> but got <" + proceed.Name.Local + "> in " + proceed.Name.Space)
	}

	l := logFor(c.config)
	io.WriteString(l, "Starting TLS handshake\n")

	var tlsConfig tls.Config
	if c.config.TLSConfig != nil {
		tlsConfig = *c.config.TLSConfig
	}
	tlsConfig.ServerName = domain
	tlsConfig.InsecureSkipVerify = true

	tlsConn := tls.Client(conn, &tlsConfig)
	if err := tlsConn.Handshake(); err != nil {
		return err
	}

	tlsState := tlsConn.ConnectionState()
	printTLSDetails(l, tlsState)

	haveCertHash := len(c.config.ServerCertificateSHA256) != 0
	if haveCertHash {
		h := sha256.New()
		h.Write(tlsState.PeerCertificates[0].Raw)
		if digest := h.Sum(nil); !bytes.Equal(digest, c.config.ServerCertificateSHA256) {
			return fmt.Errorf("xmpp: server certificate does not match expected hash (got: %x, want: %x)",
				digest, c.config.ServerCertificateSHA256)
		}
	} else {
		if len(tlsState.PeerCertificates) == 0 {
			return errors.New("xmpp: server has no certificates")
		}

		opts := x509.VerifyOptions{
			Intermediates: x509.NewCertPool(),
		}
		for _, cert := range tlsState.PeerCertificates[1:] {
			opts.Intermediates.AddCert(cert)
		}
		verifiedChains, err := tlsState.PeerCertificates[0].Verify(opts)
		if err != nil {
			return errors.New("xmpp: failed to verify TLS certificate: " + err.Error())
		}

		for i, cert := range verifiedChains[0] {
			fmt.Fprintf(l, "  certificate %d: %s\n", i, certName(cert))
		}
		leafCert := verifiedChains[0][0]

		if err := leafCert.VerifyHostname(domain); err != nil {
			if c.config.TrustedAddress {
				fmt.Fprintf(l, "Certificate fails to verify against domain in username: %s\n", err)
				host, _, err := net.SplitHostPort(address)
				if err != nil {
					return errors.New("xmpp: failed to split address when checking whether TLS certificate is valid: " + err.Error())
				}

				if err = leafCert.VerifyHostname(host); err != nil {
					return errors.New("xmpp: failed to match TLS certificate to address after failing to match to username: " + err.Error())
				}
				fmt.Fprintf(l, "Certificate matches against trusted server hostname: %s\n", host)
			} else {
				return errors.New("xmpp: failed to match TLS certificate to name: " + err.Error())
			}
		}
	}

	c.in, c.out = makeInOut(tlsConn, c.config)
	c.rawOut = tlsConn

	return nil
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

func logFor(config *Config) io.Writer {
	log := ioutil.Discard
	if config != nil && config.Log != nil {
		log = config.Log
	}

	return log
}

func authenticate(features streamFeatures, user, password string, config *Config, c *Conn) error {
	l := logFor(config)
	io.WriteString(l, "Authenticating as "+user+"\n")
	if err := c.authenticate(features, user, password); err != nil {
		return err
	}

	io.WriteString(l, "Authentication successful\n")
	return nil
}

func createAccount(user, password string, config *Config, c *Conn) error {
	if config == nil || config.CreateCallback == nil {
		return nil
	}

	io.WriteString(logFor(config), "Attempting to create account\n")
	fmt.Fprintf(c.out, "<iq type='get' id='create_1'><query xmlns='jabber:iq:register'/></iq>")
	var iq ClientIQ
	if err := c.in.DecodeElement(&iq, nil); err != nil {
		return errors.New("unmarshal <iq>: " + err.Error())
	}

	if iq.Type != "result" {
		return errors.New("xmpp: account creation failed")
	}
	var register RegisterQuery
	if err := xml.NewDecoder(bytes.NewBuffer(iq.Query)).Decode(&register); err != nil {
		return err
	}

	if len(register.Form.Type) > 0 {
		reply, err := processForm(&register.Form, register.Datas, config.CreateCallback)
		fmt.Fprintf(c.rawOut, "<iq type='set' id='create_2'><query xmlns='jabber:iq:register'>")
		if err = xml.NewEncoder(c.rawOut).Encode(reply); err != nil {
			return err
		}
		fmt.Fprintf(c.rawOut, "</query></iq>")
	} else if register.Username != nil && register.Password != nil {
		// Try the old-style registration.
		fmt.Fprintf(c.rawOut, "<iq type='set' id='create_2'><query xmlns='jabber:iq:register'><username>%s</username><password>%s</password></query></iq>", user, password)
	}

	var iq2 ClientIQ
	if err := c.in.DecodeElement(&iq2, nil); err != nil {
		return errors.New("unmarshal <iq>: " + err.Error())
	}

	if iq2.Type == "error" {
		return errors.New("xmpp: account creation failed")
	}

	return nil
}
