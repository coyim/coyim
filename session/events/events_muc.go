package events

import (
	"time"

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
	// MUCNoError is a special type that can be used as a "no error"
	// flag inside the logic of the MUC implementation
	MUCNoError MUCErrorType = iota

	MUCNotAuthorized
	MUCForbidden
	MUCItemNotFound
	MUCNotAllowed
	MUCNotAcceptable
	MUCRegistrationRequired
	MUCConflict
	MUCServiceUnavailable

	MUCMessageForbidden
	MUCMessageNotAcceptable
)

// MUCError contains information about a MUC-related
// error event
type MUCError struct {
	ErrorType MUCErrorType
	Room      jid.Bare
	Nickname  string
}

// MUCRoom contains information about the room id
type MUCRoom struct {
	Room jid.Bare
}

// MUCRoomCreated contains event information about
// the created room
type MUCRoomCreated struct {
	MUCRoom
}

// MUCRoomDestroyed contains event information about
// the destroyed room
type MUCRoomDestroyed struct {
	Reason          string
	AlternativeRoom jid.Bare
	Password        string
}

// MUCRoomRenamed contains event information about
// the renamed room's nickname
type MUCRoomRenamed struct {
	NewNickname string
}

// MUCOccupant contains basic information about
// any room's occupant
type MUCOccupant struct {
	Nickname      string
	RealJid       jid.Full
	Status        string
	StatusMessage string
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

// MUCOccupantRemoved contains information related to member removedcontains information related to self occupant which has been removed
type MUCOccupantRemoved struct {
	MUCOccupant
}

// MUCSelfOccupantRemoved contains information related to self occupant which has been removed
type MUCSelfOccupantRemoved struct{}

// MUCMessageReceived represents a received groupchat message
type MUCMessageReceived struct {
	Nickname  string
	Message   string
	Timestamp time.Time
}

// MUCLiveMessageReceived contains information about the received live message
type MUCLiveMessageReceived struct {
	MUCMessageReceived
}

// MUCDelayedMessageReceived contains information about the received delayed message
type MUCDelayedMessageReceived struct {
	MUCMessageReceived
}

// MUCDiscussionHistoryReceived contains information about full discussion history
type MUCDiscussionHistoryReceived struct {
	History *data.DiscussionHistory
}

// MUCSubjectUpdated contains the room subject will be updated
type MUCSubjectUpdated struct {
	Nickname string
	Subject  string
}

// MUCSubjectReceived contains the room subject received
type MUCSubjectReceived struct {
	Subject string
}

// MUCLoggingEnabled signifies that logging has been turned on from the room
type MUCLoggingEnabled struct{}

// MUCLoggingDisabled signifies that logging has been turned off from the room
type MUCLoggingDisabled struct{}

// MUCRoomAnonymityChanged contains information regarding to if the the room is semi or non anonymous
type MUCRoomAnonymityChanged struct {
	AnonymityLevel string
}

// MUCRoomDiscoInfoReceived contains information of the received room disco info
type MUCRoomDiscoInfoReceived struct {
	DiscoInfo data.RoomDiscoInfo
}

// MUCRoomConfigTimeout indicates that the room listing request has timeout
type MUCRoomConfigTimeout struct{}

// MUCRoomConfigChanged signifies that room configuration changed
type MUCRoomConfigChanged struct {
	Changes   []data.RoomConfigType
	DiscoInfo data.RoomDiscoInfo
}

// MUCOccupantAffiliationRoleUpdated signifies that an occupant affiliation and role was updated
type MUCOccupantAffiliationRoleUpdated struct {
	AffiliationRoleUpdate data.AffiliationRoleUpdate
}

// MUCSelfOccupantAffiliationRoleUpdated signifies that the self-occupant affiliation and role was updated
type MUCSelfOccupantAffiliationRoleUpdated struct {
	AffiliationRoleUpdate data.SelfAffiliationRoleUpdate
}

// MUCOccupantAffiliationUpdated signifies that an occupant affiliation was updated
type MUCOccupantAffiliationUpdated struct {
	AffiliationUpdate data.AffiliationUpdate
}

// MUCSelfOccupantAffiliationUpdated signifies that the self-occupant affiliation was updated
type MUCSelfOccupantAffiliationUpdated struct {
	AffiliationUpdate data.SelfAffiliationUpdate
}

// MUCOccupantRoleUpdated signifies that an occupant role was updated
type MUCOccupantRoleUpdated struct {
	RoleUpdate data.RoleUpdate
}

// MUCSelfOccupantRoleUpdated signifies that the self-occupant role was updated
type MUCSelfOccupantRoleUpdated struct {
	RoleUpdate data.SelfRoleUpdate
}

// MUCOccupantKicked contains information about the occupant kicked
type MUCOccupantKicked struct {
	RoleUpdate data.RoleUpdate
}

// MUCSelfOccupantKicked contains information about the self-occupant kicked
type MUCSelfOccupantKicked struct {
	RoleUpdate data.SelfRoleUpdate
}

func (MUCError) markAsMUCEventTypeInterface()                              {}
func (MUCRoom) markAsMUCEventTypeInterface()                               {}
func (MUCRoomCreated) markAsMUCEventTypeInterface()                        {}
func (MUCRoomDestroyed) markAsMUCEventTypeInterface()                      {}
func (MUCRoomRenamed) markAsMUCEventTypeInterface()                        {}
func (MUCOccupant) markAsMUCEventTypeInterface()                           {}
func (MUCOccupantUpdated) markAsMUCEventTypeInterface()                    {}
func (MUCOccupantJoined) markAsMUCEventTypeInterface()                     {}
func (MUCSelfOccupantJoined) markAsMUCEventTypeInterface()                 {}
func (MUCOccupantLeft) markAsMUCEventTypeInterface()                       {}
func (MUCLiveMessageReceived) markAsMUCEventTypeInterface()                {}
func (MUCDelayedMessageReceived) markAsMUCEventTypeInterface()             {}
func (MUCSubjectUpdated) markAsMUCEventTypeInterface()                     {}
func (MUCSubjectReceived) markAsMUCEventTypeInterface()                    {}
func (MUCLoggingEnabled) markAsMUCEventTypeInterface()                     {}
func (MUCLoggingDisabled) markAsMUCEventTypeInterface()                    {}
func (MUCRoomAnonymityChanged) markAsMUCEventTypeInterface()               {}
func (MUCDiscussionHistoryReceived) markAsMUCEventTypeInterface()          {}
func (MUCRoomDiscoInfoReceived) markAsMUCEventTypeInterface()              {}
func (MUCRoomConfigTimeout) markAsMUCEventTypeInterface()                  {}
func (MUCRoomConfigChanged) markAsMUCEventTypeInterface()                  {}
func (MUCOccupantRemoved) markAsMUCEventTypeInterface()                    {}
func (MUCSelfOccupantRemoved) markAsMUCEventTypeInterface()                {}
func (MUCOccupantAffiliationUpdated) markAsMUCEventTypeInterface()         {}
func (MUCSelfOccupantAffiliationUpdated) markAsMUCEventTypeInterface()     {}
func (MUCOccupantRoleUpdated) markAsMUCEventTypeInterface()                {}
func (MUCSelfOccupantRoleUpdated) markAsMUCEventTypeInterface()            {}
func (MUCOccupantAffiliationRoleUpdated) markAsMUCEventTypeInterface()     {}
func (MUCSelfOccupantAffiliationRoleUpdated) markAsMUCEventTypeInterface() {}
func (MUCOccupantKicked) markAsMUCEventTypeInterface()                     {}
func (MUCSelfOccupantKicked) markAsMUCEventTypeInterface()                 {}
