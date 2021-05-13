package muc

import (
	xmppData "github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
	. "gopkg.in/check.v1"
)

type MucRoomConfigSuite struct {
	rcf *RoomConfigForm
}

var _ = Suite(&MucRoomConfigSuite{})

func (s *MucRoomConfigSuite) SetUpSuite(c *C) {
	s.rcf = NewRoomConfigForm(&xmppData.Form{
		Fields: []xmppData.FormFieldX{
			{Var: configFieldFormType, Type: RoomConfigFieldHidden, Values: []string{"stuff"}},
			{Var: configFieldRoomName, Values: []string{"a title"}},
			{Var: configFieldRoomDescription, Values: []string{"a description"}},
			{Var: configFieldEnableLogging, Values: []string{"true"}},
			{Var: configFieldEnableArchiving, Values: []string{"true"}},
			{Var: configFieldMemberList, Values: []string{}},
			{Var: configFieldLanguage, Values: []string{"eng"}},
			{Var: configFieldPubsub, Values: []string{}},
			{Var: configFieldCanChangeSubject, Values: []string{"true"}},
			{Var: configFieldAllowInvites, Values: []string{"true"}},
			{Var: configFieldAllowMemberInvites, Values: []string{"true"}},
			{Var: configFieldAllowPM, Values: []string{"allow private messages"}},
			{Var: configFieldAllowPrivateMessages, Values: []string{"allow private messages"}},
			{Var: configFieldMaxOccupantsNumber, Values: []string{"42"}},
			{Var: configFieldIsPublic, Values: []string{"true"}},
			{Var: configFieldIsPersistent, Values: []string{"true"}},
			{Var: configFieldPresenceBroadcast, Values: []string{}},
			{Var: configFieldModerated, Values: []string{"true"}},
			{Var: configFieldMembersOnly, Values: []string{"true"}},
			{Var: configFieldPasswordProtected, Values: []string{"true"}},
			{Var: configFieldPassword, Values: []string{"a password"}},
			{Var: configFieldOwners, Values: []string{}},
			{Var: configFieldWhoIs, Values: []string{"a whois"}},
			{Var: configFieldMaxHistoryFetch, Values: []string{"43"}, Options: []xmppData.FormFieldOptionX{
				{Value: "one"},
				{Value: "two"},
			}},
			{Var: configFieldMaxHistoryLength, Values: []string{"43"}, Options: []xmppData.FormFieldOptionX{
				{Value: "one"},
				{Value: "two"},
			}},
			{Var: configFieldRoomAdmins, Values: []string{"one@foobar.com", "two@example.org"}},
			{Var: "unknown_field_name", Type: RoomConfigFieldText, Values: []string{"foo"}},
		},
	})
}

func (s *MucRoomConfigSuite) Test_NewRoomConfigForm(c *C) {
	c.Assert(s.rcf.Title, Equals, "a title")
	c.Assert(s.rcf.Description, Equals, "a description")
	c.Assert(s.rcf.Logged, Equals, true)
	c.Assert(s.rcf.OccupantsCanChangeSubject, Equals, true)
	c.Assert(s.rcf.OccupantsCanInvite, Equals, true)
	c.Assert(s.rcf.AllowPrivateMessages.Selected(), Equals, "allow private messages")
	c.Assert(s.rcf.MaxOccupantsNumber.Selected(), Equals, "42")
	c.Assert(s.rcf.Public, Equals, true)
	c.Assert(s.rcf.Persistent, Equals, true)
	c.Assert(s.rcf.Moderated, Equals, true)
	c.Assert(s.rcf.MembersOnly, Equals, true)
	c.Assert(s.rcf.PasswordProtected, Equals, true)
	c.Assert(s.rcf.Password, Equals, "a password")
	c.Assert(s.rcf.Whois.Selected(), Equals, "a whois")
	c.Assert(s.rcf.MaxHistoryFetch.Selected(), Equals, "43")
	c.Assert(s.rcf.Language, Equals, "eng")
	c.Assert(s.rcf.Admins.List(), DeepEquals, []jid.Any{jid.Parse("one@foobar.com"), jid.Parse("two@example.org")})

	res := s.rcf.GetFormData()

	c.Assert(res.Type, Equals, "submit")
	c.Assert(res.Fields, HasLen, 27)

	vals := map[string][]string{}
	for _, ff := range res.Fields {
		vals[ff.Var] = ff.Values
	}

	c.Assert(vals, DeepEquals, map[string][]string{
		"FORM_TYPE":                      {"stuff"},
		"muc#roomconfig_roomname":        {"a title"},
		"muc#roomconfig_roomdesc":        {"a description"},
		"muc#roomconfig_enablelogging":   {"true"},
		"muc#roomconfig_enablearchiving": {"true"},
		"muc#roomconfig_getmemberlist":   {},
		"muc#roomconfig_lang":            {"eng"},
		"muc#roomconfig_pubsub":          {""},
		"muc#roomconfig_changesubject":   {"true"},
		"muc#roomconfig_allowinvites":    {"true"},
		"{http://prosody.im/protocol/muc}roomconfig_allowmemberinvites": {"true"},
		"muc#roomconfig_allowpm":               {"allow private messages"},
		"allow_private_messages":               {"allow private messages"},
		"muc#roomconfig_maxusers":              {"42"},
		"muc#roomconfig_publicroom":            {"true"},
		"muc#roomconfig_persistentroom":        {"true"},
		"muc#roomconfig_presencebroadcast":     {},
		"muc#roomconfig_moderatedroom":         {"true"},
		"muc#roomconfig_membersonly":           {"true"},
		"muc#roomconfig_passwordprotectedroom": {"true"},
		"muc#roomconfig_roomsecret":            {"a password"},
		"muc#roomconfig_roomowners":            {},
		"muc#roomconfig_whois":                 {"a whois"},
		"muc#maxhistoryfetch":                  {"43"},
		"muc#roomconfig_historylength":         {"43"},
		"muc#roomconfig_roomadmins":            {"one@foobar.com", "two@example.org"},
		"unknown_field_name":                   {"foo"},
	})
}

