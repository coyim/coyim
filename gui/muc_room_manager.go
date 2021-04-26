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
// Please note that "backToPreviousStep" will be called from the UI thread too
func (u *gtkUI) joinRoom(a *account, roomID jid.Bare, backToPreviousStepProvider func()) {
	rvd := newRoomViewData()
	rvd.backToPreviousStep = backToPreviousStepProvider
	u.joinRoomWithData(a, roomID, rvd)
}

func (u *gtkUI) joinRoomWithData(a *account, roomID jid.Bare, d roomViewDataProvider) {
	if d == nil {
		d = newRoomViewData()
	}

	v := u.getOrCreateRoomView(a, roomID)

	if v.isOpen() {
		v.present()
		return
	}

	if v.isSelfOccupantInTheRoom() {
		v.switchToMainView()
	} else {
		v.backToPreviousStep = d.backToPreviousStepProvider()
		v.passwordProvider = d.passwordProvider
		v.switchToLobbyView()
	}

	v.show()
}
