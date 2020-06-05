package otr3

import (
	"encoding/binary"
	"hash"
	"io"
	"math/big"

	"github.com/coyim/constbn"
)

type dhKeyPair struct {
	pub  *big.Int
	priv secretKeyValue
}

type akeKeys struct {
	c      []byte
	m1, m2 []byte
}

// SIZE: this will always be the same size as the version.hash1Length
type macKey []byte

type sessionKeys struct {
	// SIZE: these should always be the same size as the AES used, usually 16
	sendingAESKey, receivingAESKey []byte
	sendingMACKey, receivingMACKey macKey
	// SIZE: this will be the same size as version.hash2Length
	extraKey []byte
}

type macKeyUsage struct {
	ourKeyID, theirKeyID uint32
	receivingKey         macKey
}

type macKeyHistory struct {
	items []macKeyUsage
}

func (h *macKeyHistory) deleteKeysAt(del ...int) {
	for j := len(del) - 1; j >= 0; j-- {
		l := len(h.items)
		h.items[del[j]], h.items = h.items[l-1], h.items[:l-1]
	}
}

func (h *macKeyHistory) addKeys(ourKeyID uint32, theirKeyID uint32, receivingMACKey macKey) {
	macKeys := macKeyUsage{
		ourKeyID:     ourKeyID,
		theirKeyID:   theirKeyID,
		receivingKey: receivingMACKey,
	}
	h.items = append(h.items, macKeys)
}

func (h *macKeyHistory) forgetMACKeysForOurKey(ourKeyID uint32) []macKey {
	var ret []macKey
	var del []int

	for i, k := range h.items {
		if k.ourKeyID == ourKeyID {
			ret = append(ret, k.receivingKey)
			del = append(del, i)
		}
	}

	h.deleteKeysAt(del...)

	return ret
}

func (h *macKeyHistory) forgetMACKeysForTheirKey(theirKeyID uint32) []macKey {
	var ret []macKey
	var del []int

	for i, k := range h.items {
		if k.theirKeyID == theirKeyID {
			ret = append(ret, k.receivingKey)
			del = append(del, i)
		}
	}

	h.deleteKeysAt(del...)

	return ret
}

type keyPairCounter struct {
	ourKeyID, theirKeyID     uint32
	ourCounter, theirCounter uint64
}

type counterHistory struct {
	counters []*keyPairCounter
}

func (h *counterHistory) findCounterFor(ourKeyID, theirKeyID uint32) *keyPairCounter {
	for _, c := range h.counters {
		if c.ourKeyID == ourKeyID && c.theirKeyID == theirKeyID {
			return c
		}
	}

	c := &keyPairCounter{
		ourKeyID:   ourKeyID,
		theirKeyID: theirKeyID,
	}

	h.counters = append(h.counters, c)
	return c
}

type keyManagementContext struct {
	ourKeyID, theirKeyID                        uint32
	ourCurrentDHKeys, ourPreviousDHKeys         dhKeyPair
	theirCurrentDHPubKey, theirPreviousDHPubKey *big.Int

	counterHistory counterHistory
	macKeyHistory  macKeyHistory
	oldMACKeys     []macKey
}

func (k *keyManagementContext) setTheirCurrentDHPubKey(key *big.Int) {
	k.theirCurrentDHPubKey = setBigInt(k.theirCurrentDHPubKey, key)
}

func (k *keyManagementContext) setOurCurrentDHKeys(priv secretKeyValue, pub *big.Int) {
	k.ourCurrentDHKeys.priv = setSecretKeyValue(k.ourCurrentDHKeys.priv, priv)
	k.ourCurrentDHKeys.pub = setBigInt(k.ourCurrentDHKeys.pub, pub)
}

func (k *keyManagementContext) checkMessageCounter(message dataMsg) error {
	counter := k.counterHistory.findCounterFor(message.recipientKeyID, message.senderKeyID)
	theirNextCounter := binary.BigEndian.Uint64(message.topHalfCtr[:])

	if theirNextCounter <= counter.theirCounter {
		return newOtrConflictError("counter regressed")
	}

	counter.theirCounter = theirNextCounter
	return nil
}

func (k *keyManagementContext) revealMACKeys() []macKey {
	ret := k.oldMACKeys
	k.oldMACKeys = []macKey{}
	return ret
}

func (k *keyManagementContext) generateNewDHKeyPair(randomness io.Reader) error {
	newPrivKey, err := randSizedSecret(randomness, 40)
	if err != nil {
		return err
	}

	tryLock(newPrivKey)

	k.ourPreviousDHKeys.wipe()
	k.ourPreviousDHKeys = k.ourCurrentDHKeys

	k.ourCurrentDHKeys = dhKeyPair{
		priv: newPrivKey,
		pub:  modExpPCT(g1ct, newPrivKey).GetBigInt(),
	}
	k.ourKeyID++
	return nil
}

func (k *keyManagementContext) revealMACKeysForOurPreviousKeyID() {
	keys := k.macKeyHistory.forgetMACKeysForOurKey(k.ourKeyID - 1)
	k.oldMACKeys = append(k.oldMACKeys, keys...)
}

func (c *Conversation) rotateKeys(dataMessage dataMsg) error {
	if err := c.keys.rotateOurKeys(dataMessage.recipientKeyID, c.rand()); err != nil {
		return err
	}
	c.keys.rotateTheirKey(dataMessage.senderKeyID, dataMessage.y)

	return nil
}

