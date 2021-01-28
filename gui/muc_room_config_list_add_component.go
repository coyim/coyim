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
	content         gtki.Box    `gtk-widget:"room-config-list-add-body"`
	removeAllButton gtki.Button `gtk-widget:"room-config-list-remove-all-button"`
	applyButton     gtki.Button `gtk-widget:"room-config-list-add-apply"`
	notificationBox gtki.Box    `gtk-widget:"notification-box"`

	notifications          *notificationsComponent
	dialogTitle            string
	formTitle              string
	form                   mucRoomConfigListForm
	items                  []*mucRoomConfigListFormItem
	addOccupantFormCreator func(onFieldChanged, onFieldActivate func()) mucRoomConfigListForm
	onApply                func([][]string)
}

func (u *gtkUI) newMUCRoomConfigListAddComponent(dialogTitle, formTitle string, addOccupantForm func(onFieldChanged, onFieldActivate func()) mucRoomConfigListForm, onApply func([][]string), parent gtki.Window) *mucRoomConfigListAddComponent {
	la := &mucRoomConfigListAddComponent{
		u:                      u,
		dialogTitle:            dialogTitle,
		formTitle:              formTitle,
		addOccupantFormCreator: addOccupantForm,
		onApply:                onApply,
	}

	la.initBuilder()
	la.initNotifications()
	la.initAddOccupantForm()
	la.initDefaults(parent)

	return la
}

func (la *mucRoomConfigListAddComponent) initBuilder() {
	builder := newBuilder("MUCRoomConfigListAddDialog")
	panicOnDevError(builder.bindObjects(la))

	builder.ConnectSignals(map[string]interface{}{
		"on_cancel":     la.onCancelClicked,
		"on_remove_all": la.onRemoveAllClicked,
		"on_apply":      la.onApplyClicked,
	})
}

func (la *mucRoomConfigListAddComponent) initNotifications() {
	la.notifications = la.u.newNotificationsComponent()
	la.notificationBox.Add(la.notifications.getBox())
}

func (la *mucRoomConfigListAddComponent) initDefaults(parent gtki.Window) {
	la.removeAllButton.SetSensitive(false)

	la.dialog.SetTitle(la.dialogTitle)
	la.dialog.SetTransientFor(parent)
	la.title.SetLabel(la.formTitle)
}

func (la *mucRoomConfigListAddComponent) initAddOccupantForm() {
	la.form = la.newAddOccupantForm()
	defaultItem := newMUCRoomConfigListFormItem(la.form, la.appendNewItem, nil)
	la.content.PackStart(defaultItem.getBox(), false, true, 0)
}

func (la *mucRoomConfigListAddComponent) newAddOccupantForm() mucRoomConfigListForm {
	return la.addOccupantFormCreator(
		la.enableApplyIfConditionsAreMet,
		la.onApplyClicked,
	)
}

func (la *mucRoomConfigListAddComponent) appendNewItem(values []string) {
	nextIndex := len(la.items)

	onRemove := func() {
		la.removeItemByIndex(nextIndex)
		la.enableApplyIfConditionsAreMet()
	}

	form := la.newAddOccupantForm()
	form.setValues(values)

	item := newMUCRoomConfigListFormItem(form, nil, onRemove)
	la.items = append(la.items, item)
	la.content.PackStart(item.getBox(), false, true, 0)

	la.enableApplyIfConditionsAreMet()
}

func (la *mucRoomConfigListAddComponent) removeItemByIndex(index int) {
	items := []*mucRoomConfigListFormItem{}
	for ix, itm := range la.items {
		if ix == index {
			la.content.Remove(itm.getBox())
			continue
		}
		items = append(items, itm)
	}
	la.items = items
}

func (la *mucRoomConfigListAddComponent) forEachForm(fn func(mucRoomConfigListForm)) {
	for _, itm := range la.items {
		fn(itm.form)
	}
}

func (la *mucRoomConfigListAddComponent) areAllFormsFilled() bool {
	formsAreFilled := la.form.isFilled() || len(la.items) > 0

	la.forEachForm(func(form mucRoomConfigListForm) {
		formsAreFilled = formsAreFilled && form.isFilled()
	})

	return formsAreFilled
}

