package gui

import (
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type mucJoinRoomView struct {
	u                 *gtkUI
	builder           *builder
	roomFormComponent *mucRoomFormComponent

	dialog           gtki.Dialog `gtk-widget:"join-room-dialog"`
	joinButton       gtki.Button `gtk-widget:"join-room-button"`
	spinnerBox       gtki.Box    `gtk-widget:"spinner-box"`
	notificationArea gtki.Box    `gtk-widget:"notification-area-box"`

	spinner       *spinner
	notifications *notifications
}

func newMUCJoinRoomView(u *gtkUI) *mucJoinRoomView {
	view := &mucJoinRoomView{
		u: u,
	}

	view.initBuilder()
	view.initNotifications()
	view.initRoomFormComponent()
	view.initDefaults()

	u.connectShortcutsChildWindow(view.dialog)

	return view
}

func (v *mucJoinRoomView) initBuilder() {
	v.builder = newBuilder("MUCJoinRoomDialog")
	panicOnDevError(v.builder.bindObjects(v))

	v.builder.ConnectSignals(map[string]interface{}{
		"on_close_window":     v.onCloseWindow,
		"on_roomname_changed": v.enableJoinIfConditionsAreMet,
		"on_cancel_clicked":   v.dialog.Destroy,
		"on_join_clicked":     doOnlyOnceAtATime(v.tryJoinRoom),
	})
}

func (v *mucJoinRoomView) initNotifications() {
	v.notifications = v.u.newNotifications(v.notificationArea)
}

func (v *mucJoinRoomView) initRoomFormComponent() {
	account := v.builder.get("accounts").(gtki.ComboBox)
	roomEntry := v.builder.get("room-name-entry").(gtki.Entry)
	chatServicesList := v.builder.get("chat-services-list").(gtki.ComboBoxText)
	chatServicesEntry := v.builder.get("chat-services-entry").(gtki.Entry)

	v.roomFormComponent = v.u.createMUCRoomFormComponent(&mucRoomFormData{
		errorNotifications:     v.notifications,
		connectedAccountsInput: account,
		roomNameEntry:          roomEntry,
		chatServicesInput:      chatServicesList,
		chatServicesEntry:      chatServicesEntry,
		onAccountSelected:      v.updateServicesBasedOnAccount,
		onNoAccount:            v.onNoAccountsConnected,
		onChatServiceChanged:   v.enableJoinIfConditionsAreMet,
	})
}

func (v *mucJoinRoomView) updateServicesBasedOnAccount(ca *account) {
	doInUIThread(func() {
		v.notifications.clearErrors()
		v.enableJoinIfConditionsAreMet()
	})
}

func (v *mucJoinRoomView) onNoAccountsConnected() {
	doInUIThread(v.enableJoinIfConditionsAreMet)
}

func (v *mucJoinRoomView) initDefaults() {
	v.spinner = newSpinner()
	v.spinnerBox.Add(v.spinner.getWidget())
}

func (v *mucJoinRoomView) onCloseWindow() {
	v.roomFormComponent.onDestroy()
}

func (v *mucJoinRoomView) typedRoomName() string {
	return v.roomFormComponent.currentRoomNameValue()
}

// enableJoinIfConditionsAreMet MUST be called from the UI thread
func (v *mucJoinRoomView) enableJoinIfConditionsAreMet() {
	v.joinButton.SetSensitive(v.roomFormComponent.isFilled())
}

func (v *mucJoinRoomView) beforeJoiningRoom() {
	v.notifications.clearErrors()
	v.disableJoinFields()
	v.spinner.show()
}

func (v *mucJoinRoomView) onJoinSuccess(a *account, roomID jid.Bare, roomInfo *muc.RoomListing) {
	doInUIThread(func() {
		v.spinner.hide()
		v.dialog.Hide()
		v.u.joinRoom(a, roomID, v.returnToJoinRoomView)
	})
}

func (v *mucJoinRoomView) returnToJoinRoomView() {
	v.enableJoinFields()
	v.dialog.Show()
}

func (v *mucJoinRoomView) onJoinFails(a *account, roomID jid.Bare) {
	a.log.WithField("room", roomID).Warn("The room doesn't exist")

	doInUIThread(func() {
		v.notifications.error(i18n.Local("The room doesn't exist on that service."))
		v.enableJoinFields()
		v.spinner.hide()
	})
}

func (v *mucJoinRoomView) onJoinError(a *account, roomID jid.Bare, err error) {
	a.log.WithField("room", roomID).WithError(err).Warn("An error occurred trying to find the room")

	doInUIThread(func() {
		v.notifications.error(i18n.Local("It looks like the room you are trying to connect to doesn't exist, please verify the provided information."))
		v.enableJoinFields()
		v.spinner.hide()
	})
}

func (v *mucJoinRoomView) onServiceUnavailable(a *account, roomID jid.Bare) {
	a.log.WithField("room", roomID).Warn("An error occurred trying to find the room")

	doInUIThread(func() {
		v.notifications.error(i18n.Local("We can't get access to the service, please check your Internet connection or make sure the service exists."))
		v.enableJoinFields()
		v.spinner.hide()
	})
}

func (v *mucJoinRoomView) log() coylog.Logger {
	l := v.u.log

	ca := v.roomFormComponent.currentAccount()
	if ca != nil {
		l = ca.log
	}

	l.WithField("who", "mucJoinRoomView")

	return l
}

func (v *mucJoinRoomView) validateFieldsAndGetBareIfOk() (jid.Bare, bool) {
	local := v.roomFormComponent.currentRoomName()
	if !local.Valid() {
		v.notifications.error(i18n.Local("You must provide a valid room name."))
		return nil, false
	}

	chatServiceName := v.roomFormComponent.currentService()
	if !chatServiceName.Valid() {
		v.notifications.error(i18n.Local("You must provide a valid service name."))
		return nil, false
	}

	return jid.NewBare(local, chatServiceName), true
}

func (v *mucJoinRoomView) tryJoinRoom(done func()) {
	roomID, ok := v.validateFieldsAndGetBareIfOk()
	if !ok {
		done()
		return
	}

	ca := v.roomFormComponent.currentAccount()
	if ca == nil {
		v.notifications.error(i18n.Local("No account was selected, select an account from the list or enable one."))
		return
	}

	c := v.newJoinRoomContext(ca, roomID, done)

	c.joinRoom()
}

func (v *mucJoinRoomView) isValidRoomName(name string) bool {
	return jid.ValidBareJID(name)
}

func (v *mucJoinRoomView) setSensitivityForJoin(f bool) {
	v.joinButton.SetSensitive(f)
}

func (v *mucJoinRoomView) disableJoinFields() {
	v.setSensitivityForJoin(false)
	v.roomFormComponent.disableFields()
}

func (v *mucJoinRoomView) enableJoinFields() {
	v.setSensitivityForJoin(true)
	v.roomFormComponent.enableFields()
}

func (u *gtkUI) mucShowJoinRoom() {
	view := newMUCJoinRoomView(u)

	u.connectShortcutsChildWindow(view.dialog)

	view.dialog.SetTransientFor(u.window)
	view.dialog.Show()
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
