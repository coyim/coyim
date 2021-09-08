package gui

import (
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/gtki"
)

// showCloseConfirmWindow MUST be called from the UI thread
func (v *roomView) showCloseConfirmWindow() {
	confirm := v.newRoomViewCloseWindowConfirm()
	confirm.showWindow()
}

type roomViewCloseWindowConfirm struct {
	u       *gtkUI
	room    *muc.Room
	account *account

	window gtki.Window `gtk-widget:"room-close-confirm-window"`

	log coylog.Logger
}

func (v *roomView) newRoomViewCloseWindowConfirm() *roomViewCloseWindowConfirm {
	confirm := &roomViewCloseWindowConfirm{
		u:       v.u,
		room:    v.room,
		account: v.account,
		log:     v.log.WithField("where", "roomViewCloseWindowConfirm"),
	}

	confirm.loadUIDefinition()

	return confirm
}

func (v *roomViewCloseWindowConfirm) connectUISignals(b *builder) {
	b.ConnectSignals(map[string]interface{}{
		"on_cancel":  v.onCancelClicked,
		"on_confirm": v.onConfirmClicked,
	})
}

func (v *roomViewCloseWindowConfirm) loadUIDefinition() {
	buildUserInterface("MUCRoomCloseWindowConfirm", v, v.connectUISignals)
}

// onCancelClicked MUST be called from the UI thread
func (v *roomViewCloseWindowConfirm) onCancelClicked() {
	v.closeWindow()
}

// onConfirmClicked MUST be called from the UI thread
func (v *roomViewCloseWindowConfirm) onConfirmClicked() {
	go v.tryLeaveRoom()
	v.closeWindow()
}

// tryLeaveRoom MUST NOT be called from the UI thread
func (v *roomViewCloseWindowConfirm) tryLeaveRoom() {
	onError := func(err error) {
		v.log.WithError(err).Error("An error occurred when trying to leave the room")
	}

	v.account.leaveRoom(
		v.room.ID,
		v.room.SelfOccupantNickname(),
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
