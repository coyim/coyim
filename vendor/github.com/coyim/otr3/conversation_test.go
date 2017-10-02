package otr3

import (
	"crypto/rand"
	"math/big"
	"testing"
)

func Test_receive_OTRQueryMsgRepliesWithDHCommitMessage(t *testing.T) {
	msg := []byte("?OTRv3?")
	c := newConversation(nil, fixtureRand())
	c.SetOurKeys([]PrivateKey{bobPrivateKey})
	c.Policies.add(allowV3)

	exp := messageWithHeader{
		0x00, 0x03, // protocol version
		msgTypeDHCommit,
	}

	_, enc, err := c.Receive(msg)
	assertEquals(t, err, nil)

	toSend, _ := c.decode(encodedMessage(enc[0]))

	assertDeepEquals(t, toSend[:3], exp)
}

func Test_receive_OTRQueryMsgChangesContextProtocolVersion(t *testing.T) {
	msg := []byte("?OTRv3?")
	c := newConversation(nil, fixtureRand())
	c.SetOurKeys([]PrivateKey{bobPrivateKey})
	c.Policies.add(allowV3)

	_, _, err := c.Receive(msg)

	assertEquals(t, err, nil)
	assertDeepEquals(t, c.version, otrV3{})
}

func Test_receive_verifiesMessageProtocolVersion(t *testing.T) {
	// protocol version
	msg := []byte{0x00, 0x02, 0x00, msgTypeDHKey}
	c := newConversation(otrV3{}, fixtureRand())
	c.SetOurKeys([]PrivateKey{bobPrivateKey})

	_, _, err := c.receiveDecoded(msg)

	assertEquals(t, err, errWrongProtocolVersion)
}

func Test_receive_returnsAnErrorForAnInvalidOTRMessageWithoutVersionData(t *testing.T) {
	msg := []byte{0x00}
	c := newConversation(otrV3{}, fixtureRand())
	c.SetOurKeys([]PrivateKey{bobPrivateKey})

	_, _, err := c.receiveDecoded(msg)

	assertEquals(t, err, errInvalidOTRMessage)
}

func Test_receive_ignoresAMessageWhenNoEncryptionIsActive(t *testing.T) {
	m := []byte{
		0x00, 0x03, // protocol version
		msgTypeData,
		0x00, 0x00, 0x01, 0x01,
		0x00, 0x00, 0x01, 0x01,
	}
	c := newConversation(otrV3{}, fixtureRand())
	c.SetOurKeys([]PrivateKey{bobPrivateKey})

	a, b, err := c.receiveDecoded(m)
	assertNil(t, a)
	assertNil(t, b)
	assertEquals(t, err, errMessageNotInPrivate)
}

func Test_receiveDecoded_signalsAMessageEventForADataMessageWhenNoEncryptionIsActive(t *testing.T) {
	m := []byte{
		0x00, 0x03, // protocol version
		msgTypeData,
		0x00, 0x00, 0x01, 0x01,
		0x00, 0x00, 0x01, 0x01,
	}
	c := newConversation(otrV3{}, fixtureRand())
	c.SetOurKeys([]PrivateKey{bobPrivateKey})

	c.expectMessageEvent(t, func() {
		c.receiveDecoded(m)
	}, MessageEventReceivedMessageNotInPrivate, nil, nil)
}

func Test_receive_DHCommitMessageReturnsDHKeyForOTR3(t *testing.T) {
	exp := messageWithHeader{
		0x00, 0x03, // protocol version
		msgTypeDHKey,
	}

	dhCommitAKE := fixtureConversation()
	dhCommitMsg, _ := dhCommitAKE.dhCommitMessage()
	dhCommitMsg, _ = dhCommitAKE.wrapMessageHeader(msgTypeDHCommit, dhCommitMsg)

	c := newConversation(otrV3{}, fixtureRand())
	c.Policies.add(allowV3)

	_, dhKeyMsg, err := c.receiveDecoded(dhCommitMsg)

	assertEquals(t, err, nil)
	assertDeepEquals(t, dhKeyMsg[0][:messageHeaderPrefix], exp)
}

