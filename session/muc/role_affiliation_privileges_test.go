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
	occupantAffiliation data.Affiliation
	occupantRole        data.Role
	expected            bool
}

func newTestOccupant(affiliation data.Affiliation, role data.Role) *Occupant {
	return &Occupant{
		Affiliation: affiliation,
		Role:        role,
	}
}

func (*MucOccupantRoleAffiliationPrivilegesSuite) Test_RoleModeratorAffiliationNone_CanKickAnOccupant(c *C) {
	testCases := []canKickOccupantTest{
		// Occupant: ModeratorRole
		{&data.NoneAffiliation{}, &data.ModeratorRole{}, false},
		{&data.MemberAffiliation{}, &data.ModeratorRole{}, false},
		{&data.AdminAffiliation{}, &data.ModeratorRole{}, false},
		{&data.OwnerAffiliation{}, &data.ModeratorRole{}, false},

		// Occupant: ParticipantRole
		{&data.NoneAffiliation{}, &data.ParticipantRole{}, false},
		{&data.MemberAffiliation{}, &data.ParticipantRole{}, false},
		{&data.AdminAffiliation{}, &data.ParticipantRole{}, false},
		{&data.OwnerAffiliation{}, &data.ParticipantRole{}, false},

		// Occupant: VisitorRole
		{&data.NoneAffiliation{}, &data.VisitorRole{}, false},
		{&data.MemberAffiliation{}, &data.VisitorRole{}, false},
		{&data.AdminAffiliation{}, &data.VisitorRole{}, false},
		{&data.OwnerAffiliation{}, &data.VisitorRole{}, false},

		// Occupant: NoneRole
		{&data.NoneAffiliation{}, &data.NoneRole{}, false},
		{&data.MemberAffiliation{}, &data.NoneRole{}, false},
		{&data.AdminAffiliation{}, &data.NoneRole{}, false},
		{&data.OwnerAffiliation{}, &data.NoneRole{}, false},
	}

	// Actor: NoneAffiliation - ModeratorRole
	actor := newTestOccupant(&data.NoneAffiliation{}, &data.ModeratorRole{})
	for _, scenario := range testCases {
		c.Assert(actor.CanKickOccupant(newTestOccupant(scenario.occupantAffiliation, scenario.occupantRole)), Equals, scenario.expected)
	}
}

func (*MucOccupantRoleAffiliationPrivilegesSuite) Test_RoleModeratorAffiliationMember_CanKickAnOccupant(c *C) {
	testCases := []canKickOccupantTest{
		// Occupant: ModeratorRole
		{&data.NoneAffiliation{}, &data.ModeratorRole{}, false},
		{&data.MemberAffiliation{}, &data.ModeratorRole{}, false},
		{&data.AdminAffiliation{}, &data.ModeratorRole{}, false},
		{&data.OwnerAffiliation{}, &data.ModeratorRole{}, false},

		// Occupant: ParticipantRole
		{&data.NoneAffiliation{}, &data.ParticipantRole{}, true},
		{&data.MemberAffiliation{}, &data.ParticipantRole{}, false},
		{&data.AdminAffiliation{}, &data.ParticipantRole{}, false},
		{&data.OwnerAffiliation{}, &data.ParticipantRole{}, false},

		// Occupant: VisitorRole
		{&data.NoneAffiliation{}, &data.VisitorRole{}, true},
		{&data.MemberAffiliation{}, &data.VisitorRole{}, false},
		{&data.AdminAffiliation{}, &data.VisitorRole{}, false},
		{&data.OwnerAffiliation{}, &data.VisitorRole{}, false},
	}

	// Actor: MemberAffiliation - ModeratorRole
	actor := newTestOccupant(&data.MemberAffiliation{}, &data.ModeratorRole{})
	for _, scenario := range testCases {
		c.Assert(actor.CanKickOccupant(newTestOccupant(scenario.occupantAffiliation, scenario.occupantRole)), Equals, scenario.expected)
	}
}

