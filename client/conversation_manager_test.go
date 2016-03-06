package client

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glib_mock"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/otr3"
	"github.com/twstrike/coyim/i18n"

	. "github.com/twstrike/coyim/Godeps/_workspace/src/gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

func init() {
	log.SetOutput(ioutil.Discard)
	i18n.InitLocalization(&glib_mock.Mock{})
}

type ConversationManagerSuite struct{}

var _ = Suite(&ConversationManagerSuite{})

type testSender struct {
	peer, msg string
	err       error
}

type testConvBuilder struct {
	fake *otr3.Conversation
}

func (cb *testConvBuilder) NewConversation(peer string) *otr3.Conversation {
	return cb.fake
}

func (ts *testSender) Send(peer, msg string) error {
	ts.peer = peer
	ts.msg = msg
	return ts.err
}

func (s *ConversationManagerSuite) Test_TerminateAll_willTerminate(c *C) {
	cb := &testConvBuilder{&otr3.Conversation{}}
	ts := &testSender{err: nil}
	mgr := NewConversationManager(cb, ts)
	conv, created := mgr.EnsureConversationWith("someone@whitehouse.gov")

	c.Assert(created, Equals, true)
	c.Assert(conv, Not(IsNil))

	mgr.TerminateAll()

	c.Assert(ts.msg, Equals, "")
}
