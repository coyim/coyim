package gui

import "github.com/coyim/gotk3adapter/gtki"

type roomViewToolbar struct {
	view                 gtki.Box    `gtk-widget:"roomToolbar"`
	roomNameLabel        gtki.Label  `gtk-widget:"roomNameLabel"`
	roomDescriptionLabel gtki.Label  `gtk-widget:"roomDescriptionLabel"`
	leaveRoomButton      gtki.Button `gtk-widget:"leaveRoomButton"`
}

func (v *roomView) newRoomViewToolbar() *roomViewToolbar {
	t := &roomViewToolbar{}

	t.initBuilder(v)
	t.initDefaults(v)
	t.initSubscribers(v)

	return t
}

func (t *roomViewToolbar) initBuilder(v *roomView) {
	builder := newBuilder("MUCRoomToolbar")
	panicOnDevError(builder.bindObjects(t))

	builder.ConnectSignals(map[string]interface{}{
		"on_leave_room": func() {
			// TODO: Maybe extract?
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
}

func (t *roomViewToolbar) initDefaults(v *roomView) {
	t.leaveRoomButton.SetSensitive(v.isJoined())
}

func (t *roomViewToolbar) initSubscribers(v *roomView) {
	v.subscribe("toolbar", "occupantSelfJoinedEvent", func(roomViewEventInfo) {
		t.leaveRoomButton.SetSensitive(true)
	})
}
