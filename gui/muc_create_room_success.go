package gui

import (
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

func (v *mucCreateRoomView) initCreateRoomSuccess(d *mucCreateRoomData) {
	v.success = v.newCreateRoomSuccess()
	v.success.updateJoinRoomData(d)
}

func (v *mucCreateRoomView) showSuccessView(ca *account, roomID jid.Bare) {
	v.form.reset()
	v.container.Remove(v.form.view)
	v.success.updateInfo(ca, roomID)
	v.container.Add(v.success.view)
}

type mucCreateRoomViewSuccess struct {
	ca           *account
	roomID       jid.Bare
	joinRoomData roomViewDataProvider

	icon gtki.Image `gtk-widget:"createRoomSuccessImage"`
	view gtki.Box   `gtk-widget:"createRoomSuccess"`
}

func (v *mucCreateRoomView) newCreateRoomSuccess() *mucCreateRoomViewSuccess {
	s := &mucCreateRoomViewSuccess{}

	builder := newBuilder("MUCCreateRoomSuccess")
	panicOnDevError(builder.bindObjects(s))

	builder.ConnectSignals(map[string]interface{}{
		"on_createRoom_clicked": func() {
			v.notifications.clearErrors()
			v.showCreateForm()
		},
		"on_joinRoom_clicked": func() {
			v.destroy()
			go v.joinRoom(s.ca, s.roomID, s.joinRoomData)
		},
	})

	s.icon.SetFromPixbuf(getMUCIconPixbuf("dialog_ok"))

	return s
}

func (s *mucCreateRoomViewSuccess) updateInfo(ca *account, roomID jid.Bare) {
	s.ca = ca
	s.roomID = roomID
}

func (s *mucCreateRoomViewSuccess) updateJoinRoomData(d roomViewDataProvider) {
	s.joinRoomData = d
}

func (s *mucCreateRoomViewSuccess) reset() {
	s.ca = nil
	s.roomID = nil
	s.joinRoomData = nil
}
