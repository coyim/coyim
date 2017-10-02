package otr3

import (
	"crypto/aes"
	"math/big"
	"testing"
)

func Test_tlvSerialize(t *testing.T) {
	expectedTLVBytes := []byte{0x00, 0x01, 0x00, 0x02, 0x01, 0x01}
	aTLV := tlv{
		tlvType:   0x0001,
		tlvLength: 0x0002,
		tlvValue:  []byte{0x01, 0x01},
	}
	aTLVBytes := aTLV.serialize()
	assertDeepEquals(t, aTLVBytes, expectedTLVBytes)
}

func Test_tlvDeserialize(t *testing.T) {
	aTLVBytes := []byte{0x00, 0x01, 0x00, 0x02, 0x01, 0x01}
	aTLV := tlv{}
	expectedTLV := tlv{
		tlvType:   0x0001,
		tlvLength: 0x0002,
		tlvValue:  []byte{0x01, 0x01},
	}
	err := aTLV.deserialize(aTLVBytes)
	assertEquals(t, err, nil)
	assertDeepEquals(t, aTLV, expectedTLV)
}

func Test_tlvDeserializeWithWrongType(t *testing.T) {
	aTLVBytes := []byte{0x00}
	aTLV := tlv{}
	err := aTLV.deserialize(aTLVBytes)
	assertEquals(t, err.Error(), "otr: wrong tlv type")
}

func Test_tlvDeserializeWithWrongLength(t *testing.T) {
	aTLVBytes := []byte{0x00, 0x01, 0x00}
	aTLV := tlv{}
	err := aTLV.deserialize(aTLVBytes)
	assertEquals(t, err.Error(), "otr: wrong tlv length")
}

func Test_tlvDeserializeWithWrongValue(t *testing.T) {
	aTLVBytes := []byte{0x00, 0x01, 0x00, 0x02, 0x01}
	aTLV := tlv{}
	err := aTLV.deserialize(aTLVBytes)
	assertEquals(t, err.Error(), "otr: wrong tlv value")
}

func Test_dataMsgSignWithSerializeUnsignedCache(t *testing.T) {
	m := dataMsg{
		serializeUnsignedCache: []byte{0x0, 0x0, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x1, 0x1, 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x0, 0x0, 0x0, 0x4, 0x0, 0x1, 0x2, 0x3},
	}
	macKey := macKey{0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03}
	m.sign(macKey, []byte{}, otrV3{})
	assertDeepEquals(t, m.authenticator, []byte{0x6f, 0x6, 0x76, 0x45, 0xbb, 0x94, 0x5c, 0xa2, 0xfc, 0x13, 0xa9, 0xfa, 0x58, 0xb7, 0xd7, 0x23, 0xee, 0xab, 0x62, 0xe8})
}

func Test_dataMsgSignWithoutSerializeUnsignedCache(t *testing.T) {
	m := dataMsg{
		flag:           byte(0x00),
		senderKeyID:    uint32(0x00000000),
		recipientKeyID: uint32(0x00000001),
		y:              big.NewInt(1),
		topHalfCtr:     [8]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
		encryptedMsg:   []byte{0x00, 0x01, 0x02, 0x03},
	}
	macKey := macKey{0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03}
	m.sign(macKey[:], []byte{}, otrV3{})
	assertDeepEquals(t, m.serializeUnsignedCache, []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x1, 0x1, 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x0, 0x0, 0x0, 0x4, 0x0, 0x1, 0x2, 0x3})
	assertDeepEquals(t, m.authenticator, []byte{0x36, 0x84, 0x4d, 0x5d, 0x5, 0xe7, 0x88, 0xd3, 0x46, 0x90, 0x27, 0xe4, 0x2, 0x7e, 0x2b, 0x6b, 0x10, 0x7b, 0xdd, 0x79})
}

func Test_dataMsg_deserializeUnsigned_failsWhenTopHalfCtrIsZero(t *testing.T) {
	msg := dataMsg{
		flag:           byte(0x00),
		senderKeyID:    uint32(0x00000000),
		recipientKeyID: uint32(0x00000001),
		y:              big.NewInt(1),
		topHalfCtr:     [8]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	}.serializeUnsigned()

	dataMessage := dataMsg{}
	err := dataMessage.deserializeUnsigned(msg)

	assertEquals(t, err.Error(), "otr: dataMsg.deserialize invalid topHalfCtr")
}

