package muc

import (
	"sync"

	"github.com/coyim/coyim/xmpp/jid"
)

// RoomManager contains information about each room that is currently active for a user
// When a window is closed, the room stays in this list. A room will only be removed
// from this list when the current user leaves that room.
type RoomManager struct {
	lock sync.RWMutex

	rooms map[string]*Room
}

// NewRoomManager returns a newly created room manager
func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[string]*Room),
	}
}

// GetRoom returns the room corresponding to the given identifier, or returns false if it can't be found
func (rm *RoomManager) GetRoom(ident jid.Bare) (*Room, bool) {
	rm.lock.RLock()
	defer rm.lock.RUnlock()

	r, ok := rm.rooms[ident.String()]
	return r, ok
}

// AddRoom adds the room to the manager. If the room is already in the manager, this method will return
// false
func (rm *RoomManager) AddRoom(r *Room) bool {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	_, ok := rm.rooms[r.Identity.String()]
	if ok {
		return false
	}

	rm.rooms[r.Identity.String()] = r
	return true
}

// LeaveRoom will remove the room with the given identifier from the manager. If the room doesn't exist, this method
// will return false
func (rm *RoomManager) LeaveRoom(room jid.Bare) bool {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	_, ok := rm.rooms[room.String()]
	if !ok {
		return false
	}

	delete(rm.rooms, room.String())
	return true
}
