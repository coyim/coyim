package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/gotk3adapter/gtki"
)

const (
	roomSubjectShownIconName  = "go-up-symbolic"
	roomSubjectHiddenIconName = "go-down-symbolic"
)

type roomViewToolbar struct {
	view                   gtki.Box               `gtk-widget:"room-view-toolbar"`
	roomNameLabel          gtki.Label             `gtk-widget:"room-name-label"`
	roomStatusIcon         gtki.Image             `gtk-widget:"room-status-icon"`
	roomMenuButton         gtki.MenuButton        `gtk-widget:"room-menu-button"`
	roomSubjectButton      gtki.Button            `gtk-widget:"room-subject-button"`
	roomSubjectButtonImage gtki.Image             `gtk-widget:"room-subject-button-image"`
	roomSubjectRevealer    gtki.Revealer          `gtk-widget:"room-subject-revealer"`
	roomSubjectLabel       gtki.Label             `gtk-widget:"room-subject-label"`
	roomMenu               gtki.Menu              `gtk-widget:"room-menu"`
	modifyBanMenuItem      gtki.MenuItem          `gtk-widget:"modify-ban-list-menu-item"`
	adminActionsSeparator  gtki.SeparatorMenuItem `gtk-widget:"admin-action-separator"`
	leaveRoomMenuItem      gtki.MenuItem          `gtk-widget:"leave-room-menu-item"`
	destroyRoomMenuItem    gtki.MenuItem          `gtk-widget:"destroy-room-menu-item"`
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
		"on_leave_room":          v.onLeaveRoom,
		"on_destroy_room":        v.onDestroyRoom,
		"on_modify_ban_list":     v.onModifyBanList,
		"on_toggle_room_subject": t.onToggleRoomSubject,
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
			t.selfOccupantJoinedEvent(v.room.SelfOccupant().Affiliation)
		case selfOccupantRoleUpdatedEvent:
			t.selfOccupantRoleUpdatedEvent(e.selfRoleUpdate.New)
		case selfOccupantAffiliationUpdatedEvent:
			t.selfOccupantAffiliationUpdatedEvent(e.selfAffiliationUpdate.New)
		case selfOccupantAffiliationRoleUpdatedEvent:
			t.selfOccupantRoleUpdatedEvent(e.selfAffiliationRoleUpdate.NewRole)
			t.selfOccupantAffiliationUpdatedEvent(e.selfAffiliationRoleUpdate.NewAffiliation)
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

func (t *roomViewToolbar) selfOccupantAffiliationUpdatedEvent(affiliation data.Affiliation) {
	doInUIThread(func() {
		t.updateMenuActionsBasedOn(affiliation)
	})

	if affiliation.IsBanned() {
		doInUIThread(t.disable)
	}
}

func (t *roomViewToolbar) selfOccupantJoinedEvent(affiliation data.Affiliation) {
	doInUIThread(func() {
		t.updateMenuActionsBasedOn(affiliation)
	})
}

func (t *roomViewToolbar) updateMenuActionsBasedOn(affiliation data.Affiliation) {
	t.destroyRoomMenuItem.SetVisible(affiliation.IsOwner())

	showAdminActions := affiliation.IsOwner() || affiliation.IsAdmin()
	t.modifyBanMenuItem.SetVisible(showAdminActions)
	t.adminActionsSeparator.SetVisible(showAdminActions)
}

// disable MUST be called from UI Thread
func (t *roomViewToolbar) disable() {
	mucStyles.setRoomToolbarNameLabelDisabledStyle(t.roomNameLabel)
	t.roomStatusIcon.SetFromPixbuf(getMUCIconPixbuf("room-offline"))
	t.roomMenuButton.Hide()
}

// displayRoomSubject MUST be called from the UI thread
func (t *roomViewToolbar) displayRoomSubject(subject string) {
	t.roomSubjectLabel.SetText(subject)
	t.roomSubjectRevealer.SetRevealChild(false)

	t.roomSubjectButton.Hide()
	if subject != "" {
		t.roomSubjectButton.Show()
	}
}

// onToggleRoomSubject MUST be called from the UI thread
func (t *roomViewToolbar) onToggleRoomSubject() {
	if t.roomSubjectRevealer.GetRevealChild() {
		t.onHideRoomSubject()
	} else {
		t.onShowRoomSubject()
	}
}

// onShowRoomSubject MUST be called from the UI thread
func (t *roomViewToolbar) onShowRoomSubject() {
	t.roomSubjectRevealer.SetRevealChild(true)
	t.roomSubjectButton.SetTooltipText(i18n.Local("Hide room subject"))
	t.roomSubjectButtonImage.SetFromIconName(roomSubjectShownIconName, gtki.ICON_SIZE_BUTTON)
}

// onHideRoomSubject MUST be called from the UI thread
func (t *roomViewToolbar) onHideRoomSubject() {
	t.roomSubjectRevealer.SetRevealChild(false)
	t.roomSubjectButton.SetTooltipText(i18n.Local("Show room subject"))
	t.roomSubjectButtonImage.SetFromIconName(roomSubjectHiddenIconName, gtki.ICON_SIZE_BUTTON)
}
