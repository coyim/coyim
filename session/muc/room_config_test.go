package muc

import (
	xmppData "github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
	. "gopkg.in/check.v1"
)

type MucRoomConfigSuite struct{}

var _ = Suite(&MucRoomConfigSuite{})

func (*MucRoomConfigSuite) Test_NewRoomConfigForm(c *C) {
	rcf := NewRoomConfigForm(&xmppData.Form{
		Fields: []xmppData.FormFieldX{
			{Var: "FORM_TYPE", Values: []string{"stuff"}},
			{Var: configFieldMaxHistoryFetch, Values: []string{"43"}, Options: []xmppData.FormFieldOptionX{
				{Value: "one"},
				{Value: "two"},
			}},
			{Var: configFieldMaxHistoryLength, Values: []string{"43"}, Options: []xmppData.FormFieldOptionX{
				{Value: "one"},
				{Value: "two"},
			}},
			{Var: configFieldAllowPM, Values: []string{"allow private messages"}},
			{Var: configFieldAllowPrivateMessages, Values: []string{"allow private messages"}},
			{Var: configFieldAllowMemberInvites, Values: []string{"true"}},
			{Var: configFieldAllowInvites, Values: []string{"true"}},
			{Var: configFieldCanChangeSubject, Values: []string{"true"}},
			{Var: configFieldEnableArchiving, Values: []string{"true"}},
			{Var: configFieldEnableLogging, Values: []string{"true"}},
			{Var: configFieldMemberList, Values: []string{}},
			{Var: configFieldLanguage, Values: []string{"eng"}},
			{Var: configFieldPubsub, Values: []string{}},
			{Var: configFieldMaxOccupantsNumber, Values: []string{"42"}},
			{Var: configFieldMembersOnly, Values: []string{"true"}},
			{Var: configFieldModerated, Values: []string{"true"}},
			{Var: configFieldPasswordProtected, Values: []string{"true"}},
			{Var: configFieldIsPersistent, Values: []string{"true"}},
			{Var: configFieldPresenceBroadcast, Values: []string{}},
			{Var: configFieldIsPublic, Values: []string{"true"}},
			{Var: configFieldRoomAdmins, Values: []string{"one@foobar.com", "two@example.org"}},
			{Var: configFieldRoomDescription, Values: []string{"a description"}},
			{Var: configFieldRoomName, Values: []string{"a title"}},
			{Var: configFieldOwners, Values: []string{}},
			{Var: configFieldPassword, Values: []string{"a password"}},
			{Var: configFieldWhoIs, Values: []string{"a whois"}},
		},
	})

	c.Assert(rcf.Title, Equals, "a title")
	c.Assert(rcf.Description, Equals, "a description")
	c.Assert(rcf.Logged, Equals, true)
	c.Assert(rcf.OccupantsCanChangeSubject, Equals, true)
	c.Assert(rcf.OccupantsCanInvite, Equals, true)
	c.Assert(rcf.AllowPrivateMessages.CurrentValue(), Equals, "allow private messages")
	c.Assert(rcf.MaxOccupantsNumber.CurrentValue(), Equals, "42")
	c.Assert(rcf.Public, Equals, true)
	c.Assert(rcf.Persistent, Equals, true)
	c.Assert(rcf.Moderated, Equals, true)
	c.Assert(rcf.MembersOnly, Equals, true)
	c.Assert(rcf.PasswordProtected, Equals, true)
	c.Assert(rcf.Password, Equals, "a password")
	c.Assert(rcf.Whois.CurrentValue(), Equals, "a whois")
	c.Assert(rcf.MaxHistoryFetch.CurrentValue(), Equals, "43")
	c.Assert(rcf.Language, Equals, "eng")
	c.Assert(rcf.Admins, DeepEquals, []jid.Any{jid.Parse("one@foobar.com"), jid.Parse("two@example.org")})

	res := rcf.GetFormData()

	c.Assert(res.Type, Equals, "submit")
	c.Assert(res.Fields, HasLen, 22)

	vals := map[string][]string{}
	for _, ff := range res.Fields {
		vals[ff.Var] = ff.Values
	}

	c.Assert(vals, DeepEquals, map[string][]string{
		"muc#roomconfig_enablelogging":                                  {"true"},
		"muc#roomconfig_enablearchiving":                                {"true"},
		"muc#roomconfig_changesubject":                                  {"true"},
		"muc#roomconfig_persistentroom":                                 {"true"},
		"muc#roomconfig_passwordprotectedroom":                          {"true"},
		"muc#roomconfig_lang":                                           {"eng"},
		"muc#roomconfig_roomname":                                       {"a title"},
		"muc#roomconfig_allowinvites":                                   {"true"},
		"{http://prosody.im/protocol/muc}roomconfig_allowmemberinvites": {"true"},
		"muc#roomconfig_allowpm":                                        {"allow private messages"},
		"allow_private_messages":                                        {"allow private messages"},
		"muc#roomconfig_roomsecret":                                     {"a password"},
		"muc#roomconfig_whois":                                          {"a whois"},
		"muc#roomconfig_publicroom":                                     {"true"},
		"muc#roomconfig_moderatedroom":                                  {"true"},
		"muc#maxhistoryfetch":                                           {"43"},
		"muc#roomconfig_historylength":                                  {"43"},
		"muc#roomconfig_roomadmins":                                     {"one@foobar.com", "two@example.org"},
		"FORM_TYPE":                                                     {"http://jabber.org/protocol/muc#roomconfig"},
		"muc#roomconfig_roomdesc":                                       {"a description"},
		"muc#roomconfig_maxusers":                                       {"42"},
		"muc#roomconfig_membersonly":                                    {"true"},
	})
}

