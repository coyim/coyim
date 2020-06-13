package otr3

import (
	"bytes"
	"encoding/binary"
)

// GetOurInstanceTag returns our instance tag - it computes it if none has been computed yet
func (c *Conversation) GetOurInstanceTag() uint32 {
	_ = c.generateInstanceTag()
	return c.ourInstanceTag
}

// GetTheirInstanceTag returns the peers instance tag, or 0 if none has been computed yet
func (c *Conversation) GetTheirInstanceTag() uint32 {
	return c.theirInstanceTag
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

// ExtractInstanceTags returns our and theirs instance tags from the message, and ok if the message was parsed properly
func ExtractInstanceTags(m []byte) (ours, theirs uint32, ok bool) {
	if bytes.HasPrefix(m, []byte("?OTR:")) {
		msg, err := decode(encodedMessage(m))
		if err != nil {
			return 0, 0, false
		}

		if len(msg) < otrv3HeaderLen {
			return 0, 0, false
		}

		_, senderInstanceTag, _ := ExtractWord(msg[messageHeaderPrefix:])
		_, receiverInstanceTag, _ := ExtractWord(msg)

		return receiverInstanceTag, senderInstanceTag, true
	} else if bytes.HasPrefix(m, []byte("?OTR|")) {
		if len(m) < 23 {
			return 0, 0, false
		}

		header := m[:23]
		headerPart := bytes.Split(header, fragmentSeparator)[0]
		itagParts := bytes.Split(headerPart, fragmentItagsSeparator)

		if len(itagParts) < 3 {
			return 0, 0, false
		}

		senderInstanceTag, err1 := parseItag(itagParts[1])
		if err1 != nil {
			return 0, 0, false
		}

		receiverInstanceTag, err2 := parseItag(itagParts[2])
		if err2 != nil {
			return 0, 0, false
		}

		return receiverInstanceTag, senderInstanceTag, true
	} else {
		// All other prefixes are for older versions or don't have instance tags
		return 0, 0, false
	}
}
