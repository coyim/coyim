package muc

import (
	"errors"

	"github.com/coyim/coyim/session/muc/data"
	"github.com/golang-collections/collections/set"
)

const (
	presentInRoom int = iota
	receiveMessages
	receiveOccupantPresence
	presenceToAllOccupants
	changeAvailabilityStatus
	changeRoomNickname
	sendPrivateMessages
	inviteOtherUsers
	sendMessagesToAll
	modifySubject
	kickParticipantsAndVisitors
	grantVoice
	revokeVoice
)

type rolesPrivileges struct {
	privileges *set.Set
}

func newRolesPrivileges(privileges ...int) *rolesPrivileges {
	p := &rolesPrivileges{
		privileges: set.New(),
	}

	for _, px := range privileges {
		p.privileges.Insert(px)
	}

	return p
}

func (p *rolesPrivileges) can(privilege int) bool {
	return p.privileges.Has(privilege)
}

var singletonRolesPrivileges map[string]*rolesPrivileges

func initializeRolePrivileges() {
	singletonRolesPrivileges = map[string]*rolesPrivileges{
		data.RoleNone: newRolesPrivileges(),
		data.RoleVisitor: newRolesPrivileges(
			presentInRoom,
			receiveMessages,
			receiveOccupantPresence,
			presenceToAllOccupants,
			changeAvailabilityStatus,
			changeRoomNickname,
			sendPrivateMessages,
			inviteOtherUsers,
		),
		data.RoleParticipant: newRolesPrivileges(
			presentInRoom,
			receiveMessages,
			receiveOccupantPresence,
			presenceToAllOccupants,
			changeAvailabilityStatus,
			changeRoomNickname,
			sendPrivateMessages,
			inviteOtherUsers,
			sendMessagesToAll,
			modifySubject,
		),
		data.RoleModerator: newRolesPrivileges(
			presentInRoom,
			receiveMessages,
			receiveOccupantPresence,
			presenceToAllOccupants,
			changeAvailabilityStatus,
			changeRoomNickname,
			sendPrivateMessages,
			inviteOtherUsers,
			sendMessagesToAll,
			modifySubject,
			kickParticipantsAndVisitors,
			grantVoice,
			revokeVoice,
		),
	}
}

func getRolePrivileges() map[string]*rolesPrivileges {
	if singletonRolesPrivileges == nil {
		initializeRolePrivileges()
	}

	return singletonRolesPrivileges
}

func getPrivilegesForRole(role data.Role) (*rolesPrivileges, error) {
	rp := getRolePrivileges()

	if p, ok := rp[role.Name()]; ok {
		return p, nil
	}

	return nil, errors.New("role not found in the privileges")
}

func roleCan(privilege int, role data.Role) bool {
	r, err := getPrivilegesForRole(role)
	if err != nil {
		return false
	}

	return r.can(privilege)
}

func (o *Occupant) roleHasPrivilege(privilege int) bool {
	return roleCan(privilege, o.Role)
}

// CanModifySubject returns a boolean indicating if the occupant can modify the room's subejct
func (o *Occupant) CanModifySubject() bool {
	return o.roleHasPrivilege(modifySubject)
}

// CanPresentInRoom description
func (o *Occupant) CanPresentInRoom() bool {
	return o.roleHasPrivilege(presentInRoom)
}

// CanReceiveMessage description
func (o *Occupant) CanReceiveMessage() bool {
	return o.roleHasPrivilege(receiveMessages)
}

// CanReceiveOccupantPresence description
func (o *Occupant) CanReceiveOccupantPresence() bool {
	return o.roleHasPrivilege(receiveOccupantPresence)
}

// CanBroadcastPresenceToAllOccupants description
func (o *Occupant) CanBroadcastPresenceToAllOccupants() bool {
	return o.roleHasPrivilege(presenceToAllOccupants)
}

// CanChangeAvailabilityStatus description
func (o *Occupant) CanChangeAvailabilityStatus() bool {
	return o.roleHasPrivilege(changeAvailabilityStatus)
}

// CanChangeRoomNickname description
func (o *Occupant) CanChangeRoomNickname() bool {
	return o.roleHasPrivilege(changeRoomNickname)
}

// CanSendPrivateMessages description
func (o *Occupant) CanSendPrivateMessages() bool {
	return o.roleHasPrivilege(sendPrivateMessages)
}

// CanInviteOtherUsers description
func (o *Occupant) CanInviteOtherUsers() bool {
	return o.roleHasPrivilege(sendPrivateMessages)
}

// CanSendMessagesToAll description
func (o *Occupant) CanSendMessagesToAll() bool {
	return o.roleHasPrivilege(sendMessagesToAll)
}

// CanKickParticipantsAndVisitors description
func (o *Occupant) CanKickParticipantsAndVisitors() bool {
	return o.roleHasPrivilege(kickParticipantsAndVisitors)
}

// CanGrantVoice description
func (o *Occupant) CanGrantVoice() bool {
	return o.roleHasPrivilege(grantVoice)
}

// CanRevokeVoice description
func (o *Occupant) CanRevokeVoice(oc *Occupant) bool {
	if oc.Affiliation.IsAdmin() || oc.Affiliation.IsOwner() {
		return false
	}

	if o.Role.IsModerator() {
		return true
	}

	return false
}
