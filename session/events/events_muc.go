package events

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
)

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
	Affiliation muc.Affiliation
	Role        muc.Role
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
	Affiliation muc.Affiliation
	Role        muc.Role
}

// MUCMessageReceived contains information about
// the message received
type MUCMessageReceived struct {
	Nickname string
	Subject  string
	Message  string
}

// MUCLoggingEnabled signifies that logging has been turned on from the room
type MUCLoggingEnabled struct{}

// MUCLoggingDisabled signifies that logging has been turned off from the room
type MUCLoggingDisabled struct{}

// TODO: Having a marker method for an interface implementation
// exported is not very nice at all.

// MarkAsMUCInterface implements the MUC interface
func (MUCError) MarkAsMUCInterface() {}

// MarkAsMUCInterface implements the MUC interface
func (MUCRoomCreated) MarkAsMUCInterface() {}

// MarkAsMUCInterface implements the MUC interface
func (MUCRoomRenamed) MarkAsMUCInterface() {}

// MarkAsMUCInterface implements the MUC interface
func (MUCOccupant) MarkAsMUCInterface() {}

// MarkAsMUCInterface implements the MUC interface
func (MUCOccupantUpdated) MarkAsMUCInterface() {}

// MarkAsMUCInterface implements the MUC interface
func (MUCOccupantJoined) MarkAsMUCInterface() {}

// MarkAsMUCInterface implements the MUC interface
func (MUCSelfOccupantJoined) MarkAsMUCInterface() {}

// MarkAsMUCInterface implements the MUC interface
func (MUCOccupantLeft) MarkAsMUCInterface() {}

// MarkAsMUCInterface implements the MUC interface
func (MUCMessageReceived) MarkAsMUCInterface() {}

// MarkAsMUCInterface implements the MUC interface
func (MUCLoggingEnabled) MarkAsMUCInterface() {}

// MarkAsMUCInterface implements the MUC interface
func (MUCLoggingDisabled) MarkAsMUCInterface() {}
