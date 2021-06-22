package muc

import (
	"github.com/coyim/coyim/session/muc/data"
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
			{
				Var:    configFieldFormType,
				Type:   RoomConfigFieldHidden,
				Values: []string{"stuff"},
			},
			{
				Var:    configFieldRoomName,
				Type:   RoomConfigFieldText,
				Values: []string{"a title"},
			},
			{
				Var:    configFieldRoomDescription,
				Type:   RoomConfigFieldTextMulti,
				Values: []string{"a description"},
			},
			{
				Var:    configFieldEnableLogging,
				Type:   RoomConfigFieldBoolean,
				Values: []string{"true"},
			},
			{
				Var:    configFieldMemberList,
				Type:   RoomConfigFieldJidMulti,
				Values: []string{},
			},
			{
				Var:    configFieldLanguage,
				Type:   RoomConfigFieldText,
				Values: []string{"eng"},
			},
			{
				Var:    configFieldPubsub,
				Type:   RoomConfigFieldText,
				Values: []string{},
			},
			{
				Var:    configFieldCanChangeSubject,
				Type:   RoomConfigFieldBoolean,
				Values: []string{"true"},
			},
			{
				Var:    configFieldAllowMemberInvites,
				Type:   RoomConfigFieldBoolean,
				Values: []string{"true"},
			},
			{
				Var:    configFieldAllowPrivateMessages,
				Type:   RoomConfigFieldBoolean,
				Values: []string{"false"},
			},
			{
				Var:    configFieldMaxOccupantsNumber,
				Type:   RoomConfigFieldList,
				Values: []string{"42"},
			},
			{
				Var:    configFieldIsPublic,
				Type:   RoomConfigFieldBoolean,
				Values: []string{"true"},
			},
			{
				Var:    configFieldIsPersistent,
				Type:   RoomConfigFieldBoolean,
				Values: []string{"true"},
			},
			{
				Var:    configFieldPresenceBroadcast,
				Type:   RoomConfigFieldListMulti,
				Values: []string{},
			},
			{
				Var:    configFieldModerated,
				Type:   RoomConfigFieldBoolean,
				Values: []string{"true"},
			},
			{
				Var:    configFieldMembersOnly,
				Type:   RoomConfigFieldBoolean,
				Values: []string{"true"},
			},
			{
				Var:    configFieldPasswordProtected,
				Type:   RoomConfigFieldBoolean,
				Values: []string{"true"},
			},
			{
				Var:    configFieldPassword,
				Type:   RoomConfigFieldText,
				Values: []string{"a password"},
			},
			{
				Var:    configFieldOwners,
				Type:   RoomConfigFieldJidMulti,
				Values: []string{},
			},
			{
				Var:    configFieldWhoIs,
				Type:   RoomConfigFieldList,
				Values: []string{"a whois"},
			},
			{
				Var:    configFieldMaxHistoryLength,
				Type:   RoomConfigFieldList,
				Values: []string{"43"},
				Options: []xmppData.FormFieldOptionX{
					{Value: "one"},
					{Value: "two"},
				},
			},
			{
				Var:    configFieldRoomAdmins,
				Type:   RoomConfigFieldJidMulti,
				Values: []string{"one@foobar.com", "two@example.org"},
			},
			{
				Var:    "unknown_field_name",
				Type:   RoomConfigFieldText,
				Values: []string{"foo"},
			},
		},
	})
}

func (s *MucRoomConfigSuite) stringValueFromFormField(rcft RoomConfigFieldType) (v string) {
	f, ok := s.rcf.GetKnownField(rcft)
	if ok && len(f.Value()) > 0 {
		v = f.Value()[0]
	}
	return
}

func (s *MucRoomConfigSuite) Test_NewRoomConfigForm(c *C) {
	c.Assert(s.stringValueFromFormField(RoomConfigFieldDescription), Equals, "a description")
	c.Assert(s.stringValueFromFormField(RoomConfigFieldEnableLogging), Equals, "true")
	c.Assert(s.stringValueFromFormField(RoomConfigFieldCanChangeSubject), Equals, "true")
	c.Assert(s.stringValueFromFormField(RoomConfigFieldAllowInvites), Equals, "true")
	c.Assert(s.stringValueFromFormField(RoomConfigFieldAllowPrivateMessages), Equals, "false")
	c.Assert(s.stringValueFromFormField(RoomConfigFieldMaxOccupantsNumber), Equals, "42")
	c.Assert(s.stringValueFromFormField(RoomConfigFieldIsPublic), Equals, "true")
	c.Assert(s.stringValueFromFormField(RoomConfigFieldIsPersistent), Equals, "true")
	c.Assert(s.stringValueFromFormField(RoomConfigFieldIsModerated), Equals, "true")
	c.Assert(s.stringValueFromFormField(RoomConfigFieldIsMembersOnly), Equals, "true")
	c.Assert(s.stringValueFromFormField(RoomConfigFieldIsPasswordProtected), Equals, "true")
	c.Assert(s.stringValueFromFormField(RoomConfigFieldPassword), Equals, "a password")
	c.Assert(s.stringValueFromFormField(RoomConfigFieldWhoIs), Equals, "a whois")
	c.Assert(s.stringValueFromFormField(RoomConfigFieldMaxHistoryFetch), Equals, "43")
	c.Assert(s.stringValueFromFormField(RoomConfigFieldLanguage), Equals, "eng")

	res := s.rcf.GetFormData()

	c.Assert(res.Type, Equals, "submit")
	c.Assert(res.Fields, HasLen, 23)

	vals := map[string][]string{}
	for _, ff := range res.Fields {
		vals[ff.Var] = ff.Values
	}

	c.Assert(vals, DeepEquals, map[string][]string{
		"FORM_TYPE":                    {"http://jabber.org/protocol/muc#roomconfig"},
		"muc#roomconfig_roomname":      {"a title"},
		"muc#roomconfig_roomdesc":      {"a description"},
		"muc#roomconfig_enablelogging": {"true"},
		"muc#roomconfig_getmemberlist": {},
		"muc#roomconfig_lang":          {"eng"},
		"muc#roomconfig_pubsub":        {""},
		"muc#roomconfig_changesubject": {"true"},
		"{http://prosody.im/protocol/muc}roomconfig_allowmemberinvites": {"true"},
		"allow_private_messages":               {"false"},
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
		"muc#roomconfig_historylength":         {"43"},
		"muc#roomconfig_roomadmins":            {"one@foobar.com", "two@example.org"},
		"unknown_field_name":                   {"foo"},
	})
}

