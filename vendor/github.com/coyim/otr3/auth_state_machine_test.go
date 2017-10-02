package otr3

import (
	"math/big"
	"testing"
)

func Test_conversationInitialState(t *testing.T) {
	c := newConversation(nil, fixtureRand())
	assertEquals(t, c.ake.state, authStateNone{})
}

func Test_receiveDHCommit_TransitionsFromNoneToAwaitingRevealSigAndSendDHKeyMsg(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	nextState, nextMsg, e := authStateNone{}.receiveDHCommitMessage(c, fixtureDHCommitMsgBody())

	assertEquals(t, e, nil)
	assertEquals(t, nextState, authStateAwaitingRevealSig{})
	assertEquals(t, dhMsgType(nextMsg), msgTypeDHKey)
}

func Test_receiveDHCommit_AtAuthStateNoneStoresGyAndY(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	authStateNone{}.receiveDHCommitMessage(c, fixtureDHCommitMsg())

	assertDeepEquals(t, c.ake.ourPublicValue, fixedGY())
	assertDeepEquals(t, c.ake.secretExponent, fixedY())
}

func Test_receiveDHCommit_AtAuthStateNoneStoresEncryptedGxAndHashedGx(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())

	dhCommitMsg := fixtureDHCommitMsgBody()
	newMsg, encryptedGx, _ := extractData(dhCommitMsg)
	_, hashedGx, _ := extractData(newMsg)

	authStateNone{}.receiveDHCommitMessage(c, dhCommitMsg)

	assertDeepEquals(t, c.ake.xhashedGx, hashedGx)
	assertDeepEquals(t, c.ake.encryptedGx, encryptedGx)
}

func Test_receiveDHCommit_ResendPreviousDHKeyMsgFromAwaitingRevealSig(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())

	authAwaitingRevSig, prevDHKeyMsg, err := authStateNone{}.receiveDHCommitMessage(c, fixtureDHCommitMsgBody())
	assertNil(t, err)
	assertEquals(t, authAwaitingRevSig, authStateAwaitingRevealSig{})

	nextState, msg, err := authAwaitingRevSig.receiveDHCommitMessage(c, fixtureDHCommitMsgBody())

	assertNil(t, err)
	assertEquals(t, nextState, authStateAwaitingRevealSig{})
	assertEquals(t, dhMsgType(msg), msgTypeDHKey)
	assertDeepEquals(t, prevDHKeyMsg, msg)
}

func Test_receiveDHCommit_AtAuthAwaitingRevealSigiForgetOldEncryptedGxAndHashedGx(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	c.ake.encryptedGx = []byte{0x02}                                 //some encryptedGx
	c.ake.xhashedGx = fixedSize(otrV3{}.hash2Length(), []byte{0x05}) //some hashedGx

	newDHCommitMsg := fixtureDHCommitMsgBody()
	newMsg, newEncryptedGx, _ := extractData(newDHCommitMsg)
	_, newHashedGx, _ := extractData(newMsg)

	authStateNone{}.receiveDHCommitMessage(c, fixtureDHCommitMsgBody())

	authStateAwaitingRevealSig{}.receiveDHCommitMessage(c, newDHCommitMsg)
	assertDeepEquals(t, c.ake.encryptedGx, newEncryptedGx)
	assertDeepEquals(t, c.ake.xhashedGx, newHashedGx)
}

func Test_receiveDHCommit_AtAuthAwaitingSigTransitionsToAwaitingRevSigAndSendsNewDHKeyMsg(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())

	authAwaitingRevSig, msg, err := authStateAwaitingSig{}.receiveDHCommitMessage(c, fixtureDHCommitMsgBody())
	assertNil(t, err)
	assertEquals(t, authAwaitingRevSig, authStateAwaitingRevealSig{})
	assertEquals(t, dhMsgType(msg), msgTypeDHKey)
}