func Test_dataMsgCheckSignWithoutError(t *testing.T) {
	m := dataMsg{
		serializeUnsignedCache: []byte{0x0, 0x0, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x1, 0x1, 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x0, 0x0, 0x0, 0x4, 0x0, 0x1, 0x2, 0x3},
		authenticator:          []byte{0x6f, 0x6, 0x76, 0x45, 0xbb, 0x94, 0x5c, 0xa2, 0xfc, 0x13, 0xa9, 0xfa, 0x58, 0xb7, 0xd7, 0x23, 0xee, 0xab, 0x62, 0xe8},
	}
	macKey := macKey{0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03}
	assertDeepEquals(t, m.checkSign(macKey, []byte{}, otrV3{}), nil)
}

func Test_dataMsgCheckSignWithError(t *testing.T) {
	m := dataMsg{
		serializeUnsignedCache: []byte{0x0, 0x0, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x1, 0x1, 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x0, 0x0, 0x0, 0x4, 0x0, 0x1, 0x2, 0x3},
		authenticator:          []byte{0x6e, 0x6, 0x76, 0x45, 0xbb, 0x94, 0x5c, 0xa2, 0xfc, 0x13, 0xa9, 0xfa, 0x58, 0xb7, 0xd7, 0x23, 0xee, 0xab, 0x62, 0xe8},
	}
	macKey := macKey{0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03}
	assertDeepEquals(t, m.checkSign(macKey, []byte{}, otrV3{}), newOtrConflictError("bad signature MAC in encrypted signature"))
}

func Test_dataMsgDeserialze(t *testing.T) {
	var msg []byte

	flag := byte(0x00)
	senderKeyID := uint32(0x00000000)
	recipientKeyID := uint32(0x00000001)
	y := big.NewInt(1)
	topHalfCtr := [8]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}
	encryptedMsg := []byte{0x00, 0x01, 0x02, 0x03}

	authenticator := []byte{0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03}
	oldMACKeys := []macKey{
		macKey{0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03},
		macKey{0x01, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03},
	}

	msg = append(msg, flag)

	msg = appendWord(msg, senderKeyID)
	msg = appendWord(msg, recipientKeyID)

	msg = appendMPI(msg, y)

	msg = append(msg, topHalfCtr[:]...)

	msg = appendData(msg, encryptedMsg)

	msg = append(msg, authenticator[:]...)
	revKeys := make([]byte, 0, len(oldMACKeys)*otrV3{}.hashLength())
	for _, k := range oldMACKeys {
		revKeys = append(revKeys, k[:]...)
	}
	msg = appendData(msg, revKeys)

	dataMessage := dataMsg{}
	err := dataMessage.deserialize(msg, otrV3{})
	assertEquals(t, err, nil)
	assertDeepEquals(t, dataMessage.flag, flag)
	assertDeepEquals(t, dataMessage.senderKeyID, senderKeyID)
	assertDeepEquals(t, dataMessage.recipientKeyID, recipientKeyID)
	assertDeepEquals(t, dataMessage.y, y)
	assertDeepEquals(t, dataMessage.topHalfCtr, topHalfCtr)
	assertDeepEquals(t, dataMessage.encryptedMsg, encryptedMsg)
	assertDeepEquals(t, dataMessage.authenticator, authenticator)
	assertDeepEquals(t, dataMessage.oldMACKeys, oldMACKeys)
}

func Test_dataMsgDeserialzeErrorWhenEmpty(t *testing.T) {
	var msg []byte

	dataMessage := dataMsg{}
	err := dataMessage.deserialize(msg, otrV3{})
	assertEquals(t, err.Error(), "otr: dataMsg.deserialize empty message")
}

func Test_dataMsgDeserialzeErrorWhenCorruptedSenderKeyID(t *testing.T) {
	var msg []byte

	flag := byte(0x00)
	senderKeyID := byte(0x00)
	msg = append(msg, flag)

	msg = append(msg, senderKeyID)

	dataMessage := dataMsg{}
	err := dataMessage.deserialize(msg, otrV3{})
	assertEquals(t, err.Error(), "otr: dataMsg.deserialize corrupted senderKeyID")
}

func Test_dataMsgDeserialzeErrorWhenCorruptedReceiverKeyID(t *testing.T) {
	var msg []byte

	flag := byte(0x00)
	senderKeyID := uint32(0x00000000)
	recipientKeyID := byte(0x00)
	msg = append(msg, flag)

	msg = appendWord(msg, senderKeyID)
	msg = append(msg, recipientKeyID)

	dataMessage := dataMsg{}
	err := dataMessage.deserialize(msg, otrV3{})
	assertEquals(t, err.Error(), "otr: dataMsg.deserialize corrupted recipientKeyID")
}

