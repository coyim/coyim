package gui

import "github.com/coyim/gotk3adapter/gtki"

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

	lfi.initBuilder()
	lfi.initListAdd(onAdd)
	lfi.initListRemove(onRemove)
	lfi.initDefaults()

	return lfi
}

func (lfi *mucRoomConfigListFormItem) initBuilder() {
	builder := newBuilder("MUCRoomConfigListAddFormItem")
	panicOnDevError(builder.bindObjects(lfi))
}

func (lfi *mucRoomConfigListFormItem) initListAdd(onAdd func(jid string)) {
	lfi.addButton.SetSensitive(false)
	lfi.addButton.SetVisible(false)

	if onAdd != nil {
		lfi.addButton.Connect("clicked", func() {
			onAdd(lfi.form.jid())
			lfi.form.resetAndFocusJidEntry()
		})

		lfi.form.onFieldChanged(func() {
			lfi.addButton.SetSensitive(lfi.form.isFilled())
		})

		lfi.addButton.SetVisible(true)
	}
}

func (lfi *mucRoomConfigListFormItem) initListRemove(onRemove func()) {
	lfi.removeButton.SetSensitive(false)
	lfi.removeButton.SetVisible(false)

	if onRemove != nil {
		lfi.removeButton.Connect("clicked", onRemove)
		lfi.removeButton.SetSensitive(true)
		lfi.removeButton.SetVisible(true)
	}
}

func (lfi *mucRoomConfigListFormItem) initDefaults() {
	lfi.formBox.Add(lfi.form.formView)
}

func (lfi *mucRoomConfigListFormItem) contentBox() gtki.Box {
	return lfi.box
}
