package muc

import (
	"bytes"
	"io"
	"os"

	"github.com/coyim/coyim/session/muc/data"
	xmppData "github.com/coyim/coyim/xmpp/data"
	. "gopkg.in/check.v1"
)

type MucRoomListingSuite struct{}

var _ = Suite(&MucRoomListingSuite{})

func (*MucRoomListingSuite) Test_NewRoomListing(c *C) {
	c.Assert(NewRoomListing(), Not(IsNil))
}

func (*MucRoomListingSuite) Test_RoomListing_GetDiscoInfo(c *C) {
	rl := &RoomListing{}
	dd := data.RoomDiscoInfo{
		AnonymityLevel: "stuff",
	}

	rl.RoomDiscoInfo = dd

	c.Assert(rl.GetDiscoInfo(), DeepEquals, dd)
}

func (*MucRoomListingSuite) Test_RoomListing_Updated(c *C) {
	rl := &RoomListing{}

	rl.Updated()

	called1 := false
	var data1 interface{}
	rl.OnUpdate(func(_ *RoomListing, dd interface{}) {
		called1 = true
		data1 = dd
	}, "data1")

	called2 := false
	var data2 interface{}
	rl.OnUpdate(func(_ *RoomListing, dd interface{}) {
		called2 = true
		data2 = dd
	}, "data2")

	rl.Updated()

	c.Assert(called1, Equals, true)
	c.Assert(called2, Equals, true)
	c.Assert(data1, DeepEquals, "data1")
	c.Assert(data2, DeepEquals, "data2")
}

func (*MucRoomListingSuite) Test_RoomListing_SetFeatures(c *C) {
	rl := &RoomListing{}

	rl.SetFeatures([]xmppData.DiscoveryFeature{})

	rl.SetFeatures([]xmppData.DiscoveryFeature{
		xmppData.DiscoveryFeature{Var: "http://jabber.org/protocol/muc"},
		xmppData.DiscoveryFeature{Var: "http://jabber.org/protocol/muc"},
		xmppData.DiscoveryFeature{Var: "http://jabber.org/protocol/muc#stable_id"},
		xmppData.DiscoveryFeature{Var: "http://jabber.org/protocol/muc#self-ping-optimization"},
		xmppData.DiscoveryFeature{Var: "http://jabber.org/protocol/disco#info"},
		xmppData.DiscoveryFeature{Var: "http://jabber.org/protocol/disco#items"},
		xmppData.DiscoveryFeature{Var: "urn:xmpp:mam:0"},
		xmppData.DiscoveryFeature{Var: "urn:xmpp:mam:1"},
		xmppData.DiscoveryFeature{Var: "urn:xmpp:mam:2"},
		xmppData.DiscoveryFeature{Var: "urn:xmpp:mam:tmp"},
		xmppData.DiscoveryFeature{Var: "urn:xmpp:mucsub:0"},
		xmppData.DiscoveryFeature{Var: "urn:xmpp:sid:0"},
		xmppData.DiscoveryFeature{Var: "vcard-temp"},
		xmppData.DiscoveryFeature{Var: "http://jabber.org/protocol/muc#request"},
		xmppData.DiscoveryFeature{Var: "jabber:iq:register"},
		xmppData.DiscoveryFeature{Var: "muc_semianonymous"},
		xmppData.DiscoveryFeature{Var: "muc_persistent"},
		xmppData.DiscoveryFeature{Var: "muc_unmoderated"},
		xmppData.DiscoveryFeature{Var: "muc_open"},
		xmppData.DiscoveryFeature{Var: "muc_passwordprotected"},
		xmppData.DiscoveryFeature{Var: "muc_public"},
	})

	c.Assert(rl.SupportsVoiceRequests, Equals, true)
	c.Assert(rl.AllowsRegistration, Equals, true)
	c.Assert(rl.AnonymityLevel, Equals, "semi")
	c.Assert(rl.Persistent, Equals, true)
	c.Assert(rl.Moderated, Equals, false)
	c.Assert(rl.Open, Equals, true)
	c.Assert(rl.PasswordProtected, Equals, true)
	c.Assert(rl.Public, Equals, true)

	rl.SetFeatures([]xmppData.DiscoveryFeature{
		xmppData.DiscoveryFeature{Var: "muc_nonanonymous"},
		xmppData.DiscoveryFeature{Var: "muc_temporary"},
		xmppData.DiscoveryFeature{Var: "muc_moderated"},
		xmppData.DiscoveryFeature{Var: "muc_membersonly"},
		xmppData.DiscoveryFeature{Var: "muc_unsecured"},
		xmppData.DiscoveryFeature{Var: "muc_hidden"},
	})

	c.Assert(rl.AnonymityLevel, Equals, "no")
	c.Assert(rl.Public, Equals, false)
	c.Assert(rl.PasswordProtected, Equals, false)
	c.Assert(rl.Open, Equals, false)
	c.Assert(rl.Moderated, Equals, true)
	c.Assert(rl.Persistent, Equals, false)
}

func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func (*MucRoomListingSuite) Test_RoomListing_setFeature_unknownFeature(c *C) {
	rl := &RoomListing{}
	stdout := captureStdout(func() { rl.setFeature("something unknown") })
	c.Assert(stdout, Equals, "UNKNOWN FEATURE: something unknown\n")
}

