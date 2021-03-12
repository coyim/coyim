package muc

import (
	"github.com/coyim/coyim/roster"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/coyim/xmpp/jid"

	. "gopkg.in/check.v1"
)

func newOccupantForTest(affiliation data.Affiliation, role data.Role) *Occupant {
	return &Occupant{
		Affiliation: affiliation,
		Role:        role,
	}
}

func (s *MucSuite) Test_Occupant_ChangeRoleToNone(c *C) {
	o := newOccupantForTest(nil, &data.VisitorRole{})
	o.ChangeRoleToNone()
	c.Assert(o.Role, FitsTypeOf, &data.NoneRole{})
}

func (s *MucSuite) Test_Occupant_ChangeRoleToVisitor(c *C) {
	o := newOccupantForTest(nil, &data.NoneRole{})
	o.ChangeRoleToVisitor()
	c.Assert(o.Role, FitsTypeOf, &data.VisitorRole{})
}

func (s *MucSuite) Test_Occupant_ChangeRoleToParticipant(c *C) {
	o := newOccupantForTest(nil, &data.NoneRole{})
	o.ChangeRoleToParticipant()
	c.Assert(o.Role, FitsTypeOf, &data.ParticipantRole{})
}

func (s *MucSuite) Test_Occupant_ChangeRoleToModerator(c *C) {
	o := newOccupantForTest(nil, &data.NoneRole{})
	o.ChangeRoleToModerator()
	c.Assert(o.Role, FitsTypeOf, &data.ModeratorRole{})
}

func (s *MucSuite) Test_Occupant_ChangeAffiliationToNone(c *C) {
	o := newOccupantForTest(&data.OutcastAffiliation{}, nil)
	o.ChangeAffiliationToNone()
	c.Assert(o.Affiliation, FitsTypeOf, &data.NoneAffiliation{})
}

func (s *MucSuite) Test_Occupant_Ban(c *C) {
	o := newOccupantForTest(&data.NoneAffiliation{}, nil)
	o.Ban()
	c.Assert(o.Affiliation, FitsTypeOf, &data.OutcastAffiliation{})
}

func (s *MucSuite) Test_Occupant_ChangeAffiliationToOutcast(c *C) {
	o := newOccupantForTest(&data.NoneAffiliation{}, nil)
	o.ChangeAffiliationToOutcast()
	c.Assert(o.Affiliation, FitsTypeOf, &data.OutcastAffiliation{})
}

func (s *MucSuite) Test_Occupant_ChangeAffiliationToMember(c *C) {
	o := newOccupantForTest(&data.NoneAffiliation{}, nil)
	o.ChangeAffiliationToMember()
	c.Assert(o.Affiliation, FitsTypeOf, &data.MemberAffiliation{})
}

func (s *MucSuite) Test_Occupant_ChangeAffiliationToAdmin(c *C) {
	o := newOccupantForTest(&data.NoneAffiliation{}, nil)
	o.ChangeAffiliationToAdmin()
	c.Assert(o.Affiliation, FitsTypeOf, &data.AdminAffiliation{})
}

func (s *MucSuite) Test_Occupant_ChangeAffiliationToOwner(c *C) {
	o := newOccupantForTest(&data.NoneAffiliation{}, nil)
	o.ChangeAffiliationToOwner()
	c.Assert(o.Affiliation, FitsTypeOf, &data.OwnerAffiliation{})
}

func (s *MucSuite) Test_Occupant_Update(c *C) {
	o := &Occupant{
		Nickname:    "One",
		RealJid:     jid.ParseFull("foo@bar.com/somewhere"),
		Affiliation: &data.MemberAffiliation{},
		Role:        &data.ModeratorRole{},
		Status:      &roster.Status{Status: "xa", StatusMsg: "foo"},
	}

	o.Update("Two", &data.AdminAffiliation{}, &data.ParticipantRole{}, "away", "here", nil)

	c.Assert(o.Nickname, Equals, "Two")
	c.Assert(o.RealJid, IsNil)
	c.Assert(o.Affiliation, FitsTypeOf, &data.AdminAffiliation{})
	c.Assert(o.Role, FitsTypeOf, &data.ParticipantRole{})
	c.Assert(o.Status.Status, Equals, "away")
	c.Assert(o.Status.StatusMsg, Equals, "here")
}

func (s *MucSuite) Test_Occupant_UpdateAffilition(c *C) {
	o := newOccupantForTest(&data.NoneAffiliation{}, nil)

	o.UpdateAffiliation(&data.OwnerAffiliation{})
	c.Assert(o.Affiliation.IsOwner(), Equals, true)

	o.UpdateAffiliation(&data.AdminAffiliation{})
	c.Assert(o.Affiliation.IsAdmin(), Equals, true)

	o.UpdateAffiliation(&data.MemberAffiliation{})
	c.Assert(o.Affiliation.IsMember(), Equals, true)

	o.UpdateAffiliation(&data.OutcastAffiliation{})
	c.Assert(o.Affiliation.IsBanned(), Equals, true)

	o.UpdateAffiliation(&data.NoneAffiliation{})
	c.Assert(o.Affiliation.IsNone(), Equals, true)
}

func (s *MucSuite) Test_Occupant_UpdateRole(c *C) {
	o := newOccupantForTest(nil, &data.NoneRole{})

	o.UpdateRole(&data.ModeratorRole{})
	c.Assert(o.Role.IsModerator(), Equals, true)

	o.UpdateRole(&data.ParticipantRole{})
	c.Assert(o.Role.IsParticipant(), Equals, true)

	o.UpdateRole(&data.VisitorRole{})
	c.Assert(o.Role.IsVisitor(), Equals, true)

	o.UpdateRole(&data.NoneRole{})
	c.Assert(o.Role.IsNone(), Equals, true)
}

func (s *MucSuite) Test_Occupant_HasVoice(c *C) {
	var o *Occupant

	c.Assert(o.HasVoice(), Equals, false)

	o = &Occupant{}
	c.Assert(o.HasVoice(), Equals, false)

	o.Role = &data.ParticipantRole{}
	c.Assert(o.HasVoice(), Equals, true)
}
