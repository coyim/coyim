package otr3

import (
	"bytes"
	"crypto/aes"

	/* #nosec G505*/
	"crypto/sha1"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"hash"
	"math/big"
	"strconv"
)

var otrv3FragmentationPrefix = []byte("?OTR|")

const (
	otrv3HeaderLen      = 11
	minValidInstanceTag = uint32(0x100)
)

type otrV3 struct{}

func (v otrV3) parameterLength() int {
	return 192
}

func (v otrV3) isGroupElement(n *big.Int) bool {
	return isGroupElement(n)
}

func (v otrV3) isFragmented(data []byte) bool {
	return bytes.HasPrefix(data, otrv3FragmentationPrefix) || otrV2{}.isFragmented(data)
}

func parseItag(s []byte) (uint32, error) {
	v, err := strconv.ParseInt(string(s), 16, 0)
	if err != nil {
		return 0, err
	}
	return uint32(v), nil
}

func (v otrV3) parseFragmentPrefix(c *Conversation, data []byte) (rest []byte, ignore bool, ok bool) {
	if len(data) < 23 {
		return data, false, false
	}

	header := data[:23]
	headerPart := bytes.Split(header, fragmentSeparator)[0]
	itagParts := bytes.Split(headerPart, fragmentItagsSeparator)

	if len(itagParts) < 3 {
		return data, false, false
	}

	senderInstanceTag, err1 := parseItag(itagParts[1])
	if err1 != nil {
		return data, false, false
	}

	receiverInstanceTag, err2 := parseItag(itagParts[2])
	if err2 != nil {
		return data, false, false
	}

	if err := v.verifyInstanceTags(c, senderInstanceTag, receiverInstanceTag); err != nil {
		switch err {
		case errInvalidOTRMessage:
			return data, false, false
		case errReceivedMessageForOtherInstance:
			return data, true, true
		}
	}

	return data[23:], false, true
}

func (v otrV3) fragmentPrefix(n, total int, itags uint32, itagr uint32) []byte {
	return []byte(fmt.Sprintf("%s%08x|%08x,%05d,%05d,", string(otrv3FragmentationPrefix), itags, itagr, n+1, total))
}

func (v otrV3) protocolVersion() uint16 {
	return 3
}

func (v otrV3) whitespaceTag() []byte {
	return convertToWhitespace("3")
}

func (v otrV3) messageHeader(c *Conversation, msgType byte) ([]byte, error) {
	if err := c.generateInstanceTag(); err != nil {
		return nil, err
	}

	out := AppendShort(nil, v.protocolVersion())
	out = append(out, msgType)
	out = AppendWord(out, c.ourInstanceTag)
	out = AppendWord(out, c.theirInstanceTag)
	return out, nil
}

func (c *Conversation) generateInstanceTag() error {
	if c.ourInstanceTag != 0 {
		return nil
	}

	var ret uint32
	var dst [4]byte

	for ret < minValidInstanceTag {
		if err := c.randomInto(dst[:]); err != nil {
			return err
		}

		ret = binary.BigEndian.Uint32(dst[:])
	}

	c.ourInstanceTag = ret

	return nil
}

func malformedMessage(c *Conversation) {
	c.messageEvent(MessageEventReceivedMessageMalformed)
	c.generatePotentialErrorMessage(ErrorCodeMessageMalformed)
}

func (v otrV3) verifyInstanceTags(c *Conversation, their, our uint32) error {
	if c.theirInstanceTag == 0 {
		c.theirInstanceTag = their
	}

	if our > 0 && our < minValidInstanceTag {
		malformedMessage(c)
		return errInvalidOTRMessage
	}

	if their < minValidInstanceTag {
		malformedMessage(c)
		return errInvalidOTRMessage
	}

	if (our != 0 && c.ourInstanceTag != our) ||
		(c.theirInstanceTag != their) {
		c.messageEvent(MessageEventReceivedMessageForOtherInstance)
		return errReceivedMessageForOtherInstance
	}

	return nil
}

func (v otrV3) parseMessageHeader(c *Conversation, msg []byte) ([]byte, []byte, error) {
	if len(msg) < otrv3HeaderLen {
		malformedMessage(c)
		return nil, nil, errInvalidOTRMessage
	}
	header := msg[:otrv3HeaderLen]

	msg, senderInstanceTag, _ := ExtractWord(msg[messageHeaderPrefix:])
	msg, receiverInstanceTag, _ := ExtractWord(msg)

	if err := v.verifyInstanceTags(c, senderInstanceTag, receiverInstanceTag); err != nil {
		return nil, nil, err
	}

	return header, msg, nil
}

func (v otrV3) hashInstance() hash.Hash {
	/* #nosec G401*/
	return sha1.New()
}

func (v otrV3) hash(val []byte) []byte {
	/* #nosec G401*/
	ret := sha1.Sum(val)
	return ret[:]
}

func (v otrV3) hashLength() int {
	return sha1.Size
}

func (v otrV3) hash2Instance() hash.Hash {
	return sha256.New()
}

func (v otrV3) hash2(val []byte) []byte {
	ret := sha256.Sum256(val)
	return ret[:]
}

func (v otrV3) hash2Length() int {
	return sha256.Size
}

func (v otrV3) truncateLength() int {
	return 20
}

func (v otrV3) keyLength() int {
	return aes.BlockSize
}
