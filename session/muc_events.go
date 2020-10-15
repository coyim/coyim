package session

import (
	"time"

	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
)

func (m *mucManager) publishRoomEvent(roomID jid.Bare, ev events.MUC) {
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

func (m *mucManager) occupantLeft(roomID jid.Bare, op *muc.OccupantPresenceInfo) {
	ev := events.MUCOccupantLeft{}
	ev.Nickname = op.Nickname
	ev.RealJid = op.RealJid
	ev.Affiliation = op.Affiliation
	ev.Role = op.Role

	m.publishRoomEvent(roomID, ev)
}

func (m *mucManager) occupantJoined(roomID jid.Bare, op *muc.OccupantPresenceInfo) {
	ev := events.MUCOccupantJoined{}
	ev.Nickname = op.Nickname
	ev.RealJid = op.RealJid
	ev.Affiliation = op.Affiliation
	ev.Role = op.Role

	m.publishRoomEvent(roomID, ev)
}

func (m *mucManager) occupantUpdate(roomID jid.Bare, op *muc.OccupantPresenceInfo) {
	ev := events.MUCOccupantUpdated{}
	ev.Nickname = op.Nickname
	ev.RealJid = op.RealJid
	ev.Affiliation = op.Affiliation
	ev.Role = op.Role

	m.publishRoomEvent(roomID, ev)
}

func (m *mucManager) loggingEnabled(roomID jid.Bare) {
	m.publishRoomEvent(roomID, events.MUCLoggingEnabled{})
}

func (m *mucManager) loggingDisabled(roomID jid.Bare) {
	m.publishRoomEvent(roomID, events.MUCLoggingDisabled{})
}

func (m *mucManager) selfOccupantJoined(roomID jid.Bare, op *muc.OccupantPresenceInfo) {
	ev := events.MUCSelfOccupantJoined{}
	ev.Nickname = op.Nickname
	ev.RealJid = op.RealJid
	ev.Affiliation = op.Affiliation
	ev.Role = op.Role

	m.publishRoomEvent(roomID, ev)
}

func (m *mucManager) liveMessageReceived(roomID jid.Bare, nickname, message string, timestamp time.Time) {
	ev := events.MUCLiveMessageReceived{}
	ev.Nickname = nickname
	ev.Message = message
	ev.Timestamp = timestamp

	m.publishRoomEvent(roomID, ev)
}

func (m *mucManager) delayedMessageReceived(roomID jid.Bare, nickname, message string, timestamp time.Time) {
	ev := events.MUCDelayedMessageReceived{}
	ev.Nickname = nickname
	ev.Message = message
	ev.Timestamp = timestamp

	m.publishRoomEvent(roomID, ev)
}

func (m *mucManager) subjectReceived(roomID jid.Bare, subject string) {
	ev := events.MUCSubjectReceived{}
	ev.Subject = subject

	m.publishRoomEvent(roomID, ev)
}

func (m *mucManager) subjectUpdated(roomID jid.Bare, nickname, subject string) {
	ev := events.MUCSubjectUpdated{}
	ev.Nickname = nickname
	ev.Subject = subject

	m.publishRoomEvent(roomID, ev)
}
