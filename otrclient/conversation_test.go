package otrclient

import (
	"errors"

	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/otr3"

	. "gopkg.in/check.v1"
)

type ConversationSuite struct{}

var _ = Suite(&ConversationSuite{})

func (s *ConversationSuite) Test_StartEncryptedChat_startsAnEncryptedChat(c *C) {
	ts := &testSender{err: nil}
	cb := &conversation{jid.NR("foo@bar.com"), false, ts, nil, &otr3.Conversation{}}

	cb.SetFriendlyQueryMessage("Your peer has requested a private conversation with you, but your client doesn't seem to support the OTR protocol.")
	e := cb.StartEncryptedChat()

	c.Assert(e, IsNil)
	c.Assert(ts.peer, Equals, jid.NR("foo@bar.com"))
	c.Assert(ts.msg, Equals, "?OTRv? Your peer has requested a private conversation with you, but your client doesn't seem to support the OTR protocol.")
}

func (s *ConversationSuite) Test_sendAll_returnsTheFirstErrorEncountered(c *C) {
	ts := &testSender{err: errors.New("hello")}
	cb := &conversation{jid.NR("foo@bar.com"), false, ts, nil, &otr3.Conversation{}}
	e := cb.sendAll([]otr3.ValidMessage{otr3.ValidMessage([]byte("Hello there"))})

	c.Assert(e, DeepEquals, errors.New("hello"))
}

func (s *ConversationSuite) Test_sendAll_sendsTheMessageGiven(c *C) {
	ts := &testSender{err: nil}
	cb := &conversation{jid.NR("foo@bar.com"), false, ts, nil, &otr3.Conversation{}}
	e := cb.sendAll([]otr3.ValidMessage{otr3.ValidMessage([]byte("Hello there"))})

	c.Assert(e, IsNil)
	c.Assert(ts.peer, Equals, jid.NR("foo@bar.com"))
	c.Assert(ts.msg, Equals, "Hello there")
}

func (s *ConversationSuite) Test_Send_(c *C) {
	ts := &testSender{err: nil}
	cb := &conversation{jid.NR("foo@bar.com"), false, ts, nil, &otr3.Conversation{}}
	_, e := cb.Send([]byte("Hello there"))

	c.Assert(e, IsNil)
	c.Assert(ts.peer, Equals, jid.NR("foo@bar.com"))
	c.Assert(ts.msg, Equals, "Hello there")
}
