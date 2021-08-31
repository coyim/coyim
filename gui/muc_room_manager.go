package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/jid"
)

// getOrCreateRoomView MUST be called from the UI thread
func (u *gtkUI) getOrCreateRoomView(a *account, roomID jid.Bare) *roomView {
	if v, ok := a.getRoomView(roomID); ok {
		return v
	}

	room := a.session.NewRoom(roomID)
	v := u.newRoomView(a, room)
	a.addRoomView(v)

	return v
}

// joinRoom MUST always be called from the UI thread
//
// Also, when we want to show a chat room, having a "return to" function that
// will be called from the lobby only when the user wants to "cancel" or "return"
// might be useful in some scenarios like "returning to previous step".
//
// Please note that "backToPreviousStep" will be called from the UI thread too
func (u *gtkUI) joinRoom(a *account, roomID jid.Bare, rvd roomViewDataProvider) {
	u.joinRoomWithData(a, roomID, rvd)
}

// joinRoomWithData MUST be called from the UI thread
func (u *gtkUI) joinRoomWithData(a *account, roomID jid.Bare, d roomViewDataProvider) {
	v := u.getOrCreateRoomView(a, roomID)

	if v.isOpen() {
		d.notifyError(i18n.Local("You are already in the room."))
		return
	}

	if d == nil {
		d = &roomViewData{}
	}

	if v.isSelfOccupantInTheRoom() {
		v.switchToMainView()
	} else {
		v.backToPreviousStep = d.backToPreviousStep()
		v.passwordProvider = d.passwordProvider
		v.switchToLobbyView()
	}

	d.doWhenNoErrors()
	v.show()
}
