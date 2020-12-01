package gui

import "github.com/coyim/coyim/xmpp/jid"

// leaveRoom should return the context so the caller can cancel the context early if required
func (a *account) leaveRoom(roomID jid.Bare, nickname string, onSuccess func(), onError func(error)) {
	leaveRoom := func() (<-chan bool, <-chan error) {
		return a.session.LeaveRoom(roomID, nickname)
	}

	leaveRoomSuccess := func() {
		a.removeRoomView(roomID)
		if onSuccess != nil {
			onSuccess()
		}
	}

	ctx := a.newAccountRoomOpContext("leave-room", roomID, leaveRoom, leaveRoomSuccess, onError)

	go ctx.doOperation()
}

// destroyRoom should return the context so the caller can cancel the context early if required
func (a *account) destroyRoom(roomID jid.Bare, reason string, alternativeRoomID jid.Bare, password string, onSuccess func(), onError func(error)) {
	destroyRoom := func() (<-chan bool, <-chan error) {
		return a.session.DestroyRoom(roomID, reason, alternativeRoomID, password)
	}

	ctx := a.newAccountRoomOpContext("destroy-room", roomID, destroyRoom, onSuccess, onError)

	go ctx.doOperation()
}
