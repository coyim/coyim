package session

import (
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/xmpp/jid"
)

func (s *session) mucRoomCreated(from jid.Bare) {
	ev := events.MUCRoomCreated{}

	s.publishMUCEvent(from, ev)
}

func (s *session) mucRoomRenamed(from jid.WithoutResource) {
	ev := events.MUCRoomRenamed{}

	s.publishMUCEvent(from, ev)
}

func (s *session) mucOccupantExit(from jid.WithoutResource, occupant jid.Resource) {
	ev := events.MUCOccupantExited{}
	ev.Jid = jid.NewBare(occupant, from)
	ev.Nickname = occupant

	s.publishMUCEvent(from, ev)
}

func (s *session) mucOccupantUpdate(from jid.WithoutResource, occupant jid.Resource, affiliation, role string) {
	ev := events.MUCOccupantUpdated{}
	ev.Jid = jid.NewBare(occupant, from)
	ev.Nickname = occupant
	ev.Affiliation = affiliation
	ev.Role = role

	s.publishMUCEvent(from, ev)
}

func (s *session) mucOccupantJoined(from jid.WithoutResource, occupant jid.Resource, affiliation, role string) {
	ev := events.MUCOccupantJoined{}
	ev.Jid = jid.NewBare(occupant, from)
	ev.Nickname = occupant
	ev.Affiliation = affiliation
	ev.Role = role

	s.publishMUCEvent(from, ev)
}

func (s *session) publishMUCEvent(from jid.WithoutResource, e events.MUCAny) {
	ev := events.MUC{}
	ev.From = from.(jid.Bare)
	ev.Info = e

	s.publishEvent(ev)
}
