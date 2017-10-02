package otr3

import (
	"crypto/rand"
	"testing"
)

func Test_processExtraSymmetricKeyTLV_signalsAReceivedKeyEventWithTheExtraKey(t *testing.T) {
	c := &Conversation{}
	extraKey := []byte{0x89, 0x11, 0x13, 0x66, 0xAB, 0xCD}
	x := dataMessageExtra{extraKey}

	called := false

	c.receivedKeyHandler = dynamicReceivedKeyHandler{func(usage uint32, usageData []byte, symkey []byte) {
		assertDeepEquals(t, symkey, extraKey)
		called = true
	}}

	c.processExtraSymmetricKeyTLV(tlv{tlvTypeExtraSymmetricKey, 0x04, []byte{0xAB, 0x12, 0xCD, 0x44}}, x)

	assertEquals(t, called, true)
}

func Test_processExtraSymmetricKeyTLV_signalsTheReceivedUsageData(t *testing.T) {
	c := &Conversation{}
	extraKey := []byte{0x89, 0x11, 0x13, 0x66, 0xAB, 0xCD}
	x := dataMessageExtra{extraKey}

	called := false

	c.receivedKeyHandler = dynamicReceivedKeyHandler{func(usage uint32, usageData []byte, symkey []byte) {
		assertEquals(t, usage, uint32(0xAB12CD44))
		called = true
	}}

	c.processExtraSymmetricKeyTLV(tlv{tlvTypeExtraSymmetricKey, 0x04, []byte{0xAB, 0x12, 0xCD, 0x44}}, x)

	assertEquals(t, called, true)
}

func Test_processExtraSymmetricKeyTLV_doesntSignalAnythingIfThereIsNoUsageData(t *testing.T) {
	c := &Conversation{}
	extraKey := []byte{0x89, 0x11, 0x13, 0x66, 0xAB, 0xCD}
	x := dataMessageExtra{extraKey}

	c.receivedKeyHandler = dynamicReceivedKeyHandler{func(usage uint32, usageData []byte, symkey []byte) {
		t.Errorf("Didn't expect a received key event one")
	}}

	c.processExtraSymmetricKeyTLV(tlv{tlvTypeExtraSymmetricKey, 0x00, []byte{}}, x)
}

func Test_processExtraSymmetricKeyTLV_providesExtraUsageDataIfGiven(t *testing.T) {
	c := &Conversation{}
	extraKey := []byte{0x89, 0x11, 0x13, 0x66, 0xAB, 0xCD}
	x := dataMessageExtra{extraKey}

	called := false

	c.receivedKeyHandler = dynamicReceivedKeyHandler{func(usage uint32, usageData []byte, symkey []byte) {
		assertDeepEquals(t, usageData, []byte{0x01, 0x02})
		called = true
	}}

	c.processExtraSymmetricKeyTLV(tlv{tlvTypeExtraSymmetricKey, 0x06, []byte{0xAB, 0x12, 0xCD, 0x44, 0x01, 0x02, 0x04}}, x)

	assertEquals(t, called, true)
}

func Test_processExtraSymmetricKeyTLV_alwaysReturnsNilAndNil(t *testing.T) {
	c := &Conversation{}
	x := dataMessageExtra{[]byte{0x89, 0x11, 0x13, 0x66, 0xAB, 0xCD}}

	c.receivedKeyHandler = dynamicReceivedKeyHandler{func(usage uint32, usageData []byte, symkey []byte) {
	}}

	res, err := c.processExtraSymmetricKeyTLV(tlv{tlvTypeExtraSymmetricKey, 0x06, []byte{0xAB, 0x12, 0xCD, 0x44, 0x01, 0x02, 0x04}}, x)

	assertNil(t, res)
	assertNil(t, err)
}

func Test_UseExtraSymmetricKey_returnsErrorIfWeAreNotInEncryptedMode(t *testing.T) {
	c := aliceContextAfterAKE()
	c.msgState = plainText

	_, _, err := c.UseExtraSymmetricKey(0, nil)
	assertDeepEquals(t, err, newOtrError("cannot send message in current state"))
}

func Test_UseExtraSymmetricKey_returnsErrorIfTheirKeyIDIsZero(t *testing.T) {
	c := aliceContextAfterAKE()
	c.msgState = encrypted
	c.keys.theirKeyID = 0

	_, _, err := c.UseExtraSymmetricKey(0, nil)
	assertDeepEquals(t, err, newOtrError("cannot send message in current state"))
}

func Test_UseExtraSymmetricKey_generatesADataMessageWithTheDataProvided(t *testing.T) {
	c := newConversation(otrV3{}, rand.Reader)
	c.Policies.add(allowV3)
	c.ourCurrentKey = bobPrivateKey

	_, c.keys = fixtureDataMsg(plainDataMsg{message: []byte("something")})
	c.msgState = encrypted

	_, msg, err := c.UseExtraSymmetricKey(0x1234, []byte{0xAB, 0xCD, 0xEE})
	decodedMsg, _ := c.decode(encodedMessage(msg[0]))
	assertNil(t, err)
	_, exp, e := fixtureDecryptDataMsgBase(decodedMsg)
	assertNil(t, e)

	assertEquals(t, len(exp.tlvs), 2)
	assertDeepEquals(t, exp.tlvs[0].tlvType, uint16(tlvTypeExtraSymmetricKey))
	assertDeepEquals(t, exp.tlvs[0].tlvLength, uint16(7))
	assertDeepEquals(t, exp.tlvs[0].tlvValue, []byte{0x00, 0x00, 0x12, 0x34, 0xAB, 0xCD, 0xEE})
}

func Test_UseExtraSymmetricKey_generatesADataMessageWithIgnoreUnreadableSet(t *testing.T) {
	c := newConversation(otrV3{}, rand.Reader)
	c.Policies.add(allowV3)
	c.ourCurrentKey = bobPrivateKey

	_, c.keys = fixtureDataMsg(plainDataMsg{message: []byte("something")})
	c.msgState = encrypted
	_, msg, _ := c.UseExtraSymmetricKey(0x1234, []byte{0xAB, 0xCD, 0xEE})
	decodedMsg, _ := c.decode(encodedMessage(msg[0]))
	assertEquals(t, decodedMsg[11], messageFlagIgnoreUnreadable)
}

func Test_UseExtraSymmetricKey_returnsTheGeneratedSymmetricKey(t *testing.T) {
	c := newConversation(otrV3{}, rand.Reader)
	c.Policies.add(allowV3)
	c.ourCurrentKey = bobPrivateKey

	_, c.keys = fixtureDataMsg(plainDataMsg{message: []byte("something")})
	c.msgState = encrypted
	k, _, _ := c.UseExtraSymmetricKey(0x1234, []byte{0xAB, 0xCD, 0xEE})
	assertDeepEquals(t, k, bytesFromHex("0e1810c7c62c3bace6450dcbef16af8a271b5ac93030b83e9d0d80e0641e3c18"))
}
