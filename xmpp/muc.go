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

//See: Section 4.1
type Room struct {
	ID, Service string
}

func (o *Room) JID() string {
	return fmt.Sprintf("%s@%s", o.ID, o.Service)
}

func parseRoomJID(roomJID string) *Room {
	parts := strings.SplitN(roomJID, "@", 2)
	if len(parts) < 2 {
		return nil
	}

	return &Room{ID: parts[0], Service: parts[1]}
}

//See: Section 4.1
type Occupant struct {
	Room
	Nick string
}

func (o *Occupant) JID() string {
	return fmt.Sprintf("%s/%s", o.Room.JID(), o.Nick)
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

func (m *muc) queryRoomInformation(room *Room) (*data.DiscoveryInfoQuery, error) {
	if room == nil {
		return nil, errors.New("invalid room")
	}

	return m.QueryServiceInformation(room.JID())
}

func (c *conn) enterRoom(roomID, service, nickname string) error {
	occupant := Occupant{Room: Room{ID: roomID, Service: service}, Nick: nickname}
	return c.sendPresenceWithChildren(occupant.JID(), "", "", mucSupport)
}

func (c *conn) leaveRoom(roomID, service, nickname string) error {
	occupant := Occupant{Room: Room{ID: roomID, Service: service}, Nick: nickname}
	return c.sendPresenceWithChildren(occupant.JID(), "unavailable", "", mucSupport)
}
