package otr3

import (
	"math/big"
	"testing"
)

func Test_processDisconnectedTLV_forgetAllKeysAndTransitionToFinished(t *testing.T) {
	c := &Conversation{}
	c.msgState = encrypted
	c.keys.theirCurrentDHPubKey = big.NewInt(99)

	c.processDisconnectedTLV(tlv{}, dataMessageExtra{})

	assertEquals(t, c.msgState, finished)
	assertDeepEquals(t, c.keys, keyManagementContext{})
}

func Test_processDisconnectedTLV_signalsASecurityEvent(t *testing.T) {
	c := &Conversation{}
	c.msgState = encrypted
	c.keys.theirCurrentDHPubKey = big.NewInt(99)

	c.expectSecurityEvent(t, func() {
		c.processDisconnectedTLV(tlv{}, dataMessageExtra{})
	}, GoneInsecure)
}

func Test_processDisconnectedTLV_isActuallyInsecureWhenTheEventIsSignalled(t *testing.T) {
	c := &Conversation{}
	c.msgState = encrypted
	c.keys.theirCurrentDHPubKey = big.NewInt(99)

	c.securityEventHandler = dynamicSecurityEventHandler{func(event SecurityEvent) {
		assertEquals(t, c.msgState, finished)
	}}

	c.processDisconnectedTLV(tlv{}, dataMessageExtra{})
}

func Test_processDisconnectedTLV_doesntSignalsASecurityEventIfWeWereInPlaintext(t *testing.T) {
	c := &Conversation{}
	c.msgState = plainText
	c.keys.theirCurrentDHPubKey = big.NewInt(99)

	c.doesntExpectSecurityEvent(t, func() {
		c.processDisconnectedTLV(tlv{}, dataMessageExtra{})
	})
}

func Test_processDisconnectedTLV_doesntSignalsASecurityEventIfWeAreInFinished(t *testing.T) {
	c := &Conversation{}
	c.msgState = finished
	c.keys.theirCurrentDHPubKey = big.NewInt(99)

	c.doesntExpectSecurityEvent(t, func() {
		c.processDisconnectedTLV(tlv{}, dataMessageExtra{})
	})
}

func Test_processDisconnectedTLV_wipesSMPState(t *testing.T) {
	c := &Conversation{}
	c.smp.state = smpStateExpect2{}
	c.smp.secret = bnFromHex("ABCDE56321F9A9F8E364607C8C82DECD8E8E6209E2CB952C7E649620F5286FE3")
	c.smp.s1 = fixtureSmp1()
	c.smp.s2 = fixtureSmp2()
	c.smp.s3 = fixtureSmp3()
	q := "Hello"
	c.smp.question = &q

	c.processDisconnectedTLV(tlv{}, dataMessageExtra{})

	assertNil(t, c.smp.state)
	assertNil(t, c.smp.question)
	assertNil(t, c.smp.secret)
	assertNil(t, c.smp.s1)
	assertNil(t, c.smp.s2)
	assertNil(t, c.smp.s3)
}
