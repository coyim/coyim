package muc

import (
	"strconv"

	"github.com/coyim/coyim/session/muc/data"
	xmppData "github.com/coyim/coyim/xmpp/data"
)

const (
	// configFieldFormType represents the configuration form type field
	configFieldFormType = "FORM_TYPE"
	// configFieldFormTypeValue represents the value of field type of configuration form
	configFieldFormTypeValue = "http://jabber.org/protocol/muc#roomconfig"
	// configFieldRoomName represents the var value of the "room name" configuration field
	configFieldRoomName = "muc#roomconfig_roomname"
	// configFieldRoomDescription represents the var value of the "room description" configuration field
	configFieldRoomDescription = "muc#roomconfig_roomdesc"
	// configFieldEnableLogging represents the var value of the "enable logging" configuration field
	configFieldEnableLogging = "muc#roomconfig_enablelogging"
	// configFieldEnableArchiving represents the var value of the "enable archiving" configuration field
	configFieldEnableArchiving = "muc#roomconfig_enablearchiving"
	// configFieldMessageArchiveManagement represents the var value of the "mam" configuration field
	configFieldMessageArchiveManagement = "mam"
	// configFieldMemberList represents the var value of the "get members list" configuration field
	configFieldMemberList = "muc#roomconfig_getmemberlist"
	// configFieldLanguage represents the var value of the "room language" configuration field
	configFieldLanguage = "muc#roomconfig_lang"
	// configFieldPubsub represents the var value of the "pubsub" configuration field
	configFieldPubsub = "muc#roomconfig_pubsub"
	// configFieldCanChangeSubject represents the var value of the "change subject" configuration field
	configFieldCanChangeSubject = "muc#roomconfig_changesubject"
	// configFieldAllowInvites represents the var value of the "allow invites" configuration field
	configFieldAllowInvites = "muc#roomconfig_allowinvites"
	// configFieldAllowMemberInvites represents the var value of the "allow member invites" configuration field
	configFieldAllowMemberInvites = "{http://prosody.im/protocol/muc}roomconfig_allowmemberinvites"
	// configFieldAllowPM represents the var value of the "allow private messages" configuration field
	configFieldAllowPM = "muc#roomconfig_allowpm"
	// configFieldAllowPrivateMessages represents the var value of the "allow private messages" configuration field
	configFieldAllowPrivateMessages = "allow_private_messages"
	// configFieldMaxOccupantsNumber represents the var value of the "max users" configuration field
	configFieldMaxOccupantsNumber = "muc#roomconfig_maxusers"
	// configFieldIsPublic represents the var value of the "public room" configuration field
	configFieldIsPublic = "muc#roomconfig_publicroom"
	// configFieldIsPersistent represents the var value of the "persistent room" configuration field
	configFieldIsPersistent = "muc#roomconfig_persistentroom"
	// configFieldPresenceBroadcast represents the var value of the "presence broadcast" configuration field
	configFieldPresenceBroadcast = "muc#roomconfig_presencebroadcast"
	// configFieldModerated represents the var value of the "moderated room" configuration field
	configFieldModerated = "muc#roomconfig_moderatedroom"
	// configFieldMembersOnly represents the var value of the "members only" configuration field
	configFieldMembersOnly = "muc#roomconfig_membersonly"
	// configFieldPasswordProtected represents the var value of the "password protected room" configuration field
	configFieldPasswordProtected = "muc#roomconfig_passwordprotectedroom"
	// configFieldPassword represents the var value of the "room secret" configuration field
	configFieldPassword = "muc#roomconfig_roomsecret"
	// configFieldOwners represents the var value of the "room owners" configuration field
	configFieldOwners = "muc#roomconfig_roomowners"
	// configFieldWhoIs represents the var value of the "who is" configuration field
	configFieldWhoIs = "muc#roomconfig_whois"
	// configFieldMaxHistoryFetch represents the var value of the "max history fetch" configuration field
	configFieldMaxHistoryFetch = "muc#maxhistoryfetch"
	// configFieldMaxHistoryLength represents the var value of the "history length" configuration field
	configFieldMaxHistoryLength = "muc#roomconfig_historylength"
	// configFieldRoomAdmins represents the var value of the "room admins" configuration field
	configFieldRoomAdmins = "muc#roomconfig_roomadmins"
)

