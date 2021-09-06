package muc

import (
	"sync"

	"github.com/coyim/coyim/xmpp/jid"
)

// RoomManager contains information about each room that is currently active for a user
// When a window is closed, the room stays in this list. A room will only be removed
// from this list when the current user leaves that room.
type RoomManager struct {
	rooms map[string]*Room

	lock sync.RWMutex
}

// NewRoomManager returns a newly created room manager
func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[string]*Room),
	}
}

// GetRoom returns the room corresponding to the given identifier, or returns false if it can't be found
func (rm *RoomManager) GetRoom(roomID jid.Bare) (*Room, bool) {
	rm.lock.RLock()
	defer rm.lock.RUnlock()

	r, ok := rm.rooms[roomID.String()]
	return r, ok
}

// GetAllRooms returns the occupant's active rooms
func (rm *RoomManager) GetAllRooms() []*Room {
	rm.lock.RLock()
	defer rm.lock.RUnlock()

	rooms := []*Room{}
	for _, v := range rm.rooms {
		rooms = append(rooms, v)
	}

	return rooms
}

// AddRoom adds the room to the manager. If the room is already in the manager, this method will return
// false
func (rm *RoomManager) AddRoom(r *Room) bool {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	_, ok := rm.rooms[r.ID.String()]
	if ok {
		return false
	}

	rm.rooms[r.ID.String()] = r
	return true
}

// DeleteRoom will remove the room with the given identifier from the manager.
// If there is no such room ID, DeleteRoom is a no-op.
func (rm *RoomManager) DeleteRoom(room jid.Bare) {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	delete(rm.rooms, room.String())
}
