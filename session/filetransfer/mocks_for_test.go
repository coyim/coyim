package filetransfer

import (
	"net"
	"time"

	mck "github.com/stretchr/testify/mock"
)

type mockedConn struct {
	mck.Mock
}

func (c *mockedConn) Read(b []byte) (n int, err error) {
	args := c.Called(b)
	return args.Int(0), args.Error(1)
}

func (c *mockedConn) Write(b []byte) (n int, err error) {
	args := c.Called(b)
	return args.Int(0), args.Error(1)
}

func (c *mockedConn) Close() error {
	return c.Called().Error(0)
}

func (c *mockedConn) LocalAddr() net.Addr {
	return c.Called().Get(0).(net.Addr)
}

func (c *mockedConn) RemoteAddr() net.Addr {
	return c.Called().Get(0).(net.Addr)
}

func (c *mockedConn) SetDeadline(t time.Time) error {
	return c.Called(t).Error(0)
}

func (c *mockedConn) SetReadDeadline(t time.Time) error {
	return c.Called(t).Error(0)
}

func (c *mockedConn) SetWriteDeadline(t time.Time) error {
	return c.Called(t).Error(0)
}

type mockedDialer struct {
	mck.Mock
}

func (m *mockedDialer) Dial(network, addr string) (net.Conn, error) {
	args := m.Called(network, addr)
	var ret net.Conn
	ci := args.Get(0)
	if ci != nil {
		ret = ci.(net.Conn)
	}
	return ret, args.Error(1)
}
