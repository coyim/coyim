package muc

import (
	"strconv"
	"strings"

	xmppData "github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
)

const (
	// ConfigFieldFormType represents the configuration form type field
	ConfigFieldFormType = "http://jabber.org/protocol/muc#roomconfig"
	// ConfigFieldRoomName represents the room name form field
	ConfigFieldRoomName = "muc#roomconfig_roomname"
	// ConfigFieldRoomDescription represents the room description form field
	ConfigFieldRoomDescription = "muc#roomconfig_roomdesc"
	// ConfigFieldEnableLogging represents the enable logging form field
	ConfigFieldEnableLogging = "muc#roomconfig_enablelogging"
	// ConfigFieldEnableArchiving represents the enable archiving form field
	ConfigFieldEnableArchiving = "muc#roomconfig_enablearchiving"
	// ConfigFieldMemberList represents the get member list form field
	ConfigFieldMemberList = "muc#roomconfig_getmemberlist"
	// ConfigFieldLanguage represents the room language form field
	ConfigFieldLanguage = "muc#roomconfig_lang"
	// ConfigFieldPubsub represents the pubsub form field
	ConfigFieldPubsub = "muc#roomconfig_pubsub"
	// ConfigFieldCanChangeSubject represents the change subject form field
	ConfigFieldCanChangeSubject = "muc#roomconfig_changesubject"
	// ConfigFieldAllowInvites represents the allow invites form field
	ConfigFieldAllowInvites = "muc#roomconfig_allowinvites"
	// ConfigFieldAllowMemberInvites represents the allow member invites form field (for some services)
	ConfigFieldAllowMemberInvites = "{http://prosody.im/protocol/muc}roomconfig_allowmemberinvites"
	// ConfigFieldAllowPM represents the allow private messages form field
	ConfigFieldAllowPM = "muc#roomconfig_allowpm"
	// ConfigFieldAllowPrivateMessages represent the allow private messages form fields (for some services)
	ConfigFieldAllowPrivateMessages = "allow_private_messages"
	// ConfigFieldMaxOccupantsNumber represents the max users form field
	ConfigFieldMaxOccupantsNumber = "muc#roomconfig_maxusers"
	// ConfigFieldIsPublic represents the public room form field
	ConfigFieldIsPublic = "muc#roomconfig_publicroom"
	// ConfigFieldIsPersistent represents the persistent room form field
	ConfigFieldIsPersistent = "muc#roomconfig_persistentroom"
	// ConfigFieldPresenceBroadcast represents the presence broadcast form field
	ConfigFieldPresenceBroadcast = "muc#roomconfig_presencebroadcast"
	// ConfigFieldModerated represents the moderated room form field
	ConfigFieldModerated = "muc#roomconfig_moderatedroom"
	// ConfigFieldMembersOnly represents the members only form field (for some services)
	ConfigFieldMembersOnly = "muc#roomconfig_membersonly"
	// ConfigFieldPasswordProtected represents the password protected room form field
	ConfigFieldPasswordProtected = "muc#roomconfig_passwordprotectedroom"
	// ConfigFieldPassword represents the room secret form field (for some services)
	ConfigFieldPassword = "muc#roomconfig_roomsecret"
	// ConfigFieldOwners represents the room owners list form field
	ConfigFieldOwners = "muc#roomconfig_roomowners"
	// ConfigFieldWhoIs represents the whois form field
	ConfigFieldWhoIs = "muc#roomconfig_whois"
	// ConfigFieldMaxHistoryFetch represents the max history fetch form field
	ConfigFieldMaxHistoryFetch = "muc#maxhistoryfetch"
	// ConfigFieldMaxHistoryLength represents the max history length form field (for some services)
	ConfigFieldMaxHistoryLength = "muc#roomconfig_historylength"
	// ConfigFieldRoomAdmins represents the room admins list form field
	ConfigFieldRoomAdmins = "muc#roomconfig_roomadmins"
)

// RoomConfigForm represents a room configuration form
type RoomConfigForm struct {
	MaxHistoryFetch                ConfigListSingleField
	AllowPrivateMessages           ConfigListSingleField
	OccupantsCanInvite             bool
	OccupantsCanChangeSubject      bool
	Logged                         bool
	RetrieveMembersList            ConfigListMultiField
	AssociatedPublishSubscribeNode string
	MaxOccupantsNumber             ConfigListSingleField
	MembersOnly                    bool
	Moderated                      bool
	PasswordProtected              bool
	Persistent                     bool
	PresenceBroadcast              ConfigListMultiField
	Public                         bool
	Admins                         []jid.Any
	Owners                         []jid.Any
	Password                       string
	Whois                          ConfigListSingleField

	Fields map[string]HasRoomConfigFormField
}

