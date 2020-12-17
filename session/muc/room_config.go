package muc

import (
	"errors"
	"strconv"

	"github.com/coyim/coyim/session/muc/data"
	xmppData "github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
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

// SetFormFields extract the form fields and updates the room config form properties based on each data
func (rcf *RoomConfigForm) SetFormFields(form *xmppData.Form) {
	for _, field := range form.Fields {
		rcf.setField(field)
	}
}

func (rcf *RoomConfigForm) setField(field xmppData.FormFieldX) {
	switch field.Var {
	case "muc#maxhistoryfetch":
		rcf.MaxHistoryFetch.UpdateField(formFieldSingleString(field.Values), formFieldOptionsValues(field.Options))

	case "muc#roomconfig_allowpm":
		rcf.AllowPrivateMessages.UpdateField(formFieldSingleString(field.Values), formFieldOptionsValues(field.Options))

	case "muc#roomconfig_allowinvites":
		rcf.OccupantsCanInvite = formFieldBool(field.Values)

	case "muc#roomconfig_changesubject":
		rcf.OccupantsCanChangeSubject = formFieldBool(field.Values)

	case "muc#roomconfig_enablelogging":
		rcf.Logged = formFieldBool(field.Values)

	case "muc#roomconfig_getmemberlist":
		rcf.RetrieveMembersList.UpdateField(field.Values, formFieldOptionsValues(field.Options))

	case "muc#roomconfig_lang":
		rcf.Language = formFieldSingleString(field.Values)

	case "muc#roomconfig_pubsub":
		rcf.AssociatedPublishSubscribeNode = formFieldSingleString(field.Values)

	case "muc#roomconfig_maxusers":
		rcf.MaxOccupantsNumber.UpdateField(formFieldSingleString(field.Values), formFieldOptionsValues(field.Options))

	case "muc#roomconfig_membersonly":
		rcf.MembersOnly = formFieldBool(field.Values)

	case "muc#roomconfig_moderatedroom":
		rcf.Moderated = formFieldBool(field.Values)

	case "muc#roomconfig_passwordprotectedroom":
		rcf.PasswordProtected = formFieldBool(field.Values)

	case "muc#roomconfig_persistentroom":
		rcf.Persistent = formFieldBool(field.Values)

	case "muc#roomconfig_presencebroadcast":
		rcf.PresenceBroadcast.UpdateField(field.Values, formFieldOptionsValues(field.Options))

	case "muc#roomconfig_publicroom":
		rcf.Public = formFieldBool(field.Values)

	case "muc#roomconfig_roomadmins":
		rcf.Admins = formFieldJidList(field.Values)

	case "muc#roomconfig_roomdesc":
		rcf.Description = formFieldSingleString(field.Values)

	case "muc#roomconfig_roomname":
		rcf.Title = formFieldSingleString(field.Values)

	case "muc#roomconfig_roomowners":
		rcf.Owners = formFieldJidList(field.Values)

	case "muc#roomconfig_roomsecret":
		rcf.Password = formFieldSingleString(field.Values)

	case "muc#roomconfig_whois":
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