func Test_receive_DHKeyMessageReturnsRevealSignature(t *testing.T) {
	v := otrV3{}

	msg := fixtureDHKeyMsg(v)
	c := bobContextAtAwaitingDHKey()

	_, toSend, err := c.receiveDecoded(msg)

	assertEquals(t, err, nil)
	assertDeepEquals(t, dhMsgType(toSend[0]), msgTypeRevealSig)
}

func Test_OTRisDisabledIfNoVersionIsAllowedInThePolicy(t *testing.T) {
	msg := []byte("?OTRv3?")

	c := newConversation(nil, fixtureRand())

	s, _ := c.Send(msg)
	assertDeepEquals(t, s, []ValidMessage{msg})

	_, r, err := c.Receive(msg)
	assertNil(t, err)
	assertNil(t, r)
}

func Test_Send_returnsErrorIfFaislToGenerateDataMsg(t *testing.T) {
	msg := []byte("hello")

	c := bobContextAfterAKE()
	c.msgState = encrypted
	c.Policies = policies(allowV3)
	c.keys.theirKeyID = 0
	s, err := c.Send(msg)

	assertNil(t, s)
	assertEquals(t, err, newOtrConflictError("invalid key id for remote peer"))
}

func Test_send_appendWhitespaceTagsWhenAllowedbyThePolicy(t *testing.T) {
	expectedWhitespaceTag := ValidMessage{
		0x20, 0x09, 0x20, 0x20, 0x09, 0x09, 0x09, 0x09,
		0x20, 0x09, 0x20, 0x09, 0x20, 0x09, 0x20, 0x20,
		0x20, 0x20, 0x09, 0x09, 0x20, 0x20, 0x09, 0x09,
	}

	c := &Conversation{}
	c.Policies = policies(allowV3 | sendWhitespaceTag)

	m, _ := c.Send([]byte("hello"))
	wsPos := len(m[0]) - len(expectedWhitespaceTag)
	assertDeepEquals(t, m[0][wsPos:], expectedWhitespaceTag)
}

func Test_send_doesNotAppendWhitespaceTagsWhenItsNotAllowedbyThePolicy(t *testing.T) {
	m := []byte("hello")
	c := &Conversation{}
	c.Policies = policies(allowV3)

	toSend, _ := c.Send(m)
	assertDeepEquals(t, toSend, []ValidMessage{m})
}

func Test_send_appendWhitespaceTagsIfReceivesNonDHCommitMsgBeforeSendingFirstMessage(t *testing.T) {
	hello := ValidMessage("hello")
	expectedWhitespaceTag := ValidMessage{
		0x20, 0x09, 0x20, 0x20, 0x09, 0x09, 0x09, 0x09,
		0x20, 0x09, 0x20, 0x09, 0x20, 0x09, 0x20, 0x20,
		0x20, 0x20, 0x09, 0x09, 0x20, 0x20, 0x09, 0x09,
	}

	c := &Conversation{}
	c.Policies = policies(allowV3 | sendWhitespaceTag)

	_, _, err := c.Receive(ValidMessage("hi"))
	assertNil(t, err)

	m, err := c.Send(hello)
	assertNil(t, err)
	assertDeepEquals(t, m[0][len(hello):], expectedWhitespaceTag)
}

