package otrclient

import (
	log "github.com/sirupsen/logrus"

	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/otr3"

	. "gopkg.in/check.v1"
)

type ConversationManagerSuite struct{}

var _ = Suite(&ConversationManagerSuite{})

type testSender struct {
	peer jid.Any
	msg  string
	err  error
}

type testConvBuilder struct {
	fake *otr3.Conversation
}

func (cb *testConvBuilder) NewConversation(peer jid.Any) *otr3.Conversation {
	return cb.fake
}

func (ts *testSender) Send(peer jid.Any, msg string, otr bool) error {
	ts.peer = peer
	ts.msg = msg
	return ts.err
}

func (s *ConversationManagerSuite) Test_TerminateAll_willTerminate(c *C) {
	cb := &testConvBuilder{&otr3.Conversation{}}
	ts := &testSender{err: nil}
	mgr := NewConversationManager(cb.NewConversation, ts, "blarg", func(jid.Any, *EventHandler, chan string, chan int) {}, log.New().WithFields(log.Fields{}))
	conv, created := mgr.EnsureConversationWith(jid.NR("someone@whitehouse.gov"), nil)

	c.Assert(created, Equals, true)
	c.Assert(conv, Not(IsNil))

	mgr.TerminateAll()

	c.Assert(ts.msg, Equals, "")
}
