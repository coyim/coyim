package otr3

import (
	"crypto/rand"
	"testing"
	"time"
)

func fixtureCorrectResend(c *Conversation) {
	c.resend.later(MessagePlaintext("hello"))
	c.resend.mayRetransmit = retransmitExact
	c.updateLastSent()
}

func Test_shouldRetransmit_returnsFalseIfThereIsNoLastMessage(t *testing.T) {
	c := &Conversation{}
	fixtureCorrectResend(c)
	c.resend.clear()

	assertEquals(t, c.shouldRetransmit(), false)
}

func Test_shouldRetransmit_returnTrueIfAllTheConditionsForResendingAreMet(t *testing.T) {
	c := &Conversation{}
	fixtureCorrectResend(c)

	assertEquals(t, c.shouldRetransmit(), true)
}

func Test_shouldRetransmit_returnFalseIfTheLastMessageWasSentTooFarBackInTime(t *testing.T) {
	c := &Conversation{}
	fixtureCorrectResend(c)
	c.heartbeat.lastSent = time.Now().Add(-61 * time.Second)

	assertEquals(t, c.shouldRetransmit(), false)
}

func Test_shouldRetransmit_returnTrueWhenFlagIsRetransmitWithPrefix(t *testing.T) {
	c := &Conversation{}
	fixtureCorrectResend(c)
	c.resend.mayRetransmit = retransmitWithPrefix

	assertEquals(t, c.shouldRetransmit(), true)
}

func Test_shouldRetransmit_returnFalseWhenFlagIsNoRetransmit(t *testing.T) {
	c := &Conversation{}
	fixtureCorrectResend(c)
	c.resend.mayRetransmit = noRetransmit

	assertEquals(t, c.shouldRetransmit(), false)
}

func Test_maybeRetransmit_returnsNothingWhenShouldntRetransmit(t *testing.T) {
	c := &Conversation{}
	fixtureCorrectResend(c)
	c.resend.clear()

	res, err := c.maybeRetransmit()

	assertNil(t, res)
	assertNil(t, err)
}

func Test_maybeRetransmit_createsADataMessageWithTheExactMessageWhenAskedToRetransmitExact(t *testing.T) {
	c := newConversation(otrV3{}, rand.Reader)
	c.Policies.add(allowV3)
	c.ourCurrentKey = bobPrivateKey
	c.smp.secret = bnFromHex("ABCDE56321F9A9F8E364607C8C82DECD8E8E6209E2CB952C7E649620F5286FE3")

	plain := plainDataMsg{
		message: []byte(""),
	}

	_, c.keys = fixtureDataMsg(plain)

	c.msgState = encrypted

	fixtureCorrectResend(c)
	c.resend.clear()
	c.resend.later(MessagePlaintext("Something else to think about"))

	res, err := c.maybeRetransmit()
	assertNil(t, err)
	dec := fixtureDecryptDataMsg(res[0])

	assertDeepEquals(t, MessagePlaintext(dec.message), MessagePlaintext("Something else to think about"))
	assertEquals(t, len(dec.tlvs), 1)
	assertEquals(t, dec.tlvs[0].tlvType, tlvTypePadding)
}

func Test_maybeRetransmit_createsADataMessageWithTheResendPrefixAndMessageWhenAskedToRetransmitWithPrefix(t *testing.T) {
	c := newConversation(otrV3{}, rand.Reader)
	c.Policies.add(allowV3)
	c.ourCurrentKey = bobPrivateKey
	c.smp.secret = bnFromHex("ABCDE56321F9A9F8E364607C8C82DECD8E8E6209E2CB952C7E649620F5286FE3")

	plain := plainDataMsg{
		message: []byte(""),
	}

	_, c.keys = fixtureDataMsg(plain)

	c.msgState = encrypted

	fixtureCorrectResend(c)
	c.resend.clear()
	c.resend.mayRetransmit = retransmitWithPrefix
	c.resend.later(MessagePlaintext("Something else to think about"))

	res, err := c.maybeRetransmit()
	dec := fixtureDecryptDataMsg(res[0])

	assertNil(t, err)
	assertDeepEquals(t, MessagePlaintext(dec.message), MessagePlaintext("[resent] Something else to think about"))
	assertEquals(t, len(dec.tlvs), 1)
	assertEquals(t, dec.tlvs[0].tlvType, tlvTypePadding)
}

