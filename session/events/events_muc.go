package events

import (
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/coyim/xmpp/jid"
)

// MUC is a marker interface that is used to differentiate MUC "events"
type MUC interface {
	markAsMUCEventTypeInterface()
}

// MUCErrorType represents the type of MUC error event
type MUCErrorType EventType

// MUC error event types
const (
	MUCNotAuthorized MUCErrorType = iota
	MUCForbidden
	MUCItemNotFound
	MUCNotAllowed
	MUCNotAcceptable
	MUCRegistrationRequired
	MUCConflict
	MUCServiceUnavailable
)

// MUCError contains information about a MUC-related
// error event
type MUCError struct {
	ErrorType MUCErrorType
	Room      jid.Bare
	Nickname  string
}

// MUCRoomCreated contains event information about
// the created room
type MUCRoomCreated struct {
	Room jid.Bare
}

// MUCRoomRenamed contains event information about
// the renamed room's nickname
type MUCRoomRenamed struct {
	NewNickname string
}

// MUCOccupant contains basic information about
// any room's occupant
type MUCOccupant struct {
	Nickname string
	RealJid  jid.Full
}

// TODO: Updated and Joined events need to have Status and StatusText fields

// MUCOccupantUpdated contains information about
// the updated occupant in a room
type MUCOccupantUpdated struct {
	MUCOccupant
	Affiliation data.Affiliation
	Role        data.Role
}

// MUCOccupantJoined contains information about
// the occupant that has joined to room
type MUCOccupantJoined struct {
	MUCOccupantUpdated
	Status string
}

// MUCSelfOccupantJoined contains information about
// the occupant that has joined to room
type MUCSelfOccupantJoined struct {
	MUCOccupantJoined
}

// MUCOccupantLeft contains information about
// the occupant that has left the room
type MUCOccupantLeft struct {
	MUCOccupant
	Affiliation data.Affiliation
	Role        data.Role
}

// MUCMessageReceived contains information about
// the message received
type MUCMessageReceived struct {
	Nickname string
	Subject  string
	Message  string
}

// MUCSubjectReceived contains the room subject received
type MUCSubjectReceived struct {
	Subject string
}

// MUCLoggingEnabled signifies that logging has been turned on from the room
type MUCLoggingEnabled struct{}

// MUCLoggingDisabled signifies that logging has been turned off from the room
type MUCLoggingDisabled struct{}

func (MUCError) markAsMUCEventTypeInterface()              {}
func (MUCRoomCreated) markAsMUCEventTypeInterface()        {}
func (MUCRoomRenamed) markAsMUCEventTypeInterface()        {}
func (MUCOccupant) markAsMUCEventTypeInterface()           {}
func (MUCOccupantUpdated) markAsMUCEventTypeInterface()    {}
func (MUCOccupantJoined) markAsMUCEventTypeInterface()     {}
func (MUCSelfOccupantJoined) markAsMUCEventTypeInterface() {}
func (MUCOccupantLeft) markAsMUCEventTypeInterface()       {}
func (MUCMessageReceived) markAsMUCEventTypeInterface()    {}
func (MUCSubjectReceived) markAsMUCEventTypeInterface()    {}
func (MUCLoggingEnabled) markAsMUCEventTypeInterface()     {}
func (MUCLoggingDisabled) markAsMUCEventTypeInterface()    {}
