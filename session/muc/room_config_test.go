package muc

import (
	xmppData "github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
	. "gopkg.in/check.v1"
)

type MucRoomConfigSuite struct{}

var _ = Suite(&MucRoomConfigSuite{})

func (*MucRoomConfigSuite) Test_NewRoomConfigRom(c *C) {
	rcf := NewRoomConfigRom(&xmppData.Form{
		Fields: []xmppData.FormFieldX{
			xmppData.FormFieldX{Var: "FORM_TYPE", Values: []string{"stuff"}},
			xmppData.FormFieldX{Var: configFieldMaxHistoryFetch, Values: []string{"43"}, Options: []xmppData.FormFieldOptionX{
				xmppData.FormFieldOptionX{Value: "one"},
				xmppData.FormFieldOptionX{Value: "two"},
			}},
			xmppData.FormFieldX{Var: configFieldAllowPrivateMessages, Values: []string{"allow private messages"}},
			xmppData.FormFieldX{Var: configFieldAllowInvites, Values: []string{"true"}},
			xmppData.FormFieldX{Var: configFieldCanChangeSubject, Values: []string{"true"}},
			xmppData.FormFieldX{Var: configFieldEnableLogging, Values: []string{"true"}},
			xmppData.FormFieldX{Var: configFieldMemberList, Values: []string{}},
			xmppData.FormFieldX{Var: configFieldLanguage, Values: []string{"eng"}},
			xmppData.FormFieldX{Var: configFieldPubsub, Values: []string{}},
			xmppData.FormFieldX{Var: configFieldMaxOccupantsNumber, Values: []string{"42"}},
			xmppData.FormFieldX{Var: configFieldMembersOnly, Values: []string{"true"}},
			xmppData.FormFieldX{Var: configFieldModerated, Values: []string{"true"}},
			xmppData.FormFieldX{Var: configFieldPasswordProtected, Values: []string{"true"}},
			xmppData.FormFieldX{Var: configFieldIsPersistent, Values: []string{"true"}},
			xmppData.FormFieldX{Var: configFieldPresenceBroadcast, Values: []string{}},
			xmppData.FormFieldX{Var: configFieldIsPublic, Values: []string{"true"}},
			xmppData.FormFieldX{Var: configFieldRoomAdmins, Values: []string{"one@foobar.com", "two@example.org"}},
			xmppData.FormFieldX{Var: configFieldRoomDescription, Values: []string{"a description"}},
			xmppData.FormFieldX{Var: configFieldRoomName, Values: []string{"a title"}},
			xmppData.FormFieldX{Var: configFieldOwners, Values: []string{}},
			xmppData.FormFieldX{Var: configFieldPassword, Values: []string{"a password"}},
			xmppData.FormFieldX{Var: configFieldWhoIs, Values: []string{"a whois"}},
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
	c.Assert(res.Fields, HasLen, 18)

	vals := map[string][]string{}
	for _, ff := range res.Fields {
		vals[ff.Var] = ff.Values
	}

	c.Assert(vals, DeepEquals, map[string][]string{
		"muc#roomconfig_enablelogging":         []string{"true"},
		"muc#roomconfig_changesubject":         []string{"true"},
		"muc#roomconfig_persistentroom":        []string{"true"},
		"muc#roomconfig_passwordprotectedroom": []string{"true"},
		"muc#roomconfig_lang":                  []string{"eng"},
		"muc#roomconfig_roomname":              []string{"a title"},
		"muc#roomconfig_allowinvites":          []string{"true"},
		"muc#roomconfig_allowpm":               []string{"allow private messages"},
		"muc#roomconfig_roomsecret":            []string{"a password"},
		"muc#roomconfig_whois":                 []string{"a whois"},
		"muc#roomconfig_publicroom":            []string{"true"},
		"muc#roomconfig_moderatedroom":         []string{"true"},
		"muc#maxhistoryfetch":                  []string{"43"},
		"muc#roomconfig_roomadmins":            []string{"one@foobar.com", "two@example.org"},
		"FORM_TYPE":                            []string{"http://jabber.org/protocol/muc#roomconfig"},
		"muc#roomconfig_roomdesc":              []string{"a description"},
		"muc#roomconfig_maxusers":              []string{"42"},
		"muc#roomconfig_membersonly":           []string{"true"},
	})
}
