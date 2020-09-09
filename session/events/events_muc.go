package events

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
)

// MUC is a marker interface that is used to differentiate MUC events
type MUC interface {
	mucEventMarkerFunction()
	WhichRoom() jid.Bare
}

// MUCError contains information about a MUC-related
// error event
type MUCError struct {
	ErrorType MUCErrorType
	Room      jid.Bare
	Nickname  string
}

func (MUCError) mucEventMarkerFunction() {}

// WhichRoom implements the MUC interface
func (m MUCError) WhichRoom() jid.Bare { return m.Room }

// MUCRoomCreated contains event information about
// the created room
type MUCRoomCreated struct {
	Room jid.Bare
}

func (MUCRoomCreated) mucEventMarkerFunction() {}

// WhichRoom implements the MUC interface
func (m MUCRoomCreated) WhichRoom() jid.Bare { return m.Room }

// MUCRoomRenamed contains event information about
// the renamed room's nickname
type MUCRoomRenamed struct {
	Room        jid.Bare
	NewNickname string
}

func (MUCRoomRenamed) mucEventMarkerFunction() {}

// WhichRoom implements the MUC interface
func (m MUCRoomRenamed) WhichRoom() jid.Bare { return m.Room }

// MUCOccupant contains basic information about
// any room's occupant
type MUCOccupant struct {
	Room     jid.Bare
	Nickname string
	RealJid  jid.Full
}

func (MUCOccupant) mucEventMarkerFunction() {}

// WhichRoom implements the MUC interface
func (m MUCOccupant) WhichRoom() jid.Bare { return m.Room }

// MUCOccupantUpdated contains information about
// the updated occupant in a room
type MUCOccupantUpdated struct {
	MUCOccupant
	Affiliation muc.Affiliation
	Role        muc.Role
}

func (MUCOccupantUpdated) mucEventMarkerFunction() {}

// WhichRoom implements the MUC interface
func (m MUCOccupantUpdated) WhichRoom() jid.Bare { return m.Room }

// MUCOccupantJoined contains information about
// the occupant that has joined to room
type MUCOccupantJoined struct {
	MUCOccupantUpdated
	Status string
}

func (MUCOccupantJoined) mucEventMarkerFunction() {}

// WhichRoom implements the MUC interface
func (m MUCOccupantJoined) WhichRoom() jid.Bare { return m.Room }

// MUCOccupantLeft contains information about
// the occupant that has left the room
type MUCOccupantLeft struct {
	MUCOccupant
	Affiliation muc.Affiliation
	Role        muc.Role
}

func (MUCOccupantLeft) mucEventMarkerFunction() {}

// WhichRoom implements the MUC interface
func (m MUCOccupantLeft) WhichRoom() jid.Bare { return m.Room }

// MUCMessageReceived contains information about
// the message received
type MUCMessageReceived struct {
	Room     jid.Bare
	Nickname jid.Resource
	Message  string
}

func (MUCMessageReceived) mucEventMarkerFunction() {}

// WhichRoom implements the MUC interface
func (m MUCMessageReceived) WhichRoom() jid.Bare { return m.Room }

// MUCLoggingEnabled signifies that logging has been turned on from the room
type MUCLoggingEnabled struct {
	Room jid.Bare
}

func (MUCLoggingEnabled) mucEventMarkerFunction() {}

// WhichRoom implements the MUC interface
func (m MUCLoggingEnabled) WhichRoom() jid.Bare { return m.Room }

// MUCLoggingDisabled signifies that logging has been turned off from the room
type MUCLoggingDisabled struct {
	Room jid.Bare
}

func (MUCLoggingDisabled) mucEventMarkerFunction() {}

// WhichRoom implements the MUC interface
func (m MUCLoggingDisabled) WhichRoom() jid.Bare { return m.Room }

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
