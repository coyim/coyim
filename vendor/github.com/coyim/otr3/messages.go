package otr3

import (
	"crypto/aes"
	"crypto/hmac"
	"crypto/subtle"
	"encoding/binary"
	"math/big"
)

const (
	messageFlagNormal           = byte(0x00)
	messageFlagIgnoreUnreadable = byte(0x01)

	messageHeaderPrefix = 3

	msgTypeDHCommit  = byte(0x02)
	msgTypeData      = byte(0x03)
	msgTypeDHKey     = byte(0x0A)
	msgTypeRevealSig = byte(0x11)
	msgTypeSig       = byte(0x12)
)

type message interface {
	serialize() []byte
	deserialize(msg []byte) error
}

type dhCommit struct {
	encryptedGx []byte

	// SIZE: this should always be version.hash2Length
	yhashedGx []byte
}

func (c dhCommit) serialize() []byte {
	out := AppendData(nil, c.encryptedGx)
	out = AppendData(out, c.yhashedGx)
	return out
}

func (c *dhCommit) deserialize(msg []byte) error {
	msg, g, ok := ExtractData(msg)
	_, h, ok2 := ExtractData(msg)

	if !(ok && ok2) {
		return newOtrError("corrupt DH commit message")
	}

	c.encryptedGx = g
	c.yhashedGx = h
	return nil
}

type dhKey struct {
	gy *big.Int
}

func (c dhKey) serialize() []byte {
	return AppendMPI(nil, c.gy)
}

func (c *dhKey) deserialize(msg []byte) error {
	_, gy, ok := ExtractMPI(msg)

	if !ok {
		return newOtrError("corrupt DH key message")
	}

	c.gy = gy
	return nil
}

type revealSig struct {
	// TODO: why this number here?
	r            [16]byte
	encryptedSig []byte
	macSig       []byte
}

func (c revealSig) serialize(v otrVersion) []byte {
	var out []byte
	out = AppendData(out, c.r[:])
	out = append(out, c.encryptedSig...)
	return append(out, c.macSig[:v.truncateLength()]...)
}

func (c *revealSig) deserialize(msg []byte, v otrVersion) error {
	in, r, ok := ExtractData(msg)
	okLen := len(r) == 16
	macSig, encryptedSig, ok2 := ExtractData(in)
	okLen2 := len(macSig) == v.truncateLength()

	if !(ok && ok2 && okLen && okLen2) {
		return newOtrError("corrupt reveal signature message")
	}

	copy(c.r[:], r)
	c.encryptedSig = encryptedSig
	c.macSig = macSig
	return nil
}

type sig struct {
	encryptedSig []byte
	macSig       []byte
}

func (c sig) serialize(v otrVersion) []byte {
	var out []byte
	out = append(out, c.encryptedSig...)
	return append(out, c.macSig[:v.truncateLength()]...)
}

func (c *sig) deserialize(msg []byte) error {
	macSig, encryptedSig, ok := ExtractData(msg)

	if !ok || len(macSig) != 20 {
		return newOtrError("corrupt signature message")
	}
	c.encryptedSig = encryptedSig
	c.macSig = macSig
	return nil
}

type dataMsg struct {
	flag                        byte
	senderKeyID, recipientKeyID uint32
	y                           *big.Int
	// SIZE: 8 is half of the AES block size used
	topHalfCtr             [8]byte
	encryptedMsg           []byte
	authenticator          []byte
	oldMACKeys             []macKey
	serializeUnsignedCache []byte
}

func (c *dataMsg) sign(key []byte, header []byte, v otrVersion) {
	if c.serializeUnsignedCache == nil {
		c.serializeUnsignedCache = c.serializeUnsigned()
	}
	mac := hmac.New(v.hashInstance, key)
	mac.Write(header)
	mac.Write(c.serializeUnsignedCache)
	c.authenticator = mac.Sum(nil)
}

func (c dataMsg) checkSign(key []byte, header []byte, v otrVersion) error {
	mac := hmac.New(v.hashInstance, key[:])
	mac.Write(header)
	mac.Write(c.serializeUnsignedCache)
	authenticatorCalculated := mac.Sum(nil)

	if subtle.ConstantTimeCompare(c.authenticator, authenticatorCalculated) == 0 {
		return newOtrConflictError("bad signature MAC in encrypted signature")
	}
	return nil
}

func (c dataMsg) serializeUnsigned() []byte {
	var out []byte

	out = append(out, c.flag)
	out = AppendWord(out, c.senderKeyID)
	out = AppendWord(out, c.recipientKeyID)
	out = AppendMPI(out, c.y)
	out = append(out, c.topHalfCtr[:]...)
	out = AppendData(out, c.encryptedMsg)
	return out
}