func (*MucOccupantRoleAffiliationPrivilegesSuite) Test_RoleModeratorAffiliationAdmin_CanKickAnOccupant(c *C) {
	testCases := []canKickOccupantTest{
		// Occupant: ModeratorRole
		{&data.NoneAffiliation{}, &data.ModeratorRole{}, false},
		{&data.MemberAffiliation{}, &data.ModeratorRole{}, false},
		{&data.AdminAffiliation{}, &data.ModeratorRole{}, false},
		{&data.OwnerAffiliation{}, &data.ModeratorRole{}, false},

		// Occupant: ParticipantRole
		{&data.NoneAffiliation{}, &data.ParticipantRole{}, true},
		{&data.MemberAffiliation{}, &data.ParticipantRole{}, true},
		{&data.AdminAffiliation{}, &data.ParticipantRole{}, false},
		{&data.OwnerAffiliation{}, &data.ParticipantRole{}, false},

		// Occupant: VisitorRole
		{&data.NoneAffiliation{}, &data.VisitorRole{}, true},
		{&data.MemberAffiliation{}, &data.VisitorRole{}, true},
		{&data.AdminAffiliation{}, &data.VisitorRole{}, false},
		{&data.OwnerAffiliation{}, &data.VisitorRole{}, false},
	}

	// Actor: AdminAffiliation - ModeratorRole
	actor := newTestOccupant(&data.AdminAffiliation{}, &data.ModeratorRole{})
	for _, scenario := range testCases {
		c.Assert(actor.CanKickOccupant(newTestOccupant(scenario.occupantAffiliation, scenario.occupantRole)), Equals, scenario.expected)
	}
}

func (*MucOccupantRoleAffiliationPrivilegesSuite) Test_RoleModeratorAffiliationOwner_CanKickAnOccupant(c *C) {
	testCases := []canKickOccupantTest{
		// Occupant: ModeratorRole
		{&data.NoneAffiliation{}, &data.ModeratorRole{}, false},
		{&data.MemberAffiliation{}, &data.ModeratorRole{}, false},
		{&data.AdminAffiliation{}, &data.ModeratorRole{}, false},
		{&data.OwnerAffiliation{}, &data.ModeratorRole{}, false},

		// Occupant: ParticipantRole
		{&data.NoneAffiliation{}, &data.ParticipantRole{}, true},
		{&data.MemberAffiliation{}, &data.ParticipantRole{}, true},
		{&data.AdminAffiliation{}, &data.ParticipantRole{}, true},
		{&data.OwnerAffiliation{}, &data.ParticipantRole{}, false},

		// Occupant: VisitorRole
		{&data.NoneAffiliation{}, &data.VisitorRole{}, true},
		{&data.MemberAffiliation{}, &data.VisitorRole{}, true},
		{&data.AdminAffiliation{}, &data.VisitorRole{}, true},
		{&data.OwnerAffiliation{}, &data.VisitorRole{}, false},
	}

	// Actor: OwnerAffiliation - ModeratorRole
	actor := newTestOccupant(&data.OwnerAffiliation{}, &data.ModeratorRole{})
	for _, scenario := range testCases {
		c.Assert(actor.CanKickOccupant(newTestOccupant(scenario.occupantAffiliation, scenario.occupantRole)), Equals, scenario.expected)
	}
}

func (*MucOccupantRoleAffiliationPrivilegesSuite) Test_RoleParticipantAffiliationNone_CanKickAnOccupant(c *C) {
	testCases := []canKickOccupantTest{
		// Occupant: ModeratorRole
		{&data.NoneAffiliation{}, &data.ModeratorRole{}, false},
		{&data.MemberAffiliation{}, &data.ModeratorRole{}, false},
		{&data.AdminAffiliation{}, &data.ModeratorRole{}, false},
		{&data.OwnerAffiliation{}, &data.ModeratorRole{}, false},

		// Occupant: ParticipantRole
		{&data.NoneAffiliation{}, &data.ParticipantRole{}, false},
		{&data.MemberAffiliation{}, &data.ParticipantRole{}, false},
		{&data.AdminAffiliation{}, &data.ParticipantRole{}, false},
		{&data.OwnerAffiliation{}, &data.ParticipantRole{}, false},

		// Occupant: VisitorRole
		{&data.NoneAffiliation{}, &data.VisitorRole{}, false},
		{&data.MemberAffiliation{}, &data.VisitorRole{}, false},
		{&data.AdminAffiliation{}, &data.VisitorRole{}, false},
		{&data.OwnerAffiliation{}, &data.VisitorRole{}, false},
	}

	// Actor: NoneAffiliation - ParticipantRole
	actor := newTestOccupant(&data.NoneAffiliation{}, &data.ParticipantRole{})
	for _, scenario := range testCases {
		c.Assert(actor.CanKickOccupant(newTestOccupant(scenario.occupantAffiliation, scenario.occupantRole)), Equals, scenario.expected)
	}
}

