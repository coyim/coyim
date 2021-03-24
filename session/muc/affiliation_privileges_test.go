package muc

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	. "gopkg.in/check.v1"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

type MucOccupantAffiliationPrivilegesSuite struct{}

var _ = Suite(&MucOccupantAffiliationPrivilegesSuite{})

func (*MucOccupantAffiliationPrivilegesSuite) Test_OccupantCanEnterOpenRoom(c *C) {
	o := &Occupant{}

	o.ChangeAffiliationToOutcast()
	c.Assert(o.CanEnterOpenRoom(), Equals, false)

	o.ChangeAffiliationToNone()
	c.Assert(o.CanEnterOpenRoom(), Equals, true)

	o.ChangeAffiliationToMember()
	c.Assert(o.CanEnterOpenRoom(), Equals, true)

	o.ChangeAffiliationToAdmin()
	c.Assert(o.CanEnterOpenRoom(), Equals, true)

	o.ChangeAffiliationToOwner()
	c.Assert(o.CanEnterOpenRoom(), Equals, true)
}

func (*MucOccupantAffiliationPrivilegesSuite) Test_OccupantCanRegisterWithOpenRoom(c *C) {
	o := &Occupant{}

	o.ChangeAffiliationToOutcast()
	c.Assert(o.CanRegisterWithOpenRoom(), Equals, false)

	o.ChangeAffiliationToNone()
	c.Assert(o.CanRegisterWithOpenRoom(), Equals, true)

	o.ChangeAffiliationToMember()
	c.Assert(o.CanRegisterWithOpenRoom(), Equals, false)

	o.ChangeAffiliationToAdmin()
	c.Assert(o.CanRegisterWithOpenRoom(), Equals, false)

	o.ChangeAffiliationToOwner()
	c.Assert(o.CanRegisterWithOpenRoom(), Equals, false)
}

func (*MucOccupantAffiliationPrivilegesSuite) Test_OccupantCanRetrieveMemberList(c *C) {
	o := &Occupant{}

	o.ChangeAffiliationToOutcast()
	c.Assert(o.CanRetrieveMemberList(), Equals, false)

	o.ChangeAffiliationToNone()
	c.Assert(o.CanRetrieveMemberList(), Equals, false)

	o.ChangeAffiliationToMember()
	c.Assert(o.CanRetrieveMemberList(), Equals, true)

	o.ChangeAffiliationToAdmin()
	c.Assert(o.CanRetrieveMemberList(), Equals, true)

	o.ChangeAffiliationToOwner()
	c.Assert(o.CanRetrieveMemberList(), Equals, true)
}

func (*MucOccupantAffiliationPrivilegesSuite) Test_OccupantCanEnterMembersOnlyRoom(c *C) {
	o := &Occupant{}

	o.ChangeAffiliationToOutcast()
	c.Assert(o.CanEnterMembersOnlyRoom(), Equals, false)

	o.ChangeAffiliationToNone()
	c.Assert(o.CanEnterMembersOnlyRoom(), Equals, false)

	o.ChangeAffiliationToMember()
	c.Assert(o.CanEnterMembersOnlyRoom(), Equals, true)

	o.ChangeAffiliationToAdmin()
	c.Assert(o.CanEnterMembersOnlyRoom(), Equals, true)

	o.ChangeAffiliationToOwner()
	c.Assert(o.CanEnterMembersOnlyRoom(), Equals, true)
}

func (*MucOccupantAffiliationPrivilegesSuite) Test_OccupantCanBanMembersAndUnaffiliatedUsers(c *C) {
	o := &Occupant{}

	o.ChangeAffiliationToOutcast()
	c.Assert(o.CanBanMembersAndUnaffiliatedUsers(), Equals, false)

	o.ChangeAffiliationToNone()
	c.Assert(o.CanBanMembersAndUnaffiliatedUsers(), Equals, false)

	o.ChangeAffiliationToMember()
	c.Assert(o.CanBanMembersAndUnaffiliatedUsers(), Equals, false)

	o.ChangeAffiliationToAdmin()
	c.Assert(o.CanBanMembersAndUnaffiliatedUsers(), Equals, true)

	o.ChangeAffiliationToOwner()
	c.Assert(o.CanBanMembersAndUnaffiliatedUsers(), Equals, true)
}