func Test_receiveDHCommit_AtAwaitingDHKeyIgnoreIncomingMsgAndResendOurDHCommitMsgIfOurHashIsHigher(t *testing.T) {
	ourDHCommitAKE := fixtureConversation()
	m, _ := ourDHCommitAKE.dhCommitMessage()
	ourDHMsg, _ := ourDHCommitAKE.wrapMessageHeader(msgTypeDHCommit, m)

	//make sure we store the same values when creating the DH commit
	c := newConversation(otrV3{}, fixtureRand())
	c.initAKE()
	c.ake.encryptedGx = ourDHCommitAKE.ake.encryptedGx
	c.ake.ourPublicValue = ourDHCommitAKE.ake.ourPublicValue

	// force their hashedGx to be lower than ours
	msg := fixtureDHCommitMsgBody()
	newPoint, _, _ := extractData(msg)
	newPoint[4] = 0x00

	state, newMsg, err := authStateAwaitingDHKey{}.receiveDHCommitMessage(c, msg)
	assertDeepEquals(t, err, nil)
	assertEquals(t, state, authStateAwaitingRevealSig{})
	assertDeepEquals(t, newMsg, ourDHMsg)
}

func Test_receiveDHCommit_AtAwaitingDHKeyForgetOurGxAndSendDHKeyMsgAndGoToAwaitingRevealSig(t *testing.T) {
	ourDHCommitAKE := fixtureConversation()
	ourDHCommitAKE.dhCommitMessage()

	//make sure we store the same values when creating the DH commit
	c := newConversation(otrV3{}, fixtureRand())
	c.initAKE()
	c.ake.ourPublicValue = ourDHCommitAKE.ake.ourPublicValue

	// force their hashedGx to be higher than ours
	msg := fixtureDHCommitMsgBody()
	newPoint, _, _ := extractData(msg)
	newPoint[4] = 0xFF

	state, newMsg, err := authStateAwaitingDHKey{}.receiveDHCommitMessage(c, msg)
	assertDeepEquals(t, err, nil)
	assertEquals(t, state, authStateAwaitingRevealSig{})
	assertEquals(t, dhMsgType(newMsg), msgTypeDHKey)
	assertDeepEquals(t, c.ake.ourPublicValue, fixedGY())
	assertDeepEquals(t, c.ake.secretExponent, fixedY())
}

func Test_receiveDHKey_AtAuthStateNoneOrAuthStateAwaitingRevealSigIgnoreIt(t *testing.T) {
	c := newConversation(otrV3{}, fixtureRand())
	c.initAKE()
	dhKeymsg := fixtureDHKeyMsg(otrV3{})

	states := []authState{
		authStateNone{},
		authStateAwaitingRevealSig{},
	}

	for _, s := range states {
		state, msg, err := s.receiveDHKeyMessage(c, dhKeymsg)
		assertNil(t, err)
		assertEquals(t, state, s)
		assertNil(t, msg)
	}
}

func Test_receiveDHKey_TransitionsFromAwaitingDHKeyToAwaitingSigAndSendsRevealSig(t *testing.T) {
	ourDHCommitAKE := fixtureConversation()
	ourDHCommitAKE.dhCommitMessage()

	c := bobContextAtAwaitingDHKey()
	c.sentRevealSig = false

	state, msg, err := authStateAwaitingDHKey{}.receiveDHKeyMessage(c, fixtureDHKeyMsg(otrV3{})[otrv3HeaderLen:])

	_, ok := state.(authStateAwaitingSig)
	assertNil(t, err)
	assertEquals(t, ok, true)
	assertEquals(t, dhMsgType(msg), msgTypeRevealSig)
	assertEquals(t, dhMsgVersion(msg), uint16(3))
	assertEquals(t, c.sentRevealSig, true)
}

