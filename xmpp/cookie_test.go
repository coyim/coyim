package xmpp

import (
	"errors"

	. "gopkg.in/check.v1"
)

type CookieXMPPSuite struct{}

var _ = Suite(&CookieXMPPSuite{})

func (s *CookieXMPPSuite) Test_getCookie_panicsOnReadFailure(c *C) {
	mockReader := &mockConnIOReaderWriter{err: errors.New("stuff")}
	conn := conn{
		rand: mockReader,
	}

	c.Assert(func() {
		conn.getCookie()
	}, PanicMatches, "Failed to read random bytes: EOF")
}
