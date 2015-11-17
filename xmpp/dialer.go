package xmpp

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
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

func (d *Dialer) getJIDLocalpart() string {
	parts := strings.SplitN(d.JID, "@", 2)
	return parts[0]
}

func (d *Dialer) getJIDDomainpart() string {
	//TODO: remove any existing resourcepart although our doc says it is a bare JID (without resourcepart) but it would be nice
	parts := strings.SplitN(d.JID, "@", 2)
	return parts[1]
}

// GetServer returns the "hardcoded" server chosen if available, otherwise returns the domainpart from the JID. The server contains port information
func (d *Dialer) GetServer() string {
	if d.hardcodedServer() {
		return d.ServerAddress
	}

	return net.JoinHostPort(d.getJIDDomainpart(), "5222")
}

// RegisterAccount registers an account on the server. The formCallback is used to handle XMPP forms.
func (d *Dialer) RegisterAccount(formCallback FormCallback) (*Conn, error) {
	d.Config.CreateCallback = formCallback
	return d.Dial()
}

// Dial creates a new connection to an XMPP server with the given proxy
// and authenticates as the given user.
func (d *Dialer) Dial() (*Conn, error) {
	// Starting an XMPP connectin comprises two parts:
	// - Opening a transport channel (TCP)
	// - Opening an XML stream over the transport channel

	// RFC 6120, section 3
	conn, err := d.newTCPConn()
	if err != nil {
		return nil, err
	}

	// RFC 6120, section 4
	return d.setupStream(conn)
}

// RFC 6120, Section 4.2
func (d *Dialer) setupStream(conn net.Conn) (c *Conn, err error) {
	if d.hardcodedServer() {
		d.Config.TrustedAddress = true
	}

	//JID domainpart is separated from localpart because it is used as "origin domain" for the TLS cert
	return setupStream(d.GetServer(), d.getJIDLocalpart(), d.getJIDDomainpart(), d.Password, d.Config, conn)
}

func setupStream(address, user, domain, password string, config Config, conn net.Conn) (c *Conn, err error) {
	c = new(Conn)
	c.config = config
	c.inflights = make(map[Cookie]inflight)
	c.archive = config.Archive

	c.in, c.out = makeInOut(conn, config)
	c.rawOut = conn

	features, err := c.negotiateStream(address, domain, conn)
	if err != nil {
		return nil, err
	}

	if features.InBandRegistration != nil {
		if err := createAccount(user, password, config, c); err != nil {
			return nil, err
		}
	}

	if err := c.authenticate(features, user, password); err != nil {
		return nil, ErrAuthenticationFailed
	}

	if features, err = c.sendInitialStreamHeader(domain); err != nil {
		return nil, err
	}

	// Send IQ message asking to bind to the local user name.
	// RFC 6210 section 7 states this is mandatory, so a missing features.Bind
	// is a protocol failure
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

//rfc3920 section 5.2
//TODO RFC 6120 obsoletes RFC 3920
func (c *Conn) negotiateStream(address, domain string, conn net.Conn) (features streamFeatures, err error) {
	features, err = c.sendInitialStreamHeader(domain)
	if err != nil {
		return
	}

	if !c.config.SkipTLS {
		if features.StartTLS.XMLName.Local == "" {
			err = errors.New("xmpp: server doesn't support TLS")
			return
		}

		if err = c.startTLS(address, domain, conn); err != nil {
			return
		}

		features, err = c.sendInitialStreamHeader(domain)
	}

	return
}

func makeInOut(conn io.ReadWriter, config Config) (in *xml.Decoder, out io.Writer) {
	if config.InLog != nil {
		in = xml.NewDecoder(io.TeeReader(conn, config.InLog))
	} else {
		in = xml.NewDecoder(conn)
	}

	if config.OutLog != nil {
		out = io.MultiWriter(conn, config.OutLog)
	} else {
		out = conn
	}

	return
}
