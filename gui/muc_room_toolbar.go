package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type roomViewToolbar struct {
	view             gtki.Box    `gtk-widget:"roomToolbar"`
	roomNameLabel    gtki.Label  `gtk-widget:"roomNameLabel"`
	roomSubjectLabel gtki.Label  `gtk-widget:"roomSubjectLabel"`
	leaveRoomButton  gtki.Button `gtk-widget:"leaveRoomButton"`
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
			t.onLeaveRoom(v)
		},
	})
}

func (t *roomViewToolbar) initDefaults(v *roomView) {
	t.leaveRoomButton.SetSensitive(v.isJoined())

	t.initLabelFor(t.roomNameLabel, providerWithStyle("label", style{
		"font-size":   "22px",
		"font-weight": "bold",
	}))

	t.initLabelFor(t.roomSubjectLabel, providerWithStyle("label", style{
		"font-size":  "14px",
		"font-style": "italic",
		"color":      "#666666",
	}))

	t.roomNameLabel.SetText(v.roomID().String())

	doInUIThread(func() {
		t.roomSubjectLabel.Hide()
	})

	if v.room.Subject != "" {
		t.showSubject(v.room.Subject)
	}
}

func (t *roomViewToolbar) showSubject(subject string) {
	doInUIThread(func() {
		t.roomSubjectLabel.SetText(subject)
		t.roomSubjectLabel.Show()
	})
}

func (t *roomViewToolbar) initLabelFor(label gtki.Label, cssProvider gtki.CssProvider) {
	updateWithStyle(label, cssProvider)
}

func (t *roomViewToolbar) initSubscribers(v *roomView) {
	v.subscribe("toolbar", func(ev roomViewEvent) {
		switch ev.(type) {
		case occupantSelfJoinedEvent:
			doInUIThread(func() {
				t.leaveRoomButton.SetSensitive(true)
			})
		case subjectEvent:
			t.showSubject(v.room.Subject)
		}
	})
}

func (t *roomViewToolbar) onLeaveRoom(v *roomView) {
	t.leaveRoomButton.SetSensitive(false)
	v.tryLeaveRoom(nil, func() {
		if v.isOpen() {
			doInUIThread(func() {
				t.leaveRoomButton.SetSensitive(true)
			})
		}
	})
}
