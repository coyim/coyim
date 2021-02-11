package gui

import "github.com/coyim/coyim/session/muc/data"

func getDisplayForOccupantAffiliationRoleUpdate(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	d := newAffiliationRoleUpdateDisplayData(affiliationRoleUpdate)
	return displayAffiliationUpdateMessage(d)
}

type affiliationRoleUpdateDisplayData struct {
	*affiliationUpdateDisplayData
	newRole      data.Role
	previousRole data.Role
}

func newAffiliationRoleUpdateDisplayData(affiliationRoleUpdate data.AffiliationRoleUpdate) *affiliationRoleUpdateDisplayData {
	d := &affiliationRoleUpdateDisplayData{
		affiliationUpdateDisplayData: newAffiliationUpdateDisplayData(affiliationRoleUpdate.AffiliationUpdate),
		newRole:                      affiliationRoleUpdate.RoleUpdate.New,
		previousRole:                 affiliationRoleUpdate.RoleUpdate.Previous,
	}

	if affiliationRoleUpdate.Actor != nil {
		d.actor = affiliationRoleUpdate.Actor.Nickname
		d.actorAffiliation = affiliationRoleUpdate.Actor.Affiliation
	}

	return d
}
