package xmpp

import (
	"encoding/xml"
	"io"
	"net"
	"strings"

	"github.com/twstrike/coyim/xmpp/data"
	"github.com/twstrike/coyim/xmpp/interfaces"

	"golang.org/x/net/proxy"
)

// A dialer connects and authenticates to an XMPP server
type dialer struct {
	// JID represents the user's "bare JID" as specified in RFC 6120
	JID string

	// password used to authenticate to the server
	password string

	// serverAddress associates a particular FQDN with the origin domain specified by the JID.
	serverAddress string

	// proxy configures a proxy used to connect to the server
	proxy proxy.Dialer

	// config configures the XMPP protocol
	config data.Config
}

// DialerFactory returns a new xmpp dialer
func DialerFactory() interfaces.Dialer {
	return &dialer{}
}

func (d *dialer) SetJID(v string) {
	d.JID = v
}

func (d *dialer) SetServerAddress(v string) {
	d.serverAddress = v
}

func (d *dialer) SetPassword(v string) {
	d.password = v
}

func (d *dialer) SetProxy(v proxy.Dialer) {
	d.proxy = v
}

func (d *dialer) SetConfig(v data.Config) {
	d.config = v
}

func (d *dialer) Config() data.Config {
	return d.config
}

func (d *dialer) ServerAddress() string {
	return d.serverAddress
}

func (d *dialer) hasCustomServer() bool {
	return d.serverAddress != ""
}

func (d *dialer) getJIDLocalpart() string {
	parts := strings.SplitN(d.JID, "@", 2)
	return parts[0]
}

func (d *dialer) getJIDDomainpart() string {
	//TODO: remove any existing resourcepart although our doc says it is a bare JID (without resourcepart) but it would be nice
	parts := strings.SplitN(d.JID, "@", 2)
	return parts[1]
}

// GetServer returns the "hardcoded" server chosen if available, otherwise returns the domainpart from the JID. The server contains port information
func (d *dialer) GetServer() string {
	if d.hasCustomServer() {
		return d.serverAddress
	}

	return d.getFallbackServer()
}

func (d *dialer) getFallbackServer() string {
	return net.JoinHostPort(d.getJIDDomainpart(), "5222")
}

// RegisterAccount registers an account on the server. The formCallback is used to handle XMPP forms.
func (d *dialer) RegisterAccount(formCallback data.FormCallback) (interfaces.Conn, error) {
	//TODO: notify in case the feature is not supported
	d.config.CreateCallback = formCallback
	return d.Dial()
}

// Dial creates a new connection to an XMPP server with the given proxy
// and authenticates as the given user.
func (d *dialer) Dial() (interfaces.Conn, error) {
	// Starting an XMPP connection comprises two parts:
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
func (d *dialer) setupStream(conn net.Conn) (interfaces.Conn, error) {
	if d.hasCustomServer() {
		d.config.TrustedAddress = true
	}

	c := newConn()
	c.config = d.config
	c.originDomain = d.getJIDDomainpart()
	d.bindTransport(c, conn)

	if err := d.negotiateStreamFeatures(c, conn); err != nil {
		return nil, err
	}

	go c.watchKeepAlive(conn)
	go c.watchPings()

	return c, nil
}

// RFC 6120, section 4.3.2
func (d *dialer) negotiateStreamFeatures(c interfaces.Conn, conn net.Conn) error {
	if err := c.SendInitialStreamHeader(); err != nil {
		return err
	}

	// STARTTLS MUST be the first feature to be negotiated
	if err := d.negotiateSTARTTLS(c, conn); err != nil {
		return err
	}

	if registered, err := d.negotiateInBandRegistration(c); err != nil || registered {
		return err
	}

	// SASL negotiation. RFC 6120, section 6
	if err := d.negotiateSASL(c); err != nil {
		return err
	}

	//TODO: negotiate other features

	return nil
}

func (d *dialer) bindTransport(c interfaces.Conn, conn net.Conn) {
	c.SetInOut(makeInOut(conn, d.config))
	c.SetRawOut(conn)
	c.SetKeepaliveOut(&timeoutableConn{conn, keepaliveTimeout})
}

func makeInOut(conn io.ReadWriter, config data.Config) (in *xml.Decoder, out io.Writer) {
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
