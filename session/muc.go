package session

import (
	"sync"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/data"
	xi "github.com/coyim/coyim/xmpp/interfaces"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
)

type mucManager struct {
	log          coylog.Logger
	conn         xi.Conn
	publishEvent func(ev interface{})
	roomManager  *muc.RoomManager
	roomLock     sync.Mutex
	sync.Mutex
}

func newMUCManager(log coylog.Logger, conn xi.Conn, publishEvent func(ev interface{})) *mucManager {
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

func (m *mucManager) handlePresence(stanza *data.ClientPresence) {
	from := jid.ParseFull(stanza.From)

	if stanza.Type == "error" {
		m.handleMUCErrorPresence(from, stanza)
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

		// TODO: is this only sent for own presence, or for changes for other nicknames?
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
			"method":   "selfOccupantUpdate",
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

func (m *mucManager) handleMUCErrorPresence(from jid.Full, stanza *data.ClientPresence) {
	m.publishMUCError(from, stanza.Error)
}

func isMUCPresence(stanza *data.ClientPresence) bool {
	return stanza.MUC != nil
}

func isMUCUserPresence(stanza *data.ClientPresence) bool {
	return stanza.MUCUser != nil
}

func getOccupantPresenceBasedOnItem(nickname jid.Resource, item *data.MUCUserItem) *muc.OccupantPresenceInfo {
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

func getAffiliationBasedOnItem(item *data.MUCUserItem) muc.Affiliation {
	affiliation := "none"
	if item != nil && len(item.Affiliation) > 0 {
		affiliation = item.Affiliation
	}

	return affiliationFromString(affiliation)
}

func affiliationFromString(a string) muc.Affiliation {
	affiliation, _ := muc.AffiliationFromString(a)
	return affiliation
}

func getRoleBasedOnItem(item *data.MUCUserItem) muc.Role {
	role := "none"
	if item != nil && len(item.Role) > 0 {
		role = item.Role
	}

	return roleFromString(role)
}

func roleFromString(r string) muc.Role {
	role, _ := muc.RoleFromString(r)
	return role
}

func getRealJidBasedOnItem(item *data.MUCUserItem) jid.Full {
	if item == nil || len(item.Jid) == 0 {
		return nil
	}

	return jid.ParseFull(item.Jid)
}
