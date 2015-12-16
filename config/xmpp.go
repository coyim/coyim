package config

import (
	"errors"
	"net/url"

	"../xmpp"

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
