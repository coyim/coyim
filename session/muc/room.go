package muc

import (
	"github.com/coyim/coyim/session/events"

	"github.com/coyim/coyim/xmpp/jid"
)

// Room represents a multi user chat room that a session is currently connected to.
// It contains information about the room configuration itself, and the participants in the room
type Room struct {
	ID jid.Bare

	subject           string
	subjectWasUpdated bool

	selfOccupant *Occupant
	roster       *RoomRoster

	observers *roomObservers

	// Configuration options:

	MaxHistoryFetch       int
	AllowPrivateMessages  string // This can be 'anyone', 'participants', 'moderators', 'none'
	AllowInvites          bool
	ChangeSubject         bool
	EnableLogging         bool
	GetMemberList         []string // This is a list of the roles that can get the member list, 'moderator', 'participant' or 'visitor'
	Language              string
	PubSub                string
	MaxUsers              int
	MembersOnly           bool
	ModeratedRoom         bool
	PasswordProtectedRoom bool
	PersistentRoom        bool
	PresenceBroadcast     []string // This is a list of the roles for which presence is broadcast, 'moderator', 'participant' or 'visitor'
	PublicRoom            bool
	Description           string
	Name                  string
	Whois                 string // This can either be 'moderators' or 'anyone'
}

// NewRoom returns a newly created room
func NewRoom(roomID jid.Bare) *Room {
	return &Room{
		ID:        roomID,
		roster:    newRoomRoster(),
		observers: newRoomObservers(),
	}
}

// SelfOccupant returns the self occupant of the room
func (r *Room) SelfOccupant() *Occupant {
	return r.selfOccupant
}

// SelfOccupantNickname returns the nickname of the room's self occupant
func (r *Room) SelfOccupantNickname() string {
	o := r.SelfOccupant()
	if o != nil {
		return o.Nickname
	}

	return ""
}

// Roster returns the RoomRoster for this room
func (r *Room) Roster() *RoomRoster {
	return r.roster
}

// Subscribe subscribes the observer to room events
func (r *Room) Subscribe(f func(events.MUC)) {
	r.observers.subscribe(f)
}

// Publish will publish a new room event
func (r *Room) Publish(ev events.MUC) {
	r.observers.publishEvent(ev)
}

// AddSelfOccupant set the own occupant of the room
func (r *Room) AddSelfOccupant(occupant *Occupant) {
	r.selfOccupant = occupant
}

// SelfOccupantIsJoined returns true if the self occupant is in the room, false in otherwise
func (r *Room) SelfOccupantIsJoined() bool {
	return r.selfOccupant != nil
}

// GetSubject returns the room subject
func (r *Room) GetSubject() string {
	return r.subject
}

// UpdateSubject updates the room subject and returns a boolean
// indicating if the subject was updated (true) or not (false)
func (r *Room) UpdateSubject(s string) bool {
	r.subject = s

	if r.subjectWasUpdated {
		return true
	}

	r.subjectWasUpdated = true
	return false
}
