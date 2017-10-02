package otr3

import (
	"math/big"
	"testing"
)

func Test_smpStateExpect1_goToWaitingForSecretWhenReceivesSmpMessage1(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	msg := fixtureMessage1()
	nextState, _, _ := smpStateExpect1{}.receiveMessage1(c, msg)

	assertDeepEquals(t, nextState, smpStateWaitingForSecret{msg: msg})
}

func Test_smpStateExpect1_willSendANotificationThatASecretIsNeeded(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	c.expectSMPEvent(t, func() {
		smpStateExpect1{}.receiveMessage1(c, fixtureMessage1())
	}, SMPEventAskForSecret, 25, "")
}

func Test_smpStateExpect1_willSendANotificationThatAnAnswerIsNeededIfQuestionProvided(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	msg := fixtureMessage1()
	msg.hasQuestion = true
	msg.question = "What do you think?"

	c.expectSMPEvent(t, func() {
		smpStateExpect1{}.receiveMessage1(c, msg)
	}, SMPEventAskForAnswer, 25, "What do you think?")
}

func Test_smpStateWaitingForSecret_goToExpectState3WhenReceivesContinueSmpMessage1(t *testing.T) {
	c := bobContextAfterAKE()
	c.msgState = encrypted
	c.ssid = [8]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	c.ourCurrentKey = bobPrivateKey
	c.theirKey = alicePrivateKey.PublicKey()
	c.smp.state = smpStateWaitingForSecret{msg: fixtureMessage1()}

	msg := fixtureMessage1()
	nextState, _, err := smpStateWaitingForSecret{msg: msg}.continueMessage1(c, []byte{})

	assertNil(t, err)
	assertNotNil(t, c.smp.s2)
	assertEquals(t, nextState, smpStateExpect3{})
}

func Test_smpStateExpect1_receiveMessage1_setsTheSMPQuestionIfThereWasOneInTheMessage(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	c.smp.secret = bnFromHex("ABCDE56321F9A9F8E364607C8C82DECD8E8E6209E2CB952C7E649620F5286FE3")
	msg := fixtureMessage1Q()

	smpStateExpect1{}.receiveMessage1(c, msg)
	v, ok := c.SMPQuestion()

	assertDeepEquals(t, ok, true)
	assertDeepEquals(t, v, "What's the clue?")
}

func Test_smpStateExpect1_returnsSmpMessageAbortIfReceivesUnexpectedMessage(t *testing.T) {
	state := smpStateExpect1{}
	c := newConversation(otrV3{}, fixtureRand())
	_, msg, err := state.receiveMessage2(c, smp2Message{})
	assertEquals(t, err, nil)
	assertDeepEquals(t, msg, smpMessageAbort{})

	_, msg, err = state.receiveMessage3(c, smp3Message{})
	assertEquals(t, err, nil)
	assertDeepEquals(t, msg, smpMessageAbort{})

	_, msg, err = state.receiveMessage4(c, smp4Message{})
	assertEquals(t, err, nil)
	assertDeepEquals(t, msg, smpMessageAbort{})
}

func Test_smpStateExpect1_givesAnErrorNotificationIfTheWrongMessageIsSent(t *testing.T) {
	state := smpStateExpect1{}
	c := newConversation(otrV3{}, fixtureRand())

	c.expectSMPEvent(t, func() {
		state.receiveMessage2(c, smp2Message{})
	}, SMPEventError, 0, "")

	c.expectSMPEvent(t, func() {
		state.receiveMessage3(c, smp3Message{})
	}, SMPEventError, 0, "")

	c.expectSMPEvent(t, func() {
		state.receiveMessage4(c, smp4Message{})
	}, SMPEventError, 0, "")
}

func Test_smpStateExpect2_givesAnErrorNotificationIfTheWrongMessageIsSent(t *testing.T) {
	state := smpStateExpect2{}
	c := newConversation(otrV3{}, fixtureRand())

	c.expectSMPEvent(t, func() {
		state.receiveMessage1(c, smp1Message{})
	}, SMPEventError, 0, "")

	c.expectSMPEvent(t, func() {
		state.receiveMessage3(c, smp3Message{})
	}, SMPEventError, 0, "")

	c.expectSMPEvent(t, func() {
		state.receiveMessage4(c, smp4Message{})
	}, SMPEventError, 0, "")
}