func (c *dataMsg) deserializeUnsigned(msg []byte) error {
	if len(msg) == 0 {
		return newOtrError("dataMsg.deserialize empty message")
	}
	in := msg
	c.flag = in[0]

	in = in[1:]
	var ok bool

	in, c.senderKeyID, ok = ExtractWord(in)
	if !ok {
		return newOtrError("dataMsg.deserialize corrupted senderKeyID")
	}

	in, c.recipientKeyID, ok = ExtractWord(in)
	if !ok {
		return newOtrError("dataMsg.deserialize corrupted recipientKeyID")
	}

	in, c.y, ok = ExtractMPI(in)
	if !ok {
		return newOtrError("dataMsg.deserialize corrupted y")
	}

	if len(in) < len(c.topHalfCtr) {
		return newOtrError("dataMsg.deserialize corrupted topHalfCtr")
	}

	copy(c.topHalfCtr[:], in)
	if binary.BigEndian.Uint64(c.topHalfCtr[:]) == 0 {
		return newOtrError("dataMsg.deserialize invalid topHalfCtr")
	}

	copy(c.topHalfCtr[:], in)
	in = in[len(c.topHalfCtr):]
	in, c.encryptedMsg, ok = ExtractData(in)
	if !ok {
		return newOtrError("dataMsg.deserialize corrupted encryptedMsg")
	}

	c.serializeUnsignedCache = msg[:len(msg)-len(in)]
	return nil
}

func (c dataMsg) serialize(v otrVersion) []byte {
	if c.serializeUnsignedCache == nil {
		c.serializeUnsignedCache = c.serializeUnsigned()
	}

	out := makeCopy(c.serializeUnsignedCache)
	out = append(out, c.authenticator...)

	keyLen := v.hashLength()
	revKeys := make([]byte, 0, len(c.oldMACKeys)*keyLen)
	for _, k := range c.oldMACKeys {
		revKeys = append(revKeys, k[:]...)
	}
	out = AppendData(out, revKeys)

	return out
}

func (c *dataMsg) deserialize(msg []byte, v otrVersion) error {
	if err := c.deserializeUnsigned(msg); err != nil {
		return err
	}

	msg = msg[len(c.serializeUnsignedCache):]
	c.authenticator = msg[0:v.hashLength()]
	msg = msg[len(c.authenticator):]

	var revKeysBytes []byte
	msg, revKeysBytes, ok := ExtractData(msg)
	if !ok {
		return newOtrError("dataMsg.deserialize corrupted revealMACKeys")
	}
	for len(revKeysBytes) > 0 {
		if len(revKeysBytes) < v.hashLength() {
			return newOtrError("dataMsg.deserialize corrupted revealMACKeys")
		}
		revKey := make([]byte, v.hashLength())
		copy(revKey, revKeysBytes)
		c.oldMACKeys = append(c.oldMACKeys, revKey)
		revKeysBytes = revKeysBytes[len(revKey):]
	}

	return nil
}

type plainDataMsg struct {
	message []byte
	tlvs    []tlv
}

func (c *plainDataMsg) deserialize(msg []byte) error {
	nulPos := 0
	for nulPos < len(msg) && msg[nulPos] != 0x00 {
		nulPos++
	}

	var tlvsBytes []byte
	if nulPos < len(msg) {
		c.message = msg[:nulPos]
		tlvsBytes = msg[nulPos+1:]
	} else {
		c.message = msg
	}

	for len(tlvsBytes) > 0 {
		atlv := tlv{}
		if err := atlv.deserialize(tlvsBytes); err != nil {
			return err
		}
		c.tlvs = append(c.tlvs, atlv)
		tlvsBytes = tlvsBytes[4+int(atlv.tlvLength):]
	}
	return nil
}

func (c plainDataMsg) serialize() []byte {
	out := c.message
	out = append(out, 0x00)

	if len(c.tlvs) > 0 {
		for i := range c.tlvs {
			out = AppendShort(out, c.tlvs[i].tlvType)
			out = AppendShort(out, c.tlvs[i].tlvLength)
			out = append(out, c.tlvs[i].tlvValue...)
		}
	}
	return out
}

const (
	paddingGranularity = 256
	tlvHeaderLen       = 4
	nulByteLen         = 1
)

func (c plainDataMsg) pad() plainDataMsg {
	padding := paddingGranularity - ((len(c.message) + tlvHeaderLen + nulByteLen) % paddingGranularity)

	paddingTlv := tlv{
		tlvType:   uint16(tlvTypePadding),
		tlvLength: uint16(padding),
		tlvValue:  make([]byte, padding),
	}

	c.tlvs = append(c.tlvs, paddingTlv)

	return c
}

func (c plainDataMsg) encrypt(key []byte, topHalfCtr [8]byte) []byte {
	var iv [aes.BlockSize]byte
	copy(iv[:], topHalfCtr[:])

	data := c.pad().serialize()
	dst := make([]byte, len(data))
	counterEncipher(key, iv[:], data, dst)

	wipeBytes(iv[:])
	return dst
}

func (c *plainDataMsg) decrypt(key []byte, topHalfCtr [8]byte, src []byte) error {
	var iv [aes.BlockSize]byte
	copy(iv[:], topHalfCtr[:])

	if err := counterEncipher(key, iv[:], src, src); err != nil {
		return err
	}

	wipeBytes(iv[:])

	c.deserialize(src)
	return nil
}
