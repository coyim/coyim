package gui

import (
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type createMUCRoomSuccess struct {
	ac         *account
	ident      jid.Bare
	onJoinRoom func(*account, jid.Bare)

	view gtki.Box `gtk-widget:"createRoomSuccess"`
}

func (v *createMUCRoom) initSuccessView(onJoinRoom func(*account, jid.Bare)) {
	s := &createMUCRoomSuccess{
		onJoinRoom: onJoinRoom,
	}

	panicOnDevError(v.builder.bindObjects(s))

	v.addBuilderSignals(map[string]interface{}{
		"on_createRoom_clicked": v.showCreateForm,
		"on_joinRoom_clicked":   s.onJoinRoomClick,
	})

	v.showSuccessView = func(a *account, ident jid.Bare) {
		v.form.reset()
		v.container.Remove(v.form.view)
		v.success.updateInfo(a, ident)
		v.container.Add(v.success.view)
	}

	v.success = s
}

func (s *createMUCRoomSuccess) updateInfo(a *account, ident jid.Bare) {
	s.ac = a
	s.ident = ident
}

func (s *createMUCRoomSuccess) onJoinRoomClick() {
	s.onJoinRoom(s.ac, s.ident)
}

func (s *createMUCRoomSuccess) reset() {
	s.ac = nil
	s.ident = nil
}