// NewRoomConfigForm creates a new room configuration form instance
func NewRoomConfigForm(form *xmppData.Form) *RoomConfigForm {
	cf := &RoomConfigForm{
		Fields: make(map[string]HasRoomConfigFormField),
	}

	cf.MaxHistoryFetch = newConfigListSingleField([]string{
		RoomConfigOption10,
		RoomConfigOption20,
		RoomConfigOption30,
		RoomConfigOption50,
		RoomConfigOption100,
		RoomConfigOptionNone,
	})

	cf.AllowPrivateMessages = newConfigListSingleField([]string{
		RoomConfigOptionParticipant,
		RoomConfigOptionModerators,
		RoomConfigOptionNone,
	})

	cf.RetrieveMembersList = newConfigListMultiField([]string{
		RoomConfigOptionModerator,
		RoomConfigOptionParticipant,
		RoomConfigOptionVisitor,
	})

	cf.MaxOccupantsNumber = newConfigListSingleField([]string{
		RoomConfigOption10,
		RoomConfigOption20,
		RoomConfigOption30,
		RoomConfigOption50,
		RoomConfigOption100,
		RoomConfigOptionNone,
	})

	cf.PresenceBroadcast = newConfigListMultiField([]string{
		RoomConfigOptionModerator,
		RoomConfigOptionParticipant,
		RoomConfigOptionVisitor,
	})

	cf.Whois = newConfigListSingleField([]string{
		RoomConfigOptionModerators,
		RoomConfigOptionAnyone,
	})

	cf.setFormFields(form.Fields)

	return cf
}

// GetFormData returns a representation of the room config FORM_TYPE as described in the
// XMPP specification for MUC
//
// For more information see:
// https://xmpp.org/extensions/xep-0045.html#createroom-reserved
// https://xmpp.org/extensions/xep-0045.html#example-163
func (rcf *RoomConfigForm) GetFormData() *xmppData.Form {
	fields := map[string][]string{
		"FORM_TYPE":                     {ConfigFieldFormType},
		ConfigFieldEnableLogging:        {strconv.FormatBool(rcf.Logged)},
		ConfigFieldEnableArchiving:      {strconv.FormatBool(rcf.Logged)},
		ConfigFieldCanChangeSubject:     {strconv.FormatBool(rcf.OccupantsCanChangeSubject)},
		ConfigFieldAllowInvites:         {strconv.FormatBool(rcf.OccupantsCanInvite)},
		ConfigFieldAllowMemberInvites:   {strconv.FormatBool(rcf.OccupantsCanInvite)},
		ConfigFieldAllowPM:              {rcf.AllowPrivateMessages.CurrentValue()},
		ConfigFieldAllowPrivateMessages: {rcf.AllowPrivateMessages.CurrentValue()},
		ConfigFieldMaxOccupantsNumber:   {rcf.MaxOccupantsNumber.CurrentValue()},
		ConfigFieldIsPublic:             {strconv.FormatBool(rcf.Public)},
		ConfigFieldIsPersistent:         {strconv.FormatBool(rcf.Persistent)},
		ConfigFieldModerated:            {strconv.FormatBool(rcf.Moderated)},
		ConfigFieldMembersOnly:          {strconv.FormatBool(rcf.MembersOnly)},
		ConfigFieldPasswordProtected:    {strconv.FormatBool(rcf.PasswordProtected)},
		ConfigFieldPassword:             {rcf.Password},
		ConfigFieldWhoIs:                {rcf.Whois.CurrentValue()},
		ConfigFieldMaxHistoryFetch:      {rcf.MaxHistoryFetch.CurrentValue()},
		ConfigFieldMaxHistoryLength:     {rcf.MaxHistoryFetch.CurrentValue()},
		ConfigFieldRoomAdmins:           jidListToStringList(rcf.Admins),
	}

	formFields := []xmppData.FormFieldX{}
	for name, values := range fields {
		formFields = append(formFields, xmppData.FormFieldX{
			Var:    name,
			Values: values,
		})
	}

	for _, f := range rcf.Fields {
		formFields = append(formFields, xmppData.FormFieldX{
			Var:    f.Name(),
			Values: f.ValueX(),
		})
	}

	return &xmppData.Form{
		Type:   "submit",
		Fields: formFields,
	}
}

func (rcf *RoomConfigForm) setFormFields(fields []xmppData.FormFieldX) {
	for _, field := range fields {
		rcf.fieldNames = append(rcf.fieldNames, field.Var)
		rcf.setField(field)
	}
}