func Test_smpStateExpect3_givesAnErrorNotificationIfTheWrongMessageIsSent(t *testing.T) {
	state := smpStateExpect3{}
	c := newConversation(otrV3{}, fixtureRand())

	c.expectSMPEvent(t, func() {
		state.receiveMessage1(c, smp1Message{})
	}, SMPEventError, 0, "")

	c.expectSMPEvent(t, func() {
		state.receiveMessage2(c, smp2Message{})
	}, SMPEventError, 0, "")

	c.expectSMPEvent(t, func() {
		state.receiveMessage4(c, smp4Message{})
	}, SMPEventError, 0, "")
}

func Test_smpStateExpect4_givesAnErrorNotificationIfTheWrongMessageIsSent(t *testing.T) {
	state := smpStateExpect4{}
	c := newConversation(otrV3{}, fixtureRand())

	c.expectSMPEvent(t, func() {
		state.receiveMessage1(c, smp1Message{})
	}, SMPEventError, 0, "")

	c.expectSMPEvent(t, func() {
		state.receiveMessage2(c, smp2Message{})
	}, SMPEventError, 0, "")

	c.expectSMPEvent(t, func() {
		state.receiveMessage3(c, smp3Message{})
	}, SMPEventError, 0, "")
}

func Test_smpStateExpect2_goToExpectState4WhenReceivesSmpMessage2(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	c.smp.secret = bnFromHex("ABCDE56321F9A9F8E364607C8C82DECD8E8E6209E2CB952C7E649620F5286FE3")
	c.smp.s1 = fixtureSmp1()

	msg := fixtureMessage2()
	nextState, _, err := smpStateExpect2{}.receiveMessage2(c, msg)

	assertNil(t, err)
	assertNotNil(t, c.smp.s3)
	assertEquals(t, nextState, smpStateExpect4{})
}

func Test_smpStateExpect2_sendsAnSMPEventAboutSMPProgressHere(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	c.smp.secret = bnFromHex("ABCDE56321F9A9F8E364607C8C82DECD8E8E6209E2CB952C7E649620F5286FE3")
	c.smp.s1 = fixtureSmp1()

	c.expectSMPEvent(t, func() {
		smpStateExpect2{}.receiveMessage2(c, fixtureMessage2())
	}, SMPEventInProgress, 60, "")
}

func Test_smpStateExpect2_returnsSmpMessageAbortIfReceivesUnexpectedMessage(t *testing.T) {
	state := smpStateExpect2{}
	c := newConversation(otrV3{}, fixtureRand())
	_, msg, err := state.receiveMessage1(c, smp1Message{})
	assertEquals(t, err, nil)
	assertDeepEquals(t, msg, smpMessageAbort{})

	_, msg, err = state.receiveMessage3(c, smp3Message{})
	assertEquals(t, err, nil)
	assertDeepEquals(t, msg, smpMessageAbort{})

	_, msg, err = state.receiveMessage4(c, smp4Message{})
	assertEquals(t, err, nil)
	assertDeepEquals(t, msg, smpMessageAbort{})
}

func Test_smpStateExpect3_goToExpectState1WhenReceivesSmpMessage3(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	c.smp.secret = bnFromHex("ABCDE56321F9A9F8E364607C8C82DECD8E8E6209E2CB952C7E649620F5286FE3")
	c.smp.s2 = fixtureSmp2()
	msg := fixtureMessage3()

	nextState, _, _ := smpStateExpect3{}.receiveMessage3(c, msg)

	assertEquals(t, nextState, smpStateExpect1{})
}

func Test_smpStateExpect3_wipesSMPWhenReceivesSmpMessage3(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	c.smp.s2 = fixtureSmp2()
	msg := fixtureMessage3()

	nextState, _, _ := smpStateExpect3{}.receiveMessage3(c, msg)

	assertEquals(t, nextState, smpStateExpect1{})
	assertNil(t, c.smp.secret)
	assertNil(t, c.smp.s1)
	assertNil(t, c.smp.s2)
	assertNil(t, c.smp.s3)
}

