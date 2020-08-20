package gui

import (
	"errors"

	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
)

func newRoomManager() *muc.RoomManager {
	return muc.NewRoomManager()
}

func (u *gtkUI) addRoom(a *account, ident jid.Bare) (*muc.Room, error) {
	// TODO[OB]-MUC: Something about the logic in this method seems weird

	_, exists := a.roomManager.GetRoom(ident)
	if exists {
		return nil, errors.New("the room is already in the manager")
	}

	r := u.newRoom(ident)
	if !a.roomManager.AddRoom(r) {
		return nil, errors.New("the room is already in the manager")
	}

	return r, nil
}