func (*MucRoomConfigSuite) Test_RoomConfigForm_setUnknowField(c *C) {
	cf := &RoomConfigForm{}
	unknowFields := []*RoomConfigFormField{}

	checks := []struct {
		name          string
		tp            string
		label         string
		value         []string
		expectedValue interface{}
	}{
		{
			"RoomConfigFieldText",
			RoomConfigFieldText,
			"field label",
			[]string{"bla"},
			"bla",
		},
		{
			"RoomConfigFieldTextPrivate",
			RoomConfigFieldTextPrivate,
			"field label",
			[]string{"foo"},
			"foo",
		},
		{
			"RoomConfigFieldTextMulti",
			RoomConfigFieldTextMulti,
			"field label",
			[]string{"bla foo"},
			"bla foo",
		},
		{
			"RoomConfigFieldBoolean",
			RoomConfigFieldBoolean,
			"field label",
			[]string{"true"},
			true,
		},
		{
			"RoomConfigFieldList",
			RoomConfigFieldList,
			"field label",
			[]string{"bla"},
			&configListSingleField{value: "bla"},
		},
		{
			"RoomConfigFieldListMulti",
			RoomConfigFieldListMulti,
			"field label",
			[]string{"bla", "foo", "bla1", "foo1"},
			&configListMultiField{values: []string{"bla", "foo", "bla1", "foo1"}},
		},
		{
			"RoomConfigFieldJidMulti",
			RoomConfigFieldJidMulti,
			"field label",
			[]string{"bla", "foo", "bla@domain.org", "foo@domain.org"},
			[]jid.Any{jid.Parse("bla"), jid.Parse("foo"), jid.Parse("bla@domain.org"), jid.Parse("foo@domain.org")},
		},
	}

	for _, chk := range checks {
		fieldX := xmppData.FormFieldX{
			Var:    chk.name,
			Type:   chk.tp,
			Label:  chk.label,
			Values: chk.value,
		}
		cf.setUnknowField(fieldX)
		unknowFields = append(unknowFields, roomConfigFormFieldFactory(fieldX))
		c.Assert(cf.UnknowFields, DeepEquals, unknowFields)
	}
}

