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
	roomView         *roomView
	isEditingSubject bool

	view                       gtki.Box            `gtk-widget:"room-view-toolbar"`
	roomNameLabel              gtki.Label          `gtk-widget:"room-name-label"`
	roomStatusIcon             gtki.Image          `gtk-widget:"room-status-icon"`
	roomMenuButton             gtki.MenuButton     `gtk-widget:"room-menu-button"`
	roomSubjectButton          gtki.Button         `gtk-widget:"room-subject-button"`
	roomSubjectButtonImage     gtki.Image          `gtk-widget:"room-subject-button-image"`
	roomSubjectRevealer        gtki.Revealer       `gtk-widget:"room-subject-revealer"`
	roomSubjectLabel           gtki.Label          `gtk-widget:"room-subject-label"`
	roomSubjectScrolledWindow  gtki.ScrolledWindow `gtk-widget:"room-subject-editable-content"`
	roomSubjectTextView        gtki.TextView       `gtk-widget:"room-subject-textview"`
	roomSubjectEditButton      gtki.Button         `gtk-widget:"room-edit-subject-button"`
	roomSubjectButtonBox       gtki.Box            `gtk-widget:"room-edit-subject-buttons-box"`
	roomSubjectApplyButton     gtki.Button         `gtk-widget:"room-edit-subject-apply-button"`
	securityPropertiesMenuItem gtki.MenuItem       `gtk-widget:"security-properties-menu-item"`
	configureRoomMenuItem      gtki.MenuItem       `gtk-widget:"room-configuration-menu-item"`
	modifyBanMenuItem          gtki.MenuItem       `gtk-widget:"modify-ban-list-menu-item"`
	destroyRoomMenuItem        gtki.MenuItem       `gtk-widget:"destroy-room-menu-item"`
	leaveRoomMenuItem          gtki.MenuItem       `gtk-widget:"leave-room-menu-item"`
}

func (v *roomView) newRoomViewToolbar() *roomViewToolbar {
	t := &roomViewToolbar{
		roomView: v,
	}

	t.initBuilder()
	t.initDefaults()
	t.initSubscribers()

	return t
}

func (t *roomViewToolbar) initBuilder() {
	builder := newBuilder("MUCRoomToolbar")
	panicOnDevError(builder.bindObjects(t))

	builder.ConnectSignals(map[string]interface{}{
		"on_edit_room_subject":        t.onEditRoomSubject,
		"on_cancel_room_subject_edit": t.onCancelEditSubject,
		"on_apply_room_subject_edit":  t.onApplyEditSubject,
		"on_subject_changed":          t.onRoomSubectChanged,
		"on_toggle_room_subject":      t.onToggleRoomSubject,
		"on_show_security_properties": t.roomView.showWarnings,
		"on_configure_room":           t.roomView.onConfigureRoom,
		"on_modify_ban_list":          t.roomView.onModifyBanList,
		"on_destroy_room":             t.roomView.onDestroyRoom,
		"on_leave_room":               t.roomView.onLeaveRoom,
	})
}

func (t *roomViewToolbar) initDefaults() {
	t.roomStatusIcon.SetFromPixbuf(getMUCIconPixbuf("room"))

	t.roomNameLabel.SetText(t.roomView.roomID().String())
	mucStyles.setRoomToolbarNameLabelStyle(t.roomNameLabel)

	mucStyles.setRoomToolbarSubjectLabelStyle(t.roomSubjectLabel)
}

func (t *roomViewToolbar) initSubscribers() {
	t.roomView.subscribe("toolbar", func(ev roomViewEvent) {
		switch e := ev.(type) {
		case subjectReceivedEvent:
			t.subjectReceivedEvent(e.subject)
		case subjectUpdatedEvent:
			t.subjectUpdatedEvent(e.subject)
		case roomDestroyedEvent:
			t.roomDestroyedEvent()
		case selfOccupantRemovedEvent:
			t.selfOccupantRemovedEvent()
		case occupantSelfJoinedEvent:
			t.selfOccupantJoinedEvent()
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
	t.subjectUpdatedEvent(subject)
	doInUIThread(t.handleSubjectComponents)
}

func (t *roomViewToolbar) subjectUpdatedEvent(subject string) {
	doInUIThread(func() {
		t.displayRoomSubject(subject)
		t.handleSubjectButtonVisibility()
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

func (t *roomViewToolbar) selfOccupantJoinedEvent() {
	doInUIThread(func() {
		t.updateMenuActionsBasedOn(t.roomView.room.SelfOccupant().Affiliation)
	})
}

func (t *roomViewToolbar) updateMenuActionsBasedOn(affiliation data.Affiliation) {
	t.configureRoomMenuItem.SetVisible(affiliation.IsOwner())
	t.destroyRoomMenuItem.SetVisible(affiliation.IsOwner())

	showAdminActions := affiliation.IsOwner() || affiliation.IsAdmin()
	t.modifyBanMenuItem.SetVisible(showAdminActions)
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
}

// onToggleRoomSubject MUST be called from the UI thread
func (t *roomViewToolbar) onToggleRoomSubject() {
	if t.roomSubjectRevealer.GetRevealChild() {
		t.onHideRoomSubject()
	} else {
		t.onShowRoomSubject()
	}
}

// onEditRoomSubject MUST be called from the UI thread
func (t *roomViewToolbar) onEditRoomSubject() {
	t.isEditingSubject = true
	t.toggleEditSubjectComponents(false)

	setTextViewText(t.roomSubjectTextView, t.roomSubjectLabel.GetLabel())
}

// onCancelEditSubject MUST be called from the UI thread
func (t *roomViewToolbar) onCancelEditSubject() {
	if t.isEditingSubject {
		t.toggleEditSubjectComponents(true)
		return
	}

	t.resetSubjectComponents()
}

// resetSubjectComponents MUST be called from the UI thread
func (t *roomViewToolbar) resetSubjectComponents() {
	setTextViewText(t.roomSubjectTextView, "")
	t.onToggleRoomSubject()
}

// onApplyEditSubject MUST be called from the UI thread
func (t *roomViewToolbar) onApplyEditSubject() {
	newSubject := getTextViewText(t.roomSubjectTextView)
	t.roomView.updateSubjectRoom(newSubject,
		func() {
			t.roomSubjectLabel.SetText(newSubject)
			t.toggleEditSubjectComponents(true)
		})
}

// onRoomSubectChanged MUST be called from the UI thread
func (t *roomViewToolbar) onRoomSubectChanged() {
	t.roomSubjectApplyButton.SetSensitive(t.roomView.room.GetSubject() != getTextViewText(t.roomSubjectTextView))
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

func (t *roomViewToolbar) toggleEditSubjectComponents(v bool) {
	t.roomSubjectLabel.SetVisible(v)
	t.roomSubjectScrolledWindow.SetVisible(!v)
	t.roomSubjectEditButton.SetVisible(v)
	t.roomSubjectButtonBox.SetVisible(!v)
}

// handleSubjectComponents MUST be called from the UI thread
func (t *roomViewToolbar) handleSubjectComponents() {
	t.handleSubjectButtonVisibility()

	if t.roomView.room.CanChangeSubject() {
		t.toggleEditSubjectComponents(t.roomView.room.HasSubject())
		t.roomSubjectButton.Show()
	}
}

// handleSubjectButtonVisibility MUST be called from the UI thread
func (t *roomViewToolbar) handleSubjectButtonVisibility() {
	t.roomSubjectButton.Hide()
	if t.roomView.room.HasSubject() {
		t.roomSubjectButton.Show()
	}
}
