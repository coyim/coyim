package session

import (
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
	roomManager  *muc.RoomManager
	roomLock     sync.Mutex
	sync.Mutex
}

func newMUCManager(log coylog.Logger, conn func() xi.Conn, publishEvent func(ev interface{})) *mucManager {
	m := &mucManager{
		log:          log,
		conn:         conn,
		publishEvent: publishEvent,
		roomManager:  muc.NewRoomManager(),
	}

	return m
}

// NewRoom creates a new muc room and add it to the room manager
func (s *session) NewRoom(roomID jid.Bare) *muc.Room {
	return s.muc.newRoom(roomID)
}

func (m *mucManager) newRoom(roomID jid.Bare) *muc.Room {
	m.roomLock.Lock()
	defer m.roomLock.Unlock()

	room, exists := m.roomManager.GetRoom(roomID)

	if exists {
		return room
	}

	room = muc.NewRoom(roomID)
	m.roomManager.AddRoom(room)

	return room
}

func (m *mucManager) handlePresence(stanza *xmppData.ClientPresence) {
	from := jid.ParseFull(stanza.From)

	if stanza.Type == "error" {
		m.handleMUCErrorPresence(from, stanza.Error)
		return
	}

	roomID := from.Bare()
	occupantPresence := getOccupantPresenceBasedOnItem(from.Resource(), stanza.MUCUser.Item)
	status := mucUserStatuses(stanza.MUCUser.Status)

	isOwnPresence := status.contains(MUCStatusSelfPresence)
	if !isOwnPresence && occupantPresence.RealJid == from {
		isOwnPresence = true
	}

	switch stanza.Type {
	case "unavailable":
		m.handleUnavailablePresence(roomID, occupantPresence, status)
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

func (m *mucManager) handleUnavailablePresence(roomID jid.Bare, op *muc.OccupantPresenceInfo, status mucUserStatuses) {
	switch {
	case status.isEmpty():
		m.log.WithFields(log.Fields{
			"room":        roomID,
			"occupant":    op.Nickname,
			"affiliation": op.Affiliation,
			"role":        op.Role,
		}).Debug("Parameters sent when someone leaves a room")

		m.handleOccupantLeft(roomID, op)

	case status.contains(MUCStatusBanned):
		// We got banned
		m.log.Debug("handleMUCPresence(): MUCStatusBanned")

	case status.contains(MUCStatusNewNickname):
		// Someone has changed its nickname
		m.log.Debug("handleMUCPresence(): MUCStatusNewNickname")

	case status.contains(MUCStatusBecauseKickedFrom):
		// Someone was kicked from the room
		m.log.Debug("handleMUCPresence(): MUCStatusBecauseKickedFrom")

	case status.contains(MUCStatusRemovedBecauseAffiliationChanged):
		// Removed due to an affiliation change
		m.log.Debug("handleMUCPresence(): MUCStatusRemovedBecauseAffiliationChanged")

	case status.contains(MUCStatusRemovedBecauseNotMember):
		// Removed because room is now members-only
		m.log.Debug("handleMUCPresence(): MUCStatusRemovedBecauseNotMember")

	case status.contains(MUCStatusRemovedBecauseShutdown):
		// Removes due to system shutdown
		m.log.Debug("handleMUCPresence(): MUCStatusRemovedBecauseShutdown")
	}
}

func (m *mucManager) handleMUCErrorPresence(from jid.Full, e *xmppData.StanzaError) {
	m.publishMUCError(from, e)
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

func getOccupantPresenceBasedOnItem(nickname jid.Resource, item *xmppData.MUCUserItem) *muc.OccupantPresenceInfo {
	realJid := getRealJidBasedOnItem(item)
	affiliation := getAffiliationBasedOnItem(item)
	role := getRoleBasedOnItem(item)

	op := &muc.OccupantPresenceInfo{
		Nickname:    nickname.String(),
		RealJid:     realJid,
		Affiliation: affiliation,
		Role:        role,
	}

	return op
}

func getAffiliationBasedOnItem(item *xmppData.MUCUserItem) data.Affiliation {
	affiliation := "none"
	if item != nil && len(item.Affiliation) > 0 {
		affiliation = item.Affiliation
	}

	return affiliationFromString(affiliation)
}

func affiliationFromString(a string) data.Affiliation {
	affiliation, _ := data.AffiliationFromString(a)
	return affiliation
}

func getRoleBasedOnItem(item *xmppData.MUCUserItem) data.Role {
	role := "none"
	if item != nil && len(item.Role) > 0 {
		role = item.Role
	}

	return roleFromString(role)
}

func roleFromString(r string) data.Role {
	role, _ := data.RoleFromString(r)
	return role
}

func getRealJidBasedOnItem(item *xmppData.MUCUserItem) jid.Full {
	if item == nil || len(item.Jid) == 0 {
		return nil
	}

	return jid.ParseFull(item.Jid)
}

func (m *mucManager) sendMessage(to, from, body string) error {
	msg := &xmppData.Message{
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