func Test_send_stopAppendingWhitespaceTagsIfReceivesNonDHCommitMsg(t *testing.T) {
	hello := ValidMessage("hello")
	expectedWhitespaceTag := ValidMessage{
		0x20, 0x09, 0x20, 0x20, 0x09, 0x09, 0x09, 0x09,
		0x20, 0x09, 0x20, 0x09, 0x20, 0x09, 0x20, 0x20,
		0x20, 0x20, 0x09, 0x09, 0x20, 0x20, 0x09, 0x09,
	}

	c := &Conversation{}
	c.Policies = policies(allowV3 | sendWhitespaceTag)

	m, err := c.Send(hello)
	assertNil(t, err)
	assertDeepEquals(t, m[0][len(hello):], expectedWhitespaceTag)

	_, _, err = c.Receive(ValidMessage("hi"))
	assertNil(t, err)

	m, err = c.Send(hello)
	assertNil(t, err)
	assertDeepEquals(t, m[0], hello)
}

func Test_send_dataMessageWhenItsMsgStateEncrypted(t *testing.T) {
	m := []byte("hello")
	c := bobContextAfterAKE()
	c.msgState = encrypted
	c.Policies = policies(allowV3)
	toSend, _ := c.Send(m)

	stub := bobContextAfterAKE()
	stub.msgState = encrypted
	expected, _, err := stub.createSerializedDataMessage(m, messageFlagNormal, []tlv{})

	assertDeepEquals(t, err, nil)
	assertDeepEquals(t, toSend, expected)
}

func Test_encodeWithoutFragment(t *testing.T) {
	c := newConversation(otrV2{}, fixtureRand())
	c.Policies = policies(allowV2 | allowV3 | whitespaceStartAKE)
	c.SetFragmentSize(64)

	msg := c.fragEncode([]byte("one two three"))

	expectedFragments := []ValidMessage{
		[]byte("?OTR:b25lIHR3byB0aHJlZQ==."),
	}
	assertDeepEquals(t, msg, expectedFragments)
}

func Test_encodeWithoutFragmentTooSmall(t *testing.T) {
	c := newConversation(otrV2{}, fixtureRand())
	c.Policies = policies(allowV2 | allowV3 | whitespaceStartAKE)
	c.SetFragmentSize(18)

	msg := c.fragEncode([]byte("one two three"))

	expectedFragments := []ValidMessage{
		[]byte("?OTR:b25lIHR3byB0aHJlZQ==."),
	}
	assertDeepEquals(t, msg, expectedFragments)
}

func Test_encodeWithFragment(t *testing.T) {
	c := newConversation(otrV2{}, fixtureRand())
	c.Policies = policies(allowV2 | allowV3 | whitespaceStartAKE)
	c.SetFragmentSize(22)

	msg := c.fragEncode([]byte("one two three"))

	expectedFragments := []ValidMessage{
		[]byte("?OTR,00001,00007,?OTR,"),
		[]byte("?OTR,00002,00007,:b25,"),
		[]byte("?OTR,00003,00007,lIHR,"),
		[]byte("?OTR,00004,00007,3byB,"),
		[]byte("?OTR,00005,00007,0aHJ,"),
		[]byte("?OTR,00006,00007,lZQ=,"),
		[]byte("?OTR,00007,00007,=.,"),
	}

	assertDeepEquals(t, msg, expectedFragments)
}

func Test_End_whenStateIsPlainText(t *testing.T) {
	c := newConversation(otrV2{}, fixtureRand())
	c.msgState = plainText
	msg, err := c.End()
	assertNil(t, err)
	assertNil(t, msg)
}

func Test_End_whenStateIsFinished(t *testing.T) {
	c := newConversation(otrV2{}, fixtureRand())
	c.msgState = finished
	msg, err := c.End()
	assertDeepEquals(t, c.msgState, plainText)
	assertNil(t, err)
	assertNil(t, msg)
}

func Test_End_whenStateIsEncrypted(t *testing.T) {
	bob := bobContextAfterAKE()
	bob.msgState = encrypted
	msg, err1 := bob.End()

	assertNil(t, err1)

	stub := bobContextAfterAKE()
	stub.msgState = encrypted
	expectedMsg, _, err := stub.createSerializedDataMessage(nil, messageFlagIgnoreUnreadable, []tlv{tlv{tlvType: tlvTypeDisconnected}})

	assertDeepEquals(t, err, nil)
	assertDeepEquals(t, bob.msgState, plainText)
	assertDeepEquals(t, msg, expectedMsg)
}