func Test_receiveDHKey_AtAwaitingDHKeyStoresGyAndSigKey(t *testing.T) {
	ourDHCommitAKE := fixtureConversation()
	ourDHCommitAKE.dhCommitMessage()

	c := bobContextAtAwaitingDHKey()

	_, _, err := authStateAwaitingDHKey{}.receiveDHKeyMessage(c, fixtureDHKeyMsg(otrV3{})[otrv3HeaderLen:])

	assertEquals(t, err, nil)
	assertDeepEquals(t, c.ake.theirPublicValue, fixedGY())
	assertDeepEquals(t, c.ake.sigKey.c, expectedC)
	assertDeepEquals(t, c.ake.sigKey.m1, expectedM1)
	assertDeepEquals(t, c.ake.sigKey.m2, expectedM2)
}

func Test_receiveDHKey_AtAwaitingDHKey_storesOursAndTheirDHKeys(t *testing.T) {
	ourDHCommitAKE := fixtureConversation()
	ourDHCommitAKE.dhCommitMessage()

	c := bobContextAtAwaitingDHKey()

	_, _, err := authStateAwaitingDHKey{}.receiveDHKeyMessage(c, fixtureDHKeyMsg(otrV3{})[otrv3HeaderLen:])

	assertEquals(t, err, nil)
	assertDeepEquals(t, c.ake.keys.theirCurrentDHPubKey, fixedGY())
	assertDeepEquals(t, c.ake.keys.ourCurrentDHKeys.pub, fixedGX())
	assertDeepEquals(t, c.ake.keys.ourCurrentDHKeys.priv, fixedX())

	assertEquals(t, c.ake.keys.ourKeyID, uint32(1))
	assertEquals(t, c.ake.keys.theirKeyID, uint32(0))
}

func Test_receiveDHKey_AtAuthAwaitingSigIfReceivesSameDHKeyMsgRetransmitRevealSigMsg(t *testing.T) {
	var nilB *big.Int

	ourDHCommitAKE := fixtureConversation()
	ourDHCommitAKE.dhCommitMessage()

	c := newConversation(otrV3{}, fixtureRand())
	c.initAKE()
	c.setSecretExponent(ourDHCommitAKE.ake.secretExponent)
	c.ourCurrentKey = bobPrivateKey

	assertDeepEquals(t, c.ake.theirPublicValue, nilB)

	sameDHKeyMsg := fixtureDHKeyMsg(otrV3{})[otrv3HeaderLen:]
	sigState, previousRevealSig, err := authStateAwaitingDHKey{}.receiveDHKeyMessage(c, sameDHKeyMsg)

	assertNil(t, err)
	assertDeepEquals(t, c.ake.theirPublicValue, fixedGY())

	state, msg, err := sigState.receiveDHKeyMessage(c, sameDHKeyMsg)
	_, sameStateType := state.(authStateAwaitingSig)

	assertNil(t, err)
	assertDeepEquals(t, c.ake.theirPublicValue, fixedGY())
	assertDeepEquals(t, sameStateType, true)
	assertDeepEquals(t, msg, previousRevealSig)
}

func Test_receiveDHKey_AtAuthAwaitingSigIgnoresMsgIfIsNotSameDHKeyMsg(t *testing.T) {
	newDHKeyMsg := fixtureDHKeyMsgBody(otrV3{})
	c := newConversation(otrV3{}, fixtureRand())
	c.initAKE()

	state, msg, err := authStateAwaitingSig{}.receiveDHKeyMessage(c, newDHKeyMsg)

	_, sameStateType := state.(authStateAwaitingSig)
	assertNil(t, err)
	assertNil(t, msg)
	assertDeepEquals(t, sameStateType, true)
}

func Test_receiveRevealSig_TransitionsFromAwaitingRevealSigToNoneOnSuccess(t *testing.T) {
	revealSignMsg := fixtureRevealSigMsgBody(otrV2{})

	c := aliceContextAtAwaitingRevealSig()
	c.sentRevealSig = true

	state, msg, err := authStateAwaitingRevealSig{}.receiveRevealSigMessage(c, revealSignMsg)

	assertEquals(t, err, nil)
	assertEquals(t, state, authStateNone{})
	assertEquals(t, dhMsgType(msg), msgTypeSig)
	assertEquals(t, c.sentRevealSig, false)
}