func (*MucOccupantRoleAffiliationPrivilegesSuite) Test_RoleParticipantAffiliationMember_CanKickAnOccupant(c *C) {
	testCases := []canKickOccupantTest{
		// Occupant: ModeratorRole
		{&data.NoneAffiliation{}, &data.ModeratorRole{}, false},
		{&data.MemberAffiliation{}, &data.ModeratorRole{}, false},
		{&data.AdminAffiliation{}, &data.ModeratorRole{}, false},
		{&data.OwnerAffiliation{}, &data.ModeratorRole{}, false},

		// Occupant: ParticipantRole
		{&data.NoneAffiliation{}, &data.ParticipantRole{}, false},
		{&data.MemberAffiliation{}, &data.ParticipantRole{}, false},
		{&data.AdminAffiliation{}, &data.ParticipantRole{}, false},
		{&data.OwnerAffiliation{}, &data.ParticipantRole{}, false},

		// Occupant: VisitorRole
		{&data.NoneAffiliation{}, &data.VisitorRole{}, false},
		{&data.MemberAffiliation{}, &data.VisitorRole{}, false},
		{&data.AdminAffiliation{}, &data.VisitorRole{}, false},
		{&data.OwnerAffiliation{}, &data.VisitorRole{}, false},
	}

	// Actor: MemberAffiliation - ParticipantRole
	actor := newTestOccupant(&data.MemberAffiliation{}, &data.ParticipantRole{})
	for _, scenario := range testCases {
		c.Assert(actor.CanKickOccupant(newTestOccupant(scenario.occupantAffiliation, scenario.occupantRole)), Equals, scenario.expected)
	}
}

func (*MucOccupantRoleAffiliationPrivilegesSuite) Test_RoleParticipantAffiliationAdmin_CanKickAnOccupant(c *C) {
	testCases := []canKickOccupantTest{
		// Occupant: ModeratorRole
		{&data.NoneAffiliation{}, &data.ModeratorRole{}, false},
		{&data.MemberAffiliation{}, &data.ModeratorRole{}, false},
		{&data.AdminAffiliation{}, &data.ModeratorRole{}, false},
		{&data.OwnerAffiliation{}, &data.ModeratorRole{}, false},

		// Occupant: ParticipantRole
		{&data.NoneAffiliation{}, &data.ParticipantRole{}, false},
		{&data.MemberAffiliation{}, &data.ParticipantRole{}, false},
		{&data.AdminAffiliation{}, &data.ParticipantRole{}, false},
		{&data.OwnerAffiliation{}, &data.ParticipantRole{}, false},

		// Occupant: VisitorRole
		{&data.NoneAffiliation{}, &data.VisitorRole{}, false},
		{&data.MemberAffiliation{}, &data.VisitorRole{}, false},
		{&data.AdminAffiliation{}, &data.VisitorRole{}, false},
		{&data.OwnerAffiliation{}, &data.VisitorRole{}, false},
	}

	// Actor: AdminAffiliation - ParticipantRole
	actor := newTestOccupant(&data.AdminAffiliation{}, &data.ParticipantRole{})
	for _, scenario := range testCases {
		c.Assert(actor.CanKickOccupant(newTestOccupant(scenario.occupantAffiliation, scenario.occupantRole)), Equals, scenario.expected)
	}
}

