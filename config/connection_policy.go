package config

import (
	"crypto/tls"
	"errors"
	"io"
	"net"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/proxy"

	ournet "github.com/twstrike/coyim/net"
	"github.com/twstrike/coyim/servers"
	"github.com/twstrike/coyim/xmpp"
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

	torState ournet.TorState
}

func (p *ConnectionPolicy) ensureTorFor(conf *Account) error {
	if !conf.RequireTor {
		return nil
	}

	tor := p.torState
	if tor == nil {
		tor = ournet.Tor
	}

	if !tor.Detect() {
		return ErrTorNotRunning
	}

	conf.EnsureTorProxy(tor.Address())
	return nil
}

func (p *ConnectionPolicy) buildDialerFor(conf *Account) (*xmpp.Dialer, error) {
	//Account is a bare JID
	jidParts := strings.SplitN(conf.Account, "@", 2)
	if len(jidParts) != 2 {
		return nil, errors.New("invalid username (want user@domain): " + conf.Account)
	}

	domainpart := jidParts[1]

	if err := p.ensureTorFor(conf); err != nil {
		return nil, err
	}

	certSHA256, err := conf.ServerCertificateHash()
	if err != nil {
		return nil, err
	}

	xmppConfig := xmpp.Config{
		Archive: false,

		ServerCertificateSHA256: certSHA256,
		TLSConfig:               newTLSConfig(),

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

	dialer := &xmpp.Dialer{
		JID:    conf.Account,
		Proxy:  proxy,
		Config: xmppConfig,
	}

	// Although RFC 6120, section 3.2.3 recommends to skip the SRV lookup in this
	// case, we opt for keep compatibility with existing client implementations
	// and still make the SRV lookup. This avoids preventing imported accounts to
	// use the SRV lookup.
	if len(conf.Server) > 0 && conf.Port > 0 {
		dialer.ServerAddress = net.JoinHostPort(conf.Server, strconv.Itoa(conf.Port))
	}

	if conf.RequireTor {
		server := dialer.GetServer()
		host, port, err := net.SplitHostPort(server)
		if err != nil {
			return nil, err
		}

		if hidden, ok := servers.Get(host); ok {
			dialer.ServerAddress = net.JoinHostPort(hidden.Onion, port)
		}
	}

	return dialer, nil
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
func (p *ConnectionPolicy) Connect(password string, conf *Account) (*xmpp.Conn, error) {
	dialer, err := p.buildDialerFor(conf)
	if err != nil {
		return nil, err
	}

	// We use password rather than conf.Password because the user might have not
	// stored the password, and changing conf.Password in this case will store it.
	dialer.Password = password

	return dialer.Dial()
}

// RegisterAccount register the account on the XMPP server.
func (p *ConnectionPolicy) RegisterAccount(createCallback xmpp.FormCallback, conf *Account) (*xmpp.Conn, error) {
	dialer, err := p.buildDialerFor(conf)
	if err != nil {
		return nil, err
	}

	return dialer.RegisterAccount(createCallback)
}

func newTLSConfig() *tls.Config {
	return &tls.Config{
		MinVersion: tls.VersionTLS10,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		},
	}
}
