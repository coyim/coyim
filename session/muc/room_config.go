package muc

import (
	"strconv"

	"github.com/coyim/coyim/session/muc/data"
	xmppData "github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
)

// RoomConfigForm represents a room configuration form
type RoomConfigForm struct {
	data.RoomConfig
}

// NewRoomConfigRom creates a new room configuration form instance
func NewRoomConfigRom(form *xmppData.Form) *RoomConfigForm {
	cf := &RoomConfigForm{}

	cf.SetFormFields(form)

	return cf
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
		rcf.MaxHistoryFetch = formFieldInt(field.Values)

	case "muc#roomconfig_allowpm":
		rcf.AllowPrivateMessages = formFieldRoles(field.Values)

	case "muc#roomconfig_allowinvites":
		rcf.OccupantsCanInvite = formFieldBool(field.Values)

	case "muc#roomconfig_changesubject":
		rcf.OccupantsCanChangeSubject = formFieldBool(field.Values)

	case "muc#roomconfig_enablelogging":
		rcf.Logged = formFieldBool(field.Values)

	case "muc#roomconfig_getmemberlist":
		rcf.RetrieveMembersList = field.Values

	case "muc#roomconfig_lang":
		rcf.Language = formFieldSingleString(field.Values)

	case "muc#roomconfig_pubsub":
		rcf.AssociatedPublishSubscribeNode = formFieldSingleString(field.Values)

	case "muc#roomconfig_maxusers":
		rcf.MaxOccupantsNumber = formFieldInt(field.Values)

	case "muc#roomconfig_membersonly":
		rcf.MembersOnly = formFieldBool(field.Values)

	case "muc#roomconfig_moderatedroom":
		rcf.Moderated = formFieldBool(field.Values)

	case "muc#roomconfig_passwordprotectedroom":
		rcf.PasswordProtected = formFieldBool(field.Values)

	case "muc#roomconfig_persistentroom":
		rcf.Persistent = formFieldBool(field.Values)

	case "muc#roomconfig_presencebroadcast":
		rcf.PresenceBroadcast = field.Values

	case "muc#roomconfig_publicroom":
		rcf.Persistent = formFieldBool(field.Values)

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
		rcf.Whois = field.Values
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
