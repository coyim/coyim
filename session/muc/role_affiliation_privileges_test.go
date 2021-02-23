package muc

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	. "gopkg.in/check.v1"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

type MucOccupantRoleAffiliationPrivilegesSuite struct{}

var _ = Suite(&MucOccupantRoleAffiliationPrivilegesSuite{})

func (*MucOccupantRoleAffiliationPrivilegesSuite) Test_RoleModeratorAffiliationNone_CanKickAnOccupant(c *C) {
	actor := &Occupant{}
	actor.ChangeRoleToModerator()
	actor.ChangeAffiliationToNone()

	oc := &Occupant{}

	// Testing role moderator
	oc.ChangeRoleToModerator()
	oc.ChangeAffiliationToNone()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToMember()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToAdmin()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToOwner()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	// Testing role participant
	oc.ChangeRoleToParticipant()
	oc.ChangeAffiliationToNone()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToMember()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToAdmin()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToOwner()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	// Testing role visitor
	oc.ChangeRoleToVisitor()
	oc.ChangeAffiliationToNone()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToMember()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToAdmin()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToOwner()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)
}

func (*MucOccupantRoleAffiliationPrivilegesSuite) Test_RoleModeratorAffiliationMember_CanKickAnOccupant(c *C) {
	actor := &Occupant{}
	actor.ChangeRoleToModerator()
	actor.ChangeAffiliationToMember()

	oc := &Occupant{}

	// Testing role moderator
	oc.ChangeRoleToModerator()
	oc.ChangeAffiliationToNone()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToMember()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToAdmin()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToOwner()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	// Testing role participant
	oc.ChangeRoleToParticipant()
	oc.ChangeAffiliationToNone()
	c.Assert(actor.CanKickOccupant(oc), Equals, true)

	oc.ChangeAffiliationToMember()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToAdmin()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToOwner()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	// Testing role visitor
	oc.ChangeRoleToVisitor()
	oc.ChangeAffiliationToNone()
	c.Assert(actor.CanKickOccupant(oc), Equals, true)

	oc.ChangeAffiliationToMember()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToAdmin()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToOwner()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)
}

func (*MucOccupantRoleAffiliationPrivilegesSuite) Test_RoleModeratorAffiliationAdmin_CanKickAnOccupant(c *C) {
	actor := &Occupant{}
	actor.ChangeRoleToModerator()
	actor.ChangeAffiliationToAdmin()

	oc := &Occupant{}

	// Testing role moderator
	oc.ChangeRoleToModerator()
	oc.ChangeAffiliationToNone()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToMember()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToAdmin()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToOwner()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	// Testing role participant
	oc.ChangeRoleToParticipant()
	oc.ChangeAffiliationToNone()
	c.Assert(actor.CanKickOccupant(oc), Equals, true)

	oc.ChangeAffiliationToMember()
	c.Assert(actor.CanKickOccupant(oc), Equals, true)

	oc.ChangeAffiliationToAdmin()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToOwner()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	// Testing role visitor
	oc.ChangeRoleToVisitor()
	oc.ChangeAffiliationToNone()
	c.Assert(actor.CanKickOccupant(oc), Equals, true)

	oc.ChangeAffiliationToMember()
	c.Assert(actor.CanKickOccupant(oc), Equals, true)

	oc.ChangeAffiliationToAdmin()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToOwner()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)
}

func (*MucOccupantRoleAffiliationPrivilegesSuite) Test_RoleModeratorAffiliationOwner_CanKickAnOccupant(c *C) {
	actor := &Occupant{}
	actor.ChangeRoleToModerator()
	actor.ChangeAffiliationToOwner()

	oc := &Occupant{}

	// Testing role moderator
	oc.ChangeRoleToModerator()
	oc.ChangeAffiliationToNone()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToMember()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToAdmin()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToOwner()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	// Testing role participant
	oc.ChangeRoleToParticipant()
	oc.ChangeAffiliationToNone()
	c.Assert(actor.CanKickOccupant(oc), Equals, true)

	oc.ChangeAffiliationToMember()
	c.Assert(actor.CanKickOccupant(oc), Equals, true)

	oc.ChangeAffiliationToAdmin()
	c.Assert(actor.CanKickOccupant(oc), Equals, true)

	oc.ChangeAffiliationToOwner()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	// Testing role visitor
	oc.ChangeRoleToVisitor()
	oc.ChangeAffiliationToNone()
	c.Assert(actor.CanKickOccupant(oc), Equals, true)

	oc.ChangeAffiliationToMember()
	c.Assert(actor.CanKickOccupant(oc), Equals, true)

	oc.ChangeAffiliationToAdmin()
	c.Assert(actor.CanKickOccupant(oc), Equals, true)

	oc.ChangeAffiliationToOwner()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)
}

func (*MucOccupantRoleAffiliationPrivilegesSuite) Test_RoleParticipantAffiliationNone_CanKickAnOccupant(c *C) {
	actor := &Occupant{}
	actor.ChangeRoleToParticipant()
	actor.ChangeAffiliationToNone()

	oc := &Occupant{}

	// Testing role moderator
	oc.ChangeRoleToModerator()
	oc.ChangeAffiliationToNone()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToMember()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToAdmin()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToOwner()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	// Testing role participant
	oc.ChangeRoleToParticipant()
	oc.ChangeAffiliationToNone()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToMember()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToAdmin()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToOwner()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	// Testing role visitor
	oc.ChangeRoleToVisitor()
	oc.ChangeAffiliationToNone()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToMember()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToAdmin()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToOwner()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)
}

