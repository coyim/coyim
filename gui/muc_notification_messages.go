package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc/data"
)

func getMUCNotificationMessageFrom(d interface{}) string {
	switch t := d.(type) {
	case data.AffiliationUpdate:
		return getAffiliationUpdateMessage(t)
	case data.RoleUpdate:
		return getRoleUpdateMessage(t)
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
	if roleUpdate.Actor == nil {
		if roleUpdate.Reason == "" {
			return i18n.Localf("The role of %s was changed from %s to %s.",
				roleUpdate.Nickname,
				displayNameForRole(roleUpdate.Previous),
				displayNameForRole(roleUpdate.New))
		}

		return i18n.Localf("The role of %s was changed from %s to %s because: %s.",
			roleUpdate.Nickname,
			displayNameForRole(roleUpdate.Previous),
			displayNameForRole(roleUpdate.New),
			roleUpdate.Reason)
	}

	if roleUpdate.Reason == "" {
		return i18n.Localf("The %s %s changed the role of %s from %s to %s.",
			displayNameForAffiliation(roleUpdate.Actor.Affiliation),
			roleUpdate.Actor.Nickname,
			roleUpdate.Nickname,
			displayNameForRole(roleUpdate.Previous),
			displayNameForRole(roleUpdate.New))
	}

	return i18n.Localf("The %s %s changed the role of %s from %s to %s because: %s.",
		displayNameForAffiliation(roleUpdate.Actor.Affiliation),
		roleUpdate.Actor.Nickname,
		roleUpdate.Nickname,
		displayNameForRole(roleUpdate.Previous),
		displayNameForRole(roleUpdate.New),
		roleUpdate.Reason)
}

func getAffiliationRoleUpate(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	switch {
	case affiliationRoleUpdate.NewAffiliation.IsNone() &&
		affiliationRoleUpdate.PreviousAffiliation.IsDifferentFrom(affiliationRoleUpdate.NewAffiliation):
		return getAffiliationRoleUpateForAffiliationRemoved(affiliationRoleUpdate)
	case affiliationRoleUpdate.PreviousAffiliation.IsNone() &&
		affiliationRoleUpdate.NewAffiliation.IsDifferentFrom(affiliationRoleUpdate.PreviousAffiliation):
		return getAffiliationRoleUpdateForAffiliationAdded(affiliationRoleUpdate)
	default:
		return getAffiliationRoleUpdateForUnexpectedSituation(affiliationRoleUpdate)
	}
}

func getAffiliationRoleUpateForAffiliationRemoved(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	if affiliationRoleUpdate.Actor == nil {
		if affiliationRoleUpdate.Reason == "" {
			return i18n.Localf("%s is not %s anymore. As a result, the role was changed from %s to %s.",
				affiliationRoleUpdate.Nickname,
				displayNameForAffiliationWithPreposition(affiliationRoleUpdate.PreviousAffiliation),
				displayNameForRole(affiliationRoleUpdate.PreviousRole),
				displayNameForRole(affiliationRoleUpdate.NewRole))
		}

		return i18n.Localf("%s is not %s anymore. As a result, the role was changed from %s to %s. The reason given was: %s.",
			affiliationRoleUpdate.Nickname,
			displayNameForAffiliationWithPreposition(affiliationRoleUpdate.PreviousAffiliation),
			displayNameForRole(affiliationRoleUpdate.PreviousRole),
			displayNameForRole(affiliationRoleUpdate.NewRole),
			affiliationRoleUpdate.Reason)
	}

	if affiliationRoleUpdate.Reason == "" {
		return i18n.Localf("The %s %s changed the position of %s; %s is not an %s anymore. As a result, the role was changed from %s to %s.",
			displayNameForAffiliation(affiliationRoleUpdate.Actor.Affiliation),
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname,
			affiliationRoleUpdate.Nickname,
			displayNameForAffiliation(affiliationRoleUpdate.PreviousAffiliation),
			displayNameForRole(affiliationRoleUpdate.PreviousRole),
			displayNameForRole(affiliationRoleUpdate.NewRole))
	}

	return i18n.Localf("The %s %s changed the position of %s; %s is not an %s anymore. "+
		"As a result, the role was changed from %s to %s. The reason given was: %s.",
		displayNameForAffiliation(affiliationRoleUpdate.Actor.Affiliation),
		affiliationRoleUpdate.Actor.Nickname,
		affiliationRoleUpdate.Nickname,
		affiliationRoleUpdate.Nickname,
		displayNameForAffiliation(affiliationRoleUpdate.PreviousAffiliation),
		displayNameForRole(affiliationRoleUpdate.PreviousRole),
		displayNameForRole(affiliationRoleUpdate.NewRole),
		affiliationRoleUpdate.Reason)
}

func getAffiliationRoleUpdateForAffiliationAdded(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	if affiliationRoleUpdate.Actor == nil {
		if affiliationRoleUpdate.Reason == "" {
			return i18n.Localf("The position of %s was changed to %s. As a result, the role was changed from %s to %s.",
				affiliationRoleUpdate.Nickname,
				displayNameForAffiliation(affiliationRoleUpdate.NewAffiliation),
				displayNameForRole(affiliationRoleUpdate.PreviousRole),
				displayNameForRole(affiliationRoleUpdate.NewRole))
		}

		return i18n.Localf("The position of %s was changed to %s. "+
			"As a result, the role was changed from %s to %s. The reason given was: %s.",
			affiliationRoleUpdate.Nickname,
			displayNameForAffiliation(affiliationRoleUpdate.NewAffiliation),
			displayNameForRole(affiliationRoleUpdate.PreviousRole),
			displayNameForRole(affiliationRoleUpdate.NewRole),
			affiliationRoleUpdate.Reason)
	}

	if affiliationRoleUpdate.Reason == "" {
		return i18n.Localf("The %s %s changed the position of %s to %s. "+
			"As a result, the role was changed from visitor to moderator.",
			displayNameForAffiliation(affiliationRoleUpdate.Actor.Affiliation),
			affiliationRoleUpdate.Actor.Nickname,
			affiliationRoleUpdate.Nickname,
			displayNameForAffiliation(affiliationRoleUpdate.NewAffiliation))
	}

	return i18n.Localf("The %s %s changed the position of %s to %s. "+
		"As a result, the role was changed from %s to %s. The reason given was: %s.",
		displayNameForAffiliation(affiliationRoleUpdate.Actor.Affiliation),
		affiliationRoleUpdate.Actor.Nickname,
		affiliationRoleUpdate.Nickname,
		displayNameForAffiliation(affiliationRoleUpdate.NewAffiliation),
		displayNameForRole(affiliationRoleUpdate.PreviousRole),
		displayNameForRole(affiliationRoleUpdate.NewRole),
		affiliationRoleUpdate.Reason)
}

func getAffiliationRoleUpdateForUnexpectedSituation(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	if affiliationRoleUpdate.Actor == nil {
		if affiliationRoleUpdate.Reason == "" {
			return i18n.Localf("The affiliation and the role of %s were changed.",
				affiliationRoleUpdate.Nickname)
		}

		return i18n.Localf("The affiliation and the role of %s were changed because: %s.",
			affiliationRoleUpdate.Nickname,
			affiliationRoleUpdate.Reason)
	}

	return i18n.Local("The owner louis changed the affiliation of superman. As a result, the role was changed too.")
}
