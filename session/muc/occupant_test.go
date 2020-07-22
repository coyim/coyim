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

func (s *MucSuite) Test_Occupant_ChangeRoleToNone(c *C) {
	o := &Occupant{Role: &visitorRole{}}
	o.ChangeRoleToNone()
	c.Assert(o.Role, FitsTypeOf, &noneRole{})
}

func (s *MucSuite) Test_Occupant_ChangeRoleToVisitor(c *C) {
	o := &Occupant{Role: &noneRole{}}
	o.ChangeRoleToVisitor()
	c.Assert(o.Role, FitsTypeOf, &visitorRole{})
}

func (s *MucSuite) Test_Occupant_ChangeRoleToParticipant(c *C) {
	o := &Occupant{Role: &noneRole{}}
	o.ChangeRoleToParticipant()
	c.Assert(o.Role, FitsTypeOf, &participantRole{})
}

func (s *MucSuite) Test_Occupant_ChangeRoleToModerator(c *C) {
	o := &Occupant{Role: &noneRole{}}
	o.ChangeRoleToModerator()
	c.Assert(o.Role, FitsTypeOf, &moderatorRole{})
}

func (s *MucSuite) Test_Occupant_ChangeAffiliationToNone(c *C) {
	o := &Occupant{Affiliation: &outcastAffiliation{}}
	o.ChangeAffiliationToNone()
	c.Assert(o.Affiliation, FitsTypeOf, &noneAffiliation{})
}

func (s *MucSuite) Test_Occupant_Ban(c *C) {
	o := &Occupant{Affiliation: &noneAffiliation{}}
	o.Ban()
	c.Assert(o.Affiliation, FitsTypeOf, &outcastAffiliation{})
}

func (s *MucSuite) Test_Occupant_ChangeAffiliationToOutcast(c *C) {
	o := &Occupant{Affiliation: &noneAffiliation{}}
	o.ChangeAffiliationToOutcast()
	c.Assert(o.Affiliation, FitsTypeOf, &outcastAffiliation{})
}

func (s *MucSuite) Test_Occupant_ChangeAffiliationToMember(c *C) {
	o := &Occupant{Affiliation: &noneAffiliation{}}
	o.ChangeAffiliationToMember()
	c.Assert(o.Affiliation, FitsTypeOf, &memberAffiliation{})
}

func (s *MucSuite) Test_Occupant_ChangeAffiliationToAdmin(c *C) {
	o := &Occupant{Affiliation: &noneAffiliation{}}
	o.ChangeAffiliationToAdmin()
	c.Assert(o.Affiliation, FitsTypeOf, &adminAffiliation{})
}

func (s *MucSuite) Test_Occupant_ChangeAffiliationToOwner(c *C) {
	o := &Occupant{Affiliation: &noneAffiliation{}}
	o.ChangeAffiliationToOwner()
	c.Assert(o.Affiliation, FitsTypeOf, &ownerAffiliation{})
}

func (s *MucSuite) Test_Occupant_Update(c *C) {
	o := &Occupant{
		Nick:        "One",
		Jid:         jid.R("foo@bar.com/somewhere"),
		Affiliation: &memberAffiliation{},
		Role:        &moderatorRole{},
		Status:      roster.Status{"xa", "foo"},
	}

	e := o.Update(jid.R("room1@conf.example.org/Two"), "admin", "participant", "away", "here", nil)
	c.Assert(e, IsNil)

	c.Assert(o.Nick, Equals, "Two")
	c.Assert(o.Jid, IsNil)
	c.Assert(o.Affiliation, FitsTypeOf, &adminAffiliation{})
	c.Assert(o.Role, FitsTypeOf, &participantRole{})
	c.Assert(o.Status.Status, Equals, "away")
	c.Assert(o.Status.StatusMsg, Equals, "here")

	e = o.Update(jid.R("room1@conf.example.org/Two"), "admin2", "participant", "away", "here", nil)
	c.Assert(e, ErrorMatches, "unknown affiliation string: 'admin2'")

	e = o.Update(jid.R("room1@conf.example.org/Two"), "admin", "participant2", "away", "here", nil)
	c.Assert(e, ErrorMatches, "unknown role string: 'participant2'")
}
