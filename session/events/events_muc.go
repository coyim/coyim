package events

import (
	"github.com/coyim/coyim/xmpp/jid"
)

// MUC is for publishing MUC-related session events
type MUC struct {
	From      jid.Bare
	EventType EventType
	// Contains information related to any MUC event
	Info interface{}
}

// MUCError contains information about a MUC-related
// error event
type MUCError struct{}

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
	Jid      jid.WithResource
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

// MUCOccupantExited contains information about
// the occupant that has exited from a room
type MUCOccupantExited struct {
	MUCOccupant
}

// MUC event types
const (
	MUCNotAuthorized EventType = iota
	MUCForbidden
	MUCItemNotFound
	MUCNotAllowed
	MUCNotAceptable
	MUCRegistrationRequired
	MUCConflict
	MUCServiceUnavailable
)
