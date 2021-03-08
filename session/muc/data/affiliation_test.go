package data

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
	c.Assert(res, FitsTypeOf, &NoneAffiliation{})

	res, e = AffiliationFromString("outcast")
	c.Assert(e, IsNil)
	c.Assert(res, FitsTypeOf, &OutcastAffiliation{})

	res, e = AffiliationFromString("member")
	c.Assert(e, IsNil)
	c.Assert(res, FitsTypeOf, &MemberAffiliation{})

	res, e = AffiliationFromString("admin")
	c.Assert(e, IsNil)
	c.Assert(res, FitsTypeOf, &AdminAffiliation{})

	res, e = AffiliationFromString("owner")
	c.Assert(e, IsNil)
	c.Assert(res, FitsTypeOf, &OwnerAffiliation{})

	res, e = AffiliationFromString("")
	c.Assert(res, IsNil)
	c.Assert(e, ErrorMatches, "unknown affiliation string: ''")

	res, e = AffiliationFromString("blabber")
	c.Assert(res, IsNil)
	c.Assert(e, ErrorMatches, "unknown affiliation string: 'blabber'")
}

func (s *MucSuite) Test_Affiliation_IsBanned(c *C) {
	c.Assert((&NoneAffiliation{}).IsBanned(), Equals, false)
	c.Assert((&OutcastAffiliation{}).IsBanned(), Equals, true)
	c.Assert((&MemberAffiliation{}).IsBanned(), Equals, false)
	c.Assert((&AdminAffiliation{}).IsBanned(), Equals, false)
	c.Assert((&OwnerAffiliation{}).IsBanned(), Equals, false)
}

func (s *MucSuite) Test_Affiliation_IsMember(c *C) {
	c.Assert((&NoneAffiliation{}).IsMember(), Equals, false)
	c.Assert((&OutcastAffiliation{}).IsMember(), Equals, false)
	c.Assert((&MemberAffiliation{}).IsMember(), Equals, true)
	c.Assert((&AdminAffiliation{}).IsMember(), Equals, false)
	c.Assert((&OwnerAffiliation{}).IsMember(), Equals, false)
}

func (s *MucSuite) Test_Affiliation_IsModerator(c *C) {
	c.Assert((&NoneAffiliation{}).IsModerator(), Equals, false)
	c.Assert((&OutcastAffiliation{}).IsModerator(), Equals, false)
	c.Assert((&MemberAffiliation{}).IsModerator(), Equals, false)
	c.Assert((&AdminAffiliation{}).IsModerator(), Equals, true)
	c.Assert((&OwnerAffiliation{}).IsModerator(), Equals, true)
}
