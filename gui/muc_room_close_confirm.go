package gui

import (
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/gtki"
)

// showCloseConfirmWindow MUST be called from the UI thread
func (v *roomView) showCloseConfirmWindow() {
	confirm := v.newRoomViewCloseWindowConfirm()
	confirm.showWindow()
}

type roomViewCloseWindowConfirm struct {
	roomView *roomView

	window         gtki.Window      `gtk-widget:"room-close-confirm-window"`
	icon           gtki.Image       `gtk-widget:"room-close-confirm-icon"`
	leaveRoomCheck gtki.CheckButton `gtk-widget:"room-close-confirm-leave-checkbox"`
	cancelButton   gtki.Button      `gtk-widget:"room-close-cancel-button"`
	confirmButton  gtki.Button      `gtk-widget:"room-close-confirm-button"`

	log coylog.Logger
}

func (v *roomView) newRoomViewCloseWindowConfirm() *roomViewCloseWindowConfirm {
	confirm := &roomViewCloseWindowConfirm{
		roomView: v,
		log:      v.log.WithField("where", "roomViewCloseWindowConfirm"),
	}

	confirm.loadUIDefinition()
	confirm.icon.SetFromPixbuf(getMUCIconPixbuf("dialog_question"))
	confirm.cancelButton.GrabFocus()
	confirm.window.SetTransientFor(v.mainWindow())

	v.u.connectShortcutsMucRoomWindow(confirm.window, func(_ gtki.Window) {
		confirm.closeWindow()
	})

	return confirm
}

func (v *roomViewCloseWindowConfirm) connectUISignals(b *builder) {
	b.ConnectSignals(map[string]interface{}{
		"on_leave_room_check_changed": v.onLeaveRoomCheckChanged,
		"on_cancel":                   v.onCancelClicked,
		"on_confirm":                  v.onConfirmClicked,
	})
}

func (v *roomViewCloseWindowConfirm) loadUIDefinition() {
	buildUserInterface("MUCRoomCloseWindowConfirm", v, v.connectUISignals)
}

// onLeaveRoomCheckChanged MUST be called from the UI thread
func (v *roomViewCloseWindowConfirm) onLeaveRoomCheckChanged() {
	buttonLabel := i18n.Local("Close Room")
	if v.leaveRoomCheck.GetActive() {
		buttonLabel = i18n.Local("Close & Leave Room")
	}
	v.confirmButton.SetLabel(buttonLabel)
}

// onCancelClicked MUST be called from the UI thread
func (v *roomViewCloseWindowConfirm) onCancelClicked() {
	v.closeWindow()
}

// onConfirmClicked MUST be called from the UI thread
func (v *roomViewCloseWindowConfirm) onConfirmClicked() {
	v.closeWindow()

	if v.leaveRoomCheck.GetActive() {
		go v.tryLeaveRoom()
	}

	v.roomView.close()
}

// tryLeaveRoom MUST NOT be called from the UI thread
func (v *roomViewCloseWindowConfirm) tryLeaveRoom() {
	onError := func(err error) {
		v.log.WithError(err).Error("An error occurred when trying to leave the room")
	}

	v.roomView.account.leaveRoom(
		v.roomView.room.ID,
		v.roomView.room.SelfOccupantNickname(),
		nil,
		onError,
		nil,
	)
}

// showWindow MUST be called from the UI thread
func (v *roomViewCloseWindowConfirm) showWindow() {
	v.window.Show()
}

// closeWindow MUST be called from the UI thread
func (v *roomViewCloseWindowConfirm) closeWindow() {
	v.window.Destroy()
}
