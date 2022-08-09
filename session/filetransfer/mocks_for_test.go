package filetransfer

import (
	"io/ioutil"
	"net"
	"os"
	"time"

	mck "github.com/stretchr/testify/mock"
	. "gopkg.in/check.v1"
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
	returns := m.Called(network, addr)
	var ret net.Conn
	ci := returns.Get(0)
	if ci != nil {
		ret = ci.(net.Conn)
	}
	return ret, returns.Error(1)
}

type WithTempFileSuite struct {
	file    string
	content []byte
}

func (s *WithTempFileSuite) SetUpTest(c *C) {
	tf, ex := ioutil.TempFile("", "coyim-filetransfer-42-")
	c.Assert(ex, IsNil)

	s.content = []byte(`something new`)

	_, ex = tf.Write(s.content)
	c.Assert(ex, IsNil)

	ex = tf.Close()
	c.Assert(ex, IsNil)

	s.file = tf.Name()
}

func (s *WithTempFileSuite) TearDownTest(c *C) {
	e := os.Remove(s.file)
	c.Assert(e, IsNil)
}
