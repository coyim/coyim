package gui

import (
	"github.com/coyim/coyim/session/muc"
)

func newRoomManager() *muc.RoomManager {
	return muc.NewRoomManager()
}