func (*MucRoomConfigSuite) Test_RoomConfigForm_setUnknowField(c *C) {
	cf := &RoomConfigForm{}
	fields := []*RoomConfigFormField{}

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
			newRoomConfigFieldListValue([]string{"bla"}, []*RoomConfigFieldOption{{Value: "bla"}}),
		},
		{
			"RoomConfigFieldListMulti",
			RoomConfigFieldListMulti,
			"field label",
			[]string{"bla", "foo", "bla1", "foo1"},
			&RoomConfigFieldListMultiValue{value: []string{"bla", "foo", "bla1", "foo1"}},
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
		cf.setField(fieldX)
		fields = append(fields, newRoomConfigFormField(fieldX))
		c.Assert(cf.Fields, DeepEquals, fields)
	}
}

func (*MucRoomConfigSuite) Test_newRoomConfigFormField(c *C) {
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
			[]string{"bla foo"},
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
			"bla",
		},
		{
			"RoomConfigFieldListMulti",
			RoomConfigFieldListMulti,
			"field label",
			[]string{"bla", "foo", "bla1", "foo1"},
			[]string{"bla", "foo", "bla1", "foo1"},
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
		field := newRoomConfigFormField(xmppData.FormFieldX{
			Var:    chk.name,
			Type:   chk.tp,
			Label:  chk.label,
			Values: chk.value,
		})

		c.Assert(field.Name, Equals, chk.name)
		c.Assert(field.Type, Equals, chk.tp)
		c.Assert(field.Label, Equals, chk.label)
		c.Assert(field.RawValue(), DeepEquals, chk.expectedValue)
	}
}

func (s *MucRoomConfigSuite) Test_RoomConfigForm_getFieldDataValue(c *C) {
	fields := []struct {
		fieldName     string
		expectedValue []string
		ok            bool
	}{
		{configFieldFormType, []string{"stuff"}, true},
		{configFieldRoomName, []string{"a title"}, true},
		{configFieldRoomDescription, []string{"a description"}, true},
		{configFieldEnableLogging, []string{"true"}, true},
		{configFieldEnableArchiving, []string{"true"}, true},
		{configFieldMemberList, []string{}, true},
		{configFieldLanguage, []string{"eng"}, true},
		{configFieldPubsub, []string{""}, true},
		{configFieldCanChangeSubject, []string{"true"}, true},
		{configFieldAllowInvites, []string{"true"}, true},
		{configFieldAllowMemberInvites, []string{"true"}, true},
		{configFieldAllowPM, []string{"allow private messages"}, true},
		{configFieldAllowPrivateMessages, []string{"allow private messages"}, true},
		{configFieldMaxOccupantsNumber, []string{"42"}, true},
		{configFieldIsPublic, []string{"true"}, true},
		{configFieldIsPersistent, []string{"true"}, true},
		{configFieldPresenceBroadcast, []string{}, true},
		{configFieldModerated, []string{"true"}, true},
		{configFieldMembersOnly, []string{"true"}, true},
		{configFieldPasswordProtected, []string{"true"}, true},
		{configFieldPassword, []string{"a password"}, true},
		{configFieldOwners, []string{}, true},
		{configFieldWhoIs, []string{"a whois"}, true},
		{configFieldMaxHistoryFetch, []string{"43"}, true},
		{configFieldMaxHistoryLength, []string{"43"}, true},
		{configFieldRoomAdmins, []string{"one@foobar.com", "two@example.org"}, true},
		{"unknown_field_name", []string{"foo"}, true},
		{"another_unknown_field_name", nil, false},
	}

	for _, fieldCase := range fields {
		c.Logf("Checking case %s with expected value %v", fieldCase.fieldName, fieldCase.expectedValue)
		value, ok := s.rcf.getFieldDataValue(fieldCase.fieldName)
		c.Assert(value, DeepEquals, fieldCase.expectedValue)
		c.Assert(ok, Equals, fieldCase.ok)
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
	c.Assert(formFieldOptionsValues(nil), DeepEquals, []*RoomConfigFieldOption{})
	c.Assert(formFieldOptionsValues([]xmppData.FormFieldOptionX{}), DeepEquals, []*RoomConfigFieldOption{})
	c.Assert(formFieldOptionsValues([]xmppData.FormFieldOptionX{
		{Label: "bla", Value: "bla"},
		{Label: "whatever", Value: "foo"},
		{Label: "whatever2", Value: "bla2"},
		{Label: "whatever3", Value: "foo2"},
	}), DeepEquals, []*RoomConfigFieldOption{
		{Label: "bla", Value: "bla"},
		{Label: "whatever", Value: "foo"},
		{Label: "whatever2", Value: "bla2"},
		{Label: "whatever3", Value: "foo2"},
	})
}

func (s *MucRoomConfigSuite) Test_HasKnownField(c *C) {
	s.rcf.knownFields = map[RoomConfigFieldType]bool{
		RoomConfigFieldName:                 true,
		RoomConfigFieldDescription:          true,
		RoomConfigFieldAllowPrivateMessages: true,
	}

	c.Assert(s.rcf.HasKnownField(RoomConfigFieldName), Equals, true)
	c.Assert(s.rcf.HasKnownField(roomConfigFieldUnexpected), Equals, false)
}
