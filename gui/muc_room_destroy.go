package gui

import (
	"errors"

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
	reasonTextView       gtki.TextView    `gtk-widget:"destroy-room-reason-entry"`
	alternativeRoomCheck gtki.CheckButton `gtk-widget:"destroy-room-alternative-check"`
	alternativeRoomBox   gtki.Box         `gtk-widget:"destroy-room-alternative-box"`
	alternativeRoomEntry gtki.Entry       `gtk-widget:"destroy-room-name-entry"`
	passwordEntry        gtki.Entry       `gtk-widget:"destroy-room-password-entry"`
	destroyRoomButton    gtki.Button      `gtk-widget:"destroy-room-button"`
	notificationBox      gtki.Box         `gtk-widget:"notification-area"`

	notifications *notifications
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
		"on_alternative_room_toggled": d.onAlternativeRoomToggled,
		"on_destroy":                  d.onDestroyRoom,
		"on_cancel":                   d.close,
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

	d.notifications = v.u.newNotificationsComponent()
	d.notificationBox.Add(d.notifications.widget())
}

// onDestroyRoom MUST be called from the UI thread
func (d *roomDestroyView) onDestroyRoom() {
	d.notifications.clearErrors()

	alternativeID, password, err := d.alternativeRoomInformation()
	if err != nil {
		d.notifications.error(d.friendlyMessageForAlternativeRoomError(err))
		return
	}

	b, _ := d.reasonTextView.GetBuffer()
	reason := b.GetText(b.GetStartIter(), b.GetEndIter(), false)

	go d.destroyRoom(reason, alternativeID, password)
	d.close()
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
	d.reasonTextView.SetSensitive(v)
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
