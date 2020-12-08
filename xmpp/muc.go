package xmpp

import (
	"encoding/xml"
	"errors"

	"github.com/coyim/coyim/xmpp/data"
)

func (c *conn) SendConfigurationFormRequest(to string) (*data.MUCRoomConfiguration, error) {
	reply, _, err := c.SendIQ(to, "get", data.MUCRoomConfiguration{})
	if err != nil {
		return nil, err
	}

	stanza, ok := <-reply
	if !ok {
		return nil, errors.New("xmpp: failed to receive response")
	}

	iq, ok := stanza.Value.(*data.ClientIQ)
	if !ok {
		return nil, errors.New("xmpp: failed to parse response")
	}

	if iq.Type == "error" {
		return nil, errors.New("xmpp: tag error received")
	}
	return parseRoomFormReply(iq)
}

func parseRoomFormReply(iq *data.ClientIQ) (*data.MUCRoomConfiguration, error) {
	reply := &data.MUCRoomConfiguration{}
	err := xml.Unmarshal(iq.Query, reply)
	return reply, err
}
