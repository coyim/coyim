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

	// TODO: this one feels like the coupling with the user account and the room
	// is too tight. We should think about other ways to represent this
	Joined bool

	// TODO: I don't like the name of this field. A Room has many occupants,
	// what makes this one special?
	Occupant *Occupant
	roster   *RoomRoster

	// TODO: I think this name is a bit ambigious. Maybe
	// observers or something like that instead?
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

// TODO: change the subscribers to use callback functions
// instead of channels, for more flexibility

// Subscribe subscribes the observer to room events
func (r *Room) Subscribe(c chan<- MUC) {
	r.subscribers.subscribe(c)
}

// TODO: Unsubscribe is not really necessary. The GUI will need
// to continue listening for events for as long as you haven't left
// the room

// Unsubscribe unsubscribe the observer to room events
func (r *Room) Unsubscribe(c chan<- MUC) {
	r.subscribers.unsubscribe(c)
}

// Publish will publish a new room event
func (r *Room) Publish(ev MUC) {
	r.subscribers.publishEvent(ev)
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
