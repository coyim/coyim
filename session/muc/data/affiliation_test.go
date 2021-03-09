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

func (s *MucSuite) Test_Affiliation_IsNone(c *C) {
	c.Assert((&NoneAffiliation{}).IsNone(), Equals, true)
	c.Assert((&OutcastAffiliation{}).IsNone(), Equals, false)
	c.Assert((&MemberAffiliation{}).IsNone(), Equals, false)
	c.Assert((&AdminAffiliation{}).IsNone(), Equals, false)
	c.Assert((&OwnerAffiliation{}).IsNone(), Equals, false)
}

func (s *MucSuite) Test_Affiliation_IsBanned(c *C) {
	c.Assert((&NoneAffiliation{}).IsBanned(), Equals, false)
	c.Assert((&OutcastAffiliation{}).IsBanned(), Equals, true)
	c.Assert((&MemberAffiliation{}).IsBanned(), Equals, false)
	c.Assert((&AdminAffiliation{}).IsBanned(), Equals, false)
	c.Assert((&OwnerAffiliation{}).IsBanned(), Equals, false)
}

func (s *MucSuite) Test_Affiliation_IsOutcast(c *C) {
	c.Assert((&NoneAffiliation{}).IsOutcast(), Equals, false)
	c.Assert((&OutcastAffiliation{}).IsOutcast(), Equals, true)
	c.Assert((&MemberAffiliation{}).IsOutcast(), Equals, false)
	c.Assert((&AdminAffiliation{}).IsOutcast(), Equals, false)
	c.Assert((&OwnerAffiliation{}).IsOutcast(), Equals, false)
}

func (s *MucSuite) Test_Affiliation_IsAdmin(c *C) {
	c.Assert((&NoneAffiliation{}).IsAdmin(), Equals, false)
	c.Assert((&OutcastAffiliation{}).IsAdmin(), Equals, false)
	c.Assert((&MemberAffiliation{}).IsAdmin(), Equals, false)
	c.Assert((&AdminAffiliation{}).IsAdmin(), Equals, true)
	c.Assert((&OwnerAffiliation{}).IsAdmin(), Equals, false)
}

func (s *MucSuite) Test_Affiliation_IsOwner(c *C) {
	c.Assert((&NoneAffiliation{}).IsOwner(), Equals, false)
	c.Assert((&OutcastAffiliation{}).IsOwner(), Equals, false)
	c.Assert((&MemberAffiliation{}).IsOwner(), Equals, false)
	c.Assert((&AdminAffiliation{}).IsOwner(), Equals, false)
	c.Assert((&OwnerAffiliation{}).IsOwner(), Equals, true)
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

func (s *MucSuite) Test_Affiliation_Name(c *C) {
	c.Assert((&NoneAffiliation{}).Name(), Equals, "none")
	c.Assert((&OutcastAffiliation{}).Name(), Equals, "outcast")
	c.Assert((&MemberAffiliation{}).Name(), Equals, "member")
	c.Assert((&AdminAffiliation{}).Name(), Equals, "admin")
	c.Assert((&OwnerAffiliation{}).Name(), Equals, "owner")
}

func (s *MucSuite) Test_Affiliation_IsLowerThan(c *C) {
	c.Assert((&NoneAffiliation{}).IsLowerThan(&OwnerAffiliation{}), Equals, true)
	c.Assert((&OutcastAffiliation{}).IsLowerThan(&OwnerAffiliation{}), Equals, true)
	c.Assert((&MemberAffiliation{}).IsLowerThan(&OwnerAffiliation{}), Equals, true)
	c.Assert((&AdminAffiliation{}).IsLowerThan(&OwnerAffiliation{}), Equals, true)
	c.Assert((&OwnerAffiliation{}).IsLowerThan(&OwnerAffiliation{}), Equals, false)
}

func (s *MucSuite) Test_AreAffiliationsDifferent(c *C) {
	c.Assert(AreAffiliationsDifferent(&NoneAffiliation{}, &NoneAffiliation{}), Equals, false)
	c.Assert(AreAffiliationsDifferent(&NoneAffiliation{}, &OutcastAffiliation{}), Equals, true)
	c.Assert(AreAffiliationsDifferent(&NoneAffiliation{}, &MemberAffiliation{}), Equals, true)
	c.Assert(AreAffiliationsDifferent(&NoneAffiliation{}, &AdminAffiliation{}), Equals, true)
	c.Assert(AreAffiliationsDifferent(&NoneAffiliation{}, &OwnerAffiliation{}), Equals, true)

	c.Assert(AreAffiliationsDifferent(&OutcastAffiliation{}, &NoneAffiliation{}), Equals, true)
	c.Assert(AreAffiliationsDifferent(&OutcastAffiliation{}, &OutcastAffiliation{}), Equals, false)
	c.Assert(AreAffiliationsDifferent(&OutcastAffiliation{}, &MemberAffiliation{}), Equals, true)
	c.Assert(AreAffiliationsDifferent(&OutcastAffiliation{}, &AdminAffiliation{}), Equals, true)
	c.Assert(AreAffiliationsDifferent(&OutcastAffiliation{}, &OwnerAffiliation{}), Equals, true)

	c.Assert(AreAffiliationsDifferent(&MemberAffiliation{}, &NoneAffiliation{}), Equals, true)
	c.Assert(AreAffiliationsDifferent(&MemberAffiliation{}, &OutcastAffiliation{}), Equals, true)
	c.Assert(AreAffiliationsDifferent(&MemberAffiliation{}, &MemberAffiliation{}), Equals, false)
	c.Assert(AreAffiliationsDifferent(&MemberAffiliation{}, &AdminAffiliation{}), Equals, true)
	c.Assert(AreAffiliationsDifferent(&MemberAffiliation{}, &OwnerAffiliation{}), Equals, true)

	c.Assert(AreAffiliationsDifferent(&AdminAffiliation{}, &NoneAffiliation{}), Equals, true)
	c.Assert(AreAffiliationsDifferent(&AdminAffiliation{}, &OutcastAffiliation{}), Equals, true)
	c.Assert(AreAffiliationsDifferent(&AdminAffiliation{}, &MemberAffiliation{}), Equals, true)
	c.Assert(AreAffiliationsDifferent(&AdminAffiliation{}, &AdminAffiliation{}), Equals, false)
	c.Assert(AreAffiliationsDifferent(&AdminAffiliation{}, &OwnerAffiliation{}), Equals, true)

	c.Assert(AreAffiliationsDifferent(&OwnerAffiliation{}, &NoneAffiliation{}), Equals, true)
	c.Assert(AreAffiliationsDifferent(&OwnerAffiliation{}, &OutcastAffiliation{}), Equals, true)
	c.Assert(AreAffiliationsDifferent(&OwnerAffiliation{}, &MemberAffiliation{}), Equals, true)
	c.Assert(AreAffiliationsDifferent(&OwnerAffiliation{}, &AdminAffiliation{}), Equals, true)
	c.Assert(AreAffiliationsDifferent(&OwnerAffiliation{}, &OwnerAffiliation{}), Equals, false)
}
