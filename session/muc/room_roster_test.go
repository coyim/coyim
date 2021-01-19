package muc

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"

	. "gopkg.in/check.v1"

	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/coyim/xmpp/jid"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func getRealJidFromString(realJid string) jid.Full {
	if len(realJid) > 0 {
		return jid.ParseFull(realJid)
	}
	return nil
}

func newRosterOccupantPresenceForTest(nickname string, realJid string, role data.Role, oa *OccupantAffiliationInfo, status, statusMessage string) *OccupantPresenceInfo {
	return &OccupantPresenceInfo{
		Nickname:        nickname,
		RealJid:         getRealJidFromString(realJid),
		Role:            role,
		AffiliationInfo: oa,
		Status:          status,
		StatusMessage:   statusMessage,
	}
}

func newRosterOccupantForTest(nickname string, realJid string, role data.Role, affiliation data.Affiliation) *Occupant {
	return &Occupant{
		Nickname:    nickname,
		RealJid:     getRealJidFromString(realJid),
		Role:        role,
		Affiliation: affiliation,
	}
}

func (s *MucSuite) Test_RoomRoster_AllOccupants(c *C) {
	rr := newRoomRoster()
	rr.occupants["one"] = newRosterOccupantForTest("One", "foo@bar.com/somewhere", nil, nil)
	rr.occupants["two"] = newRosterOccupantForTest("Abc", "foo@bar.com/somewhereelse", nil, nil)
	rr.occupants["three"] = newRosterOccupantForTest("Foo", "foo@bar.com/foo", nil, nil)

	oo := rr.AllOccupants()
	c.Assert(oo, HasLen, 3)
	c.Assert(oo[0].Nickname, Equals, "Abc")
	c.Assert(oo[1].Nickname, Equals, "Foo")
	c.Assert(oo[2].Nickname, Equals, "One")
}

func (s *MucSuite) Test_RoomRoster_NoRole(c *C) {
	rr := newRoomRoster()
	rr.occupants["one"] = newRosterOccupantForTest("One", "foo@bar.com/somewhere", &data.NoneRole{}, nil)
	rr.occupants["two"] = newRosterOccupantForTest("Abc", "foo@bar.com/somewhereelse", &data.ParticipantRole{}, nil)
	rr.occupants["three"] = newRosterOccupantForTest("Foo", "foo@bar.com/foo", &data.NoneRole{}, nil)

	oo := rr.NoRole()
	c.Assert(oo, HasLen, 2)
	c.Assert(oo[0].Nickname, Equals, "Foo")
	c.Assert(oo[1].Nickname, Equals, "One")
}

func (s *MucSuite) Test_RoomRoster_Visitors(c *C) {
	rr := newRoomRoster()
	rr.occupants["one"] = newRosterOccupantForTest("One", "foo@bar.com/somewhere", &data.VisitorRole{}, nil)
	rr.occupants["two"] = newRosterOccupantForTest("Abc", "foo@bar.com/somewhereelse", &data.ParticipantRole{}, nil)
	rr.occupants["three"] = newRosterOccupantForTest("Foo", "foo@bar.com/foo", &data.VisitorRole{}, nil)

	oo := rr.Visitors()
	c.Assert(oo, HasLen, 2)
	c.Assert(oo[0].Nickname, Equals, "Foo")
	c.Assert(oo[1].Nickname, Equals, "One")
}

func (s *MucSuite) Test_RoomRoster_Participants(c *C) {
	rr := newRoomRoster()
	rr.occupants["one"] = newRosterOccupantForTest("One", "foo@bar.com/somewhere", &data.ParticipantRole{}, nil)
	rr.occupants["two"] = newRosterOccupantForTest("Abc", "foo@bar.com/somewhereelse", &data.VisitorRole{}, nil)
	rr.occupants["three"] = newRosterOccupantForTest("Foo", "foo@bar.com/foo", &data.ParticipantRole{}, nil)

	oo := rr.Participants()
	c.Assert(oo, HasLen, 2)
	c.Assert(oo[0].Nickname, Equals, "Foo")
	c.Assert(oo[1].Nickname, Equals, "One")
}

