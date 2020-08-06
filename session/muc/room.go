package muc

import "github.com/coyim/coyim/xmpp/jid"

// Room represents a multi user chat room that a session is currently connected to.
// It contains information about the room configuration itself, and the participants in the room
type Room struct {
	Identity jid.Bare

	Subject string

	Opaque interface{}

	roster *RoomRoster

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
		Identity: ident,
		roster:   newRoomRoster(),
	}
}

// Roster returns the RoomRoster for this room
func (r *Room) Roster() *RoomRoster {
	return r.roster
}
