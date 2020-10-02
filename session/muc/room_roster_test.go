package muc

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"

	. "gopkg.in/check.v1"

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

func newRosterOccupantPresenceForTest(nickname string, realJid string, role Role, affiliation Affiliation, statusCode, statusMessage string) *OccupantPresenceInfo {
	return &OccupantPresenceInfo{
		Nickname:      nickname,
		RealJid:       getRealJidFromString(realJid),
		Role:          role,
		Affiliation:   affiliation,
		StatusCode:    statusCode,
		StatusMessage: statusMessage,
	}
}

func newRosterOccupantForTest(nickname string, realJid string, role Role, affiliation Affiliation) *Occupant {
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
	rr.occupants["one"] = newRosterOccupantForTest("One", "foo@bar.com/somewhere", &noneRole{}, nil)
	rr.occupants["two"] = newRosterOccupantForTest("Abc", "foo@bar.com/somewhereelse", &participantRole{}, nil)
	rr.occupants["three"] = newRosterOccupantForTest("Foo", "foo@bar.com/foo", &noneRole{}, nil)

	oo := rr.NoRole()
	c.Assert(oo, HasLen, 2)
	c.Assert(oo[0].Nickname, Equals, "Foo")
	c.Assert(oo[1].Nickname, Equals, "One")
}

func (s *MucSuite) Test_RoomRoster_Visitors(c *C) {
	rr := newRoomRoster()
	rr.occupants["one"] = newRosterOccupantForTest("One", "foo@bar.com/somewhere", &visitorRole{}, nil)
	rr.occupants["two"] = newRosterOccupantForTest("Abc", "foo@bar.com/somewhereelse", &participantRole{}, nil)
	rr.occupants["three"] = newRosterOccupantForTest("Foo", "foo@bar.com/foo", &visitorRole{}, nil)

	oo := rr.Visitors()
	c.Assert(oo, HasLen, 2)
	c.Assert(oo[0].Nickname, Equals, "Foo")
	c.Assert(oo[1].Nickname, Equals, "One")
}

func (s *MucSuite) Test_RoomRoster_Participants(c *C) {
	rr := newRoomRoster()
	rr.occupants["one"] = newRosterOccupantForTest("One", "foo@bar.com/somewhere", &participantRole{}, nil)
	rr.occupants["two"] = newRosterOccupantForTest("Abc", "foo@bar.com/somewhereelse", &visitorRole{}, nil)
	rr.occupants["three"] = newRosterOccupantForTest("Foo", "foo@bar.com/foo", &participantRole{}, nil)

	oo := rr.Participants()
	c.Assert(oo, HasLen, 2)
	c.Assert(oo[0].Nickname, Equals, "Foo")
	c.Assert(oo[1].Nickname, Equals, "One")
}

func (s *MucSuite) Test_RoomRoster_Moderators(c *C) {
	rr := newRoomRoster()
	rr.occupants["one"] = newRosterOccupantForTest("One", "foo@bar.com/somewhere", &moderatorRole{}, nil)
	rr.occupants["two"] = newRosterOccupantForTest("Abc", "foo@bar.com/somewhereelse", &visitorRole{}, nil)
	rr.occupants["three"] = newRosterOccupantForTest("Foo", "foo@bar.com/foo", &moderatorRole{}, nil)

	oo := rr.Moderators()
	c.Assert(oo, HasLen, 2)
	c.Assert(oo[0].Nickname, Equals, "Foo")
	c.Assert(oo[1].Nickname, Equals, "One")
}

func (s *MucSuite) Test_RoomRoster_NoAffiliation(c *C) {
	rr := newRoomRoster()
	rr.occupants["one"] = newRosterOccupantForTest("One", "foo@bar.com/somewhere", nil, &noneAffiliation{})
	rr.occupants["two"] = newRosterOccupantForTest("Abc", "foo@bar.com/somewhereelse", nil, &ownerAffiliation{})
	rr.occupants["three"] = newRosterOccupantForTest("Foo", "foo@bar.com/foo", nil, &noneAffiliation{})

	oo := rr.NoAffiliation()
	c.Assert(oo, HasLen, 2)
	c.Assert(oo[0].Nickname, Equals, "Foo")
	c.Assert(oo[1].Nickname, Equals, "One")
}

func (s *MucSuite) Test_RoomRoster_Banned(c *C) {
	rr := newRoomRoster()
	rr.occupants["one"] = newRosterOccupantForTest("One", "foo@bar.com/somewhere", nil, &outcastAffiliation{})
	rr.occupants["two"] = newRosterOccupantForTest("Abc", "foo@bar.com/somewhereelse", nil, &ownerAffiliation{})
	rr.occupants["three"] = newRosterOccupantForTest("Foo", "foo@bar.com/foo", nil, &outcastAffiliation{})

	oo := rr.Banned()
	c.Assert(oo, HasLen, 2)
	c.Assert(oo[0].Nickname, Equals, "Foo")
	c.Assert(oo[1].Nickname, Equals, "One")
}

func (s *MucSuite) Test_RoomRoster_Members(c *C) {
	rr := newRoomRoster()
	rr.occupants["one"] = newRosterOccupantForTest("One", "foo@bar.com/somewhere", nil, &memberAffiliation{})
	rr.occupants["two"] = newRosterOccupantForTest("Abc", "foo@bar.com/somewhereelse", nil, &ownerAffiliation{})
	rr.occupants["three"] = newRosterOccupantForTest("Foo", "foo@bar.com/foo", nil, &memberAffiliation{})

	oo := rr.Members()
	c.Assert(oo, HasLen, 2)
	c.Assert(oo[0].Nickname, Equals, "Foo")
	c.Assert(oo[1].Nickname, Equals, "One")
}

