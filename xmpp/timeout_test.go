package xmpp

import (
	"net"
	"time"

	. "gopkg.in/check.v1"
)

type TimeoutSuite struct{}

var _ = Suite(&TimeoutSuite{})

type callbackConn struct {
	readCB             func([]byte) (int, error)
	writeCB            func([]byte) (int, error)
	setReadDeadlineCB  func(time.Time) error
	setWriteDeadlineCB func(time.Time) error
}

func (c *callbackConn) Read(b []byte) (n int, err error) {
	return c.readCB(b)
}

func (c *callbackConn) Write(b []byte) (n int, err error) {
	return c.writeCB(b)
}

func (c *callbackConn) Close() error {
	return nil
}

func (c *callbackConn) LocalAddr() net.Addr {
	return nil
}

func (c *callbackConn) RemoteAddr() net.Addr {
	return nil
}

func (c *callbackConn) SetDeadline(t time.Time) error {
	return nil
}

func (c *callbackConn) SetReadDeadline(t time.Time) error {
	return c.setReadDeadlineCB(t)
}

func (c *callbackConn) SetWriteDeadline(t time.Time) error {
	return c.setWriteDeadlineCB(t)
}

func (s *TimeoutSuite) Test_timeoutableConn_Read_readsWithReadDeadline(c *C) {
	rc := &callbackConn{}
	calls := []string{}
	rc.readCB = func([]byte) (int, error) {
		calls = append(calls, "read")
		return 0, nil
	}

	tt := []time.Time{}
	rc.setReadDeadlineCB = func(t time.Time) error {
		calls = append(calls, "set-read-deadline")
		tt = append(tt, t)
		return nil
	}

	_, _ = (&timeoutableConn{rc, time.Duration(42) * time.Minute}).Read(nil)

	c.Assert(calls, DeepEquals, []string{"set-read-deadline", "read", "set-read-deadline"})
	c.Assert(tt[1], DeepEquals, time.Time{})
	c.Assert(tt[0].Before(time.Now()), Equals, false)
	c.Assert(tt[0].After(time.Now()), Equals, true)
	c.Assert(tt[0].Before(time.Now().Add(time.Duration(50)*time.Minute)), Equals, true)
}

func (s *TimeoutSuite) Test_timeoutableConn_Write_writesWithReadDeadline(c *C) {
	rc := &callbackConn{}
	calls := []string{}
	rc.writeCB = func([]byte) (int, error) {
		calls = append(calls, "write")
		return 0, nil
	}

	tt := []time.Time{}
	rc.setWriteDeadlineCB = func(t time.Time) error {
		calls = append(calls, "set-write-deadline")
		tt = append(tt, t)
		return nil
	}

	_, _ = (&timeoutableConn{rc, time.Duration(42) * time.Minute}).Write(nil)

	c.Assert(calls, DeepEquals, []string{"set-write-deadline", "write", "set-write-deadline"})
	c.Assert(tt[1], DeepEquals, time.Time{})
	c.Assert(tt[0].Before(time.Now()), Equals, false)
	c.Assert(tt[0].After(time.Now()), Equals, true)
	c.Assert(tt[0].Before(time.Now().Add(time.Duration(50)*time.Minute)), Equals, true)
}
