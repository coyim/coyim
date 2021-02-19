package muc

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	. "gopkg.in/check.v1"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

type MucOccupantRolePrivilegesSuite struct{}

var _ = Suite(&MucOccupantRolePrivilegesSuite{})

func (*MucOccupantRolePrivilegesSuite) Test_OccupantCanPresentInRoom(c *C) {
	o := &Occupant{}

	o.ChangeRoleToNone()
	c.Assert(o.CanPresentInRoom(), Equals, false)

	o.ChangeRoleToVisitor()
	c.Assert(o.CanPresentInRoom(), Equals, true)

	o.ChangeRoleToParticipant()
	c.Assert(o.CanPresentInRoom(), Equals, true)

	o.ChangeRoleToModerator()
	c.Assert(o.CanPresentInRoom(), Equals, true)
}

func (*MucOccupantRolePrivilegesSuite) Test_OccupantCanReceiveMessage(c *C) {
	o := &Occupant{}

	o.ChangeRoleToNone()
	c.Assert(o.CanReceiveMessage(), Equals, false)

	o.ChangeRoleToVisitor()
	c.Assert(o.CanReceiveMessage(), Equals, true)

	o.ChangeRoleToParticipant()
	c.Assert(o.CanReceiveMessage(), Equals, true)

	o.ChangeRoleToModerator()
	c.Assert(o.CanReceiveMessage(), Equals, true)
}

func (*MucOccupantRolePrivilegesSuite) Test_OccupantCanReceiveOccupantPresence(c *C) {
	o := &Occupant{}

	o.ChangeRoleToNone()
	c.Assert(o.CanReceiveOccupantPresence(), Equals, false)

	o.ChangeRoleToVisitor()
	c.Assert(o.CanReceiveOccupantPresence(), Equals, true)

	o.ChangeRoleToParticipant()
	c.Assert(o.CanReceiveOccupantPresence(), Equals, true)

	o.ChangeRoleToModerator()
	c.Assert(o.CanReceiveOccupantPresence(), Equals, true)
}

func (*MucOccupantRolePrivilegesSuite) Test_OccupantCanBroadcastPresenceToAllOccupants(c *C) {
	o := &Occupant{}

	o.ChangeRoleToNone()
	c.Assert(o.CanBroadcastPresenceToAllOccupants(), Equals, false)

	o.ChangeRoleToVisitor()
	c.Assert(o.CanBroadcastPresenceToAllOccupants(), Equals, true)

	o.ChangeRoleToParticipant()
	c.Assert(o.CanBroadcastPresenceToAllOccupants(), Equals, true)

	o.ChangeRoleToModerator()
	c.Assert(o.CanBroadcastPresenceToAllOccupants(), Equals, true)
}

func (*MucOccupantRolePrivilegesSuite) Test_OccupantCanChangeAvailabilityStatus(c *C) {
	o := &Occupant{}

	o.ChangeRoleToNone()
	c.Assert(o.CanChangeAvailabilityStatus(), Equals, false)

	o.ChangeRoleToVisitor()
	c.Assert(o.CanChangeAvailabilityStatus(), Equals, true)

	o.ChangeRoleToParticipant()
	c.Assert(o.CanChangeAvailabilityStatus(), Equals, true)

	o.ChangeRoleToModerator()
	c.Assert(o.CanChangeAvailabilityStatus(), Equals, true)
}

func (*MucOccupantRolePrivilegesSuite) Test_OccupantCanChangeRoomNickname(c *C) {
	o := &Occupant{}

	o.ChangeRoleToNone()
	c.Assert(o.CanChangeRoomNickname(), Equals, false)

	o.ChangeRoleToVisitor()
	c.Assert(o.CanChangeRoomNickname(), Equals, true)

	o.ChangeRoleToParticipant()
	c.Assert(o.CanChangeRoomNickname(), Equals, true)

	o.ChangeRoleToModerator()
	c.Assert(o.CanChangeRoomNickname(), Equals, true)
}

func (*MucOccupantRolePrivilegesSuite) Test_OccupantCanSendPrivateMessages(c *C) {
	o := &Occupant{}

	o.ChangeRoleToNone()
	c.Assert(o.CanSendPrivateMessages(), Equals, false)

	o.ChangeRoleToVisitor()
	c.Assert(o.CanSendPrivateMessages(), Equals, true)

	o.ChangeRoleToParticipant()
	c.Assert(o.CanSendPrivateMessages(), Equals, true)

	o.ChangeRoleToModerator()
	c.Assert(o.CanSendPrivateMessages(), Equals, true)
}

func (*MucOccupantRolePrivilegesSuite) Test_OccupantCanInviteOtherUsers(c *C) {
	o := &Occupant{}

	o.ChangeRoleToNone()
	c.Assert(o.CanInviteOtherUsers(), Equals, false)

	o.ChangeRoleToVisitor()
	c.Assert(o.CanInviteOtherUsers(), Equals, true)

	o.ChangeRoleToParticipant()
	c.Assert(o.CanInviteOtherUsers(), Equals, true)

	o.ChangeRoleToModerator()
	c.Assert(o.CanInviteOtherUsers(), Equals, true)
}

func (*MucOccupantRolePrivilegesSuite) Test_OccupantCanSendMessagesToAll(c *C) {
	o := &Occupant{}

	o.ChangeRoleToNone()
	c.Assert(o.CanSendMessagesToAll(), Equals, false)

	o.ChangeRoleToVisitor()
	c.Assert(o.CanSendMessagesToAll(), Equals, false)

	o.ChangeRoleToParticipant()
	c.Assert(o.CanSendMessagesToAll(), Equals, true)

	o.ChangeRoleToModerator()
	c.Assert(o.CanSendMessagesToAll(), Equals, true)
}

