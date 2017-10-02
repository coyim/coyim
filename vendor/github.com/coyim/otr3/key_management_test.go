package otr3

import (
	"math/big"
	"testing"
)

func Test_calculateDHSessionKeys(t *testing.T) {
	c := keyManagementContext{
		ourKeyID:   1,
		theirKeyID: 2,
		ourCurrentDHKeys: dhKeyPair{
			pub:  fixedGX(),
			priv: fixedX(),
		},
		theirPreviousDHPubKey: fixedGY(),
	}

	sendingAESKey := bytesFromHex("42e258bebf031acf442f52d6ef52d6f1")
	sendingMACKey := bytesFromHex("a45e2b122f58bbe2042f73f092329ad9b5dfe23e")
	receivingAESKey := bytesFromHex("c778c71cb63161e8e06d245e77ff6430")
	receivingMACKey := bytesFromHex("03f8034b891b1e843db5bba9a41ec68a1f5f8bbf")
	extraKey := bytesFromHex("0e1810c7c62c3bace6450dcbef16af8a271b5ac93030b83e9d0d80e0641e3c18")

	keys, err := c.calculateDHSessionKeys(1, 1, otrV3{})

	assertEquals(t, err, nil)
	assertDeepEquals(t, keys.sendingAESKey, sendingAESKey)
	assertDeepEquals(t, keys.sendingMACKey, macKey(sendingMACKey))
	assertDeepEquals(t, keys.receivingAESKey, receivingAESKey)
	assertDeepEquals(t, keys.receivingMACKey, macKey(receivingMACKey))
	assertDeepEquals(t, keys.extraKey, extraKey)
}

func Test_calculateDHSessionKeys_storesGeneratedMACKeys(t *testing.T) {
	ourKeyID := uint32(1)
	theirKeyID := uint32(2)

	c := keyManagementContext{
		ourKeyID:             ourKeyID,
		theirKeyID:           theirKeyID,
		theirCurrentDHPubKey: big.NewInt(1),
		ourCurrentDHKeys: dhKeyPair{
			priv: big.NewInt(1),
			pub:  big.NewInt(1),
		},
	}
	keys, _ := c.calculateDHSessionKeys(ourKeyID, theirKeyID, otrV3{})

	expectedMACKeys := macKeyUsage{
		ourKeyID:     ourKeyID,
		theirKeyID:   theirKeyID,
		receivingKey: keys.receivingMACKey,
	}

	assertDeepEquals(t, len(c.macKeyHistory.items), 1)
	assertDeepEquals(t, c.macKeyHistory.items[0], expectedMACKeys)
}

func Test_calculateDHSessionKeys_failsWhenOurOrTheyKeyIsUnknown(t *testing.T) {
	c := keyManagementContext{
		ourKeyID:   1,
		theirKeyID: 1,
	}

	_, err := c.calculateDHSessionKeys(2, 1, otrV3{})
	assertDeepEquals(t, err, newOtrConflictError("mismatched key id for local peer"))

	_, err = c.calculateDHSessionKeys(1, 3, otrV3{})
	assertDeepEquals(t, err, newOtrConflictError("mismatched key id for remote peer"))
}

func Test_calculateDHSessionKeys_failsWhenTheirPreviousPubliKeyIsNull(t *testing.T) {
	c := keyManagementContext{
		ourKeyID:   2,
		theirKeyID: 2,
	}
	_, err := c.calculateDHSessionKeys(2, 1, otrV3{})

	assertEquals(t, err, newOtrConflictError("no previous key for remote peer found"))
}

func Test_pickTheirKey_shouldFailsForInvalidSenderID(t *testing.T) {
	c := keyManagementContext{}

	c.theirKeyID = uint32(0)
	_, err := c.pickTheirKey(uint32(0x00000000))
	assertEquals(t, err, newOtrConflictError("invalid key id for remote peer"))

	c.theirKeyID = uint32(2)
	_, err = c.pickTheirKey(uint32(0x00000000))
	assertEquals(t, err, newOtrConflictError("invalid key id for remote peer"))

	c.theirKeyID = uint32(1)
	c.theirPreviousDHPubKey = nil
	_, err = c.pickTheirKey(uint32(0x00000000))
	assertEquals(t, err, newOtrConflictError("invalid key id for remote peer"))
}

func Test_pickOurKeys_shouldFailsForInvalidRecipientID(t *testing.T) {
	c := keyManagementContext{}

	c.ourKeyID = uint32(0x00000000)
	_, _, err := c.pickOurKeys(uint32(0x00000000))
	assertEquals(t, err, newOtrConflictError("invalid key id for local peer"))

	c.ourKeyID = uint32(3)
	_, _, err = c.pickOurKeys(uint32(0x00000001))
	assertEquals(t, err, newOtrConflictError("mismatched key id for local peer"))
}

