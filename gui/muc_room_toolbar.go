package gui

import "github.com/coyim/gotk3adapter/gtki"

type roomViewToolbar struct {
	view                 gtki.Box    `gtk-widget:"roomToolbar"`
	roomNameLabel        gtki.Label  `gtk-widget:"roomNameLabel"`
	roomDescriptionLabel gtki.Label  `gtk-widget:"roomDescriptionLabel"`
	togglePanelButton    gtki.Button `gtk-widget:"togglePanelButton"`
	leaveRoomButton      gtki.Button `gtk-widget:"leaveRoomButton"`
}

func (r *roomView) newRoomViewToolbar() *roomViewToolbar {
	t := &roomViewToolbar{}

	builder := newBuilder("MUCRoomToolbar")
	panicOnDevError(builder.bindObjects(t))

	builder.ConnectSignals(map[string]interface{}{
		"on_leave_room": r.LeaveRoom,
	})

	r.onSelfJoinReceived(func() {
		t.leaveRoomButton.SetSensitive(true)
	})

	return t
}
