package session

import (
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
)

func (m *mucManager) receiveClientMessage(stanza *data.ClientMessage) {
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

		m.liveMessageReceived(roomID, nickname, message)
	}
}

func (m *mucManager) handleSubjectReceived(stanza *data.ClientMessage) {
	l := m.log.WithFields(log.Fields{
		"from": stanza.From,
		"who":  "handleSubjectReceived",
	})

	roomID, ok := jid.TryParseBare(stanza.From)
	if !ok {
		l.Error("Error trying to get the room ID from stanza")
		return
	}

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

func getNicknameFromStanza(stanza *data.ClientMessage) string {
	from, ok := jid.TryParseFull(stanza.From)
	if ok {
		return from.Resource().String()
	}

	return ""
}

func getSubjectFromStanza(stanza *data.ClientMessage) string {
	if hasSubject(stanza) {
		return stanza.Subject.Text
	}
	return ""
}
