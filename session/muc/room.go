package muc

import (
	"time"

	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/session/muc/data"

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

	observers         *roomObservers
	discussionHistory *data.DiscussionHistory

	// Configuration options:
	properties *RoomListing
}

// NewRoom returns a newly created room
func NewRoom(roomID jid.Bare) *Room {
	return &Room{
		ID:                roomID,
		roster:            newRoomRoster(),
		observers:         newRoomObservers(),
		discussionHistory: data.NewDiscussionHistory(),
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

// IsSelfOccupantInTheRoom returns true if the self occupant is in the room, false in otherwise
func (r *Room) IsSelfOccupantInTheRoom() bool {
	return r.selfOccupant != nil
}

// IsSelfOccupantAnOwner returns a boolean indicating if the self occupant is an owner
func (r *Room) IsSelfOccupantAnOwner() bool {
	return r.IsSelfOccupantInTheRoom() && r.selfOccupant.Affiliation.IsOwner()
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

// HasHistory returns true if room has history
func (r *Room) HasHistory() bool {
	return len(r.discussionHistory.GetHistory()) > 0
}

// GetHistory returns the room history
func (r *Room) GetHistory() *data.DiscussionHistory {
	return r.discussionHistory
}

// AddHistoryMessage adds a new message in the room history
func (r *Room) AddHistoryMessage(nickname, message string, timestamp time.Time) {
	r.discussionHistory.AddMessage(nickname, message, timestamp)
}

// UpdateProperties updates the room properties
func (r *Room) UpdateProperties(properties *RoomListing) {
	r.properties = properties
}
