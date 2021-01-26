package gui

import (
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/golang-collections/collections/set"
)

// initCreateRoomForm MUST be called from the UI thread
func (v *mucCreateRoomView) initCreateRoomForm(d *mucCreateRoomData) {
	f := v.newCreateRoomForm()

	if d != nil {
		f.roomFormComponent.setCurrentAccount(d.ca)
		f.roomFormComponent.setCurrentRoomName(d.roomName)
		f.roomFormComponent.setCurrentServiceValue(d.where)
		f.roomAutoJoinCheck.SetActive(d.autoJoin)
		f.roomConfigCheck.SetActive(d.customConfig)
	}

	f.createRoom = func(ca *account, roomID jid.Bare) {
		v.createRoom(ca, roomID, func(err error) {
			f.onCreateRoomError(roomID, err)
		})
	}

	f.addCallbacks(v)

	v.form = f
}

func (v *mucCreateRoomView) showCreateForm() {
	v.success.reset()
	v.container.Remove(v.success.view)
	v.container.Add(v.form.view)
	v.form.isShown = true
}

type mucCreateRoomViewForm struct {
	isShown           bool
	builder           *builder
	roomFormComponent *mucRoomFormComponent

	view              gtki.Box         `gtk-widget:"create-room-form"`
	roomAutoJoinCheck gtki.CheckButton `gtk-widget:"autojoin-check-button"`
	roomConfigCheck   gtki.CheckButton `gtk-widget:"config-room-check-button"`
	createButton      gtki.Button      `gtk-widget:"create-room-button"`
	spinnerBox        gtki.Box         `gtk-widget:"spinner-box"`
	notificationsArea gtki.Box         `gtk-widget:"notification-area-box"`

	spinner       *spinner
	notifications *notifications

	roomNameConflictList *set.Set
	// createRoom MUST NOT be called from the UI thread
	createRoom               func(*account, jid.Bare)
	updateAutoJoinValue      func(bool)
	updateConfigureRoomValue func(bool)
	onCheckFieldsConditions  func(string, string, *account) bool

	log func(*account, jid.Bare) coylog.Logger
}

func (v *mucCreateRoomView) newCreateRoomForm() *mucCreateRoomViewForm {
	f := &mucCreateRoomViewForm{
		roomNameConflictList:     set.New(),
		updateAutoJoinValue:      v.updateAutoJoinValue,
		updateConfigureRoomValue: v.updateConfigureRoomValue,
		log:                      v.log,
	}

	f.initBuilder(v)
	f.initNotificationsAndSpinner(v)
	f.initRoomFormComponent(v)

	return f
}

func (f *mucCreateRoomViewForm) initBuilder(v *mucCreateRoomView) {
	f.builder = newBuilder("MUCCreateRoomForm")
	panicOnDevError(f.builder.bindObjects(f))

	f.builder.ConnectSignals(map[string]interface{}{
		"on_cancel":                   v.onCancel,
		"on_create":                   f.onCreateRoom,
		"on_room_name_change":         f.enableCreationIfConditionsAreMet,
		"on_room_autojoin_toggled":    f.onRoomAutoJoinToggled,
		"on_room_config_toggled":      f.onRoomConfigToggled,
		"on_chatservice_entry_change": f.enableCreationIfConditionsAreMet,
	})
}

func (f *mucCreateRoomViewForm) initRoomFormComponent(v *mucCreateRoomView) {
	account := f.builder.get("accounts").(gtki.ComboBox)
	roomEntry := f.builder.get("room-name-entry").(gtki.Entry)
	chatServicesList := f.builder.get("chat-services-list").(gtki.ComboBoxText)
	chatServicesEntry := f.builder.get("chat-services-entry").(gtki.Entry)

	f.roomFormComponent = v.u.createMUCRoomFormComponent(&mucRoomFormData{
		errorNotifications:     f.notifications,
		connectedAccountsInput: account,
		roomNameEntry:          roomEntry,
		chatServicesInput:      chatServicesList,
		chatServicesEntry:      chatServicesEntry,
		onAccountSelected:      f.updateServicesBasedOnAccount,
		onNoAccount:            f.onNoAccountsConnected,
		onChatServiceChanged:   f.enableCreationIfConditionsAreMet,
	})
}

func (f *mucCreateRoomViewForm) initNotificationsAndSpinner(v *mucCreateRoomView) {
	f.spinner = v.u.newSpinnerComponent()
	f.notifications = v.u.newNotificationsComponent()

	f.spinnerBox.Add(f.spinner.widget())
	f.notificationsArea.Add(f.notifications.getBox())
}

func (f *mucCreateRoomViewForm) onRoomAutoJoinToggled() {
	f.updateAutoJoinValue(f.roomAutoJoinCheck.GetActive())
}

func (f *mucCreateRoomViewForm) onRoomConfigToggled() {
	f.updateConfigureRoomValue(f.roomConfigCheck.GetActive())
}

