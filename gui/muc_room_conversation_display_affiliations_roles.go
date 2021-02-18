package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc/data"
)

func getDisplayForOccupantAffiliationRoleUpdate(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	d := newAffiliationUpdateDisplayData(data.AffiliationUpdate{
		Nickname: affiliationRoleUpdate.Nickname,
		Reason:   affiliationRoleUpdate.Reason,
		New:      affiliationRoleUpdate.NewAffiliation,
		Previous: affiliationRoleUpdate.PreviousAffiliation,
		Actor:    affiliationRoleUpdate.Actor,
	})

	message := displayAffiliationUpdateMessage(d, i18n.Localf("As a result, the role changed from %s to %s.",
		displayNameForRole(affiliationRoleUpdate.PreviousRole),
		displayNameForRole(affiliationRoleUpdate.NewRole)))

	return message
}

func getDisplayForSelfOccupantAffiliationRoleUpdate(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	d := newSelfAffiliationUpdateDisplayData(data.AffiliationUpdate{
		Nickname: affiliationRoleUpdate.Nickname,
		Reason:   affiliationRoleUpdate.Reason,
		New:      affiliationRoleUpdate.NewAffiliation,
		Previous: affiliationRoleUpdate.PreviousAffiliation,
		Actor:    affiliationRoleUpdate.Actor,
	})

	message := displayAffiliationUpdateMessage(d, i18n.Localf("As a result, your role changed from %s to %s.",
		displayNameForRole(affiliationRoleUpdate.PreviousRole),
		displayNameForRole(affiliationRoleUpdate.NewRole)))

	return message
}