func Test_dataMsgDeserialzeErrorWhenCorruptedY(t *testing.T) {
	var msg []byte

	flag := byte(0x00)
	senderKeyID := uint32(0x00000000)
	recipientKeyID := byte(0x00)
	y := big.NewInt(1)
	msg = append(msg, flag)

	msg = appendWord(msg, senderKeyID)
	msg = append(msg, recipientKeyID)
	mpiY := appendMPI([]byte{}, y)
	msg = append(msg, mpiY[1:]...)

	dataMessage := dataMsg{}
	err := dataMessage.deserialize(msg, otrV3{})
	assertEquals(t, err.Error(), "otr: dataMsg.deserialize corrupted y")
}

func Test_dataMsgDeserialzeErrorWhenCorruptedEncryptedMsg(t *testing.T) {
	var msg []byte

	flag := byte(0x00)
	senderKeyID := uint32(0x00000000)
	recipientKeyID := uint32(0x00000001)
	y := big.NewInt(1)
	topHalfCtr := [8]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}
	encryptedMsg := []byte{0x00, 0x01, 0x02, 0x03}

	msg = append(msg, flag)
	msg = appendWord(msg, senderKeyID)
	msg = appendWord(msg, recipientKeyID)
	msg = appendMPI(msg, y)
	msg = append(msg, topHalfCtr[:]...)
	encryptedMsgData := appendData([]byte{}, encryptedMsg)
	msg = append(msg, encryptedMsgData[1:]...)

	dataMessage := dataMsg{}
	err := dataMessage.deserialize(msg, otrV3{})
	assertEquals(t, err.Error(), "otr: dataMsg.deserialize corrupted encryptedMsg")
}

func Test_dataMsgDeserialzeErrorWhenCorruptedTopHalfCtr(t *testing.T) {
	var msg []byte

	flag := byte(0x00)
	senderKeyID := uint32(0x00000000)
	recipientKeyID := uint32(0x00000001)
	y := big.NewInt(1)
	topHalfCtr := [7]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06}

	msg = append(msg, flag)
	msg = appendWord(msg, senderKeyID)
	msg = appendWord(msg, recipientKeyID)
	msg = appendMPI(msg, y)

	msg = append(msg, topHalfCtr[:]...)

	dataMessage := dataMsg{}
	err := dataMessage.deserialize(msg, otrV3{})
	assertEquals(t, err.Error(), "otr: dataMsg.deserialize corrupted topHalfCtr")
}

func Test_dataMsgDeserialzeErrorWhenCorruptedRevealMACKeys(t *testing.T) {
	var msg []byte

	flag := byte(0x00)
	senderKeyID := uint32(0x00000000)
	recipientKeyID := uint32(0x00000001)
	y := big.NewInt(1)
	topHalfCtr := [8]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}
	encryptedMsg := []byte{0x00, 0x01, 0x02, 0x03}

	authenticator := []byte{0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03}
	oldMACKeys := []macKey{
		macKey{0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03},
		macKey{0x01, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03},
	}

	msg = append(msg, flag)

	msg = appendWord(msg, senderKeyID)
	msg = appendWord(msg, recipientKeyID)

	msg = appendMPI(msg, y)

	msg = append(msg, topHalfCtr[:]...)

	msg = appendData(msg, encryptedMsg)

	msg = append(msg, authenticator...)
	revKeys := make([]byte, 0, len(oldMACKeys)*otrV3{}.hashLength())
	for _, k := range oldMACKeys {
		revKeys = append(revKeys, k[1:]...)
	}
	msg = appendData(msg, revKeys)

	dataMessage := dataMsg{}
	err := dataMessage.deserialize(msg, otrV3{})
	assertEquals(t, err.Error(), "otr: dataMsg.deserialize corrupted revealMACKeys")
}