func Test_End_wipesSMPStateWhenGoingFromEncrypted(t *testing.T) {
	c := bobContextAfterAKE()
	c.msgState = encrypted
	c.smp.state = smpStateExpect2{}
	c.smp.secret = bnFromHex("ABCDE56321F9A9F8E364607C8C82DECD8E8E6209E2CB952C7E649620F5286FE3")
	c.smp.s1 = fixtureSmp1()
	c.smp.s2 = fixtureSmp2()
	c.smp.s3 = fixtureSmp3()
	q := "Hello"
	c.smp.question = &q

	_, e := c.End()

	assertNil(t, e)
	assertNil(t, c.smp.state)
	assertNil(t, c.smp.question)
	assertNil(t, c.smp.secret)
	assertNil(t, c.smp.s1)
	assertNil(t, c.smp.s2)
	assertNil(t, c.smp.s3)
}

func Test_End_whenStateIsEncrypted_willSignalSecurityEvent(t *testing.T) {
	bob := bobContextAfterAKE()
	bob.msgState = encrypted

	bob.expectSecurityEvent(t, func() {
		bob.End()
	}, GoneInsecure)
}

func Test_End_whenStateIsPlaintext_willNotSignalSecurityEvent(t *testing.T) {
	bob := bobContextAfterAKE()
	bob.msgState = plainText

	bob.doesntExpectSecurityEvent(t, func() {
		bob.End()
	})
}

func Test_End_whenStateIsFinished_willNotSignalSecurityEvent(t *testing.T) {
	bob := bobContextAfterAKE()
	bob.msgState = finished

	bob.doesntExpectSecurityEvent(t, func() {
		bob.End()
	})
}

func Test_End_wipesKeys(t *testing.T) {
	bob := bobContextAfterAKE()
	bob.msgState = encrypted
	bob.End()
	stub := bobContextAfterAKE()
	stub.createSerializedDataMessage(nil, messageFlagIgnoreUnreadable, []tlv{tlv{tlvType: tlvTypeDisconnected}})

	assertDeepEquals(t, dhKeyPair{}, bob.keys.ourCurrentDHKeys)
	assertDeepEquals(t, dhKeyPair{}, bob.keys.ourPreviousDHKeys)
	assertDeepEquals(t, eq(bob.keys.theirCurrentDHPubKey, big.NewInt(0)), true)
}

func Test_receive_canDecodeOTRMessagesWithoutFragments(t *testing.T) {
	c := newConversation(otrV2{}, rand.Reader)
	c.Policies.add(allowV2)

	dhCommitMsg := []byte("?OTR:AAICAAAAxPWaCOvRNycg72w2shQjcSEiYjcTh+w7rq+48UM9mpZIkpN08jtTAPcc8/9fcx9mmlVy/We+n6/G65RvobYWPoY+KD9Si41TFKku34gU4HaBbwwa7XpB/4u1gPCxY6EGe0IjthTUGK2e3qLf9YCkwJ1lm+X9kPOS/Jqu06V0qKysmbUmuynXG8T5Q8rAIRPtA/RYMqSGIvfNcZfrlJRIw6M784YtWlF3i2B6dmtjMrjH/8x5myN++Q2bxh69g6z/WX1rAFoAAAAg7Vwgf3JoiH5MdRznnS3aL66tjxQzN5qiwLtImE+KFnM=.")
	_, _, err := c.Receive(dhCommitMsg)

	assertEquals(t, err, nil)
	assertEquals(t, c.ake.state, authStateAwaitingRevealSig{})
	assertEquals(t, c.version, otrV2{})
}

