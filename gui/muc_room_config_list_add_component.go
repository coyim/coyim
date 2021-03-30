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
	form                   *roomConfigListForm
	items                  []*mucRoomConfigListFormItem
	addOccupantFormCreator func(onFieldChanged, onFieldActivate func()) *roomConfigListForm
	onApply                func(jidList []string)
}

func (u *gtkUI) newMUCRoomConfigListAddComponent(dialogTitle, formTitle string, addOccupantForm func(onFieldChanged, onFieldActivate func()) *roomConfigListForm, onApply func(jidList []string), parent gtki.Window) *mucRoomConfigListAddComponent {
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
	la.notificationBox.Add(la.notifications.contentBox())
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
	la.content.PackStart(defaultItem.contentBox(), false, true, 0)
}

func (la *mucRoomConfigListAddComponent) newAddOccupantForm() *roomConfigListForm {
	return la.addOccupantFormCreator(
		la.enableApplyIfConditionsAreMet,
		la.onApplyClicked,
	)
}

func (la *mucRoomConfigListAddComponent) appendNewItem(jid string) {
	nextIndex := len(la.items)

	onRemove := func() {
		la.removeItemByIndex(nextIndex)
		la.enableApplyIfConditionsAreMet()
	}

	form := la.newAddOccupantForm()
	form.setValue(jid)

	item := newMUCRoomConfigListFormItem(form, nil, onRemove)
	la.items = append(la.items, item)
	la.content.PackStart(item.contentBox(), false, true, 0)

	la.enableApplyIfConditionsAreMet()
}

func (la *mucRoomConfigListAddComponent) removeItemByIndex(index int) {
	items := []*mucRoomConfigListFormItem{}
	for ix, itm := range la.items {
		if ix == index {
			la.content.Remove(itm.contentBox())
			continue
		}
		items = append(items, itm)
	}
	la.items = items
}

func (la *mucRoomConfigListAddComponent) forEachForm(fn func(*roomConfigListForm)) {
	for _, itm := range la.items {
		fn(itm.form)
	}
}

func (la *mucRoomConfigListAddComponent) areAllFormsFilled() bool {
	formsAreFilled := la.form.isFilled() || len(la.items) > 0

	la.forEachForm(func(form *roomConfigListForm) {
		formsAreFilled = formsAreFilled && form.isFilled()
	})

	return formsAreFilled
}

func (la *mucRoomConfigListAddComponent) enableApplyIfConditionsAreMet() {
	la.removeAllButton.SetSensitive(len(la.items) > 0)

	v := la.areAllFormsFilled()
	la.applyButton.SetSensitive(v)
	if v {
		if la.hasMoreThanOneItem() {
			la.applyButton.SetLabel("Add all")
			return
		}
		la.applyButton.SetLabel("Add")
	}
}

func (la *mucRoomConfigListAddComponent) onCancelClicked() {
	la.close()
}

func (la *mucRoomConfigListAddComponent) onRemoveAllClicked() {
	for _, itm := range la.items {
		la.content.Remove(itm.contentBox())
	}

	la.items = nil

	la.form.reset()
	la.form.focus()

	la.enableApplyIfConditionsAreMet()
}

func (la *mucRoomConfigListAddComponent) onApplyClicked() {
	if la.isValid() {
		jidList := []string{}

		if la.form.isFilled() {
			jidList = append(jidList, la.form.jid())
		}

		la.forEachForm(func(form *roomConfigListForm) {
			jidList = append(jidList, form.jid())
		})

		la.onApply(jidList)
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

	la.forEachForm(func(form *roomConfigListForm) {
		if isValid, err = la.isFormValid(form); err != nil {
			la.notifyError(la.friendlyErrorMessage(form, err))
		}
	})

	return isValid
}

func (la *mucRoomConfigListAddComponent) isFormValid(form *roomConfigListForm) (bool, error) {
	if form.jid() == "" {
		return false, errEmptyMemberIdentifier
	}

	j := jid.Parse(form.jid())
	if !j.Valid() {
		return false, errInvalidMemberIdentifier
	}

	return form.isValid()
}

func (la *mucRoomConfigListAddComponent) hasMoreThanOneItem() bool {
	count := len(la.items)
	if la.form.isFilled() {
		count = count + 1
	}
	return count > 1
}

func (la *mucRoomConfigListAddComponent) friendlyErrorMessage(form *roomConfigListForm, err error) string {
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
	form *roomConfigListForm

	box          gtki.Box    `gtk-widget:"room-config-list-add-item-box"`
	formBox      gtki.Box    `gtk-widget:"room-config-list-add-item-form-box"`
	addButton    gtki.Button `gtk-widget:"room-config-list-add-item-button"`
	removeButton gtki.Button `gtk-widget:"room-config-list-remove-item-button"`
}

func newMUCRoomConfigListFormItem(form *roomConfigListForm, onAdd func(jid string), onRemove func()) *mucRoomConfigListFormItem {
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
			onAdd(form.jid())
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

func (lfi *mucRoomConfigListFormItem) contentBox() gtki.Box {
	return lfi.box
}
