package gui

import (
	"strings"

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
	chatServices     gtki.ComboBoxText `gtk-widget:"chatServices"`
	chatServiceEntry gtki.Entry        `gtk-widget:"chatServiceEntry"`
	room             gtki.Entry        `gtk-widget:"room"`
	cancelButton     gtki.Button       `gtk-widget:"button-cancel"`
	createButton     gtki.Button       `gtk-widget:"button-ok"`

	errorBox     *errorNotification
	notification gtki.InfoBar

	createButtonPrevText  string
	previousUpdateChannel chan bool

	fieldsToValidate map[string]*validateField
}

const (
	validationFieldAccount     = "account"
	validationFieldRoomName    = "room"
	validationFieldChatService = "service"

	invalidCharactersForRoomName    = "\"&'/:<>@ "
	invalidCharactersForChatService = "\"&'/:<>@+ "
)

type rule func() error

type validateField struct {
	widget    gtki.Widget
	rules     []rule
	isValid   bool
	lastError error
}

type errorEmptyField struct {
	fieldName string
}

type errorNotAllowedCharacters struct {
	characters string
	fieldName  string
}

type errorNotValidLocal struct {
	local string
}

type errorNotValidDomain struct {
	domain string
}

type errorNoConnectedAccount struct{}

func newErrorNotAllowedCharacters(c, n string) error {
	return &errorNotAllowedCharacters{
		characters: c,
		fieldName:  n,
	}
}

func newErrorEmptyField(n string) error {
	return &errorEmptyField{fieldName: n}
}

func newErrorValidLocal(l string) error {
	return &errorNotValidLocal{local: l}
}

func newErrorValidDomain(d string) error {
	return &errorNotValidDomain{domain: d}
}

func newErrorNoConnectedAccount() error {
	return &errorNoConnectedAccount{}
}

func (e *errorNotAllowedCharacters) Error() string {
	switch e.fieldName {
	case validationFieldRoomName:
		return i18n.Localf("The character(s) [%s] are not allowed in room name field", e.characters)
	case validationFieldChatService:
		return i18n.Localf("The character(s) [%s] are not allowed in chat service", e.characters)
	default:
		return i18n.Localf("The character(s) [%s] are not allowed", e.characters)
	}
}

func (e *errorEmptyField) Error() string {
	switch e.fieldName {
	case validationFieldRoomName:
		return i18n.Localf("The field room name cannot be empty")
	case validationFieldChatService:
		return i18n.Localf("The field chat service cannot be empty")
	default:
		return i18n.Localf("The fields cannot be empty")
	}
}

func (e *errorNotValidLocal) Error() string {
	return i18n.Localf("The room name '%s' is not valid", e.local)
}

func (e *errorNotValidDomain) Error() string {
	return i18n.Localf("The chat service name '%s' is not valid", e.domain)
}

func (e *errorNoConnectedAccount) Error() string {
	return i18n.Local("There are not connected accounts")
}

func (v *createMUCRoom) newValidateField(w gtki.Widget) *validateField {
	vi := &validateField{
		widget: w,
	}
	return vi
}

func (v *createMUCRoom) initRulesForAccount() *validateField {
	vf := v.newValidateField(v.account)

	vf.addRule(func() error {
		ac := v.ac.currentAccount()
		if ac == nil {
			return newErrorNoConnectedAccount()
		}

		return nil
	})

	return vf
}

func (v *createMUCRoom) initRulesForRoomName() *validateField {
	vf := v.newValidateField(v.room)

	vf.addRule(func() error {
		s, _ := vf.widget.(gtki.Entry).GetText()
		if s == "" {
			return newErrorEmptyField(validationFieldRoomName)
		}
		return nil
	})

	vf.addRule(func() error {
		s, _ := vf.widget.(gtki.Entry).GetText()
		return checkIfAnyCharacterInField(validationFieldRoomName, s, invalidCharactersForRoomName)
	})

	vf.addRule(func() error {
		s, _ := vf.widget.(gtki.Entry).GetText()
		return isLocalValid(s)
	})

	return vf
}