func Test_receiveRevealSig_AtAwaitingRevealSig_savesAKEKeysToConversationAndGenerateANewPairOfKeys(t *testing.T) {
	revealSignMsg := fixtureRevealSigMsgBody(otrV2{})

	c := aliceContextAtAwaitingRevealSig()

	_, _, err := authStateAwaitingRevealSig{}.receiveRevealSigMessage(c, revealSignMsg)

	assertNil(t, err)
	assertDeepEquals(t, c.keys.theirCurrentDHPubKey, fixedGX())
	assertNil(t, c.keys.theirPreviousDHPubKey)
	assertDeepEquals(t, c.keys.ourPreviousDHKeys.pub, fixedGY())
	assertDeepEquals(t, c.keys.ourPreviousDHKeys.priv, fixedY())

	//should wipe
	assertDeepEquals(t, c.ake, &ake{state: c.ake.state})

	assertEquals(t, c.keys.ourKeyID, uint32(2))
	assertEquals(t, c.keys.theirKeyID, uint32(1))
}

func Test_authStateAwaitingRevealSig_receiveRevealSigMessage_returnsErrorIfProcessRevealSigFails(t *testing.T) {
	c := newConversation(otrV2{}, fixtureRand())
	c.Policies.add(allowV2)
	_, _, err := authStateAwaitingRevealSig{}.receiveRevealSigMessage(c, []byte{0x00, 0x00})
	assertDeepEquals(t, err, newOtrError("corrupt reveal signature message"))
}

func Test_receiveRevealSig_IgnoreMessageIfNotInStateAwaitingRevealSig(t *testing.T) {
	states := []authState{
		authStateNone{},
		authStateAwaitingDHKey{},
		authStateAwaitingSig{},
	}

	revealSignMsg := fixtureRevealSigMsg(otrV2{})

	for _, s := range states {
		c := newConversation(otrV3{}, fixtureRand())
		state, msg, err := s.receiveRevealSigMessage(c, revealSignMsg)

		assertNil(t, err)
		assertNil(t, msg)
		assertDeepEquals(t, state, s)
	}
}

func Test_receiveSig_TransitionsFromAwaitingSigToNoneOnSuccess(t *testing.T) {
	sigMsg := fixtureSigMsg(otrV2{})[otrv2HeaderLen:]
	c := bobContextAtAwaitingSig()

	state, msg, err := authStateAwaitingSig{}.receiveSigMessage(c, sigMsg)

	assertNil(t, err)
	assertNil(t, msg)
	assertEquals(t, state, authStateNone{})
	assertEquals(t, c.keys.theirKeyID, uint32(1))
}

func Test_receiveSig_IgnoreMessageIfNotInStateAwaitingSig(t *testing.T) {
	states := []authState{
		authStateNone{},
		authStateAwaitingDHKey{},
		authStateAwaitingRevealSig{},
	}

	revealSignMsg := fixtureRevealSigMsg(otrV2{})[otrv2HeaderLen:]

	for _, s := range states {
		c := newConversation(otrV3{}, fixtureRand())
		state, msg, err := s.receiveSigMessage(c, revealSignMsg)

		assertNil(t, err)
		assertNil(t, msg)
		assertEquals(t, state, s)
	}
}

func Test_receiveDecoded_receiveRevealSigMessageAndSetMessageStateToEncrypted(t *testing.T) {
	c := aliceContextAtAwaitingRevealSig()
	msg := fixtureRevealSigMsg(otrV2{})
	assertEquals(t, c.msgState, plainText)

	_, _, err := c.receiveDecoded(msg)

	assertNil(t, err)
	assertEquals(t, c.msgState, encrypted)
}

