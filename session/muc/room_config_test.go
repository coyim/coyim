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
			{Var: "FORM_TYPE", Type: RoomConfigFieldHidden, Values: []string{"stuff"}},
			{Var: ConfigFieldMaxHistoryFetch, Values: []string{"43"}, Options: []xmppData.FormFieldOptionX{
				{Value: "one"},
				{Value: "two"},
			}},
			{Var: ConfigFieldMaxHistoryLength, Values: []string{"43"}, Options: []xmppData.FormFieldOptionX{
				{Value: "one"},
				{Value: "two"},
			}},
			{Var: ConfigFieldAllowPM, Values: []string{"allow private messages"}},
			{Var: ConfigFieldAllowPrivateMessages, Values: []string{"allow private messages"}},
			{Var: ConfigFieldAllowMemberInvites, Values: []string{"true"}},
			{Var: ConfigFieldAllowInvites, Values: []string{"true"}},
			{Var: ConfigFieldCanChangeSubject, Values: []string{"true"}},
			{Var: ConfigFieldEnableArchiving, Values: []string{"true"}},
			{Var: ConfigFieldEnableLogging, Values: []string{"true"}},
			{Var: ConfigFieldMemberList, Values: []string{}},
			{Var: ConfigFieldLanguage, Values: []string{"eng"}},
			{Var: ConfigFieldPubsub, Values: []string{}},
			{Var: ConfigFieldMaxOccupantsNumber, Values: []string{"42"}},
			{Var: ConfigFieldMembersOnly, Values: []string{"true"}},
			{Var: ConfigFieldModerated, Values: []string{"true"}},
			{Var: ConfigFieldPasswordProtected, Values: []string{"true"}},
			{Var: ConfigFieldIsPersistent, Values: []string{"true"}},
			{Var: ConfigFieldPresenceBroadcast, Values: []string{}},
			{Var: ConfigFieldIsPublic, Values: []string{"true"}},
			{Var: ConfigFieldRoomAdmins, Values: []string{"one@foobar.com", "two@example.org"}},
			{Var: ConfigFieldRoomDescription, Values: []string{"a description"}},
			{Var: ConfigFieldRoomName, Values: []string{"a title"}},
			{Var: ConfigFieldOwners, Values: []string{}},
			{Var: ConfigFieldPassword, Values: []string{"a password"}},
			{Var: ConfigFieldWhoIs, Values: []string{"a whois"}},
			{Var: "unknown_field_name", Type: RoomConfigFieldText, Values: []string{"foo"}},
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
	c.Assert(res.Fields, HasLen, 2)

	vals := map[string][]string{}
	for _, ff := range res.Fields {
		vals[ff.Var] = ff.Values
	}

	c.Assert(vals, DeepEquals, map[string][]string{
		"FORM_TYPE":          {"http://jabber.org/protocol/muc#roomconfig"},
		"unknown_field_name": {"foo"},
	})
}

func (*MucRoomConfigSuite) Test_RoomConfigForm_setUnknowField(c *C) {
	cf := &RoomConfigForm{}
	fields := []HasRoomConfigFormField{}

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
		cf.setFieldX(fieldX)
		fields = append(fields, roomConfigFormFieldFactory(fieldX))
		c.Assert(cf.Fields, DeepEquals, fields)
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
		{
			"unknown_field_value",
			"foo",
			"field label",
			[]string{"bla"},
			[]string{"bla"},
		},
	}

	for _, chk := range checks {
		field := roomConfigFormFieldFactory(xmppData.FormFieldX{
			Var:    chk.name,
			Type:   chk.tp,
			Label:  chk.label,
			Values: chk.value,
		})

		c.Assert(field.Name(), Equals, chk.name)
		c.Assert(field.Type(), Equals, chk.tp)
		c.Assert(field.Label(), Equals, chk.label)
		c.Assert(field.Value(), DeepEquals, chk.expectedValue)
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

func (*MucRoomConfigSuite) Test_RoomConfigForm_updateFieldValueByName(c *C) {
	cf := &RoomConfigForm{}
	fields := []HasRoomConfigFormField{}

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
		cf.setFieldX(fieldX)
		fields = append(fields, roomConfigFormFieldFactory(fieldX))
	}

	cf.UpdateFieldValueByName("foo", "something")
	c.Assert(cf.Fields, DeepEquals, fields)

	for _, f := range fields {
		if f.Name() == "RoomConfigFieldText" {
			f.SetValue("bla1")
		}
	}

	cf.UpdateFieldValueByName("RoomConfigFieldText", "bla1")
	c.Assert(cf.Fields, DeepEquals, fields)

	for _, f := range fields {
		if f.Name() == "RoomConfigFieldText" {
			f.SetValue(nil)
		}
	}

	cf.UpdateFieldValueByName("RoomConfigFieldText", nil)
	c.Assert(cf.Fields, DeepEquals, fields)
}
