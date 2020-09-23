package session

import (
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
)

func (m *mucManager) publishRoomEvent(ident jid.Bare, ev muc.MUC) {
	room, exists := m.roomManager.GetRoom(ident)
	if !exists {
		m.log.WithField("room", ident).Error("Trying to publish an event in a room that does not exist")
		return
	}
	room.Publish(ev)
}

func (m *mucManager) roomCreated(from jid.Full, room jid.Bare) {
	ev := events.MUCRoomCreated{}
	ev.Room = room

	m.publishEvent(ev)
}

func (m *mucManager) roomRenamed(from jid.Full, room jid.Bare) {
	// TODO: This should use publishRoomEvent
	ev := events.MUCRoomRenamed{}
	ev.Room = room

	m.publishEvent(ev)
}

func (m *mucManager) publishOccupantLeft(from jid.Full, room jid.Bare, occupant mucRoomOccupant) {
	// TODO: for room events, we should not have a "Room" field in the event

	ev := events.MUCOccupantLeft{}
	ev.Room = room
	ev.Nickname = occupant.nickname
	ev.RealJid = from
	ev.Affiliation = occupant.affiliation
	ev.Role = occupant.role

	m.publishRoomEvent(room, ev)
}

func (m *mucManager) publishOccupantJoined(from jid.Full, room jid.Bare, occupant mucRoomOccupant) {
	// TODO: for room events, we should not have a "Room" field in the event

	ev := events.MUCOccupantJoined{}
	ev.Room = room
	ev.Nickname = occupant.nickname
	ev.RealJid = from
	ev.Affiliation = occupant.affiliation
	ev.Role = occupant.role

	m.publishRoomEvent(room, ev)
}

func (m *mucManager) publishOccupantUpdate(from jid.Full, room jid.Bare, occupant mucRoomOccupant) {
	// TODO: for room events, we should not have a "Room" field in the event

	ev := events.MUCOccupantUpdated{}
	ev.Room = room
	ev.Nickname = occupant.nickname
	ev.RealJid = from
	ev.Affiliation = occupant.affiliation
	ev.Role = occupant.role

	m.publishRoomEvent(room, ev)
}

func (m *mucManager) publishLoggingEnabled(ident jid.Bare) {
	// TODO: for room events, we should not have a "Room" field in the event

	ev := events.MUCLoggingEnabled{}
	ev.Room = ident

	m.publishRoomEvent(ident, ev)
}

func (m *mucManager) publishLoggingDisabled(room jid.Bare) {
	// TODO: for room events, we should not have a "Room" field in the event

	ev := events.MUCLoggingDisabled{}
	ev.Room = room

	m.publishRoomEvent(room, ev)
}

func (m *mucManager) publishSelfOccupantJoined(from jid.Full, room jid.Bare, occupant mucRoomOccupant) {
	// TODO: for room events, we should not have a "Room" field in the event

	ev := events.MUCSelfOccupantJoined{}
	ev.Room = room
	ev.Nickname = occupant.nickname
	ev.RealJid = occupant.realJid
	ev.Affiliation = occupant.affiliation
	ev.Role = occupant.role

	m.publishRoomEvent(room, ev)
}

func (m *mucManager) messageReceived(room jid.Bare, nickname, subject, message string) {
	// TODO: for room events, we should not have a "Room" field in the event

	ev := events.MUCMessageReceived{}
	ev.Room = room
	ev.Nickname = nickname
	ev.Subject = subject
	ev.Message = message

	m.publishRoomEvent(room, ev)
}
