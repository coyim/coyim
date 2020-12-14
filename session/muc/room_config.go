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
		v := formFieldSingleString(field.Values)
		rcf.MaxHistoryFetch.SetValue(v)

	case "muc#roomconfig_allowpm":
		v := formFieldSingleString(field.Values)
		rcf.AllowPrivateMessages.SetValue(v)

	case "muc#roomconfig_allowinvites":
		v := formFieldBool(field.Values)
		rcf.OccupantsCanInvite.SetValue(v)

	case "muc#roomconfig_changesubject":
		v := formFieldBool(field.Values)
		rcf.OccupantsCanChangeSubject.SetValue(v)

	case "muc#roomconfig_enablelogging":
		rcf.Logged.SetValue(formFieldBool(field.Values))

	case "muc#roomconfig_getmemberlist":
		o := formFieldOptionsValues(field.Options)
		rcf.RetrieveMembersList.SetValues(field.Values)
		rcf.RetrieveMembersList.SetOptions(o)

	case "muc#roomconfig_lang":
		v := formFieldSingleString(field.Values)
		rcf.Language.SetValue(v)

	case "muc#roomconfig_pubsub":
		v := formFieldSingleString(field.Values)
		rcf.AssociatedPublishSubscribeNode.SetValue(v)

	case "muc#roomconfig_maxusers":
		v := formFieldSingleString(field.Values)
		o := formFieldOptionsValues(field.Options)
		rcf.MaxOccupantsNumber.SetValue(v)
		rcf.MaxOccupantsNumber.SetOptions(o)

	case "muc#roomconfig_membersonly":
		v := formFieldBool(field.Values)
		rcf.MembersOnly.SetValue(v)

	case "muc#roomconfig_moderatedroom":
		v := formFieldBool(field.Values)
		rcf.Moderated.SetValue(v)

	case "muc#roomconfig_passwordprotectedroom":
		v := formFieldBool(field.Values)
		rcf.PasswordProtected.SetValue(v)

	case "muc#roomconfig_persistentroom":
		v := formFieldBool(field.Values)
		rcf.Persistent.SetValue(v)

	case "muc#roomconfig_presencebroadcast":
		o := formFieldOptionsValues(field.Options)
		rcf.PresenceBroadcast.SetValues(field.Values)
		rcf.PresenceBroadcast.SetOptions(o)

	case "muc#roomconfig_publicroom":
		v := formFieldBool(field.Values)
		rcf.Persistent.SetValue(v)

	case "muc#roomconfig_roomadmins":
		v := formFieldJidList(field.Values)
		rcf.Admins.SetValues(v)

	case "muc#roomconfig_roomdesc":
		v := formFieldSingleString(field.Values)
		rcf.Description.SetValue(v)

	case "muc#roomconfig_roomname":
		v := formFieldSingleString(field.Values)
		rcf.Title.SetValue(v)

	case "muc#roomconfig_roomowners":
		v := formFieldJidList(field.Values)
		rcf.Owners.SetValues(v)

	case "muc#roomconfig_roomsecret":
		v := formFieldSingleString(field.Values)
		rcf.Password.SetValue(v)

	case "muc#roomconfig_whois":
		v := formFieldSingleString(field.Values)
		o := formFieldOptionsValues(field.Options)
		rcf.Whois.SetValue(v)
		rcf.Whois.SetOptions(o)
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