func (la *mucRoomConfigListAddComponent) enableApplyIfConditionsAreMet() {
	la.removeAllButton.SetSensitive(len(la.items) > 0)
	la.applyButton.SetSensitive(la.areAllFormsFilled())
}

func (la *mucRoomConfigListAddComponent) onCancelClicked() {
	la.close()
}

func (la *mucRoomConfigListAddComponent) onRemoveAllClicked() {
	for _, itm := range la.items {
		la.content.Remove(itm.getBox())
	}

	la.items = nil

	la.form.reset()
	la.form.focus()

	la.enableApplyIfConditionsAreMet()
}

func (la *mucRoomConfigListAddComponent) onApplyClicked() {
	if la.isValid() {
		entries := [][]string{}

		if la.form.isFilled() {
			entries = append(entries, la.form.getValues())
		}

		la.forEachForm(func(form mucRoomConfigListForm) {
			entries = append(entries, form.getValues())
		})

		la.onApply(entries)
		la.close()
	}
}

func (la *mucRoomConfigListAddComponent) isValid() bool {
	var isValid bool
	var err error

	hasNoItems := len(la.items) == 0
	if !hasNoItems && la.form.jid() != "" || hasNoItems {
		isValid, err = la.isFormValid(la.form)
		if err != nil {
			la.notifyError(la.friendlyErrorMessage(la.form, err))
			return false
		}
	}

	la.forEachForm(func(form mucRoomConfigListForm) {
		if isValid, err = la.isFormValid(form); err != nil {
			la.notifyError(la.friendlyErrorMessage(form, err))
		}
	})

	return isValid
}

func (la *mucRoomConfigListAddComponent) isFormValid(form mucRoomConfigListForm) (bool, error) {
	if form.jid() == "" {
		return false, errEmptyMemberIdentifier
	}

	j := jid.Parse(form.jid())
	if !j.Valid() {
		return false, errInvalidMemberIdentifier
	}

	return form.isValid()
}

func (la *mucRoomConfigListAddComponent) friendlyErrorMessage(form mucRoomConfigListForm, err error) string {
	switch err {
	case errEmptyMemberIdentifier:
		return i18n.Local("You must enter the member identifier (also kown as JID)")
	case errInvalidMemberIdentifier:
		return i18n.Local("You must provide a valid member identifier (also kown as JID)")
	default:
		return form.friendlyErrorMessage(err)
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

type mucRoomConfigListFormItem struct {
	form mucRoomConfigListForm

	box          gtki.Box    `gtk-widget:"room-config-list-add-item-box"`
	formBox      gtki.Box    `gtk-widget:"room-config-list-add-item-form-box"`
	addButton    gtki.Button `gtk-widget:"room-config-list-add-item-button"`
	removeButton gtki.Button `gtk-widget:"room-config-list-remove-item-button"`
}

func newMUCRoomConfigListFormItem(form mucRoomConfigListForm, onAdd func([]string), onRemove func()) *mucRoomConfigListFormItem {
	lfi := &mucRoomConfigListFormItem{
		form: form,
	}

	builder := newBuilder("MUCRoomConfigListAddFormItem")
	panicOnDevError(builder.bindObjects(lfi))

	lfi.formBox.Add(lfi.form.getFormView())

	lfi.addButton.SetSensitive(false)
	lfi.removeButton.SetSensitive(false)

	lfi.addButton.SetVisible(false)
	lfi.removeButton.SetVisible(false)

	if onAdd != nil {
		lfi.addButton.Connect("clicked", func() {
			onAdd(form.getValues())
			form.reset()
			form.focus()
		})

		form.onFieldChanged(func() {
			lfi.addButton.SetSensitive(form.isFilled())
		})

		lfi.addButton.SetVisible(true)
	}

	if onRemove != nil {
		lfi.removeButton.Connect("clicked", onRemove)
		lfi.removeButton.SetSensitive(true)
		lfi.removeButton.SetVisible(true)
	}

	return lfi
}

func (lfi *mucRoomConfigListFormItem) getBox() gtki.Box {
	return lfi.box
}