func (s *MucSuite) Test_RoomRoster_Moderators(c *C) {
	rr := newRoomRoster()
	rr.occupants["one"] = newRosterOccupantForTest("One", "foo@bar.com/somewhere", &data.ModeratorRole{}, nil)
	rr.occupants["two"] = newRosterOccupantForTest("Abc", "foo@bar.com/somewhereelse", &data.VisitorRole{}, nil)
	rr.occupants["three"] = newRosterOccupantForTest("Foo", "foo@bar.com/foo", &data.ModeratorRole{}, nil)

	oo := rr.Moderators()
	c.Assert(oo, HasLen, 2)
	c.Assert(oo[0].Nickname, Equals, "Foo")
	c.Assert(oo[1].Nickname, Equals, "One")
}

func (s *MucSuite) Test_RoomRoster_NoAffiliation(c *C) {
	rr := newRoomRoster()
	rr.occupants["one"] = newRosterOccupantForTest("One", "foo@bar.com/somewhere", nil, &data.NoneAffiliation{})
	rr.occupants["two"] = newRosterOccupantForTest("Abc", "foo@bar.com/somewhereelse", nil, &data.OwnerAffiliation{})
	rr.occupants["three"] = newRosterOccupantForTest("Foo", "foo@bar.com/foo", nil, &data.NoneAffiliation{})

	oo := rr.NoAffiliation()
	c.Assert(oo, HasLen, 2)
	c.Assert(oo[0].Nickname, Equals, "Foo")
	c.Assert(oo[1].Nickname, Equals, "One")
}

func (s *MucSuite) Test_RoomRoster_Banned(c *C) {
	rr := newRoomRoster()
	rr.occupants["one"] = newRosterOccupantForTest("One", "foo@bar.com/somewhere", nil, &data.OutcastAffiliation{})
	rr.occupants["two"] = newRosterOccupantForTest("Abc", "foo@bar.com/somewhereelse", nil, &data.OwnerAffiliation{})
	rr.occupants["three"] = newRosterOccupantForTest("Foo", "foo@bar.com/foo", nil, &data.OutcastAffiliation{})

	oo := rr.Banned()
	c.Assert(oo, HasLen, 2)
	c.Assert(oo[0].Nickname, Equals, "Foo")
	c.Assert(oo[1].Nickname, Equals, "One")
}

func (s *MucSuite) Test_RoomRoster_Members(c *C) {
	rr := newRoomRoster()
	rr.occupants["one"] = newRosterOccupantForTest("One", "foo@bar.com/somewhere", nil, &data.MemberAffiliation{})
	rr.occupants["two"] = newRosterOccupantForTest("Abc", "foo@bar.com/somewhereelse", nil, &data.OwnerAffiliation{})
	rr.occupants["three"] = newRosterOccupantForTest("Foo", "foo@bar.com/foo", nil, &data.MemberAffiliation{})

	oo := rr.Members()
	c.Assert(oo, HasLen, 2)
	c.Assert(oo[0].Nickname, Equals, "Foo")
	c.Assert(oo[1].Nickname, Equals, "One")
}

func (s *MucSuite) Test_RoomRoster_Admins(c *C) {
	rr := newRoomRoster()
	rr.occupants["one"] = newRosterOccupantForTest("One", "foo@bar.com/somewhere", nil, &data.AdminAffiliation{})
	rr.occupants["two"] = newRosterOccupantForTest("Abc", "foo@bar.com/somewhereelse", nil, &data.OwnerAffiliation{})
	rr.occupants["three"] = newRosterOccupantForTest("Foo", "foo@bar.com/foo", nil, &data.AdminAffiliation{})

	oo := rr.Admins()
	c.Assert(oo, HasLen, 2)
	c.Assert(oo[0].Nickname, Equals, "Foo")
	c.Assert(oo[1].Nickname, Equals, "One")
}