func Test_smpStateExpect3_willSendAnSMPNotificationOnProtocolSuccess(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	c.smp.secret = bnFromHex("ABCDE56321F9A9F8E364607C8C82DECD8E8E6209E2CB952C7E649620F5286FE3")
	c.smp.s2 = fixtureSmp2()

	c.expectSMPEvent(t, func() {
		smpStateExpect3{}.receiveMessage3(c, fixtureMessage3())
	}, SMPEventSuccess, 100, "")
}

func Test_smpStateExpect3_returnsSmpMessageAbortIfReceivesUnexpectedMessage(t *testing.T) {
	state := smpStateExpect3{}
	c := newConversation(otrV3{}, fixtureRand())
	_, msg, err := state.receiveMessage1(c, smp1Message{})
	assertEquals(t, err, nil)
	assertDeepEquals(t, msg, smpMessageAbort{})

	_, msg, err = state.receiveMessage2(c, smp2Message{})
	assertEquals(t, err, nil)
	assertDeepEquals(t, msg, smpMessageAbort{})

	_, msg, err = state.receiveMessage4(c, smp4Message{})
	assertEquals(t, err, nil)
	assertDeepEquals(t, msg, smpMessageAbort{})
}

func Test_smpStateExpect4_goToExpectState1WhenReceivesSmpMessage4(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	c.smp.s1 = fixtureSmp1()
	c.smp.s3 = fixtureSmp3()
	msg := fixtureMessage4()

	nextState, _, _ := smpStateExpect4{}.receiveMessage4(c, msg)

	assertEquals(t, nextState, smpStateExpect1{})
}

func Test_smpStateExpect4_willSendAnSMPNotificationOnProtocolSuccess(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	c.smp.s1 = fixtureSmp1()
	c.smp.s3 = fixtureSmp3()

	c.expectSMPEvent(t, func() {
		smpStateExpect4{}.receiveMessage4(c, fixtureMessage4())
	}, SMPEventSuccess, 100, "")
}

func Test_smpStateExpect4_wipesSMPWhenReceivesSmpMessage4(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	c.smp.s1 = fixtureSmp1()
	c.smp.s3 = fixtureSmp3()
	msg := fixtureMessage4()

	nextState, _, _ := smpStateExpect4{}.receiveMessage4(c, msg)

	assertEquals(t, nextState, smpStateExpect1{})
	assertNil(t, c.smp.secret)
	assertNil(t, c.smp.s1)
	assertNil(t, c.smp.s2)
	assertNil(t, c.smp.s3)
}

func Test_smpStateExpect4_returnsSmpMessageAbortIfReceivesUnexpectedMessage(t *testing.T) {
	state := smpStateExpect4{}
	c := newConversation(otrV3{}, fixtureRand())
	_, msg, err := state.receiveMessage1(c, smp1Message{})
	assertEquals(t, err, nil)
	assertDeepEquals(t, msg, smpMessageAbort{})

	_, msg, err = state.receiveMessage2(c, smp2Message{})
	assertEquals(t, err, nil)
	assertDeepEquals(t, msg, smpMessageAbort{})

	_, msg, err = state.receiveMessage3(c, smp3Message{})
	assertEquals(t, err, nil)
	assertDeepEquals(t, msg, smpMessageAbort{})
}

func Test_contextTransitionsFromSmpExpect1ToSmpWaitingForSecret(t *testing.T) {
	m := fixtureMessage1()
	c := newConversation(otrV3{}, fixtureRand())
	c.smp.secret = bnFromHex("ABCDE56321F9A9F8E364607C8C82DECD8E8E6209E2CB952C7E649620F5286FE3")

	c.receiveSMP(m)
	assertDeepEquals(t, c.smp.state, smpStateWaitingForSecret{msg: m})
}

func Test_contextTransitionsFromSmpExpect2ToSmpExpect4(t *testing.T) {
	m := fixtureMessage2()
	c := newConversation(otrV3{}, fixtureRand())
	c.smp.state = smpStateExpect2{}
	c.smp.s1 = fixtureSmp1()
	c.smp.secret = bnFromHex("ABCDE56321F9A9F8E364607C8C82DECD8E8E6209E2CB952C7E649620F5286FE3")

	c.receiveSMP(m)
	assertEquals(t, c.smp.state, smpStateExpect4{})
}

