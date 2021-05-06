package muc

import (
	"strconv"

	xmppData "github.com/coyim/coyim/xmpp/data"
)

const (
	// ConfigFieldFormType represents the configuration form type field
	ConfigFieldFormType = "FORM_TYPE"
	// ConfigFieldRoomName represents the var value of the "room name" configuration field
	ConfigFieldRoomName = "muc#roomconfig_roomname"
	// ConfigFieldRoomDescription represents the var value of the "room description" configuration field
	ConfigFieldRoomDescription = "muc#roomconfig_roomdesc"
	// ConfigFieldEnableLogging represents the var value of the "enable logging" configuration field
	ConfigFieldEnableLogging = "muc#roomconfig_enablelogging"
	// ConfigFieldEnableArchiving represents the var value of the "enable archiving" configuration field
	ConfigFieldEnableArchiving = "muc#roomconfig_enablearchiving"
	// ConfigFieldMemberList represents the var value of the "get members list" configuration field
	ConfigFieldMemberList = "muc#roomconfig_getmemberlist"
	// ConfigFieldLanguage represents the var value of the "room language" configuration field
	ConfigFieldLanguage = "muc#roomconfig_lang"
	// ConfigFieldPubsub represents the var value of the "pubsub" configuration field
	ConfigFieldPubsub = "muc#roomconfig_pubsub"
	// ConfigFieldCanChangeSubject represents the var value of the "change subject" configuration field
	ConfigFieldCanChangeSubject = "muc#roomconfig_changesubject"
	// ConfigFieldAllowInvites represents the var value of the "allow invites" configuration field
	ConfigFieldAllowInvites = "muc#roomconfig_allowinvites"
	// ConfigFieldAllowMemberInvites represents the var value of the "allow member invites" configuration field
	ConfigFieldAllowMemberInvites = "{http://prosody.im/protocol/muc}roomconfig_allowmemberinvites"
	// ConfigFieldAllowPM represents the var value of the "allow private messages" configuration field
	ConfigFieldAllowPM = "muc#roomconfig_allowpm"
	// ConfigFieldAllowPrivateMessages represents the var value of the "allow private messages" configuration field
	ConfigFieldAllowPrivateMessages = "allow_private_messages"
	// ConfigFieldMaxOccupantsNumber represents the var value of the "max users" configuration field
	ConfigFieldMaxOccupantsNumber = "muc#roomconfig_maxusers"
	// ConfigFieldIsPublic represents the var value of the "public room" configuration field
	ConfigFieldIsPublic = "muc#roomconfig_publicroom"
	// ConfigFieldIsPersistent represents the var value of the "persistent room" configuration field
	ConfigFieldIsPersistent = "muc#roomconfig_persistentroom"
	// ConfigFieldPresenceBroadcast represents the var value of the "presence broadcast" configuration field
	ConfigFieldPresenceBroadcast = "muc#roomconfig_presencebroadcast"
	// ConfigFieldModerated represents the var value of the "moderated room" configuration field
	ConfigFieldModerated = "muc#roomconfig_moderatedroom"
	// ConfigFieldMembersOnly represents the var value of the "members only" configuration field
	ConfigFieldMembersOnly = "muc#roomconfig_membersonly"
	// ConfigFieldPasswordProtected represents the var value of the "password protected room" configuration field
	ConfigFieldPasswordProtected = "muc#roomconfig_passwordprotectedroom"
	// ConfigFieldPassword represents the var value of the "room secret" configuration field
	ConfigFieldPassword = "muc#roomconfig_roomsecret"
	// ConfigFieldOwners represents the var value of the "room owners" configuration field
	ConfigFieldOwners = "muc#roomconfig_roomowners"
	// ConfigFieldWhoIs represents the var value of the "who is" configuration field
	ConfigFieldWhoIs = "muc#roomconfig_whois"
	// ConfigFieldMaxHistoryFetch represents the var value of the "max history fetch" configuration field
	ConfigFieldMaxHistoryFetch = "muc#maxhistoryfetch"
	// ConfigFieldMaxHistoryLength represents the var value of the "history length" configuration field
	ConfigFieldMaxHistoryLength = "muc#roomconfig_historylength"
	// ConfigFieldRoomAdmins represents the var value of the "room admins" configuration field
	ConfigFieldRoomAdmins = "muc#roomconfig_roomadmins"
)

