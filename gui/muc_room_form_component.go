package gui

import (
	"strings"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type mucRoomFormData struct {
	errorNotifications canNotifyErrors

	connectedAccountsInput gtki.ComboBox
	roomNameEntry          gtki.Entry
	chatServicesInput      gtki.ComboBoxText
	chatServicesEntry      gtki.Entry

	onAccountSelected    func(*account)
	onNoAccount          func()
	onChatServiceChanged func()
}

type mucRoomFormComponent struct {
	errorNotifications canNotifyErrors

	accountsComponent     *connectedAccountsComponent
	roomNameComponent     *roomNameComponent
	chatServicesComponent *chatServicesComponent
}

func (u *gtkUI) createMUCRoomFormComponent(d *mucRoomFormData) *mucRoomFormComponent {
	chatServicesComponent := u.createChatServicesComponent(d.chatServicesInput, d.chatServicesEntry, d.onChatServiceChanged)

	onAccountSelected := func(ca *account) {
		go chatServicesComponent.updateServicesBasedOnAccount(ca)
		go d.onAccountSelected(ca)
	}

	onNoAccount := func() {
		go chatServicesComponent.removeAll()
		go d.onNoAccount()
	}

	c := &mucRoomFormComponent{
		errorNotifications:    d.errorNotifications,
		accountsComponent:     u.createConnectedAccountsComponent(d.connectedAccountsInput, d.errorNotifications, onAccountSelected, onNoAccount),
		roomNameComponent:     u.createRoomNameComponent(d.roomNameEntry),
		chatServicesComponent: chatServicesComponent,
	}

	return c
}

func (f *mucRoomFormComponent) currentAccount() *account {
	return f.accountsComponent.currentAccount()
}

func (f *mucRoomFormComponent) setCurrentAccount(ca *account) {
	f.accountsComponent.setCurrentAccount(ca)
}

func (f *mucRoomFormComponent) currentService() jid.Domain {
	return f.chatServicesComponent.currentService()
}

func (f *mucRoomFormComponent) currentServiceValue() string {
	return f.chatServicesComponent.currentServiceValue()
}

func (f *mucRoomFormComponent) setCurrentServiceValue(s jid.Domain) {
	f.chatServicesComponent.setCurrentService(s)
}

func (f *mucRoomFormComponent) currentRoomName() jid.Local {
	return f.roomNameComponent.currentName()
}

func (f *mucRoomFormComponent) currentRoomNameValue() string {
	return f.roomNameComponent.currentNameValue()
}

func (f *mucRoomFormComponent) setCurrentRoomName(n jid.Local) {
	f.roomNameComponent.setCurrentNameValue(n.String())
}

func (f *mucRoomFormComponent) currentRoomID() jid.Bare {
	return jid.NewBareFromStrings(
		// Some services may fail with capital letters.
		// Maybe we should implement a "ToLower" method as a member of the "jid.Bare" struct
		strings.ToLower(f.currentRoomNameValue()),
		strings.ToLower(f.currentServiceValue()),
	)
}

func (f *mucRoomFormComponent) currentRoomIDValue() string {
	return f.currentRoomID().String()
}

func (f *mucRoomFormComponent) enableFields() {
	f.accountsComponent.enableAccountInput()
	f.roomNameComponent.enableNameInput()
	f.chatServicesComponent.enableServiceInput()
}

func (f *mucRoomFormComponent) disableFields() {
	f.accountsComponent.disableAccountInput()
	f.roomNameComponent.disableNameInput()
	f.chatServicesComponent.disableServiceInput()
}

func (f *mucRoomFormComponent) resetFields() {
	f.roomNameComponent.reset()
}

func (f *mucRoomFormComponent) onDestroy() {
	f.accountsComponent.onDestroy()
}

func (f *mucRoomFormComponent) notifyError(err string) {
	f.errorNotifications.notifyOnError(err)
}

func (f *mucRoomFormComponent) isFilled() bool {
	return f.accountsComponent.hasAccountValue() && f.roomNameComponent.hasNameValue() && f.chatServicesComponent.hasServiceValue()
}

func (f *mucRoomFormComponent) isEmpty() bool {
	return !f.isFilled()
}

func (f *mucRoomFormComponent) hasNoErrorsReported() bool {
	return f.accountsComponent.hasAccounts()
}

func (f *mucRoomFormComponent) isValid() bool {
	f.errorNotifications.clearErrors()

	local := f.currentRoomName()
	if !local.Valid() {
		f.notifyError(i18n.Local("You must provide a valid room name."))
		return false
	}

	domain := f.currentService()
	if !domain.Valid() {
		f.notifyError(i18n.Local("You must provide a valid service name."))
		return false
	}

	ca := f.currentAccount()
	if ca == nil {
		f.notifyError(i18n.Local("No account is selected, please select one account from the list or connect to one."))
		return false
	}

	return true
}
