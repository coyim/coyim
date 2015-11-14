package config

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"

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
	RequireTor       bool
	UseHiddenService bool

	// Logger logs connection information.
	Logger io.Writer

	// XMPPLogger logs XMPP messages
	XMPPLogger io.Writer
}

func (p *ConnectionPolicy) buildDialerFor(conf *Account) (*xmpp.Dialer, error) {
	//Account is a bare JID
	jidParts := strings.SplitN(conf.Account, "@", 2)
	if len(jidParts) != 2 {
		return nil, errors.New("invalid username (want user@domain): " + conf.Account)
	}

	domainpart := jidParts[1]

	_, torDetected := DetectTor()
	if p.RequireTor && !torDetected {
		scannedForTor = false
		return nil, ErrTorNotRunning
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

	// We ignore the configured server if it is the same as the domainpart.
	// This will avoid preventing misconfigured (and imported) accounts to use
	// the SRV lookup - which is in conformance to RFC 6120.
	if len(conf.Server) > 0 && conf.Port > 0 && (conf.Server != domainpart || conf.Port != 5222) {
		dialer.ServerAddress = net.JoinHostPort(conf.Server, strconv.Itoa(conf.Port))
	}

	if p.UseHiddenService {
		server, err := dialer.GetServer()
		if err != nil {
			return nil, err
		}

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

type rawLogger struct {
	out    io.Writer
	prefix []byte
	lock   *sync.Mutex
	other  *rawLogger
	buf    []byte
}

func (r *rawLogger) Write(data []byte) (int, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if err := r.other.flush(); err != nil {
		return 0, nil
	}

	origLen := len(data)
	for len(data) > 0 {
		if newLine := bytes.IndexByte(data, '\n'); newLine >= 0 {
			r.buf = append(r.buf, data[:newLine]...)
			data = data[newLine+1:]
		} else {
			r.buf = append(r.buf, data...)
			data = nil
		}
	}

	return origLen, nil
}

func (r *rawLogger) flush() error {
	newLine := []byte{'\n'}

	if len(r.buf) == 0 {
		return nil
	}

	if _, err := r.out.Write(r.prefix); err != nil {
		return err
	}
	if _, err := r.out.Write(r.buf); err != nil {
		return err
	}
	if _, err := r.out.Write(newLine); err != nil {
		return err
	}
	r.buf = r.buf[:0]
	return nil
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
//TODO: it is weird that conf.Password is ignored and password is used
func (p *ConnectionPolicy) Connect(password string, conf *Account) (*xmpp.Conn, error) {
	dialer, err := p.buildDialerFor(conf)
	if err != nil {
		return nil, err
	}

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
