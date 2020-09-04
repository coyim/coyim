package gui

import (
	"github.com/coyim/coyim/i18n"
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
	chatServices     gtki.ComboBoxText `gtk-widget:"chat-services"`
	chatServiceEntry gtki.Entry        `gtk-widget:"chat-service-entry"`
	roomEntry        gtki.Entry        `gtk-widget:"room-name-entry"`
	roomAutoJoin     gtki.CheckButton  `gtk-widget:"room-auto-join"`
	spinner          gtki.Spinner      `gtk-widget:"spinner"`
	cancelButton     gtki.Button       `gtk-widget:"cancel-button"`
	createButton     gtki.Button       `gtk-widget:"create-button"`

	autoJoin     bool
	errorBox     *errorNotification
	notification gtki.InfoBar

	previousUpdateChannel chan bool
}

func (u *gtkUI) newCreateMUCRoom() *createMUCRoom {
	view := &createMUCRoom{
		u:        u,
		autoJoin: true,
	}

	view.initUIBuilder()
	view.initConnectedAccountsComponent()
	view.initDefaults()

	return view
}

func (v *createMUCRoom) initUIBuilder() {
	v.builder = newBuilder("MUCCreateRoomDialog")

	panicOnDevError(v.builder.bindObjects(v))

	v.errorBox = newErrorNotification(v.notificationArea)

	v.builder.ConnectSignals(map[string]interface{}{
		"on_create_room":             v.onCreateRoom,
		"on_cancel":                  v.dialog.Destroy,
		"on_close_window":            v.onCloseWindow,
		"on_roomName_change":         v.onRoomNameChange,
		"on_roomAutoJoin_toggled":    v.onRoomAutoJoinChange,
		"on_chatServiceEntry_change": v.onChatServiceChange,
	})
}

func (v *createMUCRoom) initConnectedAccountsComponent() {
	c := v.builder.get("accounts").(gtki.ComboBox)
	v.ac = v.u.createConnectedAccountsComponent(c, v, v.updateServicesBasedOnAccount, v.onNoAccountsConnected)
}

func (v *createMUCRoom) initDefaults() {
	v.roomAutoJoin.SetActive(v.autoJoin)
}

func (v *createMUCRoom) onCloseWindow() {
	v.ac.onDestroy()
}

// disableOrEnableFields SHOULD be called from the UI thread
func (v *createMUCRoom) disableOrEnableFields(f bool) {
	v.cancelButton.SetSensitive(f)
	v.createButton.SetSensitive(f)
	v.account.SetSensitive(f)
	v.roomEntry.SetSensitive(f)
	v.chatServices.SetSensitive(f)
	v.roomAutoJoin.SetSensitive(f)
	if f {
		v.ac.enableAccountInput()
	} else {
		v.ac.disableAccountInput()
	}
}

func (v *createMUCRoom) clearFields() {
	// TODO: should we keep the service name or not?
	v.roomEntry.SetText("")
	v.enableCreationIfConditionsAreMet()
}

func (v *createMUCRoom) onCreateRoom() {
	v.clearErrors()

	roomName, err := v.roomEntry.GetText()
	if err != nil {
		v.u.log.WithError(err).Error("Something went wrong while trying to create the room")
		v.notifyOnError(i18n.Local("Could not get the room name, please try again."))
	}

	local := jid.NewLocal(roomName)
	if !local.Valid() {
		v.u.log.WithField("local", roomName).Error("Trying to create a room with an invalid local")
		v.notifyOnError(i18n.Local("You must provide a valid room name."))
		return
	}

	chatService, err := v.chatServiceEntry.GetText()
	if err != nil {
		v.u.log.WithError(err).Error("Something went wrong while trying to create the room")
		v.notifyOnError(i18n.Local("Could not get the service name, please try again."))
		return
	}

	domain := jid.NewDomain(chatService)
	if !domain.Valid() {
		v.u.log.WithField("domain", chatService).Error("Trying to create a room with an invalid domain")
		v.notifyOnError(i18n.Local("You must provide a valid service name."))
		return
	}

	roomIdentity := jid.NewBare(local, domain)

	ca := v.ac.currentAccount()
	if ca == nil {
		v.u.log.WithField("room", roomIdentity).Error("No account was selected to create the room")
		v.notifyOnError(i18n.Local("No account is selected, please select one account from the list or connect to one."))
		return
	}

	v.onBeforeToCreateARoom()

	go v.createRoomIfDoesntExist(ca, roomIdentity)
}

