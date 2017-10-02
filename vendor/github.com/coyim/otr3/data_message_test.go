package otr3

import (
	"crypto/rand"
	"testing"
	"time"
)

func Test_processTLVs_ignoresInvalidTLVMessageTypes(t *testing.T) {
	var nilT []tlv
	tlvs := []tlv{
		fixtureMessage1().tlv(),
		tlv{
			tlvType:   9,
			tlvLength: 1,
			tlvValue:  []byte{0x01},
		},
	}

	c := newConversation(otrV3{}, fixtureRand())
	c.msgState = encrypted

	toSend, err := c.processTLVs(tlvs, dataMessageExtra{})
	assertEquals(t, err, nil)
	assertDeepEquals(t, toSend, nilT)
}

func Test_processTLVs_ignoresPaddingTLV(t *testing.T) {
	var nilB []tlv

	aTLV := tlv{
		tlvType:   0,
		tlvLength: 1,
		tlvValue:  []byte{0x00},
	}

	c := newConversation(otrV3{}, fixtureRand())
	c.msgState = encrypted

	tlvs, err := c.processTLVs([]tlv{aTLV}, dataMessageExtra{})
	assertDeepEquals(t, err, nil)
	assertDeepEquals(t, tlvs, nilB)
}

func Test_genDataMsg_withKeyExchangeData(t *testing.T) {
	c := bobContextAfterAKE()
	c.msgState = encrypted
	c.keys.ourKeyID = 2
	c.keys.theirKeyID = 3
	counter := c.keys.counterHistory.findCounterFor(c.keys.ourKeyID-1, c.keys.theirKeyID)
	counter.ourCounter = 0x1011121314

	dataMsg, _, err := c.genDataMsg(nil)

	assertEquals(t, err, nil)
	assertEquals(t, dataMsg.senderKeyID, uint32(1))
	assertEquals(t, dataMsg.recipientKeyID, uint32(3))
	assertDeepEquals(t, dataMsg.y, c.keys.ourCurrentDHKeys.pub)
	assertDeepEquals(t, dataMsg.topHalfCtr, [8]byte{
		0x00, 0x00, 0x00, 0x10, 0x11, 0x12, 0x13, 0x14,
	})
	assertEquals(t, counter.ourCounter, uint64(0x1011121314+1))
}

func Test_genDataMsg_willResetMayRetransmit(t *testing.T) {
	c := bobContextAfterAKE()
	c.msgState = encrypted
	c.keys.ourKeyID = 2
	c.keys.theirKeyID = 3

	counter := c.keys.counterHistory.findCounterFor(c.keys.ourKeyID-1, c.keys.theirKeyID)
	counter.ourCounter = 0x1011121314

	c.resend.mayRetransmit = retransmitExact

	c.genDataMsg(nil)

	assertEquals(t, c.resend.mayRetransmit, noRetransmit)
}

func Test_genDataMsg_willNotResetMayRetransmitIfItEncountersAnError(t *testing.T) {
	c := newConversation(otrV3{}, rand.Reader)
	c.msgState = encrypted
	c.Rand = fixedRand([]string{})
	c.ourInstanceTag = 0
	c.resend.mayRetransmit = retransmitExact

	c.genDataMsg(nil)

	assertEquals(t, c.resend.mayRetransmit, retransmitExact)
}

func Test_genDataMsg_setsLastMessageWhenNewMessageIsPlaintext(t *testing.T) {
	msg := []byte("hello")
	c := bobContextAfterAKE()
	c.msgState = encrypted
	c.keys.ourKeyID = 2
	c.keys.theirKeyID = 3

	counter := c.keys.counterHistory.findCounterFor(c.keys.ourKeyID-1, c.keys.theirKeyID)
	counter.ourCounter = 0x1011121314

	c.genDataMsg(msg)

	assertDeepEquals(t, c.resend.pending(),
		[]messageToResend{
			messageToResend{MessagePlaintext(msg), nil},
		})
}

func Test_genDataMsg_hasEncryptedMessage(t *testing.T) {
	c := bobContextAfterAKE()
	c.msgState = encrypted

	expected := bytesFromHex("a9a31c32ce6e7ae3c3956401ddc6aed6da24cb75dc16c3c473f2b14c2dda1cd52bfd20559eedf51d275b049fdefd93af2325d28d5d2f0fb05e8524842e32d4275c69a621e5fa133977563345055fded5511a78337a6d9a213bc5319de11a578818c2edb21b510595157feea3ed93a1178021571aa21765fd974c89cdcbda8ec0afce0c0ea5901021657b959f842df47224edd5dd50d9e736ed8982580373dcd0e2f06a5421472ae2bc58cc4ea7cb2b054e22c1781b72595909b37640e28f435df98b16410c76969fa9112a114b4ab7fb5b3265aa5efa0a99b9c47097d6d42a232a223d03b7d4a8fd5e57a748d1e06ef106e265f70421b708ca85b89e92f02082")
	dataMsg, _, err := c.genDataMsg([]byte("we are awesome"))

	assertEquals(t, err, nil)
	assertDeepEquals(t, dataMsg.encryptedMsg, expected)
}

