package muc

import (
	"io/ioutil"
	"testing"

	log "github.com/sirupsen/logrus"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

func init() {
	log.SetOutput(ioutil.Discard)
}

type MucSuite struct{}

var _ = Suite(&MucSuite{})

func (s *MucSuite) Test_RoleFromString(c *C) {
	res, e := RoleFromString("none")
	c.Assert(e, IsNil)
	c.Assert(res, FitsTypeOf, &noneRole{})

	res, e = RoleFromString("visitor")
	c.Assert(e, IsNil)
	c.Assert(res, FitsTypeOf, &visitorRole{})

	res, e = RoleFromString("participant")
	c.Assert(e, IsNil)
	c.Assert(res, FitsTypeOf, &participantRole{})

	res, e = RoleFromString("moderator")
	c.Assert(e, IsNil)
	c.Assert(res, FitsTypeOf, &moderatorRole{})

	res, e = RoleFromString("")
	c.Assert(res, IsNil)
	c.Assert(e, ErrorMatches, "unknown role string: ''")

	res, e = RoleFromString("blabber")
	c.Assert(res, IsNil)
	c.Assert(e, ErrorMatches, "unknown role string: 'blabber'")
}

func (s *MucSuite) Test_Role_HasVoice(c *C) {
	c.Assert((&noneRole{}).HasVoice(), Equals, false)
	c.Assert((&visitorRole{}).HasVoice(), Equals, false)
	c.Assert((&participantRole{}).HasVoice(), Equals, true)
	c.Assert((&moderatorRole{}).HasVoice(), Equals, true)
}

func (s *MucSuite) Test_Role_WithVoice(c *C) {
	c.Assert((&noneRole{}).WithVoice(), FitsTypeOf, &participantRole{})
	c.Assert((&visitorRole{}).WithVoice(), FitsTypeOf, &participantRole{})
	c.Assert((&participantRole{}).WithVoice(), FitsTypeOf, &participantRole{})
	c.Assert((&moderatorRole{}).WithVoice(), FitsTypeOf, &moderatorRole{})
}

func (s *MucSuite) Test_Role_AsModerator(c *C) {
	c.Assert((&noneRole{}).AsModerator(), FitsTypeOf, &moderatorRole{})
	c.Assert((&visitorRole{}).AsModerator(), FitsTypeOf, &moderatorRole{})
	c.Assert((&participantRole{}).AsModerator(), FitsTypeOf, &moderatorRole{})
	c.Assert((&moderatorRole{}).AsModerator(), FitsTypeOf, &moderatorRole{})
}
