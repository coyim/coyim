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

func newRoomConfigFormTextField(f *muc.RoomConfigFormField, value *muc.RoomConfigFieldTextValue) hasRoomConfigFormField {
	field := &roomConfigFormFieldText{value: value}
	field.roomConfigFormField = newRoomConfigFormField(f, "MUCRoomConfigFormFieldText")

	panicOnDevError(field.builder.bindObjects(field))

	field.entry.SetText(value.Text())

	return field
}

// collectFieldValue MUST be called from the UI thread
func (f *roomConfigFormFieldText) collectFieldValue() {
	f.value.SetText(getEntryText(f.entry))
}