func (*MucRoomConfigSuite) Test_RoomConfigForm_GetUnknownFields(c *C) {
	cf := NewRoomConfigForm(&xmppData.Form{
		Fields: []xmppData.FormFieldX{
			{
				Label:  "field label",
				Var:    "foo",
				Type:   RoomConfigFieldHidden,
				Values: []string{"stuff"},
			},
			{
				Label:  "field label",
				Var:    "bla",
				Type:   RoomConfigFieldText,
				Values: []string{"a title"},
			},
			{
				Label:  "field label",
				Var:    "stuff",
				Type:   RoomConfigFieldTextMulti,
				Values: []string{"a description"},
			},
		},
	})

	checks := []struct {
		name          string
		tp            string
		label         string
		value         []string
		expectedValue interface{}
	}{
		{
			"foo",
			RoomConfigFieldHidden,
			"field label",
			[]string{"stuff"},
			"bla",
		},
		{
			"bla",
			RoomConfigFieldText,
			"field label",
			[]string{"a title"},
			"bla",
		},
		{
			"stuff",
			RoomConfigFieldTextMulti,
			"field label",
			[]string{"a description"},
			"foo",
		},
	}

	fields := []*RoomConfigFormField{}
	for _, chk := range checks {
		fieldX := xmppData.FormFieldX{
			Var:    chk.name,
			Type:   chk.tp,
			Label:  chk.label,
			Values: chk.value,
		}
		fields = append(fields, newRoomConfigFormField(fieldX))
	}
	c.Assert(cf.GetUnknownFields(), DeepEquals, fields)
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
		{Var: "bla", Value: "bla"},
		{Var: "whatever", Value: "foo"},
		{Var: "whatever2", Value: "bla2"},
		{Var: "whatever3", Value: "foo2"},
	}), DeepEquals, []*RoomConfigFieldOption{
		{Label: "bla", Value: "bla"},
		{Label: "whatever", Value: "foo"},
		{Label: "whatever2", Value: "bla2"},
		{Label: "whatever3", Value: "foo2"},
	})
}

func (s *MucRoomConfigSuite) Test_RoomConfigForm_HasKnownField(c *C) {
	c.Assert(s.rcf.HasKnownField(RoomConfigFieldName), Equals, true)
	c.Assert(s.rcf.HasKnownField(RoomConfigFieldDescription), Equals, true)
	c.Assert(s.rcf.HasKnownField(RoomConfigFieldUnexpected), Equals, false)
}