func Test_genDataMsg_revealOldMACKeysFromKeyManagementContext(t *testing.T) {
	oldMACKeys := []macKey{
		macKey{0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03},
		macKey{0x01, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03},
	}

	c := bobContextAfterAKE()
	c.msgState = encrypted
	c.keys.oldMACKeys = oldMACKeys

	dataMsg, _, err := c.genDataMsg(nil)

	assertEquals(t, err, nil)
	assertDeepEquals(t, dataMsg.oldMACKeys, oldMACKeys)
}

func Test_genDataMsg_returnsErrorIfFailsToCalculateDHSessionKey(t *testing.T) {
	c := newConversation(otrV3{}, rand.Reader)
	c.msgState = encrypted
	_, _, err := c.genDataMsg(nil)
	assertEquals(t, err, newOtrConflictError("invalid key id for local peer"))
}

func Test_genDataMsg_returnsErrorIfFailsToGenerateInstanceTag(t *testing.T) {
	c := bobContextAfterAKE()
	c.msgState = encrypted
	c.Rand = fixedRand([]string{})
	c.ourInstanceTag = 0

	_, _, err := c.genDataMsg(nil)
	assertEquals(t, err, errShortRandomRead)
}

func Test_processDataMessage_deserializeAndDecryptDataMsg(t *testing.T) {
	bob := newConversation(otrV3{}, rand.Reader)
	bob.msgState = encrypted
	bob.Policies.add(allowV3)
	bob.ourCurrentKey = bobPrivateKey
	bob.smp.secret = bnFromHex("ABCDE56321F9A9F8E364607C8C82DECD8E8E6209E2CB952C7E649620F5286FE3")

	plain := plainDataMsg{
		message: []byte("hello"),
	}

	var msg []byte
	msg, bob.keys = fixtureDataMsg(plain)

	bob.msgState = encrypted
	exp, _, err := bob.receiveDecoded(msg)

	assertNil(t, err)
	assertDeepEquals(t, exp, MessagePlaintext("hello"))
}

func Test_processDataMessage_willGenerateAHeartBeatEventForAnEmptyMessage(t *testing.T) {
	bob := newConversation(otrV3{}, rand.Reader)
	bob.Policies.add(allowV3)
	bob.ourCurrentKey = bobPrivateKey
	bob.smp.secret = bnFromHex("ABCDE56321F9A9F8E364607C8C82DECD8E8E6209E2CB952C7E649620F5286FE3")

	plain := plainDataMsg{
		message: []byte(""),
	}

	var msg []byte
	msg, bob.keys = fixtureDataMsg(plain)

	bob.msgState = encrypted

	bob.expectMessageEvent(t, func() {
		exp, _, _ := bob.receiveDecoded(msg)
		assertNil(t, exp)
	}, MessageEventLogHeartbeatReceived, nil, nil)
}

func Test_processDataMessage_processSMPMessage(t *testing.T) {
	bob := newConversation(otrV3{}, rand.Reader)
	bob.Policies.add(allowV3)
	bob.ourCurrentKey = bobPrivateKey

	bob.smp.state = smpStateExpect2{}
	bob.smp.s1 = fixtureSmp1()
	bob.smp.secret = bnFromHex("ABCDE56321F9A9F8E364607C8C82DECD8E8E6209E2CB952C7E649620F5286FE3")

	plain := plainDataMsg{
		tlvs: []tlv{
			fixtureMessage2().tlv(),
		},
	}

	var msg []byte
	msg, bob.keys = fixtureDataMsg(plain)

	bob.msgState = encrypted
	_, toSend, err := bob.receiveDecoded(msg)

	exp := fixtureDecryptDataMsg(toSend[0])

	assertDeepEquals(t, err, nil)
	assertDeepEquals(t, len(exp.tlvs), 2)
	assertDeepEquals(t, exp.tlvs[0].tlvType, uint16(tlvTypeSMP3))
	assertDeepEquals(t, exp.tlvs[1].tlvType, uint16(tlvTypePadding))
}