func Test_calculateAKEKeys(t *testing.T) {
	ssid, revealSigKeys, signatureKeys := calculateAKEKeys(expectedSharedSecret, otrV3{})

	assertDeepEquals(t, ssid[:], bytesFromHex("9cee5d2c7edbc86d"))
	assertDeepEquals(t, revealSigKeys.c, bytesFromHex("5745340b350364a02a0ac1467a318dcc"))
	assertDeepEquals(t, signatureKeys.c, bytesFromHex("d942cc80b66503414c05e3752d9ba5c4"))
	assertDeepEquals(t, revealSigKeys.m1, bytesFromHex("d3251498fb9d977d07392a96eafb8c048d6bc67064bd7da72aa38f20f87a2e3d"))
	assertDeepEquals(t, revealSigKeys.m2, bytesFromHex("79c101a78a6c5819547a36b4813c84a8ac553d27a5d4b58be45dd0f3a67d3ca6"))
	assertDeepEquals(t, signatureKeys.m1, bytesFromHex("b6254b8eab0ad98152949454d23c8c9b08e4e9cf423b27edc09b1975a76eb59c"))
	assertDeepEquals(t, signatureKeys.m2, bytesFromHex("954be27015eeb0455250144d906e83e7d329c49581aea634c4189a3c981184f5"))
}

func Test_rotateTheirKey_rotatesTheirKeysWhenWeReceiveANewPubKey(t *testing.T) {
	senderKey := uint32(1)
	currentPubKey := big.NewInt(9)
	receivedKey := big.NewInt(99)

	c := keyManagementContext{
		theirKeyID:           senderKey,
		theirCurrentDHPubKey: currentPubKey,
	}

	c.rotateTheirKey(senderKey, receivedKey)

	assertEquals(t, c.theirKeyID, senderKey+1)
	assertDeepEquals(t, c.theirPreviousDHPubKey, currentPubKey)
	assertDeepEquals(t, c.theirCurrentDHPubKey, receivedKey)
}

func Test_rotateTheirKey_doesNotRotateIfWeDontReceiveTheCurrentSenderKey(t *testing.T) {
	senderKey := uint32(1)
	previousPubKey := big.NewInt(8)
	currentPubKey := big.NewInt(9)
	receivedKey := big.NewInt(99)

	c := keyManagementContext{
		theirKeyID:            senderKey,
		theirPreviousDHPubKey: previousPubKey,
		theirCurrentDHPubKey:  currentPubKey,
	}

	c.rotateTheirKey(senderKey+1, receivedKey)

	assertEquals(t, c.theirKeyID, senderKey)
	assertDeepEquals(t, c.theirPreviousDHPubKey, previousPubKey)
	assertDeepEquals(t, c.theirCurrentDHPubKey, currentPubKey)
}

func Test_rotateOurKeys_rotateOurCurrentDHKeys(t *testing.T) {
	recipientKeyID := uint32(1)

	c := keyManagementContext{
		ourKeyID: recipientKeyID,
		ourCurrentDHKeys: dhKeyPair{
			pub:  fixedGX(),
			priv: fixedX(),
		},
	}

	c.rotateOurKeys(recipientKeyID, fixedRand([]string{"abcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcd"}))

	assertEquals(t, c.ourKeyID, recipientKeyID+1)
	assertDeepEquals(t, c.ourPreviousDHKeys.priv, fixedX())
	assertDeepEquals(t, c.ourPreviousDHKeys.pub, fixedGX())
	assertDeepEquals(t, c.ourCurrentDHKeys.priv, fixedY())
	assertDeepEquals(t, c.ourCurrentDHKeys.pub, fixedGY())
}

func Test_rotateOurKeys_doesNotRotateIfWeDontReceiveOurCurrentKeyID(t *testing.T) {
	var nilB *big.Int
	recipientKeyID := uint32(1)

	c := keyManagementContext{
		ourKeyID: recipientKeyID,
		ourCurrentDHKeys: dhKeyPair{
			pub:  fixedGX(),
			priv: fixedX(),
		},
	}

	c.rotateOurKeys(recipientKeyID+1, fixedRand([]string{"abcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcd"}))

	assertEquals(t, c.ourKeyID, recipientKeyID)
	assertEquals(t, c.ourPreviousDHKeys.priv, nilB)
	assertEquals(t, c.ourPreviousDHKeys.pub, nilB)
	assertDeepEquals(t, c.ourCurrentDHKeys.priv, fixedX())
	assertDeepEquals(t, c.ourCurrentDHKeys.pub, fixedGX())
}

