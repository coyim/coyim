package xmpp

import (
	"errors"
	"io"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	. "gopkg.in/check.v1"
)

type KeepaliveSuite struct{}

var _ = Suite(&KeepaliveSuite{})

func (s *KeepaliveSuite) Test_conn_sendKeepalive(c *C) {
	out := &mockConnIOReaderWriter{}
	cc := &conn{
		keepaliveOut: out,
		closed:       true,
	}

	c.Assert(cc.sendKeepalive(), Equals, true)

	cc.closed = false
	c.Assert(cc.sendKeepalive(), Equals, true)

	out.err = io.EOF
	out.errCount = 0
	c.Assert(cc.sendKeepalive(), Equals, true)

	out.err = errors.New("oh no")
	out.errCount = 0
	c.Assert(cc.sendKeepalive(), Equals, false)
}

func (s *KeepaliveSuite) Test_conn_watchKeepAlive_closed(c *C) {
	orgKeepaliveInterval := keepaliveInterval
	defer func() {
		keepaliveInterval = orgKeepaliveInterval
	}()

	keepaliveInterval = 1 * time.Millisecond

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	cc := &conn{
		closed: true,
		log:    l,
	}

	cc.watchKeepAlive()

	c.Assert(len(hook.Entries), Equals, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.InfoLevel)
	c.Assert(hook.Entries[0].Message, Equals, "xmpp: no more watching keepalives")
}

func (s *KeepaliveSuite) Test_conn_watchKeepAlive_workingOnce(c *C) {
	orgKeepaliveInterval := keepaliveInterval
	defer func() {
		keepaliveInterval = orgKeepaliveInterval
	}()

	keepaliveInterval = 1 * time.Millisecond

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	out := &mockConnIOReaderWriter{}

	cc := &conn{
		closed:       false,
		log:          l,
		keepaliveOut: out,
		out:          out,
	}

	out.errCount = 1
	out.err = errors.New("what up?")

	cc.watchKeepAlive()

	c.Assert(len(hook.Entries), Equals, 2)
	c.Assert(hook.Entries[0].Level, Equals, log.InfoLevel)
	c.Assert(hook.Entries[0].Message, Equals, "xmpp: keepalive failed")
	c.Assert(hook.Entries[1].Level, Equals, log.InfoLevel)
	c.Assert(hook.Entries[1].Message, Equals, "xmpp: no more watching keepalives")
}