func (*MucRoomListingSuite) Test_RoomListing_SetFormsData(c *C) {
	rl := &RoomListing{}
	rl.SetFormsData(nil)

	rl.SetFormsData([]xmppData.Form{
		xmppData.Form{
			Type: "result",
			Fields: []xmppData.FormFieldX{
				xmppData.FormFieldX{Var: "FORM_TYPE", Values: []string{discoInfoFieldFormType}},
				xmppData.FormFieldX{Var: discoInfoFieldLang, Values: []string{"eng"}},
			},
		},
	})

	c.Assert(rl.Language, Equals, "eng")
}

func (*MucRoomListingSuite) Test_RoomListing_updateWithFormField(c *C) {
	rl := &RoomListing{}

	rl.Language = "swe"
	rl.updateWithFormField("muc#roominfo_lang", []string{})
	c.Assert(rl.Language, Equals, "swe")
	rl.updateWithFormField("muc#roominfo_lang", []string{"en", "something"})
	c.Assert(rl.Language, Equals, "en")

	rl.OccupantsCanChangeSubject = false
	rl.updateWithFormField("muc#roominfo_changesubject", []string{"1"})
	c.Assert(rl.OccupantsCanChangeSubject, Equals, true)
	rl.updateWithFormField("muc#roominfo_changesubject", []string{"0", "1"})
	c.Assert(rl.OccupantsCanChangeSubject, Equals, false)

	rl.Logged = false
	rl.updateWithFormField("muc#roomconfig_enablelogging", []string{"1"})
	c.Assert(rl.Logged, Equals, true)
	rl.updateWithFormField("muc#roomconfig_enablelogging", []string{})
	c.Assert(rl.Logged, Equals, true)
	rl.updateWithFormField("muc#roomconfig_enablelogging", []string{"0", "1"})
	c.Assert(rl.Logged, Equals, false)

	rl.Title = "hello"
	rl.updateWithFormField("muc#roomconfig_roomname", []string{})
	c.Assert(rl.Title, Equals, "")
	rl.updateWithFormField("muc#roomconfig_roomname", []string{"something", "foo"})
	c.Assert(rl.Title, Equals, "something")

	rl.Description = "hello"
	rl.updateWithFormField("muc#roominfo_description", []string{})
	c.Assert(rl.Description, Equals, "hello")
	rl.updateWithFormField("muc#roominfo_description", []string{"something", "foo"})
	c.Assert(rl.Description, Equals, "something")

	rl.Occupants = 42
	rl.updateWithFormField("muc#roominfo_occupants", []string{})
	c.Assert(rl.Occupants, Equals, 42)
	rl.updateWithFormField("muc#roominfo_occupants", []string{"xq"})
	c.Assert(rl.Occupants, Equals, 42)
	rl.updateWithFormField("muc#roominfo_occupants", []string{"55"})
	c.Assert(rl.Occupants, Equals, 55)

	rl.MembersCanInvite = false
	rl.updateWithFormField("{http://prosody.im/protocol/muc}roomconfig_allowmemberinvites", []string{"1"})
	c.Assert(rl.MembersCanInvite, Equals, true)
	rl.updateWithFormField("{http://prosody.im/protocol/muc}roomconfig_allowmemberinvites", []string{})
	c.Assert(rl.MembersCanInvite, Equals, true)
	rl.updateWithFormField("{http://prosody.im/protocol/muc}roomconfig_allowmemberinvites", []string{"0", "1"})
	c.Assert(rl.MembersCanInvite, Equals, false)

	rl.OccupantsCanInvite = false
	rl.updateWithFormField("muc#roomconfig_allowinvites", []string{"1"})
	c.Assert(rl.OccupantsCanInvite, Equals, true)
	rl.updateWithFormField("muc#roomconfig_allowinvites", []string{})
	c.Assert(rl.OccupantsCanInvite, Equals, true)
	rl.updateWithFormField("muc#roomconfig_allowinvites", []string{"0", "1"})
	c.Assert(rl.OccupantsCanInvite, Equals, false)

	rl.AllowPrivateMessages = "somewhere"
	rl.updateWithFormField("muc#roomconfig_allowpm", []string{})
	c.Assert(rl.AllowPrivateMessages, Equals, "somewhere")
	rl.updateWithFormField("muc#roomconfig_allowpm", []string{"something", "foo"})
	c.Assert(rl.AllowPrivateMessages, Equals, "something")

	rl.ContactJid = "somewhere"
	rl.updateWithFormField("muc#roominfo_contactjid", []string{})
	c.Assert(rl.ContactJid, Equals, "somewhere")
	rl.updateWithFormField("muc#roominfo_contactjid", []string{"something", "foo"})
	c.Assert(rl.ContactJid, Equals, "something")

	rl.MaxHistoryFetch = 42
	rl.updateWithFormField("muc#maxhistoryfetch", []string{})
	c.Assert(rl.MaxHistoryFetch, Equals, 42)
	rl.updateWithFormField("muc#maxhistoryfetch", []string{"xq"})
	c.Assert(rl.MaxHistoryFetch, Equals, 42)
	rl.updateWithFormField("muc#maxhistoryfetch", []string{"55"})
	c.Assert(rl.MaxHistoryFetch, Equals, 55)
}

func (*MucRoomListingSuite) Test_RoomListing_updateWithFormField_fails(c *C) {
	rl := &RoomListing{}
	stdout := captureStdout(func() { rl.updateWithFormField("something unknown", []string{}) })
	c.Assert(stdout, Equals, "UNKNOWN FORM VAR: something unknown\n")
}
