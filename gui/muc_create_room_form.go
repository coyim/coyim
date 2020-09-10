package gui

import (
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/golang-collections/collections/set"
)

type createMUCRoomForm struct {
	log coylog.Logger
	ac  *connectedAccountsComponent

	view             gtki.Box          `gtk-widget:"createRoomForm"`
	notificationArea gtki.Box          `gtk-widget:"notificationArea"`
	account          gtki.ComboBox     `gtk-widget:"accounts"`
	chatServices     gtki.ComboBoxText `gtk-widget:"chatServices"`
	chatServiceEntry gtki.Entry        `gtk-widget:"chatServicesEntry"`
	roomEntry        gtki.Entry        `gtk-widget:"roomNameEntry"`
	roomAutoJoin     gtki.CheckButton  `gtk-widget:"autoJoinCheckButton"`
	spinner          gtki.Spinner      `gtk-widget:"createRoomFormSpinner"`
	cancelButton     gtki.Button       `gtk-widget:"createRoomFormCancelButton"`
	createButton     gtki.Button       `gtk-widget:"createRoomFormCreateButton"`

	errorBox     *errorNotification
	notification gtki.InfoBar

	previousUpdateChannel   chan bool
	roomNameConflictList    *set.Set
	createRoomIfDoesntExist func(*account, jid.Bare)
	onCheckFieldsConditions func(string, string, *account) bool
}

func (v *createMUCRoom) initForm() {
	f := &createMUCRoomForm{
		log:                  v.log,
		roomNameConflictList: set.New(),
	}

	panicOnDevError(v.builder.bindObjects(f))

	f.errorBox = newErrorNotification(f.notificationArea)

	c := v.builder.get("accounts").(gtki.ComboBox)
	f.ac = v.u.createConnectedAccountsComponent(c, f, f.updateServicesBasedOnAccount, f.onNoAccountsConnected)

	f.createRoomIfDoesntExist = func(a *account, ident jid.Bare) {
		successResult := make(chan bool)
		errResult := make(chan error)
		v.createRoomIfDoesntExist(a, ident, successResult, errResult)
		select {
		case err := <-errResult:
			switch err {
			case errCreateRoomCheckIfExistsFails:
				doInUIThread(func() {
					f.errorBox.ShowMessage(i18n.Local("Couldn't connect to the service, please verify that it exists or try again later."))
					f.hideSpinner()
					f.disableOrEnableFields(true)
				})

			case errCreateRoomAlreadyExists:
				f.roomNameConflictList.Insert(ident.String())
				doInUIThread(func() {
					f.errorBox.ShowMessage(i18n.Local("That room already exists, try again with a different name."))
					f.hideSpinner()
					f.disableOrEnableFields(true)
					f.createButton.SetSensitive(false)
				})

			case errCreateRoomFailed:
				userErr, ok := supportedCreateMUCErrors[err]
				if !ok {
					userErr = i18n.Local("Could not create the room.")
				}
				doInUIThread(func() {
					f.errorBox.ShowMessage(userErr)
				})
			}
		}
	}

	v.addBuilderSignals(map[string]interface{}{
		"on_cancel":          v.onCancel,
		"on_create_room":     f.onCreateRoom,
		"on_roomName_change": f.enableCreationIfConditionsAreMet,
		"on_roomAutoJoin_toggled": func() {
			v.updateAutoJoinValue(f.roomAutoJoin.GetActive())
		},
		"on_chatServiceEntry_change": f.enableCreationIfConditionsAreMet,
	})

	v.onAutoJoin(f.onAutoJoinChange)
	v.onDestroy(f.destroy)

	v.showCreateForm = func() {
		v.success.reset()
		v.container.Remove(v.success.view)
		v.form.reset()
		v.container.Add(v.form.view)
	}

	v.form = f
}

func (f *createMUCRoomForm) onAutoJoinChange(aj bool) {
	if aj {
		f.createButton.SetProperty("label", i18n.Local("Create Room & Join"))
	} else {
		f.createButton.SetProperty("label", i18n.Local("Create Room"))
	}
}

func (f *createMUCRoomForm) onCreateRoom() {
	f.clearErrors()

	roomName, _ := f.roomEntry.GetText()
	local := jid.NewLocal(roomName)
	if !local.Valid() {
		f.log.WithField("local", roomName).Error("Trying to create a room with an invalid local")
		f.notifyOnError(i18n.Local("You must provide a valid room name."))
		return
	}

	chatService, _ := f.chatServiceEntry.GetText()
	domain := jid.NewDomain(chatService)
	if !domain.Valid() {
		f.log.WithField("domain", chatService).Error("Trying to create a room with an invalid domain")
		f.notifyOnError(i18n.Local("You must provide a valid service name."))
		return
	}

	roomIdentity := jid.NewBare(local, domain)

	ca := f.ac.currentAccount()
	if ca == nil {
		f.log.WithField("room", roomIdentity).Error("No account was selected to create the room")
		f.notifyOnError(i18n.Local("No account is selected, please select one account from the list or connect to one."))
		return
	}

	f.onBeforeToCreateARoom()

	go f.createRoomIfDoesntExist(ca, roomIdentity)
}

