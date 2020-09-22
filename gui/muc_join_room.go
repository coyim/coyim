package gui

import (
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type mucJoinRoomView struct {
	u       *gtkUI
	builder *builder
	ac      *connectedAccountsComponent

	dialog           gtki.Dialog  `gtk-widget:"join-chat-dialog"`
	roomNameEntry    gtki.Entry   `gtk-widget:"room-name-entry"`
	chatServiceEntry gtki.Entry   `gtk-widget:"chat-service-entry"`
	joinButton       gtki.Button  `gtk-widget:"join-button"`
	spinner          gtki.Spinner `gtk-widget:"spinner"`
	notificationArea gtki.Box     `gtk-widget:"box-notification-area"`

	notification gtki.InfoBar
	errorNotif   *errorNotification
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

func (v *mucJoinRoomView) showSpinner() {
	v.spinner.Start()
	v.spinner.Show()
}

func (v *mucJoinRoomView) hideSpinner() {
	v.spinner.Stop()
	v.spinner.Hide()
}

func (v *mucJoinRoomView) typedRoomName() string {
	name, _ := v.roomNameEntry.GetText()
	return name
}

// enableJoinIfConditionsAreMet SHOULD be called from the UI thread
func (v *mucJoinRoomView) enableJoinIfConditionsAreMet() {
	roomName, _ := v.roomNameEntry.GetText()
	chatServiceName, _ := v.chatServiceEntry.GetText()

	hasAllValues := len(roomName) != 0 && len(chatServiceName) != 0 && v.ac.currentAccount() != nil
	v.joinButton.SetSensitive(hasAllValues)
}

func (v *mucJoinRoomView) onBeforeStart() {
	v.clearErrors()
	v.disableJoinFields()
	v.showSpinner()
}

func (v *mucJoinRoomView) onJoinSuccess(a *account, ident jid.Bare, roomInfo *muc.RoomListing) {
	doInUIThread(func() {
		v.hideSpinner()
		v.dialog.Hide()
		v.u.joinMultiUserChat(a, ident, v.returnWhenCancelJoining)
	})
}

func (v *mucJoinRoomView) returnWhenCancelJoining() {
	v.enableJoinFields()
	v.dialog.Show()
}

func (v *mucJoinRoomView) onJoinFails(a *account, ident jid.Bare) {
	doInUIThread(func() {
		v.notifyOnError(i18n.Local("The room doesn't exist on that service."))
		v.enableJoinFields()
		v.hideSpinner()
	})
	a.log.WithField("room", ident).Warn("The room doesn't exist")
}

func (v *mucJoinRoomView) onJoinError(a *account, ident jid.Bare, err error) {
	doInUIThread(func() {
		v.hideSpinner()
		v.enableJoinFields()
		if err != nil {
			v.notifyOnError(i18n.Local("It looks like the room you are trying to connect to doesn't exist, please verify the provided information."))
			a.log.WithField("room", ident).WithError(err).Warn("An error occurred trying to find the room")
		}
	})
}

func (v *mucJoinRoomView) onServiceUnavailable(a *account, ident jid.Bare) {
	doInUIThread(func() {
		v.hideSpinner()
		v.notifyOnError(i18n.Local("We can't get access to the service, please check your Internet connection or make sure the service exists."))
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

func (c *mucJoinRoomContext) waitToFinish(resultChannel <-chan bool, errorChannel <-chan error, roomInfo <-chan *muc.RoomListing) {
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

		ri := <-roomInfo

		c.v.onJoinSuccess(c.a, c.ident, ri)
	case err, ok := <-errorChannel:
		c.onFinishWithError(err, ok)
	}
}

func (c *mucJoinRoomContext) exec() {
	c.v.onBeforeStart()
	roomInfo := make(chan *muc.RoomListing)
	resultChannel, errorChannel := c.a.session.HasRoom(c.ident, roomInfo)
	go c.waitToFinish(resultChannel, errorChannel, roomInfo)
}

func (v *mucJoinRoomView) log() coylog.Logger {
	ca := v.ac.currentAccount()
	if ca != nil {
		return ca.log
	}
	return v.u.log
}

func (v *mucJoinRoomView) validateFieldsAndGetBareIfOk() (jid.Bare, bool) {
	roomName, err := v.roomNameEntry.GetText()
	if err != nil {
		v.log().WithField("name", roomName).Error("Trying to join a room with an invalid room name")
		v.notifyOnError(i18n.Local("Could not get the room name, please try again."))
		return nil, false
	}

	local := jid.NewLocal(roomName)
	if !local.Valid() {
		v.log().WithField("local", roomName).Error("Trying to join a room with an invalid local")
		v.notifyOnError(i18n.Local("You must provide a valid room name."))
		return nil, false
	}

	chatServiceName, err := v.chatServiceEntry.GetText()
	if err != nil {
		v.log().WithError(err).Error("Something went wrong while trying to join the room")
		v.notifyOnError(i18n.Local("Could not get the service name, please try again."))
		return nil, false
	}

	domain := jid.NewDomain(chatServiceName)
	if !domain.Valid() {
		v.log().WithField("domain", chatServiceName).Error("Trying to join a room with an invalid domain")
		v.notifyOnError(i18n.Local("You must provide a valid service name."))
		return nil, false
	}

	return jid.NewBare(local, domain), true
}

func (v *mucJoinRoomView) tryJoinRoom(done func()) {
	ident, isValid := v.validateFieldsAndGetBareIfOk()
	if !isValid {
		done()
		return
	}

	ca := v.ac.currentAccount()
	if ca == nil {
		v.notifyOnError(i18n.Local("No account was selected, select an account from the list or enable one."))
		return
	}

	c := &mucJoinRoomContext{
		a:     ca,
		v:     v,
		ident: ident,
		done:  done,
	}

	c.exec()
}

func doOnlyOnceAtATime(f func(func())) func() {
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
}

func (v *mucJoinRoomView) init() {
	v.builder = newBuilder("MUCJoinRoomDialog")

	panicOnDevError(v.builder.bindObjects(v))

	v.errorNotif = newErrorNotification(v.notificationArea)

	accountsInput := v.builder.get("accounts").(gtki.ComboBox)
	v.ac = v.u.createConnectedAccountsComponent(accountsInput, v, func(a *account) {
		doInUIThread(v.enableJoinIfConditionsAreMet)
	}, func() {
		doInUIThread(v.enableJoinIfConditionsAreMet)
	})

	v.builder.ConnectSignals(map[string]interface{}{
		"on_close_window":        v.ac.onDestroy,
		"on_roomname_changed":    v.enableJoinIfConditionsAreMet,
		"on_chatService_changed": v.enableJoinIfConditionsAreMet,
		"on_nickName_changed":    v.enableJoinIfConditionsAreMet,
		"on_cancel_clicked":      v.dialog.Destroy,
		"on_join_clicked":        doOnlyOnceAtATime(v.tryJoinRoom),
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

func (v *mucJoinRoomView) isValidRoomName(name string) bool {
	bare, isValid := jid.TryParseBare(name)
	if !isValid {
		return false
	}
	return jid.NewLocal(bare.Local().String()).Valid() && jid.NewDomain(bare.Host().String()).Valid()
}

func (v *mucJoinRoomView) disableJoinFields() {
	v.roomNameEntry.SetSensitive(false)
	v.chatServiceEntry.SetSensitive(false)
	v.joinButton.SetSensitive(false)
	v.ac.disableAccountInput()
}

func (v *mucJoinRoomView) enableJoinFields() {
	v.roomNameEntry.SetSensitive(true)
	v.chatServiceEntry.SetSensitive(true)
	v.joinButton.SetSensitive(true)
	v.ac.enableAccountInput()
}

func (u *gtkUI) mucShowJoinRoom() {
	view := newMUCJoinRoomView(u)

	u.connectShortcutsChildWindow(view.dialog)

	view.dialog.SetTransientFor(u.window)
	view.dialog.Show()
}
