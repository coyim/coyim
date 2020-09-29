package gui

import (
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type mucCreateRoomViewSuccess struct {
	ac         *account
	roomID     jid.Bare
	onJoinRoom func(*account, jid.Bare)

	view gtki.Box `gtk-widget:"createRoomSuccess"`
}

func (v *mucCreateRoomView) newCreateRoomSuccess() *mucCreateRoomViewSuccess {
	s := &mucCreateRoomViewSuccess{
		onJoinRoom: v.joinRoom,
	}

	s.initBuilder(v)

	return s
}

func (s *mucCreateRoomViewSuccess) initBuilder(v *mucCreateRoomView) {
	builder := newBuilder("MUCCreateRoomSuccess")
	panicOnDevError(builder.bindObjects(s))

	builder.ConnectSignals(map[string]interface{}{
		"on_createRoom_clicked": v.showCreateForm,
		"on_joinRoom_clicked": func() {
			s.onJoinRoom(s.ac, s.roomID)
		},
	})
}

func (v *mucCreateRoomView) initCreateRoomSuccess() *mucCreateRoomViewSuccess {
	s := v.newCreateRoomSuccess()
	return s
}

func (s *mucCreateRoomViewSuccess) showSuccessView(v *mucCreateRoomView, a *account, roomID jid.Bare) {
	v.form.reset()
	v.container.Remove(v.form.view)
	v.success.updateInfo(a, roomID)
	v.container.Add(s.view)
}

func (s *mucCreateRoomViewSuccess) updateInfo(a *account, roomID jid.Bare) {
	s.ac = a
	s.roomID = roomID
}

func (s *mucCreateRoomViewSuccess) reset() {
	s.ac = nil
	s.roomID = nil
}
