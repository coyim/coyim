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

	message := displayAffiliationUpdateMessage(d)

	message = i18n.Localf("%s %s", message, i18n.Localf("and as a result the role changed from %s to %s",
		displayNameForRole(affiliationRoleUpdate.PreviousRole),
		displayNameForRole(affiliationRoleUpdate.NewRole)))

	return message
}