func Test_receiveDecoded_receiveRevealSigMessageWillResendPotentialLastMessage(t *testing.T) {
	c := aliceContextAtAwaitingRevealSig()
	c.resend.later(MessagePlaintext("what do you think turn 2"))
	c.resend.later(MessagePlaintext("I mean, about that thing"))
	c.resend.mayRetransmit = retransmitWithPrefix
	c.updateLastSent()
	msg := fixtureRevealSigMsg(otrV2{})

	c.expectMessageEvent(t, func() {
		_, toSends, _ := c.receiveDecoded(msg)
		assertEquals(t, len(toSends), 3)
	}, MessageEventMessageResent, nil, nil)
}

func Test_receiveDecoded_receiveRevealSigMessageAndStoresTheirKeyIDAndTheirCurrentDHPubKey(t *testing.T) {
	c := aliceContextAtAwaitingRevealSig()
	msg := fixtureRevealSigMsg(otrV2{})
	assertEquals(t, c.msgState, plainText)

	_, _, err := c.receiveDecoded(msg)

	assertNil(t, err)
	assertEquals(t, c.keys.theirKeyID, uint32(1))
	assertDeepEquals(t, c.keys.theirCurrentDHPubKey, fixedGX())
	assertNil(t, c.keys.theirPreviousDHPubKey)
}

func Test_receiveDecoded_receiveDHCommitMessageAndFailsWillSignalSetupError(t *testing.T) {
	c := aliceContextAtAwaitingDHCommit()
	c.Rand = fixedRand([]string{"ABCD"})
	msg := fixtureDHCommitMsgV2()

	c.expectMessageEvent(t, func() {
		c.receiveDecoded(msg)
	}, MessageEventSetupError, nil, errShortRandomRead)
}

func Test_receiveDecoded_receiveDHKeyMessageAndFailsWillSignalSetupError(t *testing.T) {
	c := bobContextAtAwaitingDHKey()
	c.Rand = fixedRand([]string{"ABCD"})
	msg := fixtureDHKeyMsg(otrV3{})

	c.expectMessageEvent(t, func() {
		c.receiveDecoded(msg)
	}, MessageEventSetupError, nil, errShortRandomRead)
}

func Test_receiveDecoded_receiveRevealSigMessageAndFailsWillSignalSetupError(t *testing.T) {
	c := aliceContextAtAwaitingRevealSig()
	c.Rand = fixedRand([]string{"ABCD"})
	msg := fixtureRevealSigMsg(otrV2{})

	c.expectMessageEvent(t, func() {
		c.receiveDecoded(msg)
	}, MessageEventSetupError, nil, errShortRandomRead)
}

func Test_receiveDecoded_receiveSigMessageAndSetMessageStateToEncrypted(t *testing.T) {
	c := bobContextAtAwaitingSig()
	c.Rand = fixedRand([]string{"ABCD"})
	msg := fixtureSigMsg(otrV2{})

	c.expectMessageEvent(t, func() {
		c.receiveDecoded(msg)
	}, MessageEventSetupError, nil, errShortRandomRead)
}

func Test_receiveDecoded_receiveSigMessageWillResendTheLastPotentialMessage(t *testing.T) {
	c := bobContextAtAwaitingSig()
	c.resend.later(MessagePlaintext("what do you think"))
	c.resend.later(MessagePlaintext("you think, dont you?"))
	c.resend.mayRetransmit = retransmitWithPrefix
	c.updateLastSent()

	msg := fixtureSigMsg(otrV2{})

	c.expectMessageEvent(t, func() {
		_, toSends, _ := c.receiveDecoded(msg)
		assertEquals(t, len(toSends), 2) // Only the retransmit messages, nothing else
	}, MessageEventMessageResent, nil, nil)
}

