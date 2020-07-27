package xmpp

import (
	"encoding/xml"
	"errors"

	"github.com/coyim/coyim/xmpp/data"
)

// TODO: Add a RoomConfigurationQuery for create a Reserved Room
func (m *muc) LegacyOldDoNotUseCreateRoom(room *data.LegacyOldDoNotUseRoom) error {
	// Send a presence for create the room and signals supportfor MUC
	// See: 10.1.1 Create room General
	p := &data.ClientPresence{
		To:    room.JID(),
		Extra: mucSupport,
	}

	err := m.conn.sendPresence(p)
	if err != nil {
		return err
	}

	// TODO: Delete 'roomConf' and get this information from the function
	// See: 10.1.2 Creating an Instant Room
	// Information Query
	roomConf := &data.LegacyOldDoNotUseRoomConfigurationQuery{
		XMLName: xml.Name{
			Local: "query",
			Space: "http://jabber.org/protocol/muc#owner",
		},
		Form: &data.Form{
			XMLName: xml.Name{
				Local: "x",
				Space: "jabber:x:data",
			},
			Type: "submit",
		},
	}

	// Send an IQ with the room information
	reply, _, err := m.SendIQ(room.JID(), "set", roomConf)
	if err != nil {
		return err
	}

	stanza, ok := <-reply
	if !ok {
		return errors.New("xmpp: failed to receive response")
	}

	iq, ok := stanza.Value.(*data.ClientIQ)
	if !ok {
		return errors.New("xmpp: failed to parse response")
	}

	r := &data.LegacyOldDoNotUseRoomConfigurationQuery{}
	err = xml.Unmarshal(iq.Query, r)
	return err
}
