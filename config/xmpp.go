package config

import (
	"net/url"

	ournet "github.com/twstrike/coyim/net"
	"golang.org/x/net/proxy"
)

func socks5UnixProxy(u *url.URL, d proxy.Dialer) (proxy.Dialer, error) {
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
}

func genTorAutoString() string {
	s := [10]byte{}
	randomString(s[:])
	return "randomTorAuto-" + string(s[:])
}

func genTorAutoUsername() string {
	return genTorAutoString()
}

func genTorAutoPassword() string {
	return genTorAutoString()
}

func genTorAutoAuth(u *url.URL) *proxy.Auth {
	auth := &proxy.Auth{}
	if u.User != nil {
		auth.User = u.User.Username()
		if p, ok := u.User.Password(); ok {
			auth.Password = p
		}
	} else {
		auth.User = genTorAutoUsername()
		auth.Password = genTorAutoPassword()
	}
	return auth
}

func genTorAutoAddr(u *url.URL) string {
	if u.Host == "" {
		return ournet.Tor.Address()
	}

	return u.Host
}

func torAutoProxy(u *url.URL, d proxy.Dialer) (proxy.Dialer, error) {
	auth := genTorAutoAuth(u)
	addr := genTorAutoAddr(u)
	if addr == "" {
		return nil, ErrTorNotRunning
	}
	return proxy.SOCKS5("tcp", addr, auth, d)
}

func init() {
	proxy.RegisterDialerType("socks5+unix", socks5UnixProxy)
	proxy.RegisterDialerType("tor-auto", torAutoProxy)
}
