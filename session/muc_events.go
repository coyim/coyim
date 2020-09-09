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

func (m *mucManager) publishLoggingEnabled(room jid.Bare) {
	ev := events.MUC{}
	ev.Room = room
	ev.Info = events.MUCLoggingEnabled{}

	m.publishEvent(ev)
}

func (m *mucManager) publishLoggingDisabled(room jid.Bare) {
	ev := events.MUC{}
	ev.Room = room
	ev.Info = events.MUCLoggingDisabled{}

	m.publishEvent(ev)
}

// mucSelfOccupantUpdated can happen several times - every time a status code update is changed, or role or affiliation
// is updated, this can lead to the method being called. For now, it will generate a event about joining, but this
// should be cleaned up and fixed
func (m *mucManager) mucSelfOccupantUpdated(from jid.Full, room jid.Bare, occupant jid.Resource, ident jid.Full, affiliation, role string, status mucUserStatuses) {
	ev := events.MUCOccupantJoined{}
	ev.Room = room
	ev.Nickname = occupant
	ev.Jid = ident
	ev.Affiliation = affiliation
	ev.Role = role

	m.publishMUCEvent(from, ev)

	if status.contains(MUCStatusRoomLoggingEnabled) {
		m.publishLoggingEnabled(room)
	}

	if status.contains(MUCStatusRoomLoggingDisabled) {
		m.publishLoggingDisabled(room)
	}
}

func (m *mucManager) mucMessageReceived(from jid.Full, room jid.Bare, nickname jid.Resource, message string) {
	ev := events.MUCMessageReceived{}
	ev.Room = room
	ev.Nickname = nickname
	ev.Message = message

	m.publishMUCEvent(from, ev)
}

func (m *mucManager) publishMUCEvent(from jid.Full, e interface{}) {
	ev := events.MUC{}
	ev.From = from
	ev.Info = e

	m.publishEvent(ev)
}
