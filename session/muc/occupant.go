package muc

import (
	"github.com/coyim/coyim/roster"
	"github.com/coyim/coyim/xmpp/jid"
)

// Occupant contains information about a specific occupant in a specific room.
// This structure doesn't make sense without a connection to a room, since the information
// inside it depends on the room
type Occupant struct {
	// TODO: We should change this member name to be consistent
	// Nick is the nickname of the person
	Nick string
	// TODO: Maybe change to RealJid to be consistent
	// Jid is the real JID of the person, if known. Otherwise it is nil
	Jid jid.Full

	// Affiliation is the current affiliation of the occupant in the room
	Affiliation Affiliation
	// Role is the current role of the occupant in the room
	Role Role

	// Status contains the current status of the occupant in the room
	Status *roster.Status
}

// ChangeRoleToNone changes the role to the none role
func (o *Occupant) ChangeRoleToNone() {
	o.Role = &noneRole{}
}

// ChangeRoleToVisitor changes the role to the visitor role
func (o *Occupant) ChangeRoleToVisitor() {
	o.Role = &visitorRole{}
}

// ChangeRoleToParticipant changes the role to the participant role
func (o *Occupant) ChangeRoleToParticipant() {
	o.Role = &participantRole{}
}

// ChangeRoleToModerator changes the role to the moderator role
func (o *Occupant) ChangeRoleToModerator() {
	o.Role = &moderatorRole{}
}

// ChangeAffiliationToNone changes the affiliation to the none affiliation
func (o *Occupant) ChangeAffiliationToNone() {
	o.Affiliation = &noneAffiliation{}
}

// Ban is a synonym for ChangeAffiliationToOutcast
func (o *Occupant) Ban() {
	o.ChangeAffiliationToOutcast()
}

// ChangeAffiliationToOutcast changes the affiliation to the outcast affiliation
func (o *Occupant) ChangeAffiliationToOutcast() {
	o.Affiliation = &outcastAffiliation{}
}

// ChangeAffiliationToMember changes the affiliation to the member affiliation
func (o *Occupant) ChangeAffiliationToMember() {
	o.Affiliation = &memberAffiliation{}
}

// ChangeAffiliationToAdmin changes the affiliation to the admin affiliation
func (o *Occupant) ChangeAffiliationToAdmin() {
	o.Affiliation = &adminAffiliation{}
}

// ChangeAffiliationToOwner changes the affiliation to the owner affiliation
func (o *Occupant) ChangeAffiliationToOwner() {
	o.Affiliation = &ownerAffiliation{}
}

// Update will update the information in this occupant object with the given information. It returns an error if the given affiliation or role doesn't match
// a known affiliation or role.
func (o *Occupant) Update(nickname string, affiliation Affiliation, role Role, status, statusMsg string, realJid jid.Full) {
	o.Nick = nickname
	o.Jid = realJid
	o.Affiliation = affiliation
	o.Role = role
	o.Status = &roster.Status{Status: status, StatusMsg: statusMsg}
}

// UpdateStatus will update the occupant's status
func (o *Occupant) UpdateStatus(status, statusMsg string) {
	o.Status = &roster.Status{Status: status, StatusMsg: statusMsg}
}
