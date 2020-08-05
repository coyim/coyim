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
	errorNotif       *errorNotification
}

func (jrv *mucJoinRoomView) clearErrors() {
	jrv.errorNotif.Hide()
}

func (jrv *mucJoinRoomView) notifyOnError(errMessage string) {
	jrv.errorNotif.ShowMessage(errMessage)
}

func (jrv *mucJoinRoomView) init() {
	jrv.builder = newBuilder("MUCJoinRoomDialog")
	panicOnDevError(jrv.builder.bindObjects(jrv))
	jrv.errorNotif = newErrorNotification(jrv.notificationArea)
}

func (u *gtkUI) tryJoinRoom(jrv *mucJoinRoomView, a *account) {
	jrv.updateLock.Lock()

	doInUIThread(jrv.clearErrors)

	roomName, _ := jrv.txtRoomName.GetText()
	rj := jid.Parse(roomName).(jid.Bare)

	doInUIThread(func() {
		jrv.spinner.Start()
		jrv.spinner.SetVisible(true)
	})

	result := a.session.HasRoom(rj)
	go func() {
		value := <-result
		defer jrv.updateLock.Unlock()

		doInUIThread(func() {
			jrv.spinner.Stop()
			jrv.spinner.SetVisible(false)

			if !value {
				jrv.notifyOnError(i18n.Local(fmt.Sprintf("The Room \"%s\" doesn't exists", roomName)))
				a.log.Debug(fmt.Sprintf("The Room \"%s\" doesn't exists", roomName))
			} else {
				jrv.dialog.Hide()
				u.mucShowRoom(a, rj)
			}
		})
	}()
}

func (u *gtkUI) mucShowJoinRoom() {
	view := &mucJoinRoomView{}
	view.init()

	accountsInput := view.builder.get("accounts").(gtki.ComboBox)
	ac := u.createConnectedAccountsComponent(accountsInput, view, nil,
		func() {
			view.spinner.Stop()
			view.spinner.SetVisible(false)
		},
	)

	view.builder.ConnectSignals(map[string]interface{}{
		"on_close_window":        ac.onDestroy,
		"on_cancel_join_clicked": view.dialog.Destroy,
		"on_accept_join_clicked": func() {
			u.tryJoinRoom(view, ac.currentAccount())
		},
	})

	u.connectShortcutsChildWindow(view.dialog)

	view.dialog.SetTransientFor(u.window)
	view.dialog.Show()
}
