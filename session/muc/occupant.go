package muc

import "github.com/coyim/coyim/xmpp/jid"
import "github.com/coyim/coyim/roster"

type Occupant struct {
	Nick string
	Jid  jid.WithResource

	Affiliation Affiliation
	Role        Role

	Status roster.Status
}

func (o *Occupant) ChangeRoleToNone() {
	o.Role = &noneRole{}
}

func (o *Occupant) ChangeRoleToVisitor() {
	o.Role = &visitorRole{}
}

func (o *Occupant) ChangeRoleToParticipant() {
	o.Role = &participantRole{}
}

func (o *Occupant) ChangeRoleToModerator() {
	o.Role = &moderatorRole{}
}

func (o *Occupant) ChangeAffiliationToNone() {
	o.Affiliation = &noneAffiliation{}
}

func (o *Occupant) Ban() {
	o.ChangeAffiliationToOutcast()
}

func (o *Occupant) ChangeAffiliationToOutcast() {
	o.Affiliation = &outcastAffiliation{}
}

func (o *Occupant) ChangeAffiliationToMember() {
	o.Affiliation = &memberAffiliation{}
}

func (o *Occupant) ChangeAffiliationToAdmin() {
	o.Affiliation = &adminAffiliation{}
}

func (o *Occupant) ChangeAffiliationToOwner() {
	o.Affiliation = &ownerAffiliation{}
}

func (o *Occupant) Update(from jid.WithResource, affiliation, role, show, statusMsg string, realJid jid.WithResource) error {
	var err error

	o.Nick = string(from.Resource())
	o.Jid = realJid
	o.Affiliation, err = AffiliationFromString(affiliation)
	if err != nil {
		return err
	}
	o.Role, err = RoleFromString(role)
	if err != nil {
		return err
	}

	o.Status = roster.Status{show, statusMsg}

	return nil
}
