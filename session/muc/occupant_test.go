package muc

import (
	"io/ioutil"

	"github.com/coyim/coyim/roster"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"

	. "gopkg.in/check.v1"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func newOccupantForTest(affiliation Affiliation, role Role) *Occupant {
	return &Occupant{
		Affiliation: affiliation,
		Role:        role,
	}
}

func (s *MucSuite) Test_Occupant_ChangeRoleToNone(c *C) {
	o := newOccupantForTest(nil, &visitorRole{})
	o.ChangeRoleToNone()
	c.Assert(o.Role, FitsTypeOf, &noneRole{})
}

func (s *MucSuite) Test_Occupant_ChangeRoleToVisitor(c *C) {
	o := newOccupantForTest(nil, &noneRole{})
	o.ChangeRoleToVisitor()
	c.Assert(o.Role, FitsTypeOf, &visitorRole{})
}

func (s *MucSuite) Test_Occupant_ChangeRoleToParticipant(c *C) {
	o := newOccupantForTest(nil, &noneRole{})
	o.ChangeRoleToParticipant()
	c.Assert(o.Role, FitsTypeOf, &participantRole{})
}

func (s *MucSuite) Test_Occupant_ChangeRoleToModerator(c *C) {
	o := newOccupantForTest(nil, &noneRole{})
	o.ChangeRoleToModerator()
	c.Assert(o.Role, FitsTypeOf, &moderatorRole{})
}

func (s *MucSuite) Test_Occupant_ChangeAffiliationToNone(c *C) {
	o := newOccupantForTest(&outcastAffiliation{}, nil)
	o.ChangeAffiliationToNone()
	c.Assert(o.Affiliation, FitsTypeOf, &noneAffiliation{})
}

func (s *MucSuite) Test_Occupant_Ban(c *C) {
	o := newOccupantForTest(&noneAffiliation{}, nil)
	o.Ban()
	c.Assert(o.Affiliation, FitsTypeOf, &outcastAffiliation{})
}

func (s *MucSuite) Test_Occupant_ChangeAffiliationToOutcast(c *C) {
	o := newOccupantForTest(&noneAffiliation{}, nil)
	o.ChangeAffiliationToOutcast()
	c.Assert(o.Affiliation, FitsTypeOf, &outcastAffiliation{})
}

func (s *MucSuite) Test_Occupant_ChangeAffiliationToMember(c *C) {
	o := newOccupantForTest(&noneAffiliation{}, nil)
	o.ChangeAffiliationToMember()
	c.Assert(o.Affiliation, FitsTypeOf, &memberAffiliation{})
}

func (s *MucSuite) Test_Occupant_ChangeAffiliationToAdmin(c *C) {
	o := newOccupantForTest(&noneAffiliation{}, nil)
	o.ChangeAffiliationToAdmin()
	c.Assert(o.Affiliation, FitsTypeOf, &adminAffiliation{})
}

func (s *MucSuite) Test_Occupant_ChangeAffiliationToOwner(c *C) {
	o := newOccupantForTest(&noneAffiliation{}, nil)
	o.ChangeAffiliationToOwner()
	c.Assert(o.Affiliation, FitsTypeOf, &ownerAffiliation{})
}

func (s *MucSuite) Test_Occupant_Update(c *C) {
	o := &Occupant{
		Nick:        "One",
		Jid:         jid.ParseFull("foo@bar.com/somewhere"),
		Affiliation: &memberAffiliation{},
		Role:        &moderatorRole{},
		Status:      roster.Status{Status: "xa", StatusMsg: "foo"},
	}

	o.Update("Two", &adminAffiliation{}, &participantRole{}, "away", "here", nil)

	c.Assert(o.Nick, Equals, "Two")
	c.Assert(o.Jid, IsNil)
	c.Assert(o.Affiliation, FitsTypeOf, &adminAffiliation{})
	c.Assert(o.Role, FitsTypeOf, &participantRole{})
	c.Assert(o.Status.Status, Equals, "away")
	c.Assert(o.Status.StatusMsg, Equals, "here")
}
