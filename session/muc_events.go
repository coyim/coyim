package session

import (
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
)

func (m *mucManager) publishRoomEvent(roomID jid.Bare, ev muc.MUC) {
	room, exists := m.roomManager.GetRoom(roomID)
	if !exists {
		m.log.WithField("room", roomID).Error("Trying to publish an event in a room that does not exist")
		return
	}
	room.Publish(ev)
}

func (m *mucManager) roomCreated(roomID jid.Bare) {
	ev := events.MUCRoomCreated{}
	ev.Room = roomID

	m.publishEvent(ev)
}

func (m *mucManager) roomRenamed(roomID jid.Bare) {
	m.publishRoomEvent(roomID, events.MUCRoomRenamed{})
}

func (m *mucManager) occupantLeft(roomID jid.Bare, occupant *muc.Occupant) {
	ev := events.MUCOccupantLeft{}
	ev.Nickname = occupant.Nick
	ev.RealJid = occupant.Jid
	ev.Affiliation = occupant.Affiliation
	ev.Role = occupant.Role

	m.publishRoomEvent(roomID, ev)
}

func (m *mucManager) occupantJoined(roomID jid.Bare, occupant *muc.Occupant) {
	ev := events.MUCOccupantJoined{}
	ev.Nickname = occupant.Nick
	ev.RealJid = occupant.Jid
	ev.Affiliation = occupant.Affiliation
	ev.Role = occupant.Role

	m.publishRoomEvent(roomID, ev)
}

func (m *mucManager) occupantUpdate(roomID jid.Bare, occupant *muc.Occupant) {
	ev := events.MUCOccupantUpdated{}
	ev.Nickname = occupant.Nick
	ev.RealJid = occupant.Jid
	ev.Affiliation = occupant.Affiliation
	ev.Role = occupant.Role

	m.publishRoomEvent(roomID, ev)
}

func (m *mucManager) loggingEnabled(roomID jid.Bare) {
	m.publishRoomEvent(roomID, events.MUCLoggingEnabled{})
}

func (m *mucManager) loggingDisabled(roomID jid.Bare) {
	m.publishRoomEvent(roomID, events.MUCLoggingDisabled{})
}

func (m *mucManager) selfOccupantJoined(roomID jid.Bare, occupant *muc.Occupant) {
	ev := events.MUCSelfOccupantJoined{}
	ev.Nickname = occupant.Nick
	ev.RealJid = occupant.Jid
	ev.Affiliation = occupant.Affiliation
	ev.Role = occupant.Role

	m.publishRoomEvent(roomID, ev)
}

func (m *mucManager) messageReceived(roomID jid.Bare, nickname, subject, message string) {
	ev := events.MUCMessageReceived{}
	ev.Nickname = nickname
	ev.Subject = subject
	ev.Message = message

	m.publishRoomEvent(roomID, ev)
}
