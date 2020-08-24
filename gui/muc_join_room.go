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

func (v *mucJoinRoomView) clearErrors() {
	v.errorNotif.Hide()
}

func (v *mucJoinRoomView) notifyOnError(err string) {
	if v.notification != nil {
		v.notificationArea.Remove(v.notification)
	}

	v.errorNotif.ShowMessage(err)
}

func (v *mucJoinRoomView) startSpinner() {
	v.spinner.Start()
	v.spinner.SetVisible(true)
}

func (v *mucJoinRoomView) stopSpinner() {
	v.spinner.Stop()
	v.spinner.SetVisible(false)
}

func (v *mucJoinRoomView) hasValidRoomName() bool {
	v.clearErrors()
	roomName, _ := v.roomNameEntry.GetText()
	valid := jid.ValidBareJID(roomName)
	if !valid {
		if len(roomName) > 0 {
			v.notifyOnError(i18n.Localf("\"%s\" is not a valid room identification", roomName))
		}
	}
	return valid
}

func (v *mucJoinRoomView) validateInput() {
	sensitiveValue := v.hasValidRoomName()
	v.joinButton.SetSensitive(sensitiveValue)
}

func (v *mucJoinRoomView) notifyErrorServerUnavailable(a *account, roomName string) {
	v.stopSpinner()
	v.notifyOnError(i18n.Local("We can't get access to the server, please check your Internet connection or make sure the server exists."))
	a.log.WithField("room", roomName).Warn("An error ocurred trying to find a room")
}

func (v *mucJoinRoomView) tryJoinRoom(a *account) {
	// TODO[OB]-MUC: I don't think using a mutex here is a good idea
	// Since this is in the UI thread, there are probably better ways to deal with it
	v.clearErrors()

	roomName, _ := v.roomNameEntry.GetText()
	rj := jid.ParseBare(roomName)
	v.startSpinner()

	rc, ec := a.session.HasRoom(rj)
	go func() {
		select {
		case value, ok := <-rc:
			if !ok {
				doInUIThread(func() {
					v.notifyErrorServerUnavailable(a, roomName)
				})
				return
			}
			doInUIThread(func() {
				v.stopSpinner()
				if !value {
					v.notifyOnError(i18n.Localf("The room \"%s\" doesn't exist", roomName))
					a.log.WithField("room", roomName).Debug("The room doesn't exist")
				} else {
					v.dialog.Hide()
					v.u.mucShowRoom(a, rj)
				}
			})
		case err, ok := <-ec:
			if !ok {
				doInUIThread(func() {
					v.notifyErrorServerUnavailable(a, roomName)
				})
				return
			}
			doInUIThread(func() {
				v.stopSpinner()
				if err != nil {
					v.notifyOnError(i18n.Local("Looks like the server or the service doesn't exists, please verify the provided name."))
					a.log.WithField("room", roomName).WithError(err).Warn("An error occurred trying to find the room")
				}
			})
		}
	}()
}

func (v *mucJoinRoomView) init() {
	v.builder = newBuilder("MUCJoinRoomDialog")

	panicOnDevError(v.builder.bindObjects(v))

	accountsInput := v.builder.get("accounts").(gtki.ComboBox)
	ac := v.u.createConnectedAccountsComponent(accountsInput, v, nil, v.stopSpinner)

	v.builder.ConnectSignals(map[string]interface{}{
		"on_close_window":     ac.onDestroy,
		"on_roomname_changed": v.validateInput,
		"on_cancel_clicked":   v.dialog.Destroy,
		"on_join_clicked": func() {
			v.tryJoinRoom(ac.currentAccount())
		},
	})

	v.errorNotif = newErrorNotification(v.notificationArea)
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
