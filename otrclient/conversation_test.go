package otrclient

import (
	"crypto/rand"
	"errors"

	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/otr3"

	. "gopkg.in/check.v1"
)

type ConversationSuite struct{}

var _ = Suite(&ConversationSuite{})

func (s *ConversationSuite) Test_StartEncryptedChat_startsAnEncryptedChat(c *C) {
	ts := &testSender{err: nil}
	cb := &conversation{jid.NR("foo@bar.com"), false, ts, nil, &otr3.Conversation{}, 0, nil, nil}

	cb.SetFriendlyQueryMessage("Your peer has requested a private conversation with you, but your client doesn't seem to support the OTR protocol.")
	e := cb.StartEncryptedChat()

	c.Assert(e, IsNil)
	c.Assert(ts.peer, Equals, jid.NR("foo@bar.com"))
	c.Assert(ts.msg, Equals, "?OTRv? Your peer has requested a private conversation with you, but your client doesn't seem to support the OTR protocol.")
}

func (s *ConversationSuite) Test_sendAll_returnsTheFirstErrorEncountered(c *C) {
	ts := &testSender{err: errors.New("hello")}
	cb := &conversation{jid.NR("foo@bar.com"), false, ts, nil, &otr3.Conversation{}, 0, nil, nil}
	e := cb.sendAll([]otr3.ValidMessage{otr3.ValidMessage([]byte("Hello there"))})

	c.Assert(e, DeepEquals, errors.New("hello"))
}

func (s *ConversationSuite) Test_sendAll_sendsTheMessageGiven(c *C) {
	ts := &testSender{err: nil}
	cb := &conversation{jid.NR("foo@bar.com"), false, ts, nil, &otr3.Conversation{}, 0, nil, nil}
	e := cb.sendAll([]otr3.ValidMessage{otr3.ValidMessage([]byte("Hello there"))})

	c.Assert(e, IsNil)
	c.Assert(ts.peer, Equals, jid.NR("foo@bar.com"))
	c.Assert(ts.msg, Equals, "Hello there")
}

func (s *ConversationSuite) Test_Send_(c *C) {
	ts := &testSender{err: nil}
	cb := &conversation{jid.NR("foo@bar.com"), false, ts, nil, &otr3.Conversation{}, 0, nil, nil}
	_, e := cb.Send([]byte("Hello there"))

	c.Assert(e, IsNil)
	c.Assert(ts.peer, Equals, jid.NR("foo@bar.com"))
	c.Assert(ts.msg, Equals, "Hello there")
}

func (s *ConversationSuite) Test_Conversation_EventHandler_returnsTheEventHandler(c *C) {
	ev := &EventHandler{}
	conv := &conversation{eh: ev}

	c.Assert(conv.EventHandler(), Equals, ev)
}

type errReader struct {
	e error
}

func (r *errReader) Read([]byte) (int, error) {
	return 0, r.e
}

func setUpCompletedConversationState() (alice, bob *otr3.Conversation) {
	alice = &otr3.Conversation{Rand: rand.Reader}
	alice.Policies.AllowV3()
	alice.SetOurKeys([]otr3.PrivateKey{alicePrivateKey})

	bob = &otr3.Conversation{Rand: rand.Reader}
	bob.Policies.AllowV3()
	bob.SetOurKeys([]otr3.PrivateKey{bobPrivateKey})

	msg := []byte("?OTRv3?")

	var toSend []otr3.ValidMessage

	_, toSend, _ = bob.Receive(msg)
	_, toSend, _ = alice.Receive(toSend[0])
	_, toSend, _ = bob.Receive(toSend[0])
	_, toSend, _ = alice.Receive(toSend[0])
	_, _, _ = bob.Receive(toSend[0])

	return
}

func (s *ConversationSuite) Test_Conversation_Send_failsOnSending(c *C) {
	alice, bob := setUpCompletedConversationState()
	toSend, _ := bob.End()
	_, _, _ = alice.Receive(toSend[0])

	conv := &conversation{Conversation: alice}

	_, e := conv.Send([]byte("hello"))
	c.Assert(e, ErrorMatches, "otr: cannot send message because secure conversation has finished")
}

