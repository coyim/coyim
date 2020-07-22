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

func (s *MucSuite) Test_RoomRoster_AllOccupants(c *C) {
	rr := newRoomRoster()
	rr.occupants["one"] = &Occupant{Nick: "One", Jid: jid.R("foo@bar.com/somewhere")}
	rr.occupants["two"] = &Occupant{Nick: "Abc", Jid: jid.R("foo@bar.com/somewhereelse")}
	rr.occupants["three"] = &Occupant{Nick: "Foo", Jid: jid.R("foo@bar.com/foo")}

	oo := rr.AllOccupants()
	c.Assert(oo, HasLen, 3)
	c.Assert(oo[0].Nick, Equals, "Abc")
	c.Assert(oo[1].Nick, Equals, "Foo")
	c.Assert(oo[2].Nick, Equals, "One")
}

func (s *MucSuite) Test_RoomRoster_NoRole(c *C) {
	rr := newRoomRoster()
	rr.occupants["one"] = &Occupant{Nick: "One", Jid: jid.R("foo@bar.com/somewhere"), Role: &noneRole{}}
	rr.occupants["two"] = &Occupant{Nick: "Abc", Jid: jid.R("foo@bar.com/somewhereelse"), Role: &participantRole{}}
	rr.occupants["three"] = &Occupant{Nick: "Foo", Jid: jid.R("foo@bar.com/foo"), Role: &noneRole{}}

	oo := rr.NoRole()
	c.Assert(oo, HasLen, 2)
	c.Assert(oo[0].Nick, Equals, "Foo")
	c.Assert(oo[1].Nick, Equals, "One")
}

func (s *MucSuite) Test_RoomRoster_Visitors(c *C) {
	rr := newRoomRoster()
	rr.occupants["one"] = &Occupant{Nick: "One", Jid: jid.R("foo@bar.com/somewhere"), Role: &visitorRole{}}
	rr.occupants["two"] = &Occupant{Nick: "Abc", Jid: jid.R("foo@bar.com/somewhereelse"), Role: &participantRole{}}
	rr.occupants["three"] = &Occupant{Nick: "Foo", Jid: jid.R("foo@bar.com/foo"), Role: &visitorRole{}}

	oo := rr.Visitors()
	c.Assert(oo, HasLen, 2)
	c.Assert(oo[0].Nick, Equals, "Foo")
	c.Assert(oo[1].Nick, Equals, "One")
}

func (s *MucSuite) Test_RoomRoster_Participants(c *C) {
	rr := newRoomRoster()
	rr.occupants["one"] = &Occupant{Nick: "One", Jid: jid.R("foo@bar.com/somewhere"), Role: &participantRole{}}
	rr.occupants["two"] = &Occupant{Nick: "Abc", Jid: jid.R("foo@bar.com/somewhereelse"), Role: &visitorRole{}}
	rr.occupants["three"] = &Occupant{Nick: "Foo", Jid: jid.R("foo@bar.com/foo"), Role: &participantRole{}}

	oo := rr.Participants()
	c.Assert(oo, HasLen, 2)
	c.Assert(oo[0].Nick, Equals, "Foo")
	c.Assert(oo[1].Nick, Equals, "One")
}

func (s *MucSuite) Test_RoomRoster_Moderators(c *C) {
	rr := newRoomRoster()
	rr.occupants["one"] = &Occupant{Nick: "One", Jid: jid.R("foo@bar.com/somewhere"), Role: &moderatorRole{}}
	rr.occupants["two"] = &Occupant{Nick: "Abc", Jid: jid.R("foo@bar.com/somewhereelse"), Role: &visitorRole{}}
	rr.occupants["three"] = &Occupant{Nick: "Foo", Jid: jid.R("foo@bar.com/foo"), Role: &moderatorRole{}}

	oo := rr.Moderators()
	c.Assert(oo, HasLen, 2)
	c.Assert(oo[0].Nick, Equals, "Foo")
	c.Assert(oo[1].Nick, Equals, "One")
}

