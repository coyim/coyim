package gui

import (
	"github.com/coyim/coyim/xmpp/jid"
)

func (u *gtkUI) getOrCreateRoomView(a *account, roomID jid.Bare) *roomView {
	v, exists := a.getRoomView(roomID)
	if !exists {
		v = newRoomView(u, a, roomID)
		a.addRoomView(v)
	}
	return v
}

var defaultJoinRoomData roomViewDataProvider

func getDefaultJoinRoomData() roomViewDataProvider {
	if defaultJoinRoomData == nil {
		defaultJoinRoomData = newRoomViewData()
	}
	return defaultJoinRoomData
}

// joinRoom MUST always be called from the UI thread
//
// Also, when we want to show a chat room, having a "return to" function that
// will be called from the lobby only when the user wants to "cancel" or "return"
// might be useful in some scenarios like "returning to previous step".
//
// Please note that "returnTo" will be called from the UI thread too
func (u *gtkUI) joinRoom(a *account, roomID jid.Bare, returnTo func()) {
	rvd := newRoomViewData()
	rvd.onReturn = returnTo
	u.joinRoomWithData(a, roomID, rvd)
}

func (u *gtkUI) joinRoomWithData(a *account, roomID jid.Bare, d roomViewDataProvider) {
	if d == nil {
		d = getDefaultJoinRoomData()
	}

	v := u.getOrCreateRoomView(a, roomID)

	if v.isOpen() {
		v.present()
		return
	}

	if v.isSelfOccupantInTheRoom() {
		v.switchToMainView()
	} else {
		v.returnTo = d.returnTo()
		v.passwordProvider = d.passwordProvider
		v.switchToLobbyView()
	}

	v.show()
}
