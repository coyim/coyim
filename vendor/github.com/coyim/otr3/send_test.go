package otr3

import (
	"bytes"
	"testing"
)

func Test_sendDHCommit_resetsAKEKeyContext(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())

	_, err := c.sendDHCommit()

	assertNil(t, err)
	assertDeepEquals(t, c.ake.keys, keyManagementContext{})
}

func Test_Send_signalsMessageEventIfTryingToSendWithoutEncryptedChannel(t *testing.T) {
	m := []byte("hello")
	c := bobContextAfterAKE()
	c.msgState = plainText
	c.Policies = policies(allowV3 | requireEncryption)

	c.expectMessageEvent(t, func() {
		c.Send(m)
	}, MessageEventEncryptionRequired, nil, nil)
}

func Test_Send_signalsMessageEventIfTryingToSendOnAFinishedChannel(t *testing.T) {
	m := []byte("hello")
	c := bobContextAfterAKE()
	c.msgState = finished
	c.Policies = policies(allowV3 | requireEncryption)

	c.expectMessageEvent(t, func() {
		c.Send(m)
	}, MessageEventConnectionEnded, nil, nil)
}

func Test_Send_signalsEncryptionErrorMessageEventIfSomethingWentWrong(t *testing.T) {
	msg := []byte("hello")

	c := bobContextAfterAKE()
	c.msgState = encrypted
	c.Policies = policies(allowV3)
	c.keys.theirKeyID = 0

	c.expectMessageEvent(t, func() {
		c.Send(msg)
	}, MessageEventEncryptionError, nil, nil)
}

func Test_Send_callsErrorMessageHandlerAndReturnsTheResultAsAnOTRErrorMessage(t *testing.T) {
	msg := []byte("hello")

	c := bobContextAfterAKE()
	c.msgState = encrypted
	c.Policies = policies(allowV3)
	c.keys.theirKeyID = 0

	c.errorMessageHandler = dynamicErrorMessageHandler{
		func(error ErrorCode) []byte {
			if error == ErrorCodeEncryptionError {
				return []byte("snowflake happened")
			}
			return []byte("nova happened")
		}}

	msgs, _ := c.Send(msg)
	assertDeepEquals(t, msgs[0], ValidMessage("?OTR Error: snowflake happened"))
}

func Test_Send_saveLastMessageWhenMsgIsPlainTextAndEncryptedIsExpected(t *testing.T) {
	m := []byte("hello")
	c := bobContextAfterAKE()
	c.msgState = plainText
	c.Policies = policies(allowV3 | requireEncryption)

	c.Send(m)

	assertDeepEquals(t, c.resend.pending(),
		[]messageToResend{
			messageToResend{MessagePlaintext(m), nil},
		})
}

func Test_Send_saveLastMessageWhenMsgIsPlainTextAndEncryptedIsExpected_AndAddsAnOpaqueValueForEachMessage(t *testing.T) {
	m := []byte("hello")
	m2 := []byte("hello again?")
	c := bobContextAfterAKE()
	c.msgState = plainText
	c.Policies = policies(allowV3 | requireEncryption)

	c.Send(m, 42, "hello")
	c.Send(m2, 15, "something")

	assertDeepEquals(t, c.resend.pending(),
		[]messageToResend{
			messageToResend{MessagePlaintext(m), []interface{}{42, "hello"}},
			messageToResend{MessagePlaintext(m2), []interface{}{15, "something"}},
		})
}

func Test_Send_setsMayRetransmitFlagToExpectExactResending(t *testing.T) {
	m := []byte("hello")
	c := bobContextAfterAKE()
	c.msgState = plainText
	c.Policies = policies(allowV3 | requireEncryption)

	c.Send(m)

	assertEquals(t, c.resend.mayRetransmit, retransmitExact)
}

func captureStderr(f func()) string {
	originalStdErr := standardErrorOutput
	bt := bytes.NewBuffer(make([]byte, 0, 200))
	standardErrorOutput = bt

	f()

	defer func() {
		standardErrorOutput = originalStdErr
	}()

	return bt.String()
}

func Test_Send_printsDebugStatementToStderrIfGivenMagicString(t *testing.T) {
	m := []byte("hel?OTR!lo")
	c := bobContextAfterAKE()
	c.theirKey = alicePrivateKey.PublicKey()
	c.debug = true

	var ret []ValidMessage
	ss := captureStderr(func() {
		ret, _ = c.Send(m)
	})
	assertNil(t, ret)
	assertDeepEquals(t, ss, `Context:

  Our instance:   00000101
  Their instance: 00000101

  Msgstate: 0 (PLAINTEXT)

  Protocol version: 3
  OTR offer: NOT

  Auth info:
    State: 0 (NONE)
    Our keyid:   2
    Their keyid: 1
    Their fingerprint: 0BB01C360424522E94EE9C346CE877A1A4288B2F
    Proto version = 3

  SM state:
    Next expected: 0 (EXPECT1)
    Received_Q: 0
`)
}
