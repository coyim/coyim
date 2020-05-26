package otr3

func (c *Conversation) processExtraSymmetricKeyTLV(t tlv, x dataMessageExtra) (toSend *tlv, err error) {
	rest, usage, ok := ExtractWord(t.tlvValue[:t.tlvLength])
	if ok {
		c.receivedSymKey(usage, rest, x.key)
	}
	return nil, nil
}

// UseExtraSymmetricKey takes a usage parameter and optional usageData and returns the current symmetric key
// and a set of messages to send in order to ask the peer to use the same symmetric key for the usage defined
func (c *Conversation) UseExtraSymmetricKey(usage uint32, usageData []byte) ([]byte, []ValidMessage, error) {
	if c.msgState != encrypted ||
		c.keys.theirKeyID == 0 {
		return nil, nil, newOtrError("cannot send message in current state")
	}

	t := tlv{
		tlvType:   tlvTypeExtraSymmetricKey,
		tlvLength: 4 + uint16(len(usageData)),
		tlvValue:  append(AppendWord(nil, usage), usageData...),
	}

	toSend, x, err := c.createSerializedDataMessage(nil, messageFlagIgnoreUnreadable, []tlv{t})
	return x.key, toSend, err
}

// ReceivedKeyHandler is an interface that will be invoked when an extra key is received
type ReceivedKeyHandler interface {
	// ReceivedSymmetricKey will be called when a TLV requesting the use of a symmetric key is received
	ReceivedSymmetricKey(usage uint32, usageData []byte, symkey []byte)
}

type dynamicReceivedKeyHandler struct {
	eh func(usage uint32, usageData []byte, symkey []byte)
}

func (d dynamicReceivedKeyHandler) ReceivedSymmetricKey(usage uint32, usageData []byte, symkey []byte) {
	d.eh(usage, usageData, symkey)
}

func (c *Conversation) receivedSymKey(usage uint32, usageData []byte, symkey []byte) {
	if c.receivedKeyHandler != nil {
		c.receivedKeyHandler.ReceivedSymmetricKey(usage, usageData, symkey)
	}
}
