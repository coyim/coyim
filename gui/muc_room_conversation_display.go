package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc/data"
)

func getDisplayRoomSubjectForNickname(nickname, subject string) string {
	if nickname == "" {
		return i18n.Localf("Someone has updated the room subject to: \"%s\"", subject)
	}

	return i18n.Localf("%s updated the room subject to: \"%s\"", nickname, subject)
}

func getDisplayRoomSubject(subject string) string {
	if subject == "" {
		return i18n.Local("The room does not have a subject")
	}

	return i18n.Localf("The room subject is \"%s\"", subject)
}

func displayNameForRole(role data.Role) string {
	switch {
	case role.IsModerator():
		return i18n.Local("moderator")
	case role.IsParticipant():
		return i18n.Local("participant")
	case role.IsVisitor():
		return i18n.Local("visitor")
	case role.IsNone():
		return i18n.Local("removed")
	default: // Other values get the default treatment
		return ""
	}
}

func displayNameForRoleWithPreposition(role data.Role) string {
	switch {
	case role.IsModerator():
		return i18n.Local("a moderator")
	case role.IsParticipant():
		return i18n.Local("a participant")
	case role.IsVisitor():
		return i18n.Local("a visitor")
	case role.IsNone():
		return i18n.Local("removed")
	default:
		return ""
	}
}