// RoomConfigForm represents a room configuration form
type RoomConfigForm struct {
	Title                          string
	Description                    string
	Language                       string
	Password                       string
	AssociatedPublishSubscribeNode string

	OccupantsCanInvite        bool
	OccupantsCanChangeSubject bool
	Logged                    bool
	MembersOnly               bool
	Moderated                 bool
	PasswordProtected         bool
	Persistent                bool
	Public                    bool

	MaxHistoryFetch      *RoomConfigFieldListValue
	AllowPrivateMessages *RoomConfigFieldListValue
	MaxOccupantsNumber   *RoomConfigFieldListValue
	Whois                *RoomConfigFieldListValue

	RetrieveMembersList *RoomConfigFieldListMultiValue
	PresenceBroadcast   *RoomConfigFieldListMultiValue

	Admins *RoomConfigFieldJidMultiValue
	Owners *RoomConfigFieldJidMultiValue

	formType   string
	fieldNames map[string]int
	Fields     []*RoomConfigFormField
}

// NewRoomConfigForm creates a new room configuration form instance
func NewRoomConfigForm(form *xmppData.Form) *RoomConfigForm {
	cf := &RoomConfigForm{
		fieldNames: map[string]int{},
	}

	cf.initListSingleValueFields()
	cf.initListMultiValueFields()
	cf.initJidMultiValueFields()

	cf.setFormFields(form.Fields)

	return cf
}

func (rcf *RoomConfigForm) initListSingleValueFields() {
	rcf.MaxHistoryFetch = newRoomConfigFieldListValue(nil, maxHistoryFetchDefaultOptions)
	rcf.AllowPrivateMessages = newRoomConfigFieldListValue(nil, allowPrivateMessagesDefaultOptions)
	rcf.MaxOccupantsNumber = newRoomConfigFieldListValue(nil, maxOccupantsNumberDefaultOptions)
	rcf.Whois = newRoomConfigFieldListValue(nil, whoisDefaultOptions)
}

func (rcf *RoomConfigForm) initListMultiValueFields() {
	rcf.RetrieveMembersList = newRoomConfigFieldListMultiValue(nil, retrieveMembersListDefaultOptions)
	rcf.PresenceBroadcast = newRoomConfigFieldListMultiValue(nil, presenceBroadcastDefaultOptions)
}

func (rcf *RoomConfigForm) initJidMultiValueFields() {
	rcf.Admins = newRoomConfigFieldJidMultiValue(nil)
	rcf.Owners = newRoomConfigFieldJidMultiValue(nil)
}

func (rcf *RoomConfigForm) setFormFields(fields []xmppData.FormFieldX) {
	for idx, field := range fields {
		if field.Var != "" {
			rcf.fieldNames[field.Var] = idx
			rcf.setField(field)
		}
	}
}

// GetFormData returns a representation of the room config FORM_TYPE as described in the
// XMPP specification for MUC
//
// For more information see:
// https://xmpp.org/extensions/xep-0045.html#createroom-reserved
// https://xmpp.org/extensions/xep-0045.html#example-163
func (rcf *RoomConfigForm) GetFormData() *xmppData.Form {
	formFields := []xmppData.FormFieldX{}

	for fieldName := range rcf.fieldNames {
		var values []string

		switch fieldName {
		case ConfigFieldFormType:
			values = []string{rcf.formType}

		case ConfigFieldRoomName:
			values = []string{rcf.Title}

		case ConfigFieldRoomDescription:
			values = []string{rcf.Description}

		case ConfigFieldEnableLogging, ConfigFieldEnableArchiving:
			values = []string{strconv.FormatBool(rcf.Logged)}

		case ConfigFieldMemberList:
			values = rcf.RetrieveMembersList.Value()

		case ConfigFieldLanguage:
			values = []string{rcf.Language}

		case ConfigFieldPubsub:
			values = []string{rcf.AssociatedPublishSubscribeNode}

		case ConfigFieldCanChangeSubject:
			values = []string{strconv.FormatBool(rcf.OccupantsCanChangeSubject)}

		case ConfigFieldAllowInvites, ConfigFieldAllowMemberInvites:
			values = []string{strconv.FormatBool(rcf.OccupantsCanInvite)}

		case ConfigFieldAllowPM, ConfigFieldAllowPrivateMessages:
			values = rcf.AllowPrivateMessages.Value()

		case ConfigFieldMaxOccupantsNumber:
			values = rcf.MaxOccupantsNumber.Value()

		case ConfigFieldIsPublic:
			values = []string{strconv.FormatBool(rcf.Public)}

		case ConfigFieldIsPersistent:
			values = []string{strconv.FormatBool(rcf.Persistent)}

		case ConfigFieldPresenceBroadcast:
			values = rcf.PresenceBroadcast.Value()

		case ConfigFieldModerated:
			values = []string{strconv.FormatBool(rcf.Moderated)}

		case ConfigFieldMembersOnly:
			values = []string{strconv.FormatBool(rcf.MembersOnly)}

		case ConfigFieldPasswordProtected:
			values = []string{strconv.FormatBool(rcf.PasswordProtected)}

		case ConfigFieldPassword:
			values = []string{rcf.Password}

		case ConfigFieldOwners:
			values = rcf.Owners.Value()

		case ConfigFieldWhoIs:
			values = rcf.Whois.Value()

		case ConfigFieldMaxHistoryFetch, ConfigFieldMaxHistoryLength:
			values = rcf.MaxHistoryFetch.Value()

		case ConfigFieldRoomAdmins:
			values = rcf.Admins.Value()
		}

		formFields = append(formFields, xmppData.FormFieldX{
			Var:    fieldName,
			Values: values,
		})
	}

	for _, f := range rcf.Fields {
		formFields = append(formFields, xmppData.FormFieldX{
			Var:    f.Name,
			Values: f.Value(),
		})
	}

	return &xmppData.Form{
		Type:   "submit",
		Fields: formFields,
	}
}