func Test_revealMACKeys_ForgotOldKeysAfterBeenCalled(t *testing.T) {
	oldMACKeys := []macKey{
		macKey{0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03},
	}
	c := keyManagementContext{
		oldMACKeys: oldMACKeys,
	}

	maKeys := c.revealMACKeys()

	assertDeepEquals(t, maKeys, oldMACKeys)
	assertDeepEquals(t, c.oldMACKeys, []macKey{})
}

func Test_rotateTheirKey_revealAllMACKeysAssociatedWithTheirPreviousPubKey(t *testing.T) {
	k1 := macKey{0x01, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	k2 := macKey{0x02, 0x02, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	expectedMACKeys := []macKey{
		k2, k2,
	}

	c := keyManagementContext{
		theirKeyID:            2,
		theirPreviousDHPubKey: big.NewInt(1),
	}

	c.macKeyHistory = macKeyHistory{
		items: []macKeyUsage{
			macKeyUsage{theirKeyID: 1, receivingKey: k2},
			macKeyUsage{theirKeyID: 2, receivingKey: k1},
			macKeyUsage{theirKeyID: 1, receivingKey: k2},
		},
	}

	c.rotateTheirKey(2, big.NewInt(2))

	assertDeepEquals(t, c.oldMACKeys, expectedMACKeys)
	assertDeepEquals(t, len(c.macKeyHistory.items), 1)
	assertDeepEquals(t, c.macKeyHistory.items[0].receivingKey, k1)
}

func Test_rotateOurKey_revealAllMACKeysAssociatedWithOurPreviousPubKey(t *testing.T) {
	k1 := macKey{0x01, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	k2 := macKey{0x02, 0x02, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	expectedMACKeys := []macKey{
		k2, k2,
	}

	c := keyManagementContext{
		ourKeyID: 2,
		ourPreviousDHKeys: dhKeyPair{
			priv: big.NewInt(1),
			pub:  big.NewInt(2),
		},
	}

	c.macKeyHistory = macKeyHistory{
		items: []macKeyUsage{
			macKeyUsage{ourKeyID: 1, receivingKey: k2},
			macKeyUsage{ourKeyID: 2, receivingKey: k1},
			macKeyUsage{ourKeyID: 1, receivingKey: k2},
		},
	}

	c.rotateOurKeys(2, fixedRand([]string{"abcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcd"}))

	assertDeepEquals(t, c.oldMACKeys, expectedMACKeys)
	assertDeepEquals(t, len(c.macKeyHistory.items), 1)
	assertDeepEquals(t, c.macKeyHistory.items[0].receivingKey, k1)
}

func Test_checkMessageCounter_messageIsInvalidWhenCounterIsNotLargerThanTheLastReceived(t *testing.T) {
	c := keyManagementContext{}
	ctr := c.counterHistory.findCounterFor(1, 1)
	ctr.theirCounter = 2

	msg := dataMsg{
		senderKeyID:    1,
		recipientKeyID: 1,
	}
	msg.topHalfCtr[7] = 2

	err := c.checkMessageCounter(msg)
	assertEquals(t, err, newOtrConflictError("counter regressed"))
	assertEquals(t, ctr.theirCounter, uint64(2))

	msg.topHalfCtr[7] = 1
	err = c.checkMessageCounter(msg)
	assertEquals(t, err, newOtrConflictError("counter regressed"))
	assertEquals(t, ctr.theirCounter, uint64(2))
}

func Test_checkMessageCounter_messageIsValidWhenCounterIsLargerThanTheLastReceived(t *testing.T) {
	c := keyManagementContext{}

	ctr := c.counterHistory.findCounterFor(1, 1)
	ctr.theirCounter = 2

	msg := dataMsg{
		senderKeyID:    1,
		recipientKeyID: 1,
	}
	msg.topHalfCtr[7] = 3
	err := c.checkMessageCounter(msg)
	assertEquals(t, err, nil)
	assertEquals(t, ctr.theirCounter, uint64(3))

	ctr = c.counterHistory.findCounterFor(2, 1)

	msg = dataMsg{
		senderKeyID:    1,
		recipientKeyID: 2,
	}

	msg.topHalfCtr[7] = 1
	err = c.checkMessageCounter(msg)
	assertEquals(t, err, nil)
	assertEquals(t, ctr.theirCounter, uint64(1))
}

func Test_generateNewDHKeypair_wipesPreviousDHKeysBeforePointingToCurrentDHKeys(t *testing.T) {
	prevPrivKey := big.NewInt(1)
	prevPubKey := big.NewInt(2)

	c := keyManagementContext{
		ourPreviousDHKeys: dhKeyPair{
			priv: prevPrivKey,
			pub:  prevPubKey,
		},
	}

	c.generateNewDHKeyPair(fixedRand([]string{"abcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcd"}))

	assertEquals(t, prevPrivKey.Int64(), int64(0))
	assertEquals(t, prevPubKey.Int64(), int64(0))
}
