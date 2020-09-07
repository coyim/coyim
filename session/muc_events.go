package session

import (
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/xmpp/jid"
)

func (m *mucManager) mucRoomCreated(from jid.Full, room jid.Bare) {
	ev := events.MUCRoomCreated{}
	ev.Room = room

	m.publishMUCEvent(from, ev)
}

func (m *mucManager) mucRoomRenamed(from jid.Full, room jid.Bare) {
	ev := events.MUCRoomRenamed{}
	ev.Room = room

	m.publishMUCEvent(from, ev)
}

func (m *mucManager) mucOccupantLeft(from jid.Full, room jid.Bare, occupant jid.Resource, affiliation, role string) {
	ev := events.MUCOccupantLeft{}
	ev.Room = room
	ev.Nickname = occupant
	ev.Jid = from
	ev.Affiliation = affiliation
	ev.Role = role

	m.publishMUCEvent(from, ev)
}

func (m *mucManager) mucOccupantUpdate(from jid.Full, room jid.Bare, occupant jid.Resource, affiliation, role string) {
	ev := events.MUCOccupantUpdated{}
	ev.Room = room
	ev.Nickname = occupant
	ev.Jid = from
	ev.Affiliation = affiliation
	ev.Role = role

	m.publishMUCEvent(from, ev)
}

func (m *mucManager) mucOccupantJoined(from jid.Full, room jid.Bare, occupant jid.Resource, ident jid.Full, affiliation, role string) {
	ev := events.MUCOccupantJoined{}
	ev.Room = room
	ev.Nickname = occupant
	ev.Jid = ident
	ev.Affiliation = affiliation
	ev.Role = role

	m.publishMUCEvent(from, ev)
}

func (m *mucManager) mucMessageReceived(from jid.Full, room jid.Bare, nickname jid.Resource, subject, body string) {
	ev := events.MUCMessageReceived{}
	ev.From = from
	ev.Room = room
	ev.Nickname = nickname
	ev.Subject = subject
	ev.Body = body

	m.publishMUCEvent(from, ev)
}

func (m *mucManager) publishMUCEvent(from jid.Full, e interface{}) {
	ev := events.MUC{}
	ev.From = from
	ev.Info = e

	m.publishEvent(ev)
}