func Test_processDataMessage_returnsErrorIfSomethingGoesWrongWithDeserialize(t *testing.T) {
	c := newConversation(otrV3{}, rand.Reader)
	c.msgState = encrypted
	_, _, err := c.processDataMessage([]byte{}, []byte{})

	assertEquals(t, err.Error(), "otr: dataMsg.deserialize empty message")
}

func Test_processDataMessage_returnsErrorIfDataMessageHasWrongCounter(t *testing.T) {
	c := newConversation(otrV3{}, rand.Reader)
	c.ourCurrentKey = bobPrivateKey

	var msg []byte
	msg, c.keys = fixtureDataMsg(plainDataMsg{})

	ctr := c.keys.counterHistory.findCounterFor(1, 1)
	ctr.theirCounter = 2 // force a bigger counter

	c.msgState = encrypted
	_, _, err := c.receiveDecoded(msg)

	assertDeepEquals(t, err, newOtrConflictError("counter regressed"))
}

func Test_processDataMessage_signalsThatMessageIsUnreadableForAGPGConflictError(t *testing.T) {
	c := newConversation(otrV3{}, rand.Reader)
	c.ourCurrentKey = bobPrivateKey

	var msg []byte
	msg, c.keys = fixtureDataMsg(plainDataMsg{})

	ctr := c.keys.counterHistory.findCounterFor(1, 1)
	ctr.theirCounter = 2 // force a bigger counter

	c.msgState = encrypted

	c.expectMessageEvent(t, func() {
		c.receiveDecoded(msg)
	}, MessageEventReceivedMessageUnreadable, nil, nil)
}

func Test_Receive_returnsACustomErrorMessageIfOneIsAvailable(t *testing.T) {
	c := newConversation(otrV3{}, rand.Reader)
	c.ourCurrentKey = bobPrivateKey

	var msg []byte
	msg, c.keys = fixtureDataMsg(plainDataMsg{})

	ctr := c.keys.counterHistory.findCounterFor(1, 1)
	ctr.theirCounter = 2 // force a bigger counter

	c.msgState = encrypted

	c.errorMessageHandler = dynamicErrorMessageHandler{
		func(error ErrorCode) []byte {
			if error == ErrorCodeMessageUnreadable {
				return []byte("nova happened")
			}
			return []byte("white hole happened")
		}}
	c.receiveDecoded(msg)
	ts, _ := c.withInjections(nil, nil)
	assertDeepEquals(t, string(ts[0]), "?OTR Error: nova happened")
}

func Test_processDataMessage_signalsThatMessageIsMalformedIfSomeOtherErrorHappens(t *testing.T) {
	c := newConversation(otrV3{}, fixedRand([]string{"ABCD"}))
	c.ourCurrentKey = bobPrivateKey
	var msg []byte
	msg, c.keys = fixtureDataMsg(plainDataMsg{message: []byte("Making sure this isn't a heartbeat message")})
	c.msgState = encrypted
	c.keys.ourKeyID = 1
	c.expectMessageEvent(t, func() {
		c.receiveDecoded(msg)
	}, MessageEventReceivedMessageMalformed, nil, nil)
}

func Test_processDataMessage_callsErrorMessageHandlerAndReturnsTheResultAsAnOTRErrorMessageForAnError(t *testing.T) {
	c := newConversation(otrV3{}, fixedRand([]string{"ABCD"}))
	c.ourCurrentKey = bobPrivateKey
	var msg []byte
	msg, c.keys = fixtureDataMsg(plainDataMsg{message: []byte("Making sure this isn't a heartbeat message")})
	c.msgState = encrypted

	c.errorMessageHandler = dynamicErrorMessageHandler{
		func(error ErrorCode) []byte {
			if error == ErrorCodeMessageMalformed {
				return []byte("sunflower happened")
			}
			return []byte("dandelion happened")
		}}

	c.keys.ourKeyID = 1
	c.receiveDecoded(msg)
	ts, _ := c.withInjections(nil, nil)
	assertDeepEquals(t, string(ts[0]), "?OTR Error: sunflower happened")
}