func Test_contextTransitionsFromSmpExpect3ToSmpExpect1(t *testing.T) {
	m := fixtureMessage3()
	c := newConversation(otrV3{}, fixtureRand())
	c.smp.state = smpStateExpect3{}
	c.smp.s2 = fixtureSmp2()
	c.smp.secret = bnFromHex("ABCDE56321F9A9F8E364607C8C82DECD8E8E6209E2CB952C7E649620F5286FE3")

	c.receiveSMP(m)
	assertEquals(t, c.smp.state, smpStateExpect1{})
}

func Test_contextTransitionsFromSmpExpect4ToSmpExpect1(t *testing.T) {
	m := fixtureMessage4()
	c := newConversation(otrV3{}, fixtureRand())
	c.smp.state = smpStateExpect4{}
	c.smp.s1 = fixtureSmp1()
	c.smp.s3 = fixtureSmp3()

	c.receiveSMP(m)
	assertEquals(t, c.smp.state, smpStateExpect1{})
}

func Test_contextUnexpectedMessageTransitionsToSmpExpected1(t *testing.T) {
	m := fixtureMessage1()

	c := newConversation(otrV3{}, fixtureRand())
	c.smp.state = smpStateExpect3{}
	toSend, err := c.receiveSMP(m)

	assertNil(t, err)
	assertEquals(t, c.smp.state, smpStateExpect1{})
	assertDeepEquals(t, *toSend, smpMessageAbort{}.tlv())
}

func Test_smpStateExpect1_receiveMessage1_abortsSMPIfVerifySMP1ReturnsError(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())

	s, m, err := smpStateExpect1{}.receiveMessage1(c, smp1Message{g2a: big.NewInt(1)})

	assertNil(t, err)
	assertEquals(t, s, smpStateExpect1{})
	assertDeepEquals(t, m, smpMessageAbort{})
}

func Test_smpStateExpect1_receiveMessage1_signalsCheatingIfVerifySMP1Fails(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())

	c.expectSMPEvent(t, func() {
		smpStateExpect1{}.receiveMessage1(c, smp1Message{g2a: big.NewInt(1)})
	}, SMPEventCheated, 0, "")
}

func Test_smp1Message_receivedMessage_abortsSMPIfFailsToVerifyMessage1(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	c.smp.state = smpStateExpect1{}
	m := smp1Message{g2a: big.NewInt(1)}
	ret, err := m.receivedMessage(c)

	assertNil(t, err)
	assertDeepEquals(t, ret, smpMessageAbort{})
}

func Test_smpStateWaitingForSecret_continueMessage1_abortsSMPIfgenerateSMP2Fails(t *testing.T) {
	c := bobContextAfterAKE()
	c.Rand = fixedRand([]string{"ABCD"})
	c.msgState = encrypted
	c.ssid = [8]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	c.ourCurrentKey = bobPrivateKey
	c.theirKey = alicePrivateKey.PublicKey()
	c.smp.state = smpStateWaitingForSecret{msg: fixtureMessage1()}

	s, m, err := smpStateWaitingForSecret{msg: fixtureMessage1()}.continueMessage1(c, []byte("hello world"))

	assertNil(t, err)
	assertEquals(t, s, smpStateExpect1{})
	assertEquals(t, m, smpMessageAbort{})
}

func Test_smpStateExpect2_receiveMessage2_abortsSMPIfVerifySMPReturnsError(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	c.smp.s1 = fixtureSmp1()
	c.smp.secret = bnFromHex("ABCDE56321F9A9F8E364607C8C82DECD8E8E6209E2CB952C7E649620F5286FE3")
	s, m, err := smpStateExpect2{}.receiveMessage2(c, smp2Message{g2b: big.NewInt(1)})

	assertNil(t, err)
	assertEquals(t, s, smpStateExpect1{})
	assertDeepEquals(t, m, smpMessageAbort{})
}