func (s *MucSuite) Test_RoomRoster_NoAffiliation(c *C) {
	rr := newRoomRoster()
	rr.occupants["one"] = &Occupant{Nick: "One", Jid: jid.R("foo@bar.com/somewhere"), Affiliation: &noneAffiliation{}}
	rr.occupants["two"] = &Occupant{Nick: "Abc", Jid: jid.R("foo@bar.com/somewhereelse"), Affiliation: &ownerAffiliation{}}
	rr.occupants["three"] = &Occupant{Nick: "Foo", Jid: jid.R("foo@bar.com/foo"), Affiliation: &noneAffiliation{}}

	oo := rr.NoAffiliation()
	c.Assert(oo, HasLen, 2)
	c.Assert(oo[0].Nick, Equals, "Foo")
	c.Assert(oo[1].Nick, Equals, "One")
}

func (s *MucSuite) Test_RoomRoster_Banned(c *C) {
	rr := newRoomRoster()
	rr.occupants["one"] = &Occupant{Nick: "One", Jid: jid.R("foo@bar.com/somewhere"), Affiliation: &outcastAffiliation{}}
	rr.occupants["two"] = &Occupant{Nick: "Abc", Jid: jid.R("foo@bar.com/somewhereelse"), Affiliation: &ownerAffiliation{}}
	rr.occupants["three"] = &Occupant{Nick: "Foo", Jid: jid.R("foo@bar.com/foo"), Affiliation: &outcastAffiliation{}}

	oo := rr.Banned()
	c.Assert(oo, HasLen, 2)
	c.Assert(oo[0].Nick, Equals, "Foo")
	c.Assert(oo[1].Nick, Equals, "One")
}

func (s *MucSuite) Test_RoomRoster_Members(c *C) {
	rr := newRoomRoster()
	rr.occupants["one"] = &Occupant{Nick: "One", Jid: jid.R("foo@bar.com/somewhere"), Affiliation: &memberAffiliation{}}
	rr.occupants["two"] = &Occupant{Nick: "Abc", Jid: jid.R("foo@bar.com/somewhereelse"), Affiliation: &ownerAffiliation{}}
	rr.occupants["three"] = &Occupant{Nick: "Foo", Jid: jid.R("foo@bar.com/foo"), Affiliation: &memberAffiliation{}}

	oo := rr.Members()
	c.Assert(oo, HasLen, 2)
	c.Assert(oo[0].Nick, Equals, "Foo")
	c.Assert(oo[1].Nick, Equals, "One")
}

func (s *MucSuite) Test_RoomRoster_Admins(c *C) {
	rr := newRoomRoster()
	rr.occupants["one"] = &Occupant{Nick: "One", Jid: jid.R("foo@bar.com/somewhere"), Affiliation: &adminAffiliation{}}
	rr.occupants["two"] = &Occupant{Nick: "Abc", Jid: jid.R("foo@bar.com/somewhereelse"), Affiliation: &ownerAffiliation{}}
	rr.occupants["three"] = &Occupant{Nick: "Foo", Jid: jid.R("foo@bar.com/foo"), Affiliation: &adminAffiliation{}}

	oo := rr.Admins()
	c.Assert(oo, HasLen, 2)
	c.Assert(oo[0].Nick, Equals, "Foo")
	c.Assert(oo[1].Nick, Equals, "One")
}

func (s *MucSuite) Test_RoomRoster_Owners(c *C) {
	rr := newRoomRoster()
	rr.occupants["one"] = &Occupant{Nick: "One", Jid: jid.R("foo@bar.com/somewhere"), Affiliation: &ownerAffiliation{}}
	rr.occupants["two"] = &Occupant{Nick: "Abc", Jid: jid.R("foo@bar.com/somewhereelse"), Affiliation: &adminAffiliation{}}
	rr.occupants["three"] = &Occupant{Nick: "Foo", Jid: jid.R("foo@bar.com/foo"), Affiliation: &ownerAffiliation{}}

	oo := rr.Owners()
	c.Assert(oo, HasLen, 2)
	c.Assert(oo[0].Nick, Equals, "Foo")
	c.Assert(oo[1].Nick, Equals, "One")
}

