package gui

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigFormFieldText struct {
	*roomConfigFormField

	entry gtki.Entry `gtk-widget:"room-config-text-field-entry"`
}

func newRoomConfigFormTextField(f muc.HasRoomConfigFormField) hasRoomConfigFormField {
	field := &roomConfigFormFieldText{}
	field.roomConfigFormField = newRoomConfigFormField(f, "MUCRoomConfigFormFieldText")

	panicOnDevError(field.builder.bindObjects(field))

	if text, ok := field.fieldValue().(string); ok {
		field.entry.SetText(text)
	}

	return field
}
