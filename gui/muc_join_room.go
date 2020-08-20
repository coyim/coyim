package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type mucJoinRoomView struct {
	builder *builder

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

func (jrv *mucJoinRoomView) startSpinner() {
	jrv.spinner.Start()
	jrv.spinner.SetVisible(true)
}

func (jrv *mucJoinRoomView) stopSpinner() {
	jrv.spinner.Stop()
	jrv.spinner.SetVisible(false)
}

func (u *gtkUI) tryJoinRoom(jrv *mucJoinRoomView, a *account) {
	// TODO[OB]-MUC: I don't think using a mutex here is a good idea
	// Since this is in the UI thread, there are probably better ways to deal with it
	jrv.clearErrors()

	roomName, _ := jrv.txtRoomName.GetText()
	rj, ok := jid.Parse(roomName).(jid.Bare)
	if !ok {
		jrv.notifyOnError(i18n.Localf("\"%s\" is not a valid room identification", roomName))
		return
	}

	jrv.startSpinner()

	rc, ec := a.session.HasRoom(rj)
	go func() {
		select {
		case value, ok := <-rc:
			if !ok {
				doInUIThread(func() {
					jrv.stopSpinner()
					jrv.notifyOnError(i18n.Localf("An error ocurred trying to find the room \"%s\"", roomName))
					a.log.WithField("Room", roomName).Warn("An error ocurred trying to find a room")
				})
				return
			}
			doInUIThread(func() {
				jrv.stopSpinner()
				if !value {
					jrv.notifyOnError(i18n.Localf("The Room \"%s\" doesn't exist", roomName))
					a.log.WithField("Room", roomName).Debug("The Room doesn't exist")
				} else {
					jrv.dialog.Hide()
					u.mucShowRoom(a, rj)
				}
			})
		case err, ok := <-ec:
			if !ok {
				return
			}
			doInUIThread(func() {
				jrv.stopSpinner()
				if err != nil {
					jrv.notifyOnError(i18n.Localf("An error occurred trying to find the room \"%s\"", roomName))
					a.log.WithField("Room", roomName).WithError(err).Warn("Error occurred trying to find the Room")
				}
			})
		}
	}()
}

func (u *gtkUI) mucShowJoinRoom() {
	view := &mucJoinRoomView{}
	view.init()

	accountsInput := view.builder.get("accounts").(gtki.ComboBox)
	ac := u.createConnectedAccountsComponent(accountsInput, view, nil,
		func() {
			view.stopSpinner()
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

	view.dialog.Show()
}
