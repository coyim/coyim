package session

import (
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/xmpp/jid"
)

func (m *mucManager) mucRoomCreated(from jid.Bare) {
	ev := events.MUCRoomCreated{}

	m.publishMUCEvent(from, ev)
}

func (m *mucManager) mucRoomRenamed(from jid.WithoutResource) {
	ev := events.MUCRoomRenamed{}

	m.publishMUCEvent(from, ev)
}

func (m *mucManager) mucOccupantExit(from jid.WithoutResource, occupant jid.Resource) {
	ev := events.MUCOccupantExited{}
	ev.Jid = jid.NewBare(occupant, from)
	ev.Nickname = occupant

	m.publishMUCEvent(from, ev)
}

func (m *mucManager) mucOccupantUpdate(from jid.WithoutResource, occupant jid.Resource, affiliation, role string) {
	ev := events.MUCOccupantUpdated{}
	ev.Jid = jid.NewBare(occupant, from)
	ev.Nickname = occupant
	ev.Affiliation = affiliation
	ev.Role = role

	m.publishMUCEvent(from, ev)
}

func (m *mucManager) mucOccupantJoined(from jid.WithoutResource, occupant jid.Resource, affiliation, role string) {
	ev := events.MUCOccupantJoined{}
	ev.Jid = jid.NewBare(occupant, from)
	ev.Nickname = occupant
	ev.Affiliation = affiliation
	ev.Role = role

	m.publishMUCEvent(from, ev)
}

func (m *mucManager) publishMUCEvent(from jid.WithoutResource, e interface{}) {
	ev := events.MUC{}
	ev.From = from.(jid.Bare)
	ev.Info = e

	m.publishEvent(ev)
}

func (m *mucManager) publishEvent(ev interface{}) {
	if m.publishEv == nil {
		panic("programmer error: muc manager \"publishEv\" function not defined")
	}

	m.publishEv(ev)
}
