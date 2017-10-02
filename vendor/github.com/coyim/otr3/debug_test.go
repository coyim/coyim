package otr3

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"testing"
)

func Test_dumpSMP_dumpsTheCurrentSMPState(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	c.smp.state = smpStateExpect2{}
	c.smp.s1 = fixtureSmp1()
	c.smp.secret = bnFromHex("ABCDE56321F9A9F8E364607C8C82DECD8E8E6209E2CB952C7E649620F5286FE3")
	q := "Blarg"
	c.smp.question = &q

	bt := bytes.NewBuffer(make([]byte, 0, 200))
	c.dumpSMP(bufio.NewWriter(bt))
	assertDeepEquals(t, bt.String(), `  SM state:
    Next expected: 2 (EXPECT2)
    Received_Q: 1
`)
}

func Test_dumpSMP_dumpsTheCurrentSMPStateWithQuestion(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	c.smp.state = smpStateExpect2{}
	c.smp.s1 = fixtureSmp1()
	c.smp.secret = bnFromHex("ABCDE56321F9A9F8E364607C8C82DECD8E8E6209E2CB952C7E649620F5286FE3")

	bt := bytes.NewBuffer(make([]byte, 0, 200))
	c.dumpSMP(bufio.NewWriter(bt))
	assertDeepEquals(t, bt.String(), `  SM state:
    Next expected: 2 (EXPECT2)
    Received_Q: 0
`)
}

func Test_identity_isCorrectForAllSMPStates(t *testing.T) {
	assertEquals(t, smpStateExpect1{}.identity(), 0)
	assertEquals(t, smpStateWaitingForSecret{}.identity(), 1)
	assertEquals(t, smpStateExpect2{}.identity(), 2)
	assertEquals(t, smpStateExpect3{}.identity(), 3)
	assertEquals(t, smpStateExpect4{}.identity(), 4)
}

func Test_identityString_isCorrectForAllSMPStates(t *testing.T) {
	assertEquals(t, smpStateExpect1{}.identityString(), "EXPECT1")
	assertEquals(t, smpStateWaitingForSecret{}.identityString(), "EXPECT1_WQ")
	assertEquals(t, smpStateExpect2{}.identityString(), "EXPECT2")
	assertEquals(t, smpStateExpect3{}.identityString(), "EXPECT3")
	assertEquals(t, smpStateExpect4{}.identityString(), "EXPECT4")
}

func Test_dumpAKE_dumpsTheCurrentAKEState(t *testing.T) {
	c := aliceContextAtAwaitingRevealSig()
	c.theirKey = bobPrivateKey.PublicKey()
	bt := bytes.NewBuffer(make([]byte, 0, 200))
	c.dumpAKE(bufio.NewWriter(bt))
	assertDeepEquals(t, bt.String(), `  Auth info:
    State: 2 (AWAITING_REVEALSIG)
    Our keyid:   0
    Their keyid: 0
    Their fingerprint: 8798FAA7735267FB8457733098482E94096D4ABD
    Proto version = 2
`)
}

func Test_identity_isCorrectForAllAKEStates(t *testing.T) {
	assertEquals(t, authStateNone{}.identity(), 0)
	assertEquals(t, authStateAwaitingDHKey{}.identity(), 1)
	assertEquals(t, authStateAwaitingRevealSig{}.identity(), 2)
	assertEquals(t, authStateAwaitingSig{}.identity(), 3)
}

func Test_identityString_isCorrectForAllAKEStates(t *testing.T) {
	assertEquals(t, authStateNone{}.identityString(), "NONE")
	assertEquals(t, authStateAwaitingDHKey{}.identityString(), "AWAITING_DHKEY")
	assertEquals(t, authStateAwaitingRevealSig{}.identityString(), "AWAITING_REVEALSIG")
	assertEquals(t, authStateAwaitingSig{}.identityString(), "AWAITING_SIG")
}

func Test_dump_dumpsAllKindsOfConversationState(t *testing.T) {
	c := bobContextAfterAKE()
	c.ake = nil
	c.msgState = encrypted
	c.whitespaceState = whitespaceSent
	c.theirInstanceTag = 0x102

	bt := bytes.NewBuffer(make([]byte, 0, 200))
	c.dump(bufio.NewWriter(bt))
	assertDeepEquals(t, bt.String(), `Context:

  Our instance:   00000101
  Their instance: 00000102

  Msgstate: 1 (ENCRYPTED)

  Protocol version: 3
  OTR offer: ACCEPTED

  Auth info: NULL

  SM state:
    Next expected: 0 (EXPECT1)
    Received_Q: 0
`)
}

