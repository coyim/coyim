package muc

import (
	"io/ioutil"

	"github.com/coyim/coyim/session/muc/data"
	log "github.com/sirupsen/logrus"
	. "gopkg.in/check.v1"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

type MucOccupantRoleAffiliationPrivilegesSuite struct{}

var _ = Suite(&MucOccupantRoleAffiliationPrivilegesSuite{})

type canKickOccupantTest struct {
	actor    *Occupant
	occupant *Occupant
	expected bool
}

func newTestOccupant(affiliation data.Affiliation, role data.Role) *Occupant {
	return &Occupant{
		Affiliation: affiliation,
		Role:        role,
	}
}

func (*MucOccupantRoleAffiliationPrivilegesSuite) Test_RoleModeratorAffiliationNone_CanKickAnOccupant(c *C) {
	testCases := []canKickOccupantTest{
		// Actor: Moderator - None
		{
			actor:    newTestOccupant(&data.NoneAffiliation{}, &data.ModeratorRole{}),
			occupant: newTestOccupant(&data.NoneAffiliation{}, &data.ModeratorRole{}),
			expected: false,
		},
		{
			actor:    newTestOccupant(&data.NoneAffiliation{}, &data.ModeratorRole{}),
			occupant: newTestOccupant(&data.MemberAffiliation{}, &data.ModeratorRole{}),
			expected: false,
		},
		{
			actor:    newTestOccupant(&data.NoneAffiliation{}, &data.ModeratorRole{}),
			occupant: newTestOccupant(&data.AdminAffiliation{}, &data.ModeratorRole{}),
			expected: false,
		},
		{
			actor:    newTestOccupant(&data.NoneAffiliation{}, &data.ModeratorRole{}),
			occupant: newTestOccupant(&data.OwnerAffiliation{}, &data.ModeratorRole{}),
			expected: false,
		},

		{
			actor:    newTestOccupant(&data.NoneAffiliation{}, &data.ModeratorRole{}),
			occupant: newTestOccupant(&data.NoneAffiliation{}, &data.ParticipantRole{}),
			expected: false,
		},
		{
			actor:    newTestOccupant(&data.NoneAffiliation{}, &data.ModeratorRole{}),
			occupant: newTestOccupant(&data.MemberAffiliation{}, &data.ParticipantRole{}),
			expected: false,
		},
		{
			actor:    newTestOccupant(&data.NoneAffiliation{}, &data.ModeratorRole{}),
			occupant: newTestOccupant(&data.AdminAffiliation{}, &data.ParticipantRole{}),
			expected: false,
		},
		{
			actor:    newTestOccupant(&data.NoneAffiliation{}, &data.ModeratorRole{}),
			occupant: newTestOccupant(&data.OwnerAffiliation{}, &data.ParticipantRole{}),
			expected: false,
		},

		{
			actor:    newTestOccupant(&data.NoneAffiliation{}, &data.ModeratorRole{}),
			occupant: newTestOccupant(&data.NoneAffiliation{}, &data.VisitorRole{}),
			expected: false,
		},
		{
			actor:    newTestOccupant(&data.NoneAffiliation{}, &data.ModeratorRole{}),
			occupant: newTestOccupant(&data.MemberAffiliation{}, &data.VisitorRole{}),
			expected: false,
		},
		{
			actor:    newTestOccupant(&data.NoneAffiliation{}, &data.ModeratorRole{}),
			occupant: newTestOccupant(&data.AdminAffiliation{}, &data.VisitorRole{}),
			expected: false,
		},
		{
			actor:    newTestOccupant(&data.NoneAffiliation{}, &data.ModeratorRole{}),
			occupant: newTestOccupant(&data.OwnerAffiliation{}, &data.VisitorRole{}),
			expected: false,
		},

		{
			actor:    newTestOccupant(&data.NoneAffiliation{}, &data.ModeratorRole{}),
			occupant: newTestOccupant(&data.NoneAffiliation{}, &data.NoneRole{}),
			expected: false,
		},
		{
			actor:    newTestOccupant(&data.NoneAffiliation{}, &data.ModeratorRole{}),
			occupant: newTestOccupant(&data.MemberAffiliation{}, &data.NoneRole{}),
			expected: false,
		},
		{
			actor:    newTestOccupant(&data.NoneAffiliation{}, &data.ModeratorRole{}),
			occupant: newTestOccupant(&data.AdminAffiliation{}, &data.NoneRole{}),
			expected: false,
		},
		{
			actor:    newTestOccupant(&data.NoneAffiliation{}, &data.ModeratorRole{}),
			occupant: newTestOccupant(&data.OwnerAffiliation{}, &data.NoneRole{}),
			expected: false,
		},
	}

	for _, scenario := range testCases {
		c.Assert(scenario.actor.CanKickOccupant(scenario.occupant), Equals, scenario.expected)
	}
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
