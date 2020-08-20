package gui

import (
	"strings"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type createMUCRoom struct {
	u       *gtkUI
	builder *builder
	ac      *connectedAccountsComponent

	dialog           gtki.Dialog       `gtk-widget:"create-chat-dialog"`
	notificationArea gtki.Box          `gtk-widget:"notification-area"`
	account          gtki.ComboBox     `gtk-widget:"accounts"`
	chatServices     gtki.ComboBoxText `gtk-widget:"chatServices"`
	chatServiceEntry gtki.Entry        `gtk-widget:"chatServiceEntry"`
	room             gtki.Entry        `gtk-widget:"room"`
	cancelButton     gtki.Button       `gtk-widget:"button-cancel"`
	createButton     gtki.Button       `gtk-widget:"button-ok"`

	errorBox     *errorNotification
	notification gtki.InfoBar

	createButtonPrevText  string
	previousUpdateChannel chan bool
	supportedErrors       map[session.CreateRoomError]string
}

func (u *gtkUI) newCreateMUCRoom() *createMUCRoom {
	view := &createMUCRoom{
		u: u,
	}

	view.initSupportedMappingErrors()
	view.initUIBuilder()
	view.initConnectedAccountsComponent()

	return view
}

func (v *createMUCRoom) initSupportedMappingErrors() {
	m := map[session.CreateRoomError]string{}

	m[session.ErrInvalidInformationQueryRequest] = i18n.Local("Couldn't send the information query to the server, please try again.")
	m[session.ErrUnexpectedResponse] = i18n.Local("The connection to the server can't be established.")
	m[session.ErrInformationQueryResponse] = i18n.Local("You don't have the permissions to create a room on the server or the room is already created.")

	v.supportedErrors = m
}

func (v *createMUCRoom) initUIBuilder() {
	v.builder = newBuilder("MUCCreateRoom")

	panicOnDevError(v.builder.bindObjects(v))

	v.errorBox = newErrorNotification(v.notificationArea)

	v.builder.ConnectSignals(map[string]interface{}{
		"on_create_room":              v.onCreateRoom,
		"on_cancel":                   v.dialog.Destroy,
		"on_close_window":             v.onCloseWindow,
		"on_room_changed":             v.disableCreationIfAnyFieldIsEmpty,
		"on_chatServiceEntry_changed": v.disableCreationIfAnyFieldIsEmpty,
	})
}

func (v *createMUCRoom) initConnectedAccountsComponent() {
	c := v.builder.get("accounts").(gtki.ComboBox)
	v.ac = v.u.createConnectedAccountsComponent(c, v, v.updateServicesBasedOnAccount, v.onNoAccountsConnected)
}

func (v *createMUCRoom) onCloseWindow() {
	v.ac.onDestroy()
}

func (v *createMUCRoom) disableOrEnableFields(f bool) {
	v.cancelButton.SetSensitive(f)
	v.createButton.SetSensitive(f)
	v.account.SetSensitive(f)
	v.room.SetSensitive(f)
	v.chatServices.SetSensitive(f)
}

func (v *createMUCRoom) getRoomID() jid.Bare {
	roomName, err := v.room.GetText()
	if err != nil {
		v.u.log.WithError(err).Error("Something went wrong while trying to create the room")
		doInUIThread(func() {
			v.errorBox.ShowMessage(i18n.Local("Could not get the room name, please try again."))
		})
		return nil
	}

	service := v.chatServices.GetActiveText()
	if !jid.ValidLocal(strings.TrimSpace(roomName)) || !jid.ValidDomain(strings.TrimSpace(roomName)) {
		doInUIThread(func() {
			v.errorBox.ShowMessage(i18n.Localf("%s@%s is not a valid room name.", roomName, service))
		})
		return nil
	}

	return jid.NewBare(jid.NewLocal(strings.TrimSpace(roomName)), jid.NewDomain(strings.TrimSpace(service)))
}

func (v *createMUCRoom) onCreateRoom() {
	v.clearErrors()

	ca := v.ac.currentAccount()
	if ca == nil {
		v.errorBox.ShowMessage(i18n.Local("No account selected, please select one account from the list or connect to one."))
		return
	}

	roomIdentity := v.getRoomID()
	if roomIdentity != nil {
		v.onBeforeToCreateARoom()
		go v.createRoom(ca, roomIdentity)
	}
}

func (v *createMUCRoom) onBeforeToCreateARoom() {
	v.disableOrEnableFields(false)
	v.createButtonPrevText, _ = v.createButton.GetLabel()
	_ = v.createButton.SetProperty("label", i18n.Local("Creating room..."))
}

func (v *createMUCRoom) createRoom(ca *account, ident jid.Bare) {
	ec := ca.session.CreateRoom(ident)

	go func() {
		isRoomCreated := v.doSomethingWhenFinish(ca, ec)
		v.onCreateRoomFinished(isRoomCreated, ca, ident)
	}()
}

func (v *createMUCRoom) doSomethingWhenFinish(ca *account, ec <-chan error) bool {
	err, ok := <-ec
	if !ok {
		return true
	}

	if err != nil {
		ca.log.WithError(err).Error("Something went wrong while trying to create the room")

		userErr, ok := v.supportedErrors[err]
		if !ok {
			userErr = i18n.Local("Could not create the new room")
		}

		doInUIThread(func() {
			v.errorBox.ShowMessage(userErr)
		})
	}

	return false
}

func (v *createMUCRoom) onCreateRoomFinished(created bool, ca *account, ident jid.Bare) {
	doInUIThread(func() {
		v.disableOrEnableFields(true)
		_ = v.createButton.SetProperty("label", v.createButtonPrevText)
	})

	if created {
		doInUIThread(func() {
			v.u.mucShowRoom(ca, ident)
			v.dialog.Destroy()
		})
	}
}

func (v *createMUCRoom) disableCreationIfAnyFieldIsEmpty() {
	setEnabled(v.createButton, v.areAllFieldsFilled())
}

func (v *createMUCRoom) areAllFieldsFilled() bool {
	accountVal := ""
	if len(v.ac.accountsList) > 0 {
		ac := v.ac.currentAccount()
		if ac != nil {
			accountVal = ac.Account()
		}
	}

	serviceVal := v.chatServices.GetActiveText()
	roomVal, _ := v.room.GetText()

	return accountVal != "" && serviceVal != "" && roomVal != ""
}

func (v *createMUCRoom) notifyOnError(err string) {
	if v.notification != nil {
		v.notificationArea.Remove(v.notification)
	}
	v.errorBox.ShowMessage(err)
}

func (v *createMUCRoom) clearErrors() {
	v.errorBox.Hide()
}

func (v *createMUCRoom) onNoAccountsConnected() {
	doInUIThread(v.chatServices.RemoveAll)
}

func (v *createMUCRoom) updateServicesBasedOnAccount(acc *account) {
	doInUIThread(v.clearErrors)
	go v.updateChatServicesBasedOnAccount(acc)
}

func (v *createMUCRoom) updateChatServicesBasedOnAccount(ac *account) {
	if v.previousUpdateChannel != nil {
		v.previousUpdateChannel <- true
	}

	v.previousUpdateChannel = make(chan bool)

	csc, ec, endEarly := ac.session.GetChatServices(jid.ParseDomain(ac.Account()))

	go v.updateChatServices(ac, csc, ec, endEarly)
}

func (v *createMUCRoom) updateChatServices(ac *account, csc <-chan jid.Domain, ec <-chan error, endEarly func()) {
	hadAny := false

	var typedService string
	doInUIThread(func() {
		typedService, _ = v.chatServiceEntry.GetText()
	})

	doInUIThread(v.chatServices.RemoveAll)

	defer v.onUpdateChatServicesFinished(hadAny, typedService)

	for {
		select {
		case <-v.previousUpdateChannel:
			doInUIThread(v.chatServices.RemoveAll)
			endEarly()
			return
		case err, _ := <-ec:
			if err != nil {
				ac.log.WithError(err).Error("Something went wrong trying to get chat services")
			}
			return
		case cs, ok := <-csc:
			if !ok {
				return
			}

			hadAny = true
			doInUIThread(func() {
				v.chatServices.AppendText(cs.String())
			})
		}
	}
}

func (v *createMUCRoom) onUpdateChatServicesFinished(hadAny bool, typedService string) {
	if hadAny && typedService == "" {
		doInUIThread(func() {
			v.chatServices.SetActive(0)
		})
	}
	v.previousUpdateChannel = nil
}

func setEnabled(w gtki.Widget, enable bool) {
	w.SetSensitive(enable)
}

func (u *gtkUI) mucCreateChatRoom() {
	view := u.newCreateMUCRoom()
	view.dialog.SetTransientFor(u.window)
	view.dialog.Show()
}