func (v *createMUCRoom) onBeforeToCreateARoom() {
	v.showSpinner()
	v.disableOrEnableFields(false)
}

func (v *createMUCRoom) createRoomIfDoesntExist(ca *account, ident jid.Bare) {
	erc, ec := ca.session.HasRoom(ident)
	go func() {
		select {
		case err, _ := <-ec:
			if err != nil {
				ca.log.WithError(err).Error("Error trying to validate if room exists")
				doInUIThread(func() {
					v.errorBox.ShowMessage(i18n.Local("Couldn't connect to the service, please verify that it exists or try again later."))
					v.hideSpinner()
					v.disableOrEnableFields(true)
				})
			}

		case er, _ := <-erc:
			if !er {
				ec := ca.session.CreateRoom(ident)
				go func() {
					isRoomCreated := v.listenToRoomCreation(ca, ec)
					v.onCreateRoomFinished(isRoomCreated, ca, ident)
				}()
				return
			}

			doInUIThread(func() {
				v.errorBox.ShowMessage(i18n.Local("That room already exists, try again with a different name."))
				v.hideSpinner()
				v.disableOrEnableFields(true)
			})
		}
	}()
}

func (v *createMUCRoom) listenToRoomCreation(ca *account, ec <-chan error) bool {
	err, ok := <-ec
	if !ok {
		return true
	}

	if err != nil {
		ca.log.WithError(err).Error("Something went wrong while trying to create the room")

		userErr, ok := supportedCreateMUCErrors[err]
		if !ok {
			userErr = i18n.Local("Could not create the new room.")
		}

		doInUIThread(func() {
			v.errorBox.ShowMessage(userErr)
		})
	}

	return false
}

func (v *createMUCRoom) onCreateRoomFinished(created bool, ca *account, ident jid.Bare) {
	if created {
		if !v.autoJoin {
			doInUIThread(func() {
				v.notifyOnError(i18n.Local("The room has been created."))
				v.disableOrEnableFields(true)
				v.hideSpinner()
				v.clearFields()
			})
			return
		}

		doInUIThread(func() {
			v.u.mucShowRoom(ca, ident)
			v.dialog.Destroy()
		})
	}
}

func (v *createMUCRoom) onRoomNameChange() {
	v.enableCreationIfConditionsAreMet()
}

func (v *createMUCRoom) onRoomAutoJoinChange() {
	v.autoJoin = v.roomAutoJoin.GetActive()
	if v.autoJoin {
		v.createButton.SetProperty("label", i18n.Local("Create Room & Join"))
	} else {
		v.createButton.SetProperty("label", i18n.Local("Create Room"))
	}
}

func (v *createMUCRoom) onChatServiceChange() {
	v.enableCreationIfConditionsAreMet()
}

func (v *createMUCRoom) enableCreationIfConditionsAreMet() {
	roomName, _ := v.roomEntry.GetText()
	chatService, _ := v.chatServiceEntry.GetText()
	currentAccount := v.ac.currentAccount()

	hasAllValues := len(roomName) != 0 && len(chatService) != 0 && currentAccount != nil
	v.createButton.SetSensitive(hasAllValues)
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
	doInUIThread(func() {
		v.enableCreationIfConditionsAreMet()
		v.chatServices.RemoveAll()
	})
}

func (v *createMUCRoom) updateServicesBasedOnAccount(acc *account) {
	doInUIThread(v.enableCreationIfConditionsAreMet)
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
		v.chatServices.RemoveAll()
		v.showSpinner()
	})

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

	doInUIThread(v.hideSpinner)

	v.previousUpdateChannel = nil
}

func (v *createMUCRoom) showSpinner() {
	v.spinner.Start()
	v.spinner.Show()
}

func (v *createMUCRoom) hideSpinner() {
	v.spinner.Stop()
	v.spinner.Hide()
}

func setEnabled(w gtki.Widget, enable bool) {
	w.SetSensitive(enable)
}

func (u *gtkUI) mucCreateChatRoom() {
	view := u.newCreateMUCRoom()

	u.connectShortcutsChildWindow(view.dialog)

	view.dialog.SetTransientFor(u.window)
	view.dialog.Show()
}
