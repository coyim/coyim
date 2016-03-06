package otr3

// StartAuthenticate should be called when the user wants to initiate authentication with a peer.
// The authentication uses an optional question message and a shared secret. The authentication will proceed
// until the event handler reports that SMP is complete, that a secret is needed or that SMP has failed.
func (c *Conversation) StartAuthenticate(question string, mutualSecret []byte) ([]ValidMessage, error) {
	c.smp.ensureSMP()

	tlvs, err := c.smp.state.startAuthenticate(c, question, mutualSecret)

	if err != nil {
		return nil, err
	}

	msgs, _, err := c.createSerializedDataMessage(nil, messageFlagIgnoreUnreadable, tlvs)
	return msgs, err
}

// ProvideAuthenticationSecret should be called when the peer has started an authentication request, and the UI has been notified that a secret is needed
// It is only valid to call this function if the current SMP state is waiting for a secret to be provided. The return is the potential messages to send.
func (c *Conversation) ProvideAuthenticationSecret(mutualSecret []byte) ([]ValidMessage, error) {
	t, err := c.continueSMP(mutualSecret)
	if err != nil {
		return nil, err
	}

	msgs, _, err := c.createSerializedDataMessage(nil, messageFlagIgnoreUnreadable, []tlv{*t})
	return msgs, err
}

func (c *Conversation) potentialAuthError(toSend []messageWithHeader, err error) ([]messageWithHeader, error) {
	if err != nil {
		c.messageEventWithError(MessageEventSetupError, err)
	}

	return toSend, err
}
