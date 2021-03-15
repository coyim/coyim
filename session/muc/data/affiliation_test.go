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

func (s *MucSuite) Test_Affiliation_IsLowerAffiliationThan(c *C) {
	c.Assert((&NoneAffiliation{}).IsLowerAffiliationThan(&OwnerAffiliation{}), Equals, true)
	c.Assert((&OutcastAffiliation{}).IsLowerAffiliationThan(&OwnerAffiliation{}), Equals, true)
	c.Assert((&MemberAffiliation{}).IsLowerAffiliationThan(&OwnerAffiliation{}), Equals, true)
	c.Assert((&AdminAffiliation{}).IsLowerAffiliationThan(&OwnerAffiliation{}), Equals, true)
	c.Assert((&OwnerAffiliation{}).IsLowerAffiliationThan(&OwnerAffiliation{}), Equals, false)
}

func (s *MucSuite) Test_Affiliation_IsDifferentFrom(c *C) {
	none := &NoneAffiliation{}
	outcast := &OutcastAffiliation{}
	member := &MemberAffiliation{}
	admin := &AdminAffiliation{}
	owner := &OwnerAffiliation{}

	c.Assert(none.IsDifferentFrom(none), Equals, false)
	c.Assert(none.IsDifferentFrom(outcast), Equals, true)
	c.Assert(none.IsDifferentFrom(member), Equals, true)
	c.Assert(none.IsDifferentFrom(admin), Equals, true)
	c.Assert(none.IsDifferentFrom(owner), Equals, true)

	c.Assert(outcast.IsDifferentFrom(none), Equals, true)
	c.Assert(outcast.IsDifferentFrom(outcast), Equals, false)
	c.Assert(outcast.IsDifferentFrom(member), Equals, true)
	c.Assert(outcast.IsDifferentFrom(admin), Equals, true)
	c.Assert(outcast.IsDifferentFrom(owner), Equals, true)

	c.Assert(member.IsDifferentFrom(none), Equals, true)
	c.Assert(member.IsDifferentFrom(outcast), Equals, true)
	c.Assert(member.IsDifferentFrom(member), Equals, false)
	c.Assert(member.IsDifferentFrom(admin), Equals, true)
	c.Assert(member.IsDifferentFrom(owner), Equals, true)

	c.Assert(admin.IsDifferentFrom(none), Equals, true)
	c.Assert(admin.IsDifferentFrom(outcast), Equals, true)
	c.Assert(admin.IsDifferentFrom(member), Equals, true)
	c.Assert(admin.IsDifferentFrom(admin), Equals, false)
	c.Assert(admin.IsDifferentFrom(owner), Equals, true)

	c.Assert(owner.IsDifferentFrom(none), Equals, true)
	c.Assert(owner.IsDifferentFrom(outcast), Equals, true)
	c.Assert(owner.IsDifferentFrom(member), Equals, true)
	c.Assert(owner.IsDifferentFrom(admin), Equals, true)
	c.Assert(owner.IsDifferentFrom(owner), Equals, false)
}
