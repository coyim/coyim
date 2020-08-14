package gui

import (
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
	jrv.errorNotif.Hide()
}

func (jrv *mucJoinRoomView) notifyOnError(err string) {
	if jrv.notification != nil {
		jrv.notificationArea.Remove(jrv.notification)
	}

	jrv.errorNotif.ShowMessage(err)
}

func (jrv *mucJoinRoomView) init() {
	jrv.builder = newBuilder("MUCJoinRoomDialog")
	panicOnDevError(jrv.builder.bindObjects(jrv))
	jrv.errorNotif = newErrorNotification(jrv.notificationArea)
}

func (u *gtkUI) tryJoinRoom(jrv *mucJoinRoomView, a *account) {
	jrv.updateLock.Lock()

	jrv.clearErrors()

	roomName, _ := jrv.txtRoomName.GetText()
	rj, ok := jid.Parse(roomName).(jid.Bare)
	if !ok {
		if len(roomName) == 0 {
			jrv.notifyOnError(i18n.Localf("Please specify a valid Room Name"))
		} else {
			jrv.notifyOnError(i18n.Localf("The Room \"%s\" is not a valid JID Bare format", roomName))
		}
		jrv.updateLock.Unlock()
		return
	}

	jrv.spinner.Start()
	jrv.spinner.SetVisible(true)

	r := a.session.HasRoom(rj)
	go func() {
		value := <-r
		defer jrv.updateLock.Unlock()

		jrv.spinner.Stop()
		jrv.spinner.SetVisible(false)
		doInUIThread(func() {
			if !value {
				jrv.notifyOnError(i18n.Localf("The Room \"%s\" doesn't exists", roomName))
				a.log.WithField("Room", roomName).Debug("The Room doesn't exists")
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
