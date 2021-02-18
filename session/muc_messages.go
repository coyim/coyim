package session

import (
	"time"

	xmppData "github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
)

func (m *mucManager) receiveClientMessage(stanza *xmppData.ClientMessage) {
	m.log.WithField("stanza", stanza).Debug("handleMUCReceivedClientMessage()")

	if hasSubject(stanza) {
		m.dhManager.handleClientPresence(stanza)
		m.handleSubjectReceived(stanza)
	}

	switch {
	case isDelayedMessage(stanza):
		m.handleMessageReceived(stanza, m.receiveDelayedMessage)
	case isLiveMessage(stanza):
		m.handleMessageReceived(stanza, m.liveMessageReceived)
	case isRoomConfigUpdate(stanza):
		m.handleRoomConfigUpdate(stanza)
	}
}

func (m *mucManager) receiveDelayedMessage(roomID jid.Bare, nickname, message string, timestamp time.Time) {
	m.dhLock.Lock()
	defer m.dhLock.Unlock()

	dh, ok := m.dhManager.getHistory(roomID)
	if !ok {
		dh = m.dhManager.addHistory(roomID)
	}

	dh.AddMessage(nickname, message, timestamp)
}

// The discussion history MUST happen only one time in the events flow of XMPP's MUC
// This should be done in a proper way, maybe in the pending "state machine" pattern
// that we want to implement later, when that happens, this method should be fine
func (m *mucManager) handleDiscussionHistory(stanza *xmppData.ClientMessage) {
	roomID := m.retrieveRoomID(stanza.From, "handleDiscussionHistory")
	dh, exists := m.dhManager.getHistory(roomID)
	if exists {
		m.discussionHistoryReceived(roomID, dh)
	}
}

func (m *mucManager) handleSubjectReceived(stanza *xmppData.ClientMessage) {
	l := m.log.WithFields(log.Fields{
		"from": stanza.From,
		"who":  "handleSubjectReceived",
	})

	roomID := m.retrieveRoomID(stanza.From, "handleSubjectReceived")
	room, ok := m.roomManager.GetRoom(roomID)
	if !ok {
		l.WithField("room", roomID).Error("Error trying to read the subject of room")
		return
	}

	s := getSubjectFromStanza(stanza)
	updated := room.UpdateSubject(s)
	if updated {
		m.subjectUpdated(roomID, getNicknameFromStanza(stanza), s)
		return
	}

	m.subjectReceived(roomID, s)
}

func (m *mucManager) handleMessageReceived(stanza *xmppData.ClientMessage, h func(jid.Bare, string, string, time.Time)) {
	roomID, nickname := m.retrieveRoomIDAndNickname(stanza.From)
	h(roomID, nickname, stanza.Body, retrieveMessageTime(stanza))
}

func bodyHasContent(stanza *xmppData.ClientMessage) bool {
	return stanza.Body != ""
}

func isDelayedMessage(stanza *xmppData.ClientMessage) bool {
	return stanza.Delay != nil
}

func isLiveMessage(stanza *xmppData.ClientMessage) bool {
	return bodyHasContent(stanza) && !isDelayedMessage(stanza)
}

func hasSubject(stanza *xmppData.ClientMessage) bool {
	return stanza.Subject != nil
}

func hasMUCUserExtension(stanza *xmppData.ClientMessage) bool {
	return stanza.MUCUser != nil
}

func getNicknameFromStanza(stanza *xmppData.ClientMessage) string {
	from, ok := jid.TryParseFull(stanza.From)
	if ok {
		return from.Resource().String()
	}

	return ""
}

func getSubjectFromStanza(stanza *xmppData.ClientMessage) string {
	if hasSubject(stanza) {
		return stanza.Subject.Text
	}

	return ""
}
