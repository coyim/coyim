package tls

import (
	gotls "crypto/tls"
	"net"
)

type Factory func(net.Conn, *gotls.Config) Conn

func Real(c net.Conn, conf *gotls.Config) Conn {
	return gotls.Client(c, conf)
}
