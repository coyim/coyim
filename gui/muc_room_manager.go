package gui

import (
	"errors"

	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
)

func newRoomManager() *muc.RoomManager {
	return muc.NewRoomManager()
}

func (a *account) joinRoom(u *gtkUI, rjid jid.Bare) (*muc.Room, error) {
	return u.addRoom(a, rjid)
}

func (u *gtkUI) addRoom(a *account, ident jid.Bare) (*muc.Room, error) {
	_, exists := a.roomManager.GetRoom(ident)
	if exists {
		return nil, errors.New("the room is already in the manager")
	}

	r := u.newRoom(a, ident)
	if !a.roomManager.AddRoom(r) {
		return nil, errors.New("the room is already in the manager")
	}

	view := u.viewForRoom(r)

	go u.observeMUCRoomEvents(view)
	a.session.Subscribe(view.events)

	return r, nil
}