func Test_identityString_isCorrectForAllMessageStates(t *testing.T) {
	assertEquals(t, plainText.identityString(), "PLAINTEXT")
	assertEquals(t, encrypted.identityString(), "ENCRYPTED")
	assertEquals(t, finished.identityString(), "FINISHED")
	assertEquals(t, msgState(99).identityString(), "INVALID")
}

func Test_otrOffer_isCorrectForNotSent(t *testing.T) {
	c := &Conversation{whitespaceState: whitespaceNotSent}
	assertEquals(t, c.otrOffer(), "NOT")
}

func Test_otrOffer_isCorrectForRejected(t *testing.T) {
	c := &Conversation{whitespaceState: whitespaceRejected}
	assertEquals(t, c.otrOffer(), "REJECTED")
}

func Test_otrOffer_isCorrectForInvalid(t *testing.T) {
	c := &Conversation{whitespaceState: whitespaceState(99)}
	assertEquals(t, c.otrOffer(), "INVALID")
}

func Test_otrOffer_isCorrectForSentAndAccepted(t *testing.T) {
	c := &Conversation{whitespaceState: whitespaceSent, msgState: encrypted}
	assertEquals(t, c.otrOffer(), "ACCEPTED")
}

func Test_otrOffer_isCorrectForSentAndNotAcceptedYet(t *testing.T) {
	c := &Conversation{whitespaceState: whitespaceSent, msgState: plainText}
	assertEquals(t, c.otrOffer(), "SENT")
}

func Test_setDebug_setsTheDebugFlag(t *testing.T) {
	c := &Conversation{}
	c.SetDebug(true)
	assertTrue(t, c.debug)
	c.SetDebug(false)
	assertFalse(t, c.debug)
}

