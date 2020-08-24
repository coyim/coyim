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
	doInUIThread(func() {
		v.clearErrors()
		v.startSpinner()
	})
}

func (v *mucJoinRoomView) onJoinSuccess(a *account, ident jid.Bare, s bool) {
	doInUIThread(v.stopSpinner)

	if !s {
		doInUIThread(func() {
			v.notifyOnError(i18n.Localf("The room \"%s\" doesn't exist", ident))
		})
		a.log.WithField("room", ident).Debug("The room doesn't exist")
		return
	}

	doInUIThread(func() {
		v.dialog.Hide()
		v.u.mucShowRoom(a, ident)
	})
}

func (v *mucJoinRoomView) onJoinError(a *account, ident jid.Bare, err error) {
	doInUIThread(func() {
		v.stopSpinner()
		if err != nil {
			v.notifyOnError(i18n.Local("Looks like the server or the service doesn't exists, please verify the provided name."))
			a.log.WithField("room", ident).WithError(err).Warn("An error occurred trying to find the room")
		}
	})
}

func (v *mucJoinRoomView) onJoinServerUnavailable(a *account, ident jid.Bare) {
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
	// onBeforeStart IS called from the UI thread
	onBeforeStart func()
	// onSuccess is NOT called from the UI thread
	onSuccess func(*account, jid.Bare, bool)
	// onError is NOT called from the UI thread
	onError func(*account, jid.Bare, error)
	// onServiceUnavailable is NOT called from the UI thread
	onServiceUnavailable func(*account, jid.Bare)
}

func (c *mucJoinRoomContext) onFinishWithResult(s, isChannelClosed bool) {
	if !isChannelClosed {
		c.onServiceUnavailable(c.a, c.ident)
		return
	}
	c.onSuccess(c.a, c.ident, s)
}

func (c *mucJoinRoomContext) onFinishWithError(err error, isErrorChannelClosed bool) {
	if !isErrorChannelClosed {
		c.onServiceUnavailable(c.a, c.ident)
		return
	}
	c.onError(c.a, c.ident, err)
}

func (c *mucJoinRoomContext) waitToFinish(resultChannel <-chan bool, errorChannel <-chan error) {
	select {
	case value, ok := <-resultChannel:
		c.onFinishWithResult(value, ok)
	case err, ok := <-errorChannel:
		c.onFinishWithError(err, ok)
	}
}

func (c *mucJoinRoomContext) exec() {
	c.onBeforeStart()
	resultChannel, errorChannel := c.a.session.HasRoom(c.ident)
	go c.waitToFinish(resultChannel, errorChannel)
}

func (v *mucJoinRoomView) tryJoinRoom() {
	// TODO[OB]-MUC: I don't think using a mutex here is a good idea
	// Since this is in the UI thread, there are probably better ways to deal with it
	c := &mucJoinRoomContext{
		a:                    v.ac.currentAccount(),
		v:                    v,
		ident:                jid.ParseBare(v.typedRoomName()),
		onBeforeStart:        v.onBeforeStart,
		onSuccess:            v.onJoinSuccess,
		onError:              v.onJoinError,
		onServiceUnavailable: v.onJoinServerUnavailable,
	}

	c.exec()
}

func (v *mucJoinRoomView) init() {
	v.builder = newBuilder("MUCJoinRoomDialog")

	panicOnDevError(v.builder.bindObjects(v))

	v.errorNotif = newErrorNotification(v.notificationArea)

	accountsInput := v.builder.get("accounts").(gtki.ComboBox)
	v.ac = v.u.createConnectedAccountsComponent(accountsInput, v, func(a *account) {
		doInUIThread(v.validateInput)
	}, func() {
		doInUIThread(v.stopSpinner)
	})

	v.builder.ConnectSignals(map[string]interface{}{
		"on_close_window":     v.ac.onDestroy,
		"on_roomname_changed": v.validateInput,
		"on_cancel_clicked":   v.dialog.Destroy,
		"on_join_clicked":     v.tryJoinRoom,
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

func (u *gtkUI) mucShowJoinRoom() {
	view := newMUCJoinRoomView(u)
	view.dialog.Show()
}