func (rcf *RoomConfigForm) setField(field xmppData.FormFieldX) {
	switch field.Var {
	case ConfigFieldFormType:
		rcf.formType = formFieldSingleString(field.Values)

	case ConfigFieldMaxHistoryFetch, ConfigFieldMaxHistoryLength:
		rcf.MaxHistoryFetch.SetValue(formFieldSingleString(field.Values))
		rcf.MaxHistoryFetch.SetOptions(formFieldOptionsValues(field.Options))

	case ConfigFieldAllowPM, ConfigFieldAllowPrivateMessages:
		rcf.AllowPrivateMessages.SetValue(formFieldSingleString(field.Values))
		rcf.AllowPrivateMessages.SetOptions(formFieldOptionsValues(field.Options))

	case ConfigFieldAllowInvites, ConfigFieldAllowMemberInvites:
		rcf.OccupantsCanInvite = formFieldBool(field.Values)

	case ConfigFieldCanChangeSubject:
		rcf.OccupantsCanChangeSubject = formFieldBool(field.Values)

	case ConfigFieldEnableLogging, ConfigFieldEnableArchiving:
		rcf.Logged = formFieldBool(field.Values)

	case ConfigFieldMemberList:
		rcf.RetrieveMembersList.SetValue(field.Values)
		rcf.RetrieveMembersList.SetOptions(formFieldOptionsValues(field.Options))

	case ConfigFieldLanguage:
		rcf.Language = formFieldSingleString(field.Values)

	case ConfigFieldPubsub:
		rcf.AssociatedPublishSubscribeNode = formFieldSingleString(field.Values)

	case ConfigFieldMaxOccupantsNumber:
		rcf.MaxOccupantsNumber.SetValue(formFieldSingleString(field.Values))
		rcf.MaxOccupantsNumber.SetOptions(formFieldOptionsValues(field.Options))

	case ConfigFieldMembersOnly:
		rcf.MembersOnly = formFieldBool(field.Values)

	case ConfigFieldModerated:
		rcf.Moderated = formFieldBool(field.Values)

	case ConfigFieldPasswordProtected:
		rcf.PasswordProtected = formFieldBool(field.Values)

	case ConfigFieldIsPersistent:
		rcf.Persistent = formFieldBool(field.Values)

	case ConfigFieldPresenceBroadcast:
		rcf.PresenceBroadcast.SetValue(field.Values)
		rcf.PresenceBroadcast.SetOptions(formFieldOptionsValues(field.Options))

	case ConfigFieldIsPublic:
		rcf.Public = formFieldBool(field.Values)

	case ConfigFieldRoomAdmins:
		rcf.Admins.SetValue(field.Values)

	case ConfigFieldRoomDescription:
		rcf.Description = formFieldSingleString(field.Values)

	case ConfigFieldRoomName:
		rcf.Title = formFieldSingleString(field.Values)

	case ConfigFieldOwners:
		rcf.Owners.SetValue(field.Values)

	case ConfigFieldPassword:
		rcf.Password = formFieldSingleString(field.Values)

	case ConfigFieldWhoIs:
		rcf.Whois.SetValue(formFieldSingleString(field.Values))
		rcf.Whois.SetOptions(formFieldOptionsValues(field.Options))

	default:
		if field.Type != RoomConfigFieldHidden && field.Type != RoomConfigFieldFixed {
			rcf.Fields = append(rcf.Fields, newRoomConfigFormField(field))
		}
	}
}

// UpdateFieldValueByName finds a field from the unknown fields list by their name and updates their value
func (rcf *RoomConfigForm) UpdateFieldValueByName(name string, value interface{}) {
	for _, field := range rcf.Fields {
		if field.Name == name {
			field.SetValue(value)
			return
		}
	}
}

func formFieldBool(values []string) bool {
	if len(values) > 0 {
		v, err := strconv.ParseBool(values[0])
		if err == nil {
			return v
		}
	}
	return false
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
