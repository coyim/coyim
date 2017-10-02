package otr3

import "testing"

func Test_StartAuthenticate_failsIfWeAreNotCurrentlyEncrypted(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	c.msgState = plainText

	_, e := c.StartAuthenticate("", []byte("hello world"))
	assertEquals(t, e, errCantAuthenticateWithoutEncryption)
}

func Test_StartAuthenticate_failsIfThereIsntEnoughRandomness(t *testing.T) {
	c := bobContextAfterAKE()
	c.Rand = fixedRand([]string{"ABCD"})
	c.msgState = encrypted
	c.ssid = [8]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	c.ourCurrentKey = alicePrivateKey
	c.theirKey = bobPrivateKey.PublicKey()

	_, e := c.StartAuthenticate("", []byte("hello world"))
	assertEquals(t, e, errShortRandomRead)
}

func Test_StartAuthenticate_generatesAnSMPSecretFromTheSharedSecret(t *testing.T) {
	c := bobContextAfterAKE()
	c.msgState = encrypted
	c.ssid = [8]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	c.ourCurrentKey = alicePrivateKey
	c.theirKey = bobPrivateKey.PublicKey()

	_, e := c.StartAuthenticate("", []byte("hello world"))
	assertEquals(t, e, nil)
	assertDeepEquals(t, c.smp.secret, bnFromHex("3D7264BD983B8CA53CB365444844816F7D2453580B552EEE45CD09CA13614A5"))
}

func Test_StartAuthenticate_generatesAndReturnsTheFirstSMPMessageToSend(t *testing.T) {
	c := bobContextAfterAKE()
	c.msgState = encrypted
	c.ssid = [8]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	c.ourCurrentKey = bobPrivateKey
	c.theirKey = alicePrivateKey.PublicKey()

	msg, e := c.StartAuthenticate("", []byte("hello world"))
	assertEquals(t, e, nil)
	assertEquals(t, c.smp.state, smpStateExpect2{})

	dec, _ := c.decode(encodedMessage(msg[0]))
	_, messageBody, _ := c.parseMessageHeader(dec)
	assertDeepEquals(t, len(messageBody), 1361)
}

func Test_StartAuthenticate_generatesAndSetsTheFirstMessageOnTheConversation(t *testing.T) {
	c := bobContextAfterAKE()
	c.msgState = encrypted
	c.ssid = [8]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	c.ourCurrentKey = bobPrivateKey
	c.theirKey = alicePrivateKey.PublicKey()
	c.smp.s1 = nil

	c.StartAuthenticate("", []byte("hello world"))

	assertNotNil(t, c.smp.s1)
	assertEquals(t, c.smp.s1.msg.hasQuestion, false)
	assertEquals(t, c.smp.s1.msg.tlv().tlvType, tlvTypeSMP1)
}

func Test_StartAuthenticate_generatesAn1QMessageIfAQuestionIsGiven(t *testing.T) {
	c := bobContextAfterAKE()
	c.msgState = encrypted
	c.ssid = [8]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	c.ourCurrentKey = bobPrivateKey
	c.theirKey = alicePrivateKey.PublicKey()
	c.smp.s1 = nil

	c.StartAuthenticate("Where did we meet?", []byte("hello world"))

	assertNotNil(t, c.smp.s1)
	assertEquals(t, c.smp.s1.msg.hasQuestion, true)
	assertEquals(t, c.smp.s1.msg.question, "Where did we meet?")
	assertEquals(t, c.smp.s1.msg.tlv().tlvType, tlvTypeSMP1WithQuestion)
}

func Test_StartAuthenticate_generatesAnAbortMessageTLVIfWeAreInAnSMPStateAlready(t *testing.T) {
	c := bobContextAfterAKE()
	c.msgState = encrypted
	c.ssid = [8]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	c.ourCurrentKey = bobPrivateKey
	c.theirKey = alicePrivateKey.PublicKey()
	c.smp.s1 = nil
	c.smp.state = smpStateExpect3{}

	msg, e := c.StartAuthenticate("", []byte("hello world"))
	assertEquals(t, e, nil)
	dec, _ := c.decode(encodedMessage(msg[0]))
	_, messageBody, _ := c.parseMessageHeader(dec)
	assertDeepEquals(t, len(messageBody), 1369)
}

