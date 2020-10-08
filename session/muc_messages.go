package session

import (
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
)

func (m *mucManager) receivedClientMessage(stanza *data.ClientMessage) {
	m.log.WithField("stanza", stanza).Debug("handleMUCReceivedClientMessage()")

	if isLiveMessage(stanza) {
		from := jid.ParseFull(stanza.From)
		roomID := from.Bare()
		nickname := from.Resource().String()
		message := stanza.Body
		subject := ""

		if hasSubject(stanza) {
			m.handleSubjectReceived(stanza)
		}

		m.log.WithFields(log.Fields{
			"roomID":   roomID,
			"message":  message,
			"subject":  subject,
			"nickname": nickname,
		}).Info("MUC message received")

		m.messageReceived(roomID, nickname, subject, message)
		return
	}

	if hasSubject(stanza) {
		m.handleSubjectReceived(stanza)
	}
}

func (m *mucManager) handleSubjectReceived(stanza *data.ClientMessage) {
	roomID := jid.ParseBare(stanza.From)
	room, ok := m.roomManager.GetRoom(roomID)
	if !ok {
		m.log.WithFields(log.Fields{
			"roomID": roomID,
		}).Error("Error trying to read the subject of room")
	}

	room.Subject = stanza.Subject.Text

	m.subjectReceived(roomID, room.Subject)
}

func bodyHasContent(stanza *data.ClientMessage) bool {
	return len(stanza.Body) > 0
}

func isMessageDelayed(stanza *data.ClientMessage) bool {
	return stanza.Delay != nil
}

func isLiveMessage(stanza *data.ClientMessage) bool {
	return bodyHasContent(stanza) && !isMessageDelayed(stanza)
}

func subjectHasContent(stanza *data.ClientMessage) bool {
	return stanza.Subject != nil && len(stanza.Subject.Text) > 0
}

func hasSubject(stanza *data.ClientMessage) bool {
	return !bodyHasContent(stanza) && subjectHasContent(stanza)
}
