package muc

import (
	"fmt"
)

type roomsFakeServer struct {
	rooms map[string]*room
}

func (r *roomsFakeServer) addRoom(id string, room *room) {
	if room.rosterItem == nil {
		room.rosterItem = &rosterItem{
			id: id,
		}
	}

	r.rooms[id] = room
}

func (r *roomsFakeServer) byID(id string) (*room, error) {
	if room, ok := r.rooms[id]; ok {
		return room, nil
	}
	return nil, fmt.Errorf("roomt %s not found", id)
}

func (u *gtkUI) initRooms() {
	s := &roomsFakeServer{
		rooms: map[string]*room{},
	}

	rooms := fakeRooms()
	for id, r := range rooms {
		s.addRoom(id, r)
	}

	u.roomsServer = s
}
