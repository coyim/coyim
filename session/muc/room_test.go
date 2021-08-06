package muc

import (
	"time"

	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/coyim/xmpp/jid"

	. "gopkg.in/check.v1"
)

func (s *MucSuite) Test_NewRoom_createsRoom(c *C) {
	r := NewRoom(jid.ParseBare("foo@bar.com"))
	c.Assert(r, Not(IsNil))
	c.Assert(r.ID, DeepEquals, jid.ParseBare("foo@bar.com"))
}

func (s *MucSuite) Test_Room_SelfOccupant(c *C) {
	r := NewRoom(jid.ParseBare("foo@bar.com"))

	c.Assert(r.SelfOccupant(), IsNil)

	vv := &Occupant{}
	c.Assert(r.IsSelfOccupantInTheRoom(), Equals, false)
	c.Assert(r.IsSelfOccupantAnOwner(), Equals, false)
	r.AddSelfOccupant(vv)
	c.Assert(r.IsSelfOccupantInTheRoom(), Equals, true)
	vv.Affiliation = &data.NoneAffiliation{}
	c.Assert(r.IsSelfOccupantAnOwner(), Equals, false)
	vv.Affiliation = &data.OwnerAffiliation{}
	c.Assert(r.IsSelfOccupantAnOwner(), Equals, true)

	c.Assert(r.SelfOccupant(), Equals, vv)
}

func (s *MucSuite) Test_Room_SelfOccupantNickname(c *C) {
	r := NewRoom(jid.ParseBare("foo@bar.com"))

	c.Assert(r.SelfOccupantNickname(), Equals, "")

	vv := &Occupant{Nickname: "something"}
	r.selfOccupant = vv

	c.Assert(r.SelfOccupantNickname(), Equals, "something")
}

func (s *MucSuite) Test_Room_UpdateSubject(c *C) {
	r := NewRoom(jid.ParseBare("foo@bar.com"))
	c.Assert(r.HasSubject(), Equals, false)

	c.Assert(r.UpdateSubject("one"), Equals, false)
	c.Assert(r.HasSubject(), Equals, true)

	c.Assert(r.UpdateSubject("two"), Equals, true)
	c.Assert(r.GetSubject(), Equals, "two")
}

func (s *MucSuite) Test_Room_Roster(c *C) {
	r := NewRoom(jid.ParseBare("foo@bar.com"))
	vv := &RoomRoster{}
	r.roster = vv
	c.Assert(r.Roster(), Equals, vv)
}

func (s *MucSuite) Test_Room_Subscribe(c *C) {
	r := NewRoom(jid.ParseBare("foo@bar.com"))
	r.Subscribe(func(events.MUC) {})
	c.Assert(r.observers.observers, HasLen, 1)
}

func (s *MucSuite) Test_Room_Publish(c *C) {
	r := NewRoom(jid.ParseBare("foo@bar.com"))
	called1 := false
	r.Subscribe(func(events.MUC) {
		called1 = true
	})
	called2 := false
	r.Subscribe(func(events.MUC) {
		called2 = true
	})
	r.Publish(&events.MUCOccupantUpdated{})
	c.Assert(called1, Equals, true)
	c.Assert(called2, Equals, true)
}

func (s *MucSuite) Test_Room_History(c *C) {
	r := NewRoom(jid.ParseBare("foo@bar.com"))
	c.Assert(r.HasHistory(), Equals, false)

	r.AddHistoryMessage("juanito", "test message", time.Now())
	c.Assert(r.HasHistory(), Equals, true)
	roomHistory := r.GetDiscussionHistory()
	c.Assert(roomHistory.GetHistory(), HasLen, 1)

}

func (s *MucSuite) Test_Room_SetProperties(c *C) {
	r := NewRoom(jid.ParseBare("foo@bar.com"))
	c.Assert(r.properties, Equals, data.RoomDiscoInfo{})

	r.SetProperties(data.RoomDiscoInfo{
		Language:           "en",
		AllowsRegistration: true,
	})

	c.Assert(r.properties.Language, Equals, "en")
	c.Assert(r.properties.AllowsRegistration, Equals, true)
}

func (s *MucSuite) Test_Room_CanChangeSubject(c *C) {
	r := NewRoom(jid.ParseBare("foo@bar.com"))
	r.AddSelfOccupant(newTestOccupant(&data.OwnerAffiliation{}, &data.VisitorRole{}))

	r.SetProperties(data.RoomDiscoInfo{
		Language:                  "en",
		OccupantsCanChangeSubject: false,
	})

	c.Assert(r.SubjectCanBeChanged(), Equals, false)

	r.properties.OccupantsCanChangeSubject = true
	c.Assert(r.SubjectCanBeChanged(), Equals, false)

	r.AddSelfOccupant(newTestOccupant(&data.OwnerAffiliation{}, &data.ParticipantRole{}))
	c.Assert(r.SubjectCanBeChanged(), Equals, true)
}
