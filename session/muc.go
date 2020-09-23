package session

import (
	"strconv"
	"sync"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/data"
	xi "github.com/coyim/coyim/xmpp/interfaces"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
)

const (
	// MUCStatusJIDPublic inform user that any occupant is
	// allowed to see the user's full JID
	MUCStatusJIDPublic = 100
	// MUCStatusAffiliationChanged inform user that his or
	// her affiliation changed while not in the room
	MUCStatusAffiliationChanged = 101
	// MUCStatusUnavailableShown inform occupants that room
	// now shows unavailable members
	MUCStatusUnavailableShown = 102
	// MUCStatusUnavailableNotShown inform occupants that room
	// now does not show unavailable members
	MUCStatusUnavailableNotShown = 103
	// MUCStatusConfigurationChanged inform occupants that a
	// non-privacy-related room configuration change has occurred
	MUCStatusConfigurationChanged = 104
	// MUCStatusSelfPresence inform user that presence refers
	// to one of its own room occupants
	MUCStatusSelfPresence = 110
	// MUCStatusRoomLoggingEnabled inform occupants that room
	// logging is now enabled
	MUCStatusRoomLoggingEnabled = 170
	// MUCStatusRoomLoggingDisabled inform occupants that room
	// logging is now disabled
	MUCStatusRoomLoggingDisabled = 171
	// MUCStatusRoomNonAnonymous inform occupants that the room
	// is now non-anonymous
	MUCStatusRoomNonAnonymous = 172
	// MUCStatusRoomSemiAnonymous inform occupants that the room
	// is now semi-anonymous
	MUCStatusRoomSemiAnonymous = 173
	// MUCStatusRoomFullyAnonymous inform occupants that the room
	// is now fully-anonymous
	MUCStatusRoomFullyAnonymous = 174
	// MUCStatusRoomCreated inform user that a new room has
	// been created
	MUCStatusRoomCreated = 201
	// MUCStatusNicknameAssigned inform user that the service has
	// assigned or modified the occupant's roomnick
	MUCStatusNicknameAssigned = 210
	// MUCStatusBanned inform user that he or she has been banned
	// from the room
	MUCStatusBanned = 301
	// MUCStatusNewNickname inform all occupants of new room nickname
	MUCStatusNewNickname = 303
	// MUCStatusBecauseKickedFrom inform user that he or she has been
	// kicked from the room
	MUCStatusBecauseKickedFrom = 307
	// MUCStatusRemovedBecauseAffiliationChanged inform user that
	// he or she is being removed from the room because of an
	// affiliation change
	MUCStatusRemovedBecauseAffiliationChanged = 321
	// MUCStatusRemovedBecauseNotMember inform user that he or she
	// is being removed from the room because the room has been
	// changed to members-only and the user is not a member
	MUCStatusRemovedBecauseNotMember = 322
	// MUCStatusRemovedBecauseShutdown inform user that he or she
	// is being removed from the room because of a system shutdown
	MUCStatusRemovedBecauseShutdown = 332
)

