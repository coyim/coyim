package gui

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigFieldDescription struct {
	*roomConfigFormField
	value *muc.RoomConfigFieldTextValue

	textView gtki.TextView `gtk-widget:"room-config-text-multi-field-textview"`
}

func newRoomConfigFieldDescription(fieldInfo roomConfigFieldTextInfo, value *muc.RoomConfigFieldTextValue) hasRoomConfigFormField {
	field := &roomConfigFieldDescription{value: value}
	field.roomConfigFormField = newRoomConfigFormField(fieldInfo, "MUCRoomConfigFormFieldTextMulti")

	panicOnDevError(field.builder.bindObjects(field))

	tb, _ := g.gtk.TextBufferNew(nil)
	field.textView.SetBuffer(tb)

	tb.SetText(value.Text())

	return field
}

// updateFieldValue MUST be called from the UI thread
func (f *roomConfigFieldDescription) updateFieldValue() {
	f.value.SetText(getTextViewText(f.textView))
}
