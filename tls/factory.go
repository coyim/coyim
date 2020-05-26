package tls

import (
	gotls "crypto/tls"
	"net"
)

// Factory represents a function that can create a new TLS connection
type Factory func(net.Conn, *gotls.Config) Conn

// Real is a function to get a real Golang TLS connection
func Real(c net.Conn, conf *gotls.Config) Conn {
	return gotls.Client(c, conf)
}
