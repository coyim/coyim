package otr3

import (
	"bytes"
	"testing"
)

func Test_extractWhitespaceTag_removesTagFromMessage(t *testing.T) {
	p := policies(allowV2)
	expectedTag := genWhitespaceTag(p)

	messages := []ValidMessage{
		ValidMessage(string(expectedTag) + "hi there"),
		ValidMessage("hi" + string(expectedTag) + " there"),
		ValidMessage("hi there" + string(expectedTag)),
	}

	for _, m := range messages {
		plain, versions := extractWhitespaceTag(m)

		assertDeepEquals(t, plain, MessagePlaintext("hi there"))
		assertEquals(t, versions, 1<<2)
	}
}

func Test_processWhitespaceTag_shouldNotStartAKEIfPolicyDoesNotAllow(t *testing.T) {
	c := &Conversation{}
	// the policy explicitly is missing whitespaceStartAKE
	c.Policies = policies(allowV2)
	c.ensureAKE()
	assertEquals(t, c.ake.state, authStateNone{})

	expectedTag := genWhitespaceTag(c.Policies)
	m := ValidMessage("hi" + string(expectedTag) + " there")

	plain, toSend, err := c.processWhitespaceTag(m)

	assertNil(t, err)
	assertNil(t, toSend)
	assertDeepEquals(t, plain, MessagePlaintext("hi there"))
	assertEquals(t, c.ake.state, authStateNone{})
}

func Test_genWhitespace_forV2(t *testing.T) {
	hLen := len(whitespaceTagHeader)
	p := policies(allowV2)
	tag := genWhitespaceTag(p)

	assertDeepEquals(t, tag[:hLen], whitespaceTagHeader)
	assertDeepEquals(t, tag[hLen:], otrV2{}.whitespaceTag())
}

func Test_genWhitespace_forV3(t *testing.T) {
	hLen := len(whitespaceTagHeader)
	p := policies(allowV3)
	tag := genWhitespaceTag(p)

	assertDeepEquals(t, tag[:hLen], whitespaceTagHeader)
	assertDeepEquals(t, tag[hLen:], otrV3{}.whitespaceTag())
}

func Test_genWhitespace_forV2AndV3(t *testing.T) {
	hLen := len(whitespaceTagHeader)
	tLen := 8

	p := policies(allowV2 | allowV3)
	tag := genWhitespaceTag(p)

	assertDeepEquals(t, tag[:hLen], whitespaceTagHeader)
	assertDeepEquals(t, tag[hLen:hLen+tLen], otrV2{}.whitespaceTag())
	assertDeepEquals(t, tag[hLen+tLen:], otrV3{}.whitespaceTag())
}

func Test_receive_acceptsV2WhitespaceTagAndStartsAKE(t *testing.T) {
	c := newConversation(nil, fixtureRand())
	c.ourKeys = []PrivateKey{alicePrivateKey}
	c.Policies = policies(allowV2 | whitespaceStartAKE)

	msg := genWhitespaceTag(policies(allowV2))

	_, enc, err := c.Receive(msg)
	toSend, _ := c.decode(encodedMessage(enc[0]))

	assertEquals(t, err, nil)
	assertEquals(t, dhMsgType(toSend), msgTypeDHCommit)
	assertEquals(t, dhMsgVersion(toSend), uint16(2))
}

func Test_receive_ignoresV2WhitespaceTagIfThePolicyDoesNotHaveWhitespaceStartAKE(t *testing.T) {
	c := newConversation(nil, fixtureRand())
	c.ourKeys = []PrivateKey{alicePrivateKey}
	c.Policies = policies(allowV2)

	msg := genWhitespaceTag(policies(allowV2))
	_, enc, err := c.Receive(msg)

	assertNil(t, err)
	assertNil(t, enc)
}

func Test_receive_failsWhenReceivesV2WhitespaceTagIfV2IsNotInThePolicy(t *testing.T) {
	c := newConversation(nil, fixtureRand())
	c.ourKeys = []PrivateKey{alicePrivateKey}
	c.Policies = policies(allowV3 | whitespaceStartAKE)

	msg := genWhitespaceTag(policies(allowV2))

	_, toSend, err := c.Receive(msg)

	assertEquals(t, err, errUnsupportedOTRVersion)
	assertNil(t, toSend)
}

func Test_receive_acceptsV3WhitespaceTagAndStartsAKE(t *testing.T) {
	c := newConversation(nil, fixtureRand())
	c.ourKeys = []PrivateKey{alicePrivateKey}
	c.Policies = policies(allowV2 | allowV3 | whitespaceStartAKE)

	msg := genWhitespaceTag(policies(allowV2 | allowV3))

	_, enc, err := c.Receive(msg)
	toSend, _ := c.decode(encodedMessage(enc[0]))

	assertEquals(t, err, nil)
	assertEquals(t, dhMsgType(toSend), msgTypeDHCommit)
	assertEquals(t, dhMsgVersion(toSend), uint16(3))
}

func Test_receive_whiteSpaceTagWillSignalSetupErrorIfSomethingFails(t *testing.T) {
	c := newConversation(nil, fixedRand([]string{"ABCD"}))
	c.ourKeys = []PrivateKey{alicePrivateKey}
	c.Policies = policies(allowV2 | allowV3 | whitespaceStartAKE)
	msg := genWhitespaceTag(policies(allowV2 | allowV3))

	c.expectMessageEvent(t, func() {
		c.Receive(msg)
	}, MessageEventSetupError, nil, errShortRandomRead)
}

func Test_receive_ignoresV3WhitespaceTagIfThePolicyDoesNotHaveWhitespaceStartAKE(t *testing.T) {
	c := newConversation(nil, fixtureRand())
	c.ourKeys = []PrivateKey{alicePrivateKey}
	c.Policies = policies(allowV2 | allowV3)

	msg := genWhitespaceTag(policies(allowV3))

	_, toSend, err := c.Receive(msg)

	assertEquals(t, err, nil)
	assertNil(t, toSend)
}

func Test_receive_failsWhenReceivesV3WhitespaceTagIfV3IsNotInThePolicy(t *testing.T) {
	c := newConversation(nil, fixtureRand())
	c.ourKeys = []PrivateKey{alicePrivateKey}
	c.Policies = policies(allowV2 | whitespaceStartAKE)

	msg := genWhitespaceTag(policies(allowV3))
	_, toSend, err := c.Receive(msg)

	assertEquals(t, err, errUnsupportedOTRVersion)
	assertNil(t, toSend)
}

func Test_stopAppendingWhitespaceTagsAfterReceivingAPlainMessage(t *testing.T) {
	c := &Conversation{}
	c.ourKeys = []PrivateKey{alicePrivateKey}
	c.Policies = policies(allowV3 | sendWhitespaceTag)

	toSend, err := c.Send([]byte("hi"))
	assertEquals(t, err, nil)
	assertEquals(t, bytes.Contains(toSend[0], whitespaceTagHeader), true)

	toSend, err = c.Send([]byte("are you having fun?"))
	assertEquals(t, err, nil)
	assertEquals(t, bytes.Contains(toSend[0], whitespaceTagHeader), true)

	_, _, err = c.Receive([]byte("no"))
	assertEquals(t, err, nil)

	toSend, err = c.Send([]byte("ok, gotcha"))
	assertEquals(t, err, nil)
	assertEquals(t, bytes.Contains(toSend[0], whitespaceTagHeader), false)

	toSend, err = c.Send([]byte("see ya"))
	assertEquals(t, err, nil)
	assertEquals(t, bytes.Contains(toSend[0], whitespaceTagHeader), false)
}