func Test_maybeRetransmit_createsADataMessageWithTheCustomResendPrefixAndMessageWhenAskedToRetransmitWithPrefix(t *testing.T) {
	c := newConversation(otrV3{}, rand.Reader)
	c.Policies.add(allowV3)
	c.ourCurrentKey = bobPrivateKey
	c.smp.secret = bnFromHex("ABCDE56321F9A9F8E364607C8C82DECD8E8E6209E2CB952C7E649620F5286FE3")

	plain := plainDataMsg{
		message: []byte(""),
	}

	_, c.keys = fixtureDataMsg(plain)

	c.msgState = encrypted

	fixtureCorrectResend(c)
	c.resend.clear()
	c.resend.mayRetransmit = retransmitWithPrefix
	c.resend.later(MessagePlaintext("Something much more to think about"))
	c.resend.messageTransform = func(msg []byte) []byte {
		return append(append([]byte("<resend>"), msg...), []byte("</resend>")...)
	}

	res, err := c.maybeRetransmit()
	dec := fixtureDecryptDataMsg(res[0])

	assertNil(t, err)
	assertDeepEquals(t, MessagePlaintext(dec.message), MessagePlaintext("<resend>Something much more to think about</resend>"))
	assertEquals(t, len(dec.tlvs), 1)
	assertEquals(t, dec.tlvs[0].tlvType, tlvTypePadding)
}

func Test_maybeRetransmit_updatesLastSentWhenSendingAMessage(t *testing.T) {
	c := newConversation(otrV3{}, rand.Reader)
	c.Policies.add(allowV3)
	c.ourCurrentKey = bobPrivateKey
	c.smp.secret = bnFromHex("ABCDE56321F9A9F8E364607C8C82DECD8E8E6209E2CB952C7E649620F5286FE3")

	plain := plainDataMsg{
		message: []byte(""),
	}

	_, c.keys = fixtureDataMsg(plain)

	c.msgState = encrypted

	fixtureCorrectResend(c)

	setSent := time.Now().Add(-30 * time.Second)
	c.heartbeat.lastSent = setSent

	c.maybeRetransmit()

	assertNotEquals(t, c.heartbeat.lastSent, setSent)
}

func Test_maybeRetransmit_returnsErrorIfWeFailAtGeneratingDataMsg(t *testing.T) {
	c := newConversation(otrV3{}, rand.Reader)
	c.Policies.add(allowV3)
	c.ourCurrentKey = bobPrivateKey
	c.smp.secret = bnFromHex("ABCDE56321F9A9F8E364607C8C82DECD8E8E6209E2CB952C7E649620F5286FE3")

	plain := plainDataMsg{
		message: []byte(""),
	}

	_, c.keys = fixtureDataMsg(plain)
	c.msgState = encrypted
	fixtureCorrectResend(c)
	c.keys.ourKeyID = 0
	_, err := c.maybeRetransmit()

	assertEquals(t, err, newOtrConflictError("invalid key id for local peer"))
}

func Test_maybeRetransmit_signalsMessageEventWhenResendingMessage(t *testing.T) {
	c := newConversation(otrV3{}, rand.Reader)
	c.Policies.add(allowV3)
	c.ourCurrentKey = bobPrivateKey
	c.smp.secret = bnFromHex("ABCDE56321F9A9F8E364607C8C82DECD8E8E6209E2CB952C7E649620F5286FE3")

	plain := plainDataMsg{
		message: []byte(""),
	}

	_, c.keys = fixtureDataMsg(plain)

	c.msgState = encrypted

	fixtureCorrectResend(c)
	c.resend.mayRetransmit = retransmitWithPrefix

	c.expectMessageEvent(t, func() {
		c.maybeRetransmit()
	}, MessageEventMessageResent, nil, nil)
}

func Test_maybeRetransmit_signalMessageEventWhenSendingMessageExact(t *testing.T) {
	c := newConversation(otrV3{}, rand.Reader)
	c.Policies.add(allowV3)
	c.ourCurrentKey = bobPrivateKey
	c.smp.secret = bnFromHex("ABCDE56321F9A9F8E364607C8C82DECD8E8E6209E2CB952C7E649620F5286FE3")

	plain := plainDataMsg{
		message: []byte(""),
	}

	_, c.keys = fixtureDataMsg(plain)

	c.msgState = encrypted

	fixtureCorrectResend(c)
	c.resend.mayRetransmit = retransmitExact

	c.expectMessageEvent(t, func() {
		c.maybeRetransmit()
	}, MessageEventMessageSent, nil, nil)
}
