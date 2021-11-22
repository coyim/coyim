package gui

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigFormFieldText struct {
	*roomConfigFormField
	value *muc.RoomConfigFieldTextValue

	entry gtki.Entry `gtk-widget:"room-config-text-field-entry"`
}

func newRoomConfigFormTextField(ft muc.RoomConfigFieldType, fieldInfo roomConfigFieldTextInfo, value *muc.RoomConfigFieldTextValue, onShowValidationErrors func(), onHideValidationErrors func()) *roomConfigFormFieldText {
	field := &roomConfigFormFieldText{value: value}
	field.roomConfigFormField = newRoomConfigFormField(ft, fieldInfo, "MUCRoomConfigFormFieldText", onShowValidationErrors, onHideValidationErrors)

	panicOnDevError(field.builder.bindObjects(field))
	field.builder.ConnectSignals(map[string]interface{}{
		"on_field_entry_change": field.onFieldEntryChanged,
	})

	field.entry.SetText(value.Text())

	return field
}

// onFieldEntryChanged MUST be called from the UI thread
func (f *roomConfigFormFieldText) onFieldEntryChanged() {
	if f.isValid() {
		f.hideValidationErrors()
		return
	}

	f.showValidationErrors()
}

// isValid implements the hasRoomConfigFormField interface
func (f *roomConfigFormFieldText) isValid() bool {
	v := getEntryText(f.entry)

	validator, ok := roomConfigFieldValidator[f.field]
	if ok {
		return validator(v)
	}

	return true
}

// updateFieldValue MUST be called from the UI thread
func (f *roomConfigFormFieldText) updateFieldValue() {
	f.value.SetText(getEntryText(f.entry))
}

// refreshContent implements the hasRoomConfigFormField interface
func (f *roomConfigFormFieldText) refreshContent() {
	doInUIThread(func() {
		f.description.SetVisible(f.description.GetLabel() != "")
	})
}

// focusWidget implements the hasRoomConfigFormField interface
func (f *roomConfigFormFieldText) focusWidget() focusable {
	return f.entry
}
