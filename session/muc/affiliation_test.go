package muc

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"

	. "gopkg.in/check.v1"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func (s *MucSuite) Test_AffiliationFromString(c *C) {
	res, e := AffiliationFromString("none")
	c.Assert(e, IsNil)
	c.Assert(res, FitsTypeOf, &noneAffiliation{})

	res, e = AffiliationFromString("outcast")
	c.Assert(e, IsNil)
	c.Assert(res, FitsTypeOf, &outcastAffiliation{})

	res, e = AffiliationFromString("member")
	c.Assert(e, IsNil)
	c.Assert(res, FitsTypeOf, &memberAffiliation{})

	res, e = AffiliationFromString("admin")
	c.Assert(e, IsNil)
	c.Assert(res, FitsTypeOf, &adminAffiliation{})

	res, e = AffiliationFromString("owner")
	c.Assert(e, IsNil)
	c.Assert(res, FitsTypeOf, &ownerAffiliation{})

	res, e = AffiliationFromString("")
	c.Assert(res, IsNil)
	c.Assert(e, ErrorMatches, "unknown affiliation string: ''")

	res, e = AffiliationFromString("blabber")
	c.Assert(res, IsNil)
	c.Assert(e, ErrorMatches, "unknown affiliation string: 'blabber'")
}

func (s *MucSuite) Test_Affiliation_IsBanned(c *C) {
	c.Assert((&noneAffiliation{}).IsBanned(), Equals, false)
	c.Assert((&outcastAffiliation{}).IsBanned(), Equals, true)
	c.Assert((&memberAffiliation{}).IsBanned(), Equals, false)
	c.Assert((&adminAffiliation{}).IsBanned(), Equals, false)
	c.Assert((&ownerAffiliation{}).IsBanned(), Equals, false)
}

func (s *MucSuite) Test_Affiliation_IsMember(c *C) {
	c.Assert((&noneAffiliation{}).IsMember(), Equals, false)
	c.Assert((&outcastAffiliation{}).IsMember(), Equals, false)
	c.Assert((&memberAffiliation{}).IsMember(), Equals, true)
	c.Assert((&adminAffiliation{}).IsMember(), Equals, true)
	c.Assert((&ownerAffiliation{}).IsMember(), Equals, true)
}

func (s *MucSuite) Test_Affiliation_IsModerator(c *C) {
	c.Assert((&noneAffiliation{}).IsModerator(), Equals, false)
	c.Assert((&outcastAffiliation{}).IsModerator(), Equals, false)
	c.Assert((&memberAffiliation{}).IsModerator(), Equals, false)
	c.Assert((&adminAffiliation{}).IsModerator(), Equals, true)
	c.Assert((&ownerAffiliation{}).IsModerator(), Equals, true)
}
