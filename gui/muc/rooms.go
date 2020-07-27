package muc

import (
	"fmt"

	"github.com/coyim/gotk3adapter/gtki"
)

type roomsFakeServer struct {
	rooms map[string]*room
}

type room struct {
	id      string
	name    string
	status  peerStatus
	members *members
}

type members struct {
	widget gtki.ScrolledWindow `gtk-widget:"room-members"`
	model  gtki.ListStore      `gtk-widget:"room-members-model"`
	view   gtki.TreeView       `gtk-widget:"room-members-tree"`
}

func (r *roomsFakeServer) addRoom(id string, room *room) {
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
