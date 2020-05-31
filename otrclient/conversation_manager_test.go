package otrclient

import (
	"io/ioutil"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/glib_mock"
	"github.com/coyim/otr3"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

func init() {
	log.SetOutput(ioutil.Discard)
	i18n.InitLocalization(&glib_mock.Mock{})
}

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

func (ts *testSender) Send(peer jid.Any, msg string) error {
	ts.peer = peer
	ts.msg = msg
	return ts.err
}

func (s *ConversationManagerSuite) Test_TerminateAll_willTerminate(c *C) {
	cb := &testConvBuilder{&otr3.Conversation{}}
	ts := &testSender{err: nil}
	mgr := NewConversationManager(cb.NewConversation, ts, "blarg", func(jid.Any, *EventHandler, chan string, chan int) {}, log.New().WithFields(log.Fields{}))
	conv, created := mgr.EnsureConversationWith(jid.NR("someone@whitehouse.gov"))

	c.Assert(created, Equals, true)
	c.Assert(conv, Not(IsNil))

	mgr.TerminateAll()

	c.Assert(ts.msg, Equals, "")
}
