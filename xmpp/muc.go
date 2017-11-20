package xmpp

import (
	"errors"
	"fmt"
	"strings"

	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/interfaces"
)

const (
	mucSupport = "<x xmlns='http://jabber.org/protocol/muc'/>"
	mucNS      = "http://jabber.org/protocol/muc"
)

func parseRoomJID(roomJID string) *data.Room {
	parts := strings.SplitN(roomJID, "@", 2)
	if len(parts) < 2 {
		return nil
	}

	return &data.Room{ID: parts[0], Service: parts[1]}
}

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
	r := parseRoomJID(room)
	if r == nil {
		return nil, errors.New("invalid room")
	}

	//TODO: this error is useless when it says ("expected query, got error")
	//It should give us a OTR error
	query, err := m.queryRoomInformation(r)
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
