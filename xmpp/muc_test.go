package xmpp

import (
	"bytes"
	"encoding/xml"
	"strings"

	"github.com/coyim/coyim/xmpp/data"
	. "gopkg.in/check.v1"
)

type MUCSuite struct{}

var _ = Suite(&MUCSuite{})

func (s *MUCSuite) Test_CanJoinRoom(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	conn := conn{
		log:  testLogger(),
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
		log:  testLogger(),
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

func (s *MUCSuite) Test_CanRequestRoomConfigForm(c *C) {

	//See Example 165. Service Sends Configuration Form to Owner
	expectedResponse := `<iq xmlns='jabber:client' from='coven@chat.shakespeare.lit'
    id='1'
    to='crone1@shakespeare.lit/desktop'
    type='result'>
  <query xmlns='http://jabber.org/protocol/muc#owner'>
    <x xmlns='jabber:x:data' type='form'>
      <title>Configuration for "coven" Room</title>
      <instructions>
        Complete this form to modify the
        configuration of your room.
      </instructions>
      <field
          type='hidden'
          var='FORM_TYPE'>
        <value>http://jabber.org/protocol/muc#roomconfig</value>
      </field>
      <field
          label='Natural-Language Room Name'
          type='text-single'
          var='muc#roomconfig_roomname'>
        <value>A Dark Cave</value>
      </field>
      <field
          label='Short Description of Room'
          type='text-single'
          var='muc#roomconfig_roomdesc'>
        <value>The place for all good witches!</value>
      </field>
      <!-- There is more in the example, but we removed in favor of brevity -->
    </x>
  </query>
</iq>`

	mockIn := xml.NewDecoder(strings.NewReader(expectedResponse))
	mockOut := &mockConnIOReaderWriter{}

	conn := conn{
		log: testLogger(),
		in:  mockIn,
		out: mockOut,

		jid: "crone1@shakespeare.lit/desktop",

		inflights: make(map[data.Cookie]inflight),
		rand:      bytes.NewBuffer([]byte{1, 0, 0, 0, 0, 0, 0, 0}),
	}

	go func() {
		for len(conn.inflights) == 0 {
		}
		conn.Next()
	}()

	result, err := conn.GetChatContext().RequestRoomConfigForm(&data.Room{ID: "coven", Service: "chat.shakespeare.lit"})

	c.Assert(err, IsNil)

	c.Assert(result, DeepEquals, &data.Form{
		XMLName:      xml.Name{Space: "jabber:x:data", Local: "x"},
		Type:         "form",
		Title:        "Configuration for \"coven\" Room",
		Instructions: "\n        Complete this form to modify the\n        configuration of your room.\n      ",
		Fields: []data.FormFieldX{
			data.FormFieldX{
				XMLName: xml.Name{
					Space: "jabber:x:data", Local: "field",
				},
				Var:    "FORM_TYPE",
				Type:   "hidden",
				Values: []string{"http://jabber.org/protocol/muc#roomconfig"},
			},
			data.FormFieldX{
				XMLName: xml.Name{Space: "jabber:x:data", Local: "field"},
				Var:     "muc#roomconfig_roomname",
				Type:    "text-single",
				Label:   "Natural-Language Room Name",
				Values:  []string{"A Dark Cave"},
			},
			data.FormFieldX{
				XMLName: xml.Name{Space: "jabber:x:data", Local: "field"},
				Var:     "muc#roomconfig_roomdesc",
				Type:    "text-single",
				Label:   "Short Description of Room",
				Values:  []string{"The place for all good witches!"},
			},
		},
	})

	c.Assert(string(mockOut.write), Equals,
		`<iq xmlns='jabber:client' `+
			`to='coven@chat.shakespeare.lit' `+
			`from='crone1@shakespeare.lit/desktop' `+
			`type='get' `+
			`id='1'`+
			`>`+
			`<query xmlns="http://jabber.org/protocol/muc#owner"></query>`+
			`</iq>`)
}

func (s *MUCSuite) Test_CanUpdateRoomConfig(c *C) {
	//See Example 159. Owner Submits Configuration Form
	expectedRequest := `<iq xmlns='jabber:client' ` +
		`to='coven@chat.shakespeare.lit' ` +
		`from='crone1@shakespeare.lit/desktop' ` +
		`type='set' ` +
		`id='1'` +
		`>` +
		`<query xmlns="http://jabber.org/protocol/muc#owner">` +
		`<x xmlns="jabber:x:data" type="submit">` +
		`<field var="FORM_TYPE">` +
		`<value>http://jabber.org/protocol/muc#roomconfig</value>` +
		`</field>` +
		`<field var="muc#roomconfig_roomname">` +
		`<value>A Dark Cave</value>` +
		`</field>` +
		`<field var="muc#roomconfig_roomdesc">` +
		`<value>The place for all good witches!</value>` +
		`</field>` +
		`<field var="muc#roomconfig_roomadmins">` +
		`<value>wiccarocks@shakespeare.lit</value>` +
		`<value>hecate@shakespeare.lit</value>` +
		`</field>` +
		`</x>` +
		`</query>` +
		`</iq>`

	mockOut := &mockConnIOReaderWriter{}

	conn := conn{
		log: testLogger(),
		out: mockOut,

		jid: "crone1@shakespeare.lit/desktop",

		inflights: make(map[data.Cookie]inflight),
		rand:      bytes.NewBuffer([]byte{1, 0, 0, 0, 0, 0, 0, 0}),
	}

	err := conn.GetChatContext().UpdateRoomConfig(
		&data.Room{ID: "coven", Service: "chat.shakespeare.lit"},
		&data.Form{
			Type: "submit",
			Fields: []data.FormFieldX{
				data.FormFieldX{
					Var:    "FORM_TYPE",
					Values: []string{"http://jabber.org/protocol/muc#roomconfig"},
				},
				data.FormFieldX{
					Var:    "muc#roomconfig_roomname",
					Values: []string{"A Dark Cave"},
				},
				data.FormFieldX{
					Var:    "muc#roomconfig_roomdesc",
					Values: []string{"The place for all good witches!"},
				},
				data.FormFieldX{
					Var: "muc#roomconfig_roomadmins",
					Values: []string{
						"wiccarocks@shakespeare.lit",
						"hecate@shakespeare.lit",
					},
				},
			},
		})

	c.Assert(err, IsNil)
	c.Assert(string(mockOut.write), Equals, expectedRequest)
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
