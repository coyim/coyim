package session

import (
	"github.com/coyim/coyim/xmpp/data"
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
		if s.Code == c {
			return true
		}
	}
	return false
}

func (mus mucUserStatuses) isEmpty() bool {
	return len(mus) == 0
}