func Test_smp2Message_receivedMessage_abortsSMPIfUnderlyingPrimitiveHasErrors(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	c.smp.state = smpStateExpect2{}
	c.smp.s1 = fixtureSmp1()
	c.smp.secret = bnFromHex("ABCDE56321F9A9F8E364607C8C82DECD8E8E6209E2CB952C7E649620F5286FE3")
	ret, err := smp2Message{g2b: big.NewInt(1)}.receivedMessage(c)

	assertNil(t, err)
	assertDeepEquals(t, ret, smpMessageAbort{})
}

func Test_smpStateExpect2_receiveMessage2_abortsSMPIfgenerateSMPFails(t *testing.T) {
	c := newConversation(otrV3{}, fixedRand([]string{"ABCD"}))
	c.smp.s1 = fixtureSmp1()
	c.smp.secret = bnFromHex("ABCDE56321F9A9F8E364607C8C82DECD8E8E6209E2CB952C7E649620F5286FE3")
	s, m, err := smpStateExpect2{}.receiveMessage2(c, fixtureMessage2())

	assertNil(t, err)
	assertEquals(t, s, smpStateExpect1{})
	assertDeepEquals(t, m, smpMessageAbort{})
}

func Test_smpStateExpect3_receiveMessage3_abortsSMPIfVerifySMPReturnsError(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	c.smp.secret = bnFromHex("ABCDE56321F9A9F8E364607C8C82DECD8E8E6209E2CB952C7E649620F5286FE3")
	c.smp.s2 = fixtureSmp2()
	s, m, err := smpStateExpect3{}.receiveMessage3(c, smp3Message{pa: big.NewInt(1)})

	assertNil(t, err)
	assertEquals(t, s, smpStateExpect1{})
	assertDeepEquals(t, m, smpMessageAbort{})
}

func Test_smp3Message_receivedMessage_abortsSMPIfUnderlyingPrimitiveDoes(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	c.smp.state = smpStateExpect3{}
	c.smp.secret = bnFromHex("ABCDE56321F9A9F8E364607C8C82DECD8E8E6209E2CB952C7E649620F5286FE3")
	c.smp.s2 = fixtureSmp2()
	ret, err := smp3Message{pa: big.NewInt(1)}.receivedMessage(c)

	assertNil(t, err)
	assertDeepEquals(t, ret, smpMessageAbort{})
}

func Test_smpStateExpect3_receiveMessage3_abortsSMPIfProtocolFails(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	c.smp.secret = bnFromHex("ABCDE56321F9A9F8E364607C8C82DECD8E8E6209E2CB952C7E649620F5286FE3")
	c.smp.s2 = fixtureSmp2()
	c.smp.s2.b3 = sub(c.smp.s2.b3, big.NewInt(1))
	s, m, err := smpStateExpect3{}.receiveMessage3(c, fixtureMessage3())

	assertNil(t, err)
	assertEquals(t, s, smpStateExpect1{})
	assertEquals(t, m, smpMessageAbort{})
}

func Test_smpStateExpect3_receiveMessage3_willSendAnSMPNotificationOnProtocolFailure(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	c.smp.secret = bnFromHex("ABCDE56321F9A9F8E364607C8C82DECD8E8E6209E2CB952C7E649620F5286FE3")
	c.smp.s2 = fixtureSmp2()
	c.smp.s2.b3 = sub(c.smp.s2.b3, big.NewInt(1))

	c.expectSMPEvent(t, func() {
		smpStateExpect3{}.receiveMessage3(c, fixtureMessage3())
	}, SMPEventFailure, 100, "")

}

func Test_smpStateExpect3_receiveMessage3_abortsSMPIfCantGenerateFinalParameters(t *testing.T) {
	c := newConversation(otrV3{}, fixedRand([]string{"ABCD"}))
	c.smp.secret = bnFromHex("ABCDE56321F9A9F8E364607C8C82DECD8E8E6209E2CB952C7E649620F5286FE3")
	c.smp.s2 = fixtureSmp2()
	s, m, err := smpStateExpect3{}.receiveMessage3(c, fixtureMessage3())

	assertNil(t, err)
	assertEquals(t, s, smpStateExpect1{})
	assertDeepEquals(t, m, smpMessageAbort{})
}

