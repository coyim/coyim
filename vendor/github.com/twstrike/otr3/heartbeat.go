package otr3

import "time"

// How long after sending a packet should we wait to send a heartbeat?
const heartbeatInterval = 60 * time.Second

type heartbeatContext struct {
	lastSent time.Time
}

func (c *Conversation) updateLastSent() {
	c.heartbeat.lastSent = time.Now()
}

func (c *Conversation) maybeHeartbeat(plain MessagePlaintext, toSend messageWithHeader, err error) (MessagePlaintext, []messageWithHeader, error) {
	if err != nil {
		return nil, nil, err
	}
	tsExtra, e := c.potentialHeartbeat(plain)
	return plain, compactMessagesWithHeader(toSend, tsExtra), e
}

func (c *Conversation) potentialHeartbeat(plain MessagePlaintext) (toSend messageWithHeader, err error) {
	if plain == nil {
		return
	}

	now := time.Now()
	if !c.heartbeat.lastSent.Before(now.Add(-heartbeatInterval)) {
		return
	}

	dataMsg, _, err := c.genDataMsgWithFlag(nil, messageFlagIgnoreUnreadable)
	if err != nil {
		return nil, err
	}

	toSend, err = c.wrapMessageHeader(msgTypeData, dataMsg.serialize(c.version))
	if err != nil {
		return nil, err
	}

	c.updateLastSent()
	c.messageEvent(MessageEventLogHeartbeatSent)
	return
}