func (s *ConversationSuite) Test_Conversation_Receive_getsTheMessage(c *C) {
	ts := &testSender{err: nil}
	alice, bob := setUpCompletedConversationState()
	toSend, _ := bob.Send(otr3.ValidMessage("hello there"))

	conv := &conversation{Conversation: alice, s: ts}

	plain, e := conv.Receive([]byte(toSend[0]))
	c.Assert(string(plain), Equals, "hello there")
	c.Assert(e, IsNil)
}

func (s *ConversationSuite) Test_Conversation_Receive_failsWhenOTR3fails(c *C) {
	ts := &testSender{err: nil}
	alice, bob := setUpCompletedConversationState()
	toSend, _ := bob.Send(otr3.ValidMessage("hello there"))
	toSendEnd, _ := bob.End()
	_, _, _ = alice.Receive(toSendEnd[0])

	conv := &conversation{Conversation: alice, s: ts}

	plain, e := conv.Receive([]byte(toSend[0]))
	c.Assert(string(plain), Equals, "")
	c.Assert(e, ErrorMatches, "otr: message not in private")
}

func (s *ConversationSuite) Test_Conversation_ProvideAuthenticationSecret(c *C) {
	ts := &testSender{err: errors.New("return marker")}
	alice, bob := setUpCompletedConversationState()
	toSend, _ := bob.StartAuthenticate("q", []byte("a"))
	_, _, _ = alice.Receive(toSend[0])

	conv := &conversation{Conversation: alice, s: ts}

	e := conv.ProvideAuthenticationSecret([]byte("a secret"))
	c.Assert(e, ErrorMatches, "return marker")
}

func (s *ConversationSuite) Test_Conversation_ProvideAuthenticationSecret_fails(c *C) {
	alice, _ := setUpCompletedConversationState()
	_, _ = alice.StartAuthenticate("q", []byte("a"))

	conv := &conversation{Conversation: alice}

	e := conv.ProvideAuthenticationSecret([]byte("a secret"))
	c.Assert(e, ErrorMatches, "otr: not expected SMP secret to be provided now")
}

func (s *ConversationSuite) Test_Conversation_StartAuthenticate_works(c *C) {
	ts := &testSender{err: errors.New("return marker")}
	alice, _ := setUpCompletedConversationState()

	conv := &conversation{Conversation: alice, s: ts}

	e := conv.StartAuthenticate("q", []byte("a secret"))
	c.Assert(e, ErrorMatches, "return marker")
}

func (s *ConversationSuite) Test_Conversation_StartAuthenticate_fails(c *C) {
	ts := &testSender{err: errors.New("return marker")}
	alice, bob := setUpCompletedConversationState()
	toSendEnd, _ := bob.End()
	_, _, _ = alice.Receive(toSendEnd[0])

	conv := &conversation{Conversation: alice, s: ts}

	e := conv.StartAuthenticate("q", []byte("a secret"))
	c.Assert(e, ErrorMatches, "otr: can't authenticate a peer without a secure conversation established")
}

func (s *ConversationSuite) Test_Conversation_AbortAuthentication_works(c *C) {
	ts := &testSender{err: errors.New("return marker")}
	alice, bob := setUpCompletedConversationState()
	toSend, _ := bob.StartAuthenticate("q", []byte("a"))
	_, _, _ = alice.Receive(toSend[0])

	conv := &conversation{Conversation: alice, s: ts}

	e := conv.AbortAuthentication()
	c.Assert(e, ErrorMatches, "return marker")
}

func (s *ConversationSuite) Test_Conversation_AbortAuthentication_fails(c *C) {
	ts := &testSender{err: errors.New("return marker")}
	alice, bob := setUpCompletedConversationState()
	toSendEnd, _ := bob.End()
	_, _, _ = alice.Receive(toSendEnd[0])

	conv := &conversation{Conversation: alice, s: ts}

	e := conv.AbortAuthentication()
	c.Assert(e, ErrorMatches, "otr: cannot send message in unencrypted state")
}