func (*MucOccupantRolePrivilegesSuite) Test_OccupantCanModifySubject(c *C) {
	o := &Occupant{}

	o.ChangeRoleToNone()
	c.Assert(o.CanModifySubject(), Equals, false)

	o.ChangeRoleToVisitor()
	c.Assert(o.CanModifySubject(), Equals, false)

	o.ChangeRoleToParticipant()
	c.Assert(o.CanModifySubject(), Equals, true)

	o.ChangeRoleToModerator()
	c.Assert(o.CanModifySubject(), Equals, true)
}

func (*MucOccupantRolePrivilegesSuite) Test_OccupantCanKickParticipantsAndVisitors(c *C) {
	o := &Occupant{}

	o.ChangeRoleToNone()
	c.Assert(o.CanKickParticipantsAndVisitors(), Equals, false)

	o.ChangeRoleToVisitor()
	c.Assert(o.CanKickParticipantsAndVisitors(), Equals, false)

	o.ChangeRoleToParticipant()
	c.Assert(o.CanKickParticipantsAndVisitors(), Equals, false)

	o.ChangeRoleToModerator()
	c.Assert(o.CanKickParticipantsAndVisitors(), Equals, true)
}

func (*MucOccupantRolePrivilegesSuite) Test_OccupantCanGrantVoice(c *C) {
	o := &Occupant{}

	o.ChangeRoleToNone()
	c.Assert(o.CanGrantVoice(), Equals, false)

	o.ChangeRoleToVisitor()
	c.Assert(o.CanGrantVoice(), Equals, false)

	o.ChangeRoleToParticipant()
	c.Assert(o.CanGrantVoice(), Equals, false)

	o.ChangeRoleToModerator()
	c.Assert(o.CanGrantVoice(), Equals, true)
}

func (*MucOccupantRolePrivilegesSuite) Test_OccupantCanRevokeVoice_OfNoneRole(c *C) {
	o := &Occupant{}

	oc := &Occupant{}
	oc.ChangeAffiliationToNone()
	oc.ChangeRoleToNone()

	o.ChangeRoleToNone()
	c.Assert(o.CanRevokeVoice(oc), Equals, false)

	o.ChangeRoleToVisitor()
	c.Assert(o.CanRevokeVoice(oc), Equals, false)

	o.ChangeRoleToParticipant()
	c.Assert(o.CanRevokeVoice(oc), Equals, false)

	o.ChangeRoleToModerator()
	c.Assert(o.CanRevokeVoice(oc), Equals, true)
}

func (*MucOccupantRolePrivilegesSuite) Test_OccupantCanRevokeVoice_OfVisitorRole(c *C) {
	o := &Occupant{}

	oc := &Occupant{}
	oc.ChangeAffiliationToNone()
	oc.ChangeRoleToVisitor()

	o.ChangeRoleToNone()
	c.Assert(o.CanRevokeVoice(oc), Equals, false)

	o.ChangeRoleToVisitor()
	c.Assert(o.CanRevokeVoice(oc), Equals, false)

	o.ChangeRoleToParticipant()
	c.Assert(o.CanRevokeVoice(oc), Equals, false)

	o.ChangeRoleToModerator()
	c.Assert(o.CanRevokeVoice(oc), Equals, true)
}

func (*MucOccupantRolePrivilegesSuite) Test_OccupantCanRevokeVoice_OfParticipantRole(c *C) {
	o := &Occupant{}

	oc := &Occupant{}
	oc.ChangeAffiliationToNone()
	oc.ChangeRoleToParticipant()

	o.ChangeRoleToNone()
	c.Assert(o.CanRevokeVoice(oc), Equals, false)

	o.ChangeRoleToVisitor()
	c.Assert(o.CanRevokeVoice(oc), Equals, false)

	o.ChangeRoleToParticipant()
	c.Assert(o.CanRevokeVoice(oc), Equals, false)

	o.ChangeRoleToModerator()
	c.Assert(o.CanRevokeVoice(oc), Equals, true)
}

func (*MucOccupantRolePrivilegesSuite) Test_OccupantCanRevokeVoice_OfAdminOccupant(c *C) {
	o := &Occupant{}

	oc := &Occupant{}
	oc.ChangeAffiliationToAdmin()

	o.ChangeRoleToNone()
	c.Assert(o.CanRevokeVoice(oc), Equals, false)

	o.ChangeRoleToVisitor()
	c.Assert(o.CanRevokeVoice(oc), Equals, false)

	o.ChangeRoleToParticipant()
	c.Assert(o.CanRevokeVoice(oc), Equals, false)

	o.ChangeRoleToModerator()
	c.Assert(o.CanRevokeVoice(oc), Equals, false)
}

func (*MucOccupantRolePrivilegesSuite) Test_OccupantCanRevokeVoice_OfOwnerOccupant(c *C) {
	o := &Occupant{}

	oc := &Occupant{}
	oc.ChangeAffiliationToOwner()

	o.ChangeRoleToNone()
	c.Assert(o.CanRevokeVoice(oc), Equals, false)

	o.ChangeRoleToVisitor()
	c.Assert(o.CanRevokeVoice(oc), Equals, false)

	o.ChangeRoleToParticipant()
	c.Assert(o.CanRevokeVoice(oc), Equals, false)

	o.ChangeRoleToModerator()
	c.Assert(o.CanRevokeVoice(oc), Equals, false)
}
