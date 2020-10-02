package muc

import (
	"sync"

	"github.com/coyim/coyim/session/events"

	"github.com/coyim/coyim/xmpp/jid"
)

// Room represents a multi user chat room that a session is currently connected to.
// It contains information about the room configuration itself, and the participants in the room
type Room struct {
	ID      jid.Bare
	Subject string

	// TODO: this one feels like the coupling with the user account and the room
	// is too tight. We should think about other ways to represent this
	Joined bool

	// TODO: I don't like the name of this field. A Room has many occupants,
	// what makes this one special?
	Occupant *Occupant
	roster   *RoomRoster

	observers *roomObservers

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
func NewRoom(roomID jid.Bare) *Room {
	return &Room{
		ID:        roomID,
		roster:    newRoomRoster(),
		observers: newRoomObservers(),
	}
}

// Roster returns the RoomRoster for this room
func (r *Room) Roster() *RoomRoster {
	return r.roster
}

// TODO: change the subscribers to use callback functions
// instead of channels, for more flexibility

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
	// TODO: this logic probably belongs in the muc manager
	// where we can do different things directly depending on
	// whether the self occupant has joined or not

	r.once.Do(func() {
		r.Joined = true
		r.Occupant = occupant
	})
}
