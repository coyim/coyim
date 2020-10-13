package session

import (
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
)

func (m *mucManager) receivedClientMessage(stanza *data.ClientMessage) {
	m.log.WithField("stanza", stanza).Debug("handleMUCReceivedClientMessage()")

	if hasSubject(stanza) {
		m.handleSubjectReceived(stanza)
	}

	if isLiveMessage(stanza) {
		from := jid.ParseFull(stanza.From)
		roomID := from.Bare()
		nickname := from.Resource().String()
		message := stanza.Body

		m.log.WithFields(log.Fields{
			"roomID":   roomID,
			"message":  message,
			"nickname": nickname,
		}).Info("MUC message received")

		m.messageReceived(roomID, nickname, message)
	}
}

func (m *mucManager) handleSubjectReceived(stanza *data.ClientMessage) {
	l := m.log.WithField("from", stanza.From)

	roomID, ok := jid.TryParseBare(stanza.From)
	if !ok {
		l.Error("Error trying to get the room from stanza")
		return
	}

	room, ok := m.roomManager.GetRoom(roomID)
	if !ok {
		l.WithField("room", roomID).Error("Error trying to read the subject of room")
		return
	}

	subjectReceived := getSubjectFromStanza(stanza)
	room.Subject.Text = subjectReceived

	if !room.Subject.Received {
		room.Subject.Received = true
		m.subjectReceived(roomID, subjectReceived)
		return
	}

	from, ok := jid.TryParseFull(stanza.From)
	if !ok {
		l.WithField("room", roomID).Error("Error trying to get the nickname")
		return
	}
	nickname := from.Resource().String()
	m.subjectUpdated(roomID, nickname, subjectReceived)
}

func bodyHasContent(stanza *data.ClientMessage) bool {
	return stanza.Body != ""
}

func isMessageDelayed(stanza *data.ClientMessage) bool {
	return stanza.Delay != nil
}

func isLiveMessage(stanza *data.ClientMessage) bool {
	return bodyHasContent(stanza) && !isMessageDelayed(stanza)
}

func hasSubject(stanza *data.ClientMessage) bool {
	return stanza.Subject != nil
}

func getSubjectFromStanza(stanza *data.ClientMessage) string {
	if hasSubject(stanza) {
		return stanza.Subject.Text
	}
	return ""
}
