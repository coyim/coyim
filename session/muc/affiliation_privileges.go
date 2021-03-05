package muc

import (
	"github.com/coyim/coyim/session/muc/data"
)

const (
	enterOpenRoom privilege = iota
	registerWithOpenRoom
	retrieveMemberList
	enterMembersOnlyRoom
	banMembersAndUnaffiliatedUsers
	editMemberList
	assignAndRemoveModeratorRole
	editAdminList
	editOwnerList
	changeRoomConfiguration
	destroyRoom
)

func definedPrivilegesForAffiliations() map[string]*privileges {
	return map[string]*privileges{
		data.AffiliationOutcast: newPrivileges(),
		data.AffiliationNone: newPrivileges(
			enterOpenRoom,
			registerWithOpenRoom,
		),
		data.AffiliationMember: newPrivileges(
			enterOpenRoom,
			retrieveMemberList,
			enterMembersOnlyRoom,
		),
		data.AffiliationAdmin: newPrivileges(
			enterOpenRoom,
			retrieveMemberList,
			enterMembersOnlyRoom,
			banMembersAndUnaffiliatedUsers,
			editMemberList,
			assignAndRemoveModeratorRole,
		),
		data.AffiliationOwner: newPrivileges(
			enterOpenRoom,
			retrieveMemberList,
			enterMembersOnlyRoom,
			banMembersAndUnaffiliatedUsers,
			editMemberList,
			assignAndRemoveModeratorRole,
			editAdminList,
			editOwnerList,
			changeRoomConfiguration,
			destroyRoom,
		),
	}
}

func affiliationCan(privilege privilege, affiliation data.Affiliation) bool {
	affiliationPrivileges := definedPrivilegesForAffiliation(affiliation)
	return affiliationPrivileges.can(privilege)
}

func (o *Occupant) affiliationHasPrivilege(privilege privilege) bool {
	return affiliationCan(privilege, o.Affiliation)
}

// CanEnterOpenRoom returns true if the occupant can enter to an open room
// based on the occupant's affiliation
func (o *Occupant) CanEnterOpenRoom() bool {
	return o.affiliationHasPrivilege(enterOpenRoom)
}

// CanRegisterWithOpenRoom returns true if the occupant can register with open room
// based on the occupant's affiliation
func (o *Occupant) CanRegisterWithOpenRoom() bool {
	return o.affiliationHasPrivilege(registerWithOpenRoom)
}

// CanRetrieveMemberList returns true if the occupant can retrieve the members list
// based on the occupant's affiliation
func (o *Occupant) CanRetrieveMemberList() bool {
	return o.affiliationHasPrivilege(retrieveMemberList)
}

// CanEnterMembersOnlyRoom returns true if the occupant can enter to a members only room
// based on the occupant's affiliation
func (o *Occupant) CanEnterMembersOnlyRoom() bool {
	return o.affiliationHasPrivilege(enterMembersOnlyRoom)
}

// CanBanMembersAndUnaffiliatedUsers returns true if the occupant can ban members and unaffiliated users
// based on the occupant's affiliation
func (o *Occupant) CanBanMembersAndUnaffiliatedUsers() bool {
	return o.affiliationHasPrivilege(banMembersAndUnaffiliatedUsers)
}

// CanEditMemberList returns true if the occupant can edit the members list
// based on the occupant's affiliation
func (o *Occupant) CanEditMemberList() bool {
	return o.affiliationHasPrivilege(editMemberList)
}

// CanAssignAndRemoveModeratorRole returns true if the occupant can assign and remove moderator role
// based on the occupant's affiliation
func (o *Occupant) CanAssignAndRemoveModeratorRole() bool {
	return o.affiliationHasPrivilege(assignAndRemoveModeratorRole)
}

// CanEditAdminList returns true if the occupant can edit the admin list
// based on the occupant's affiliation
func (o *Occupant) CanEditAdminList() bool {
	return o.affiliationHasPrivilege(editAdminList)
}

// CanEditOwnerList returns true if the occupant can edit the owners list
// based on the occupant's affiliation
func (o *Occupant) CanEditOwnerList() bool {
	return o.affiliationHasPrivilege(editOwnerList)
}

// CanChangeRoomConfiguration returns true if the occupant can change the room configuration
// based on the occupant's affiliation
func (o *Occupant) CanChangeRoomConfiguration() bool {
	return o.affiliationHasPrivilege(changeRoomConfiguration)
}

// CanDestroyRoom returns true if the occupant can destroy the room
// based on the occupant's affiliation
func (o *Occupant) CanDestroyRoom() bool {
	return o.affiliationHasPrivilege(destroyRoom)
}

// CanChangeAffiliation returns a boolean indicating if the occupant can change the affiliation of the
// given occupant based on the occupant's affiliation
func (o *Occupant) CanChangeAffiliation(oc *Occupant) bool {
	if o.Affiliation.IsOwner() {
		return true
	}

	if o.Affiliation.IsAdmin() && (oc.Affiliation.IsOwner() || oc.Affiliation.IsAdmin()) {
		return false
	}

	return o.Affiliation.IsAdmin() && (oc.Affiliation.IsMember() || oc.Affiliation.IsNone() || oc.Affiliation.IsBanned())
}
