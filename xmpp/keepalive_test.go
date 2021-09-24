package xmpp

import (
	"errors"
	"io"
	"sort"
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

	ll := len(hook.Entries)
	// We can get either two or three log entries here, depending on whether the sending
	// of the stream error happens or not - if it happens, we will get a log saying
	// the connection is already closed.
	c.Assert(ll == 2 || ll == 3, Equals, true)
	c.Assert(hook.Entries[0].Level, Equals, log.InfoLevel)
	c.Assert(hook.Entries[1].Level, Equals, log.InfoLevel)

	messages := []string{hook.Entries[0].Message, hook.Entries[1].Message}
	if ll == 3 {
		messages = append(messages, hook.Entries[2].Message)
	}

	sort.Strings(messages)
	c.Assert(messages[0:2], DeepEquals, []string{
		"xmpp: keepalive failed",
		"xmpp: no more watching keepalives",
	})
}
