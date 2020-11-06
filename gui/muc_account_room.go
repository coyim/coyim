package gui

import "github.com/coyim/coyim/xmpp/jid"

// leaveRoom should return the context so the caller can cancel the context early if required
func (a *account) leaveRoom(roomID jid.Bare, nickname string, onSuccess func(), onError func(error)) {
	leaveRoomCb := func() (<-chan bool, <-chan error, func()) {
		ok, err := a.session.LeaveRoom(roomID, nickname)
		return ok, err, nil
	}

	leaveRoomSuccess := func() {
		a.removeRoomView(roomID)
		if onSuccess != nil {
			onSuccess()
		}
	}

	controller := a.newRoomOpController("leave-room", leaveRoomCb, leaveRoomSuccess, onError)
	ctx := a.newAccountRoomOpContext("leave-room", roomID, controller)

	go ctx.doOperation()
}

// destroyRoom should return the context so the caller can cancel the context early if required
func (a *account) destroyRoom(roomID jid.Bare, alternateID jid.Bare, reason string, onSuccess func(), onError func(error)) {
	destroyRoomCb := func() (<-chan bool, <-chan error, func()) {
		return a.session.DestroyRoom(roomID, alternateID, reason)
	}

	controller := a.newRoomOpController("destroy-room", destroyRoomCb, onSuccess, onError)
	ctx := a.newAccountRoomOpContext("destroy-room", roomID, controller)

	go ctx.doOperation()
}