func (s *MucRoomConfigSuite) Test_RoomConfigForm_GetKnownField(c *C) {
	var field *RoomConfigFormField
	var ok bool

	field, ok = s.rcf.GetKnownField(RoomConfigFieldName)
	c.Assert(field, DeepEquals, &RoomConfigFormField{
		Name:        "muc#roomconfig_roomname",
		Label:       "",
		Type:        "text-single",
		Description: "",
		value:       &RoomConfigFieldTextValue{"a title"},
	})
	c.Assert(ok, Equals, true)

	field, ok = s.rcf.GetKnownField(RoomConfigFieldDescription)
	c.Assert(field, DeepEquals, &RoomConfigFormField{
		Name:        "muc#roomconfig_roomdesc",
		Label:       "",
		Type:        "text-multi",
		Description: "",
		value:       &RoomConfigFieldTextMultiValue{[]string{"a description"}},
	})
	c.Assert(ok, Equals, true)

	field, ok = s.rcf.GetKnownField(RoomConfigFieldUnexpected)
	c.Assert(field, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *MucRoomConfigSuite) Test_RoomConfigForm_getKnownFieldValue(c *C) {
	value, ok := s.rcf.getKnownFieldValue(configFieldRoomName)
	c.Assert(value, DeepEquals, []string{"a title"})
	c.Assert(ok, Equals, true)

	value, ok = s.rcf.getKnownFieldValue("foo")
	c.Assert(value, IsNil)
	c.Assert(ok, Equals, false)

}

func (s *MucRoomConfigSuite) Test_RoomConfigForm_GetOccupantsByAffiliation(c *C) {
	ownerAffiliation := &data.OwnerAffiliation{}
	s.rcf.UpdateRoomOccupantsByAffiliation(ownerAffiliation, []*RoomOccupantItem{
		{
			Jid:         jid.Parse("jid"),
			Affiliation: ownerAffiliation,
			Reason:      "bla",
		},
		{
			Jid:         jid.Parse("jid@foo.org"),
			Affiliation: ownerAffiliation,
			Reason:      "foo",
		},
		{
			Jid:         jid.Parse("1234"),
			Affiliation: ownerAffiliation,
			Reason:      "foo123",
		},
	})

	adminAffiliation := &data.AdminAffiliation{}
	s.rcf.UpdateRoomOccupantsByAffiliation(adminAffiliation, []*RoomOccupantItem{
		{
			Jid:         jid.Parse("batman@cave.org"),
			Affiliation: adminAffiliation,
			Reason:      "boom",
		},
	})

	outcastAffiliation := &data.OutcastAffiliation{}
	s.rcf.UpdateRoomOccupantsByAffiliation(outcastAffiliation, []*RoomOccupantItem{
		{
			Jid:         jid.Parse("robin"),
			Affiliation: outcastAffiliation,
			Reason:      "123456",
		},
	})

	c.Assert(s.rcf.GetOccupantsByAffiliation(ownerAffiliation), HasLen, 3)
	c.Assert(s.rcf.GetOccupantsByAffiliation(adminAffiliation), HasLen, 1)
	c.Assert(s.rcf.GetOccupantsByAffiliation(outcastAffiliation), HasLen, 1)
}

func (s *MucRoomConfigSuite) Test_RoomConfigForm_ConfigureRoomAsPersistent(c *C) {
	f := s.rcf.knownFields[RoomConfigFieldIsPersistent].value.(*RoomConfigFieldBooleanValue)
	f.SetBoolean(false)
	c.Assert(f.value, Equals, false)
	s.rcf.ConfigureRoomAsPersistent()
	c.Assert(f.value, Equals, true)
}

func (s *MucRoomConfigSuite) Test_RoomConfigForm_GetConfiguredPassword(c *C) {
	c.Assert(s.rcf.GetConfiguredPassword(), Equals, "a password")
}

func (s *MucRoomConfigSuite) Test_RoomConfigForm_GetRoomOccupantsToUpdate(c *C) {
	roomOccupants := s.rcf.GetRoomOccupantsToUpdate()
	c.Assert(roomOccupants, HasLen, 5)

	owners := []*RoomOccupantItem{}
	ownerAffiliation := &data.OwnerAffiliation{}
	for _, o := range roomOccupants {
		if o.Affiliation == ownerAffiliation {
			owners = append(owners, o)
		}
	}
	c.Assert(owners, DeepEquals, []*RoomOccupantItem{
		{
			Jid:         jid.Parse("jid"),
			Affiliation: ownerAffiliation,
			Reason:      "bla",
		},
		{
			Jid:         jid.Parse("jid@foo.org"),
			Affiliation: ownerAffiliation,
			Reason:      "foo",
		},
		{
			Jid:         jid.Parse("1234"),
			Affiliation: ownerAffiliation,
			Reason:      "foo123",
		},
	})

	admins := []*RoomOccupantItem{}
	adminAffiliation := &data.AdminAffiliation{}
	for _, o := range roomOccupants {
		if o.Affiliation == adminAffiliation {
			admins = append(admins, o)
		}
	}
	c.Assert(admins, DeepEquals, []*RoomOccupantItem{
		{
			Jid:         jid.Parse("batman@cave.org"),
			Affiliation: adminAffiliation,
			Reason:      "boom",
		},
	})

	banned := []*RoomOccupantItem{}
	outcastAffiliation := &data.OutcastAffiliation{}
	for _, o := range roomOccupants {
		if o.Affiliation == outcastAffiliation {
			banned = append(banned, o)
		}
	}
	c.Assert(banned, DeepEquals, []*RoomOccupantItem{
		{
			Jid:         jid.Parse("robin"),
			Affiliation: outcastAffiliation,
			Reason:      "123456",
		},
	})
}

func (s *MucRoomConfigSuite) Test_RoomConfigForm_updateValueOfPasswordProtectedField(c *C) {
	field, _ := s.rcf.GetKnownField(RoomConfigFieldIsPasswordProtected)
	c.Assert(field.Value(), DeepEquals, []string{"true"})

	f := s.rcf.knownFields[RoomConfigFieldPassword].value.(*RoomConfigFieldTextValue)
	f.SetText("")

	s.rcf.updateValueOfPasswordProtectedField()
	c.Assert(field.Value(), DeepEquals, []string{"false"})

}
