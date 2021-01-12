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

	notifications          *notifications
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
	la.notificationBox.Add(la.notifications.widget())
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
	la.content.PackStart(defaultItem.widget(), false, true, 0)
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
	la.content.PackStart(item.widget(), false, true, 0)

	la.enableApplyIfConditionsAreMet()
}

func (la *mucRoomConfigListAddComponent) removeItemByIndex(index int) {
	items := []*mucRoomConfigListFormItem{}
	for ix, itm := range la.items {
		if ix == index {
			la.content.Remove(itm.widget())
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
		la.content.Remove(itm.widget())
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
	box  gtki.Box
}

func newMUCRoomConfigListFormItem(form mucRoomConfigListForm, onAdd func([]string), onRemove func()) *mucRoomConfigListFormItem {
	item := &mucRoomConfigListFormItem{form: form}

	var err error
	item.box, err = g.gtk.BoxNew(gtki.HorizontalOrientation, 12)
	if err != nil {
		panic(err)
	}

	item.box.PackStart(form.getFormView(), false, true, 0)

	item.box.SetHExpand(true)
	item.box.SetVExpand(false)

	item.box.SetHAlign(gtki.ALIGN_FILL)
	item.box.SetVAlign(gtki.ALIGN_START)

	if onAdd != nil {
		addButton, _ := g.gtk.ButtonNewWithLabel(i18n.Local("Add"))
		addButton.SetSensitive(false)
		addButton.Connect("clicked", func() {
			onAdd(form.getValues())
			form.reset()
			form.focus()
		})
		form.onFieldChanged(func() {
			addButton.SetSensitive(form.isFilled())
		})
		item.box.PackEnd(addButton, false, false, 0)
	}

	if onRemove != nil {
		removeButton, _ := g.gtk.ButtonNewWithLabel(i18n.Local("Remove"))
		removeButton.Connect("clicked", onRemove)
		item.box.PackEnd(removeButton, false, false, 0)
	}

	item.box.ShowAll()

	return item
}

func (lfi *mucRoomConfigListFormItem) widget() gtki.Box {
	return lfi.box
}
