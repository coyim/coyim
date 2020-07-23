package muc

import (
	"sync"

	"github.com/coyim/coyim/xmpp/jid"
)

// RoomListing contains the information about a room for listing it
type RoomListing struct {
	Service     jid.Any
	ServiceName string
	Jid         jid.Bare
	Name        string

	SupportsVoiceRequests     bool
	AllowsRegistration        bool
	Anonymity                 string
	Persistent                bool
	Moderated                 bool
	Open                      bool
	PasswordProtected         bool
	Public                    bool
	Language                  string
	OccupantsCanChangeSubject bool
	Description               string
	Occupants                 int
	MembersCanInvite          bool

	lockUpdates sync.RWMutex
	onUpdates   []func()
}

// NewRoomListing creates and returns a new room listing
func NewRoomListing() *RoomListing {
	return &RoomListing{}
}

// OnUpdate takes a function and some data, and when this room listing is updated, that function
// will be called with the current room listing and the associated data
func (rl *RoomListing) OnUpdate(f func(*RoomListing, interface{}), data interface{}) {
	rl.lockUpdates.Lock()
	defer rl.lockUpdates.Unlock()

	rl.onUpdates = append(rl.onUpdates, func() {
		f(rl, data)
	})
}

// Updated should be called after a room listing has been updated, to notify observers of the update
func (rl *RoomListing) Updated() {
	rl.lockUpdates.RLock()
	defer rl.lockUpdates.RUnlock()

	for _, f := range rl.onUpdates {
		f()
	}
}