func (*MucOccupantAffiliationPrivilegesSuite) Test_OccupantCanEditMemberList(c *C) {
	o := &Occupant{}

	o.ChangeAffiliationToOutcast()
	c.Assert(o.CanEditMemberList(), Equals, false)

	o.ChangeAffiliationToNone()
	c.Assert(o.CanEditMemberList(), Equals, false)

	o.ChangeAffiliationToMember()
	c.Assert(o.CanEditMemberList(), Equals, false)

	o.ChangeAffiliationToAdmin()
	c.Assert(o.CanEditMemberList(), Equals, true)

	o.ChangeAffiliationToOwner()
	c.Assert(o.CanEditMemberList(), Equals, true)
}

func (*MucOccupantAffiliationPrivilegesSuite) Test_OccupantCanAssignAndRemoveModeratorRole(c *C) {
	o := &Occupant{}

	o.ChangeAffiliationToOutcast()
	c.Assert(o.CanAssignAndRemoveModeratorRole(), Equals, false)

	o.ChangeAffiliationToNone()
	c.Assert(o.CanAssignAndRemoveModeratorRole(), Equals, false)

	o.ChangeAffiliationToMember()
	c.Assert(o.CanAssignAndRemoveModeratorRole(), Equals, false)

	o.ChangeAffiliationToAdmin()
	c.Assert(o.CanAssignAndRemoveModeratorRole(), Equals, true)

	o.ChangeAffiliationToOwner()
	c.Assert(o.CanAssignAndRemoveModeratorRole(), Equals, true)
}

func (*MucOccupantAffiliationPrivilegesSuite) Test_OccupantCanEditAdminList(c *C) {
	o := &Occupant{}

	o.ChangeAffiliationToOutcast()
	c.Assert(o.CanEditAdminList(), Equals, false)

	o.ChangeAffiliationToNone()
	c.Assert(o.CanEditAdminList(), Equals, false)

	o.ChangeAffiliationToMember()
	c.Assert(o.CanEditAdminList(), Equals, false)

	o.ChangeAffiliationToAdmin()
	c.Assert(o.CanEditAdminList(), Equals, false)

	o.ChangeAffiliationToOwner()
	c.Assert(o.CanEditAdminList(), Equals, true)
}

func (*MucOccupantAffiliationPrivilegesSuite) Test_OccupantCanEditOwnerList(c *C) {
	o := &Occupant{}

	o.ChangeAffiliationToOutcast()
	c.Assert(o.CanEditOwnerList(), Equals, false)

	o.ChangeAffiliationToNone()
	c.Assert(o.CanEditOwnerList(), Equals, false)

	o.ChangeAffiliationToMember()
	c.Assert(o.CanEditOwnerList(), Equals, false)

	o.ChangeAffiliationToAdmin()
	c.Assert(o.CanEditOwnerList(), Equals, false)

	o.ChangeAffiliationToOwner()
	c.Assert(o.CanEditOwnerList(), Equals, true)
}

func (*MucOccupantAffiliationPrivilegesSuite) Test_OccupantCanChangeRoomConfiguration(c *C) {
	o := &Occupant{}

	o.ChangeAffiliationToOutcast()
	c.Assert(o.CanChangeRoomConfiguration(), Equals, false)

	o.ChangeAffiliationToNone()
	c.Assert(o.CanChangeRoomConfiguration(), Equals, false)

	o.ChangeAffiliationToMember()
	c.Assert(o.CanChangeRoomConfiguration(), Equals, false)

	o.ChangeAffiliationToAdmin()
	c.Assert(o.CanChangeRoomConfiguration(), Equals, false)

	o.ChangeAffiliationToOwner()
	c.Assert(o.CanChangeRoomConfiguration(), Equals, true)
}

func (*MucOccupantAffiliationPrivilegesSuite) Test_OccupantCanDestroyRoom(c *C) {
	o := &Occupant{}

	o.ChangeAffiliationToOutcast()
	c.Assert(o.CanDestroyRoom(), Equals, false)

	o.ChangeAffiliationToNone()
	c.Assert(o.CanDestroyRoom(), Equals, false)

	o.ChangeAffiliationToMember()
	c.Assert(o.CanDestroyRoom(), Equals, false)

	o.ChangeAffiliationToAdmin()
	c.Assert(o.CanDestroyRoom(), Equals, false)

	o.ChangeAffiliationToOwner()
	c.Assert(o.CanDestroyRoom(), Equals, true)
}

