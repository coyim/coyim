package otr3

import (
	"crypto/rand"
	"testing"
	"time"
)

func Test_conversation_SMPStateMachineStartsAtSmpExpect1(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	assertEquals(t, c.smp.state, smpStateExpect1{})
}

func Test_receive_generatesErrorIfDoesNotHaveASecureChannel(t *testing.T) {
	states := []msgState{
		plainText, finished,
	}
	c := bobContextAfterAKE()
	c.msgState = encrypted
	smpMsg := fixtureMessage1()
	dataMsg, _, _ := c.genDataMsg(nil, smpMsg.tlv())
	m := dataMsg.serialize(c.version)
	m, _ = c.wrapMessageHeader(msgTypeData, m)
	for _, s := range states {
		c.msgState = s
		c.expectMessageEvent(t, func() {
			_, _, err := c.receiveDecoded(m)
			assertEquals(t, err, errMessageNotInPrivate)
		}, MessageEventReceivedMessageNotInPrivate, nil, nil)
	}
}

func Test_receive_doesntGenerateErrorIfThereIsNoSecureChannelButTheMessageIsIGNORE_UNREADABLE(t *testing.T) {
	states := []msgState{
		plainText, finished,
	}
	c := bobContextAfterAKE()
	c.msgState = encrypted
	smpMsg := fixtureMessage1()
	dataMsg, _, _ := c.genDataMsgWithFlag(nil, messageFlagIgnoreUnreadable, smpMsg.tlv())
	m, _ := c.wrapMessageHeader(msgTypeData, dataMsg.serialize(c.version))

	for _, s := range states {
		c.msgState = s
		_, _, err := c.receiveDecoded(m)
		assertNil(t, err)
	}
}

func Test_AKE_forVersion3And2InThePolicy(t *testing.T) {
	alice := &Conversation{Rand: rand.Reader}
	alice.SetOurKeys([]PrivateKey{alicePrivateKey})
	alice.ourCurrentKey = alicePrivateKey
	alice.Policies = policies(allowV2 | allowV3)

	bob := &Conversation{Rand: rand.Reader}
	bob.SetOurKeys([]PrivateKey{bobPrivateKey})
	bob.ourCurrentKey = bobPrivateKey
	bob.Policies = policies(allowV2 | allowV3)

	var toSend []ValidMessage
	var err error
	msg := alice.QueryMessage()

	//Alice send Bob queryMsg
	_, toSend, err = bob.Receive(msg)
	assertEquals(t, err, nil)
	assertEquals(t, bob.ake.state, authStateAwaitingDHKey{})

	//Bob send Alice DHCommit
	_, toSend, err = alice.Receive(toSend[0])
	assertEquals(t, alice.ake.state, authStateAwaitingRevealSig{})
	assertEquals(t, err, nil)

	//Alice send Bob DHKey
	_, toSend, err = bob.Receive(toSend[0])
	m, _ := bob.decode(encodedMessage(toSend[0]))
	assertEquals(t, err, nil)
	assertDeepEquals(t, bob.ake.state, authStateAwaitingSig{revealSigMsg: m})

	//Bob send Alice RevealSig
	_, toSend, err = alice.Receive(toSend[0])
	assertEquals(t, err, nil)
	assertEquals(t, alice.ake.state, authStateNone{})

	//Alice send Bob Sig
	_, toSend, err = bob.Receive(toSend[0])
	assertEquals(t, err, nil)
	assertEquals(t, bob.ake.state, authStateNone{})

	// "When starting a private Conversation [...],
	// generate two DH key pairs for yourself, and set our_keyid = 2"
	assertEquals(t, alice.keys.ourKeyID, uint32(2))
	assertEquals(t, alice.keys.ourCurrentDHKeys.priv.BitLen() > 0, true)
	assertEquals(t, alice.keys.ourCurrentDHKeys.pub.BitLen() > 0, true)
	assertEquals(t, alice.keys.ourPreviousDHKeys.priv.BitLen() > 0, true)
	assertEquals(t, alice.keys.ourPreviousDHKeys.pub.BitLen() > 0, true)

	assertEquals(t, bob.keys.ourKeyID, uint32(2))
	assertEquals(t, bob.keys.ourCurrentDHKeys.priv.BitLen() > 0, true)
	assertEquals(t, bob.keys.ourCurrentDHKeys.pub.BitLen() > 0, true)
	assertEquals(t, bob.keys.ourPreviousDHKeys.priv.BitLen() > 0, true)
	assertEquals(t, bob.keys.ourPreviousDHKeys.pub.BitLen() > 0, true)
}