func (s *MucSuite) Test_RoomRoster_OccupantsByRole(c *C) {
	rr := newRoomRoster()
	rr.occupants["one"] = &Occupant{Nick: "One", Jid: jid.R("foo@bar.com/somewhere"), Role: &noneRole{}}
	rr.occupants["two"] = &Occupant{Nick: "Abc", Jid: jid.R("foo@bar.com/somewhereelse"), Role: &noneRole{}}
	rr.occupants["three"] = &Occupant{Nick: "Foo", Jid: jid.R("foo@bar.com/foo"), Role: &participantRole{}}
	rr.occupants["four"] = &Occupant{Nick: "Bar", Jid: jid.R("foo@bar.com/bar"), Role: &visitorRole{}}
	rr.occupants["five"] = &Occupant{Nick: "Baz", Jid: jid.R("foo@bar.com/baz"), Role: &participantRole{}}
	rr.occupants["six"] = &Occupant{Nick: "Quux", Jid: jid.R("foo@bar.com/quu"), Role: &moderatorRole{}}

	r1, r2, r3, r4 := rr.OccupantsByRole()
	c.Assert(r1, HasLen, 2)
	c.Assert(r1[0].Nick, Equals, "Abc")
	c.Assert(r1[1].Nick, Equals, "One")

	c.Assert(r2, HasLen, 1)
	c.Assert(r2[0].Nick, Equals, "Bar")

	c.Assert(r3, HasLen, 2)
	c.Assert(r3[0].Nick, Equals, "Baz")
	c.Assert(r3[1].Nick, Equals, "Foo")

	c.Assert(r4, HasLen, 1)
	c.Assert(r4[0].Nick, Equals, "Quux")
}

func (s *MucSuite) Test_RoomRoster_OccupantsByAffiliation(c *C) {
	rr := newRoomRoster()
	rr.occupants["one"] = &Occupant{Nick: "One", Jid: jid.R("foo@bar.com/somewhere"), Affiliation: &noneAffiliation{}}
	rr.occupants["two"] = &Occupant{Nick: "Abc", Jid: jid.R("foo@bar.com/somewhereelse"), Affiliation: &noneAffiliation{}}
	rr.occupants["three"] = &Occupant{Nick: "Foo", Jid: jid.R("foo@bar.com/foo"), Affiliation: &memberAffiliation{}}
	rr.occupants["four"] = &Occupant{Nick: "Bar", Jid: jid.R("foo@bar.com/bar"), Affiliation: &outcastAffiliation{}}
	rr.occupants["five"] = &Occupant{Nick: "Baz", Jid: jid.R("foo@bar.com/baz"), Affiliation: &memberAffiliation{}}
	rr.occupants["six"] = &Occupant{Nick: "Quux", Jid: jid.R("foo@bar.com/quu"), Affiliation: &adminAffiliation{}}

	r1, r2, r3, r4, r5 := rr.OccupantsByAffiliation()
	c.Assert(r1, HasLen, 2)
	c.Assert(r1[0].Nick, Equals, "Abc")
	c.Assert(r1[1].Nick, Equals, "One")

	c.Assert(r2, HasLen, 1)
	c.Assert(r2[0].Nick, Equals, "Bar")

	c.Assert(r3, HasLen, 2)
	c.Assert(r3[0].Nick, Equals, "Baz")
	c.Assert(r3[1].Nick, Equals, "Foo")

	c.Assert(r4, HasLen, 1)
	c.Assert(r4[0].Nick, Equals, "Quux")

	c.Assert(r5, HasLen, 0)
}

func (s *MucSuite) Test_RoomRoster_UpdateNick(c *C) {
	rr := newRoomRoster()
	e := rr.UpdateNick(jid.R("foo@somewhere.com/bello"), "something")
	c.Assert(e, ErrorMatches, "no such occupant known in this room")

	occ := &Occupant{Nick: "bello"}
	rr.occupants["foo@somewhere.com/bello"] = occ

	e = rr.UpdateNick(jid.R("foo@somewhere.com/bello"), "baxa")
	c.Assert(e, IsNil)
	c.Assert(occ.Nick, Equals, "baxa")
	newOcc, ok := rr.occupants["foo@somewhere.com/baxa"]
	c.Assert(ok, Equals, true)
	c.Assert(newOcc, Equals, occ)
}