func Test_receiveDecoded_receiveSigMessageAndFailsWillSignalSetupError(t *testing.T) {
	c := bobContextAtAwaitingSig()
	msg := fixtureSigMsg(otrV2{})
	assertEquals(t, c.msgState, plainText)

	_, _, err := c.receiveDecoded(msg)

	assertNil(t, err)
	assertEquals(t, c.msgState, encrypted)
}

func Test_receiveDecoded_receiveSigMessageAndStoresTheirKeyIDAndTheirCurrentDHPubKey(t *testing.T) {
	var nilBigInt *big.Int

	c := bobContextAtAwaitingSig()

	msg := fixtureSigMsg(otrV2{})
	assertEquals(t, c.msgState, plainText)

	_, _, err := c.receiveDecoded(msg)

	assertNil(t, err)
	assertEquals(t, c.keys.theirKeyID, uint32(1))
	assertDeepEquals(t, c.keys.theirCurrentDHPubKey, fixedGY())
	assertEquals(t, c.keys.theirPreviousDHPubKey, nilBigInt)
}

func Test_authStateAwaitingDHKey_receiveDHKeyMessage_returnsErrorIfprocessDHKeyReturnsError(t *testing.T) {
	ourDHCommitAKE := fixtureConversation()
	ourDHCommitAKE.dhCommitMessage()

	c := newConversation(otrV3{}, fixtureRand())
	c.initAKE()
	c.setSecretExponent(ourDHCommitAKE.ake.secretExponent)
	c.ourCurrentKey = bobPrivateKey

	_, _, err := authStateAwaitingDHKey{}.receiveDHKeyMessage(c, []byte{0x00, 0x02})

	assertDeepEquals(t, err, newOtrError("corrupt DH key message"))
}

func Test_authStateAwaitingDHKey_receiveDHKeyMessage_returnsErrorIfrevealSigMessageReturnsError(t *testing.T) {
	ourDHCommitAKE := fixtureConversation()
	ourDHCommitAKE.dhCommitMessage()

	c := newConversation(otrV3{}, fixedRand([]string{"ABCD"}))
	c.initAKE()
	c.setSecretExponent(ourDHCommitAKE.ake.secretExponent)
	c.ourCurrentKey = bobPrivateKey

	sameDHKeyMsg := fixtureDHKeyMsgBody(otrV3{})
	_, _, err := authStateAwaitingDHKey{}.receiveDHKeyMessage(c, sameDHKeyMsg)

	assertEquals(t, err, errShortRandomRead)
}

func Test_authStateAwaitingSig_receiveDHKeyMessage_returnsErrorIfprocessDHKeyReturnsError(t *testing.T) {
	ourDHCommitAKE := fixtureConversation()
	ourDHCommitAKE.dhCommitMessage()

	c := newConversation(otrV3{}, fixtureRand())
	c.initAKE()
	c.setSecretExponent(ourDHCommitAKE.ake.secretExponent)
	c.ourCurrentKey = bobPrivateKey

	_, _, err := authStateAwaitingSig{}.receiveDHKeyMessage(c, []byte{0x01, 0x02})

	assertEquals(t, err, newOtrError("corrupt DH key message"))
}

func Test_authStateAwaitingSig_receiveSigMessage_returnsErrorIfProcessSigFails(t *testing.T) {
	c := newConversation(otrV2{}, fixtureRand())
	c.Policies.add(allowV2)
	_, _, err := authStateAwaitingSig{}.receiveSigMessage(c, []byte{0x00, 0x00})
	assertEquals(t, err, newOtrError("corrupt signature message"))
}

func Test_authStateAwaitingRevealSig_receiveDHCommitMessage_returnsErrorIfProcessDHCommitOrGenerateCommitInstanceTagsFailsFails(t *testing.T) {
	ourDHCommitAKE := fixtureConversation()
	ourDHCommitAKE.dhCommitMessage()

	c := newConversation(otrV3{}, fixtureRand())
	c.initAKE()
	c.ake.theirPublicValue = ourDHCommitAKE.ake.ourPublicValue

	_, _, err := authStateAwaitingRevealSig{}.receiveDHCommitMessage(c, []byte{0x00, 0x00})
	assertEquals(t, err, newOtrError("corrupt DH commit message"))
}

