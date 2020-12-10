package gui

import "github.com/coyim/coyim/xmpp/jid"

// leaveRoom should return the context so the caller can cancel the context early if required
func (a *account) leaveRoom(roomID jid.Bare, nickname string, onSuccess func(), onError func(error), onDone func()) {
	leaveRoom := func() (<-chan bool, <-chan error) {
		return a.session.LeaveRoom(roomID, nickname)
	}

	leaveRoomSuccess := func() {
		a.removeRoomView(roomID)
		if onSuccess != nil {
			onSuccess()
		}
	}

	ctx := a.newAccountRoomOpContext(
		"leave-room",
		roomID,
		leaveRoom,
		leaveRoomSuccess,
		onError,
		onDone,
	)

	go ctx.doOperation()
}