// RoomConfigForm represents a room configuration form
type RoomConfigForm struct {
	knownFields   map[RoomConfigFieldType]*RoomConfigFormField
	unknownFields []*RoomConfigFormField
	occupants     map[data.Affiliation][]*RoomOccupantItem
}

// NewRoomConfigForm creates a new room configuration form instance
func NewRoomConfigForm(form *xmppData.Form) *RoomConfigForm {
	cf := &RoomConfigForm{
		occupants:   map[data.Affiliation][]*RoomOccupantItem{},
		knownFields: map[RoomConfigFieldType]*RoomConfigFormField{},
	}

	cf.setFormFields(form.Fields)

	return cf
}

func (rcf *RoomConfigForm) setFormFields(fields []xmppData.FormFieldX) {
	for _, field := range fields {
		if field.Var != "" {
			if key, isKnown := getKnownRoomConfigFieldKey(field.Var); isKnown {
				rcf.knownFields[key] = newRoomConfigFormField(field)
			} else if field.Type != RoomConfigFieldFixed && field.Var != configFieldFormType {
				rcf.unknownFields = append(rcf.unknownFields, newRoomConfigFormField(field))
			}
		}
	}
}

// HasKnownField cheks if the filed was defined from the form
func (rcf *RoomConfigForm) HasKnownField(k RoomConfigFieldType) bool {
	_, ok := rcf.knownFields[k]
	return ok
}

// GetKnownField returns the known form field for the given key
func (rcf *RoomConfigForm) GetKnownField(k RoomConfigFieldType) (*RoomConfigFormField, bool) {
	if rcf.HasKnownField(k) {
		return rcf.knownFields[k], true
	}
	return nil, false
}

// GetConfiguredPassword returns the configured password in the room configuration form
func (rcf *RoomConfigForm) GetConfiguredPassword() (pwd string) {
	field, ok := rcf.GetKnownField(RoomConfigFieldPassword)
	if ok {
		pwd = field.Value()[0]
	}
	return
}

// GetUnknownFields returns the known form field for the given key
func (rcf *RoomConfigForm) GetUnknownFields() []*RoomConfigFormField {
	return rcf.unknownFields
}

// GetRoomOccupants returns all occupants in the room configuration form
func (rcf *RoomConfigForm) GetRoomOccupants() map[data.Affiliation][]*RoomOccupantItem {
	return rcf.occupants
}

// GetFormData returns a representation of the room config FORM_TYPE as described in the
// XMPP specification for MUC
//
// For more information see:
// https://xmpp.org/extensions/xep-0045.html#createroom-reserved
// https://xmpp.org/extensions/xep-0045.html#example-163
func (rcf *RoomConfigForm) GetFormData() *xmppData.Form {
	formFields := []xmppData.FormFieldX{{Var: configFieldFormType, Values: []string{configFieldFormTypeValue}}}

	for _, field := range rcf.knownFields {
		formFields = append(formFields, xmppData.FormFieldX{
			Var:    field.Name,
			Values: field.Value(),
		})
	}

	for _, field := range rcf.unknownFields {
		formFields = append(formFields, xmppData.FormFieldX{
			Var:    field.Name,
			Values: field.Value(),
		})
	}

	return &xmppData.Form{
		Type:   "submit",
		Fields: formFields,
	}
}

func (rcf *RoomConfigForm) getKnownFieldValue(fieldName string) ([]string, bool) {
	for _, field := range rcf.knownFields {
		if field.Name == fieldName {
			return field.Value(), true
		}
	}
	return nil, false
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
