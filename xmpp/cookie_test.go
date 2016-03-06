package xmpp

import (
	"errors"

	. "github.com/twstrike/coyim/Godeps/_workspace/src/gopkg.in/check.v1"
)

type CookieXmppSuite struct{}

var _ = Suite(&CookieXmppSuite{})

func (s *CookieXmppSuite) Test_getCookie_panicsOnReadFailure(c *C) {
	mockReader := &mockConnIOReaderWriter{err: errors.New("stuff")}
	conn := conn{
		rand: mockReader,
	}

	c.Assert(func() {
		conn.getCookie()
	}, PanicMatches, "Failed to read random bytes: EOF")
}
