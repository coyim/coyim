package gui

import (
	"errors"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session"
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
	destroyRoom           func(jid.Bare, string, func(), func(error))

	dialog               gtki.Dialog      `gtk-widget:"destroy-room-dialog"`
	reasonEntry          gtki.Entry       `gtk-widget:"destroy-room-reason-entry"`
	alternativeRoomCheck gtki.CheckButton `gtk-widget:"destroy-room-alternative-check"`
	alternativeRoomLabel gtki.Label       `gtk-widget:"destroy-room-name-label"`
	alternativeRoomEntry gtki.Entry       `gtk-widget:"destroy-room-name-entry"`
	chatServicesLabel    gtki.Label       `gtk-widget:"destroy-room-service-label"`
	chatServicesBox      gtki.Box         `gtk-widget:"chat-services-box"`
	destroyRoomButton    gtki.Button      `gtk-widget:"destroy-room-button"`
	spinnerBox           gtki.Box         `gtk-widget:"destroy-room-spinner-box"`
	notificationBox      gtki.Box         `gtk-widget:"notification-area"`

	spinner      *spinner
	notification *notifications

	cancelChannel chan bool
}

func (v *roomView) newRoomDestroyView() *roomDestroyView {
	d := &roomDestroyView{
		destroyRoom: v.tryDestroyRoom,
	}

	d.initBuilder()
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
		"on_cancel_clicked":           d.onCancel,
		"on_dialog_destroyed":         d.onDialogDestroy,
	})
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

	reason := d.getReason()

	alternativeID, err := d.tryParseAlternativeRoomID()
	if err != nil {
		d.notification.error(d.getMessageForAlternativeRoomError(err))
		d.enableFieldsAndHideSpinner()
		return
	}

	d.destroyRoom(alternativeID, reason, d.onDestroySuccess, d.onDestroyFails)
}

// onAlternativeRoomToggled MUST be called from the UI thread
func (d *roomDestroyView) onAlternativeRoomToggled() {
	v := d.alternativeRoomCheck.GetActive()

	setFieldVisibility(d.alternativeRoomLabel, v)
	setFieldVisibility(d.alternativeRoomEntry, v)
	setFieldVisibility(d.chatServicesLabel, v)
	setFieldVisibility(d.chatServicesBox, v)
}

func setFieldVisibility(w gtki.Widget, v bool) {
	w.SetVisible(v)
}

func (d *roomDestroyView) getMessageForAlternativeRoomError(err error) string {
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
		return i18n.Local("The room identification is not valid")
	}
}

// onDestroySuccess MUST NOT be called from the UI thread
func (d *roomDestroyView) onDestroySuccess() {
	doInUIThread(d.close)
}

// onDestroyFails MUST NOT be called from the UI thread
func (d *roomDestroyView) onDestroyFails(err error) {
	doInUIThread(func() {
		d.enableFields()
		d.notification.error(d.getMessageForDestroyError(err))
	})
}

func (d *roomDestroyView) getMessageForDestroyError(err error) string {
	switch err {
	case session.ErrDestroyRoomInvalidIQResponse:
		return i18n.Local("We were able to connect to the room service, " +
			"but we received an invalid response from it. Please try again.")
	case session.ErrDestroyRoomForbidden:
		return i18n.Local("You don't have the permission to destroy this room. " +
			"Please get in contact with one of the room owners.")
	case session.ErrDestroyRoomUnknown:
		return i18n.Local("The room's service responded with an unknow error, " +
			"so, the room can't be destroyed. Please try again.")
	case session.ErrDestroyRoomNoResult:
		return i18n.Local("We were able to send the request to destroy the room, " +
			"but the service responded with an unknow result. Please contact the " +
			"room's administrator.")
	default:
		return i18n.Local("An error occurred while destroying the room, please try again.")
	}
}

// onCancel MUST be called from the UI thread
func (d *roomDestroyView) onCancel() {
	d.cancelActiveRequest()

	d.spinner.hide()
	d.close()
}

// onDialogDestroy MUST be called from the UI thread
func (d *roomDestroyView) onDialogDestroy() {
	d.cancelActiveRequest()
}

// cancelActiveRequest MUST be called from the UI thread
func (d *roomDestroyView) cancelActiveRequest() {
	if d.cancelChannel != nil {
		d.cancelChannel <- true
	}
}

// getReason MUST be called from the UI thread
func (d *roomDestroyView) getReason() string {
	t, _ := d.reasonEntry.GetText()
	return t
}

// tryParseAlternativeRoomID MUST be called from the UI thread
//
// This should be "alternative venue" as the protocol says, but
// we prefer to use "alternative room id" in this context
// in order to have a better understanding of what this field means
func (d *roomDestroyView) tryParseAlternativeRoomID() (jid.Bare, error) {
	rn, _ := d.alternativeRoomEntry.GetText()
	s := d.chatServicesComponent.currentServiceValue()

	ch := newAlternativeRoomChecker(rn, s)

	// We don't really need to continue if the user hasn't entered
	// anything in the room name and the service, because the alternative
	// room is always optional according to the protocol
	if ch.shouldBypassChecking() {
		return nil, nil
	}

	roomID, err := ch.alternativeRoomID()
	if err != nil {
		return nil, err
	}

	return roomID, nil
}

type alternativeRoomChecker struct {
	roomName, service string
}

func newAlternativeRoomChecker(rn, s string) *alternativeRoomChecker {
	return &alternativeRoomChecker{rn, s}
}

func (c *alternativeRoomChecker) shouldBypassChecking() bool {
	return c.roomName == "" && c.service == ""
}

func (c *alternativeRoomChecker) alternativeRoomID() (jid.Bare, error) {
	err := c.doAllChecks()
	if err != nil {
		return nil, err
	}
	return jid.NewBare(jid.NewLocal(c.roomName), jid.NewDomain(c.service)), nil
}

func (c *alternativeRoomChecker) doAllChecks() error {
	rules := []func() (bool, error){
		c.validateRoomName,
		c.validateServiceName,
	}

	for _, ch := range rules {
		invalid, err := ch()
		if invalid {
			return err
		}
	}

	return nil
}

func (c *alternativeRoomChecker) validateRoomName() (bool, error) {
	if c.roomName == "" && c.service != "" {
		return true, errEmptyRoomName
	}

	l := jid.NewLocal(c.roomName)
	if !l.Valid() {
		return true, errInvalidRoomName
	}

	return false, nil
}

func (c *alternativeRoomChecker) validateServiceName() (bool, error) {
	if c.roomName != "" && c.service == "" {
		return true, errEmptyServiceName
	}

	d := jid.NewDomain(c.service)
	if !d.Valid() {
		return true, errInvalidServiceName
	}

	return false, nil
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