func (k *keyManagementContext) rotateOurKeys(recipientKeyID uint32, randomness io.Reader) error {
	if recipientKeyID == k.ourKeyID {
		k.revealMACKeysForOurPreviousKeyID()
		return k.generateNewDHKeyPair(randomness)
	}
	return nil
}

func (k *keyManagementContext) revealMACKeysForTheirPreviousKeyID() {
	keys := k.macKeyHistory.forgetMACKeysForTheirKey(k.theirKeyID - 1)
	k.oldMACKeys = append(k.oldMACKeys, keys...)
}

func (k *keyManagementContext) rotateTheirKey(senderKeyID uint32, pubDHKey *big.Int) {
	if senderKeyID == k.theirKeyID {
		k.revealMACKeysForTheirPreviousKeyID()

		k.theirPreviousDHPubKey = k.theirCurrentDHPubKey
		k.theirCurrentDHPubKey = pubDHKey
		k.theirKeyID++
	}
}

func (k *keyManagementContext) calculateDHSessionKeys(ourKeyID, theirKeyID uint32, v otrVersion) (sessionKeys, error) {
	var ret sessionKeys

	ourPrivKey, ourPubKey, err := k.pickOurKeys(ourKeyID)
	if err != nil {
		return ret, err
	}

	theirPubKey, err := k.pickTheirKey(theirKeyID)
	if err != nil {
		return ret, err
	}

	ret = calculateDHSessionKeys(ourPrivKey, ourPubKey, theirPubKey, v)
	k.macKeyHistory.addKeys(ourKeyID, theirKeyID, ret.receivingMACKey)

	return ret, nil
}

func calculateDHSessionKeys(ourPrivKey secretKeyValue, ourPubKey, theirPubKey *big.Int, v otrVersion) sessionKeys {
	var ret sessionKeys
	var sendbyte, recvbyte byte

	if gt(ourPubKey, theirPubKey) {
		//we are high end
		sendbyte, recvbyte = 0x01, 0x02
	} else {
		//we are low end
		sendbyte, recvbyte = 0x02, 0x01
	}

	s := modExpCT(new(constbn.Int).SetBigInt(theirPubKey), ourPrivKey, pct).GetBigInt()
	secbytes := AppendMPI(nil, s)

	sha := v.hashInstance()

	ret.sendingAESKey = h(sendbyte, secbytes, sha)[:v.keyLength()]
	ret.receivingAESKey = h(recvbyte, secbytes, sha)[:v.keyLength()]

	ret.sendingMACKey = v.hash(ret.sendingAESKey)
	ret.receivingMACKey = v.hash(ret.receivingAESKey)

	ret.extraKey = h(0xFF, secbytes, v.hash2Instance())

	ret.lock()

	return ret
}

func (k *keyManagementContext) pickOurKeys(ourKeyID uint32) (privKey secretKeyValue, pubKey *big.Int, err error) {
	if ourKeyID == 0 || k.ourKeyID == 0 {
		return nil, nil, newOtrConflictError("invalid key id for local peer")
	}

	switch ourKeyID {
	case k.ourKeyID:
		privKey, pubKey = k.ourCurrentDHKeys.priv, k.ourCurrentDHKeys.pub
	case k.ourKeyID - 1:
		privKey, pubKey = k.ourPreviousDHKeys.priv, k.ourPreviousDHKeys.pub
	default:
		err = newOtrConflictError("mismatched key id for local peer")
	}

	return privKey, pubKey, err
}

func (k *keyManagementContext) pickTheirKey(theirKeyID uint32) (pubKey *big.Int, err error) {
	if theirKeyID == 0 || k.theirKeyID == 0 {
		return nil, newOtrConflictError("invalid key id for remote peer")
	}

	switch theirKeyID {
	case k.theirKeyID:
		pubKey = k.theirCurrentDHPubKey
	case k.theirKeyID - 1:
		if k.theirPreviousDHPubKey == nil {
			err = newOtrConflictError("no previous key for remote peer found")
		} else {
			pubKey = k.theirPreviousDHPubKey
		}
	default:
		err = newOtrConflictError("mismatched key id for remote peer")
	}

	return pubKey, err
}

func calculateAKEKeys(s *big.Int, v otrVersion) (ssid [8]byte, revealSigKeys, signatureKeys akeKeys) {
	secbytes := AppendMPI(nil, s)
	sha := v.hash2Instance()
	keys := h(0x01, secbytes, sha)

	copy(ssid[:], h(0x00, secbytes, sha)[:8])
	// SIZE: 16 is the size of the AES key used
	revealSigKeys.c = keys[:16]
	signatureKeys.c = keys[16:]
	revealSigKeys.m1 = h(0x02, secbytes, sha)
	revealSigKeys.m2 = h(0x03, secbytes, sha)
	signatureKeys.m1 = h(0x04, secbytes, sha)
	signatureKeys.m2 = h(0x05, secbytes, sha)

	revealSigKeys.lock()
	signatureKeys.lock()

	return
}

// h1() and h2() are the same
func h(b byte, secbytes []byte, h hash.Hash) []byte {
	h.Reset()
	_, _ = h.Write([]byte{b})
	_, _ = h.Write(secbytes[:])
	return h.Sum(nil)
}
