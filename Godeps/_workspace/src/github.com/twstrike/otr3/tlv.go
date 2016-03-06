package otr3

import "bytes"

const tlvHeaderLength = 4

const (
	tlvTypePadding           = uint16(0x00)
	tlvTypeDisconnected      = uint16(0x01)
	tlvTypeSMP1              = uint16(0x02)
	tlvTypeSMP2              = uint16(0x03)
	tlvTypeSMP3              = uint16(0x04)
	tlvTypeSMP4              = uint16(0x05)
	tlvTypeSMPAbort          = uint16(0x06)
	tlvTypeSMP1WithQuestion  = uint16(0x07)
	tlvTypeExtraSymmetricKey = uint16(0x08)
)

type tlvHandler func(*Conversation, tlv, dataMessageExtra) (*tlv, error)

var tlvHandlers = make([]tlvHandler, 9)

func initTLVHandlers() {
	tlvHandlers[tlvTypePadding] = func(c *Conversation, t tlv, x dataMessageExtra) (*tlv, error) {
		return c.processPaddingTLV(t, x)
	}
	tlvHandlers[tlvTypeDisconnected] = func(c *Conversation, t tlv, x dataMessageExtra) (*tlv, error) {
		return c.processDisconnectedTLV(t, x)
	}
	tlvHandlers[tlvTypeSMP1] = func(c *Conversation, t tlv, x dataMessageExtra) (*tlv, error) {
		return c.processSMPTLV(t, x)
	}
	tlvHandlers[tlvTypeSMP2] = func(c *Conversation, t tlv, x dataMessageExtra) (*tlv, error) {
		return c.processSMPTLV(t, x)
	}
	tlvHandlers[tlvTypeSMP3] = func(c *Conversation, t tlv, x dataMessageExtra) (*tlv, error) {
		return c.processSMPTLV(t, x)
	}
	tlvHandlers[tlvTypeSMP4] = func(c *Conversation, t tlv, x dataMessageExtra) (*tlv, error) {
		return c.processSMPTLV(t, x)
	}
	tlvHandlers[tlvTypeSMPAbort] = func(c *Conversation, t tlv, x dataMessageExtra) (*tlv, error) {
		return c.processSMPTLV(t, x)
	}
	tlvHandlers[tlvTypeSMP1WithQuestion] = func(c *Conversation, t tlv, x dataMessageExtra) (*tlv, error) {
		return c.processSMPTLV(t, x)
	}
	tlvHandlers[tlvTypeExtraSymmetricKey] = func(c *Conversation, t tlv, x dataMessageExtra) (*tlv, error) {
		return c.processExtraSymmetricKeyTLV(t, x)
	}
}

func messageHandlerForTLV(t tlv) (tlvHandler, error) {
	if t.tlvType >= uint16(len(tlvHandlers)) {
		return nil, newOtrError("unexpected TLV type")
	}
	return tlvHandlers[t.tlvType], nil
}

type tlv struct {
	tlvType   uint16
	tlvLength uint16
	tlvValue  []byte
}

func (c tlv) serialize() []byte {
	out := appendShort([]byte{}, c.tlvType)
	out = appendShort(out, c.tlvLength)
	return append(out, c.tlvValue...)
}

func (c *tlv) deserialize(tlvsBytes []byte) error {
	var ok bool
	tlvsBytes, c.tlvType, ok = extractShort(tlvsBytes)
	if !ok {
		return newOtrError("wrong tlv type")
	}
	tlvsBytes, c.tlvLength, ok = extractShort(tlvsBytes)
	if !ok {
		return newOtrError("wrong tlv length")
	}
	if len(tlvsBytes) < int(c.tlvLength) {
		return newOtrError("wrong tlv value")
	}
	c.tlvValue = tlvsBytes[:int(c.tlvLength)]
	return nil
}

func (c tlv) isSMPMessage() bool {
	return c.tlvType >= tlvTypeSMP1 && c.tlvType <= tlvTypeSMP1WithQuestion
}

func (c tlv) smpMessage() (smpMessage, bool) {
	switch c.tlvType {
	case tlvTypeSMP1:
		return toSmpMessage1(c)
	case tlvTypeSMP1WithQuestion:
		return toSmpMessage1Q(c)
	case tlvTypeSMP2:
		return toSmpMessage2(c)
	case tlvTypeSMP3:
		return toSmpMessage3(c)
	case tlvTypeSMP4:
		return toSmpMessage4(c)
	case tlvTypeSMPAbort:
		return toSmpMessageAbort(c)
	}

	return nil, false
}

func toSmpMessage1(t tlv) (msg smp1Message, ok bool) {
	_, mpis, ok := extractMPIs(t.tlvValue)
	if !ok || len(mpis) < 6 {
		return msg, false
	}
	msg.g2a = mpis[0]
	msg.c2 = mpis[1]
	msg.d2 = mpis[2]
	msg.g3a = mpis[3]
	msg.c3 = mpis[4]
	msg.d3 = mpis[5]
	return msg, true
}

func toSmpMessage1Q(t tlv) (msg smp1Message, ok bool) {
	nulPos := bytes.IndexByte(t.tlvValue, 0)
	if nulPos == -1 {
		return msg, false
	}
	question := string(t.tlvValue[:nulPos])
	t.tlvValue = t.tlvValue[(nulPos + 1):]
	msg, ok = toSmpMessage1(t)
	msg.hasQuestion = true
	msg.question = question
	return msg, ok
}

func toSmpMessage2(t tlv) (msg smp2Message, ok bool) {
	_, mpis, ok := extractMPIs(t.tlvValue)
	if !ok || len(mpis) < 11 {
		return msg, false
	}
	msg.g2b = mpis[0]
	msg.c2 = mpis[1]
	msg.d2 = mpis[2]
	msg.g3b = mpis[3]
	msg.c3 = mpis[4]
	msg.d3 = mpis[5]
	msg.pb = mpis[6]
	msg.qb = mpis[7]
	msg.cp = mpis[8]
	msg.d5 = mpis[9]
	msg.d6 = mpis[10]
	return msg, true
}

func toSmpMessage3(t tlv) (msg smp3Message, ok bool) {
	_, mpis, ok := extractMPIs(t.tlvValue)
	if !ok || len(mpis) < 8 {
		return msg, false
	}
	msg.pa = mpis[0]
	msg.qa = mpis[1]
	msg.cp = mpis[2]
	msg.d5 = mpis[3]
	msg.d6 = mpis[4]
	msg.ra = mpis[5]
	msg.cr = mpis[6]
	msg.d7 = mpis[7]
	return msg, true
}

func toSmpMessage4(t tlv) (msg smp4Message, ok bool) {
	_, mpis, ok := extractMPIs(t.tlvValue)
	if !ok || len(mpis) < 3 {
		return msg, false
	}
	msg.rb = mpis[0]
	msg.cr = mpis[1]
	msg.d7 = mpis[2]
	return msg, true
}

func toSmpMessageAbort(t tlv) (msg smpMessageAbort, ok bool) {
	return msg, true
}
