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

// joinRoom MUST always be called from the UI thread
//
// Also, when we want to show a chat room, having a "return to" function that
// will be called from the lobby only when the user wants to "cancel" or "return"
// might be useful in some scenarios like "returning to previous step".
//
// Please note that "returnTo" will be called from the UI thread too
func (u *gtkUI) joinRoom(a *account, roomID jid.Bare, returnTo func()) {
	v := u.getOrCreateRoomView(a, roomID)

	if v.isOpen() {
		v.present()
		return
	}

	if v.isSelfOccupantJoined() {
		v.switchToMainView()
	} else {
		v.returnTo = returnTo
		v.switchToLobbyView()
	}

	v.show()
}