func (s *MucSuite) Test_RoomRoster_UpdatePresence_unavailable(c *C) {
	rr := newRoomRoster()

	j, l, e := rr.UpdatePresence(jid.R("foo@somewhere.com/bello"), "unavailable", "none", "none", "away", "101", "gone", nil)
	c.Assert(j, Equals, false)
	c.Assert(l, Equals, false)
	c.Assert(e, ErrorMatches, "no such occupant known in this room")

	occ := &Occupant{Nick: "bello"}
	rr.occupants["foo@somewhere.com/bello"] = occ

	j, l, e = rr.UpdatePresence(jid.R("foo@somewhere.com/bello"), "unavailable", "none", "none", "away", "101", "gone", nil)
	c.Assert(j, Equals, false)
	c.Assert(l, Equals, true)
	c.Assert(e, IsNil)
	_, ok := rr.occupants["foo@somewhere.com/bello"]
	c.Assert(ok, Equals, false)
}

func (s *MucSuite) Test_RoomRoster_UpdatePresence_bad_type(c *C) {
	rr := newRoomRoster()

	j, l, e := rr.UpdatePresence(jid.R("foo@somewhere.com/bello"), "hungry", "none", "none", "away", "101", "gone", nil)
	c.Assert(j, Equals, false)
	c.Assert(l, Equals, false)
	c.Assert(e, ErrorMatches, "incorrect presence type sent to room roster: 'hungry'")
}

func (s *MucSuite) Test_RoomRoster_UpdatePresence_new(c *C) {
	rr := newRoomRoster()

	j, l, e := rr.UpdatePresence(jid.R("foo@somewhere.com/bello"), "", "none", "none", "away", "101", "gone", jid.R("foo@example.org/test1"))
	c.Assert(j, Equals, true)
	c.Assert(l, Equals, false)
	c.Assert(e, IsNil)

	occ, ok := rr.occupants["foo@somewhere.com/bello"]
	c.Assert(ok, Equals, true)
	c.Assert(occ, Not(IsNil))
	c.Assert(occ.Nick, Equals, "bello")
	c.Assert(occ.Role, FitsTypeOf, &noneRole{})
	c.Assert(occ.Affiliation, FitsTypeOf, &noneAffiliation{})
	c.Assert(occ.Status.Status, Equals, "away")
	c.Assert(occ.Status.StatusMsg, Equals, "gone")
	c.Assert(occ.Jid, Equals, jid.R("foo@example.org/test1"))

	j, l, e = rr.UpdatePresence(jid.R("foo@somewhere.com/bello2"), "", "nonex", "none", "away", "101", "gone", jid.R("foo@example.org/test2"))
	c.Assert(j, Equals, false)
	c.Assert(l, Equals, false)
	c.Assert(e, ErrorMatches, "unknown affiliation string: 'nonex'")
}

func (s *MucSuite) Test_RoomRoster_UpdatePresence_update(c *C) {
	rr := newRoomRoster()

	occ := &Occupant{Nick: "bello"}
	rr.occupants["foo@somewhere.com/bello"] = occ

	j, l, e := rr.UpdatePresence(jid.R("foo@somewhere.com/bello"), "", "none", "none", "away", "101", "gone", jid.R("foo@example.org/test1"))
	c.Assert(j, Equals, false)
	c.Assert(l, Equals, false)
	c.Assert(e, IsNil)

	c.Assert(occ.Affiliation, FitsTypeOf, &noneAffiliation{})
	c.Assert(occ.Role, FitsTypeOf, &noneRole{})
	c.Assert(occ.Status.Status, Equals, "away")
	c.Assert(occ.Status.StatusMsg, Equals, "gone")
	c.Assert(occ.Jid, Equals, jid.R("foo@example.org/test1"))

	j, l, e = rr.UpdatePresence(jid.R("foo@somewhere.com/bello"), "", "none", "nonexx", "away", "101", "gone", jid.R("foo@example.org/test2"))
	c.Assert(j, Equals, false)
	c.Assert(l, Equals, false)
	c.Assert(e, ErrorMatches, "unknown role string: 'nonexx'")
}
