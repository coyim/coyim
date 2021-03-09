package xmpp

import (
	"errors"

	"github.com/coyim/coyim/xmpp/data"
	. "gopkg.in/check.v1"
)

type CookieXMPPSuite struct{}

var _ = Suite(&CookieXMPPSuite{})

func (s *CookieXMPPSuite) Test_getCookie_panicsOnReadFailure(c *C) {
	mockReader := &mockConnIOReaderWriter{err: errors.New("stuff")}
	conn := conn{
		log:  testLogger(),
		rand: mockReader,
	}

	c.Assert(func() {
		conn.getCookie()
	}, PanicMatches, "Failed to read random bytes: EOF")
}

func (s *CookieXMPPSuite) Test_conn_cancelInFlights_works(c *C) {
	conn := conn{
		log: testLogger(),
	}

	conn.inflights = map[data.Cookie]inflight{
		data.Cookie(42): inflight{replyChan: make(chan data.Stanza)},
	}

	conn.cancelInflights()

	c.Assert(conn.inflights, HasLen, 0)
}