func (s *ConversationSuite) Test_Conversation_OurFingerprint_failsIfNoKeyExist(c *C) {
	alice := &otr3.Conversation{}
	conv := &conversation{Conversation: alice}

	c.Assert(conv.OurFingerprint(), IsNil)
}

func (s *ConversationSuite) Test_Conversation_OurFingerprint_works(c *C) {
	alice, _ := setUpCompletedConversationState()
	conv := &conversation{Conversation: alice}

	c.Assert(conv.OurFingerprint(), DeepEquals, []byte{0xb, 0xb0, 0x1c, 0x36, 0x4, 0x24, 0x52, 0x2e, 0x94, 0xee, 0x9c, 0x34, 0x6c, 0xe8, 0x77, 0xa1, 0xa4, 0x28, 0x8b, 0x2f})
}

func (s *ConversationSuite) Test_Conversation_TheirFingerprint_failsIfNoKeyExist(c *C) {
	alice := &otr3.Conversation{}
	conv := &conversation{Conversation: alice}

	c.Assert(conv.TheirFingerprint(), IsNil)
}

func (s *ConversationSuite) Test_Conversation_TheirFingerprint_works(c *C) {
	alice, _ := setUpCompletedConversationState()
	conv := &conversation{Conversation: alice}

	c.Assert(conv.TheirFingerprint(), DeepEquals, []byte{0x87, 0x98, 0xfa, 0xa7, 0x73, 0x52, 0x67, 0xfb, 0x84, 0x57, 0x73, 0x30, 0x98, 0x48, 0x2e, 0x94, 0x9, 0x6d, 0x4a, 0xbd})
}

func (s *ConversationSuite) Test_Conversation_CreateExtraSymmetricKey_returnsErrorWhenFailing(c *C) {
	ts := &testSender{err: errors.New("return marker")}
	alice, bob := setUpCompletedConversationState()
	toSendEnd, _ := bob.End()
	_, _, _ = alice.Receive(toSendEnd[0])

	conv := &conversation{Conversation: alice, s: ts}

	_, e := conv.CreateExtraSymmetricKey()
	c.Assert(e, ErrorMatches, "otr: cannot send message in current state")
}

func (s *ConversationSuite) Test_Conversation_CreateExtraSymmetricKey_succeeds(c *C) {
	ts := &testSender{err: errors.New("return marker")}
	alice, _ := setUpCompletedConversationState()

	conv := &conversation{Conversation: alice, s: ts}

	_, e := conv.CreateExtraSymmetricKey()
	c.Assert(e, ErrorMatches, "return marker")
}

func (s *ConversationSuite) Test_Conversation_ReceivedSymmetricKey(c *C) {
	conv := &conversation{}
	conv.ReceivedSymmetricKey(0x01, []byte("foo"), []byte("bar"))
	c.Assert(conv.lastExtraKeyUsage, Equals, uint32(0x01))
	c.Assert(conv.lastExtraKeyUsageData, DeepEquals, []byte("foo"))
	c.Assert(conv.lastExtraKeySymkey, DeepEquals, []byte("bar"))
}

func (s *ConversationSuite) Test_Conversation_GetAndWipeLastExtraKey(c *C) {
	conv := &conversation{}
	conv.ReceivedSymmetricKey(0x01, []byte("foo"), []byte("bar"))
	ret1, ret2, ret3 := conv.GetAndWipeLastExtraKey()

	c.Assert(ret1, Equals, uint32(0x01))
	c.Assert(ret2, DeepEquals, []byte("foo"))
	c.Assert(ret3, DeepEquals, []byte("bar"))

	c.Assert(conv.lastExtraKeyUsage, Equals, uint32(0x00))
	c.Assert(conv.lastExtraKeyUsageData, IsNil)
	c.Assert(conv.lastExtraKeySymkey, IsNil)
}