func (s *MucSuite) Test_RoomRoster_Owners(c *C) {
	rr := newRoomRoster()
	rr.occupants["one"] = newRosterOccupantForTest("One", "foo@bar.com/somewhere", nil, &data.OwnerAffiliation{})
	rr.occupants["two"] = newRosterOccupantForTest("Abc", "foo@bar.com/somewhereelse", nil, &data.AdminAffiliation{})
	rr.occupants["three"] = newRosterOccupantForTest("Foo", "foo@bar.com/foo", nil, &data.OwnerAffiliation{})

	oo := rr.Owners()
	c.Assert(oo, HasLen, 2)
	c.Assert(oo[0].Nickname, Equals, "Foo")
	c.Assert(oo[1].Nickname, Equals, "One")
}

func (s *MucSuite) Test_RoomRoster_OccupantsByRole(c *C) {
	rr := newRoomRoster()
	rr.occupants["one"] = newRosterOccupantForTest("One", "foo@bar.com/somewhere", &data.NoneRole{}, nil)
	rr.occupants["two"] = newRosterOccupantForTest("Abc", "foo@bar.com/somewhereelse", &data.NoneRole{}, nil)
	rr.occupants["three"] = newRosterOccupantForTest("Foo", "foo@bar.com/foo", &data.ParticipantRole{}, nil)
	rr.occupants["four"] = newRosterOccupantForTest("Bar", "foo@bar.com/bar", &data.VisitorRole{}, nil)
	rr.occupants["five"] = newRosterOccupantForTest("Baz", "foo@bar.com/baz", &data.ParticipantRole{}, nil)
	rr.occupants["six"] = newRosterOccupantForTest("Quux", "foo@bar.com/quu", &data.ModeratorRole{}, nil)

	r1, r2, r3, r4 := rr.OccupantsByRole()
	c.Assert(r1, HasLen, 2)
	c.Assert(r1[0].Nickname, Equals, "Abc")
	c.Assert(r1[1].Nickname, Equals, "One")

	c.Assert(r2, HasLen, 1)
	c.Assert(r2[0].Nickname, Equals, "Bar")

	c.Assert(r3, HasLen, 2)
	c.Assert(r3[0].Nickname, Equals, "Baz")
	c.Assert(r3[1].Nickname, Equals, "Foo")

	c.Assert(r4, HasLen, 1)
	c.Assert(r4[0].Nickname, Equals, "Quux")
}

func (s *MucSuite) Test_RoomRoster_OccupantsByAffiliation(c *C) {
	rr := newRoomRoster()
	rr.occupants["one"] = newRosterOccupantForTest("One", "foo@bar.com/somewhere", nil, &data.NoneAffiliation{})
	rr.occupants["two"] = newRosterOccupantForTest("Abc", "foo@bar.com/somewhereelse", nil, &data.NoneAffiliation{})
	rr.occupants["three"] = newRosterOccupantForTest("Foo", "foo@bar.com/foo", nil, &data.MemberAffiliation{})
	rr.occupants["four"] = newRosterOccupantForTest("Bar", "foo@bar.com/bar", nil, &data.OutcastAffiliation{})
	rr.occupants["five"] = newRosterOccupantForTest("Baz", "foo@bar.com/baz", nil, &data.MemberAffiliation{})
	rr.occupants["six"] = newRosterOccupantForTest("Quux", "foo@bar.com/quu", nil, &data.AdminAffiliation{})

	r1, r2, r3, r4, r5 := rr.OccupantsByAffiliation()
	c.Assert(r1, HasLen, 2)
	c.Assert(r1[0].Nickname, Equals, "Abc")
	c.Assert(r1[1].Nickname, Equals, "One")

	c.Assert(r2, HasLen, 1)
	c.Assert(r2[0].Nickname, Equals, "Bar")

	c.Assert(r3, HasLen, 2)
	c.Assert(r3[0].Nickname, Equals, "Baz")
	c.Assert(r3[1].Nickname, Equals, "Foo")

	c.Assert(r4, HasLen, 1)
	c.Assert(r4[0].Nickname, Equals, "Quux")

	c.Assert(r5, HasLen, 0)
}

func (s *MucSuite) Test_RoomRoster_UpdateNickname(c *C) {
	rr := newRoomRoster()
	e := rr.UpdateNickname("bello", "something")
	c.Assert(e, ErrorMatches, "no such occupant known in this room")

	occ := newRosterOccupantForTest("bello", "", nil, nil)
	rr.occupants["bello"] = occ

	e = rr.UpdateNickname("bello", "baxa")
	c.Assert(e, IsNil)
	c.Assert(occ.Nickname, Equals, "baxa")
	newOcc, ok := rr.occupants["baxa"]
	c.Assert(ok, Equals, true)
	c.Assert(newOcc, Equals, occ)
}

