package net

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"golang.org/x/net/proxy"
)

var (
	defaultTorHost  = "127.0.0.1"
	defaultTorPorts = []string{"9050", "9150"}
	timeout         = 30 * time.Second
)

//TorState informs the state of Tor
type TorState interface {
	Detect() bool
	Address() string
	IsConnectionOverTor(proxy.Dialer) bool
}

// Tor is the default state manager for Tor
var Tor TorState = &defaultTorManager{}

type defaultTorManager struct {
	addr     string
	detected bool

	torHost  string
	torPorts []string
}

func (m *defaultTorManager) Detect() bool {
	torHost := m.torHost
	if len(torHost) == 0 {
		torHost = defaultTorHost
	}

	torPorts := m.torPorts
	if len(m.torPorts) == 0 {
		torPorts = defaultTorPorts
	}

	var found bool
	m.addr, found = detectTor(torHost, torPorts)
	m.detected = found
	return found
}

func (m *defaultTorManager) Address() string {
	if !m.detected {
		m.Detect()
	}

	return m.addr
}

func detectTor(host string, ports []string) (string, bool) {
	for _, port := range ports {
		addr := net.JoinHostPort(host, port)
		conn, err := net.DialTimeout("tcp", addr, timeout)
		if err != nil {
			continue
		}

		defer conn.Close()
		return addr, true
	}

	return "", false
}

// CheckTorResult represents the JSON result from a check tor request
type CheckTorResult struct {
	IsTor bool
	IP    string
}

// IsConnectionOverTor will make a connection to the check.torproject page to see if we're using Tor or not
func (*defaultTorManager) IsConnectionOverTor(d proxy.Dialer) bool {
	if d == nil {
		d = proxy.Direct
	}

	c := &http.Client{Transport: &http.Transport{Dial: func(network, addr string) (net.Conn, error) {
		return d.Dial(network, addr)
	}}}

	resp, err := c.Get("https://check.torproject.org/api/ip")
	if err != nil {
		log.WithError(err).Warn("Got error when trying to check tor")
		return false
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Warn("Got error when trying to check tor")
		return false
	}

	v := CheckTorResult{}
	err = json.Unmarshal(content, &v)
	if err != nil {
		log.WithError(err).Warn("Got error when trying to check tor")
		return false
	}

	return v.IsTor
}
