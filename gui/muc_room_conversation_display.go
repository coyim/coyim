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

func displayActorWithAffiliation(actor string, affiliation data.Affiliation) string {
	if affiliation != nil {
		return i18n.Localf("The %s %s", displayNameForAffiliation(affiliation), actor)
	}
	return actor
}

func displayNameForAffiliation(affiliation data.Affiliation) string {
	switch {
	case affiliation.IsAdmin():
		return i18n.Local("administrator")
	case affiliation.IsOwner():
		return i18n.Local("owner")
	case affiliation.IsBanned():
		return i18n.Local("outcast")
	case affiliation.IsMember():
		return i18n.Local("member")
	default: // Other values get the default treatment
		return ""
	}
}

func displayNameForAffiliationWithPreposition(affiliation data.Affiliation) string {
	switch {
	case affiliation.IsAdmin():
		return i18n.Local("an administrator")
	case affiliation.IsOwner():
		return i18n.Local("an owner")
	case affiliation.IsMember():
		return i18n.Local("a member")
	case affiliation.IsBanned():
		return i18n.Local("a banned")
	default: // Other values get the default treatment
		return ""
	}
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