func (v *createMUCRoom) initRulesForChatService() *validateField {
	vi := v.newValidateField(v.chatServiceEntry)

	vi.addRule(func() error {
		s, _ := vi.widget.(gtki.Entry).GetText()
		if s == "" {
			return newErrorEmptyField(validationFieldChatService)
		}
		return nil
	})

	vi.addRule(func() error {
		s, _ := vi.widget.(gtki.Entry).GetText()
		return checkIfAnyCharacterInField(validationFieldChatService, s, invalidCharactersForChatService)
	})

	vi.addRule(func() error {
		s, _ := vi.widget.(gtki.Entry).GetText()
		return isDomainValid(s)
	})

	return vi
}

func (u *gtkUI) newCreateMUCRoom() *createMUCRoom {
	view := &createMUCRoom{
		u:                u,
		fieldsToValidate: make(map[string]*validateField),
	}

	view.initUIBuilder()
	view.initConnectedAccountsComponent()
	view.initValidationRules()

	return view
}

func (vi *validateField) addRule(r rule) {
	vi.rules = append(vi.rules, r)
}

func (v *createMUCRoom) initValidationRules() {
	v.fieldsToValidate[validationFieldAccount] = v.initRulesForAccount()
	v.fieldsToValidate[validationFieldRoomName] = v.initRulesForRoomName()
	v.fieldsToValidate[validationFieldChatService] = v.initRulesForChatService()
}

func isLocalValid(s string) error {
	if jid.ValidLocal(s) {
		return nil
	}
	return newErrorValidLocal(s)
}

func isDomainValid(s string) error {
	if jid.ValidDomain(s) {
		return nil
	}
	return newErrorValidDomain(s)
}

func checkIfAnyCharacterInField(n, s, pattern string) error {
	var sb strings.Builder
	for _, c := range strings.Split(s, "") {
		if strings.ContainsAny(c, pattern) && !strings.ContainsAny(c, sb.String()) {
			sb.WriteString(c)
		}
	}

	if sb.Len() > 0 {
		return newErrorNotAllowedCharacters(sb.String(), n)
	}

	return nil
}

func (vi *validateField) validate() {
	for _, f := range vi.rules {
		err := f()
		if err != nil {
			vi.isValid = false
			vi.lastError = err
			return
		}
	}

	vi.isValid = true
	vi.lastError = nil
}

func (v *createMUCRoom) doValidationFor(k string, onValidate func()) {
	vi, ok := v.fieldsToValidate[k]
	if !ok {
		return
	}

	vi.validate()
	if onValidate != nil {
		onValidate()
	}
}

func (v *createMUCRoom) validateAll(onValid, onInvalid func()) {
	isValid := true
	for _, vi := range v.fieldsToValidate {
		vi.validate()
		if !vi.isValid {
			isValid = false
		}
	}

	if isValid {
		onValid()
		return
	}

	onInvalid()
}