func Test_processDataMessage_shouldNotRotateKeysWhenDecryptFails(t *testing.T) {
	bob := newConversation(otrV3{}, rand.Reader)
	bob.Policies.add(allowV3)
	bob.ourCurrentKey = bobPrivateKey

	var msg []byte
	msg, bob.keys = fixtureDataMsg(plainDataMsg{})

	bob.keys.ourKeyID = 1 //force key rotation
	msg[len(msg)-8] = 0   //force check sign failure
	bobCurrentDHKeys := bob.keys.ourCurrentDHKeys
	bobPreviousDHKeys := bob.keys.ourPreviousDHKeys
	aliceCurrentKey := bob.keys.theirCurrentDHPubKey
	alicePreviousKey := bob.keys.theirPreviousDHPubKey

	bob.msgState = encrypted
	_, _, err := bob.receiveDecoded(msg)

	assertDeepEquals(t, err, newOtrConflictError("bad signature MAC in encrypted signature"))
	assertDeepEquals(t, bobCurrentDHKeys, bob.keys.ourCurrentDHKeys)
	assertDeepEquals(t, bobPreviousDHKeys, bob.keys.ourPreviousDHKeys)

	assertDeepEquals(t, aliceCurrentKey, bob.keys.theirCurrentDHPubKey)
	assertDeepEquals(t, alicePreviousKey, bob.keys.theirPreviousDHPubKey)
}

func Test_processDataMessage_rotateOurKeysAfterDecryptingTheMessage(t *testing.T) {
	bob := newConversation(otrV3{}, rand.Reader)
	bob.Policies.add(allowV3)
	bob.ourCurrentKey = bobPrivateKey

	var msg []byte
	msg, bob.keys = fixtureDataMsg(plainDataMsg{})

	bob.keys.ourKeyID = 1
	bobCurrentDHKeys := bob.keys.ourCurrentDHKeys

	bob.msgState = encrypted
	_, toSend, err := bob.receiveDecoded(msg)

	assertNil(t, err)
	assertNil(t, toSend)
	assertDeepEquals(t, bobCurrentDHKeys, bob.keys.ourPreviousDHKeys)
	assertEquals(t, eq(bobCurrentDHKeys.pub, bob.keys.ourCurrentDHKeys.pub), false)
	assertEquals(t, eq(bobCurrentDHKeys.priv, bob.keys.ourCurrentDHKeys.priv), false)
}

func Test_processDataMessage_willReturnAHeartbeatMessageAfterAPlainTextMessage(t *testing.T) {
	bob := newConversation(otrV3{}, rand.Reader)
	bob.Policies.add(allowV3)
	bob.ourCurrentKey = bobPrivateKey
	bob.heartbeat.lastSent = time.Now().Add(-61 * time.Second)

	var msg []byte
	msg, bob.keys = fixtureDataMsg(plainDataMsg{message: []byte("something")})

	bob.keys.ourKeyID = 1

	bob.msgState = encrypted
	_, toSend, err := bob.receiveDecoded(msg)

	assertDeepEquals(t, err, nil)

	header, exp, e := fixtureDecryptDataMsgBase(toSend[0])
	assertNil(t, e)
	assertDeepEquals(t, header, bytesFromHex("0003030000010100000101"))
	assertDeepEquals(t, exp.message, []byte{})
	assertDeepEquals(t, len(exp.tlvs), 1)
	assertDeepEquals(t, exp.tlvs[0].tlvType, tlvTypePadding)
	assertDeepEquals(t, exp.tlvs[0].tlvLength, uint16(0xFB))
}

func Test_processDataMessage_rotateTheirKeysAfterDecryptingTheMessage(t *testing.T) {
	bob := newConversation(otrV3{}, rand.Reader)
	bob.Policies.add(allowV3)
	bob.ourCurrentKey = bobPrivateKey

	var msg []byte
	msg, bob.keys = fixtureDataMsg(plainDataMsg{})

	bob.keys.theirKeyID = 1
	aliceCurrentDHPubKey := bob.keys.theirCurrentDHPubKey

	bob.msgState = encrypted
	_, toSend, err := bob.receiveDecoded(msg)

	assertNil(t, err)
	assertNil(t, toSend)
	assertDeepEquals(t, aliceCurrentDHPubKey, bob.keys.theirPreviousDHPubKey)
	assertEquals(t, eq(aliceCurrentDHPubKey, bob.keys.theirCurrentDHPubKey), false)
}

