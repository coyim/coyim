package gui

import (
	"fmt"
	"sync"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type mucJoinRoomView struct {
	builder *builder

	generation int
	updateLock sync.RWMutex

	dialog           gtki.Dialog  `gtk-widget:"join-room"`
	txtRoomName      gtki.Entry   `gtk-widget:"textRoomName"`
	spinner          gtki.Spinner `gtk-widget:"spinner"`
	notificationArea gtki.Box     `gtk-widget:"boxNotificationArea"`
	notification     gtki.InfoBar
	errorNotif       *errorNotification
}

func (jrv *mucJoinRoomView) clearErrors() {
	jrv.errorNotif = newErrorNotification(jrv.notificationArea)
}

func (jrv *mucJoinRoomView) notifyOnError(errMessage string) {
	doInUIThread(func() {
		if jrv.notification != nil {
			jrv.notificationArea.Remove(jrv.notification)
		}

		jrv.errorNotif.ShowMessage(errMessage)
	})
}

func (jrv *mucJoinRoomView) init() {
	jrv.builder = newBuilder("MUCJoinRoomDialog")
	panicOnDevError(jrv.builder.bindObjects(jrv))
	jrv.errorNotif = newErrorNotification(jrv.notificationArea)
}

// tryJoinRoom find the room information and make the join to the room
func (u *gtkUI) tryJoinRoom(jrv *mucJoinRoomView, a *account) {
	jrv.updateLock.Lock()

	doInUIThread(jrv.clearErrors)

	roomName, _ := jrv.txtRoomName.GetText()
	rj := jid.Parse(roomName).(jid.Bare)

	doInUIThread(func() {
		jrv.spinner.Start()
		jrv.spinner.SetVisible(true)
	})

	value := a.session.HasRoom(rj)

	doInUIThread(func() {
		jrv.spinner.Stop()
		jrv.spinner.SetVisible(false)
	})

	jrv.updateLock.Unlock()

	if !value {
		jrv.notifyOnError(i18n.Local(fmt.Sprintf("The Room \"%s\" doesn't exists", roomName)))
		a.log.Debug(fmt.Sprintf("The Room \"%s\" doesn't exists", roomName))
	} else {
		doInUIThread(func() {
			u.mucShowRoom(a, rj)
			jrv.dialog.Hide()
		})
	}
}

//
// Custom GTK Events
//

func (jrv *mucJoinRoomView) onShowWindow() {

}

// mucJoinRoom should be called from the UI thread
func (u *gtkUI) mucShowJoinRoom() {
	view := &mucJoinRoomView{}
	view.init()

	accountsInput := view.builder.get("accounts").(gtki.ComboBox)
	ac := u.createConnectedAccountsComponent(accountsInput, view,
		func(*account) {},
		func() {
			view.spinner.Stop()
			view.spinner.SetVisible(false)
		},
	)

	view.builder.ConnectSignals(map[string]interface{}{
		"on_close_window": func() {},
		"on_show_window": func() {
			view.onShowWindow()
		},
		"on_cancel_join_clicked": view.dialog.Destroy,
		"on_accept_join_clicked": func() {
			u.tryJoinRoom(view, ac.currentAccount())
		},
	})

	u.connectShortcutsChildWindow(view.dialog)

	view.dialog.SetTransientFor(u.window)
	view.dialog.Show()
}