func Test_SMP_CompleteDebug(t *testing.T) {
	alice := &Conversation{Rand: rand.Reader}
	alice.ourCurrentKey = alicePrivateKey
	alice.SetOurKeys([]PrivateKey{alicePrivateKey})

	alice.Policies = policies(allowV3)

	bob := &Conversation{Rand: rand.Reader}
	bob.ourCurrentKey = bobPrivateKey
	bob.SetOurKeys([]PrivateKey{bobPrivateKey})
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

	bt := bytes.NewBuffer(make([]byte, 0, 200))

	// Should not they all be "Next expected: 0 (EXPECT1)" ?
	// The spec says: This is the default state when SMP has not yet begun
	// This is probably a consequence of needing ensureSMP() to initialize

	bt.Reset()
	alice.dumpSMP(bufio.NewWriter(bt))
	assertDeepEquals(t, bt.String(), `  SM state:
    Received_Q: 0
`)

	bt.Reset()
	bob.dumpSMP(bufio.NewWriter(bt))
	assertDeepEquals(t, bt.String(), `  SM state:
    Received_Q: 0
`)

	bobMessages, err = bob.StartAuthenticate("", []byte("secret"))
	assertNil(t, err)
	assertEquals(t, bob.smp.state, smpStateExpect2{})

	bt.Reset()
	bob.dumpSMP(bufio.NewWriter(bt))
	assertDeepEquals(t, bt.String(), `  SM state:
    Next expected: 2 (EXPECT2)
    Received_Q: 0
`)

	_, aliceMessages, err = alice.Receive(bobMessages[0])
	assertNil(t, err)

	// this is an internal state
	_, ok := alice.smp.state.(smpStateWaitingForSecret)
	assertEquals(t, ok, true)

	bt.Reset()
	alice.dumpSMP(bufio.NewWriter(bt))

	// This feels wrong.
	// I read "Alice is expecting a SMP1Q" which is not true because she has already received a SMP1Q
	// BTW, libotr otrl_sm_dump() does not have it
	assertDeepEquals(t, bt.String(), `  SM state:
    Next expected: 1 (EXPECT1_WQ)
    Received_Q: 0
`)

	aliceMessages, err = alice.ProvideAuthenticationSecret([]byte("secret"))
	assertNil(t, err)
	assertEquals(t, alice.smp.state, smpStateExpect3{})

	bt.Reset()
	alice.dumpSMP(bufio.NewWriter(bt))
	assertDeepEquals(t, bt.String(), `  SM state:
    Next expected: 3 (EXPECT3)
    Received_Q: 0
`)

	_, bobMessages, err = bob.Receive(aliceMessages[0])
	assertNil(t, err)
	assertEquals(t, bob.smp.state, smpStateExpect4{})

	bt.Reset()
	bob.dumpSMP(bufio.NewWriter(bt))
	assertDeepEquals(t, bt.String(), `  SM state:
    Next expected: 4 (EXPECT4)
    Received_Q: 0
`)

	_, aliceMessages, err = alice.Receive(bobMessages[0])
	assertNil(t, err)
	assertEquals(t, alice.smp.state, smpStateExpect1{})

	bt.Reset()
	alice.dumpSMP(bufio.NewWriter(bt))
	assertDeepEquals(t, bt.String(), `  SM state:
    Next expected: 0 (EXPECT1)
    Received_Q: 0
`)

	_, bobMessages, err = bob.Receive(aliceMessages[0])
	assertNil(t, err)
	assertEquals(t, bob.smp.state, smpStateExpect1{})

	bt.Reset()
	alice.dumpSMP(bufio.NewWriter(bt))
	assertDeepEquals(t, bt.String(), `  SM state:
    Next expected: 0 (EXPECT1)
    Received_Q: 0
`)

	bobMessages, err = bob.StartAuthenticate("What is the secret?", []byte("secret"))
	assertNil(t, err)
	assertEquals(t, bob.smp.state, smpStateExpect2{})

	bt.Reset()
	bob.dumpSMP(bufio.NewWriter(bt))
	assertDeepEquals(t, bt.String(), `  SM state:
    Next expected: 2 (EXPECT2)
    Received_Q: 0
`)

	_, aliceMessages, err = alice.Receive(bobMessages[0])
	assertNil(t, err)

	// this is an internal state
	_, ok = alice.smp.state.(smpStateWaitingForSecret)
	assertEquals(t, ok, true)

	bt.Reset()
	alice.dumpSMP(bufio.NewWriter(bt))

	// This feels wrong.
	// I read "Alice is expecting a SMP1Q" which is not true because she has already received a SMP1Q
	// BTW, libotr otrl_sm_dump() does not have it
	assertDeepEquals(t, bt.String(), `  SM state:
    Next expected: 1 (EXPECT1_WQ)
    Received_Q: 1
`)

	aliceMessages, err = alice.ProvideAuthenticationSecret([]byte("secret"))
	assertNil(t, err)
	assertEquals(t, alice.smp.state, smpStateExpect3{})

	bt.Reset()
	alice.dumpSMP(bufio.NewWriter(bt))
	assertDeepEquals(t, bt.String(), `  SM state:
    Next expected: 3 (EXPECT3)
    Received_Q: 1
`)

	_, bobMessages, err = bob.Receive(aliceMessages[0])
	assertNil(t, err)
	assertEquals(t, bob.smp.state, smpStateExpect4{})

	bt.Reset()
	bob.dumpSMP(bufio.NewWriter(bt))
	assertDeepEquals(t, bt.String(), `  SM state:
    Next expected: 4 (EXPECT4)
    Received_Q: 0
`)

	_, aliceMessages, err = alice.Receive(bobMessages[0])
	assertNil(t, err)
	assertEquals(t, alice.smp.state, smpStateExpect1{})

	bt.Reset()
	alice.dumpSMP(bufio.NewWriter(bt))
	assertDeepEquals(t, bt.String(), `  SM state:
    Next expected: 0 (EXPECT1)
    Received_Q: 0
`)

	_, bobMessages, err = bob.Receive(aliceMessages[0])
	assertNil(t, err)
	assertEquals(t, bob.smp.state, smpStateExpect1{})

	bt.Reset()
	alice.dumpSMP(bufio.NewWriter(bt))
	assertDeepEquals(t, bt.String(), `  SM state:
    Next expected: 0 (EXPECT1)
    Received_Q: 0
`)
}
