package otr3

// Receive handles a message from a peer. It returns a human readable message and zero or more messages to send back to the peer.
func (c *Conversation) Receive(m ValidMessage) (plain MessagePlaintext, toSend []ValidMessage, err error) {
	return c.receiveUnit(m, true)
}

// Receive handles a message from a peer. It returns a human readable message and zero or more messages to send back to the peer.
func (c *Conversation) receiveUnit(m ValidMessage, forgetFragments bool) (plain MessagePlaintext, toSend []ValidMessage, err error) {
	message := makeCopy(m)
	defer wipeBytes(message)

	if !c.Policies.isOTREnabled() {
		return c.receiveWithoutOTR(message)
	}

	msgType := guessMessageType(message)
	var messagesToSend []messageWithHeader
	shouldForgetFragment := true
	switch msgType {
	case msgGuessError:
		return c.withInjectionsPlain(c.receiveErrorMessage(message))
	case msgGuessQuery:
		messagesToSend, err = c.receiveQueryMessage(message)
	case msgGuessTaggedPlaintext:
		plain, messagesToSend, err = c.receiveTaggedPlaintext(message)
	case msgGuessNotOTR:
		plain, messagesToSend, err = c.receivePlaintext(message)
	case msgGuessV1KeyExch:
		return nil, nil, errUnsupportedOTRVersion
	case msgGuessFragment:
		shouldForgetFragment = false
		c.fragmentationContext, err = c.receiveFragment(c.fragmentationContext, message)
		if fragmentsFinished(c.fragmentationContext) {
			return c.withInjectionsPlain(c.receiveUnit(c.fragmentationContext.frag, false))
		}
	case msgGuessUnknown:
		c.messageEvent(MessageEventReceivedMessageUnrecognized)
	case msgGuessDHCommit, msgGuessDHKey, msgGuessRevealSig, msgGuessSignature, msgGuessData:
		plain, messagesToSend, err = c.receiveEncoded(encodedMessage(message))
	}

	if shouldForgetFragment && forgetFragments {
		c.fragmentationContext = forgetFragment()
	}

	return c.withInjectionsPlain(c.toSendEncoded(plain, messagesToSend, err))
}

func (c *Conversation) receiveWithoutOTR(message ValidMessage) (MessagePlaintext, []ValidMessage, error) {
	return MessagePlaintext(message), nil, nil
}

func withoutPotentialSpaceStart(msg []byte) []byte {
	if len(msg) > 0 && msg[0] == ' ' {
		return msg[1:]
	}
	return msg
}

func (c *Conversation) receiveErrorMessage(message ValidMessage) (plain MessagePlaintext, toSend []ValidMessage, err error) {
	msg := MessagePlaintext(makeCopy(message[len(errorMarker):]))

	if c.Policies.has(errorStartAKE) {
		toSend = []ValidMessage{c.QueryMessage()}
	}

	if c.msgState == encrypted {
		c.updateMayRetransmitTo(retransmitWithPrefix)
	}

	c.messageEventWithMessage(MessageEventReceivedMessageGeneralError, withoutPotentialSpaceStart(msg))
	return
}

func (c *Conversation) encodeAndCombine(toSend []messageWithHeader) []ValidMessage {
	var result []ValidMessage

	for _, ts := range toSend {
		result = append(result, c.fragEncode(ts)...)
	}

	return result
}

func (c *Conversation) toSendEncoded(plain MessagePlaintext, toSend []messageWithHeader, err error) (MessagePlaintext, []ValidMessage, error) {
	if err != nil || len(toSend) == 0 || len(toSend[0]) == 0 {
		return plain, nil, err
	}

	return plain, c.encodeAndCombine(toSend), err
}

func (c *Conversation) receiveEncoded(message encodedMessage) (MessagePlaintext, []messageWithHeader, error) {
	decodedMessage, err := c.decode(message)
	if err != nil {
		return nil, nil, err
	}
	return c.receiveDecoded(decodedMessage)
}

func (c *Conversation) checkPlaintextPolicies(plain MessagePlaintext) {
	if c.whitespaceState == whitespaceSent {
		c.whitespaceState = whitespaceRejected
	}

	if c.msgState != plainText || c.Policies.has(requireEncryption) {
		c.messageEventWithMessage(MessageEventReceivedMessageUnencrypted, plain)
	}
}

func (c *Conversation) receivePlaintext(message ValidMessage) (plain MessagePlaintext, toSend []messageWithHeader, err error) {
	p := makeCopy(message)
	plain = MessagePlaintext(p)
	c.checkPlaintextPolicies(plain)
	return
}

func (c *Conversation) receiveTaggedPlaintext(message ValidMessage) (plain MessagePlaintext, toSend []messageWithHeader, err error) {
	plain, toSend, err = c.processWhitespaceTag(message)
	c.checkPlaintextPolicies(plain)
	return
}

func removeOTRMsgEnvelope(msg encodedMessage) []byte {
	return msg[len(msgMarker) : len(msg)-1]
}

func (c *Conversation) decode(encoded encodedMessage) (messageWithHeader, error) {
	encoded = removeOTRMsgEnvelope(encoded)
	msg, err := b64decode(encoded)

	if err != nil {
		return nil, errInvalidOTRMessage
	}

	return msg, nil
}

func (c *Conversation) receiveDecoded(message messageWithHeader) (plain MessagePlaintext, toSend []messageWithHeader, err error) {
	if err = c.checkVersion(message); err != nil {
		return
	}

	var messageHeader, messageBody []byte
	if messageHeader, messageBody, err = c.parseMessageHeader(message); err != nil {
		if err == errReceivedMessageForOtherInstance {
			err = nil
		}
		return
	}

	msgType := messageHeader[2]
	switch msgType {
	case msgTypeData:
		return c.receiveDataMessage(messageHeader, messageBody)
	default:
		return c.receiveAKEMessage(msgType, messageBody)
	}
}

func (c *Conversation) receiveAKEMessage(msgType byte, messageBody []byte) (plain MessagePlaintext, toSend []messageWithHeader, err error) {
	toSend, err = c.potentialAuthError(c.processAKE(msgType, messageBody))
	return
}

func (c *Conversation) receiveDataMessage(messageHeader, messageBody []byte) (plain MessagePlaintext, toSend []messageWithHeader, err error) {
	plain, toSend, err = c.maybeHeartbeat(c.processDataMessage(messageHeader, messageBody))
	if err != nil {
		c.notifyDataMessageError(err)
	}

	return
}

func (c *Conversation) notifyDataMessageError(err error) {
	var e ErrorCode

	if err == errMessageNotInPrivate {
		return
	}

	if isConflict(err) {
		c.messageEvent(MessageEventReceivedMessageUnreadable)
		e = ErrorCodeMessageUnreadable
	} else {
		c.messageEvent(MessageEventReceivedMessageMalformed)
		e = ErrorCodeMessageMalformed
	}

	c.generatePotentialErrorMessage(e)
}
