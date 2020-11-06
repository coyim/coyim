package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type roomViewToolbar struct {
	menu gtki.Popover

	view             gtki.Box        `gtk-widget:"room-view-toolbar"`
	roomNameLabel    gtki.Label      `gtk-widget:"room-name-label"`
	roomSubjectLabel gtki.Label      `gtk-widget:"room-subject-label"`
	roomStatusIcon   gtki.Image      `gtk-widget:"room-status-icon"`
	roomMenu         gtki.MenuButton `gtk-widget:"room-menu"`
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
}

func (t *roomViewToolbar) initDefaults(v *roomView) {
	t.roomStatusIcon.SetFromPixbuf(getMUCIconPixbuf("room"))

	t.roomNameLabel.SetText(v.roomID().String())
	updateWithStyle(t.roomNameLabel, providerWithStyle("label", style{
		"font-size":   "22px",
		"font-weight": "bold",
	}))

	t.displayRoomSubject(v.room.GetSubject())
	updateWithStyle(t.roomSubjectLabel, providerWithStyle("label", style{
		"font-size":  "14px",
		"font-style": "italic",
		"color":      "#666666",
	}))

	t.roomMenu.SetPopover(v.getRoomMenuWidget())
}

func (t *roomViewToolbar) initSubscribers(v *roomView) {
	v.subscribe("toolbar", func(ev roomViewEvent) {
		switch e := ev.(type) {
		case subjectReceivedEvent:
			t.subjectReceivedEvent(e.subject)
		case subjectUpdatedEvent:
			t.subjectReceivedEvent(e.subject)
		}
	})
}

func (t *roomViewToolbar) subjectReceivedEvent(subject string) {
	doInUIThread(func() {
		t.displayRoomSubject(subject)
	})
}

// displayRoomSubject MUST be called from the UI thread
func (t *roomViewToolbar) displayRoomSubject(subject string) {
	t.roomSubjectLabel.SetText(subject)
	if subject == "" {
		t.roomSubjectLabel.Hide()
	} else {
		t.roomSubjectLabel.Show()
	}
}
