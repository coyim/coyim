package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc/data"
)

func getDisplayForSelfOccupantAffiliationRoleUpdate(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	m := getMUCNotificationMessageFrom(data.AffiliationUpdate{
		Nickname: affiliationRoleUpdate.Nickname,
		Reason:   affiliationRoleUpdate.Reason,
		New:      affiliationRoleUpdate.NewAffiliation,
		Previous: affiliationRoleUpdate.PreviousAffiliation,
		Actor:    affiliationRoleUpdate.Actor,
	})

	// TODO: This functionality will be removed by `getMUCNotificationMessageFrom`
	// when it supports changing role and affiliation
	message := i18n.Localf("%s As a result, your role changed from %s to %s.", m,
		displayNameForRole(affiliationRoleUpdate.PreviousRole),
		displayNameForRole(affiliationRoleUpdate.NewRole))

	return message
}
