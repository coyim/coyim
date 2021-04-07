package gui

import (
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomViewToolbar struct {
	view                  gtki.Box        `gtk-widget:"room-view-toolbar"`
	roomNameLabel         gtki.Label      `gtk-widget:"room-name-label"`
	roomStatusIcon        gtki.Image      `gtk-widget:"room-status-icon"`
	roomMenuButton        gtki.MenuButton `gtk-widget:"room-menu-button"`
	roomSubjectBox        gtki.Box        `gtk-widget:"room-subject-box"`
	roomSubjectLabel      gtki.Label      `gtk-widget:"room-subject-label"`
	roomSubjectHideButton gtki.Button     `gtk-widget:"room-subject-hide-button"`
	roomSubjectShowButton gtki.Button     `gtk-widget:"room-subject-show-button"`
	roomMenu              gtki.Menu       `gtk-widget:"room-menu"`
	leaveRoomMenuItem     gtki.MenuItem   `gtk-widget:"leave-room-menu-item"`
	destroyRoomMenuItem   gtki.MenuItem   `gtk-widget:"destroy-room-menu-item"`
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
		"on_leave_room":        v.onLeaveRoom,
		"on_destroy_room":      v.onDestroyRoom,
		"on_show_room_subject": t.onShowRoomSubject,
		"on_hide_room_subject": t.onHideRoomSubject,
	})
}

func (t *roomViewToolbar) initDefaults(v *roomView) {
	t.roomStatusIcon.SetFromPixbuf(getMUCIconPixbuf("room"))

	t.roomNameLabel.SetText(v.roomID().String())
	mucStyles.setRoomToolbarNameLabelStyle(t.roomNameLabel)

	t.displayRoomSubject(v.room.GetSubject())
	mucStyles.setRoomToolbarSubjectLabelStyle(t.roomSubjectLabel)
}

func (t *roomViewToolbar) initSubscribers(v *roomView) {
	v.subscribe("toolbar", func(ev roomViewEvent) {
		switch e := ev.(type) {
		case subjectReceivedEvent:
			t.subjectReceivedEvent(e.subject)
		case subjectUpdatedEvent:
			t.subjectReceivedEvent(e.subject)
		case roomDestroyedEvent:
			t.roomDestroyedEvent()
		case selfOccupantRemovedEvent:
			t.selfOccupantRemovedEvent()
		case occupantSelfJoinedEvent:
			t.selfOccupantJoinedEvent(v.isSelfOccupantAnOwner())
		case selfOccupantRoleUpdatedEvent:
			t.selfOccupantRoleUpdatedEvent(e.selfRoleUpdate.New)
		}
	})
}

func (t *roomViewToolbar) subjectReceivedEvent(subject string) {
	doInUIThread(func() {
		t.displayRoomSubject(subject)
	})
}

func (t *roomViewToolbar) roomDestroyedEvent() {
	doInUIThread(t.disable)
}

func (t *roomViewToolbar) selfOccupantRemovedEvent() {
	doInUIThread(t.disable)
}

func (t *roomViewToolbar) selfOccupantRoleUpdatedEvent(role data.Role) {
	if role.IsNone() {
		doInUIThread(t.disable)
	}
}

func (t *roomViewToolbar) selfOccupantJoinedEvent(isOwner bool) {
	t.destroyRoomMenuItem.SetVisible(isOwner)
}

// disable MUST be called from UI Thread
func (t *roomViewToolbar) disable() {
	mucStyles.setRoomToolbarNameLabelDisabledStyle(t.roomNameLabel)
	t.roomStatusIcon.SetFromPixbuf(getMUCIconPixbuf("room-offline"))
	t.roomMenuButton.Hide()
}

// displayRoomSubject MUST be called from the UI thread
func (t *roomViewToolbar) displayRoomSubject(subject string) {
	t.roomSubjectShowButton.Hide()
	t.roomSubjectLabel.SetText(subject)

	if subject == "" {
		t.roomSubjectBox.Hide()
	} else {
		t.roomSubjectBox.Show()
	}
}

func (t *roomViewToolbar) onShowRoomSubject() {
	t.roomSubjectShowButton.Hide()
	t.roomSubjectBox.Show()
}

func (t *roomViewToolbar) onHideRoomSubject() {
	t.roomSubjectShowButton.Show()
	t.roomSubjectBox.Hide()
}
