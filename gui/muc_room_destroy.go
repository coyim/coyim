package gui

import (
	"errors"

	"github.com/coyim/coyim/coylog"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

func (v *roomView) onDestroyRoom() {
	d := v.newRoomDestroyView()
	d.show()
}

var (
	errEmptyServiceName   = errors.New("empty service name")
	errEmptyRoomName      = errors.New("empty room name")
	errInvalidRoomName    = errors.New("invalid room name")
	errInvalidServiceName = errors.New("invalid service name")
)

type roomDestroyView struct {
	builder               *builder
	chatServicesComponent *chatServicesComponent
	destroyRoom           func(reason string, alternativeID jid.Bare, password string)

	dialog               gtki.Dialog      `gtk-widget:"destroy-room-dialog"`
	reasonEntry          gtki.TextView    `gtk-widget:"destroy-room-reason-entry"`
	alternativeRoomCheck gtki.CheckButton `gtk-widget:"destroy-room-alternative-check"`
	alternativeRoomBox   gtki.Box         `gtk-widget:"destroy-room-alternative-box"`
	alternativeRoomLabel gtki.Label       `gtk-widget:"destroy-room-name-label"`
	alternativeRoomEntry gtki.Entry       `gtk-widget:"destroy-room-name-entry"`
	chatServicesLabel    gtki.Label       `gtk-widget:"destroy-room-service-label"`
	passwordLabel        gtki.Label       `gtk-widget:"destroy-room-password-label"`
	passwordEntry        gtki.Entry       `gtk-widget:"destroy-room-password-entry"`
	destroyRoomButton    gtki.Button      `gtk-widget:"destroy-room-button"`
	spinnerBox           gtki.Box         `gtk-widget:"destroy-room-spinner-box"`
	notificationBox      gtki.Box         `gtk-widget:"notification-area"`

	spinner      *spinner
	notification *notifications
}

func (v *roomView) newRoomDestroyView() *roomDestroyView {
	d := &roomDestroyView{}

	d.initBuilder()
	d.initDestroyContext(v)
	d.initChatServices(v)
	d.initDefaults(v)

	return d
}

func (d *roomDestroyView) initBuilder() {
	d.builder = newBuilder("MUCRoomDestroyDialog")
	panicOnDevError(d.builder.bindObjects(d))

	d.builder.ConnectSignals(map[string]interface{}{
		"on_destroy_clicked":          d.onDestroyRoom,
		"on_alternative_room_toggled": d.onAlternativeRoomToggled,
		"on_cancel_clicked":           d.close,
	})
}

func (d *roomDestroyView) initDestroyContext(v *roomView) {
	d.destroyRoom = func(reason string, alternativeID jid.Bare, password string) {
		d.close()
		ctx := v.newDestroyContext(reason, alternativeID, password)
		ctx.destroyRoom()
	}
}

func (d *roomDestroyView) initChatServices(v *roomView) {
	chatServicesList := d.builder.get("chat-services-list").(gtki.ComboBoxText)
	chatServicesEntry := d.builder.get("chat-services-entry").(gtki.Entry)
	d.chatServicesComponent = v.u.createChatServicesComponent(chatServicesList, chatServicesEntry, nil)
	go d.chatServicesComponent.updateServicesBasedOnAccount(v.account)
}

func (d *roomDestroyView) initDefaults(v *roomView) {
	d.dialog.SetTransientFor(v.window)

	d.spinner = newSpinner()
	d.spinnerBox.Add(d.spinner.getWidget())

	d.notification = v.u.newNotifications(d.notificationBox)
}

// onDestroyRoom MUST be called from the UI thread
func (d *roomDestroyView) onDestroyRoom() {
	d.disableFieldsAndShowSpinner()

	b, _ := d.reasonEntry.GetBuffer()
	reason := b.GetText(b.GetStartIter(), b.GetEndIter(), false)

	alternativeID, password, err := d.alternativeRoomInformation()
	if err != nil {
		d.notification.error(d.friendlyMessageForAlternativeRoomError(err))
		d.enableFieldsAndHideSpinner()
		return
	}

	d.destroyRoom(reason, alternativeID, password)
}

func (d *roomDestroyView) alternativeRoomInformation() (jid.Bare, string, error) {
	if !d.alternativeRoomCheck.GetActive() {
		return nil, "", nil
	}

	alternativeID, err := d.tryParseAlternativeRoomID()
	if err != nil {
		return nil, "", err
	}

	password, _ := d.passwordEntry.GetText()

	return alternativeID, password, nil
}

// onAlternativeRoomToggled MUST be called from the UI thread
func (d *roomDestroyView) onAlternativeRoomToggled() {
	v := d.alternativeRoomCheck.GetActive()
	d.alternativeRoomBox.SetVisible(v)
	d.resetAlternativeRoomFields()
}

func (d *roomDestroyView) resetAlternativeRoomFields() {
	d.alternativeRoomEntry.SetText("")
	d.passwordEntry.SetText("")
	d.chatServicesComponent.resetToDefault()
}

func (d *roomDestroyView) friendlyMessageForAlternativeRoomError(err error) string {
	switch err {
	case errEmptyServiceName:
		return i18n.Local("You must provide a service name")
	case errEmptyRoomName:
		return i18n.Local("You must provide a room name")
	case errInvalidRoomName:
		return i18n.Local("You must provide a valid room name")
	case errInvalidServiceName:
		return i18n.Local("You must provide a valid service name")
	default:
		return i18n.Local("You must provide a valid service and room name")
	}
}

// tryParseAlternativeRoomID MUST be called from the UI thread
//
// This should be "alternative venue" as the protocol says, but
// we prefer to use "alternative room id" in this context
// in order to have a better understanding of what this field means
func (d *roomDestroyView) tryParseAlternativeRoomID() (jid.Bare, error) {
	rn, _ := d.alternativeRoomEntry.GetText()
	s := d.chatServicesComponent.currentServiceValue()

	// We don't really need to continue if the user hasn't entered
	// anything in the room name and the service, because the alternative
	// room is always optional according to the protocol
	if rn == "" && s == "" {
		return nil, nil
	}

	r, err := d.alternativeRoomID()
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (d *roomDestroyView) alternativeRoomID() (r jid.Bare, err error) {
	l, err := d.validateRoomName()
	if err != nil {
		return
	}

	s, err := d.validateServiceName()
	if err != nil {
		return
	}

	r = jid.NewBare(l, s)
	return
}

func (d *roomDestroyView) validateRoomName() (l jid.Local, err error) {
	rn, _ := d.alternativeRoomEntry.GetText()

	if rn == "" {
		err = errEmptyRoomName
		return
	}

	l = jid.NewLocal(rn)
	if !l.Valid() {
		err = errInvalidRoomName
	}

	return
}

func (d *roomDestroyView) validateServiceName() (s jid.Domain, err error) {
	if !d.chatServicesComponent.hasServiceValue() {
		err = errEmptyServiceName
		return
	}

	s = jid.NewDomain(d.chatServicesComponent.currentServiceValue())
	if !s.Valid() {
		err = errInvalidServiceName
	}

	return
}

// disableFields MUST be called from the UI thread
func (d *roomDestroyView) disableFields() {
	d.setSensitivityForAllFields(false)
}

// enableFields MUST be called from the UI thread
func (d *roomDestroyView) enableFields() {
	d.setSensitivityForAllFields(true)
}

// setSensitivityForAllFields MUST be called from the UI thread
func (d *roomDestroyView) setSensitivityForAllFields(v bool) {
	d.reasonEntry.SetSensitive(v)
	d.alternativeRoomEntry.SetSensitive(v)
	d.destroyRoomButton.SetSensitive(v)
}

// show MUST be called from the UI thread
func (d *roomDestroyView) show() {
	d.dialog.Show()
}

// close MUST be called from the UI thread
func (d *roomDestroyView) close() {
	d.dialog.Destroy()
}

// disableFieldsAndShowSpinner MUST be called from the UI thread
func (d *roomDestroyView) disableFieldsAndShowSpinner() {
	d.disableFields()
	d.spinner.show()
}

// enableFieldsAndHideSpinner MUST be called from the UI thread
func (d *roomDestroyView) enableFieldsAndHideSpinner() {
	d.enableFields()
	d.spinner.hide()
}

type roomDestroyContext struct {
	roomID        jid.Bare
	reason        string
	alternativeID jid.Bare
	password      string
	destroy       func(reason string, alternativeID jid.Bare, password string, onSuccess func(), onError func(error))
	log           coylog.Logger
}

func (v *roomView) newDestroyContext(reason string, alternativeID jid.Bare, password string) *roomDestroyContext {
	return &roomDestroyContext{
		roomID:        v.roomID(),
		reason:        reason,
		alternativeID: alternativeID,
		password:      password,
		destroy:       v.tryDestroyRoom,
		log:           v.log,
	}
}

func (dc *roomDestroyContext) destroyRoom() {
	dc.destroy(dc.reason, dc.alternativeID, dc.password, dc.onDestroySuccess, dc.onDestroyFails)
}

func (dc *roomDestroyContext) onDestroySuccess() {
	dc.log.Info("The room has been destroyed")
}

func (dc *roomDestroyContext) onDestroyFails(err error) {
	doInUIThread(func() {
		rd := newDestroyError(dc.roomID, err, dc.destroyRoom)
		rd.show()
	})
}
