package xmpp

import (
	"net"
	"time"
)

type timeoutableConn struct {
	net.Conn
	duration time.Duration
}

func (c *timeoutableConn) Read(b []byte) (n int, err error) {
	deadline := time.Now().Add(c.duration)
	_ = c.Conn.SetReadDeadline(deadline)
	n, err = c.Conn.Read(b)
	_ = c.Conn.SetReadDeadline(time.Time{})

	return
}

func (c *timeoutableConn) Write(b []byte) (n int, err error) {
	deadline := time.Now().Add(c.duration)
	_ = c.Conn.SetWriteDeadline(deadline)
	n, err = c.Conn.Write(b)
	_ = c.Conn.SetWriteDeadline(time.Time{})

	return
}