func (s *MucSuite) Test_RoomRoster_UpdatePresence_unavailable(c *C) {
	rr := newRoomRoster()

	oa := &OccupantAffiliationInfo{
		Affiliation: &data.NoneAffiliation{},
	}

	o := newRosterOccupantPresenceForTest("bello", "", &data.NoneRole{}, oa, "away", "gone")

	j, l, e := rr.UpdatePresence(o, "unavailable")
	c.Assert(j, Equals, false)
	c.Assert(l, Equals, false)
	c.Assert(e, ErrorMatches, "no such occupant known in this room")

	occ := newRosterOccupantForTest("bello", "", nil, nil)
	rr.occupants["bello"] = occ

	j, l, e = rr.UpdatePresence(o, "unavailable")
	c.Assert(j, Equals, false)
	c.Assert(l, Equals, true)
	c.Assert(e, IsNil)
	_, ok := rr.occupants["bello"]
	c.Assert(ok, Equals, false)
}

func (s *MucSuite) Test_RoomRoster_UpdatePresence_bad_type(c *C) {
	rr := newRoomRoster()

	oa := &OccupantAffiliationInfo{
		Affiliation: &data.NoneAffiliation{},
	}

	o := newRosterOccupantPresenceForTest("bello", "", &data.NoneRole{}, oa, "away", "gone")

	j, l, e := rr.UpdatePresence(o, "hungry")
	c.Assert(j, Equals, false)
	c.Assert(l, Equals, false)
	c.Assert(e, ErrorMatches, "incorrect presence type sent to room roster: 'hungry'")
}

func (s *MucSuite) Test_RoomRoster_UpdatePresence_new(c *C) {
	rr := newRoomRoster()

	oa := &OccupantAffiliationInfo{
		Affiliation: &data.NoneAffiliation{},
	}
	o := newRosterOccupantPresenceForTest("bello", "foo@example.org/test1", &data.NoneRole{}, oa, "away", "gone")

	j, l, e := rr.UpdatePresence(o, "")
	c.Assert(j, Equals, true)
	c.Assert(l, Equals, false)
	c.Assert(e, IsNil)

	occ, ok := rr.occupants["bello"]
	c.Assert(ok, Equals, true)
	c.Assert(occ, Not(IsNil))
	c.Assert(occ.Nickname, Equals, "bello")
	c.Assert(occ.Role, FitsTypeOf, &data.NoneRole{})
	c.Assert(occ.Affiliation, FitsTypeOf, &data.NoneAffiliation{})
	c.Assert(occ.Status.Status, Equals, "away")
	c.Assert(occ.Status.StatusMsg, Equals, "gone")
	c.Assert(occ.RealJid, Equals, jid.ParseFull("foo@example.org/test1"))
}

func (s *MucSuite) Test_RoomRoster_UpdatePresence_update(c *C) {
	rr := newRoomRoster()

	occ := newRosterOccupantForTest("bello", "", nil, nil)
	rr.occupants["bello"] = occ

	oa := &OccupantAffiliationInfo{
		Affiliation: &data.NoneAffiliation{},
	}
	o := newRosterOccupantPresenceForTest("bello", "foo@example.org/test1", &data.NoneRole{}, oa, "away", "gone")

	j, l, e := rr.UpdatePresence(o, "")

	occ = rr.occupants["bello"]

	c.Assert(j, Equals, false)
	c.Assert(l, Equals, false)
	c.Assert(e, IsNil)

	c.Assert(occ.Affiliation, FitsTypeOf, &data.NoneAffiliation{})
	c.Assert(occ.Role, FitsTypeOf, &data.NoneRole{})
	c.Assert(occ.Status.Status, Equals, "away")
	c.Assert(occ.Status.StatusMsg, Equals, "gone")
	c.Assert(occ.RealJid, Equals, jid.ParseFull("foo@example.org/test1"))
}
