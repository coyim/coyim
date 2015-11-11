package config

import (
	"crypto/tls"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/twstrike/coyim/xmpp"

	"golang.org/x/net/proxy"
)

func init() {
	proxy.RegisterDialerType("socks5+unix", func(u *url.URL, d proxy.Dialer) (proxy.Dialer, error) {
		var auth *proxy.Auth
		if u.User != nil {
			auth = &proxy.Auth{
				User: u.User.Username(),
			}

			if p, ok := u.User.Password(); ok {
				auth.Password = p
			}
		}

		return proxy.SOCKS5("unix", u.Path, auth, d)
	})
}

// ResolveXMPPServerOverTor resolves the XMPP service from a domain using Tor
//TODO: remove me once config assistant goes away
func ResolveXMPPServerOverTor(domain string) ([]string, error) {
	dnsProxy, err := NewTorProxy()
	if err != nil {
		return nil, errors.New("Failed to resolve XMPP server: " + err.Error())
	}

	ret, err := xmpp.ResolveProxy(dnsProxy, domain)
	if err != nil {
		return nil, errors.New("Failed to resolve XMPP server: " + err.Error())
	}

	return ret, nil
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
				Timeout: 30 * time.Second,
			}
		}

		if dialer, err = proxy.FromURL(u, dialer); err != nil {
			err = errors.New("Failed to parse " + proxies[i] + " as a proxy: " + err.Error())
			return
		}
	}

	return
}

// NewXMPPConn creates a new XMPP connection based on the given information
func NewXMPPConn(conf *Account, password string, createCallback xmpp.FormCallback, logger io.Writer) (*xmpp.Conn, error) {
	parts := strings.SplitN(conf.Account, "@", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid username (want user@domain): " + conf.Account)
	}

	domain := parts[1]
	addrTrusted := false

	if len(conf.Server) > 0 && conf.Port > 0 {
		addrTrusted = true
	} else {
		if len(conf.Proxies) > 0 && len(detectTor()) == 0 {
			return nil, errors.New("Cannot connect via a proxy without Server and Port being set in the config file as an SRV lookup would leak information.")
		}
	}

	var certSHA256 []byte
	var err error
	if len(conf.ServerCertificateSHA256) > 0 {
		certSHA256, err = hex.DecodeString(conf.ServerCertificateSHA256)
		if err != nil {
			return nil, errors.New("Failed to parse ServerCertificateSHA256 (should be hex string): " + err.Error())
		}

		if len(certSHA256) != 32 {
			return nil, errors.New("ServerCertificateSHA256 is not 32 bytes long")
		}
	}

	xmppConfig := xmpp.Config{
		Log:                     logger,
		CreateCallback:          createCallback,
		TrustedAddress:          addrTrusted,
		Archive:                 false,
		ServerCertificateSHA256: certSHA256,
		TLSConfig:               newTLSConfig(),
	}

	domainRoot, err := rootCAFor(domain)
	if err != nil {
		//alert(term, "Tried to add CACert root for jabber.ccc.de but failed: "+err.Error())
	}

	if domainRoot != nil {
		//alert(term, "Temporarily trusting only CACert root for CCC Jabber server")
		xmppConfig.TLSConfig.RootCAs = domainRoot
	}

	//TODO: uncomment me
	//Also, move this defered functions
	//if len(conf.RawLogFile) > 0 {
	//	rawLog, err := os.OpenFile(conf.RawLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	//	if err != nil {
	//		return nil, errors.New("Failed to open raw log file: " + err.Error())
	//	}

	//	lock := new(sync.Mutex)
	//	in := rawLogger{
	//		out:    rawLog,
	//		prefix: []byte("<- "),
	//		lock:   lock,
	//	}
	//	out := rawLogger{
	//		out:    rawLog,
	//		prefix: []byte("-> "),
	//		lock:   lock,
	//	}
	//	in.other, out.other = &out, &in

	//	xmppConfig.InLog = &in
	//	xmppConfig.OutLog = &out

	//	defer in.flush()
	//	defer out.flush()
	//}

	proxy, err := buildProxyChain(conf.Proxies)
	if err != nil {
		return nil, err
	}

	dialer := xmpp.Dialer{
		JID:      conf.Account,
		Password: password,
		Proxy:    proxy,
		Config:   xmppConfig,
	}

	if len(conf.Server) > 0 && conf.Port > 0 {
		dialer.ServerAddress = fmt.Sprintf("%s:%d", conf.Server, conf.Port)
	}

	return dialer.Dial()
}

func newTLSConfig() *tls.Config {
	return &tls.Config{
		MinVersion: tls.VersionTLS10,
		CipherSuites: []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		},
	}
}
