package muc

import (
	"errors"
	"strconv"

	"github.com/coyim/coyim/session/muc/data"
	xmppData "github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
)

const (
	// ConfigFieldFormType represents `FORM_TYPE` field in room configuration form
	ConfigFieldFormType = "FORM_TYPE"
	// ConfigFieldRoomName represents `muc#roomconfig_roomname` field in room configuration form
	ConfigFieldRoomName = "muc#roomconfig_roomname"
	// ConfigFieldRoomDescription represents `muc#roomconfig_roomdesc` field in room configuration form
	ConfigFieldRoomDescription = "muc#roomconfig_roomdesc"
	// ConfigFieldEnableLogging represents `muc#roomconfig_enablelogging` field in room configuration form
	ConfigFieldEnableLogging = "muc#roomconfig_enablelogging"
	// ConfigFieldMemberList represents `muc#roomconfig_getmemberlist` field in room configuration form
	ConfigFieldMemberList = "muc#roomconfig_getmemberlist"
	// ConfigFieldLanguage represents `muc#roomconfig_lang` field in room configuration form
	ConfigFieldLanguage = "muc#roomconfig_lang"
	// ConfigFieldPubsub represents `muc#roomconfig_pubsub` field in room configuration form
	ConfigFieldPubsub = "muc#roomconfig_pubsub"
	// ConfigFieldCanChangeSubject represents `muc#roomconfig_changesubject` field in room configuration form
	ConfigFieldCanChangeSubject = "muc#roomconfig_changesubject"
	// ConfigFieldAllowInvites represents `muc#roomconfig_allowinvites` field in room configuration form
	ConfigFieldAllowInvites = "muc#roomconfig_allowinvites"
	// ConfigFieldAllowPrivateMessages represents `muc#roomconfig_allowpm` field in room configuration form
	ConfigFieldAllowPrivateMessages = "muc#roomconfig_allowpm"
	// ConfigFieldMaxOccupantsNumber represents `muc#roomconfig_maxusers` field in room configuration form
	ConfigFieldMaxOccupantsNumber = "muc#roomconfig_maxusers"
	// ConfigFieldIsPublic represents `muc#roomconfig_publicroom` field in room configuration form
	ConfigFieldIsPublic = "muc#roomconfig_publicroom"
	// ConfigFieldIsPersistent represents `muc#roomconfig_persistentroom` field in room configuration form
	ConfigFieldIsPersistent = "muc#roomconfig_persistentroom"
	// ConfigFieldPresenceBroadcast represents `muc#roomconfig_presencebroadcast` field in room configuration form
	ConfigFieldPresenceBroadcast = "muc#roomconfig_presencebroadcast"
	// ConfigFieldModerated represents `muc#roomconfig_moderatedroom` field in room configuration form
	ConfigFieldModerated = "muc#roomconfig_moderatedroom"
	// ConfigFieldMembersOnly represents `muc#roomconfig_membersonly` field in room configuration form
	ConfigFieldMembersOnly = "muc#roomconfig_membersonly"
	// ConfigFieldPasswordProtected represents `muc#roomconfig_passwordprotectedroom` field in room configuration form
	ConfigFieldPasswordProtected = "muc#roomconfig_passwordprotectedroom"
	// ConfigFieldPassword represents `muc#roomconfig_roomsecret` field in room configuration form
	ConfigFieldPassword = "muc#roomconfig_roomsecret"
	// ConfigFieldOwners represents `muc#roomconfig_roomowners` field in room configuration form
	ConfigFieldOwners = "muc#roomconfig_roomowners"
	// ConfigFieldWhoIs represents `muc#roomconfig_whois` field in room configuration form
	ConfigFieldWhoIs = "muc#roomconfig_whois"
	// ConfigFieldMaxHistoryFetch represents `muc#maxhistoryfetch` field in room configuration form
	ConfigFieldMaxHistoryFetch = "muc#maxhistoryfetch"
	// ConfigFieldRoomAdmins represents `muc#roomconfig_roomadmins` field in room configuration form
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
	Language                       string
	AssociatedPublishSubscribeNode string
	MaxOccupantsNumber             ConfigListSingleField
	MembersOnly                    bool
	Moderated                      bool
	PasswordProtected              bool
	Persistent                     bool
	PresenceBroadcast              ConfigListMultiField
	Public                         bool
	Admins                         []jid.Any
	Description                    string
	Title                          string
	Owners                         []jid.Any
	Password                       string
	Whois                          ConfigListSingleField
}

