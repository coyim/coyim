package session

import (
	"strconv"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
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
	s *session
}

func newMUCManager(s *session) *mucManager {
	m := &mucManager{
		s: s,
	}
	return m
}

func (m *mucManager) isMUCPresence(stanza *data.ClientPresence) bool {
	return stanza.MUC != nil
}

func (m *mucManager) isMUCUserPresence(stanza *data.ClientPresence) bool {
	return stanza.MUCUser != nil
}

func (m *mucManager) handleMUCPresence(stanza *data.ClientPresence) {
	from := jid.Parse(stanza.From).(jid.WithResource)

	room, occupant := from.PotentialSplit()
	status := stanza.MUCUser.Status

	isOwnPresence := userStatusContains(status, MUCStatusSelfPresence)
	if !isOwnPresence && stanza.MUCUser.Item.Jid == from.String() {
		isOwnPresence = true
	}

	switch stanza.Type {
	case "unavailable":
		m.mucOccupantExit(room, occupant)

		if userStatusContains(status, MUCStatusBanned) {
			// We got banned
			log.Debug("handleMUCPresence(): MUCStatusBanned")
		}

		if userStatusContains(status, MUCStatusNewNickname) {
			// Someone has changed its nickname
			log.Debug("handleMUCPresence(): MUCStatusNewNickname")
		}

		if userStatusContains(status, MUCStatusBecauseKickedFrom) {
			// Someone was kicked from the room
			log.Debug("handleMUCPresence(): MUCStatusBecauseKickedFrom")
		}

		if userStatusContains(status, MUCStatusRemovedBecauseAffiliationChanged) {
			// Removed due to an affiliation change
			log.Debug("handleMUCPresence(): MUCStatusRemovedBecauseAffiliationChanged")
		}

		if userStatusContains(status, MUCStatusRemovedBecauseNotMember) {
			// Removed because room is now members-only
			log.Debug("handleMUCPresence(): MUCStatusRemovedBecauseNotMember")
		}

		if userStatusContains(status, MUCStatusRemovedBecauseShutdown) {
			// Removes due to system shutdown
			log.Debug("handleMUCPresence(): MUCStatusRemovedBecauseShutdown")
		}
	case "":
		affiliation := stanza.MUCUser.Item.Affiliation
		role := stanza.MUCUser.Item.Role

		if isOwnPresence {
			m.mucOccupantJoined(room, occupant, affiliation, role)
		} else {
			m.mucOccupantUpdate(room, occupant, affiliation, role)
		}

		if userStatusContains(status, MUCStatusNicknameAssigned) {
			m.mucRoomRenamed(room)
		}
	}
}

func userStatusContains(status []data.MUCUserStatus, c int) bool {
	for _, s := range status {
		code, _ := strconv.Atoi(s.Code)
		if code == code {
			return true
		}
	}
	return false
}

func (s *session) hasSomeConferenceService(identities []data.DiscoveryIdentity) bool {
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
		s.log.WithError(err).Error("Error getting the information query for the service:", di.Jid)
		return false
	}
	return s.hasSomeConferenceService(iq.Identities)
}

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

//GetChatServices offers the chat services from a xmpp server.
func (s *session) GetChatServices(server jid.Domain) (<-chan jid.Domain, <-chan error, func()) {
	ctx := s.createChatServiceReceivalContext()
	go ctx.fetchChatServices(server)
	return ctx.resultsChannel, ctx.errorChannel, ctx.end
}