func (*MucOccupantAffiliationPrivilegesSuite) Test_OccupantCanChangeAffiliation(c *C) {
	o := &Occupant{}
	o.ChangeAffiliationToNone()

	oc := &Occupant{}
	oc.ChangeAffiliationToNone()

	c.Assert(o.CanChangeAffiliation(oc), Equals, false)

	o.ChangeAffiliationToOwner()
	c.Assert(o.CanChangeAffiliation(oc), Equals, true)

	o.ChangeAffiliationToAdmin()
	c.Assert(o.CanChangeAffiliation(oc), Equals, true)

	o.ChangeAffiliationToMember()
	c.Assert(o.CanChangeAffiliation(oc), Equals, false)

	o.ChangeAffiliationToNone()
	c.Assert(o.CanChangeAffiliation(oc), Equals, false)

	o.ChangeAffiliationToOutcast()
	c.Assert(o.CanChangeAffiliation(oc), Equals, false)
}

func (*MucOccupantAffiliationPrivilegesSuite) Test_OccupantCanChangeAffiliation_OfNoAffiliation(c *C) {
	o := &Occupant{}

	oc := &Occupant{}
	oc.ChangeAffiliationToNone()

	o.ChangeAffiliationToNone()
	c.Assert(o.CanChangeAffiliation(oc), Equals, false)

	o.ChangeAffiliationToMember()
	c.Assert(o.CanChangeAffiliation(oc), Equals, false)

	o.ChangeAffiliationToAdmin()
	c.Assert(o.CanChangeAffiliation(oc), Equals, true)

	o.ChangeAffiliationToOwner()
	c.Assert(o.CanChangeAffiliation(oc), Equals, true)
}

func (*MucOccupantAffiliationPrivilegesSuite) Test_OccupantCanChangeAffiliation_OfAMember(c *C) {
	o := &Occupant{}

	oc := &Occupant{}
	oc.ChangeAffiliationToMember()

	o.ChangeAffiliationToNone()
	c.Assert(o.CanChangeAffiliation(oc), Equals, false)

	o.ChangeAffiliationToMember()
	c.Assert(o.CanChangeAffiliation(oc), Equals, false)

	o.ChangeAffiliationToAdmin()
	c.Assert(o.CanChangeAffiliation(oc), Equals, true)

	o.ChangeAffiliationToOwner()
	c.Assert(o.CanChangeAffiliation(oc), Equals, true)
}

func (*MucOccupantAffiliationPrivilegesSuite) Test_OccupantCanChangeAffiliation_OfAnAdmin(c *C) {
	o := &Occupant{}

	oc := &Occupant{}
	oc.ChangeAffiliationToAdmin()

	o.ChangeAffiliationToNone()
	c.Assert(o.CanChangeAffiliation(oc), Equals, false)

	o.ChangeAffiliationToMember()
	c.Assert(o.CanChangeAffiliation(oc), Equals, false)

	o.ChangeAffiliationToAdmin()
	c.Assert(o.CanChangeAffiliation(oc), Equals, false)

	o.ChangeAffiliationToOwner()
	c.Assert(o.CanChangeAffiliation(oc), Equals, true)
}

func (*MucOccupantAffiliationPrivilegesSuite) Test_OccupantCanChangeAffiliation_OfAnOwner(c *C) {
	o := &Occupant{}

	oc := &Occupant{}
	oc.ChangeAffiliationToOwner()

	o.ChangeAffiliationToNone()
	c.Assert(o.CanChangeAffiliation(oc), Equals, false)

	o.ChangeAffiliationToMember()
	c.Assert(o.CanChangeAffiliation(oc), Equals, false)

	o.ChangeAffiliationToAdmin()
	c.Assert(o.CanChangeAffiliation(oc), Equals, false)

	o.ChangeAffiliationToOwner()
	c.Assert(o.CanChangeAffiliation(oc), Equals, true)
}

