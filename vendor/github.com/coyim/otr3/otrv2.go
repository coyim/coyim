package otr3

import (
	"bytes"
	"crypto/aes"

	/* #nosec G505*/
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"hash"
	"math/big"
)

var otrv2FragmentationPrefix = []byte("?OTR,")

const otrv2HeaderLen = 3

type otrV2 struct{}

func (v otrV2) parameterLength() int {
	return 16
}

func (v otrV2) isGroupElement(n *big.Int) bool {
	return true
}

func (v otrV2) isFragmented(data []byte) bool {
	return bytes.HasPrefix(data, otrv2FragmentationPrefix)
}

func (v otrV2) parseFragmentPrefix(c *Conversation, data []byte) (rest []byte, ignore bool, ok bool) {
	if len(data) < 5 {
		return data, false, false
	}

	return data[5:], false, true
}

func (v otrV2) fragmentPrefix(n, total int, itags uint32, itagr uint32) []byte {
	return []byte(fmt.Sprintf("%s%05d,%05d,", string(otrv2FragmentationPrefix), n+1, total))
}

func (v otrV2) protocolVersion() uint16 {
	return 2
}

func (v otrV2) whitespaceTag() []byte {
	return convertToWhitespace("2")
}

func (v otrV2) messageHeader(c *Conversation, msgType byte) ([]byte, error) {
	out := AppendShort(nil, v.protocolVersion())
	out = append(out, msgType)
	return out, nil
}

func (v otrV2) parseMessageHeader(c *Conversation, msg []byte) ([]byte, []byte, error) {
	if len(msg) < otrv2HeaderLen {
		return nil, nil, errInvalidOTRMessage
	}
	return msg[:otrv2HeaderLen], msg[otrv2HeaderLen:], nil
}

func (v otrV2) hashInstance() hash.Hash {
	/* #nosec G401*/
	return sha1.New()
	// return sha3.New256()
}

func (v otrV2) hash(val []byte) []byte {
	/* #nosec G401*/
	ret := sha1.Sum(val)
	return ret[:]
	// ret := sha3.Sum256(val)
	// return ret[:]
}

func (v otrV2) hashLength() int {
	/* #nosec G401*/
	return sha1.Size
	// return 32
}

func (v otrV2) hash2Instance() hash.Hash {
	return sha256.New()
	// return sha3.New256()
}

func (v otrV2) hash2(val []byte) []byte {
	ret := sha256.Sum256(val)
	return ret[:]
	// ret := sha3.Sum256(val)
	// return ret[:]
}

func (v otrV2) hash2Length() int {
	return sha256.Size
	// return 32
}

func (v otrV2) truncateLength() int {
	return 20
}

func (v otrV2) keyLength() int {
	return aes.BlockSize
}
