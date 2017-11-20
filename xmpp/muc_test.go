package xmpp

import (
	"encoding/xml"

	"github.com/coyim/coyim/xmpp/data"
	. "gopkg.in/check.v1"
)

type MUCSuite struct{}

var _ = Suite(&MUCSuite{})

func (s *MUCSuite) Test_CanJoinRoom(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	conn := conn{
		out:  mockOut,
		rand: &mockConnIOReaderWriter{read: []byte("123555111654")},
	}

	err := conn.GetChatContext().EnterRoom(&data.Occupant{
		Room:   data.Room{ID: "coyim", Service: "chat.coy.im"},
		Handle: "i_am_coy",
	})
	c.Assert(err, IsNil)
	c.Assert(string(mockOut.write), Equals, `<presence xmlns="jabber:client" `+
		`id="3544672884359377457" `+
		`to="coyim@chat.coy.im/i_am_coy">`+
		`<x xmlns='http://jabber.org/protocol/muc'/>`+
		`</presence>`)
}

func (s *MUCSuite) Test_CanLeaveRoom(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	conn := conn{
		out:  mockOut,
		rand: &mockConnIOReaderWriter{read: []byte("123555111654")},
	}

	err := conn.GetChatContext().LeaveRoom(&data.Occupant{
		Room:   data.Room{ID: "coyim", Service: "chat.coy.im"},
		Handle: "i_am_coy",
	})
	c.Assert(err, IsNil)
	c.Assert(string(mockOut.write), Equals, `<presence xmlns="jabber:client" `+
		`id="3544672884359377457" `+
		`to="coyim@chat.coy.im/i_am_coy" `+
		`type="unavailable">`+
		`<x xmlns='http://jabber.org/protocol/muc'/>`+
		`</presence>`)
}

func (s *MUCSuite) Test_parseRoomInfo(c *C) {

	//Sample data.DiscoveryInfoQuery: We need to parse Forms[], RoomType, Rooms[],
	query := &data.DiscoveryInfoQuery{
		XMLName: xml.Name{Space: "http://jabber.org/protocol/disco#info", Local: "query"},
		Node:    "",
		Identities: []data.DiscoveryIdentity{
			data.DiscoveryIdentity{
				XMLName:  xml.Name{Space: "http://jabber.org/protocol/disco#info", Local: "identity"},
				Lang:     "",
				Category: "conference",
				Type:     "text",
				Name:     "coyim-test",
			},
		},
		Features: []data.DiscoveryFeature{
			data.DiscoveryFeature{XMLName: xml.Name{Space: "http://jabber.org/protocol/disco#info", Local: "feature"},
				Var: "http://jabber.org/protocol/muc",
			},
			data.DiscoveryFeature{XMLName: xml.Name{Space: "http://jabber.org/protocol/disco#info", Local: "feature"},
				Var: "muc_unsecured",
			},
			data.DiscoveryFeature{XMLName: xml.Name{Space: "http://jabber.org/protocol/disco#info", Local: "feature"},
				Var: "muc_unmoderated",
			},
			data.DiscoveryFeature{XMLName: xml.Name{Space: "http://jabber.org/protocol/disco#info", Local: "feature"},
				Var: "muc_open",
			},
			data.DiscoveryFeature{XMLName: xml.Name{Space: "http://jabber.org/protocol/disco#info", Local: "feature"},
				Var: "muc_temporary",
			},
			data.DiscoveryFeature{XMLName: xml.Name{Space: "http://jabber.org/protocol/disco#info", Local: "feature"},
				Var: "muc_public",
			},
			data.DiscoveryFeature{XMLName: xml.Name{Space: "http://jabber.org/protocol/disco#info", Local: "feature"},
				Var: "muc_semianonymous",
			},
		},
		Forms: []data.Form{
			data.Form{XMLName: xml.Name{Space: "jabber:x:data", Local: "x"},
				Type:         "result",
				Title:        "",
				Instructions: "",
				Fields: []data.FormFieldX{
					data.FormFieldX{XMLName: xml.Name{Space: "jabber:x:data", Local: "field"},
						Desc: "", Var: "FORM_TYPE", Type: "hidden", Label: "",
						Required: (*data.FormFieldRequiredX)(nil),
						Values:   []string{"http://jabber.org/protocol/muc#roominfo"},
						Options:  []data.FormFieldOptionX(nil),
						Media:    []data.FormFieldMediaX(nil),
					},
					data.FormFieldX{XMLName: xml.Name{Space: "jabber:x:data", Local: "field"},
						Desc: "", Var: "muc#roominfo_description", Type: "text-single", Label: "Description",
						Required: (*data.FormFieldRequiredX)(nil),
						Values:   []string{"CoyIM testing room"},
						Options:  []data.FormFieldOptionX(nil),
						Media:    []data.FormFieldMediaX(nil),
					},
					data.FormFieldX{XMLName: xml.Name{Space: "jabber:x:data", Local: "field"},
						Desc: "", Var: "muc#roominfo_occupants", Type: "text-single", Label: "Number of occupants",
						Required: (*data.FormFieldRequiredX)(nil),
						Values:   []string{"1"},
						Options:  []data.FormFieldOptionX(nil),
						Media:    []data.FormFieldMediaX(nil),
					},
				},
			},
		},
		ResultSet: (*data.ResultSet)(nil),
	}

	roomInfo := parseRoomInformation(query)

	c.Assert(roomInfo.Description, Equals, "CoyIM testing room")
	c.Assert(roomInfo.Occupants, Equals, 1)
	c.Assert(roomInfo.PasswordProtected, Equals, false)
	c.Assert(roomInfo.Moderated, Equals, false)
	c.Assert(roomInfo.Open, Equals, true)
	c.Assert(roomInfo.Persistent, Equals, false)
	c.Assert(roomInfo.Public, Equals, true)
	c.Assert(roomInfo.SemiAnonymous, Equals, true)
}