func (*MucOccupantAffiliationPrivilegesSuite) Test_Occupant_isOwner(c *C) {
	o := &Occupant{}

	o.ChangeAffiliationToOwner()
	c.Assert(o.isOwner(), Equals, true)

	o.ChangeAffiliationToAdmin()
	c.Assert(o.isOwner(), Equals, false)

	o.ChangeAffiliationToMember()
	c.Assert(o.isOwner(), Equals, false)

	o.ChangeAffiliationToOutcast()
	c.Assert(o.isOwner(), Equals, false)

	o.ChangeAffiliationToNone()
	c.Assert(o.isOwner(), Equals, false)
}

func (*MucOccupantAffiliationPrivilegesSuite) Test_Occupant_notOwner(c *C) {
	o := &Occupant{}

	o.ChangeAffiliationToOwner()
	c.Assert(o.isNotOwner(), Equals, false)

	o.ChangeAffiliationToAdmin()
	c.Assert(o.isNotOwner(), Equals, true)

	o.ChangeAffiliationToMember()
	c.Assert(o.isNotOwner(), Equals, true)

	o.ChangeAffiliationToOutcast()
	c.Assert(o.isNotOwner(), Equals, true)

	o.ChangeAffiliationToNone()
	c.Assert(o.isNotOwner(), Equals, true)
}

func (*MucOccupantAffiliationPrivilegesSuite) Test_Occupant_isAdmin(c *C) {
	o := &Occupant{}

	o.ChangeAffiliationToOwner()
	c.Assert(o.isAdmin(), Equals, false)

	o.ChangeAffiliationToAdmin()
	c.Assert(o.isAdmin(), Equals, true)

	o.ChangeAffiliationToMember()
	c.Assert(o.isAdmin(), Equals, false)

	o.ChangeAffiliationToOutcast()
	c.Assert(o.isAdmin(), Equals, false)

	o.ChangeAffiliationToNone()
	c.Assert(o.isAdmin(), Equals, false)
}

func (*MucOccupantAffiliationPrivilegesSuite) Test_Occupant_notAnAdmin(c *C) {
	o := &Occupant{}

	o.ChangeAffiliationToOwner()
	c.Assert(o.isNotAdmin(), Equals, true)

	o.ChangeAffiliationToAdmin()
	c.Assert(o.isNotAdmin(), Equals, false)

	o.ChangeAffiliationToMember()
	c.Assert(o.isNotAdmin(), Equals, true)

	o.ChangeAffiliationToOutcast()
	c.Assert(o.isNotAdmin(), Equals, true)

	o.ChangeAffiliationToNone()
	c.Assert(o.isNotAdmin(), Equals, true)
}

func (*MucOccupantAffiliationPrivilegesSuite) Test_Occupant_isOwnerOrAdmin(c *C) {
	o := &Occupant{}

	o.ChangeAffiliationToOwner()
	c.Assert(o.isOwnerOrAdmin(), Equals, true)

	o.ChangeAffiliationToAdmin()
	c.Assert(o.isOwnerOrAdmin(), Equals, true)

	o.ChangeAffiliationToMember()
	c.Assert(o.isOwnerOrAdmin(), Equals, false)

	o.ChangeAffiliationToOutcast()
	c.Assert(o.isOwnerOrAdmin(), Equals, false)

	o.ChangeAffiliationToNone()
	c.Assert(o.isOwnerOrAdmin(), Equals, false)
}

func (*MucOccupantAffiliationPrivilegesSuite) Test_Occupant_notOwnerOrAdmin(c *C) {
	o := &Occupant{}

	o.ChangeAffiliationToOwner()
	c.Assert(o.isNotOwnerOrAdmin(), Equals, false)

	o.ChangeAffiliationToAdmin()
	c.Assert(o.isNotOwnerOrAdmin(), Equals, false)

	o.ChangeAffiliationToMember()
	c.Assert(o.isNotOwnerOrAdmin(), Equals, true)

	o.ChangeAffiliationToOutcast()
	c.Assert(o.isNotOwnerOrAdmin(), Equals, true)

	o.ChangeAffiliationToNone()
	c.Assert(o.isNotOwnerOrAdmin(), Equals, true)
}