func (s *MucSuite) Test_RoomRoster_Admins(c *C) {
	rr := newRoomRoster()
	rr.occupants["one"] = newRosterOccupantForTest("One", "foo@bar.com/somewhere", nil, &adminAffiliation{})
	rr.occupants["two"] = newRosterOccupantForTest("Abc", "foo@bar.com/somewhereelse", nil, &ownerAffiliation{})
	rr.occupants["three"] = newRosterOccupantForTest("Foo", "foo@bar.com/foo", nil, &adminAffiliation{})

	oo := rr.Admins()
	c.Assert(oo, HasLen, 2)
	c.Assert(oo[0].Nickname, Equals, "Foo")
	c.Assert(oo[1].Nickname, Equals, "One")
}

func (s *MucSuite) Test_RoomRoster_Owners(c *C) {
	rr := newRoomRoster()
	rr.occupants["one"] = newRosterOccupantForTest("One", "foo@bar.com/somewhere", nil, &ownerAffiliation{})
	rr.occupants["two"] = newRosterOccupantForTest("Abc", "foo@bar.com/somewhereelse", nil, &adminAffiliation{})
	rr.occupants["three"] = newRosterOccupantForTest("Foo", "foo@bar.com/foo", nil, &ownerAffiliation{})

	oo := rr.Owners()
	c.Assert(oo, HasLen, 2)
	c.Assert(oo[0].Nickname, Equals, "Foo")
	c.Assert(oo[1].Nickname, Equals, "One")
}

func (s *MucSuite) Test_RoomRoster_OccupantsByRole(c *C) {
	rr := newRoomRoster()
	rr.occupants["one"] = newRosterOccupantForTest("One", "foo@bar.com/somewhere", &noneRole{}, nil)
	rr.occupants["two"] = newRosterOccupantForTest("Abc", "foo@bar.com/somewhereelse", &noneRole{}, nil)
	rr.occupants["three"] = newRosterOccupantForTest("Foo", "foo@bar.com/foo", &participantRole{}, nil)
	rr.occupants["four"] = newRosterOccupantForTest("Bar", "foo@bar.com/bar", &visitorRole{}, nil)
	rr.occupants["five"] = newRosterOccupantForTest("Baz", "foo@bar.com/baz", &participantRole{}, nil)
	rr.occupants["six"] = newRosterOccupantForTest("Quux", "foo@bar.com/quu", &moderatorRole{}, nil)

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
	rr.occupants["one"] = newRosterOccupantForTest("One", "foo@bar.com/somewhere", nil, &noneAffiliation{})
	rr.occupants["two"] = newRosterOccupantForTest("Abc", "foo@bar.com/somewhereelse", nil, &noneAffiliation{})
	rr.occupants["three"] = newRosterOccupantForTest("Foo", "foo@bar.com/foo", nil, &memberAffiliation{})
	rr.occupants["four"] = newRosterOccupantForTest("Bar", "foo@bar.com/bar", nil, &outcastAffiliation{})
	rr.occupants["five"] = newRosterOccupantForTest("Baz", "foo@bar.com/baz", nil, &memberAffiliation{})
	rr.occupants["six"] = newRosterOccupantForTest("Quux", "foo@bar.com/quu", nil, &adminAffiliation{})

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

	o := newRosterOccupantPresenceForTest("bello", "", &noneRole{}, &noneAffiliation{}, "away", "gone")

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

	o := newRosterOccupantPresenceForTest("bello", "", &noneRole{}, &noneAffiliation{}, "away", "gone")

	j, l, e := rr.UpdatePresence(o, "hungry")
	c.Assert(j, Equals, false)
	c.Assert(l, Equals, false)
	c.Assert(e, ErrorMatches, "incorrect presence type sent to room roster: 'hungry'")
}

func (s *MucSuite) Test_RoomRoster_UpdatePresence_new(c *C) {
	rr := newRoomRoster()

	o := newRosterOccupantPresenceForTest("bello", "foo@example.org/test1", &noneRole{}, &noneAffiliation{}, "away", "gone")

	j, l, e := rr.UpdatePresence(o, "")
	c.Assert(j, Equals, true)
	c.Assert(l, Equals, false)
	c.Assert(e, IsNil)

	occ, ok := rr.occupants["bello"]
	c.Assert(ok, Equals, true)
	c.Assert(occ, Not(IsNil))
	c.Assert(occ.Nickname, Equals, "bello")
	c.Assert(occ.Role, FitsTypeOf, &noneRole{})
	c.Assert(occ.Affiliation, FitsTypeOf, &noneAffiliation{})
	c.Assert(occ.Status.Status, Equals, "away")
	c.Assert(occ.Status.StatusMsg, Equals, "gone")
	c.Assert(occ.RealJid, Equals, jid.ParseFull("foo@example.org/test1"))
}

func (s *MucSuite) Test_RoomRoster_UpdatePresence_update(c *C) {
	rr := newRoomRoster()

	occ := newRosterOccupantForTest("bello", "", nil, nil)
	rr.occupants["bello"] = occ

	o := newRosterOccupantPresenceForTest("bello", "foo@example.org/test1", &noneRole{}, &noneAffiliation{}, "away", "gone")

	j, l, e := rr.UpdatePresence(o, "")

	occ = rr.occupants["bello"]

	c.Assert(j, Equals, false)
	c.Assert(l, Equals, false)
	c.Assert(e, IsNil)

	c.Assert(occ.Affiliation, FitsTypeOf, &noneAffiliation{})
	c.Assert(occ.Role, FitsTypeOf, &noneRole{})
	c.Assert(occ.Status.Status, Equals, "away")
	c.Assert(occ.Status.StatusMsg, Equals, "gone")
	c.Assert(occ.RealJid, Equals, jid.ParseFull("foo@example.org/test1"))
}
