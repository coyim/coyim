package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc/data"
)

type roleUpdateDisplayer interface {
	displayForRoleChanged() string
	updateReason() string
}

func getDisplayForSelfOccupantRoleUpdate(roleUpdate data.RoleUpdate) string {
	d := newSelfRoleUpdateDisplayData(roleUpdate)
	return displayRoleUpdateMessage(d)
}

func displayRoleUpdateMessage(d roleUpdateDisplayer) (message string) {
	message = d.displayForRoleChanged()

	if reason := d.updateReason(); reason != "" {
		message = i18n.Localf("%s because: %s", message, reason)
	}

	return message
}

type roleUpdateDisplayData struct {
	nickname         string
	newRole          data.Role
	previousRole     data.Role
	actor            string
	actorAffiliation data.Affiliation
	reason           string
}

func newRoleUpdateDisplayData(roleUpdate data.RoleUpdate) *roleUpdateDisplayData {
	d := &roleUpdateDisplayData{
		nickname:     roleUpdate.Nickname,
		newRole:      roleUpdate.New,
		previousRole: roleUpdate.Previous,
		reason:       roleUpdate.Reason,
	}

	if roleUpdate.Actor != nil {
		d.actor = roleUpdate.Actor.Nickname
		d.actorAffiliation = roleUpdate.Actor.Affiliation
	}

	return d
}

type selfRoleUpdateDisplayData struct {
	*roleUpdateDisplayData
}

func newSelfRoleUpdateDisplayData(roleUpdate data.RoleUpdate) *selfRoleUpdateDisplayData {
	return &selfRoleUpdateDisplayData{
		newRoleUpdateDisplayData(roleUpdate),
	}
}

func (d *roleUpdateDisplayData) displayForRoleChanged() string {
	if d.actor == "" {
		return i18n.Localf("The role of %s was changed from %s to %s", d.nickname,
			displayNameForRole(d.previousRole),
			displayNameForRole(d.newRole))
	}
	return i18n.Localf("%s changed the role of %s from %s to %s",
		displayActorWithAffiliation(d.actor, d.actorAffiliation),
		d.nickname,
		displayNameForRole(d.previousRole),
		displayNameForRole(d.newRole),
	)
}

func (d *selfRoleUpdateDisplayData) displayForRoleChanged() string {
	if d.actor == "" {
		return i18n.Localf("Your role was changed from %s to %s",
			displayNameForRole(d.previousRole),
			displayNameForRole(d.newRole))
	}
	return i18n.Localf("%s changed your role from %s to %s",
		displayActorWithAffiliation(d.actor, d.actorAffiliation),
		displayNameForRole(d.previousRole),
		displayNameForRole(d.newRole),
	)
}

func (d *roleUpdateDisplayData) updateReason() string {
	return d.reason
}

func displayNameForRole(role data.Role) string {
	switch {
	case role.IsModerator():
		return i18n.Local("moderator")
	case role.IsParticipant():
		return i18n.Local("participant")
	case role.IsVisitor():
		return i18n.Local("visitor")
	default:
		return ""
	}
}
