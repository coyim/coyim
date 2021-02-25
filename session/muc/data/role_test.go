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
