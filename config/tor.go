package config

import (
	"net"
	"net/url"

	ournet "github.com/twstrike/coyim/net"
)

func newTorProxy(torAddress string) string {
	host, port, _ := net.SplitHostPort(torAddress)

	user := [10]byte{}
	pass := [10]byte{}

	var credentials *url.Userinfo
	if randomString(user[:]) == nil && randomString(pass[:]) == nil {
		credentials = url.UserPassword("randomTor-"+string(user[:]), "randomTor-"+string(pass[:]))
	}

	proxy := url.URL{
		Scheme: "socks5",
		User:   credentials,
		Host:   net.JoinHostPort(host, port),
	}

	return proxy.String()
}

// TODO[ola] figure out this one
func (a *Account) ensureTorProxy(torAddress string) {
	if !a.RequireTor {
		return
	}

	for _, proxy := range a.Proxies {
		p, err := url.Parse(proxy)
		if err != nil {
			continue
		}

		//Already configured
		if p.Host == torAddress {
			return
		}
	}

	//Tor refuses to connect to any other proxy at localhost/127.0.0.1 in the
	//chain, so we remove them
	allowedProxies := make([]string, 0, len(a.Proxies))
	for _, proxy := range a.Proxies {
		p := ournet.ParseProxy(proxy)
		if p.Host != "localhost" && p.Host != "127.0.0.1" {
			allowedProxies = append(allowedProxies, proxy)
		}
	}

	torProxy := newTorProxy(torAddress)
	allowedProxies = append(allowedProxies, torProxy)
	a.Proxies = allowedProxies
}
