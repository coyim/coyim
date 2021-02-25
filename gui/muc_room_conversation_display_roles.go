package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc/data"
)

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
