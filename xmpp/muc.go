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

//See: Section 6.4
func (m *muc) queryRoomInformation(room *Room) (*data.DiscoveryInfoQuery, error) {
	if room == nil {
		return nil, errors.New("invalid room")
	}

	//Sample data.DiscoveryInfoQuery: We need to parse Forms[], RoomType, Rooms[],
	//info := &data.DiscoveryInfoQuery{
	//	XMLName: xml.Name{Space: "http://jabber.org/protocol/disco#info", Local: "query"},
	//	Node:    "",
	//	Identities: []data.DiscoveryIdentity{
	//		data.DiscoveryIdentity{
	//			XMLName:  xml.Name{Space: "http://jabber.org/protocol/disco#info", Local: "identity"},
	//			Lang:     "",
	//			Category: "conference",
	//			Type:     "text",
	//			Name:     "coyim-test",
	//		},
	//	},
	//	Features: []data.DiscoveryFeature{
	//		data.DiscoveryFeature{XMLName: xml.Name{Space: "http://jabber.org/protocol/disco#info", Local: "feature"},
	//			Var: "http://jabber.org/protocol/muc",
	//		},
	//		data.DiscoveryFeature{XMLName: xml.Name{Space: "http://jabber.org/protocol/disco#info", Local: "feature"},
	//			Var: "muc_unsecured",
	//		},
	//		data.DiscoveryFeature{XMLName: xml.Name{Space: "http://jabber.org/protocol/disco#info", Local: "feature"},
	//			Var: "muc_unmoderated",
	//		},
	//		data.DiscoveryFeature{XMLName: xml.Name{Space: "http://jabber.org/protocol/disco#info", Local: "feature"},
	//			Var: "muc_open",
	//		},
	//		data.DiscoveryFeature{XMLName: xml.Name{Space: "http://jabber.org/protocol/disco#info", Local: "feature"},
	//			Var: "muc_temporary",
	//		},
	//		data.DiscoveryFeature{XMLName: xml.Name{Space: "http://jabber.org/protocol/disco#info", Local: "feature"},
	//			Var: "muc_public",
	//		},
	//		data.DiscoveryFeature{XMLName: xml.Name{Space: "http://jabber.org/protocol/disco#info", Local: "feature"},
	//			Var: "muc_semianonymous",
	//		},
	//	},
	//	Forms: []data.Form{
	//		data.Form{XMLName: xml.Name{Space: "jabber:x:data", Local: "x"},
	//			Type:         "result",
	//			Title:        "",
	//			Instructions: "",
	//			Fields: []data.FormFieldX{
	//				data.FormFieldX{XMLName: xml.Name{Space: "jabber:x:data", Local: "field"},
	//					Desc: "", Var: "FORM_TYPE", Type: "hidden", Label: "",
	//					Required: (*data.FormFieldRequiredX)(nil),
	//					Values:   []string{"http://jabber.org/protocol/muc#roominfo"},
	//					Options:  []data.FormFieldOptionX(nil),
	//					Media:    []data.FormFieldMediaX(nil),
	//				},
	//				data.FormFieldX{XMLName: xml.Name{Space: "jabber:x:data", Local: "field"},
	//					Desc: "", Var: "muc#roominfo_description", Type: "text-single", Label: "Description",
	//					Required: (*data.FormFieldRequiredX)(nil),
	//					Values:   []string{""},
	//					Options:  []data.FormFieldOptionX(nil),
	//					Media:    []data.FormFieldMediaX(nil),
	//				},
	//				data.FormFieldX{XMLName: xml.Name{Space: "jabber:x:data", Local: "field"},
	//					Desc: "", Var: "muc#roominfo_occupants", Type: "text-single", Label: "Number of occupants",
	//					Required: (*data.FormFieldRequiredX)(nil),
	//					Values:   []string{"1"},
	//					Options:  []data.FormFieldOptionX(nil),
	//					Media:    []data.FormFieldMediaX(nil),
	//				},
	//			},
	//		},
	//	},
	//	ResultSet: (*data.ResultSet)(nil),
	//}

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
