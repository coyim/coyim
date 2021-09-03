package muc

import (
	"sync"
	"time"

	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/session/muc/data"

	"github.com/coyim/coyim/xmpp/jid"
)

// Room represents a multi user chat room that a session is currently connected to.
// It contains information about the room configuration itself, and the participants in the room
type Room struct {
	ID jid.Bare

	subject       string
	subjectIsNew  bool
	subjectLocker sync.Mutex

	selfOccupant *Occupant
	roster       *RoomRoster

	observers         *roomObservers
	discussionHistory *data.DiscussionHistory

	properties data.RoomDiscoInfo
}

// NewRoom returns a newly created room
func NewRoom(roomID jid.Bare) *Room {
	return &Room{
		ID:                roomID,
		subjectIsNew:      true,
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

// HasSubject returns true if the room has subject
func (r *Room) HasSubject() bool {
	return r.subject != ""
}

// UpdateSubject updates the room subject and returns a boolean
// indicating if the subject was updated (true) or not (false)
func (r *Room) UpdateSubject(s string) bool {
	r.subjectLocker.Lock()
	defer r.subjectLocker.Unlock()

	r.subject = s
	isUpdated := !r.subjectIsNew
	r.subjectIsNew = false

	return isUpdated
}

// GetDiscussionHistory returns the room history
func (r *Room) GetDiscussionHistory() *data.DiscussionHistory {
	return r.discussionHistory
}

// AddHistoryMessage adds a new chat message in the room history
func (r *Room) AddHistoryMessage(nickname, message string, timestamp time.Time) {
	r.discussionHistory.AddMessage(nickname, message, timestamp, data.Chat)
}

// AddMessage adds a new message in the room history with a specific message type
func (r *Room) AddMessage(messageData *data.DelayedMessage) {
	r.discussionHistory.AddMessage(messageData.Nickname, messageData.Message, messageData.Timestamp, messageData.MessageType)
}

// SetProperties replaces the current room properties by new ones
func (r *Room) SetProperties(p data.RoomDiscoInfo) {
	r.properties = p
}

// SubjectCanBeChanged returns true if the subject of the room can be changed,
// specifically by the self occupant of the room.
func (r *Room) SubjectCanBeChanged() bool {
	occupantsCanChangeSubject := r.properties.OccupantsCanChangeSubject

	roomSelfOccupant := r.SelfOccupant()
	if roomSelfOccupant != nil {
		role := roomSelfOccupant.Role
		return !role.IsVisitor() && (role.IsModerator() || occupantsCanChangeSubject)
	}

	return false
}

// OnStatusConnected sets the room's initial values when a new connection is established
func (r *Room) OnStatusConnected() {
	r.roster.reset()

	r.subjectIsNew = true

	r.Publish(events.MUCSelfOccupantConnected{})
}

// OnStatusDisconnected sets appropriate values when an occupant loses connection to the room
func (r *Room) OnStatusDisconnected() {
	roomSelfOccupant := r.SelfOccupant()
	if roomSelfOccupant != nil {
		roomSelfOccupant.ChangeAffiliationToNone()
		roomSelfOccupant.ChangeRoleToNone()
	}

	r.Publish(events.MUCSelfOccupantDisconnected{})
}
