package xmpp

import (
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/interfaces"
)

const (
	mucSupport = "<x xmlns='http://jabber.org/protocol/muc'/>"
	mucNS      = "http://jabber.org/protocol/muc"
)

func (c *conn) GetChatContext() interfaces.Chat {
	return &muc{c}
}

type muc struct {
	*conn
}

//See: Section "6.2 Discovering the Features Supported by a MUC Service"
func (m *muc) CheckForSupport(entity string) bool {
	return m.HasSupportTo(entity, mucNS)
}

//See: Section "6.3 Discovering Rooms"
func (m *muc) QueryRooms(entity string) ([]data.DiscoveryItem, error) {
	query, err := m.QueryServiceItems(entity)
	if err != nil {
		return nil, err
	}

	return query.DiscoveryItems, nil
}

//See: Section "6.4 Querying for Room Information"
func (m *muc) QueryRoomInformation(room string) (*data.RoomInfo, error) {
	jid := ParseJID(room)
	if jid.LocalPart == "" || jid.DomainPart == "" {
		return nil, errors.New("invalid room")
	}

	//TODO: this error is useless when it says ("expected query, got error")
	//It should give us a xmpp error
	query, err := m.queryRoomInformation(&data.Room{
		ID:      jid.LocalPart,
		Service: jid.DomainPart,
	})

	if err != nil {
		return nil, err
	}

	return parseRoomInformation(query), nil
}

func parseRoomInfoForm(forms []data.Form) data.RoomInfoForm {
	ret := data.RoomInfoForm{}
	parseForms(&ret, forms)
	return ret
}

func parseRoomType(features []data.DiscoveryFeature) data.RoomType {
	ret := data.RoomType{}

	for _, f := range features {
		switch f.Var {
		case "muc_public":
			ret.Public = true
		case "muc_open":
			ret.Open = true
		case "muc_moderated":
			ret.Moderated = true
		case "muc_semianonymous":
			ret.SemiAnonymous = true
		case "muc_passwordprotected":
			ret.PasswordProtected = true
		case "muc_persistenc":
			ret.Persistent = true
		}
	}

	return ret
}

func parseRoomInformation(query *data.DiscoveryInfoQuery) *data.RoomInfo {
	return &data.RoomInfo{
		RoomInfoForm: parseRoomInfoForm(query.Forms[:]),
		RoomType:     parseRoomType(query.Features),
	}
}

func (m *muc) queryRoomInformation(room *data.Room) (*data.DiscoveryInfoQuery, error) {
	return m.QueryServiceInformation(room.JID())
}

//See: Section "7.2.2 Basic MUC Protocol"
func (m *muc) EnterRoom(occupant *data.Occupant) error {
	//TODO: Implement section "7.2.1 Groupchat 1.0 Protocol"?
	return m.sendPresence(&data.ClientPresence{
		To:    occupant.JID(),
		Extra: mucSupport,
	})
}

//See: Section "7.14 Exiting a Room"
func (c *muc) LeaveRoom(occupant *data.Occupant) error {
	return c.sendPresence(&data.ClientPresence{
		To:    occupant.JID(),
		Type:  "unavailable",
		Extra: mucSupport,
	})
}

//See: Section "7.4 Sending a Message to All Occupants"
func (m *muc) SendChatMessage(msg string, to *data.Room) error {
	//TODO: How to disable archive for chat messages?
	//TODO: Can we just use the same conn.Send() with a different type?
	_, err := fmt.Fprintf(m.out, "<message "+
		"to='%s' "+
		"from='%s' "+
		"type='groupchat'>"+
		"<body>%s</body>"+
		"</message>",
		xmlEscape(to.JID()), xmlEscape(m.conn.jid), xmlEscape(msg))
	return err
}

//See: Section "10.2 Subsequent Room Configuration"
func (m *muc) RequestRoomConfigForm(room *data.Room) (*data.Form, error) {
	reply, _, err := m.SendIQ(room.JID(), "get", &data.RoomConfigurationQuery{})

	stanza, ok := <-reply
	if !ok {
		return nil, errors.New("xmpp: failed to receive response")
	}

	iq, ok := stanza.Value.(*data.ClientIQ)
	if !ok {
		return nil, errors.New("xmpp: failed to parse response")
	}

	r := &data.RoomConfigurationQuery{}
	err = xml.Unmarshal(iq.Query, r)
	return r.Form, err
}

//See: Section "10.2 Subsequent Room Configuration"
func (m *muc) UpdateRoomConfig(room *data.Room, form *data.Form) error {
	_, _, err := m.SendIQ(room.JID(), "set", &data.RoomConfigurationQuery{
		Form: form,
	})

	return err
}