func Test_receive_ignoresMessagesWithWrongInstanceTags(t *testing.T) {
	bob := newConversation(otrV3{}, rand.Reader)
	bob.Policies.add(allowV3)
	bob.ourCurrentKey = bobPrivateKey

	var msg []byte
	msg, bob.keys = fixtureDataMsg(plainDataMsg{})

	bob.ourInstanceTag = 0x1000 // different than the fixture
	bob.keys.ourKeyID = 1       //this would force key rotation
	plain, toSend, err := bob.Receive(bob.fragEncode(msg)[0])
	assertNil(t, plain)
	assertNil(t, toSend)
	assertNil(t, err)
}

func Test_receive_doesntDisplayErrorMessageToTheUser(t *testing.T) {
	msg := []byte("?OTR Error:You are wrong")
	c := &Conversation{}
	c.Policies.add(allowV3)
	plain, toSend, err := c.Receive(msg)

	assertNil(t, err)
	assertNil(t, plain)
	assertNil(t, toSend)
}

func Test_receive_doesntDisplayErrorMessageToTheUserAndStartAKE(t *testing.T) {
	msg := []byte("?OTR Error:You are wrong")
	c := &Conversation{}
	c.Policies.add(allowV3)
	c.Policies.add(errorStartAKE)
	plain, toSend, err := c.Receive(msg)

	assertEquals(t, err, nil)
	assertNil(t, plain)
	assertDeepEquals(t, toSend[0], ValidMessage("?OTRv3?"))
}

func Test_Conversation_GetTheirKey_getsTheirKey(t *testing.T) {
	c := &Conversation{theirKey: bobPrivateKey.PublicKey()}
	assertEquals(t, c.GetTheirKey(), bobPrivateKey.PublicKey())
}

func Test_Conversation_GetSSID_getsTheSSID(t *testing.T) {
	c := &Conversation{ssid: [8]byte{0xAB, 0xCD, 0xAB, 0xCD, 0xDD, 0xDD, 0xCC, 0xC0}}
	assertEquals(t, c.GetSSID(), [8]byte{0xAB, 0xCD, 0xAB, 0xCD, 0xDD, 0xDD, 0xCC, 0xC0})
}

func Test_Conversation_SetSMPEventHandler_setSMPEventHandler(t *testing.T) {
	c := &Conversation{}
	ev := CombineSMPEventHandlers()
	c.SetSMPEventHandler(ev)
	assertDeepEquals(t, c.smpEventHandler, ev)
}

func Test_Conversation_SetErrorMessageHandler_setsErrorMessageHandler(t *testing.T) {
	c := &Conversation{}
	ev := CombineErrorMessageHandlers()
	c.SetErrorMessageHandler(ev)
	assertDeepEquals(t, c.errorMessageHandler, ev)
}

func Test_Conversation_SetMessageEventHandler_setsMessageEventHandler(t *testing.T) {
	c := &Conversation{}
	ev := CombineMessageEventHandlers()
	c.SetMessageEventHandler(ev)
	assertDeepEquals(t, c.messageEventHandler, ev)
}

func Test_Conversation_SetSecurityEventHandler_setsSecurityEventHandler(t *testing.T) {
	c := &Conversation{}
	ev := CombineSecurityEventHandlers()
	c.SetSecurityEventHandler(ev)
	assertDeepEquals(t, c.securityEventHandler, ev)
}

func Test_Conversation_InitializeInstanceTag_initializesTheInstanceTag(t *testing.T) {
	c := &Conversation{}
	ret := c.InitializeInstanceTag(42)
	assertEquals(t, c.ourInstanceTag, uint32(42))
	assertEquals(t, ret, uint32(42))
}

func Test_Conversation_InitializeInstanceTag_initializesTheInstanceTagFromRandomnessIfNoneProvided(t *testing.T) {
	c := &Conversation{
		Rand: fixtureRand(),
	}
	ret := c.InitializeInstanceTag(0)
	assertEquals(t, c.ourInstanceTag, uint32(0xabcdabcd))
	assertEquals(t, ret, uint32(0xabcdabcd))
}
