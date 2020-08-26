package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type mucJoinRoomView struct {
	u       *gtkUI
	builder *builder
	ac      *connectedAccountsComponent

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

func (v *mucJoinRoomView) typedRoomName() string {
	name, _ := v.roomNameEntry.GetText()
	return name
}

func (v *mucJoinRoomView) isValidRoomName(name string) bool {
	return jid.ValidBareJID(name)
}

func (v *mucJoinRoomView) hasValidRoomName() bool {
	if len(v.ac.accounts) == 0 {
		return false
	}

	v.clearErrors()

	roomName := v.typedRoomName()
	if len(roomName) == 0 {
		return false
	}

	if !v.isValidRoomName(roomName) {
		v.notifyOnError(i18n.Localf("\"%s\" is not a valid room identification", roomName))
		return false
	}

	return true
}

func (v *mucJoinRoomView) validateInput() {
	v.joinButton.SetSensitive(v.hasValidRoomName())
}

func (v *mucJoinRoomView) onBeforeStart() {
	v.clearErrors()
	v.disableJoinFields()
	v.startSpinner()
}

func (v *mucJoinRoomView) onJoinSuccess(a *account, ident jid.Bare) {
	doInUIThread(func() {
		v.stopSpinner()
		v.dialog.Hide()
		v.u.mucShowRoom(a, ident)
	})
}

func (v *mucJoinRoomView) onJoinFails(a *account, ident jid.Bare) {
	doInUIThread(func() {
		v.notifyOnError(i18n.Localf("The room \"%s\" doesn't exist", ident))
		v.enableJoinFields()
	})
	a.log.WithField("room", ident).Warn("The room doesn't exist")
}

func (v *mucJoinRoomView) onJoinError(a *account, ident jid.Bare, err error) {
	doInUIThread(func() {
		v.stopSpinner()
		v.enableJoinFields()
		if err != nil {
			v.notifyOnError(i18n.Local("Looks like the room you are trying to connect to doesn't exists, please verify the provided name."))
			a.log.WithField("room", ident).WithError(err).Warn("An error occurred trying to find the room")
		}
	})
}

func (v *mucJoinRoomView) onServiceUnavailable(a *account, ident jid.Bare) {
	doInUIThread(func() {
		v.stopSpinner()
		v.notifyOnError(i18n.Local("We can't get access to the server, please check your Internet connection or make sure the server exists."))
	})
	a.log.WithField("room", ident).Warn("An error ocurred trying to find the room")
}

type mucJoinRoomContext struct {
	a     *account
	v     *mucJoinRoomView
	ident jid.Bare
	done  func()
}

func (c *mucJoinRoomContext) onFinishWithError(err error, isErrorChannelClosed bool) {
	if !isErrorChannelClosed {
		c.v.onServiceUnavailable(c.a, c.ident)
		return
	}
	c.v.onJoinError(c.a, c.ident, err)
}

func (c *mucJoinRoomContext) waitToFinish(resultChannel <-chan bool, errorChannel <-chan error) {
	defer doInUIThread(c.done)

	select {
	case value, ok := <-resultChannel:
		if !ok {
			c.v.onServiceUnavailable(c.a, c.ident)
			return
		}

		if !value {
			c.v.onJoinFails(c.a, c.ident)
			return
		}

		c.v.onJoinSuccess(c.a, c.ident)
	case err, ok := <-errorChannel:
		c.onFinishWithError(err, ok)
	}
}

func (c *mucJoinRoomContext) exec() {
	c.v.onBeforeStart()
	resultChannel, errorChannel := c.a.session.HasRoom(c.ident)
	go c.waitToFinish(resultChannel, errorChannel)
}

func (v *mucJoinRoomView) tryJoinRoom(done func()) {
	ca := v.ac.currentAccount()
	if ca == nil {
		v.notifyOnError(i18n.Local("No account was selected, select an account from the list or enable one."))
		return
	}

	c := &mucJoinRoomContext{
		a:     ca,
		v:     v,
		ident: jid.ParseBare(v.typedRoomName()),
		done:  done,
	}

	c.exec()
}

func doOnlyOnceAtATime(f func(func())) func() {
	return func() func() {
		isDoing := false
		return func() {
			if isDoing {
				return
			}
			isDoing = true
			// The "done" function should be called ONLY from the UI thread,
			// in other cases it's not "safe" executing it.
			f(func() {
				isDoing = false
			})
		}
	}()
}

func (v *mucJoinRoomView) init() {
	v.builder = newBuilder("MUCJoinRoomDialog")

	panicOnDevError(v.builder.bindObjects(v))

	v.errorNotif = newErrorNotification(v.notificationArea)

	accountsInput := v.builder.get("accounts").(gtki.ComboBox)
	v.ac = v.u.createConnectedAccountsComponent(accountsInput, v, func(a *account) {
		doInUIThread(v.validateInput)
	}, func() {
		doInUIThread(func() {
			v.stopSpinner()
			v.validateInput()
		})
	})

	v.builder.ConnectSignals(map[string]interface{}{
		"on_close_window":     v.ac.onDestroy,
		"on_roomname_changed": v.validateInput,
		"on_cancel_clicked":   v.dialog.Destroy,
		"on_join_clicked":     doOnlyOnceAtATime(v.tryJoinRoom),
	})
}

func newMUCJoinRoomView(u *gtkUI) *mucJoinRoomView {
	view := &mucJoinRoomView{
		u: u,
	}

	view.init()

	u.connectShortcutsChildWindow(view.dialog)

	return view
}

func (v *mucJoinRoomView) disableJoinFields() {
	v.joinButton.SetSensitive(false)
	v.roomNameEntry.SetSensitive(false)
	v.ac.disableAccountInput()
}

func (v *mucJoinRoomView) enableJoinFields() {
	v.joinButton.SetSensitive(true)
	v.roomNameEntry.SetSensitive(true)
	v.ac.enableAccountInput()
}

func (u *gtkUI) mucShowJoinRoom() {
	view := newMUCJoinRoomView(u)
	view.dialog.Show()
}
