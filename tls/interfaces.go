package tls

import (
	gotls "crypto/tls"
	"net"
)

type Conn interface {
	net.Conn

	Handshake() error
	ConnectionState() gotls.ConnectionState
}
