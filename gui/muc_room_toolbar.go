package gui

import "github.com/coyim/gotk3adapter/gtki"

type roomViewToolbar struct {
	view                 gtki.Box    `gtk-widget:"roomToolbar"`
	roomNameLabel        gtki.Label  `gtk-widget:"roomNameLabel"`
	roomDescriptionLabel gtki.Label  `gtk-widget:"roomDescriptionLabel"`
	togglePanelButton    gtki.Button `gtk-widget:"togglePanelButton"`
	leaveRoomButton      gtki.Button `gtk-widget:"leaveRoomButton"`
}

func (v *roomView) newRoomViewToolbar() *roomViewToolbar {
	t := &roomViewToolbar{}

	builder := newBuilder("MUCRoomToolbar")
	panicOnDevError(builder.bindObjects(t))

	builder.ConnectSignals(map[string]interface{}{
		"on_leave_room": func() {
			t.leaveRoomButton.SetSensitive(false)
			v.tryLeaveRoom(nil, func() {
				if v.isOpen() {
					doInUIThread(func() {
						t.leaveRoomButton.SetSensitive(true)
					})
				}
			})
		},
	})

	t.leaveRoomButton.SetSensitive(v.isJoined())
	v.subscribe("toolbar", occupantSelfJoined, func(roomViewEventInfo) {
		t.leaveRoomButton.SetSensitive(true)
	})

	return t
}
