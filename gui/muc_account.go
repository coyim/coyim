package gui

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
)

func (a *account) getRoomView(roomID jid.Bare) (*roomView, bool) {
	// TODO: I think mucRoomsLock should be fine for this field
	a.multiUserChatRoomsLock.RLock()
	defer a.multiUserChatRoomsLock.RUnlock()

	// TODO: This one too, mucRooms
	v, ok := a.multiUserChatRooms[roomID.String()]
	if !ok {
		a.log.WithField("room", roomID).Debug("getRoomView(): trying to get a not connected room")
	}

	return v, ok
}

func (a *account) addRoomView(v *roomView) {
	a.multiUserChatRoomsLock.Lock()
	defer a.multiUserChatRoomsLock.Unlock()

	a.multiUserChatRooms[v.roomID().String()] = v
}

func (a *account) removeRoomView(roomID jid.Bare) {
	a.multiUserChatRoomsLock.Lock()
	defer a.multiUserChatRoomsLock.Unlock()

	_, exists := a.multiUserChatRooms[roomID.String()]
	if !exists {
		return
	}

	delete(a.multiUserChatRooms, roomID.String())
}

func (a *account) newRoomModel(roomID jid.Bare) *muc.Room {
	return a.session.NewRoom(roomID)
}

func (a *account) leaveRoom(roomID jid.Bare, nickname string, onSuccess func(), onError func(error)) {
	ok, anyError := a.session.LeaveRoom(roomID, nickname)

	go func() {
		select {
		case <-ok:
			a.removeRoomView(roomID)
			if onSuccess != nil {
				onSuccess()
			}
		case err := <-anyError:
			// TODO: Would it be possible for us to log this on the room logger?
			a.log.WithError(err).Error("An error occurred while trying to leave the room.")
			if onError != nil {
				onError(err)
			}
		}
	}()
}
