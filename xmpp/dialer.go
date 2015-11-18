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

	c = NewConn(nil, nil, "")
	c.config = d.Config
	d.bindTransport(c, conn)

	originDomain := d.getJIDDomainpart()

	features, err := d.negotiateStreamFeatures(c, conn)
	if err != nil {
		return nil, err
	}

	go c.watchKeepAlive(conn)

	// Resource binding. RFC 6120, section 7
	// This is mandatory, so a missing features.Bind is a protocol failure
	fmt.Fprintf(c.out, "<iq type='set' id='bind_1'><bind xmlns='%s'/></iq>", NsBind)
	var iq ClientIQ
	if err = c.in.DecodeElement(&iq, nil); err != nil {
		return nil, errors.New("unmarshal <iq>: " + err.Error())
	}
	c.jid = iq.Bind.Jid // our local id

	if features.Session != nil {
		// The server needs a session to be established. See RFC 3921,
		// section 3.
		fmt.Fprintf(c.out, "<iq to='%s' type='set' id='sess_1'><session xmlns='%s'/></iq>", originDomain, NsSession)
		if err = c.in.DecodeElement(&iq, nil); err != nil {
			return nil, errors.New("xmpp: unmarshal <iq>: " + err.Error())
		}
		if iq.Type != "result" {
			return nil, errors.New("xmpp: session establishment failed")
		}
	}

	return c, nil
}

// RFC 6120, section 4.3.2
func (d *Dialer) negotiateStreamFeatures(c *Conn, conn net.Conn) (features streamFeatures, err error) {
	originDomain := d.getJIDDomainpart()

	features, err = c.sendInitialStreamHeader(originDomain)
	if err != nil {
		return
	}

	// STARTTLS is the first feature to be negotiated
	features, err = d.negotiateSTARTTLS(features, c, conn)
	if err != nil {
		return
	}

	err = d.negotiateInBandRegistration(features, c)
	if err != nil {
		return
	}

	// SASL negotiation. RFC 6120, section 6
	features, err = d.negotiateSASL(features, c)
	if err != nil {
		return
	}

	//TODO: negotiate other features

	return
}

func (d *Dialer) bindTransport(c *Conn, conn net.Conn) {
	c.in, c.out = makeInOut(conn, d.Config)
	c.rawOut = conn
	c.keepaliveOut = &timeoutableConn{conn, keepaliveTimeout}
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
