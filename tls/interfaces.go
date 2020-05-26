package tls

import (
	gotls "crypto/tls"
	"net"
)

// Conn represents a TLS interface
type Conn interface {
	net.Conn

	Handshake() error
	ConnectionState() gotls.ConnectionState
}
