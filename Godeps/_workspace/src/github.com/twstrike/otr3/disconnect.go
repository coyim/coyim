package otr3

import "time"

func (c *Conversation) processDisconnectedTLV(t tlv, x dataMessageExtra) (toSend *tlv, err error) {
	previousMsgState := c.msgState

	defer c.signalSecurityEventIf(previousMsgState == encrypted, GoneInsecure)
	c.lastMessageStateChange = time.Now()
	c.msgState = finished
	c.smp.wipe()

	c.keys = keyManagementContext{}

	return nil, nil
}
