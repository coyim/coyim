package session

import (
	"errors"
	"sync"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/session/muc/data"
	xmppData "github.com/coyim/coyim/xmpp/data"
	xi "github.com/coyim/coyim/xmpp/interfaces"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
)

type mucManager struct {
	log          coylog.Logger
	conn         func() xi.Conn
	publishEvent func(ev interface{})

	roomInfos     map[jid.Bare]*muc.RoomListing
	roomInfosLock sync.Mutex

	roomManager *muc.RoomManager
	roomLock    sync.Mutex

	roomConfigChangesHandlers     map[int]func(jid.Bare)
	roomConfigChangesHandlersLock sync.Mutex
}

func newMUCManager(log coylog.Logger, conn func() xi.Conn, publishEvent func(ev interface{})) *mucManager {
	m := &mucManager{
		log:          log,
		conn:         conn,
		publishEvent: publishEvent,
		roomManager:  muc.NewRoomManager(),
		roomInfos:    make(map[jid.Bare]*muc.RoomListing),
	}

	return m
}

// GetRoom will look up in the room manager for the room
// with the given ID and returns it
func (s *session) GetRoom(roomID jid.Bare) (*muc.Room, bool) {
	return s.muc.getRoom(roomID)
}

// NewRoom creates a new muc room and add it to the room manager
func (s *session) NewRoom(roomID jid.Bare) *muc.Room {
	return s.muc.newRoom(roomID)
}

func (m *mucManager) getRoom(roomID jid.Bare) (*muc.Room, bool) {
	m.roomLock.Lock()
	defer m.roomLock.Unlock()

	room, ok := m.roomManager.GetRoom(roomID)
	return room, ok
}

func (m *mucManager) newRoom(roomID jid.Bare) *muc.Room {
	if room, ok := m.getRoom(roomID); ok {
		return room
	}

	m.roomLock.Lock()
	defer m.roomLock.Unlock()

	room := muc.NewRoom(roomID)
	m.roomManager.AddRoom(room)

	return room
}

func (m *mucManager) addRoomInfo(roomID jid.Bare, rl *muc.RoomListing) {
	m.roomInfosLock.Lock()
	defer m.roomInfosLock.Unlock()

	m.roomInfos[roomID] = rl
	if room, exists := m.roomManager.GetRoom(roomID); exists {
		room.SetProperties(rl.RoomDiscoInfo)
	}
}

func (m *mucManager) getRoomInfo(roomID jid.Bare) (*muc.RoomListing, bool) {
	m.roomInfosLock.Lock()
	defer m.roomInfosLock.Unlock()

	rl, ok := m.roomInfos[roomID]

	return rl, ok
}

func (m *mucManager) handlePresence(stanza *xmppData.ClientPresence) {
	from := jid.ParseFull(stanza.From)
	roomID := from.Bare()
	occupantPresence := getOccupantPresenceBasedOnStanza(from.Resource(), stanza)
	status := mucUserStatuses(stanza.MUCUser.Status)

	isOwnPresence := status.contains(MUCStatusSelfPresence)
	if !isOwnPresence && occupantPresence.RealJid == from {
		isOwnPresence = true
	}

	switch stanza.Type {
	case "unavailable":
		m.handleUnavailablePresence(roomID, occupantPresence, status, stanza)
	case "":
		if isOwnPresence {
			m.handleSelfOccupantUpdate(roomID, occupantPresence, status)
		} else {
			m.handleOccupantUpdate(roomID, occupantPresence)
		}

		if status.contains(MUCStatusNicknameAssigned) {
			m.roomRenamed(roomID)
		}
	}
}

// handleSelfOccupantUpdate can happen several times - every time a status code update is
// changed, or role or affiliation is updated, this can lead to the method being called.
// For now, it will generate a event about joining, but this should be cleaned up and fixed
func (m *mucManager) handleSelfOccupantUpdate(roomID jid.Bare, op *muc.OccupantPresenceInfo, status mucUserStatuses) {
	m.selfOccupantUpdate(roomID, op)

	if status.contains(MUCStatusRoomLoggingEnabled) {
		m.loggingEnabled(roomID)
	}

	if status.contains(MUCStatusRoomLoggingDisabled) {
		m.loggingDisabled(roomID)
	}
}

func (m *mucManager) selfOccupantUpdate(roomID jid.Bare, op *muc.OccupantPresenceInfo) {
	room, exists := m.roomManager.GetRoom(roomID)
	if !exists {
		m.log.WithFields(log.Fields{
			"room":     roomID,
			"occupant": op.Nickname,
			"who":      "selfOccupantUpdate",
		}).Error("trying to join to an unavailable room")
		// TODO: This will only happen when the room disappeared AFTER trying to join, but before we could
		// finish the join. We should figure out the right way of handling this situation
		return
	}

	exists = m.existOccupantInRoster(room, op.Nickname)

	o := m.updateOccupantAndReturn(room, op)

	if !exists {
		room.AddSelfOccupant(o)
		m.selfOccupantJoined(roomID, op)
	}
}

func (m *mucManager) existOccupantInRoster(room *muc.Room, nickname string) bool {
	_, exist := room.Roster().GetOccupant(nickname)
	return exist
}

func (m *mucManager) updateOccupantAndReturn(room *muc.Room, op *muc.OccupantPresenceInfo) *muc.Occupant {
	m.handleOccupantUpdate(room.ID, op)
	o, _ := room.Roster().GetOccupant(op.Nickname)
	return o
}

