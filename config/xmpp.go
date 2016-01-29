package config

import (
	"net/url"

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