func Test_dataMsgDeserialzeErrorWhenCorruptedRevealMACKeyEnding(t *testing.T) {
	var msg []byte

	flag := byte(0x00)
	senderKeyID := uint32(0x00000000)
	recipientKeyID := uint32(0x00000001)
	y := big.NewInt(1)
	topHalfCtr := [8]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}
	encryptedMsg := []byte{0x00, 0x01, 0x02, 0x03}

	authenticator := []byte{0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03}
	oldMACKeys := []macKey{
		macKey{0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03},
		macKey{0x01, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03},
	}

	msg = append(msg, flag)

	msg = appendWord(msg, senderKeyID)
	msg = appendWord(msg, recipientKeyID)

	msg = appendMPI(msg, y)

	msg = append(msg, topHalfCtr[:]...)

	msg = appendData(msg, encryptedMsg)

	msg = append(msg, authenticator[:]...)
	revKeys := make([]byte, 0, len(oldMACKeys)*otrV3{}.hashLength())
	revKeys = append(revKeys, oldMACKeys[0][:]...)
	revKeys = append(revKeys, oldMACKeys[1][:otrV3{}.hashLength()-1]...)
	msg = appendData(msg, revKeys)

	dataMessage := dataMsg{}
	err := dataMessage.deserialize(msg, otrV3{})
	assertEquals(t, err.Error(), "otr: dataMsg.deserialize corrupted revealMACKeys")
}

func Test_plainDataMsgShouldDeserializeOneTLV(t *testing.T) {
	plain := []byte("helloworld")
	atlvBytes := []byte{0x00, 0x01, 0x00, 0x02, 0x01, 0x01}
	msg := append(plain, 0x00)
	msg = append(msg, atlvBytes...)
	aDataMsg := plainDataMsg{}
	err := aDataMsg.deserialize(msg)
	atlv := tlv{
		tlvType:   0x0001,
		tlvLength: 0x0002,
		tlvValue:  []byte{0x01, 0x01},
	}

	assertEquals(t, err, nil)
	assertDeepEquals(t, aDataMsg.message, plain)
	assertDeepEquals(t, aDataMsg.tlvs[0], atlv)
}

func Test_plainDataMsgShouldDeserializeMultiTLV(t *testing.T) {
	plain := []byte("helloworld")
	atlvBytes := []byte{0x00, 0x01, 0x00, 0x02, 0x01, 0x01}
	btlvBytes := []byte{0x00, 0x02, 0x00, 0x05, 0x01, 0x01, 0x01, 0x01, 0x01}
	msg := append(plain, 0x00)
	msg = append(msg, atlvBytes...)
	msg = append(msg, btlvBytes...)
	aDataMsg := plainDataMsg{}
	err := aDataMsg.deserialize(msg)
	atlv := tlv{
		tlvType:   0x0001,
		tlvLength: 0x0002,
		tlvValue:  []byte{0x01, 0x01},
	}

	btlv := tlv{
		tlvType:   0x0002,
		tlvLength: 0x0005,
		tlvValue:  []byte{0x01, 0x01, 0x01, 0x01, 0x01},
	}

	assertEquals(t, err, nil)
	assertDeepEquals(t, aDataMsg.message, plain)
	assertDeepEquals(t, aDataMsg.tlvs[0], atlv)
	assertDeepEquals(t, aDataMsg.tlvs[1], btlv)
}

func Test_plainDataMsgShouldDeserializeNoTLV(t *testing.T) {
	plain := []byte("helloworld")
	aDataMsg := plainDataMsg{}
	err := aDataMsg.deserialize(plain)
	assertEquals(t, err, nil)
	assertDeepEquals(t, aDataMsg.message, plain)
	assertDeepEquals(t, len(aDataMsg.tlvs), 0)
}

func Test_plainDataMsgShouldSerialize(t *testing.T) {
	plain := []byte("helloworld")
	atlvBytes := []byte{0x00, 0x01, 0x00, 0x02, 0x01, 0x01}
	btlvBytes := []byte{0x00, 0x02, 0x00, 0x05, 0x01, 0x01, 0x01, 0x01, 0x01}
	msg := append(plain, 0x00)
	msg = append(msg, atlvBytes...)
	msg = append(msg, btlvBytes...)
	aDataMsg := plainDataMsg{}
	atlv := tlv{
		tlvType:   0x0001,
		tlvLength: 0x0002,
		tlvValue:  []byte{0x01, 0x01},
	}

	btlv := tlv{
		tlvType:   0x0002,
		tlvLength: 0x0005,
		tlvValue:  []byte{0x01, 0x01, 0x01, 0x01, 0x01},
	}
	aDataMsg.message = plain
	aDataMsg.tlvs = []tlv{atlv, btlv}

	assertDeepEquals(t, aDataMsg.serialize(), msg)
}