func (rcf *RoomConfigForm) setField(field xmppData.FormFieldX) {
	switch field.Var {
	case ConfigFieldMaxHistoryFetch, ConfigFieldMaxHistoryLength:
		rcf.MaxHistoryFetch.UpdateField(formFieldSingleString(field.Values), formFieldOptionsValues(field.Options))

	case ConfigFieldAllowPM, ConfigFieldAllowPrivateMessages:
		rcf.AllowPrivateMessages.UpdateField(formFieldSingleString(field.Values), formFieldOptionsValues(field.Options))

	case ConfigFieldAllowInvites, ConfigFieldAllowMemberInvites:
		rcf.OccupantsCanInvite = formFieldBool(field.Values)

	case ConfigFieldCanChangeSubject:
		rcf.OccupantsCanChangeSubject = formFieldBool(field.Values)

	case ConfigFieldEnableLogging, ConfigFieldEnableArchiving:
		rcf.Logged = formFieldBool(field.Values)

	case ConfigFieldMemberList:
		rcf.RetrieveMembersList.UpdateField(field.Values, formFieldOptionsValues(field.Options))

	case ConfigFieldLanguage:
		rcf.setFieldX(ConfigFieldLanguage, field)

	case ConfigFieldPubsub:
		rcf.AssociatedPublishSubscribeNode = formFieldSingleString(field.Values)

	case ConfigFieldMaxOccupantsNumber:
		rcf.MaxOccupantsNumber.UpdateField(formFieldSingleString(field.Values), formFieldOptionsValues(field.Options))

	case ConfigFieldMembersOnly:
		rcf.MembersOnly = formFieldBool(field.Values)

	case ConfigFieldModerated:
		rcf.Moderated = formFieldBool(field.Values)

	case ConfigFieldPasswordProtected:
		rcf.PasswordProtected = formFieldBool(field.Values)

	case ConfigFieldIsPersistent:
		rcf.Persistent = formFieldBool(field.Values)

	case ConfigFieldPresenceBroadcast:
		rcf.PresenceBroadcast.UpdateField(field.Values, formFieldOptionsValues(field.Options))

	case ConfigFieldIsPublic:
		rcf.Public = formFieldBool(field.Values)

	case ConfigFieldRoomAdmins:
		rcf.Admins = formFieldJidList(field.Values)

	case ConfigFieldRoomDescription:
		rcf.setFieldX(ConfigFieldRoomDescription, field)

	case ConfigFieldRoomName:
		rcf.setFieldX(ConfigFieldRoomName, field)

	case ConfigFieldOwners:
		rcf.Owners = formFieldJidList(field.Values)

	case ConfigFieldPassword:
		rcf.Password = formFieldSingleString(field.Values)

	case ConfigFieldWhoIs:
		rcf.Whois.UpdateField(formFieldSingleString(field.Values), formFieldOptionsValues(field.Options))

	default:
		rcf.setFieldX(field.Var, field)
	}
}

// UpdateFieldValueByName finds a field from the unknown fields list by their name and updates their value
func (rcf *RoomConfigForm) UpdateFieldValueByName(name string, value interface{}) {
	for _, field := range rcf.Fields {
		if field.Name() == name {
			field.SetValue(value)
			return
		}
	}
}

func (rcf *RoomConfigForm) setFieldX(id string, field xmppData.FormFieldX) {
	if field.Type != RoomConfigFieldHidden && field.Type != RoomConfigFieldFixed {
		rcf.Fields[id] = roomConfigFormFieldFactory(field)
	}
}

func roomConfigFormFieldFactory(field xmppData.FormFieldX) HasRoomConfigFormField {
	f := newRoomConfigFormField(field.Var, field.Type, field.Label, field.Desc)

	switch field.Type {
	case RoomConfigFieldText, RoomConfigFieldTextPrivate:
		f.SetValue(formFieldSingleString(field.Values))

	case RoomConfigFieldTextMulti:
		f.SetValue(strings.Join(field.Values, "\n"))

	case RoomConfigFieldBoolean:
		f.SetValue(formFieldBool(field.Values))

	case RoomConfigFieldList:
		ls := newConfigListSingleField(nil)
		ls.UpdateField(formFieldSingleString(field.Values), formFieldOptionsValues(field.Options))
		f.SetValue(ls)

	case RoomConfigFieldListMulti:
		lm := newConfigListMultiField(nil)
		lm.UpdateField(field.Values, formFieldOptionsValues(field.Options))
		f.SetValue(lm)

	case RoomConfigFieldJidMulti:
		f.SetValue(formFieldJidList(field.Values))

	default:
		f.SetValue(field.Values)
	}

	return f
}

// GetStringValue returns the value of a string type field
func (rcf *RoomConfigForm) GetStringValue(identifier string) string {
	return rcf.Fields[identifier].ValueX()[0]
}

// UpdateFieldValue updates the value of a field based on their identifier
func (rcf *RoomConfigForm) UpdateFieldValue(identifier string, value interface{}) {
	rcf.Fields[identifier].SetValue(value)
}

func formFieldBool(values []string) bool {
	return len(values) > 0 && (strings.ToLower(values[0]) == "true" || values[0] == "1")
}

func formFieldSingleString(values []string) string {
	if len(values) > 0 {
		return values[0]
	}
	return ""
}

func formFieldOptionsValues(options []xmppData.FormFieldOptionX) (list []string) {
	for _, o := range options {
		list = append(list, o.Value)
	}

	return list
}

func formFieldJidList(values []string) (list []jid.Any) {
	for _, v := range values {
		list = append(list, jid.Parse(v))
	}
	return list
}

func jidListToStringList(jidList []jid.Any) (result []string) {
	for _, j := range jidList {
		result = append(result, j.String())
	}
	return result
}