func (*MucRoomConfigSuite) Test_roomConfigFormFieldFactory(c *C) {
	checks := []struct {
		name          string
		tp            string
		label         string
		value         []string
		expectedValue interface{}
	}{
		{
			"RoomConfigFieldText",
			RoomConfigFieldText,
			"field label",
			[]string{"bla"},
			"bla",
		},
		{
			"RoomConfigFieldTextPrivate",
			RoomConfigFieldTextPrivate,
			"field label",
			[]string{"foo"},
			"foo",
		},
		{
			"RoomConfigFieldTextMulti",
			RoomConfigFieldTextMulti,
			"field label",
			[]string{"bla foo"},
			"bla foo",
		},
		{
			"RoomConfigFieldBoolean",
			RoomConfigFieldBoolean,
			"field label",
			[]string{"true"},
			true,
		},
		{
			"RoomConfigFieldList",
			RoomConfigFieldList,
			"field label",
			[]string{"bla"},
			&configListSingleField{value: "bla"},
		},
		{
			"RoomConfigFieldListMulti",
			RoomConfigFieldListMulti,
			"field label",
			[]string{"bla", "foo", "bla1", "foo1"},
			&configListMultiField{values: []string{"bla", "foo", "bla1", "foo1"}},
		},
		{
			"RoomConfigFieldJidMulti",
			RoomConfigFieldJidMulti,
			"field label",
			[]string{"bla", "foo", "bla@domain.org", "foo@domain.org"},
			[]jid.Any{jid.Parse("bla"), jid.Parse("foo"), jid.Parse("bla@domain.org"), jid.Parse("foo@domain.org")},
		},
	}

	for _, chk := range checks {
		field := roomConfigFormFieldFactory(xmppData.FormFieldX{
			Var:    chk.name,
			Type:   chk.tp,
			Label:  chk.label,
			Values: chk.value,
		})

		c.Assert(field.Name, Equals, chk.name)
		c.Assert(field.Type, Equals, chk.tp)
		c.Assert(field.Label, Equals, chk.label)
		c.Assert(field.Value, DeepEquals, chk.expectedValue)
	}
}

func (*MucRoomConfigSuite) Test_formFieldBool(c *C) {
	c.Assert(formFieldBool(nil), Equals, false)
	c.Assert(formFieldBool([]string(nil)), Equals, false)
	c.Assert(formFieldBool([]string{"true"}), Equals, true)
	c.Assert(formFieldBool([]string{"false"}), Equals, false)
	c.Assert(formFieldBool([]string{"true", "false"}), Equals, true)
}

func (*MucRoomConfigSuite) Test_formFieldSingleString(c *C) {
	c.Assert(formFieldSingleString(nil), Equals, "")
	c.Assert(formFieldSingleString([]string(nil)), Equals, "")
	c.Assert(formFieldSingleString([]string{"bla"}), Equals, "bla")
	c.Assert(formFieldSingleString([]string{"bla", "foo"}), Equals, "bla")
	c.Assert(formFieldSingleString([]string{""}), Equals, "")
}

func (*MucRoomConfigSuite) Test_formFieldOptionsValues(c *C) {
	c.Assert(formFieldOptionsValues(nil), IsNil)
	c.Assert(formFieldOptionsValues([]xmppData.FormFieldOptionX{}), DeepEquals, []string(nil))
	c.Assert(formFieldOptionsValues([]xmppData.FormFieldOptionX{
		{Label: "bla", Value: "bla"},
		{Label: "whatever", Value: "foo"},
		{Label: "whatever2", Value: "bla2"},
		{Label: "whatever3", Value: "foo2"},
	}), DeepEquals, []string{"bla", "foo", "bla2", "foo2"})
}

func (*MucRoomConfigSuite) Test_formFieldJidList(c *C) {
	c.Assert(formFieldJidList(nil), IsNil)
	c.Assert(formFieldJidList([]string{}), DeepEquals, []jid.Any(nil))
	c.Assert(formFieldJidList([]string{"bla"}), DeepEquals, []jid.Any{jid.Parse("bla")})
	c.Assert(formFieldJidList([]string{"bla", "foo@domain.org"}), DeepEquals, []jid.Any{jid.Parse("bla"), jid.Parse("foo@domain.org")})
}

func (*MucRoomConfigSuite) Test_jidListToStringList(c *C) {
	c.Assert(jidListToStringList(nil), IsNil)
	c.Assert(jidListToStringList([]jid.Any{}), DeepEquals, []string(nil))
	c.Assert(jidListToStringList([]jid.Any{jid.Parse("bla")}), DeepEquals, []string{"bla"})
	c.Assert(jidListToStringList([]jid.Any{jid.Parse("foo@domain.org"), jid.Parse("foo"), jid.Parse("bla@domain.org")}), DeepEquals, []string{"foo@domain.org", "foo", "bla@domain.org"})
}
