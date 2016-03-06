package client

import (
	"errors"

	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/otr3"
)
import . "github.com/twstrike/coyim/Godeps/_workspace/src/gopkg.in/check.v1"

type ConversationSuite struct{}

var _ = Suite(&ConversationSuite{})

func (s *ConversationSuite) Test_StartEncryptedChat_startsAnEncryptedChat(c *C) {
	cb := &conversation{"foo@bar.com", &otr3.Conversation{}}
	ts := &testSender{err: nil}
	e := cb.StartEncryptedChat(ts)

	c.Assert(e, IsNil)
	c.Assert(ts.peer, Equals, "foo@bar.com")
	c.Assert(ts.msg, Equals, "?OTRv? Your peer has requested a private conversation with you, but your client doesn't seem to support the OTR protocol.")
}

func (s *ConversationSuite) Test_sendAll_returnsTheFirstErrorEncountered(c *C) {
	cb := &conversation{"foo@bar.com", &otr3.Conversation{}}
	ts := &testSender{err: errors.New("hello")}
	e := cb.sendAll(ts, []otr3.ValidMessage{otr3.ValidMessage([]byte("Hello there"))})

	c.Assert(e, DeepEquals, errors.New("hello"))
}

func (s *ConversationSuite) Test_sendAll_sendsTheMessageGiven(c *C) {
	cb := &conversation{"foo@bar.com", &otr3.Conversation{}}
	ts := &testSender{err: nil}
	e := cb.sendAll(ts, []otr3.ValidMessage{otr3.ValidMessage([]byte("Hello there"))})

	c.Assert(e, IsNil)
	c.Assert(ts.peer, Equals, "foo@bar.com")
	c.Assert(ts.msg, Equals, "Hello there")
}

func (s *ConversationSuite) Test_Send_(c *C) {
	cb := &conversation{"foo@bar.com", &otr3.Conversation{}}
	ts := &testSender{err: nil}
	e := cb.Send(ts, []byte("Hello there"))

	c.Assert(e, IsNil)
	c.Assert(ts.peer, Equals, "foo@bar.com")
	c.Assert(ts.msg, Equals, "Hello there")
}