func Test_ProvideAuthenticationSecret_failsIfWeAreNotCurrentlyEncrypted(t *testing.T) {
	c := bobContextAfterAKE()
	c.msgState = plainText
	c.ssid = [8]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	c.ourCurrentKey = bobPrivateKey
	c.theirKey = alicePrivateKey.PublicKey()
	c.smp.state = smpStateWaitingForSecret{msg: fixtureMessage1()}

	_, e := c.ProvideAuthenticationSecret([]byte("hello world"))
	assertEquals(t, e, errCantAuthenticateWithoutEncryption)
}

func Test_ProvideAuthenticationSecret_generatesAnSMPSecretFromTheSharedSecret(t *testing.T) {
	c := bobContextAfterAKE()
	c.msgState = encrypted
	c.ssid = [8]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	c.ourCurrentKey = bobPrivateKey
	c.theirKey = alicePrivateKey.PublicKey()
	c.smp.state = smpStateWaitingForSecret{msg: fixtureMessage1()}

	_, e := c.ProvideAuthenticationSecret([]byte("hello world"))
	assertEquals(t, e, nil)
	assertDeepEquals(t, c.smp.secret, bnFromHex("3D7264BD983B8CA53CB365444844816F7D2453580B552EEE45CD09CA13614A5"))
}

func Test_ProvideAuthenticationSecret_failsAndAbortsIfWeAreNotWaitingForASecret(t *testing.T) {
	c := bobContextAfterAKE()
	c.msgState = encrypted
	c.ssid = [8]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	c.ourCurrentKey = bobPrivateKey
	c.theirKey = alicePrivateKey.PublicKey()
	c.smp.state = smpStateExpect3{}

	_, e := c.ProvideAuthenticationSecret([]byte("hello world"))
	assertEquals(t, e, errNotWaitingForSMPSecret)
}

func Test_ProvideAuthenticationSecret_continuesWithMessageProcessingIfInTheRightState(t *testing.T) {
	c := bobContextAfterAKE()
	c.msgState = encrypted
	c.ssid = [8]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	c.ourCurrentKey = bobPrivateKey
	c.theirKey = alicePrivateKey.PublicKey()
	c.smp.state = smpStateWaitingForSecret{msg: fixtureMessage1()}

	msg, e := c.ProvideAuthenticationSecret([]byte("hello world"))
	assertNil(t, e)
	dec, _ := c.decode(encodedMessage(msg[0]))
	_, messageBody, _ := c.parseMessageHeader(dec)
	assertDeepEquals(t, len(messageBody), 2181)
}

func Test_ProvideAuthenticationSecret_setsTheNextMessageState(t *testing.T) {
	c := bobContextAfterAKE()
	c.msgState = encrypted
	c.ssid = [8]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	c.ourCurrentKey = bobPrivateKey
	c.theirKey = alicePrivateKey.PublicKey()
	c.smp.state = smpStateWaitingForSecret{msg: fixtureMessage1()}

	_, e := c.ProvideAuthenticationSecret([]byte("hello world"))
	assertNil(t, e)

	assertDeepEquals(t, c.smp.state, smpStateExpect3{})
}

func Test_ProvideAuthenticationSecret_returnsFailureFromContinueSMP(t *testing.T) {
	c := bobContextAfterAKE()
	c.msgState = plainText
	c.ssid = [8]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	c.ourCurrentKey = bobPrivateKey
	c.theirKey = alicePrivateKey.PublicKey()
	c.smp.state = smpStateWaitingForSecret{msg: fixtureMessage1()}

	_, e := c.ProvideAuthenticationSecret([]byte("hello world"))
	assertEquals(t, e, errCantAuthenticateWithoutEncryption)
}

func Test_AbortAuthentication_generatesSMPAbortMessage(t *testing.T) {
	c := bobContextAfterAKE()
	c.msgState = encrypted
	c.ssid = [8]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	c.ourCurrentKey = bobPrivateKey
	c.theirKey = alicePrivateKey.PublicKey()
	c.smp.state = smpStateWaitingForSecret{msg: fixtureMessage1()}

	msgs, e := c.AbortAuthentication()

	assertDeepEquals(t, c.smp.state, smpStateExpect1{})
	assertNil(t, e)
	dec, _ := c.decode(encodedMessage(msgs[0]))
	_, messageBody, _ := c.parseMessageHeader(dec)
	assertEquals(t, len(messageBody), 505)
}

func Test_AbortAuthentication_generatesErrorWhenNoEncryptedChannelExists(t *testing.T) {
	c := bobContextAfterAKE()
	c.msgState = plainText
	c.ssid = [8]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	c.ourCurrentKey = bobPrivateKey
	c.theirKey = alicePrivateKey.PublicKey()

	_, e := c.AbortAuthentication()

	assertEquals(t, e, errCannotSendUnencrypted)
}
