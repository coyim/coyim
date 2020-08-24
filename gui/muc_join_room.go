package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type mucJoinRoomView struct {
	u       *gtkUI
	builder *builder

	dialog           gtki.Dialog  `gtk-widget:"join-room"`
	roomNameEntry    gtki.Entry   `gtk-widget:"roomNameEntry"`
	joinButton       gtki.Button  `gtk-widget:"joinButton"`
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

func (jrv *mucJoinRoomView) startSpinner() {
	jrv.spinner.Start()
	jrv.spinner.SetVisible(true)
}

func (jrv *mucJoinRoomView) stopSpinner() {
	jrv.spinner.Stop()
	jrv.spinner.SetVisible(false)
}

func (jrv *mucJoinRoomView) hasValidRoomName() bool {
	jrv.clearErrors()
	roomName, _ := jrv.roomNameEntry.GetText()
	valid := jid.ValidBareJID(roomName)
	if !valid {
		if len(roomName) > 0 {
			jrv.notifyOnError(i18n.Localf("\"%s\" is not a valid room identification", roomName))
		}
	}
	return valid
}

func (jrv *mucJoinRoomView) validateInput() {
	sensitiveValue := jrv.hasValidRoomName()
	jrv.joinButton.SetSensitive(sensitiveValue)
}

func (jrv *mucJoinRoomView) notifyErrorServerUnavailable(a *account, roomName string) {
	jrv.stopSpinner()
	jrv.notifyOnError(i18n.Local("We can't get access to the server, please check your Internet connection or make sure the server exists."))
	a.log.WithField("room", roomName).Warn("An error ocurred trying to find a room")
}

func (jrv *mucJoinRoomView) tryJoinRoom(a *account) {
	// TODO[OB]-MUC: I don't think using a mutex here is a good idea
	// Since this is in the UI thread, there are probably better ways to deal with it
	jrv.clearErrors()

	roomName, _ := jrv.roomNameEntry.GetText()
	rj := jid.ParseBare(roomName)
	jrv.startSpinner()

	rc, ec := a.session.HasRoom(rj)
	go func() {
		select {
		case value, ok := <-rc:
			if !ok {
				doInUIThread(func() {
					jrv.notifyErrorServerUnavailable(a, roomName)
				})
				return
			}
			doInUIThread(func() {
				jrv.stopSpinner()
				if !value {
					jrv.notifyOnError(i18n.Localf("The room \"%s\" doesn't exist", roomName))
					a.log.WithField("room", roomName).Debug("The room doesn't exist")
				} else {
					jrv.dialog.Hide()
					jrv.u.mucShowRoom(a, rj)
				}
			})
		case err, ok := <-ec:
			if !ok {
				doInUIThread(func() {
					jrv.notifyErrorServerUnavailable(a, roomName)
				})
				return
			}
			doInUIThread(func() {
				jrv.stopSpinner()
				if err != nil {
					jrv.notifyOnError(i18n.Local("Looks like the server or the service doesn't exists, please verify the provided name."))
					a.log.WithField("room", roomName).WithError(err).Warn("An error occurred trying to find the room")
				}
			})
		}
	}()
}

func (jrv *mucJoinRoomView) init() {
	jrv.builder = newBuilder("MUCJoinRoomDialog")

	panicOnDevError(jrv.builder.bindObjects(jrv))

	accountsInput := jrv.builder.get("accounts").(gtki.ComboBox)
	ac := jrv.u.createConnectedAccountsComponent(accountsInput, jrv, nil, jrv.stopSpinner)

	jrv.builder.ConnectSignals(map[string]interface{}{
		"on_close_window":     ac.onDestroy,
		"on_roomname_changed": jrv.validateInput,
		"on_cancel_clicked":   jrv.dialog.Destroy,
		"on_join_clicked": func() {
			jrv.tryJoinRoom(ac.currentAccount())
		},
	})

	jrv.errorNotif = newErrorNotification(jrv.notificationArea)
}

func newMUCJoinRoomView(u *gtkUI) *mucJoinRoomView {
	view := &mucJoinRoomView{
		u: u,
	}

	view.init()

	u.connectShortcutsChildWindow(view.dialog)

	return view
}

func (u *gtkUI) mucShowJoinRoom() {
	view := newMUCJoinRoomView(u)
	view.dialog.Show()
}
