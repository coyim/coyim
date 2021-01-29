package muc

import (
	"github.com/coyim/coyim/roster"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/coyim/xmpp/jid"
)

// Occupant contains information about a specific occupant in a specific room.
// This structure doesn't make sense without a connection to a room, since the information
// inside it depends on the room
type Occupant struct {
	// Nickname is the nickname of the person
	Nickname string

	// RealJid is the real JID of the person, if known. Otherwise it is nil
	RealJid jid.Full

	// Affiliation is the current affiliation of the occupant in the room
	Affiliation data.Affiliation

	// Role is the current role of the occupant in the room
	Role data.Role

	// Status contains the current status of the occupant in the room
	Status *roster.Status
}

// ChangeRoleToNone changes the role to the none role
func (o *Occupant) ChangeRoleToNone() {
	o.Role = &data.NoneRole{}
}

// ChangeRoleToVisitor changes the role to the visitor role
func (o *Occupant) ChangeRoleToVisitor() {
	o.Role = &data.VisitorRole{}
}

// ChangeRoleToParticipant changes the role to the participant role
func (o *Occupant) ChangeRoleToParticipant() {
	o.Role = &data.ParticipantRole{}
}

// ChangeRoleToModerator changes the role to the moderator role
func (o *Occupant) ChangeRoleToModerator() {
	o.Role = &data.ModeratorRole{}
}

// ChangeAffiliationToNone changes the affiliation to the none affiliation
func (o *Occupant) ChangeAffiliationToNone() {
	o.Affiliation = &data.NoneAffiliation{}
}

// Ban is a synonym for ChangeAffiliationToOutcast
func (o *Occupant) Ban() {
	o.ChangeAffiliationToOutcast()
}

// ChangeAffiliationToOutcast changes the affiliation to the outcast affiliation
func (o *Occupant) ChangeAffiliationToOutcast() {
	o.Affiliation = &data.OutcastAffiliation{}
}

// ChangeAffiliationToMember changes the affiliation to the member affiliation
func (o *Occupant) ChangeAffiliationToMember() {
	o.Affiliation = &data.MemberAffiliation{}
}

// ChangeAffiliationToAdmin changes the affiliation to the admin affiliation
func (o *Occupant) ChangeAffiliationToAdmin() {
	o.Affiliation = &data.AdminAffiliation{}
}

// ChangeAffiliationToOwner changes the affiliation to the owner affiliation
func (o *Occupant) ChangeAffiliationToOwner() {
	o.Affiliation = &data.OwnerAffiliation{}
}

// Update will update the information in this occupant object with the given information.
func (o *Occupant) Update(nickname string, affiliation data.Affiliation, role data.Role, status, statusMessage string, realJid jid.Full) {
	o.Nickname = nickname
	o.RealJid = realJid
	o.Role = role

	o.UpdateAffiliation(affiliation)
	o.UpdateStatus(status, statusMessage)
}

// UpdateAffiliation will update the occupant's affiliation
func (o *Occupant) UpdateAffiliation(a data.Affiliation) {
	switch a.(type) {
	case *data.OwnerAffiliation:
		o.ChangeAffiliationToOwner()
	case *data.AdminAffiliation:
		o.ChangeAffiliationToAdmin()
	case *data.MemberAffiliation:
		o.ChangeAffiliationToMember()
	case *data.OutcastAffiliation:
		o.ChangeAffiliationToOutcast()
	case *data.NoneAffiliation:
		o.ChangeAffiliationToNone()
	}
}

// UpdateRole will update the occupant's role
func (o *Occupant) UpdateRole(a data.Role) {
	switch a.(type) {
	case *data.ModeratorRole:
		o.ChangeRoleToModerator()
	case *data.ParticipantRole:
		o.ChangeRoleToParticipant()
	case *data.VisitorRole:
		o.ChangeRoleToVisitor()
	case *data.NoneRole:
		o.ChangeRoleToNone()
	}
}

// UpdateStatus will update the occupant's status
func (o *Occupant) UpdateStatus(status, statusMsg string) {
	o.Status = &roster.Status{Status: status, StatusMsg: statusMsg}
}

// HasVoice returns a boolean indicating if the occupant has voice based on the role
func (o *Occupant) HasVoice() bool {
	if o == nil || o.Role == nil {
		return false
	}

	return o.Role.HasVoice()
}
