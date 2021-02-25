package config

import (
	"errors"
	"io"
	"net"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/proxy"

	"github.com/coyim/coyim/coylog"
	ournet "github.com/coyim/coyim/net"
	"github.com/coyim/coyim/servers"
	ourtls "github.com/coyim/coyim/tls"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/interfaces"
)

var (
	// ErrTorNotRunning is the error returned when Tor is required by the policy
	// but it was not found to be running (on port 9050 or 9051).
	ErrTorNotRunning = errors.New("Tor is not running")
)

// ConnectionPolicy represents a policy to connect to XMPP servers
type ConnectionPolicy struct {
	// Logger logs connection information.
	Logger io.Writer

	// XMPPLogger logs XMPP messages
	XMPPLogger io.Writer

	DialerFactory interfaces.DialerFactory

	torState ournet.TorState

	Log coylog.Logger
}

func (p *ConnectionPolicy) initTorState() {
	if p.torState == nil {
		p.torState = ournet.Tor
	}
}

func (p *ConnectionPolicy) isTorRunning() error {
	p.initTorState()

	if !p.torState.Detect() {
		return ErrTorNotRunning
	}

	return nil
}

// HasTorAuto check if account has proxy with prefix "tor-auto://"
func (a *Account) HasTorAuto() bool {
	for _, px := range a.Proxies {
		if strings.HasPrefix(px, "tor-auto://") {
			return true
		}
	}
	return false
}

var buildDialerFor = func(p *ConnectionPolicy, conf *Account, verifier ourtls.Verifier) (interfaces.Dialer, error) {
	return p.buildDialerFor(conf, verifier)
}

func (p *ConnectionPolicy) buildDialerFor(conf *Account, verifier ourtls.Verifier) (interfaces.Dialer, error) {
	//Account is a bare JID
	jidParts := strings.SplitN(conf.Account, "@", 2)
	if len(jidParts) != 2 {
		return nil, errors.New("invalid username (want user@domain): " + conf.Account)
	}

	domainpart := jidParts[1]

	p.initTorState()

	hasTorAuto := conf.HasTorAuto()

	if hasTorAuto {
		if err := p.isTorRunning(); err != nil {
			return nil, err
		}
	}

	xmppConfig := data.Config{
		Archive: false,

		TLSConfig: newTLSConfig(),

		Log: p.Logger,
	}

	xmppConfig.InLog, xmppConfig.OutLog = buildInOutLogs(p.XMPPLogger)

	domainRoot, err := rootCAFor(domainpart)
	if err != nil {
		//alert(term, "Tried to add CACert root for jabber.ccc.de but failed: "+err.Error())
		return nil, err
	}

	if domainRoot != nil {
		//alert(term, "Temporarily trusting only CACert root for CCC Jabber server")
		xmppConfig.TLSConfig.RootCAs = domainRoot
	}

	proxy, err := buildProxyChain(conf.Proxies)
	if err != nil {
		return nil, err
	}

	dialer := p.DialerFactory(verifier, ourtls.Real)
	dialer.SetLogger(p.Log)
	dialer.SetJID(conf.Account)
	dialer.SetProxy(proxy)
	dialer.SetConfig(xmppConfig)

	// Although RFC 6120, section 3.2.3 recommends to skip the SRV lookup in this
	// case, we opt for keep compatibility with existing client implementations
	// and still make the SRV lookup. This avoids preventing imported accounts to
	// use the SRV lookup.
	if len(conf.Server) > 0 && conf.Port > 0 {
		dialer.SetServerAddress(net.JoinHostPort(conf.Server, strconv.Itoa(conf.Port)))
	}

	server := dialer.GetServer()
	// This error can't logically happen, so we will ignore it here
	host, port, _ := net.SplitHostPort(server)

	known, ok := servers.Get(host)
	if ok {
		dialer.SetKnown(&known)
	}

	if hasTorAuto || p.torState.IsConnectionOverTor(proxy) {
		if hidden, ok := servers.GetOnion(host); ok {
			dialer.SetServerAddress(net.JoinHostPort(hidden, port))
		}
	}

	return dialer, nil
}

var torDetect = func() bool {
	return ournet.Tor.Detect()
}

// CreateTorProxy returns a dialer that uses the Tor connection if available
func (a *Account) CreateTorProxy() (proxy.Dialer, error) {
	if a.HasTorAuto() {
		if !torDetect() {
			return nil, ErrTorNotRunning
		}
	}

	return buildProxyChain(a.Proxies)
}

func buildProxyChain(proxies []string) (dialer proxy.Dialer, err error) {
	for i := len(proxies) - 1; i >= 0; i-- {
		u, e := url.Parse(proxies[i])
		if e != nil {
			err = errors.New("Failed to parse " + proxies[i] + " as a URL: " + e.Error())
			return
		}

		if dialer == nil {
			dialer = &net.Dialer{
				Timeout: 60 * time.Second,
			}
		}

		if dialer, err = proxy.FromURL(u, dialer); err != nil {
			err = errors.New("Failed to parse " + proxies[i] + " as a proxy: " + err.Error())
			return
		}
	}

	return
}

func buildInOutLogs(rawLog io.Writer) (io.Writer, io.Writer) {
	if rawLog == nil {
		return nil, nil
	}

	lock := new(sync.Mutex)
	in := rawLogger{
		out:    rawLog,
		prefix: []byte("<- "),
		lock:   lock,
	}
	out := rawLogger{
		out:    rawLog,
		prefix: []byte("-> "),
		lock:   lock,
	}
	in.other, out.other = &out, &in

	go in.flush()
	go out.flush()

	return &in, &out
}

// Connect to the server and authenticates with the password
func (p *ConnectionPolicy) Connect(password, resource string, conf *Account, verifier ourtls.Verifier) (interfaces.Conn, error) {
	dialer, err := buildDialerFor(p, conf, verifier)
	if err != nil {
		return nil, err
	}

	// We use password rather than conf.Password because the user might have not
	// stored the password, and changing conf.Password in this case will store it.
	dialer.SetPassword(password)
	dialer.SetResource(resource)
	dialer.SetShouldConnectTLS(conf.ConnectTLS)
	dialer.SetShouldSendALPN(conf.SetALPN)

	return dialer.Dial()
}

// RegisterAccount register the account on the XMPP server.
func (p *ConnectionPolicy) RegisterAccount(createCallback data.FormCallback, conf *Account, verifier ourtls.Verifier) (interfaces.Conn, error) {
	dialer, err := buildDialerFor(p, conf, verifier)
	if err != nil {
		return nil, err
	}

	conn, err := dialer.RegisterAccount(createCallback)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
