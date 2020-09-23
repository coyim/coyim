package gui

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
)

func (a *account) getRoomView(ident jid.Bare) (*roomView, bool) {
	a.multiUserChatRoomsLock.RLock()
	defer a.multiUserChatRoomsLock.RUnlock()

	v, ok := a.multiUserChatRooms[ident.String()]
	if !ok {
		a.log.WithField("room", ident).Debug("getRoomView(): trying to get a not connected room")
	}

	return v, ok
}

func (a *account) addRoomView(v *roomView) {
	a.multiUserChatRoomsLock.Lock()
	defer a.multiUserChatRoomsLock.Unlock()

	a.multiUserChatRooms[v.identity().String()] = v
}

func (a *account) removeRoomView(ident jid.Bare) {
	a.multiUserChatRoomsLock.Lock()
	defer a.multiUserChatRoomsLock.Unlock()

	_, exists := a.multiUserChatRooms[ident.String()]
	if !exists {
		return
	}

	delete(a.multiUserChatRooms, ident.String())
}

func (a *account) newRoomModel(ident jid.Bare) *muc.Room {
	return a.session.NewRoom(ident)
}

func (a *account) leaveRoom(ident jid.Bare, nickname string, onSuccess func(), onError func(error)) {
	ok, anyError := a.session.LeaveRoom(ident, nickname)

	go func() {
		select {
		case <-ok:
			a.removeRoomView(ident)
			if onSuccess != nil {
				onSuccess()
			}
		case err := <-anyError:
			a.log.WithError(err).Error("An error occurred while trying to leave the room.")
			if onError != nil {
				onError(err)
			}
		}
	}()
}