func (f *mucCreateRoomViewForm) onCreateRoomError(roomID jid.Bare, err error) {
	doInUIThread(f.hideSpinnerAndEnableFields)

	switch err {
	case errCreateRoomCheckIfExistsFails:
		doInUIThread(f.onCreateRoomCheckIfExistsFails)

	case errCreateRoomAlreadyExists:
		f.roomNameConflictList.Insert(roomID.String())
		doInUIThread(f.onCreateRoomAlreadyExists)

	case errCreateRoomFailed:
		doInUIThread(func() {
			f.onCreateRoomFailed(err)
		})
	}
}

// hideSpinnerAndEnableFields MUST be called from the UI thread
func (f *mucCreateRoomViewForm) hideSpinnerAndEnableFields() {
	f.spinner.hide()
	f.enableFields()
}

// onCreateRoomCheckIfExistsFails MUST be called from the UI thread
func (f *mucCreateRoomViewForm) onCreateRoomCheckIfExistsFails() {
	f.notifications.error(i18n.Local("Couldn't connect to the service, please verify that it exists or try again later."))
	f.hideSpinnerAndEnableFields()
}

// onCreateRoomAlreadyExists MUST be called from the UI thread
func (f *mucCreateRoomViewForm) onCreateRoomAlreadyExists() {
	f.notifications.error(i18n.Local("That room already exists, try again with a different name."))
	f.hideSpinnerAndEnableFields()
	f.disableCreateRoomButton()
}

func (f *mucCreateRoomViewForm) onCreateRoomFailed(err error) {
	displayErr, ok := supportedCreateMUCErrors[err]
	if ok {
		f.notifications.error(displayErr)
	} else {
		f.notifications.error(i18n.Local("Could not create the room."))
	}
}

func (f *mucCreateRoomViewForm) addCallbacks(v *mucCreateRoomView) {
	v.onCreateOptionChange.add(func() {
		f.onCreateOptionsChange(v.autoJoin, v.configureRoom)
	})

	v.onDestroy.add(f.destroy)
}

func (f *mucCreateRoomViewForm) setCreateRoomButtonLabel(l string) {
	f.createButton.SetProperty("label", l)
}

func (f *mucCreateRoomViewForm) onCreateOptionsChange(autoJoin, configRoom bool) {
	f.setCreateRoomButtonLabel(labelForCreateRoomButton(autoJoin, configRoom))
}

func (f *mucCreateRoomViewForm) onCreateRoom() {
	if f.roomFormComponent.isValid() {
		f.beforeCreatingTheRoom()
		go f.createRoom(f.roomFormComponent.currentAccount(), f.roomFormComponent.currentRoomID())
	}
}

func (f *mucCreateRoomViewForm) beforeCreatingTheRoom() {
	f.spinner.show()
	f.disableFields()
	f.setFieldsSensitive(false)
}

func (f *mucCreateRoomViewForm) destroy() {
	f.isShown = false
	f.roomFormComponent.onDestroy()
}

func (f *mucCreateRoomViewForm) clearFields() {
	f.roomFormComponent.resetFields()
	f.enableCreationIfConditionsAreMet()
}

func (f *mucCreateRoomViewForm) reset() {
	f.spinner.hide()
	f.enableFields()
	f.clearFields()
}

func (f *mucCreateRoomViewForm) enableFields() {
	f.roomFormComponent.enableFields()
	f.setFieldsSensitive(true)
}

func (f *mucCreateRoomViewForm) disableFields() {
	f.roomFormComponent.disableFields()
	f.setFieldsSensitive(false)
}

func (f *mucCreateRoomViewForm) setFieldsSensitive(v bool) {
	f.createButton.SetSensitive(v)
	f.roomAutoJoinCheck.SetSensitive(v)
	f.roomConfigCheck.SetSensitive(v)
}

func (f *mucCreateRoomViewForm) updateServicesBasedOnAccount(ca *account) {
	doInUIThread(func() {
		f.notifications.clearErrors()
		f.enableCreationIfConditionsAreMet()
	})
}

func (f *mucCreateRoomViewForm) onNoAccountsConnected() {
	doInUIThread(f.enableCreationIfConditionsAreMet)
}

func (f *mucCreateRoomViewForm) enableCreationIfConditionsAreMet() {
	if f.roomFormComponent.hasNoErrorsReported() {
		f.notifications.clearErrors()
	}

	f.disableCreateRoomButton()

	if f.roomFormComponent.isEmpty() || f.checkIfRoomNameHasConflict() {
		return
	}

	f.enableCreateRoomButton()
}

func (f *mucCreateRoomViewForm) checkIfRoomNameHasConflict() bool {
	if f.roomNameConflictList.Has(f.roomFormComponent.currentRoomIDValue()) {
		f.notifications.error(i18n.Local("That room already exists, try again with a different name."))
		return true
	}
	return false
}

func (f *mucCreateRoomViewForm) enableCreateRoomButton() {
	f.createButton.SetSensitive(true)
}

func (f *mucCreateRoomViewForm) disableCreateRoomButton() {
	f.createButton.SetSensitive(false)
}

func setEnabled(w gtki.Widget, enable bool) {
	w.SetSensitive(enable)
}

func labelForCreateRoomButton(autoJoin, configRoom bool) string {
	if configRoom {
		return i18n.Local("Configure Room")
	}

	if autoJoin {
		return i18n.Local("Create Room & Join")
	}

	return i18n.Local("Create Room")
}