func (*MucOccupantRoleAffiliationPrivilegesSuite) Test_RoleParticipantAffiliationMember_CanKickAnOccupant(c *C) {
	actor := &Occupant{}
	actor.ChangeRoleToParticipant()
	actor.ChangeAffiliationToMember()

	oc := &Occupant{}

	// Testing role moderator
	oc.ChangeRoleToModerator()
	oc.ChangeAffiliationToNone()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToMember()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToAdmin()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToOwner()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	// Testing role participant
	oc.ChangeRoleToParticipant()
	oc.ChangeAffiliationToNone()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToMember()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToAdmin()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToOwner()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	// Testing role visitor
	oc.ChangeRoleToVisitor()
	oc.ChangeAffiliationToNone()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToMember()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToAdmin()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToOwner()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)
}

func (*MucOccupantRoleAffiliationPrivilegesSuite) Test_RoleParticipantAffiliationAdmin_CanKickAnOccupant(c *C) {
	actor := &Occupant{}
	actor.ChangeRoleToParticipant()
	actor.ChangeAffiliationToAdmin()

	oc := &Occupant{}

	// Testing role moderator
	oc.ChangeRoleToModerator()
	oc.ChangeAffiliationToNone()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToMember()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToAdmin()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToOwner()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	// Testing role participant
	oc.ChangeRoleToParticipant()
	oc.ChangeAffiliationToNone()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToMember()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToAdmin()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToOwner()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	// Testing role visitor
	oc.ChangeRoleToVisitor()
	oc.ChangeAffiliationToNone()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToMember()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToAdmin()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToOwner()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)
}

func (*MucOccupantRoleAffiliationPrivilegesSuite) Test_RoleParticipantAffiliationOwner_CanKickAnOccupant(c *C) {
	actor := &Occupant{}
	actor.ChangeRoleToParticipant()
	actor.ChangeAffiliationToOwner()

	oc := &Occupant{}

	// Testing role moderator
	oc.ChangeRoleToModerator()
	oc.ChangeAffiliationToNone()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToMember()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToAdmin()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToOwner()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	// Testing role participant
	oc.ChangeRoleToParticipant()
	oc.ChangeAffiliationToNone()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToMember()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToAdmin()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToOwner()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	// Testing role visitor
	oc.ChangeRoleToVisitor()
	oc.ChangeAffiliationToNone()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToMember()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToAdmin()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToOwner()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)
}

func (*MucOccupantRoleAffiliationPrivilegesSuite) Test_RoleVisitorAffiliationNone_CanKickAnOccupant(c *C) {
	actor := &Occupant{}
	actor.ChangeRoleToVisitor()
	actor.ChangeAffiliationToNone()

	oc := &Occupant{}

	// Testing role moderator
	oc.ChangeRoleToModerator()
	oc.ChangeAffiliationToNone()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToMember()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToAdmin()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToOwner()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	// Testing role participant
	oc.ChangeRoleToParticipant()
	oc.ChangeAffiliationToNone()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToMember()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToAdmin()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToOwner()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	// Testing role visitor
	oc.ChangeRoleToVisitor()
	oc.ChangeAffiliationToNone()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToMember()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToAdmin()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToOwner()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)
}

func (*MucOccupantRoleAffiliationPrivilegesSuite) Test_RoleVisitorAffiliationMember_CanKickAnOccupant(c *C) {
	actor := &Occupant{}
	actor.ChangeRoleToVisitor()
	actor.ChangeAffiliationToMember()

	oc := &Occupant{}

	// Testing role moderator
	oc.ChangeRoleToModerator()
	oc.ChangeAffiliationToNone()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToMember()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToAdmin()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToOwner()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	// Testing role participant
	oc.ChangeRoleToParticipant()
	oc.ChangeAffiliationToNone()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToMember()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToAdmin()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToOwner()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	// Testing role visitor
	oc.ChangeRoleToVisitor()
	oc.ChangeAffiliationToNone()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToMember()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToAdmin()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToOwner()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)
}

func (*MucOccupantRoleAffiliationPrivilegesSuite) Test_RoleVisitorAffiliationAdmin_CanKickAnOccupant(c *C) {
	actor := &Occupant{}
	actor.ChangeRoleToVisitor()
	actor.ChangeAffiliationToAdmin()

	oc := &Occupant{}

	// Testing role moderator
	oc.ChangeRoleToModerator()
	oc.ChangeAffiliationToNone()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToMember()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToAdmin()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToOwner()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	// Testing role participant
	oc.ChangeRoleToParticipant()
	oc.ChangeAffiliationToNone()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToMember()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToAdmin()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToOwner()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	// Testing role visitor
	oc.ChangeRoleToVisitor()
	oc.ChangeAffiliationToNone()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToMember()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToAdmin()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToOwner()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)
}

func (*MucOccupantRoleAffiliationPrivilegesSuite) Test_RoleVisitorAffiliationOwner_CanKickAnOccupant(c *C) {
	actor := &Occupant{}
	actor.ChangeRoleToVisitor()
	actor.ChangeAffiliationToOwner()

	oc := &Occupant{}

	// Testing role moderator
	oc.ChangeRoleToModerator()
	oc.ChangeAffiliationToNone()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToMember()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToAdmin()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToOwner()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	// Testing role participant
	oc.ChangeRoleToParticipant()
	oc.ChangeAffiliationToNone()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToMember()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToAdmin()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToOwner()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	// Testing role visitor
	oc.ChangeRoleToVisitor()
	oc.ChangeAffiliationToNone()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToMember()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToAdmin()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)

	oc.ChangeAffiliationToOwner()
	c.Assert(actor.CanKickOccupant(oc), Equals, false)
}
