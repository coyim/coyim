package gui

import (
	"errors"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

var (
	errEmptyMemberIdentifier   = errors.New("empty member identifier (jid)")
	errInvalidMemberIdentifier = errors.New("invalid member identifier (jid)")
)

type mucRoomConfigListAddComponent struct {
	u *gtkUI

	dialog          gtki.Dialog `gtk-widget:"room-config-list-add-dialog"`
	title           gtki.Label  `gtk-widget:"room-config-list-add-title"`
	content         gtki.Box    `gtk-widget:"room-config-list-add-content"`
	applyButton     gtki.Button `gtk-widget:"room-config-list-add-apply"`
	notificationBox gtki.Box    `gtk-widget:"notification-box"`

	notifications *notifications
	dialogTitle   string
	formTitle     string
	form          mucRoomConfigListForm
	onApply       func(...string)
}

func (u *gtkUI) newMUCRoomConfigListAddComponent(dialogTitle, formTitle string, addOccupantForm func(onFieldChanged, onFieldActivate func()) mucRoomConfigListForm, onApply func(...string), parent gtki.Window) *mucRoomConfigListAddComponent {
	la := &mucRoomConfigListAddComponent{
		u:           u,
		dialogTitle: dialogTitle,
		formTitle:   formTitle,
		onApply:     onApply,
	}

	la.initBuilder()
	la.initNotifications()
	la.initAddOccupantForm(addOccupantForm)
	la.initDefaults(parent)

	return la
}

func (la *mucRoomConfigListAddComponent) initBuilder() {
	builder := newBuilder("MUCRoomConfigListAddDialog")
	panicOnDevError(builder.bindObjects(la))

	builder.ConnectSignals(map[string]interface{}{
		"on_cancel": la.onCancelClicked,
		"on_apply":  la.onApplyClicked,
	})
}

func (la *mucRoomConfigListAddComponent) initNotifications() {
	la.notifications = la.u.newNotificationsComponent()
	la.notificationBox.Add(la.notifications.widget())
}

func (la *mucRoomConfigListAddComponent) initDefaults(parent gtki.Window) {
	la.dialog.SetTitle(la.dialogTitle)
	la.dialog.SetTransientFor(parent)
	la.title.SetLabel(la.formTitle)
}

func (la *mucRoomConfigListAddComponent) initAddOccupantForm(addOccupantForm func(onFieldChanged, onFieldActivate func()) mucRoomConfigListForm) {
	la.form = addOccupantForm(
		la.onAddOccupantFormFieldChanged,
		la.onApplyClicked,
	)

	la.content.Add(la.form.getFormView())
}

func (la *mucRoomConfigListAddComponent) onAddOccupantFormFieldChanged() {
	la.applyButton.SetSensitive(la.form.isFilled())
}

func (la *mucRoomConfigListAddComponent) onCancelClicked() {
	la.close()
}

func (la *mucRoomConfigListAddComponent) onApplyClicked() {
	if la.isValid() {
		la.onApply(la.form.getValues()...)
		la.close()
	}
}

func (la *mucRoomConfigListAddComponent) isValid() bool {
	if ok, err := la.isFormValid(); !ok {
		la.notifyError(la.friendlyErrorMessage(err))
		return false
	}
	return true
}

func (la *mucRoomConfigListAddComponent) isFormValid() (bool, error) {
	if la.form.jid() == "" {
		return false, errEmptyMemberIdentifier
	}

	j := jid.Parse(la.form.jid())
	if !j.Valid() {
		return false, errInvalidMemberIdentifier
	}

	return la.form.isValid()
}

func (la *mucRoomConfigListAddComponent) friendlyErrorMessage(err error) string {
	switch err {
	case errEmptyMemberIdentifier:
		return i18n.Local("You must enter the member identifier (JID)")
	case errInvalidMemberIdentifier:
		return i18n.Local("You must provide a valid member identifier (JID)")
	default:
		return la.form.friendlyErrorMessage(err)
	}
}

func (la *mucRoomConfigListAddComponent) notifyError(err string) {
	la.notifications.notifyOnError(err)
}

func (la *mucRoomConfigListAddComponent) close() {
	la.dialog.Destroy()
}

func (la *mucRoomConfigListAddComponent) show() {
	la.dialog.Show()
}
