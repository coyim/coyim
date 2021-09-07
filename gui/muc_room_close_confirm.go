package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/gtki"
)

// showCloseConfirmWindow MUST be called from the UI thread
func (v *roomView) showCloseConfirmWindow() {
	confirm := v.newRoomViewCloseWindowConfirm()
	confirm.showWindow()
}

type roomViewCloseWindowConfirm struct {
	u *gtkUI

	headerLabel  gtki.Label  `gtk-widget:"room-close-confirm-header"`
	messageLabel gtki.Label  `gtk-widget:"room-close-confirm-message"`
	window       gtki.Window `gtk-widget:"room-close-confirm-window"`
}

func (v *roomView) newRoomViewCloseWindowConfirm() *roomViewCloseWindowConfirm {
	confirm := &roomViewCloseWindowConfirm{
		u: v.u,
	}

	confirm.loadUIDefinition()

	confirm.messageLabel.SetText(i18n.Localf("Hi, we have seen that you closed the room %s.\n"+
		"What action do you would like to do?", v.roomID()))

	mucStyles.setRoomCloseWindowConfirmHeaderStyle(confirm.headerLabel)

	return confirm
}

func (v *roomViewCloseWindowConfirm) connectUISignals(b *builder) {
	b.ConnectSignals(map[string]interface{}{
		"on_return_to_room": v.onReturnToRoomClicked,
		"on_leave_room":     v.onLeaveRoomClicked,
		"on_keep_in_room":   v.onKeepInRoomClicked,
	})
}

func (v *roomViewCloseWindowConfirm) loadUIDefinition() {
	buildUserInterface("MUCRoomCloseWindowConfirm", v, v.connectUISignals)
}

// onReturnToRoomClicked MUST be called from the UI thread
func (v *roomViewCloseWindowConfirm) onReturnToRoomClicked() {
	v.closeWindow()
}

// onLeaveRoomClicked MUST be called from the UI thread
func (v *roomViewCloseWindowConfirm) onLeaveRoomClicked() {
	v.closeWindow()
}

// onKeepInRoomClicked MUST be called from the UI thread
func (v *roomViewCloseWindowConfirm) onKeepInRoomClicked() {
	v.closeWindow()
}

// showWindow MUST be called from the UI thread
func (v *roomViewCloseWindowConfirm) showWindow() {
	v.window.Show()
}

// closeWindow MUST be called from the UI thread
func (v *roomViewCloseWindowConfirm) closeWindow() {
	v.window.Destroy()
}
