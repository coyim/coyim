package otr3

import "testing"

func Test_generateSMPSecretGeneratesASecret(t *testing.T) {
	aliceFingerprint := bytesFromHex("0102030405060708090A0B0C0D0E0F1011121314")
	bobFingerprint := bytesFromHex("3132333435363738393A3B3C3D3E3F4041424344")
	ssid := bytesFromHex("FFF1D1E412345668")
	secret := []byte("this is something secret")
	result := generateSMPSecret(aliceFingerprint, bobFingerprint, ssid, secret, otrV3{})
	assertDeepEquals(t, result, bnFromHex("D9B2E56321F9A9F8E364607C8C82DECD8E8E6209E2CB952C7E649620F5286FE3"))
}

func Test_SMPQuestion_returnsTheCurrentSMPQuestion(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	q := "Are all greeks liars?"
	c.smp.question = &q
	res, ok := c.SMPQuestion()
	assertEquals(t, ok, true)
	assertDeepEquals(t, res, "Are all greeks liars?")
}

func Test_SMPQuestion_returnsNotOKIfThereIsNoQuestion(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	_, ok := c.SMPQuestion()
	assertEquals(t, ok, false)
}
