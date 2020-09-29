package gui

import (
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type mucJoinRoomView struct {
	u  *gtkUI
	ac *connectedAccountsComponent

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

// TODO: Maybe we should extract some spinner logic into a component - we do this in a lot of places

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

// enableJoinIfConditionsAreMet MUST be called from the UI thread
func (v *mucJoinRoomView) enableJoinIfConditionsAreMet() {
	roomName, _ := v.roomNameEntry.GetText()
	chatServiceName, _ := v.chatServiceEntry.GetText()

	hasAllValues := len(roomName) != 0 && len(chatServiceName) != 0 && v.ac.currentAccount() != nil
	v.joinButton.SetSensitive(hasAllValues)
}

// TODO: not sure what "before start" means in this method name.

func (v *mucJoinRoomView) beforeJoiningRoom() {
	v.clearErrors()
	v.disableJoinFields()
	v.showSpinner()
}

func (v *mucJoinRoomView) onJoinSuccess(a *account, roomID jid.Bare, roomInfo *muc.RoomListing) {
	doInUIThread(func() {
		v.hideSpinner()
		v.dialog.Hide()
		v.u.joinRoom(a, roomID, v.returnWhenCancelJoining)
	})
}

// TODO: This method name might also be nicer and more understandable
// This will be called when going back from the lobby, so maybe something about that?

func (v *mucJoinRoomView) returnWhenCancelJoining() {
	v.enableJoinFields()
	v.dialog.Show()
}

func (v *mucJoinRoomView) onJoinFails(a *account, roomID jid.Bare) {
	a.log.WithField("room", roomID).Warn("The room doesn't exist")
	doInUIThread(func() {
		v.notifyOnError(i18n.Local("The room doesn't exist on that service."))
		v.enableJoinFields()
		v.hideSpinner()
	})
}

func (v *mucJoinRoomView) onJoinError(a *account, roomID jid.Bare, err error) {
	doInUIThread(func() {
		v.hideSpinner()
		v.enableJoinFields()
		// TODO: This should not be necessary. We should analyze and check IF and why it could
		// happen that error is sent in as nil
		if err != nil {
			v.notifyOnError(i18n.Local("It looks like the room you are trying to connect to doesn't exist, please verify the provided information."))
			a.log.WithField("room", roomID).WithError(err).Warn("An error occurred trying to find the room")
		}
	})
}

func (v *mucJoinRoomView) onServiceUnavailable(a *account, roomID jid.Bare) {
	a.log.WithField("room", roomID).Warn("An error occurred trying to find the room")
	doInUIThread(func() {
		v.hideSpinner()
		v.notifyOnError(i18n.Local("We can't get access to the service, please check your Internet connection or make sure the service exists."))
	})
}

type mucJoinRoomContext struct {
	a      *account
	v      *mucJoinRoomView
	roomID jid.Bare
	done   func()
}

func (c *mucJoinRoomContext) onFinishWithError(err error, errorReceived bool) {
	if errorReceived {
		c.v.onJoinError(c.a, c.roomID, err)
		return
	}
	c.v.onServiceUnavailable(c.a, c.roomID)
}

func (c *mucJoinRoomContext) waitToFinish(result <-chan bool, errors <-chan error, roomInfo <-chan *muc.RoomListing) {
	defer doInUIThread(c.done)

	select {
	case value, ok := <-result:
		if !ok {
			c.v.onServiceUnavailable(c.a, c.roomID)
			return
		}

		if !value {
			c.v.onJoinFails(c.a, c.roomID)
			return
		}

		ri := <-roomInfo

		c.v.onJoinSuccess(c.a, c.roomID, ri)
	case err, ok := <-errors:
		c.onFinishWithError(err, !ok)
	}
}

func (c *mucJoinRoomContext) joinRoom() {
	c.v.beforeJoiningRoom()
	roomInfo := make(chan *muc.RoomListing)
	result, errors := c.a.session.HasRoom(c.roomID, roomInfo)
	go c.waitToFinish(result, errors, roomInfo)
}

func (v *mucJoinRoomView) log() coylog.Logger {
	l := v.u.log

	ca := v.ac.currentAccount()
	if ca != nil {
		l = ca.log
	}

	l.WithField("who", "mucJoinRoomView")

	return l
}

func (v *mucJoinRoomView) validateFieldsAndGetBareIfOk() (jid.Bare, bool) {
	roomName, _ := v.roomNameEntry.GetText()
	local := jid.NewLocal(roomName)
	if !local.Valid() {
		v.log().WithField("local", roomName).Error("Trying to join a room with an invalid local")
		v.notifyOnError(i18n.Local("You must provide a valid room name."))
		return nil, false
	}

	chatServiceName, _ := v.chatServiceEntry.GetText()
	domain := jid.NewDomain(chatServiceName)
	if !domain.Valid() {
		v.log().WithField("domain", chatServiceName).Error("Trying to join a room with an invalid domain")
		v.notifyOnError(i18n.Local("You must provide a valid service name."))
		return nil, false
	}

	return jid.NewBare(local, domain), true
}

func (v *mucJoinRoomView) tryJoinRoom(done func()) {
	roomID, ok := v.validateFieldsAndGetBareIfOk()
	if !ok {
		done()
		return
	}

	ca := v.ac.currentAccount()
	if ca == nil {
		v.notifyOnError(i18n.Local("No account was selected, select an account from the list or enable one."))
		return
	}

	c := &mucJoinRoomContext{
		a:      ca,
		v:      v,
		roomID: roomID,
		done:   done,
	}

	c.joinRoom()
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
	builder := newBuilder("MUCJoinRoomDialog")

	panicOnDevError(builder.bindObjects(v))

	v.errorNotif = newErrorNotification(v.notificationArea)

	accountsInput := builder.get("accounts").(gtki.ComboBox)
	v.ac = v.u.createConnectedAccountsComponent(accountsInput, v, func(a *account) {
		doInUIThread(v.enableJoinIfConditionsAreMet)
	}, func() {
		doInUIThread(v.enableJoinIfConditionsAreMet)
	})

	builder.ConnectSignals(map[string]interface{}{
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
	return jid.ValidBareJID(name)
}

func (v *mucJoinRoomView) setSensitivityForAllFields(f bool) {
	v.roomNameEntry.SetSensitive(f)
	v.chatServiceEntry.SetSensitive(f)
	v.joinButton.SetSensitive(f)
}

func (v *mucJoinRoomView) disableJoinFields() {
	v.setSensitivityForAllFields(false)
	v.ac.disableAccountInput()
}

func (v *mucJoinRoomView) enableJoinFields() {
	v.setSensitivityForAllFields(true)
	v.ac.enableAccountInput()
}

func (u *gtkUI) mucShowJoinRoom() {
	view := newMUCJoinRoomView(u)

	u.connectShortcutsChildWindow(view.dialog)

	view.dialog.SetTransientFor(u.window)
	view.dialog.Show()
}