// NewRoomConfigRom creates a new room configuration form instance
func NewRoomConfigRom(form *xmppData.Form) *RoomConfigForm {
	cf := &RoomConfigForm{}

	cf.MaxHistoryFetch = newConfigListSingleField([]string{
		RoomConfigOption50,
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

	cf.SetFormFields(form)

	return cf
}

// Submit will send the configuration form data to the service that provided the initial configuration for the room.
//
// Please note that the configuration process can be canceled but if the room owner
// cancels the initial configuration, the service will destroy the room and will send
// an unavailable presence to the room owner
//
// For more information see:
// https://xmpp.org/extensions/xep-0045.html#createroom-reserved
// https://xmpp.org/extensions/xep-0045.html#example-163
func (rcf *RoomConfigForm) Submit() error {
	return errors.New("RoomConfigForm: Submit() not yet implemented")
}

// GetFormData description
func (rcf *RoomConfigForm) GetFormData() *xmppData.Form {
	fields := map[string][]string{
		ConfigFieldFormType:             []string{"http://jabber.org/protocol/muc#roomconfig"},
		ConfigFieldRoomName:             []string{rcf.Title},
		ConfigFieldRoomDescription:      []string{rcf.Description},
		ConfigFieldEnableLogging:        []string{strconv.FormatBool(rcf.Logged)},
		ConfigFieldCanChangeSubject:     []string{strconv.FormatBool(rcf.OccupantsCanChangeSubject)},
		ConfigFieldAllowInvites:         []string{strconv.FormatBool(rcf.OccupantsCanInvite)},
		ConfigFieldAllowPrivateMessages: []string{rcf.AllowPrivateMessages.CurrentValue()},
		ConfigFieldMaxOccupantsNumber:   []string{rcf.MaxOccupantsNumber.CurrentValue()},
		ConfigFieldIsPublic:             []string{strconv.FormatBool(rcf.Public)},
		ConfigFieldIsPersistent:         []string{strconv.FormatBool(rcf.Persistent)},
		ConfigFieldModerated:            []string{strconv.FormatBool(rcf.Moderated)},
		ConfigFieldMembersOnly:          []string{strconv.FormatBool(rcf.MembersOnly)},
		ConfigFieldPasswordProtected:    []string{strconv.FormatBool(rcf.PasswordProtected)},
		ConfigFieldPassword:             []string{rcf.Password},
		ConfigFieldWhoIs:                []string{rcf.Whois.CurrentValue()},
		ConfigFieldMaxHistoryFetch:      []string{rcf.MaxHistoryFetch.CurrentValue()},
		ConfigFieldLanguage:             []string{rcf.Language},
		ConfigFieldRoomAdmins:           jidListToStringList(rcf.Admins),
	}

	formFields := []xmppData.FormFieldX{}
	for name, values := range fields {
		formFields = append(formFields, xmppData.FormFieldX{
			Var:    name,
			Values: values,
		})
	}

	return &xmppData.Form{
		Type:   "submit",
		Fields: formFields,
	}
}

func jidListToStringList(jidList []jid.Any) (result []string) {
	for _, j := range jidList {
		result = append(result, j.String())
	}
	return result
}

// SetFormFields extract the form fields and updates the room config form properties based on each data
func (rcf *RoomConfigForm) SetFormFields(form *xmppData.Form) {
	for _, field := range form.Fields {
		rcf.setField(field)
	}
}

func (rcf *RoomConfigForm) setField(field xmppData.FormFieldX) {
	switch field.Var {
	case ConfigFieldMaxHistoryFetch:
		rcf.MaxHistoryFetch.UpdateField(formFieldSingleString(field.Values), formFieldOptionsValues(field.Options))

	case ConfigFieldAllowPrivateMessages:
		rcf.AllowPrivateMessages.UpdateField(formFieldSingleString(field.Values), formFieldOptionsValues(field.Options))

	case ConfigFieldAllowInvites:
		rcf.OccupantsCanInvite = formFieldBool(field.Values)

	case ConfigFieldCanChangeSubject:
		rcf.OccupantsCanChangeSubject = formFieldBool(field.Values)

	case ConfigFieldEnableLogging:
		rcf.Logged = formFieldBool(field.Values)

	case ConfigFieldMemberList:
		rcf.RetrieveMembersList.UpdateField(field.Values, formFieldOptionsValues(field.Options))

	case ConfigFieldLanguage:
		rcf.Language = formFieldSingleString(field.Values)

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
		rcf.Description = formFieldSingleString(field.Values)

	case ConfigFieldRoomName:
		rcf.Title = formFieldSingleString(field.Values)

	case ConfigFieldOwners:
		rcf.Owners = formFieldJidList(field.Values)

	case ConfigFieldPassword:
		rcf.Password = formFieldSingleString(field.Values)

	case ConfigFieldWhoIs:
		rcf.Whois.UpdateField(formFieldSingleString(field.Values), formFieldOptionsValues(field.Options))
	}
}

func formFieldBool(values []string) bool {
	return len(values) > 0 && (values[0] == "true" || values[0] == "1")
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

func formFieldInt(values []string) int {
	if len(values) > 0 {
		res, e := strconv.Atoi(values[0])
		if e == nil {
			return res
		}
	}
	return 0
}

func formFieldRoles(values []string) (roles []data.Role) {
	for _, v := range values {
		r, err := data.RoleFromString(v)
		if err == nil {
			roles = append(roles, r)
		}
	}
	return roles
}

func formFieldJidList(values []string) (list []jid.Any) {
	for _, v := range values {
		list = append(list, jid.Parse(v))
	}
	return list
}