func (m *mucManager) handleUnavailablePresence(roomID jid.Bare, op *muc.OccupantPresenceInfo, status mucUserStatuses, stanza *xmppData.ClientPresence) {
	switch {
	case status.isEmpty():
		m.handleOccupantLeft(roomID, op)

	case status.contains(MUCStatusBanned):
		// We got banned
		m.handleOccupantBanned(roomID, op)

	case status.contains(MUCStatusNewNickname):
		// Someone has changed its nickname
		m.log.Debug("handleMUCPresence(): MUCStatusNewNickname")

	case status.contains(MUCStatusBecauseKickedFrom):
		m.handleOccupantKick(roomID, op)

	case status.contains(MUCStatusRemovedBecauseAffiliationChanged):
		// Removed due to an affiliation change
		m.log.Debug("handleMUCPresence(): MUCStatusRemovedBecauseAffiliationChanged")

	case status.contains(MUCStatusRemovedBecauseNotMember):
		// Removed because room is now members-only
		m.handleNonMembersRemoved(roomID, op)

	case status.contains(MUCStatusRemovedBecauseShutdown):
		// Removes due to system shutdown
		m.log.Debug("handleMUCPresence(): MUCStatusRemovedBecauseShutdown")

	default:
		m.handleOccupantUnavailable(roomID, op, stanza.MUCUser)
	}
}

func (m *mucManager) handleMUCErrorPresence(from string, e *xmppData.StanzaError) {
	if jid, ok := jid.TryParseFull(from); ok {
		m.publishMUCError(jid.Bare(), jid.Resource().String(), e)
		return
	}

	if roomID, ok := jid.TryParseBare(from); ok {
		m.publishMUCError(roomID, "", e)
		return
	}

	m.log.WithFields(log.Fields{
		"from":         from,
		"where":        "handleMUCErrorPresence",
		"stanza error": e,
	}).Error(errors.New("error trying to publish an error received. It couldn't parse room identifier"))
}

func (m *mucManager) handleMUCErrorMessage(roomID jid.Bare, e *xmppData.StanzaError) {
	m.publishMUCMessageError(roomID, e)
}

func isMUCPresence(stanza *xmppData.ClientPresence) bool {
	return stanza.MUC != nil
}

func isMUCUserPresence(stanza *xmppData.ClientPresence) bool {
	return stanza.MUCUser != nil
}

func getOccupantPresenceBasedOnStanza(nickname jid.Resource, stanza *xmppData.ClientPresence) *muc.OccupantPresenceInfo {
	item := stanza.MUCUser.Item

	opi := &muc.OccupantPresenceInfo{
		Nickname:        nickname.String(),
		AffiliationRole: getAffiliationRoleBasedOnItem(item),
		Status:          stanza.Show,
		StatusMessage:   stanza.Status,
	}

	if realJid, ok := getRealJidBasedOnItem(item); ok {
		opi.RealJid = realJid
	}

	return opi
}

func getAffiliationRoleBasedOnItem(item *xmppData.MUCUserItem) *muc.OccupantAffiliationRole {
	a := data.AffiliationNone
	r := data.RoleNone

	if item != nil {
		if item.Affiliation != "" {
			a = item.Affiliation
		}

		if item.Role != "" {
			r = item.Role
		}
	}

	return &muc.OccupantAffiliationRole{
		Affiliation: affiliationFromString(a),
		Role:        roleFromString(r),
		Actor:       getActorNicknameBasedOnItem(item),
		Reason:      item.Reason,
	}
}

func getActorNicknameBasedOnItem(item *xmppData.MUCUserItem) string {
	if item != nil && item.Actor != nil {
		return item.Actor.Nickname
	}
	return ""
}

func getReasonBasedOnItem(item *xmppData.MUCUserItem) string {
	if item != nil {
		return item.Reason
	}
	return ""
}

func affiliationFromString(a string) data.Affiliation {
	affiliation, _ := data.AffiliationFromString(a)
	return affiliation
}

func roleFromString(r string) data.Role {
	role, _ := data.RoleFromString(r)
	return role
}

func getRealJidBasedOnItem(item *xmppData.MUCUserItem) (jid.Full, bool) {
	if item != nil {
		return jid.TryParseFull(item.Jid)
	}
	return nil, false
}

func (m *mucManager) sendMessage(to, from, body string) error {
	msg := &xmppData.ClientMessage{
		To:   to,
		From: from,
		Body: body,
		Type: "groupchat",
	}

	err := m.conn().SendMessage(msg)
	if err != nil {
		m.log.WithError(err).Error("The message could not be send")
		return err
	}

	return nil
}

func (m *mucManager) updateRoomSubject(roomID jid.Bare, actor, subject string) error {
	err := m.conn().SendMessage(&xmppData.ClientMessage{
		To:   roomID.String(),
		From: actor,
		Subject: &xmppData.Subject{
			Text: subject,
		},
		Type: "groupchat",
	})

	if err != nil {
		m.log.WithError(err).Error("The message could not be send")
		return err
	}

	return nil
}

func (m *mucManager) retrieveRoomID(from string, who string) jid.Bare {
	roomID, ok := jid.TryParseBare(from)
	if !ok {
		m.log.WithFields(log.Fields{
			"who":  who,
			"from": from,
		}).Error("Error trying to get the room ID from stanza's from")
		return nil
	}

	return roomID
}

func (m *mucManager) retrieveRoomIDAndNickname(from string) (jid.Bare, string) {
	f, ok := jid.TryParseFull(from)
	if !ok {
		m.log.WithFields(log.Fields{
			"who":  "retrieveRoomIDAndNickname",
			"from": from,
		}).Error("Error trying to get the room ID and the nickname from stanza's from")
		return nil, ""
	}

	return f.Bare(), f.Resource().String()
}
