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

	notifications *notificationsComponent
	dialogTitle   string
	formTitle     string
	form          *roomConfigListForm
	formItems     []*mucRoomConfigListFormItem
	onApply       func(jidList []string)
	jidList       []string
}

func (u *gtkUI) newMUCRoomConfigListAddComponent(dialogTitle, formTitle string, onApply func(jidList []string), parent gtki.Window, jidList []string) *mucRoomConfigListAddComponent {
	la := &mucRoomConfigListAddComponent{
		u:           u,
		dialogTitle: dialogTitle,
		formTitle:   formTitle,
		onApply:     onApply,
		jidList:     jidList,
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
		"on_cancel":     la.close,
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
	defaultItem := newMUCRoomConfigListFormItem(la.form, la.appendNewFormItem, nil)
	la.content.PackStart(defaultItem.contentBox(), false, true, 0)
}

// newAddOccupantForm MUST be called from the UI thread
func (la *mucRoomConfigListAddComponent) newAddOccupantForm() *roomConfigListForm {
	return newRoomConfigListForm(
		la.enableApplyIfConditionsAreMet,
		la.onApplyClicked,
	)
}

// appendNewFormItem MUST be called from the UI thread
func (la *mucRoomConfigListAddComponent) appendNewFormItem(jid string) {
	nextIndex := len(la.formItems)

	if la.existJidInList(jid) || la.jidAlreadyInserted(jid) {
		return
	}

	onRemove := func() {
		la.removeItemByIndex(nextIndex)
		la.enableApplyIfConditionsAreMet()
	}

	form := la.newAddOccupantForm()
	form.setJid(jid)

	item := newMUCRoomConfigListFormItem(form, nil, onRemove)
	la.formItems = append(la.formItems, item)
	la.content.PackStart(item.contentBox(), false, true, 0)

	la.enableApplyIfConditionsAreMet()
}

func (la *mucRoomConfigListAddComponent) existJidInList(jid string) bool {
	for _, item := range la.items {
		if jid == item.form.jid() {
			return true
		}
	}
	return false
}

func (la *mucRoomConfigListAddComponent) jidAlreadyInserted(jid string) bool {
	for _, l := range la.jidList {
		if jid == l {
			return true
		}
	}
	return false
}

// removeItemByIndex MUST be called from the UI thread
func (la *mucRoomConfigListAddComponent) removeItemByIndex(index int) {
	items := []*mucRoomConfigListFormItem{}

	for ix, itm := range la.formItems {
		if ix == index {
			la.content.Remove(itm.contentBox())
		} else {
			items = append(items, itm)
		}
	}

	la.formItems = items
}

// forEachForm MUST be called from the UI thread
func (la *mucRoomConfigListAddComponent) forEachForm(fn func(*roomConfigListForm)) {
	for _, itm := range la.formItems {
		fn(itm.form)
	}
}

// areAllFormsFilled MUST be called from the UI thread
func (la *mucRoomConfigListAddComponent) areAllFormsFilled() bool {
	formsAreFilled := la.form.isFilled() || la.hasItems()

	la.forEachForm(func(form *roomConfigListForm) {
		formsAreFilled = formsAreFilled && form.isFilled()
	})

	return formsAreFilled
}

// enableApplyIfConditionsAreMet MUST be called from the UI thread
func (la *mucRoomConfigListAddComponent) enableApplyIfConditionsAreMet() {
	la.removeAllButton.SetSensitive(la.hasItems())

	v := la.areAllFormsFilled()
	la.applyButton.SetSensitive(v)

	if la.hasItems() {
		la.applyButton.SetLabel("Add all")
	} else {
		la.applyButton.SetLabel("Add")
	}
}

// onRemoveAllClicked MUST be called from the UI thread
func (la *mucRoomConfigListAddComponent) onRemoveAllClicked() {
	itms := la.formItems
	la.formItems = nil

	for _, itm := range itms {
		la.content.Remove(itm.contentBox())
	}

	la.form.resetAndFocusJidEntry()
	la.enableApplyIfConditionsAreMet()
}

// onApplyClicked MUST be called from the UI thread
func (la *mucRoomConfigListAddComponent) onApplyClicked() {
	if la.existJidInList(la.form.jid()) || la.jidAlreadyInserted(la.form.jid()) {
		return
	}

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

// isValid MUST be called from the UI thread
func (la *mucRoomConfigListAddComponent) isValid() bool {
	var isValid bool
	var err error

	if la.hasItems() && la.form.isFilled() || !la.hasItems() {
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

// isFormValid MUST be called from the UI thread
func (la *mucRoomConfigListAddComponent) isFormValid(form *roomConfigListForm) (bool, error) {
	if form.jid() == "" {
		return false, errEmptyMemberIdentifier
	}

	if !jid.Parse(form.jid()).Valid() {
		return false, errInvalidMemberIdentifier
	}

	return form.isValid()
}

// notifyError MUST be called from the UI thread
func (la *mucRoomConfigListAddComponent) notifyError(err string) {
	la.notifications.notifyOnError(err)
}

// close MUST be called from the UI thread
func (la *mucRoomConfigListAddComponent) close() {
	la.dialog.Destroy()
}

// show MUST be called from the UI thread
func (la *mucRoomConfigListAddComponent) show() {
	la.dialog.Show()
}

func (la *mucRoomConfigListAddComponent) hasItems() bool {
	return len(la.formItems) > 0
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