func Test_AKE_withVersion3ButWithoutVersion2InThePolicy(t *testing.T) {
	alice := &Conversation{Rand: rand.Reader}
	alice.SetOurKeys([]PrivateKey{alicePrivateKey})
	alice.ourCurrentKey = alicePrivateKey
	alice.Policies = policies(allowV3)

	bob := &Conversation{Rand: rand.Reader}
	bob.SetOurKeys([]PrivateKey{bobPrivateKey})
	bob.ourCurrentKey = bobPrivateKey
	bob.Policies = policies(allowV3)

	var toSend []ValidMessage
	var err error
	msg := alice.QueryMessage()

	//Alice send Bob queryMsg
	_, toSend, err = bob.Receive(msg)
	assertEquals(t, err, nil)
	assertEquals(t, bob.ake.state, authStateAwaitingDHKey{})

	//Bob send Alice DHCommit
	_, toSend, err = alice.Receive(toSend[0])
	assertEquals(t, alice.ake.state, authStateAwaitingRevealSig{})
	assertEquals(t, err, nil)

	//Alice send Bob DHKey
	_, toSend, err = bob.Receive(toSend[0])
	m, _ := bob.decode(encodedMessage(toSend[0]))
	assertEquals(t, err, nil)
	assertDeepEquals(t, bob.ake.state, authStateAwaitingSig{revealSigMsg: m})

	//Bob send Alice RevealSig
	_, toSend, err = alice.Receive(toSend[0])
	assertEquals(t, err, nil)
	assertEquals(t, alice.ake.state, authStateNone{})

	//Alice send Bob Sig
	_, toSend, err = bob.Receive(toSend[0])
	assertEquals(t, err, nil)
	assertEquals(t, bob.ake.state, authStateNone{})

	// "When starting a private Conversation [...],
	// generate two DH key pairs for yourself, and set our_keyid = 2"
	assertEquals(t, alice.keys.ourKeyID, uint32(2))
	assertEquals(t, alice.keys.ourCurrentDHKeys.priv.BitLen() > 0, true)
	assertEquals(t, alice.keys.ourCurrentDHKeys.pub.BitLen() > 0, true)
	assertEquals(t, alice.keys.ourPreviousDHKeys.priv.BitLen() > 0, true)
	assertEquals(t, alice.keys.ourPreviousDHKeys.pub.BitLen() > 0, true)

	assertEquals(t, bob.keys.ourKeyID, uint32(2))
	assertEquals(t, bob.keys.ourCurrentDHKeys.priv.BitLen() > 0, true)
	assertEquals(t, bob.keys.ourCurrentDHKeys.pub.BitLen() > 0, true)
	assertEquals(t, bob.keys.ourPreviousDHKeys.priv.BitLen() > 0, true)
	assertEquals(t, bob.keys.ourPreviousDHKeys.pub.BitLen() > 0, true)
}

