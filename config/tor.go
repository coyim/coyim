package config

import (
	"net"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

var (
	torHost            = "127.0.0.1"
	torPorts           = []string{"9050", "9150"}
	detectedTorAddress = ""
	scannedForTor      = false
)

// TODO: set up Tor config correctly here, instead of from session

// DetectTor detects a Tor service running in the machine.
func DetectTor() (string, bool) {
	detectedTorAddress = ""
	detectTor()

	return detectedTorAddress, len(detectedTorAddress) != 0
}

// NewTorProxy creates a new proxy using the Tor service detected at the machine.
func NewTorProxy() (proxy.Dialer, error) {
	u, err := url.Parse(newTorProxy(detectTor()))
	if err != nil {
		return nil, err
	}

	return proxy.FromURL(u, proxy.Direct)
}

func detectTor() string {
	if scannedForTor {
		return detectedTorAddress
	}

	detectedTorAddress = ""
	for _, port := range torPorts {
		addr := net.JoinHostPort(torHost, port)
		conn, err := net.DialTimeout("tcp", addr, 30*time.Second)
		if err != nil {
			continue
		}

		detectedTorAddress = addr
		conn.Close()
	}

	scannedForTor = true
	return detectedTorAddress
}

func newTorProxy(torAddress string) string {
	host, port, _ := net.SplitHostPort(torAddress)

	user := [10]byte{}
	pass := [10]byte{}

	var credentials *url.Userinfo
	if randomString(user[:]) == nil && randomString(pass[:]) == nil {
		credentials = url.UserPassword(string(user[:]), string(pass[:]))
	}

	proxy := url.URL{
		Scheme: "socks5",
		User:   credentials,
		Host:   net.JoinHostPort(host, port),
	}

	return proxy.String()
}