func Test_smpStateExpect4_receiveMessage4_abortsSMPIfVerifySMPReturnsError(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	c.smp.s1 = fixtureSmp1()
	c.smp.s3 = fixtureSmp3()
	s, m, err := smpStateExpect4{}.receiveMessage4(c, smp4Message{rb: big.NewInt(1)})

	assertNil(t, err)
	assertEquals(t, s, smpStateExpect1{})
	assertDeepEquals(t, m, smpMessageAbort{})
}

func Test_smpStateExpect4_receiveMessage4_abortsSMPIfProtocolFails(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	c.smp.s1 = fixtureSmp1()
	c.smp.s3 = fixtureSmp3()
	c.smp.s3.papb = sub(c.smp.s3.papb, big.NewInt(1))
	s, m, err := smpStateExpect4{}.receiveMessage4(c, fixtureMessage4())

	assertNil(t, err)
	assertEquals(t, s, smpStateExpect1{})
	assertEquals(t, m, smpMessageAbort{})
}

func Test_smpStateExpect4_receiveMessage4_willSendAnSMPNotificationOnProtocolFailure(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	c.smp.s1 = fixtureSmp1()
	c.smp.s3 = fixtureSmp3()
	c.smp.s3.papb = sub(c.smp.s3.papb, big.NewInt(1))

	c.expectSMPEvent(t, func() {
		smpStateExpect4{}.receiveMessage4(c, fixtureMessage4())
	}, SMPEventFailure, 100, "")
}

func Test_smp4Message_receivedMessage_abortsSMPIfTheUnderlyingPrimitiveDoes(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	c.smp.state = smpStateExpect4{}
	c.smp.s1 = fixtureSmp1()
	c.smp.s3 = fixtureSmp3()

	ret, err := smp4Message{rb: big.NewInt(1)}.receivedMessage(c)
	assertNil(t, err)
	assertDeepEquals(t, ret, smpMessageAbort{})
}

func Test_receive_returnsAnyErrorThatOccurs(t *testing.T) {
	m := fixtureMessage2()
	c := newConversation(otrV3{}, fixedRand([]string{"ABCD"}))
	c.smp.s1 = fixtureSmp1()
	c.smp.state = smpStateExpect2{}
	c.smp.secret = bnFromHex("ABCDE56321F9A9F8E364607C8C82DECD8E8E6209E2CB952C7E649620F5286FE3")

	ret, err := c.receiveSMP(m)
	assertNil(t, err)
	assertDeepEquals(t, *ret, smpMessageAbort{}.tlv())
}

func Test_smpStateExpect1_String_returnsTheCorrectString(t *testing.T) {
	assertEquals(t, smpStateExpect1{}.String(), "SMPSTATE_EXPECT1")
}

func Test_smpStateExpect2_String_returnsTheCorrectString(t *testing.T) {
	assertEquals(t, smpStateExpect2{}.String(), "SMPSTATE_EXPECT2")
}

func Test_smpStateExpect3_String_returnsTheCorrectString(t *testing.T) {
	assertEquals(t, smpStateExpect3{}.String(), "SMPSTATE_EXPECT3")
}

func Test_smpStateExpect4_String_returnsTheCorrectString(t *testing.T) {
	assertEquals(t, smpStateExpect4{}.String(), "SMPSTATE_EXPECT4")
}

func Test_smpMessageAbort_receivedMessage_setsTheNewState(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	c.smp.state = smpStateExpect2{}
	ret, err := smpMessageAbort{}.receivedMessage(c)
	assertDeepEquals(t, ret, nil)
	assertDeepEquals(t, err, nil)
	assertDeepEquals(t, c.smp.state, smpStateExpect1{})
}

func Test_smpMessageAbort_receivedMessage_sendsAnSMPEventAboutTheAbort(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	c.smp.state = smpStateExpect2{}

	c.expectSMPEvent(t, func() {
		smpMessageAbort{}.receivedMessage(c)
	}, SMPEventAbort, 0, "")
}

func Test_restartSMP_createsSMPAbort(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	c.smp.state = smpStateExpect2{}

	assertDeepEquals(t, c.restartSMP(), smpMessageAbort{}.tlv())
	assertDeepEquals(t, c.smp.state, smpStateExpect1{})
}
