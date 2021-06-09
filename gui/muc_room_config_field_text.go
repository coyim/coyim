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

func newRoomConfigFormTextField(fieldInfo roomConfigFieldTextInfo, value *muc.RoomConfigFieldTextValue) *roomConfigFormFieldText {
	field := &roomConfigFormFieldText{value: value}
	field.roomConfigFormField = newRoomConfigFormField(fieldInfo, "MUCRoomConfigFormFieldText")

	panicOnDevError(field.builder.bindObjects(field))

	field.entry.SetText(value.Text())

	return field
}

// updateFieldValue MUST be called from the UI thread
func (f *roomConfigFormFieldText) updateFieldValue() {
	f.value.SetText(getEntryText(f.entry))
}
