package gui

import (
	"sync"

	"github.com/coyim/gotk3adapter/gtki"
)

type mucJoinPrivateRoomView struct {
	builder *builder

	generation int
	updateLock sync.RWMutex

	dialog              gtki.Dialog `gtk-widget:"MUCJoinPrivateRoom"`
	privateRoomNameText gtki.Entry  `gtk-widget:"txtPrivateRoomName"`
	notificationArea    gtki.Box    `gtk-widget:"boxNotificationArea"`
	notification        gtki.InfoBar
	errorNotif          *errorNotification
}

func (jrv *mucJoinPrivateRoomView) clearErrors() {
	jrv.errorNotif.Hide()
}

func (jrv *mucJoinPrivateRoomView) notifyOnError(errMessage string) {
	doInUIThread(func() {
		if jrv.notification != nil {
			jrv.notificationArea.Remove(jrv.notification)
		}

		jrv.errorNotif.ShowMessage(errMessage)
	})
}

func (jrv *mucJoinPrivateRoomView) init() {
	jrv.builder = newBuilder("MUCJoinPrivateRoomDialog")
	panicOnDevError(jrv.builder.bindObjects(jrv))
	//jrv.serviceGroups = make(map[string]gtki.TreeIter)
	jrv.errorNotif = newErrorNotification(jrv.notificationArea)
}

func (jrv *mucJoinPrivateRoomView) onShowWindow() {

}

func (jrv *mucJoinPrivateRoomView) onBtnJoinClicked() {

}

// mucJoinRoom should be called from the UI thread
func (u *gtkUI) mucShowJoinPrivateRoom() {
	view := &mucJoinPrivateRoomView{}

	view.init()

	view.builder.ConnectSignals(map[string]interface{}{
		"on_close_window_signal": func() {},
		"on_show_window_signal": func() {
			view.onShowWindow()
		},
		"on_btn_cancel_clicked_signal": view.dialog.Destroy,
		"on_btn_join_clicked_signal": func() {
			view.onBtnJoinClicked()
		},
	})

	u.connectShortcutsChildWindow(view.dialog)

	view.dialog.SetTransientFor(u.window)
	view.dialog.Show()
}