func Test_processDataMessage_ignoresTLVsWhenFailsToRotateKeys(t *testing.T) {
	bob := newConversation(otrV3{}, fixedRand([]string{}))
	bob.Policies.add(allowV3)
	bob.ourCurrentKey = bobPrivateKey

	// setup state for receiving a SMP message 2
	bob.smp.state = smpStateExpect2{}
	bob.smp.s1 = fixtureSmp1()
	bob.smp.secret = bnFromHex("ABCDE56321F9A9F8E364607C8C82DECD8E8E6209E2CB952C7E649620F5286FE3")

	var msg []byte
	plain := plainDataMsg{}
	plain.tlvs = append(plain.tlvs, fixtureMessage2().tlv())
	msg, bob.keys = fixtureDataMsg(plain)

	bob.keys.theirKeyID = 1 //forces our key rotation
	bobCurrentDHKeys := bob.keys.ourCurrentDHKeys

	bob.msgState = encrypted
	bob.smp.state = smpStateExpect1{}
	bob.keys.ourKeyID = 1
	_, toSend, err := bob.receiveDecoded(msg)

	assertDeepEquals(t, err, errShortRandomRead)
	assertNil(t, toSend)
	assertDeepEquals(t, bobCurrentDHKeys, bob.keys.ourCurrentDHKeys)
}

func Test_processDataMessage_returnErrorWhenOurKeyIDUnexpected(t *testing.T) {
	datamsg := bytesFromHex("0003030000010100000101000000000100000001000000c03a3ca02c03bef84c7596504b7b2dee2820500bf51107e4447cfd2fddd8132a29668ef7cb3f56ff75f80e9d5a3c34e4aaa45a63beee83c058d21653e45d56ad04f6493545ad5bc3441f9a1a23fdf5ea0d812f3dfa02de9742ee9b1779dd1d84bf1bf06700a05779ff1a730c51ecdce34d251317dacdcbe865f12c2bf8e4a8a15cc10975184a7509e3f82244c8594d3df18b411648dc059cf341c50ab0d3981f186519ca3104609e89a5f4be44047068c5ba33d2b1de0e9b7d5e6aa67c148f57d70000000000000001000001007104b8684860d2eacc0d653ca9696171f5d7b03d90a06fd46305c041ab4af8313826ca82f8fc43c755c56dd62fa025822e72d9566a32fe88f189e0fb1b07128a37db49350392470cdd57f280f565ab775d58af6f5d8efca39126192efefe1f98bdfd2135b1c6ce8e68d8d3bfd50eae34187191524492193d20dd75d6b04a1e7d90fe1e71a9843b720df310119c1db82928c11308d93ed508641e73b6d579eefbcb432ab2ebf2b15a3b1c8baca86d5008c81286705b9368abec0d5cf4b6e2289be1040b5ac172cbc81f7a594d721cafd50e7cfdc2616c6d59cf445f885d8e80980a73f6a55a34be9e90b7ec25f757e212fa2b79c4c56d922a804168bfeca75199dbede31d8101018586d1f992afdd80117cf84d1000000000")
	bob := newConversation(otrV3{}, rand.Reader)
	bob.Policies.add(allowV2)
	bob.Policies.add(allowV3)
	bob.ourCurrentKey = bobPrivateKey
	bob.theirKey = alicePrivateKey.PublicKey()
	bob.keys.ourKeyID = 3
	bob.keys.theirKeyID = 1
	bob.keys.ourPreviousDHKeys.priv = bnFromHex("28cea443a1ddeae5c39fd9061a429243eeb52f9f963dcb483a77ec9ed201f8eb3e898fb645657f27")
	bob.keys.ourPreviousDHKeys.pub = bnFromHex("e291f2e06da00d59c9666d80d6c511a0bd9ae54d916b65db7e72f70904ae05d55259df42fb7b29d11babf11e78cd584d0f137ca1187b4f920e0fbef85c0e5f4b55bf907ea6e119dcfa7e339e72d6b52e874dc46afedd9290360659928ad30f504dad43160946dbd9de7748d18417c223790e528a6f13bf25285318416ccfed0bceafbca70dce832ca8216a654c49ac29dc6af098e7e2744a1dfaf7d2643eb1b3787c4c1db4f649096c3241f69165f965a290651304e23fd2422dae180796d52f")
	bob.keys.theirCurrentDHPubKey = bnFromHex("da61b77be39426456fecfd6df16645bd2c967bc1a27b165dbf77fea4753ece7a8b938532395bbd1def2890a2792f1854c2d736ee27139356b3bb2583afa4c96a9083209d9f2bb1caeb6fe5ee608715ae6dc1c470e38b895e48e0532af5388c8e591d9ebe361f118ad54d8640f24fa54fdb1d07594d496150554094e5ec4bcfcc6b1b4b058b679824306ad7ae481a25d0758cc01c29c281ce33ac2f58d6eaa99985f855e9ce667ff287b4d27d7c73a7717277546d17e8dd5539861bc26fa04c1b")

	bob.msgState = encrypted
	_, _, err := bob.receiveDecoded(datamsg)

	assertDeepEquals(t, err, newOtrConflictError("mismatched key id for local peer"))
}