func Test_authStateNone_receiveDHCommitMessage_returnsErrorIfgenerateCommitMsgInstanceTagsFails(t *testing.T) {
	ourDHCommitAKE := fixtureConversation()
	ourDHCommitAKE.dhCommitMessage()

	c := newConversation(otrV3{}, fixtureRand())
	c.initAKE()
	c.ake.theirPublicValue = ourDHCommitAKE.ake.ourPublicValue

	_, _, err := authStateNone{}.receiveDHCommitMessage(c, []byte{0x00, 0x00})
	assertEquals(t, err, newOtrError("corrupt DH commit message"))
}

func Test_authStateNone_receiveDHCommitMessage_returnsErrorIfdhKeyMessageFails(t *testing.T) {
	ourDHCommitAKE := fixtureConversation()
	ourDHCommitAKE.dhCommitMessage()

	c := newConversation(otrV2{}, fixedRand([]string{"ABCD"}))
	c.initAKE()
	c.ake.theirPublicValue = ourDHCommitAKE.ake.ourPublicValue

	_, _, err := authStateNone{}.receiveDHCommitMessage(c, []byte{0x00, 0x00, 0x00})
	assertEquals(t, err, errShortRandomRead)
}

func Test_authStateNone_receiveDHCommitMessage_returnsErrorIfProcessDHCommitFails(t *testing.T) {
	ourDHCommitAKE := fixtureConversation()
	ourDHCommitAKE.dhCommitMessage()

	c := newConversation(otrV2{}, fixtureRand())
	c.initAKE()
	c.ake.theirPublicValue = ourDHCommitAKE.ake.ourPublicValue

	_, _, err := authStateNone{}.receiveDHCommitMessage(c, []byte{0x00, 0x00})
	assertEquals(t, err, newOtrError("corrupt DH commit message"))
}

func Test_authStateAwaitingDHKey_receiveDHCommitMessage_failsIfMsgDoesntHaveHeader(t *testing.T) {
	ourDHCommitAKE := fixtureConversation()
	ourDHCommitAKE.dhCommitMessage()

	c := newConversation(otrV2{}, fixtureRand())
	c.initAKE()
	c.ake.theirPublicValue = ourDHCommitAKE.ake.ourPublicValue

	_, _, err := authStateAwaitingDHKey{}.receiveDHCommitMessage(c, []byte{0x00, 0x00})
	assertEquals(t, err, errInvalidOTRMessage)
}

func Test_authStateAwaitingDHKey_receiveDHCommitMessage_failsIfCantExtractFirstPart(t *testing.T) {
	ourDHCommitAKE := fixtureConversation()
	ourDHCommitAKE.dhCommitMessage()

	c := newConversation(otrV2{}, fixtureRand())
	c.initAKE()
	c.ake.theirPublicValue = ourDHCommitAKE.ake.ourPublicValue

	_, _, err := authStateAwaitingDHKey{}.receiveDHCommitMessage(c, []byte{0x00, 0x00, 0x00, 0x01})
	assertEquals(t, err, errInvalidOTRMessage)
}

func Test_authStateAwaitingDHKey_receiveDHCommitMessage_failsIfCantExtractSecondPart(t *testing.T) {
	ourDHCommitAKE := fixtureConversation()
	ourDHCommitAKE.dhCommitMessage()

	c := newConversation(otrV2{}, fixtureRand())
	c.initAKE()
	c.ake.theirPublicValue = ourDHCommitAKE.ake.ourPublicValue

	_, _, err := authStateAwaitingDHKey{}.receiveDHCommitMessage(c, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x01, 0x00, 0x01, 0x02})
	assertEquals(t, err, errInvalidOTRMessage)
}