func (*MucOccupantRoleAffiliationPrivilegesSuite) Test_RoleParticipantAffiliationOwner_CanKickAnOccupant(c *C) {
	testCases := []canKickOccupantTest{
		// Occupant: ModeratorRole
		{&data.NoneAffiliation{}, &data.ModeratorRole{}, false},
		{&data.MemberAffiliation{}, &data.ModeratorRole{}, false},
		{&data.AdminAffiliation{}, &data.ModeratorRole{}, false},
		{&data.OwnerAffiliation{}, &data.ModeratorRole{}, false},

		// Occupant: ParticipantRole
		{&data.NoneAffiliation{}, &data.ParticipantRole{}, false},
		{&data.MemberAffiliation{}, &data.ParticipantRole{}, false},
		{&data.AdminAffiliation{}, &data.ParticipantRole{}, false},
		{&data.OwnerAffiliation{}, &data.ParticipantRole{}, false},

		// Occupant: VisitorRole
		{&data.NoneAffiliation{}, &data.VisitorRole{}, false},
		{&data.MemberAffiliation{}, &data.VisitorRole{}, false},
		{&data.AdminAffiliation{}, &data.VisitorRole{}, false},
		{&data.OwnerAffiliation{}, &data.VisitorRole{}, false},
	}

	// Actor: OwnerAffiliation - ParticipantRole
	actor := newTestOccupant(&data.OwnerAffiliation{}, &data.ParticipantRole{})
	for _, scenario := range testCases {
		c.Assert(actor.CanKickOccupant(newTestOccupant(scenario.occupantAffiliation, scenario.occupantRole)), Equals, scenario.expected)
	}
}

func (*MucOccupantRoleAffiliationPrivilegesSuite) Test_RoleVisitorAffiliationNone_CanKickAnOccupant(c *C) {
	testCases := []canKickOccupantTest{
		// Occupant: ModeratorRole
		{&data.NoneAffiliation{}, &data.ModeratorRole{}, false},
		{&data.MemberAffiliation{}, &data.ModeratorRole{}, false},
		{&data.AdminAffiliation{}, &data.ModeratorRole{}, false},
		{&data.OwnerAffiliation{}, &data.ModeratorRole{}, false},

		// Occupant: ParticipantRole
		{&data.NoneAffiliation{}, &data.ParticipantRole{}, false},
		{&data.MemberAffiliation{}, &data.ParticipantRole{}, false},
		{&data.AdminAffiliation{}, &data.ParticipantRole{}, false},
		{&data.OwnerAffiliation{}, &data.ParticipantRole{}, false},

		// Occupant: VisitorRole
		{&data.NoneAffiliation{}, &data.VisitorRole{}, false},
		{&data.MemberAffiliation{}, &data.VisitorRole{}, false},
		{&data.AdminAffiliation{}, &data.VisitorRole{}, false},
		{&data.OwnerAffiliation{}, &data.VisitorRole{}, false},
	}

	// Actor: NoneAffiliation - VisitorRole
	actor := newTestOccupant(&data.NoneAffiliation{}, &data.VisitorRole{})
	for _, scenario := range testCases {
		c.Assert(actor.CanKickOccupant(newTestOccupant(scenario.occupantAffiliation, scenario.occupantRole)), Equals, scenario.expected)
	}
}

func (*MucOccupantRoleAffiliationPrivilegesSuite) Test_RoleVisitorAffiliationMember_CanKickAnOccupant(c *C) {
	testCases := []canKickOccupantTest{
		// Occupant: ModeratorRole
		{&data.NoneAffiliation{}, &data.ModeratorRole{}, false},
		{&data.MemberAffiliation{}, &data.ModeratorRole{}, false},
		{&data.AdminAffiliation{}, &data.ModeratorRole{}, false},
		{&data.OwnerAffiliation{}, &data.ModeratorRole{}, false},

		// Occupant: ParticipantRole
		{&data.NoneAffiliation{}, &data.ParticipantRole{}, false},
		{&data.MemberAffiliation{}, &data.ParticipantRole{}, false},
		{&data.AdminAffiliation{}, &data.ParticipantRole{}, false},
		{&data.OwnerAffiliation{}, &data.ParticipantRole{}, false},

		// Occupant: VisitorRole
		{&data.NoneAffiliation{}, &data.VisitorRole{}, false},
		{&data.MemberAffiliation{}, &data.VisitorRole{}, false},
		{&data.AdminAffiliation{}, &data.VisitorRole{}, false},
		{&data.OwnerAffiliation{}, &data.VisitorRole{}, false},
	}

	// Actor: MemberAffiliation - VisitorRole
	actor := newTestOccupant(&data.MemberAffiliation{}, &data.VisitorRole{})
	for _, scenario := range testCases {
		c.Assert(actor.CanKickOccupant(newTestOccupant(scenario.occupantAffiliation, scenario.occupantRole)), Equals, scenario.expected)
	}
}

