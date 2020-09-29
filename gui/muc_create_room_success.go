package gui

import (
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type mucCreateRoomViewSuccess struct {
	ca         *account
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
			s.onJoinRoom(s.ca, s.roomID)
		},
	})
}

func (v *mucCreateRoomView) initCreateRoomSuccess() *mucCreateRoomViewSuccess {
	s := v.newCreateRoomSuccess()
	return s
}

func (s *mucCreateRoomViewSuccess) showSuccessView(v *mucCreateRoomView, ca *account, roomID jid.Bare) {
	v.form.reset()
	v.container.Remove(v.form.view)
	v.success.updateInfo(ca, roomID)
	v.container.Add(s.view)
}

func (s *mucCreateRoomViewSuccess) updateInfo(ca *account, roomID jid.Bare) {
	s.ca = ca
	s.roomID = roomID
}

func (s *mucCreateRoomViewSuccess) reset() {
	s.ca = nil
	s.roomID = nil
}
