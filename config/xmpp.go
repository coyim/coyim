package config

import (
	"errors"
	"net"
	"net/url"
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

	ret, err := xmpp.ResolveSRVWithProxy(dnsProxy, domain)
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
