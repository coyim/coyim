package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
)

// getOrCreateRoomView MUST be called from the UI thread
func (u *gtkUI) getOrCreateRoomView(a *account, room *muc.Room) *roomView {
	if v, ok := a.getRoomView(room.ID); ok {
		return v
	}

	v := u.newRoomView(a, room)
	a.addRoomView(v)

	return v
}

// joinRoom MUST NOT  be called from the UI thread
//
// Also, when we want to show a chat room, having a "return to" function that
// will be called from the lobby only when the user wants to "cancel" or "return"
// might be useful in some scenarios like "returning to previous step".
//
// Please note that "backToPreviousStep" will be called from the UI thread too
func (u *gtkUI) joinRoom(a *account, roomID jid.Bare, rvd roomViewDataProvider) {
	room := a.session.NewRoom(roomID)
	doInUIThread(func() {
		u.joinRoomWithData(a, room, rvd)
	})
}

// joinRoomWithData MUST be called from the UI thread
func (u *gtkUI) joinRoomWithData(a *account, room *muc.Room, d roomViewDataProvider) {
	v := u.getOrCreateRoomView(a, room)

	if v.isOpen() {
		d.notifyError(i18n.Local("You are already in the room."))
		v.present()
		return
	}

	if d == nil {
		d = &roomViewData{}
	}

	isOpeningRoomAgain := false
	if v.isSelfOccupantInTheRoom() {
		v.notifications.other(roomNotificationOptions{
			message:   i18n.Local("You were already connected to this room."),
			closeable: true,
		})

		isOpeningRoomAgain = true
		v.switchToMainView()
	} else {
		v.backToPreviousStep = d.backToPreviousStep()
		v.passwordProvider = d.passwordProvider
		v.switchToLobbyView()
	}

	d.doWhenNoErrorOccurred()
	v.show()

	if isOpeningRoomAgain {
		v.publishEvent(reopenRoomEvent{
			history: v.room.GetDiscussionHistory(),
			subject: v.room.GetSubject(),
		})
	}
}
