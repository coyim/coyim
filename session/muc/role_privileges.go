package muc

import "github.com/coyim/coyim/session/muc/data"

type roleNumberType int

const (
	roleTypeNone roleNumberType = iota
	roleTypeVisitor
	roleTypeParticipant
	roleTypeModerator
)

func roleNumberTypeFrom(r data.Role) roleNumberType {
	switch {
	case r.IsModerator():
		return roleTypeModerator
	case r.IsParticipant():
		return roleTypeParticipant
	case r.IsVisitor():
		return roleTypeVisitor
	}
	return roleTypeNone
}

type privilege int

const (
	presentInRoom privilege = iota
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

var rolePrivileges = [][]bool{
	{false /*none*/, true /*visitor*/, true /*participant*/, true /*moderator*/},   //presentInRoom
	{false /*none*/, true /*visitor*/, true /*participant*/, true /*moderator*/},   //receiveMessages
	{false /*none*/, true /*visitor*/, true /*participant*/, true /*moderator*/},   //receiveOccupantPresence
	{false /*none*/, true /*visitor*/, true /*participant*/, true /*moderator*/},   //presenceToAllOccupants
	{false /*none*/, true /*visitor*/, true /*participant*/, true /*moderator*/},   //changeAvailabilityStatus
	{false /*none*/, true /*visitor*/, true /*participant*/, true /*moderator*/},   //changeRoomNickname
	{false /*none*/, true /*visitor*/, true /*participant*/, true /*moderator*/},   //sendPrivateMessages
	{false /*none*/, true /*visitor*/, true /*participant*/, true /*moderator*/},   //inviteOtherUsers
	{false /*none*/, false /*visitor*/, true /*participant*/, true /*moderator*/},  //sendMessagesToAll
	{false /*none*/, false /*visitor*/, true /*participant*/, true /*moderator*/},  //modifySubject
	{false /*none*/, false /*visitor*/, false /*participant*/, true /*moderator*/}, //kickParticipantsAndVisitors
	{false /*none*/, false /*visitor*/, false /*participant*/, true /*moderator*/}, //grantVoice
	{false /*none*/, false /*visitor*/, false /*participant*/, true /*moderator*/}, //revokeVoice
}

func (o *Occupant) roleHasPrivilege(p privilege) bool {
	return rolePrivileges[p][roleNumberTypeFrom(o.Role)]
}

// CanPresentInRoom returns a boolean indicating if the occupant can be present in the room
// based on the occupant's role
func (o *Occupant) CanPresentInRoom() bool {
	return o.roleHasPrivilege(presentInRoom)
}

// CanReceiveMessage returns a boolean indicating if the occupant can receive messages
// based on the occupant's role
func (o *Occupant) CanReceiveMessage() bool {
	return o.roleHasPrivilege(receiveMessages)
}

// CanReceiveOccupantPresence returns a boolean indicating if the occupant can receive occupant presence
// based on the occupant's role
func (o *Occupant) CanReceiveOccupantPresence() bool {
	return o.roleHasPrivilege(receiveOccupantPresence)
}

// CanBroadcastPresenceToAllOccupants returns a boolean indicating if the occupant can send a broadcast
// presence to all occupants based on the occupant's role
func (o *Occupant) CanBroadcastPresenceToAllOccupants() bool {
	return o.roleHasPrivilege(presenceToAllOccupants)
}

// CanChangeAvailabilityStatus returns a boolean indicating if the occupant can change their availability status
// based on the occupant's role
func (o *Occupant) CanChangeAvailabilityStatus() bool {
	return o.roleHasPrivilege(changeAvailabilityStatus)
}

// CanChangeRoomNickname returns a boolean indicating if the occupant can change the room nickname
// based on the occupant's role
func (o *Occupant) CanChangeRoomNickname() bool {
	return o.roleHasPrivilege(changeRoomNickname)
}

// CanSendPrivateMessages returns a boolean indicating if the occupant can send private messages
// based on the occupant's role
func (o *Occupant) CanSendPrivateMessages() bool {
	return o.roleHasPrivilege(sendPrivateMessages)
}

// CanInviteOtherUsers returns a boolean indicating if the occupant can invite other users
// based on the occupant's role
func (o *Occupant) CanInviteOtherUsers() bool {
	return o.roleHasPrivilege(sendPrivateMessages)
}

// CanSendMessagesToAll returns a boolean indicating if the occupant can send messages to all
// based on the occupant's role
func (o *Occupant) CanSendMessagesToAll() bool {
	return o.roleHasPrivilege(sendMessagesToAll)
}

// CanModifySubject returns a boolean indicating if the occupant can modify the room's subject
// based on the occupant's role
func (o *Occupant) CanModifySubject() bool {
	return o.roleHasPrivilege(modifySubject)
}

// CanKickParticipantsAndVisitors returns a boolean indicating if the occupant can kick participants and
// visitors
// based on the occupant's role
func (o *Occupant) CanKickParticipantsAndVisitors() bool {
	return o.roleHasPrivilege(kickParticipantsAndVisitors)
}

// CanGrantVoice returns a boolean indicating if the occupant can grant voice
// based on the occupant's role
func (o *Occupant) CanGrantVoice() bool {
	return o.roleHasPrivilege(grantVoice)
}

// CanRevokeVoice returns a boolean indicating if the occupant can revoke voice
// based on the occupant's role
func (o *Occupant) CanRevokeVoice(oc *Occupant) bool {
	if oc.Affiliation.IsAdmin() || oc.Affiliation.IsOwner() {
		return false
	}

	if o.Role.IsModerator() {
		return true
	}

	return false
}

// CanChangeRole returns a boolean indicating if the occupant can change the role of the
// given occupant based on the occupant's role and affiliation
func (o *Occupant) CanChangeRole(oc *Occupant) bool {
	return (o.isOwner() || o.isAdmin()) && (!oc.isOwner() && !oc.isAdmin())
}

// CanKickOccupant returns a boolean indicating if the occupant can kick another occupant
// based on the occupant's role
func (o *Occupant) CanKickOccupant(oc *Occupant) bool {
	if !o.roleHasPrivilege(kickParticipantsAndVisitors) || !(oc.Role.IsParticipant() || oc.Role.IsVisitor()) {
		return false
	}

	if o.isOwner() {
		return !oc.isOwner()
	}

	return !oc.isOwner() && !oc.isAdmin()
}