func Test_processDataMessageShouldExtractData(t *testing.T) {
	var toSend []ValidMessage
	var err error

	alice := &Conversation{Rand: rand.Reader}
	alice.Policies = policies(allowV2 | allowV3)
	alice.SetOurKeys([]PrivateKey{alicePrivateKey})

	bob := &Conversation{Rand: rand.Reader}
	bob.Policies = policies(allowV2 | allowV3)
	bob.SetOurKeys([]PrivateKey{bobPrivateKey})

	msg := []byte("?OTRv3?")

	//Alice send Bob queryMsg
	_, toSend, err = bob.Receive(msg)
	assertNil(t, err)
	assertEquals(t, bob.ake.state, authStateAwaitingDHKey{})
	assertEquals(t, bob.version, otrV3{})

	//Bob send Alice DHCommit
	_, toSend, err = alice.Receive(toSend[0])
	assertEquals(t, alice.ake.state, authStateAwaitingRevealSig{})
	assertNil(t, err)

	//Alice send Bob DHKey
	_, toSend, err = bob.Receive(toSend[0])
	m, _ := bob.decode(encodedMessage(toSend[0]))
	assertNil(t, err)
	assertDeepEquals(t, bob.ake.state, authStateAwaitingSig{revealSigMsg: m})

	//Bob send Alice RevealSig
	_, toSend, err = alice.Receive(toSend[0])
	assertNil(t, err)
	assertEquals(t, alice.ake.state, authStateNone{})

	//Alice send Bob Sig
	_, toSend, err = bob.Receive(toSend[0])
	assertNil(t, err)
	assertEquals(t, bob.ake.state, authStateNone{})

	// Alice sends a message to bob
	msg = []byte("hello")
	dataMsg, _, _ := alice.genDataMsg(msg)
	m, _ = alice.wrapMessageHeader(msgTypeData, dataMsg.serialize(alice.version))

	bob.updateLastSent()
	plain, ret, err := bob.receiveDecoded(m)

	assertNil(t, err)
	assertDeepEquals(t, plain, MessagePlaintext(msg))
	assertNil(t, ret)
}

func Test_startingAKE_shouldNotBreakTheEncryptedChannel(t *testing.T) {
	var hello = []byte("hello")
	var toSend []ValidMessage
	var err error

	alice := &Conversation{Rand: rand.Reader}
	alice.Policies = policies(allowV2 | allowV3)
	alice.SetOurKeys([]PrivateKey{alicePrivateKey})

	bob := &Conversation{Rand: rand.Reader}
	bob.Policies = policies(allowV2 | allowV3)
	bob.SetOurKeys([]PrivateKey{bobPrivateKey})

	//Alice send Bob queryMsg
	_, toSend, err = bob.Receive(alice.QueryMessage())
	assertNil(t, err)
	assertEquals(t, bob.ake.state, authStateAwaitingDHKey{})
	assertEquals(t, bob.version, otrV3{})

	//Bob send Alice DHCommit
	_, toSend, err = alice.Receive(toSend[0])
	assertEquals(t, alice.ake.state, authStateAwaitingRevealSig{})
	assertNil(t, err)

	//Alice send Bob DHKey
	_, toSend, err = bob.Receive(toSend[0])
	dec, _ := bob.decode(encodedMessage(toSend[0]))
	assertNil(t, err)
	assertDeepEquals(t, bob.ake.state, authStateAwaitingSig{revealSigMsg: dec})

	//Bob send Alice RevealSig
	_, toSend, err = alice.Receive(toSend[0])
	assertNil(t, err)
	assertEquals(t, alice.ake.state, authStateNone{})

	//Alice send Bob Sig
	_, toSend, err = bob.Receive(toSend[0])
	assertNil(t, err)
	assertEquals(t, bob.ake.state, authStateNone{})

	// Alice sends a message to bob
	m, err := alice.Send(hello)

	bob.updateLastSent()
	plain, ret, err := bob.Receive(m[0])

	assertNil(t, err)
	assertDeepEquals(t, plain, MessagePlaintext(hello))
	assertNil(t, ret)

	//
	// Restarts the AKE
	//

	//Alice send Bob queryMsg
	bob.lastMessageStateChange = time.Time{}
	bob.ake.lastStateChange = time.Time{}
	_, toSend, err = bob.Receive(alice.QueryMessage())
	assertNil(t, err)
	assertEquals(t, bob.ake.state, authStateAwaitingDHKey{})
	assertEquals(t, bob.version, otrV3{})

	//Alice sends Bob a data message
	m, err = alice.Send(hello)
	assertNil(t, err)

	plain, ret, err = bob.Receive(m[0])
	assertNil(t, err)
	assertDeepEquals(t, plain, MessagePlaintext(hello))
	assertNil(t, ret)
}
