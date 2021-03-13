package data

import (
	. "gopkg.in/check.v1"
)

type MucSuite struct{}

var _ = Suite(&MucSuite{})

func (s *MucSuite) Test_RoleFromString(c *C) {
	res, e := RoleFromString("none")
	c.Assert(e, IsNil)
	c.Assert(res, FitsTypeOf, &NoneRole{})

	res, e = RoleFromString("visitor")
	c.Assert(e, IsNil)
	c.Assert(res, FitsTypeOf, &VisitorRole{})

	res, e = RoleFromString("participant")
	c.Assert(e, IsNil)
	c.Assert(res, FitsTypeOf, &ParticipantRole{})

	res, e = RoleFromString("moderator")
	c.Assert(e, IsNil)
	c.Assert(res, FitsTypeOf, &ModeratorRole{})

	res, e = RoleFromString("")
	c.Assert(res, IsNil)
	c.Assert(e, ErrorMatches, "unknown role string: ''")

	res, e = RoleFromString("blabber")
	c.Assert(res, IsNil)
	c.Assert(e, ErrorMatches, "unknown role string: 'blabber'")
}

func (s *MucSuite) Test_Role_HasVoice(c *C) {
	c.Assert((&NoneRole{}).HasVoice(), Equals, false)
	c.Assert((&VisitorRole{}).HasVoice(), Equals, false)
	c.Assert((&ParticipantRole{}).HasVoice(), Equals, true)
	c.Assert((&ModeratorRole{}).HasVoice(), Equals, true)
}

func (s *MucSuite) Test_Role_WithVoice(c *C) {
	c.Assert((&NoneRole{}).WithVoice(), FitsTypeOf, &ParticipantRole{})
	c.Assert((&VisitorRole{}).WithVoice(), FitsTypeOf, &ParticipantRole{})
	c.Assert((&ParticipantRole{}).WithVoice(), FitsTypeOf, &ParticipantRole{})
	c.Assert((&ModeratorRole{}).WithVoice(), FitsTypeOf, &ModeratorRole{})
}

func (s *MucSuite) Test_Role_AsModerator(c *C) {
	c.Assert((&NoneRole{}).AsModerator(), FitsTypeOf, &ModeratorRole{})
	c.Assert((&VisitorRole{}).AsModerator(), FitsTypeOf, &ModeratorRole{})
	c.Assert((&ParticipantRole{}).AsModerator(), FitsTypeOf, &ModeratorRole{})
	c.Assert((&ModeratorRole{}).AsModerator(), FitsTypeOf, &ModeratorRole{})
}

func (s *MucSuite) Test_Role_Name(c *C) {
	c.Assert((&NoneRole{}).Name(), Equals, "none")
	c.Assert((&VisitorRole{}).Name(), Equals, "visitor")
	c.Assert((&ParticipantRole{}).Name(), Equals, "participant")
	c.Assert((&ModeratorRole{}).Name(), Equals, "moderator")
}

func (s *MucSuite) Test_Role_IsModerator(c *C) {
	c.Assert((&NoneRole{}).IsModerator(), Equals, false)
	c.Assert((&VisitorRole{}).IsModerator(), Equals, false)
	c.Assert((&ParticipantRole{}).IsModerator(), Equals, false)
	c.Assert((&ModeratorRole{}).IsModerator(), Equals, true)
}

func (s *MucSuite) Test_Role_IsParticipant(c *C) {
	c.Assert((&NoneRole{}).IsParticipant(), Equals, false)
	c.Assert((&VisitorRole{}).IsParticipant(), Equals, false)
	c.Assert((&ParticipantRole{}).IsParticipant(), Equals, true)
	c.Assert((&ModeratorRole{}).IsParticipant(), Equals, false)
}

func (s *MucSuite) Test_Role_IsVisitor(c *C) {
	c.Assert((&NoneRole{}).IsVisitor(), Equals, false)
	c.Assert((&VisitorRole{}).IsVisitor(), Equals, true)
	c.Assert((&ParticipantRole{}).IsVisitor(), Equals, false)
	c.Assert((&ModeratorRole{}).IsVisitor(), Equals, false)
}

func (s *MucSuite) Test_Role_IsNone(c *C) {
	c.Assert((&NoneRole{}).IsNone(), Equals, true)
	c.Assert((&VisitorRole{}).IsNone(), Equals, false)
	c.Assert((&ParticipantRole{}).IsNone(), Equals, false)
	c.Assert((&ModeratorRole{}).IsNone(), Equals, false)
}

func (s *MucSuite) Test_Role_IsDifferentFrom(c *C) {
	c.Assert((&NoneRole{}).IsDifferentFrom(&NoneRole{}), Equals, false)
	c.Assert((&NoneRole{}).IsDifferentFrom(&VisitorRole{}), Equals, true)
	c.Assert((&NoneRole{}).IsDifferentFrom(&ParticipantRole{}), Equals, true)
	c.Assert((&NoneRole{}).IsDifferentFrom(&ModeratorRole{}), Equals, true)

	c.Assert((&VisitorRole{}).IsDifferentFrom(&NoneRole{}), Equals, true)
	c.Assert((&VisitorRole{}).IsDifferentFrom(&VisitorRole{}), Equals, false)
	c.Assert((&VisitorRole{}).IsDifferentFrom(&ParticipantRole{}), Equals, true)
	c.Assert((&VisitorRole{}).IsDifferentFrom(&ModeratorRole{}), Equals, true)

	c.Assert((&ParticipantRole{}).IsDifferentFrom(&NoneRole{}), Equals, true)
	c.Assert((&ParticipantRole{}).IsDifferentFrom(&VisitorRole{}), Equals, true)
	c.Assert((&ParticipantRole{}).IsDifferentFrom(&ParticipantRole{}), Equals, false)
	c.Assert((&ParticipantRole{}).IsDifferentFrom(&ModeratorRole{}), Equals, true)

	c.Assert((&ModeratorRole{}).IsDifferentFrom(&NoneRole{}), Equals, true)
	c.Assert((&ModeratorRole{}).IsDifferentFrom(&VisitorRole{}), Equals, true)
	c.Assert((&ModeratorRole{}).IsDifferentFrom(&ParticipantRole{}), Equals, true)
	c.Assert((&ModeratorRole{}).IsDifferentFrom(&ModeratorRole{}), Equals, false)
}
