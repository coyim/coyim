package otr3

import (
	"bufio"
	"bytes"
)

// Send takes a human readable message from the local user, possibly encrypts
// it and returns zero or more messages to send to the peer.
func (c *Conversation) Send(m ValidMessage, trace ...interface{}) ([]ValidMessage, error) {
	message := makeCopy(m)
	defer wipeBytes(message)

	if !c.Policies.isOTREnabled() {
		return []ValidMessage{makeCopy(message)}, nil
	}

	if c.debug && bytes.Index(message, []byte(debugString)) != -1 {
		c.dump(bufio.NewWriter(standardErrorOutput))
		return nil, nil
	}

	switch c.msgState {
	case plainText:
		return c.withInjections(c.sendMessageOnPlaintext(message, trace...))
	case encrypted:
		return c.withInjections(c.sendMessageOnEncrypted(message))
	case finished:
		c.messageEvent(MessageEventConnectionEnded)
		return c.withInjections(nil, newOtrError("cannot send message because secure conversation has finished"))
	}

	return c.withInjections(nil, newOtrError("cannot send message in current state"))
}

func (c *Conversation) sendMessageOnPlaintext(message ValidMessage, trace ...interface{}) ([]ValidMessage, error) {
	if c.Policies.has(requireEncryption) {
		c.messageEvent(MessageEventEncryptionRequired, trace...)
		c.updateLastSent()
		c.updateMayRetransmitTo(retransmitExact)
		c.lastMessage(MessagePlaintext(makeCopy(message)), trace...)
		return []ValidMessage{c.QueryMessage()}, nil
	}

	return []ValidMessage{makeCopy(c.appendWhitespaceTag(message))}, nil
}

func (c *Conversation) sendMessageOnEncrypted(message ValidMessage) ([]ValidMessage, error) {
	result, _, err := c.createSerializedDataMessage(message, messageFlagNormal, []tlv{})
	if err != nil {
		c.messageEvent(MessageEventEncryptionError)
		c.generatePotentialErrorMessage(ErrorCodeEncryptionError)
	}

	return result, err
}

func (c *Conversation) sendDHCommit() (toSend messageWithHeader, err error) {
	c.ake.wipe(true)
	c.ake = nil

	toSend, err = c.dhCommitMessage()
	if err != nil {
		return
	}

	toSend, err = c.wrapMessageHeader(msgTypeDHCommit, toSend)
	if err != nil {
		return nil, err
	}

	c.ake.state = authStateAwaitingDHKey{}

	return
}
