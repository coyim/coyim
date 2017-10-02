package otr3

import (
	"crypto/rand"
	"testing"
)

func TestFullSMPHandshake(t *testing.T) {
	secret := bnFromHex("ABCDE56321F9A9F8E364607C8C82DECD8E8E6209E2CB952C7E649620F5286FE3")
	alice := newConversation(otrV3{}, rand.Reader)
	bob := newConversation(otrV3{}, rand.Reader)

	// Alice -> Bob
	// Stores: x, a2, and a3
	// Sends: g2a, c2, D2, g3a, c3 and D3
	s1, _ := alice.generateSMP1()

	//Bob
	err := bob.verifySMP1(s1.msg)
	assertDeepEquals(t, err, nil)

	// Bob -> Alice
	// Stores: g3a, g2, g3, b3, Pb and Qb
	// Sends: g2b, c2, D2, g3b, c3, D3, Pb, Qb, cP, D5 and D6
	s2, _ := bob.generateSMP2(secret, s1.msg)

	// Alice
	err = alice.verifySMP2(&s1, s2.msg)
	assertDeepEquals(t, err, nil)

	// Alice -> Bob
	// Stores: g3b, (Pa / Pb), (Qa / Qb) and Ra
	// Sends: Pa, Qa, cP, D5, D6, Ra, cR and D7
	s3, _ := alice.generateSMP3(secret, s1, s2.msg)

	// Bob
	err = bob.verifySMP3(&s2, s3.msg)
	assertDeepEquals(t, err, nil)

	err = bob.verifySMP3ProtocolSuccess(&s2, s3.msg)
	assertDeepEquals(t, err, nil)

	// Bob -> Alice
	// Stores: ???
	// Sends: Rb, cR and D7
	s4, _ := bob.generateSMP4(secret, s2, s3.msg)

	// Alice
	err = alice.verifySMP4(&s3, s4.msg)
	assertDeepEquals(t, err, nil)

	err = alice.verifySMP4ProtocolSuccess(&s1, &s3, s4.msg)
	assertDeepEquals(t, err, nil)
}

func Test_SMP_Full(t *testing.T) {
	alice := &Conversation{Rand: rand.Reader}
	alice.ourKeys = []PrivateKey{alicePrivateKey}
	alice.Policies = policies(allowV3)

	bob := &Conversation{Rand: rand.Reader}
	bob.ourKeys = []PrivateKey{bobPrivateKey}
	bob.Policies = policies(allowV3)

	var err error
	var aliceMessages []ValidMessage
	var bobMessages []ValidMessage

	aliceMessages = append(bobMessages, alice.QueryMessage())

	for len(aliceMessages)+len(bobMessages) > 0 {
		bobMessages = nil
		for _, m := range aliceMessages {
			_, bobMessages, err = bob.Receive(m)
			assertNil(t, err)
		}

		aliceMessages = nil
		for _, m := range bobMessages {
			_, aliceMessages, err = alice.Receive(m)
			assertNil(t, err)
		}
	}

	assertEquals(t, bob.IsEncrypted(), true)
	assertEquals(t, alice.IsEncrypted(), true)

	bobMessages, err = bob.StartAuthenticate("", []byte("secret"))
	assertNil(t, err)
	assertEquals(t, bob.smp.state, smpStateExpect2{})

	_, aliceMessages, err = alice.Receive(bobMessages[0])
	assertNil(t, err)

	// this is an internal state
	_, ok := alice.smp.state.(smpStateWaitingForSecret)
	assertEquals(t, ok, true)

	aliceMessages, err = alice.ProvideAuthenticationSecret([]byte("secret"))
	assertNil(t, err)
	assertEquals(t, alice.smp.state, smpStateExpect3{})

	_, bobMessages, err = bob.Receive(aliceMessages[0])
	assertNil(t, err)
	assertEquals(t, bob.smp.state, smpStateExpect4{})

	_, aliceMessages, err = alice.Receive(bobMessages[0])
	assertNil(t, err)
	assertEquals(t, alice.smp.state, smpStateExpect1{})

	_, bobMessages, err = bob.Receive(aliceMessages[0])
	assertNil(t, err)
	assertEquals(t, bob.smp.state, smpStateExpect1{})

}
