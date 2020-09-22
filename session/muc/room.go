package muc

import (
	"sync"

	"github.com/coyim/coyim/xmpp/jid"
)

// Room represents a multi user chat room that a session is currently connected to.
// It contains information about the room configuration itself, and the participants in the room
type Room struct {
	Identity jid.Bare
	Subject  string
	Joined   bool

	Occupant    *Occupant
	roster      *RoomRoster
	subscribers *roomSubscribers

	once sync.Once

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
func NewRoom(ident jid.Bare) *Room {
	return &Room{
		Identity:    ident,
		roster:      newRoomRoster(),
		subscribers: newRoomSubscribers(),
	}
}

// Roster returns the RoomRoster for this room
func (r *Room) Roster() *RoomRoster {
	return r.roster
}

// Subscribe subscribes the observer to room events
func (r *Room) Subscribe(c chan<- MUC) {
	r.subscribers.subscribe(c)
}

// Publish will publish a new room event
func (r *Room) Publish(ev MUC) {
	r.subscribers.publishEvent(ev)
}

// AddSelfOccupant set the own occupant of the room
func (r *Room) AddSelfOccupant(occupant *Occupant) {
	r.once.Do(func() {
		r.Joined = true
		r.Occupant = occupant
	})
}