func (*MucOccupantRoleAffiliationPrivilegesSuite) Test_RoleVisitorAffiliationAdmin_CanKickAnOccupant(c *C) {
	testCases := []canKickOccupantTest{
		// Occupant: ModeratorRole
		{&data.NoneAffiliation{}, &data.ModeratorRole{}, false},
		{&data.MemberAffiliation{}, &data.ModeratorRole{}, false},
		{&data.AdminAffiliation{}, &data.ModeratorRole{}, false},
		{&data.OwnerAffiliation{}, &data.ModeratorRole{}, false},

		// Occupant: ParticipantRole
		{&data.NoneAffiliation{}, &data.ParticipantRole{}, false},
		{&data.MemberAffiliation{}, &data.ParticipantRole{}, false},
		{&data.AdminAffiliation{}, &data.ParticipantRole{}, false},
		{&data.OwnerAffiliation{}, &data.ParticipantRole{}, false},

		// Occupant: VisitorRole
		{&data.NoneAffiliation{}, &data.VisitorRole{}, false},
		{&data.MemberAffiliation{}, &data.VisitorRole{}, false},
		{&data.AdminAffiliation{}, &data.VisitorRole{}, false},
		{&data.OwnerAffiliation{}, &data.VisitorRole{}, false},
	}

	// Actor: AdminAffiliation - VisitorRole
	actor := newTestOccupant(&data.AdminAffiliation{}, &data.VisitorRole{})
	for _, scenario := range testCases {
		c.Assert(actor.CanKickOccupant(newTestOccupant(scenario.occupantAffiliation, scenario.occupantRole)), Equals, scenario.expected)
	}
}

func (*MucOccupantRoleAffiliationPrivilegesSuite) Test_RoleVisitorAffiliationOwner_CanKickAnOccupant(c *C) {
	testCases := []canKickOccupantTest{
		// Occupant: ModeratorRole
		{&data.NoneAffiliation{}, &data.ModeratorRole{}, false},
		{&data.MemberAffiliation{}, &data.ModeratorRole{}, false},
		{&data.AdminAffiliation{}, &data.ModeratorRole{}, false},
		{&data.OwnerAffiliation{}, &data.ModeratorRole{}, false},

		// Occupant: ParticipantRole
		{&data.NoneAffiliation{}, &data.ParticipantRole{}, false},
		{&data.MemberAffiliation{}, &data.ParticipantRole{}, false},
		{&data.AdminAffiliation{}, &data.ParticipantRole{}, false},
		{&data.OwnerAffiliation{}, &data.ParticipantRole{}, false},

		// Occupant: VisitorRole
		{&data.NoneAffiliation{}, &data.VisitorRole{}, false},
		{&data.MemberAffiliation{}, &data.VisitorRole{}, false},
		{&data.AdminAffiliation{}, &data.VisitorRole{}, false},
		{&data.OwnerAffiliation{}, &data.VisitorRole{}, false},
	}

	// Actor: OwnerAffiliation - VisitorRole
	actor := newTestOccupant(&data.OwnerAffiliation{}, &data.VisitorRole{})
	for _, scenario := range testCases {
		c.Assert(actor.CanKickOccupant(newTestOccupant(scenario.occupantAffiliation, scenario.occupantRole)), Equals, scenario.expected)
	}
}