type mucManager struct {
	log          coylog.Logger
	conn         xi.Conn
	publishEvent func(ev interface{})
	roomManager  *muc.RoomManager
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
func (s *session) NewRoom(ident jid.Bare) *muc.Room {
	// TODO: This must be made concurrency safe

	room, exists := s.muc.roomManager.GetRoom(ident)

	if exists {
		return room
	}

	room = muc.NewRoom(ident)
	s.muc.roomManager.AddRoom(room)

	return room
}

func isMUCPresence(stanza *data.ClientPresence) bool {
	return stanza.MUC != nil
}

func isMUCUserPresence(stanza *data.ClientPresence) bool {
	return stanza.MUCUser != nil
}

func getOccupantBasedOnItem(from jid.Full, item *data.MUCUserItem) *mucRoomOccupant {
	affiliation, role := getAffiliationAndRoleBasedOnItem(item)
	realJid := getRealJidBasedOnItem(item)
	return newMUCRoomOccupant(from.Resource(), affiliation, role, realJid)
}

func getAffiliationAndRoleBasedOnItem(item *data.MUCUserItem) (string, string) {
	// TODO: refactor to separate out getters for affiliation and role

	affiliation := "none"
	role := "none"
	if item != nil {
		if len(item.Affiliation) > 0 {
			affiliation = item.Affiliation
		}
		if len(item.Role) > 0 {
			role = item.Role
		}
	}
	return affiliation, role
}

func getRealJidBasedOnItem(item *data.MUCUserItem) jid.Full {
	if item == nil || len(item.Jid) == 0 {
		return nil
	}

	return jid.ParseFull(item.Jid)
}

func (m *mucManager) handlePresence(stanza *data.ClientPresence) {
	from := jid.ParseFull(stanza.From)

	if stanza.Type == "error" {
		m.handleMUCErrorPresence(from, stanza)
		return
	}

	roomID := from.Bare()
	occupant := getOccupantBasedOnItem(from, stanza.MUCUser.Item)
	status := mucUserStatuses(stanza.MUCUser.Status)

	isOwnPresence := status.contains(MUCStatusSelfPresence)
	if !isOwnPresence && occupant.sameFrom(from) {
		isOwnPresence = true
	}

	switch stanza.Type {
	case "unavailable":
		m.handleUnavailablePresence(roomID, occupant, status)
	case "":
		if isOwnPresence {
			m.handleSelfOccupantUpdate(roomID, occupant, status)
		} else {
			m.handleOccupantUpdate(roomID, occupant)
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
func (m *mucManager) handleSelfOccupantUpdate(roomID jid.Bare, occupant *mucRoomOccupant, status mucUserStatuses) {
	m.handleOccupantUpdate(roomID, occupant)
	// TODO: This is a bit confusing since the selfOccupantUpdate method can be called more than
	// once - and it's only the first time it is called that it actually means the person joined
	m.selfOccupantJoin(roomID, occupant)

	if status.contains(MUCStatusRoomLoggingEnabled) {
		m.loggingEnabled(roomID)
	}

	if status.contains(MUCStatusRoomLoggingDisabled) {
		m.loggingDisabled(roomID)
	}
}

func (m *mucManager) selfOccupantJoin(roomID jid.Bare, occupant *mucRoomOccupant) {
	room, exists := m.roomManager.GetRoom(roomID)
	if !exists {
		m.log.WithFields(log.Fields{
			"room":     roomID,
			"occupant": occupant.nickname,
			"method":   "selfOccupantJoin",
		}).Error("trying to join to an unavailable room")
		// TODO: This will only happen when the room disappeared AFTER trying to join, but before we could
		// finish the join. We should figure out the right way of handling this situation
		return
	}

	// TODO: This should be simplified by only using the nickname to identify peers in the roster
	o, exists := room.Roster().GetOccupantByIdentity(roomID.WithResource(jid.NewResource(occupant.nickname)).String())
	if exists {
		room.AddSelfOccupant(o)
		m.selfOccupantJoined(roomID, occupant)
	}
}

func (m *mucManager) handleUnavailablePresence(roomID jid.Bare, occupant *mucRoomOccupant, status mucUserStatuses) {
	switch {
	case status.isEmpty():
		m.log.WithFields(log.Fields{
			"room":        roomID,
			"occupant":    occupant.nickname,
			"affiliation": occupant.affiliation,
			"role":        occupant.role,
		}).Debug("Parameters sent when someone leaves a room")

		m.handleOccupantLeft(roomID, occupant)

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

type mucUserStatuses []data.MUCUserStatus

// contains will return true if the list of MUC user statuses contains ALL of the given argument statuses
func (mus mucUserStatuses) contains(c ...int) bool {
	for _, cc := range c {
		if !mus.containsOne(cc) {
			return false
		}
	}
	return true
}

// containsAny will return true if the list of MUC user statuses contains ANY of the given argument statuses
func (mus mucUserStatuses) containsAny(c ...int) bool {
	for _, cc := range c {
		if mus.containsOne(cc) {
			return true
		}
	}
	return false
}

// containsOne will return true if the list of MUC user statuses contains ONLY the given argument status
func (mus mucUserStatuses) containsOne(c int) bool {
	for _, s := range mus {
		code, _ := strconv.Atoi(s.Code)
		if code == c {
			return true
		}
	}
	return false
}

func (mus mucUserStatuses) isEmpty() bool {
	return len(mus) == 0
}

func hasSomeConferenceService(identities []data.DiscoveryIdentity) bool {
	for _, identity := range identities {
		if identity.Category == "conference" && identity.Type == "text" {
			return true
		}
	}
	return false
}

func (s *session) hasSomeChatService(di data.DiscoveryItem) bool {
	iq, err := s.conn.QueryServiceInformation(di.Jid)
	if err != nil {
		s.log.WithField("jid", di.Jid).WithError(err).Error("Error getting the information query for the service")
		return false
	}
	return hasSomeConferenceService(iq.Identities)
}

// TODO: Ola named this badly. We should change the name
type chatServiceReceivalContext struct {
	sync.RWMutex

	resultsChannel chan jid.Domain
	errorChannel   chan error

	s *session
}

func (c *chatServiceReceivalContext) end() {
	c.Lock()
	defer c.Unlock()
	if c.resultsChannel != nil {
		close(c.resultsChannel)
		close(c.errorChannel)
		c.resultsChannel = nil
		c.errorChannel = nil
	}
}

func (s *session) createChatServiceReceivalContext() *chatServiceReceivalContext {
	result := &chatServiceReceivalContext{}

	result.resultsChannel = make(chan jid.Domain)
	result.errorChannel = make(chan error)
	result.s = s

	return result
}

func (c *chatServiceReceivalContext) fetchChatServices(server jid.Domain) {
	defer c.end()
	items, err := c.s.conn.QueryServiceItems(server.String())
	if err != nil {
		c.RLock()
		defer c.RUnlock()
		if c.errorChannel != nil {
			c.errorChannel <- err
		}
		return
	}
	for _, item := range items.DiscoveryItems {
		if c.s.hasSomeChatService(item) {
			c.RLock()
			defer c.RUnlock()
			if c.resultsChannel == nil {
				return
			}
			c.resultsChannel <- jid.Parse(item.Jid).Host()
		}
	}
}

// TODO: Maybe move the below function and related functionality to its own file?

// GetChatServices offers the chat services from a xmpp server.
func (s *session) GetChatServices(server jid.Domain) (<-chan jid.Domain, <-chan error, func()) {
	ctx := s.createChatServiceReceivalContext()
	go ctx.fetchChatServices(server)
	return ctx.resultsChannel, ctx.errorChannel, ctx.end
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

func (m *mucManager) receivedClientMessage(stanza *data.ClientMessage) {
	m.log.WithField("stanza", stanza).Debug("handleMUCReceivedClientMessage()")

	if isLiveMessage(stanza) {
		from := jid.ParseFull(stanza.From)
		room := from.Bare()
		nickname := from.Resource().String()
		message := stanza.Body
		subject := ""

		if stanza.Subject != nil {
			subject = stanza.Subject.Text
		}

		m.log.WithFields(log.Fields{
			"room":     room,
			"message":  message,
			"subject":  subject,
			"nickname": nickname,
		}).Info("MUC message received")

		m.messageReceived(room, nickname, subject, message)
	}
}
