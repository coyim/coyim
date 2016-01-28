package config

import (
	"net"
	"net/url"

	ournet "github.com/twstrike/coyim/net"
	"golang.org/x/net/proxy"
)

var (
	detectedTorAddress = ""
	scannedForTor      = false
)

// NewTorProxy creates a new proxy using the Tor service detected at the machine.
func NewTorProxy() (proxy.Dialer, error) {
	u, err := url.Parse(newTorProxy(ournet.Tor.Address()))
	if err != nil {
		return nil, err
	}

	return proxy.FromURL(u, proxy.Direct)
}

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
