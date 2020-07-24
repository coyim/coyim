package muc

import (
	"fmt"

	"github.com/coyim/gotk3adapter/gtki"
)

type mucRoomsFakeServer struct {
	rooms map[string]*mucRoom
}

type mucRoom struct {
	id      string
	name    string
	status  mucPeerStatus
	members *mucMembers
}

type mucMembers struct {
	widget gtki.ScrolledWindow `gtk-widget:"room-members"`
	model  gtki.ListStore      `gtk-widget:"room-members-model"`
	view   gtki.TreeView       `gtk-widget:"room-members-tree"`
}

func (r *mucRoomsFakeServer) addRoom(id string, room *mucRoom) {
	r.rooms[id] = room
}

func (r *mucRoomsFakeServer) byID(id string) (*mucRoom, error) {
	if room, ok := r.rooms[id]; ok {
		return room, nil
	}
	return nil, fmt.Errorf("roomt %s not found", id)
}

func (m *mucUI) initRooms() {
	s := &mucRoomsFakeServer{
		rooms: map[string]*mucRoom{},
	}

	rooms := fakeRooms()
	for id, r := range rooms {
		s.addRoom(id, r)
	}

	m.roomsServer = s
}