func Test_authStateNone_String_returnsTheCorrectString(t *testing.T) {
	assertEquals(t, authStateNone{}.String(), "AUTHSTATE_NONE")
}

func Test_authStateAwaitingDHKey_String_returnsTheCorrectString(t *testing.T) {
	assertEquals(t, authStateAwaitingDHKey{}.String(), "AUTHSTATE_AWAITING_DHKEY")
}

func Test_authStateAwaitingRevealSig_String_returnsTheCorrectString(t *testing.T) {
	assertEquals(t, authStateAwaitingRevealSig{}.String(), "AUTHSTATE_AWAITING_REVEALSIG")
}

func Test_authStateAwaitingSig_String_returnsTheCorrectString(t *testing.T) {
	assertEquals(t, authStateAwaitingSig{}.String(), "AUTHSTATE_AWAITING_SIG")
}

func Test_akeHasFinished_willSignalThatWeAreTalkingToOurselvesIfWeAre(t *testing.T) {
	c := bobContextAfterAKE()
	c.ourCurrentKey = bobPrivateKey
	c.theirKey = bobPrivateKey.PublicKey()

	c.expectMessageEvent(t, func() {
		c.akeHasFinished()
	}, MessageEventMessageReflected, nil, nil)
}

func Test_akeHasFinished_willSignalThatWeHaveGoneSecureIfWeHave(t *testing.T) {
	c := bobContextAfterAKE()
	c.ourCurrentKey = bobPrivateKey
	c.theirKey = alicePrivateKey.PublicKey()
	c.msgState = plainText

	c.expectSecurityEvent(t, func() {
		c.akeHasFinished()
	}, GoneSecure)
}

func Test_akeHasFinished_willSignalThatWeHaveGoneSecureIfWeWereFinished(t *testing.T) {
	c := bobContextAfterAKE()
	c.ourCurrentKey = bobPrivateKey
	c.theirKey = alicePrivateKey.PublicKey()
	c.msgState = plainText

	c.expectSecurityEvent(t, func() {
		c.akeHasFinished()
	}, GoneSecure)
}

func Test_akeHasFinished_willSignalThatWeHaveGoneSecureIfWeHaveRefreshed(t *testing.T) {
	c := bobContextAfterAKE()
	c.ourCurrentKey = bobPrivateKey
	c.theirKey = alicePrivateKey.PublicKey()
	c.msgState = encrypted

	c.expectSecurityEvent(t, func() {
		c.akeHasFinished()
	}, StillSecure)
}

func Test_akeHasFinished_wipesAKEKeys(t *testing.T) {
	c := &Conversation{}
	c.ourCurrentKey = bobPrivateKey
	c.theirKey = bobPrivateKey.PublicKey()

	revKey := akeKeys{
		c:  fixedSize(16, []byte{1, 2, 3}),
		m1: fixedSize(32, []byte{4, 5, 6}),
		m2: fixedSize(32, []byte{7, 8, 9}),
	}

	sigKey := akeKeys{
		c:  fixedSize(16, []byte{3, 2, 1}),
		m1: fixedSize(32, []byte{6, 5, 4}),
		m2: fixedSize(32, []byte{9, 8, 7}),
	}

	c.ake = &ake{
		secretExponent:   big.NewInt(1),
		ourPublicValue:   big.NewInt(2),
		theirPublicValue: big.NewInt(2),
		revealKey:        revKey,
		sigKey:           sigKey,
		r:                [16]byte{1, 2, 3},
		encryptedGx:      []byte{1, 2, 3},
		xhashedGx:        fixedSize(otrV3{}.hash2Length(), []byte{1, 2, 3}),
		state:            authStateNone{},
	}

	c.akeHasFinished()

	assertDeepEquals(t, *c.ake, ake{state: c.ake.state})
}