func (f *createMUCRoomForm) onBeforeToCreateARoom() {
	f.showSpinner()
	f.disableOrEnableFields(false)
}

func (f *createMUCRoomForm) destroy() {
	f.ac.onDestroy()
}

func (f *createMUCRoomForm) notifyOnError(err string) {
	if f.notification != nil {
		f.notificationArea.Remove(f.notification)
	}
	f.errorBox.ShowMessage(err)
}

func (f *createMUCRoomForm) clearErrors() {
	f.errorBox.Hide()
}

func (f *createMUCRoomForm) clearFields() {
	// TODO: should we keep the service name?
	f.roomEntry.SetText("")
	f.enableCreationIfConditionsAreMet()
}

func (f *createMUCRoomForm) reset() {
	f.spinner.Stop()
	f.disableOrEnableFields(true)
	f.clearFields()
}

// disableOrEnableFields SHOULD be called from the UI thread
func (f *createMUCRoomForm) disableOrEnableFields(v bool) {
	f.cancelButton.SetSensitive(v)
	f.createButton.SetSensitive(v)
	f.account.SetSensitive(v)
	f.roomEntry.SetSensitive(v)
	f.chatServices.SetSensitive(v)
	f.roomAutoJoin.SetSensitive(v)
	if v {
		f.ac.enableAccountInput()
	} else {
		f.ac.disableAccountInput()
	}
}

func (f *createMUCRoomForm) updateServicesBasedOnAccount(acc *account) {
	doInUIThread(f.enableCreationIfConditionsAreMet)
	go f.updateChatServicesBasedOnAccount(acc)
}

func (f *createMUCRoomForm) onNoAccountsConnected() {
	doInUIThread(func() {
		f.enableCreationIfConditionsAreMet()
		f.chatServices.RemoveAll()
	})
}

func (f *createMUCRoomForm) enableCreationIfConditionsAreMet() {
	f.clearErrors()

	roomName, _ := f.roomEntry.GetText()
	chatService, _ := f.chatServiceEntry.GetText()
	currentAccount := f.ac.currentAccount()

	s := len(roomName) != 0 && len(chatService) != 0 && currentAccount != nil
	if s {
		ident := jid.NewBare(jid.NewLocal(roomName), jid.NewDomain(chatService))
		if ident.Valid() {
			if f.roomNameConflictList.Has(ident.String()) {
				f.errorBox.ShowMessage(i18n.Local("That room already exists, try again with a different name."))
				s = false
			}
		}
	}

	f.createButton.SetSensitive(s)
}

func (f *createMUCRoomForm) updateChatServicesBasedOnAccount(ac *account) {
	if f.previousUpdateChannel != nil {
		f.previousUpdateChannel <- true
	}

	f.previousUpdateChannel = make(chan bool)

	csc, ec, endEarly := ac.session.GetChatServices(jid.ParseDomain(ac.Account()))

	go f.updateChatServices(ac, csc, ec, endEarly)
}

func (f *createMUCRoomForm) updateChatServices(ac *account, csc <-chan jid.Domain, ec <-chan error, endEarly func()) {
	hadAny := false

	var typedService string
	doInUIThread(func() {
		typedService, _ = f.chatServiceEntry.GetText()
		f.chatServices.RemoveAll()
		f.showSpinner()
	})

	defer f.onUpdateChatServicesFinished(hadAny, typedService)

	for {
		select {
		case <-f.previousUpdateChannel:
			doInUIThread(f.chatServices.RemoveAll)
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
				f.chatServices.AppendText(cs.String())
			})
		}
	}
}

func (f *createMUCRoomForm) onUpdateChatServicesFinished(hadAny bool, typedService string) {
	if hadAny && typedService == "" {
		doInUIThread(func() {
			f.chatServices.SetActive(0)
		})
	}

	doInUIThread(f.hideSpinner)

	f.previousUpdateChannel = nil
}

func (f *createMUCRoomForm) showSpinner() {
	f.spinner.Start()
	f.spinner.Show()
}

func (f *createMUCRoomForm) hideSpinner() {
	f.spinner.Stop()
	f.spinner.Hide()
}

func setEnabled(w gtki.Widget, enable bool) {
	w.SetSensitive(enable)
}
