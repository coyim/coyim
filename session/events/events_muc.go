package events

import (
	"github.com/coyim/coyim/xmpp/jid"
)

// MUC is for publishing MUC-related session events
type MUC struct {
	From jid.Full
	Room jid.Bare
	// Contains information related to any MUC event
	Info interface{}
}

// MUCError contains information about a MUC-related
// error event
type MUCError struct {
	ErrorType MUCErrorType
	Room      jid.Bare
}

// MUCRoomCreated contains event information about
// the created room
type MUCRoomCreated struct {
	MUC
}

// MUCRoomRenamed contains event information about
// the renamed room's nickname
type MUCRoomRenamed struct {
	MUC
}

// MUCOccupant contains basic information about
// any room's occupant
type MUCOccupant struct {
	MUC
	Nickname jid.Resource
	Jid      jid.Full
}

// MUCOccupantUpdated contains information about
// the updated occupant in a room
type MUCOccupantUpdated struct {
	MUCOccupant
	Affiliation string
	Role        string
}

// MUCOccupantJoined contains information about
// the occupant that has joined to room
type MUCOccupantJoined struct {
	MUCOccupantUpdated
	Status string
}

// MUCOccupantLeft contains information about
// the occupant that has left the room
type MUCOccupantLeft struct {
	MUCOccupant
	Affiliation string
	Role        string
}

// MUCMessageReceived contains information about
// the message received
type MUCMessageReceived struct {
	MUC
	Nickname jid.Resource
	Message  string
}

// MUCLoggingEnabled signifies that logging has been turned on from the room
type MUCLoggingEnabled struct {
	MUC
}

// MUCLoggingDisabled signifies that logging has been turned off from the room
type MUCLoggingDisabled struct {
	MUC
}

// MUCErrorType represents the type of MUC error event
type MUCErrorType EventType

// MUC event types
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