func (v *createMUCRoom) initUIBuilder() {
	v.builder = newBuilder("MUCCreateRoom")

	panicOnDevError(v.builder.bindObjects(v))

	v.errorBox = newErrorNotification(v.notificationArea)

	v.builder.ConnectSignals(map[string]interface{}{
		"on_create_room":              v.onCreateRoom,
		"on_cancel":                   v.dialog.Destroy,
		"on_close_window":             v.onCloseWindow,
		"on_room_changed":             v.onRoomNameChanged,
		"on_chatServiceEntry_changed": v.onChatServiceChanged,
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
	if !jid.ValidLocal(roomName) || !jid.ValidDomain(service) {
		doInUIThread(func() {
			v.errorBox.ShowMessage(i18n.Localf("The room identity \"%s@%s\" is not valid.", roomName, service))
		})
		return nil
	}

	return jid.NewBare(jid.NewLocal(roomName), jid.NewDomain(service))
}

// createRoom should be called only when all validations
// has passed succesfully
func (v *createMUCRoom) createRoom() {
	ca := v.ac.currentAccount()
	if ca == nil {
		v.errorBox.ShowMessage(i18n.Local("No account is selected, please select one account from the list or connect to one."))
		return
	}

	roomIdentity := v.getRoomID()
	if roomIdentity != nil {
		v.onBeforeToCreateARoom()
		go v.createRoomIfDoesntExist(ca, roomIdentity)
	}
}

func (v *createMUCRoom) onCreateRoom() {
	v.validateAll(v.createRoom, v.checkValidationsAndShowErrorsIfAny)
}

func (v *createMUCRoom) onBeforeToCreateARoom() {
	v.disableOrEnableFields(false)
	v.createButtonPrevText, _ = v.createButton.GetLabel()
	_ = v.createButton.SetProperty("label", i18n.Local("Creating room..."))
}

func (v *createMUCRoom) afterRoomIsCreated() {
	doInUIThread(func() {
		v.disableOrEnableFields(true)
		_ = v.createButton.SetProperty("label", v.createButtonPrevText)
	})
}

func (v *createMUCRoom) createRoomIfDoesntExist(ca *account, ident jid.Bare) {
	erc, ec := ca.session.HasRoom(ident)
	go func() {
		defer v.afterRoomIsCreated()

		select {
		case err, _ := <-ec:
			if err != nil {
				ca.log.WithError(err).Error("Error trying to validate if room exists")
				doInUIThread(func() {
					v.errorBox.ShowMessage(i18n.Local("Could not connect with the server, please try again later."))
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
				v.errorBox.ShowMessage(i18n.Local("The room already exists."))
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
		doInUIThread(func() {
			v.u.mucShowRoom(ca, ident)
			v.dialog.Destroy()
		})
	}
}

func (v *createMUCRoom) checkValidationsAndShowErrorsIfAny() {
	v.clearErrors()

	errorMessages := []string{}
	itHasValidationErrors := false
	for _, vi := range v.fieldsToValidate {
		if !vi.isValid {
			itHasValidationErrors = true
			if vi.lastError != nil {
				errorMessages = append(errorMessages, vi.lastError.Error())
			}
		}
	}

	setEnabled(v.createButton, !itHasValidationErrors)
	if len(errorMessages) > 0 {
		if len(errorMessages) == 1 {
			v.errorBox.ShowMessage(strings.Join(errorMessages, ""))
			return
		}
		v.errorBox.ShowMessage(i18n.Localf("The following errors were found:\n• %s", strings.Join(errorMessages, "\n• ")))
	}
}

func (v *createMUCRoom) onRoomNameChanged() {
	v.handleRoomNameEntered()
	v.doValidationFor(validationFieldRoomName, v.checkValidationsAndShowErrorsIfAny)
}

func (v *createMUCRoom) onChatServiceChanged() {
	v.doValidationFor(validationFieldChatService, v.checkValidationsAndShowErrorsIfAny)
}

func (v *createMUCRoom) handleRoomNameEntered() {
	s, _ := v.room.GetText()
	ri := strings.SplitN(s, "@", 2)
	if len(ri) >= 2 {
		v.room.SetText(ri[0])
		if v.chatServices.GetActiveText() == "" {
			v.chatServiceEntry.SetText(ri[1])
		}
		v.chatServices.SetProperty("is_focus", true)
	}
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
		v.chatServices.RemoveAll()
		v.doValidationFor(validationFieldAccount, v.checkValidationsAndShowErrorsIfAny)
	})
}

func (v *createMUCRoom) updateServicesBasedOnAccount(acc *account) {
	doInUIThread(func() {
		v.doValidationFor(validationFieldAccount, v.checkValidationsAndShowErrorsIfAny)
	})
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

	u.connectShortcutsChildWindow(view.dialog)

	view.dialog.SetTransientFor(u.window)
	view.dialog.Show()
}
