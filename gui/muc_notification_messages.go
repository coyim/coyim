package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc/data"
)

func getMUCNotificationMessageFrom(d interface{}) string {
	switch t := d.(type) {
	case data.AffiliationUpdate:
		return getAffiliationUpdateMessage(t)
	default:
		return ""
	}
}

func getAffiliationUpdateMessage(affiliationUpdate data.AffiliationUpdate) string {
	switch {
	case affiliationUpdate.New.IsNone():
		return getAffiliationRemovedMessage(affiliationUpdate)

	case affiliationUpdate.New.IsBanned():
		return getAffiliationBannedMessage(affiliationUpdate)

	case affiliationUpdate.Previous.IsNone():
		return getAffiliationAddedMessage(affiliationUpdate)

	default:
		return getAffiliationChangedMessage(affiliationUpdate)
	}
}

func getAffiliationRemovedMessage(affiliationUpdate data.AffiliationUpdate) string {
	if affiliationUpdate.Actor == nil {
		if affiliationUpdate.Reason == "" {
			return i18n.Localf("%s is not %s anymore.",
				affiliationUpdate.Nickname,
				displayNameForAffiliationWithPreposition(affiliationUpdate.Previous))
		}
		return i18n.Localf("%s is not %s anymore because: %s.",
			affiliationUpdate.Nickname,
			displayNameForAffiliationWithPreposition(affiliationUpdate.Previous),
			affiliationUpdate.Reason)
	}

	if affiliationUpdate.Reason == "" {
		return i18n.Localf("The %s %s changed the position of %s; %s is not %s anymore.",
			displayNameForAffiliation(affiliationUpdate.Actor.Affiliation),
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname,
			affiliationUpdate.Nickname,
			displayNameForAffiliationWithPreposition(affiliationUpdate.Previous))
	}

	return i18n.Localf("The %s %s changed the position of %s; %s is not %s anymore because: %s.",
		displayNameForAffiliation(affiliationUpdate.Actor.Affiliation),
		affiliationUpdate.Actor.Nickname,
		affiliationUpdate.Nickname,
		affiliationUpdate.Nickname,
		displayNameForAffiliationWithPreposition(affiliationUpdate.Previous),
		affiliationUpdate.Reason)
}

func getAffiliationBannedMessage(affiliationUpdate data.AffiliationUpdate) string {
	if affiliationUpdate.Actor == nil {
		if affiliationUpdate.Reason == "" {
			return i18n.Localf("%s was banned from the room.", affiliationUpdate.Nickname)
		}

		return i18n.Localf("%s was banned from the room because: %s.",
			affiliationUpdate.Nickname,
			affiliationUpdate.Reason)
	}

	if affiliationUpdate.Reason == "" {
		return i18n.Localf("The %s %s banned %s from the room.",
			displayNameForAffiliation(affiliationUpdate.Actor.Affiliation),
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname)
	}

	return i18n.Localf("The %s %s banned %s from the room because: %s.",
		displayNameForAffiliation(affiliationUpdate.Actor.Affiliation),
		affiliationUpdate.Actor.Nickname,
		affiliationUpdate.Nickname,
		affiliationUpdate.Reason)
}

func getAffiliationAddedMessage(affiliationUpdate data.AffiliationUpdate) string {
	if affiliationUpdate.Actor == nil {
		if affiliationUpdate.Reason == "" {
			return i18n.Localf("%s is now %s.",
				affiliationUpdate.Nickname,
				displayNameForAffiliationWithPreposition(affiliationUpdate.New))
		}

		return i18n.Localf("%s is now %s because: %s.",
			affiliationUpdate.Nickname,
			displayNameForAffiliationWithPreposition(affiliationUpdate.New),
			affiliationUpdate.Reason)
	}

	if affiliationUpdate.Reason == "" {
		return i18n.Localf("The %s %s changed the position of %s; %s is now %s.",
			displayNameForAffiliation(affiliationUpdate.Actor.Affiliation),
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname,
			affiliationUpdate.Nickname,
			displayNameForAffiliationWithPreposition(affiliationUpdate.New))
	}

	return i18n.Localf("The %s %s changed the position of %s; %s is now %s because: %s.",
		displayNameForAffiliation(affiliationUpdate.Actor.Affiliation),
		affiliationUpdate.Actor.Nickname,
		affiliationUpdate.Nickname,
		affiliationUpdate.Nickname,
		displayNameForAffiliationWithPreposition(affiliationUpdate.New),
		affiliationUpdate.Reason)
}

func getAffiliationChangedMessage(affiliationUpdate data.AffiliationUpdate) string {
	if affiliationUpdate.Actor == nil {
		if affiliationUpdate.Reason == "" {
			return i18n.Localf("The position of %s was changed from %s to %s.",
				affiliationUpdate.Nickname,
				displayNameForAffiliation(affiliationUpdate.Previous),
				displayNameForAffiliation(affiliationUpdate.New))
		}

		return i18n.Localf("The position of %s was changed from %s to %s because: %s.",
			affiliationUpdate.Nickname,
			displayNameForAffiliation(affiliationUpdate.Previous),
			displayNameForAffiliation(affiliationUpdate.New),
			affiliationUpdate.Reason)
	}

	if affiliationUpdate.Reason == "" {
		return i18n.Localf("The %s %s changed the position of %s from %s to %s.",
			displayNameForAffiliation(affiliationUpdate.Actor.Affiliation),
			affiliationUpdate.Actor.Nickname,
			affiliationUpdate.Nickname,
			displayNameForAffiliation(affiliationUpdate.Previous),
			displayNameForAffiliation(affiliationUpdate.New))
	}

	return i18n.Localf("The %s %s changed the position of %s from %s to %s because: %s.",
		displayNameForAffiliation(affiliationUpdate.Actor.Affiliation),
		affiliationUpdate.Actor.Nickname,
		affiliationUpdate.Nickname,
		displayNameForAffiliation(affiliationUpdate.Previous),
		displayNameForAffiliation(affiliationUpdate.New),
		affiliationUpdate.Reason)
}

func getRoleUpdateMessage(roleUpdate data.RoleUpdate) string {
	return ""
}

func getAffiliationRoleUpate(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	return ""
}