func Test_plainDataMsgShouldSerializeWithoutTLVs(t *testing.T) {
	plain := []byte("helloworld")
	expected := append(plain, 0x00)

	dataMsg := plainDataMsg{
		message: plain,
	}

	assertDeepEquals(t, dataMsg.serialize(), expected)
}

func Test_encrypt_EncryptsPlainMessageUsingSendingAESKeyAndCounter(t *testing.T) {
	plain := plainDataMsg{
		message: []byte("we are awesome"),
	}

	var sendingAESKey [aes.BlockSize]byte
	topHalfCtr := [8]byte{}
	copy(sendingAESKey[:], bytesFromHex("42e258bebf031acf442f52d6ef52d6f1"))
	expectedEncrypted := bytesFromHex("4f0de18011633ed0264ccc1840d64f4cf8f0c91ef78890ab82edef36cb38210bb80760585ff43d736a9ff3e4bb05fc088fa34c2f21012988d539ebc839e9bc97633f4c42de15ea5c3c55a2b9940ca35015ded14205b9df78f936cb1521aedbea98df7dc03c116570ba8d034abc8e2d23185d2ce225845f38c08cb2aae192d66d601c1bc86149c98e8874705ae365b31cda76d274429de5e07b93f0ff29152716980a63c31b7bda150b222ba1d373f786d5f59f580d4f690a71d7fc620e0a3b05d692221ddeebac98d6ed16272e7c4596de27fb104ad747aa9a3ad9d3bc4f988af0beb21760df06047e267af0109baceb0f363bcaff7b205f2c42b3cb67a942f2")

	encrypted := plain.encrypt(sendingAESKey[:], topHalfCtr)

	assertDeepEquals(t, encrypted, expectedEncrypted)
}

func Test_encrypt_EncryptsPlainMessageUsingSendingAESKeyAndCounterNotZero(t *testing.T) {
	plain := plainDataMsg{
		message: []byte("we are awesome"),
	}

	var sendingAESKey [aes.BlockSize]byte
	topHalfCtr := [8]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}
	copy(sendingAESKey[:], bytesFromHex("42e258bebf031acf442f52d6ef52d6f1"))
	expectedEncrypted := bytesFromHex("2dccced4937a337e01bc2ed969b4f60d3ab0a4844aef0a02ebc5c6f09f71a7819687cdbcf2a912be1e8ceda086d188ce3e0bbfecaa77a050a5ed9f98f0c6590579e4d1fb9f753102955dcfc5535af3906ff7d62490362e6e89e28c3b41081f2ce3e8c2ea154a582ff7a1449e7ad8abf295b5e3f8fb80e9b6482fc3bae869ccdb9144f0242604ddee924f388c308c6ce123b5ae22a93ac7c315b13019d474134dd9fd15334fade1b6737b11f79a3cfeed8dd18d72739436ebb560ecdca71a9a67c7b97c2526119a4b1323a6de7c70dffaf7229d798aaea4a692410a139249305d3059685b6ecd0760323ea16db9e02497f5657d1a5d82e09df0088e572b5d0bd7")

	encrypted := plain.encrypt(sendingAESKey[:], topHalfCtr)

	assertDeepEquals(t, encrypted, expectedEncrypted)
}

func Test_pad_PlainMessageUsingTLV0(t *testing.T) {
	plain := plainDataMsg{
		message: []byte("123456"),
		tlvs: []tlv{
			smpMessageAbort{}.tlv(),
		},
	}

	paddedMessage := plain.pad()

	assertEquals(t, len(paddedMessage.tlvs), 2)
	assertEquals(t, paddedMessage.tlvs[1].tlvLength, uint16(245))
}

func Test_dataMsg_serializeExposesOldMACKeys(t *testing.T) {
	macKey1 := bytesFromHex("a45e2b122f58bbe2042f73f092329ad9b5dfe23e")
	macKey2 := bytesFromHex("e55a2b111f60bbe1041f73f003333ad9a5dfe22a")

	keyLen := otrV3{}.hashLength()

	m := dataMsg{
		y:          big.NewInt(0x01),
		oldMACKeys: []macKey{macKey1, macKey2},
	}
	msg := m.serialize(otrV3{})
	MACsIndex := len(msg) - 2*keyLen - 4

	_, expectedData, _ := extractData(msg[MACsIndex:])
	assertDeepEquals(t, expectedData[:keyLen], macKey1[:])
	assertDeepEquals(t, expectedData[keyLen:], macKey2[:])
}
