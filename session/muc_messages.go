package session

import (
	"strconv"
	"time"

	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/session/muc/data"
	xmppData "github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
)

func (m *mucManager) receiveClientMessage(stanza *xmppData.ClientMessage) {
	m.log.WithField("stanza", stanza).Debug("handleMUCReceivedClientMessage()")

	if hasSubject(stanza) {
		m.handleDiscussionHistoryOneTime(stanza)
		m.handleSubjectReceived(stanza)
	}

	switch {
	case isDelayedMessage(stanza):
		m.handleMessageReceived(stanza, m.receiveDelayedMessage)
	case isLiveMessage(stanza):
		m.handleMessageReceived(stanza, m.liveMessageReceived)
	case hasMucUserExtension(stanza):
		m.handleMUCUserExtension(stanza)
	}
}

func (m *mucManager) receiveDelayedMessage(roomID jid.Bare, nickname, message string, timestamp time.Time) {
	dh, ok := m.getDiscussionHistory(roomID)
	if !ok {
		dh = m.addNewDiscussionHistory(roomID)
	}

	dh.AddMessage(nickname, message, timestamp)
}

// getDiscussionHistory returns the discussion history for the given room and a boolean indicating if it was found or not
func (m *mucManager) getDiscussionHistory(roomID jid.Bare) (*data.DiscussionHistory, bool) {
	h, ok := m.discussionHistory[roomID]
	return h, ok
}

func (m *mucManager) addNewDiscussionHistory(roomID jid.Bare) *data.DiscussionHistory {
	m.discussionHistoryLock.Lock()
	defer m.discussionHistoryLock.Unlock()

	dh := data.NewDiscussionHistory()
	m.discussionHistory[roomID] = dh

	return dh
}

// The discussion history MUST happen only one time in the events flow of XMPP's MUC
// This should be done in a proper way, maybe in the pending "state machine" pattern
// that we want to implement later, when that happens, this method should be fine
func (m *mucManager) handleDiscussionHistory(stanza *xmppData.ClientMessage) {
	roomID := m.retrieveRoomID(stanza.From)
	dh, exists := m.getDiscussionHistory(roomID)
	if !exists {
		m.log.WithField("room", roomID).Warn("Trying to get a not available discussion history for the given room")
		return
	}

	m.discussionHistoryReceived(roomID, dh)
}

func (m *mucManager) handleSubjectReceived(stanza *xmppData.ClientMessage) {
	l := m.log.WithFields(log.Fields{
		"from": stanza.From,
		"who":  "handleSubjectReceived",
	})

	roomID := m.retrieveRoomID(stanza.From)
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

func (m *mucManager) handleMUCUserExtension(stanza *xmppData.ClientMessage) {
	roomID, ok := jid.TryParseBare(stanza.From)
	if !ok {
		m.log.WithFields(log.Fields{
			"stanza": stanza,
			"who":    "handleMucUserExtension",
		}).Error("Error trying to get room ID from stanza")
		return
	}

	for _, status := range stanza.MucUserExtension.Status {
		switch status.Code {
		case strconv.Itoa(MUCStatusRoomLoggingEnabled):
			m.loggingEnabled(roomID)
		case strconv.Itoa(MUCStatusRoomLoggingDisabled):
			m.loggingDisabled(roomID)
		case strconv.Itoa(MUCStatusRoomNonAnonymous):
			m.nonAnonymousRoom(roomID)
		case strconv.Itoa(MUCStatusRoomSemiAnonymous):
			m.semiAnonymousRoom(roomID)
		case strconv.Itoa(MUCStatusConfigurationChanged):
			m.handleMUCRoomConfigurationChanged(roomID)
		default:
			m.log.WithField("status code:", status.Code).Warn("Unknown status code received")
		}
	}
}

func (m *mucManager) handleMUCRoomConfigurationChanged(roomID jid.Bare) {
	roomInfo := make(chan *muc.RoomListing)
	go m.getRoom(roomID, roomInfo)

	ri := <-roomInfo
	m.roomConfigurationChanged(roomID, ri)
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

func hasMucUserExtension(stanza *xmppData.ClientMessage) bool {
	return stanza.MucUserExtension != nil
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
