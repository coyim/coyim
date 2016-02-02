package net

import (
	"net"
	"time"
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
