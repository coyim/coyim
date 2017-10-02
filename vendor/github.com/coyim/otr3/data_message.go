package otr3

import "encoding/binary"

type dataMessageExtra struct {
	key []byte
}

func (c *Conversation) genDataMsg(message []byte, tlvs ...tlv) (dataMsg, dataMessageExtra, error) {
	return c.genDataMsgWithFlag(message, messageFlagNormal, tlvs...)
}

func (c *Conversation) genDataMsgWithFlag(message []byte, flag byte, tlvs ...tlv) (dataMsg, dataMessageExtra, error) {
	if c.msgState != encrypted {
		return dataMsg{}, dataMessageExtra{}, errCannotSendUnencrypted
	}

	keys, err := c.keys.calculateDHSessionKeys(c.keys.ourKeyID-1, c.keys.theirKeyID, c.version)
	if err != nil {
		return dataMsg{}, dataMessageExtra{}, err
	}

	topHalfCtr := [8]byte{}
	counter := c.keys.counterHistory.findCounterFor(c.keys.ourKeyID-1, c.keys.theirKeyID)
	if counter.ourCounter == 0 {
		counter.ourCounter = 1
	}

	binary.BigEndian.PutUint64(topHalfCtr[:], counter.ourCounter)
	counter.ourCounter++

	plain := plainDataMsg{
		message: message,
		tlvs:    tlvs,
	}

	encrypted := plain.encrypt(keys.sendingAESKey[:], topHalfCtr)

	header, err := c.messageHeader(msgTypeData)
	if err != nil {
		return dataMsg{}, dataMessageExtra{}, err
	}

	dataMessage := dataMsg{
		flag:           flag,
		senderKeyID:    c.keys.ourKeyID - 1,
		recipientKeyID: c.keys.theirKeyID,
		y:              c.keys.ourCurrentDHKeys.pub,
		topHalfCtr:     topHalfCtr,
		encryptedMsg:   encrypted,
		oldMACKeys:     c.keys.revealMACKeys(),
	}

	// fmt.Printf("sendingMACKey: len: %d %X\n", len(keys.sendingMACKey), keys.sendingMACKey)
	dataMessage.sign(keys.sendingMACKey, header, c.version)

	c.updateMayRetransmitTo(noRetransmit)
	c.lastMessage(message)

	x := dataMessageExtra{keys.extraKey[:]}

	return dataMessage, x, nil
}

func extractDataMessageFlag(msg []byte) byte {
	if len(msg) == 0 {
		return messageFlagNormal
	}
	return msg[0]
}

func (c *Conversation) createSerializedDataMessage(msg []byte, flag byte, tlvs []tlv) ([]ValidMessage, dataMessageExtra, error) {
	dataMsg, x, err := c.genDataMsgWithFlag(msg, flag, tlvs...)
	if err != nil {
		return nil, dataMessageExtra{}, err
	}

	res, err := c.wrapMessageHeader(msgTypeData, dataMsg.serialize(c.version))
	if err != nil {
		return nil, dataMessageExtra{}, err
	}

	c.updateLastSent()
	return c.fragEncode(res), x, nil
}

func (c *Conversation) fragEncode(msg messageWithHeader) []ValidMessage {
	return c.fragment(c.encode(msg), c.fragmentSize)
}

func (c *Conversation) encode(msg messageWithHeader) encodedMessage {
	return append(append(msgMarker, b64encode(msg)...), '.')
}

func (c *Conversation) processDataMessage(header, msg []byte) (plain MessagePlaintext, toSend messageWithHeader, err error) {
	ignoreUnreadable := (extractDataMessageFlag(msg) & messageFlagIgnoreUnreadable) == messageFlagIgnoreUnreadable
	plain, toSend, err = c.processDataMessageWithRawErrors(header, msg)
	if err != nil && ignoreUnreadable {
		err = nil
	}
	return
}

// processDataMessageWithRawErrors receives a decoded incoming message and returns the plain text inside that message
// and a data message (with header) generated in response to any TLV contained in the incoming message.
// The header and message compose the decoded incoming message.
func (c *Conversation) processDataMessageWithRawErrors(header, msg []byte) (plain MessagePlaintext, toSend messageWithHeader, err error) {
	dataMessage := dataMsg{}

	if c.msgState != encrypted {
		err = errMessageNotInPrivate
		c.messageEvent(MessageEventReceivedMessageNotInPrivate)
		return
	}

	if err = dataMessage.deserialize(msg, c.version); err != nil {
		return
	}

	if err = c.keys.checkMessageCounter(dataMessage); err != nil {
		return
	}

	sessionKeys, err := c.keys.calculateDHSessionKeys(dataMessage.recipientKeyID, dataMessage.senderKeyID, c.version)
	if err != nil {
		return
	}

	if err = dataMessage.checkSign(sessionKeys.receivingMACKey, header, c.version); err != nil {
		return
	}

	p := plainDataMsg{}
	//this can't return an error since receivingAESKey is a AES-128 key
	p.decrypt(sessionKeys.receivingAESKey[:], dataMessage.topHalfCtr, dataMessage.encryptedMsg)

	plain = makeCopy(p.message)
	if len(plain) == 0 {
		plain = nil
		c.messageEvent(MessageEventLogHeartbeatReceived)
	}

	err = c.rotateKeys(dataMessage)
	if err != nil {
		return
	}

	var tlvs []tlv

	tlvs, err = c.processTLVs(p.tlvs, dataMessageExtra{sessionKeys.extraKey})
	if err != nil {
		return
	}

	if len(tlvs) > 0 {
		var reply dataMsg
		reply, _, err = c.genDataMsgWithFlag(nil, decideFlagFrom(tlvs), tlvs...)
		if err != nil {
			return
		}

		toSend, err = c.wrapMessageHeader(msgTypeData, reply.serialize(c.version))
		if err != nil {
			return
		}
	}

	return
}

func decideFlagFrom(tlvs []tlv) byte {
	flag := byte(0x00)
	for _, t := range tlvs {
		if t.tlvType >= tlvTypeSMP1 && t.tlvType <= tlvTypeSMP1WithQuestion {
			flag = messageFlagIgnoreUnreadable
		}

	}
	return flag
}

func (c *Conversation) processSMPTLV(t tlv, x dataMessageExtra) (toSend *tlv, err error) {
	c.smp.ensureSMP()

	smpMessage, ok := t.smpMessage()
	if !ok {
		return nil, newOtrError("corrupt data message")
	}

	return c.receiveSMP(smpMessage)
}

func (c *Conversation) processTLVs(tlvs []tlv, x dataMessageExtra) ([]tlv, error) {
	var retTLVs []tlv

	for _, t := range tlvs {
		mh, e := messageHandlerForTLV(t)
		if e != nil {
			continue
		}

		toSend, err := mh(c, t, x)
		if err != nil {
			//We assume this will only happen if the message was sent by a
			//malicious/broken client and it's reasonable to stop processing the
			//remaining TLVs and consider the entire TLVs block as corrupted.
			//Any valid SMP TLV processed before the error can potentially cause a side
			//effect on the SMP state machine and we wont reply (take the bait).
			return nil, err
		}

		if toSend != nil {
			retTLVs = append(retTLVs, *toSend)
		}
	}

	return retTLVs, nil
}
