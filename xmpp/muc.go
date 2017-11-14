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

//See: Section 4.2
type RoomType struct {
	Public bool
	//vs Hidden bool

	Open bool
	//vs MembersOnly bool

	Moderated bool
	//vs Unmoderated bool

	SemiAnonymous bool
	//vs NonAnonymous bool

	PasswordProtected bool
	//vs Unsecured bool

	Persistent bool
	//vs Temporary bool
}

//See: Section 15.3 Service Discovery Features
type MUCFeatures struct {
	Rooms []string
	RoomType
}

// See: Section 15.5.4 muc#roominfo FORM_TYPE
type RoomInfoForm struct {
	MaxHistoryFetch string   `form-field:"muc#maxhistoryfetch"`
	ContactJID      []string `form-field:"muc#roominfo_contactjid"`
	Description     string   `form-field:"muc#roominfo_description"`
	Language        string   `form-field:"muc#roominfo_language"`
	LDAPGroup       string   `form-field:"muc#roominfo_ldapgroup"`
	Logs            string   `form-field:"muc#roominfo_logs"`
	Occupants       int      `form-field:"muc#roominfo_occupants"`
	Subject         string   `form-field:"muc#roominfo_subject"`
	SubjectMod      bool     `form-field:"muc#roominfo_subjectmod"`
}

//TODO: Ahh, naming
type RoomInfo struct {
	RoomInfoForm `form-type:"http://jabber.org/protocol/muc#roominfo"`
	RoomType
}

func (c *conn) GetChatContext() interfaces.Chat {
	return &muc{c}
}

type muc struct {
	*conn
}

//See: Section 6.2
func (m *muc) CheckForSupport(entity string) bool {
	return m.HasSupportTo(entity, mucNS)
}

func (m *muc) QueryRoomInformation(room string) (*data.DiscoveryInfoQuery, error) {
	r := parseRoomJID(room)
	return m.queryRoomInformation(r)
}

func parseRoomInfoForm(forms []data.Form) RoomInfoForm {
	ret := RoomInfoForm{}
	parseForms(&ret, forms)
	return ret
}

func parseRoomType(features []data.DiscoveryFeature) RoomType {
	ret := RoomType{}

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

func parseRoomInformation(query *data.DiscoveryInfoQuery) *RoomInfo {
	return &RoomInfo{
		RoomInfoForm: parseRoomInfoForm(query.Forms[:]),
		RoomType:     parseRoomType(query.Features),
	}
}

//See: Section 6.4
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
