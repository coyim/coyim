package gui

import (
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/golang-collections/collections/set"
)

type mucCreateRoomViewForm struct {
	isShown           bool
	builder           *builder
	roomFormComponent *mucRoomFormComponent

	view              gtki.Box         `gtk-widget:"create-room-form"`
	roomAutoJoinCheck gtki.CheckButton `gtk-widget:"autojoin-check-button"`
	roomConfigCheck   gtki.CheckButton `gtk-widget:"config-room-check-button"`
	createButton      gtki.Button      `gtk-widget:"create-room-button"`
	spinnerBox        gtki.Box         `gtk-widget:"spinner-box"`
	notificationArea  gtki.Box         `gtk-widget:"notification-area-box"`

	spinner       *spinner
	notifications *notifications

	roomNameConflictList     *set.Set
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
	f.initNotifications(v)
	f.initRoomFormComponent(v)
	f.initDefaults(v)

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

func (f *mucCreateRoomViewForm) initNotifications(v *mucCreateRoomView) {
	f.notifications = v.u.newNotifications(f.notificationArea)
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

func (f *mucCreateRoomViewForm) initDefaults(v *mucCreateRoomView) {
	f.spinner = newSpinner()
	f.spinnerBox.Add(f.spinner.getWidget())
}

func (v *mucCreateRoomView) initCreateRoomForm() *mucCreateRoomViewForm {
	f := v.newCreateRoomForm()

	f.createRoom = func(ca *account, roomID jid.Bare) {
		errors := make(chan error)
		v.createRoom(ca, roomID, errors)
		go f.listenToCreateError(roomID, errors)
	}

	f.addCallbacks(v)

	return f
}

func (f *mucCreateRoomViewForm) onRoomAutoJoinToggled() {
	f.updateAutoJoinValue(f.roomAutoJoinCheck.GetActive())
}

func (f *mucCreateRoomViewForm) onRoomConfigToggled() {
	f.updateConfigureRoomValue(f.roomConfigCheck.GetActive())
}

func (f *mucCreateRoomViewForm) listenToCreateError(roomID jid.Bare, errors chan error) {
	err := <-errors

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

func (f *mucCreateRoomViewForm) onCreateRoomCheckIfExistsFails() {
	f.notifications.error(i18n.Local("Couldn't connect to the service, please verify that it exists or try again later."))
	f.spinner.hide()
	f.enableFields()
}

func (f *mucCreateRoomViewForm) onCreateRoomAlreadyExists() {
	f.notifications.error(i18n.Local("That room already exists, try again with a different name."))
	f.spinner.hide()
	f.enableFields()
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

func (f *mucCreateRoomViewForm) showCreateForm(v *mucCreateRoomView) {
	v.success.reset()
	v.container.Remove(v.success.view)
	f.reset()
	v.container.Add(f.view)
	f.isShown = true
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
