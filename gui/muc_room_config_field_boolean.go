package gui

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigFormFieldBoolean struct {
	*roomConfigFormField
	value *muc.RoomConfigFieldBooleanValue

	grid   gtki.Grid       `gtk-widget:"room-config-field-boolean-grid"`
	field  gtki.ListBoxRow `gtk-widget:"room-config-boolean-field"`
	toggle gtki.Switch     `gtk-widget:"room-config-field-boolean"`
}

func newRoomConfigFormFieldBoolean(fieldType muc.RoomConfigFieldType, fieldInfo roomConfigFieldTextInfo, value *muc.RoomConfigFieldBooleanValue, onShowValidationErrors func(), onHideValidationErrors func()) hasRoomConfigFormField {
	field := &roomConfigFormFieldBoolean{value: value}
	field.roomConfigFormField = newRoomConfigFormField(fieldType, fieldInfo, "MUCRoomConfigFormFieldBoolean", onShowValidationErrors, onHideValidationErrors)

	panicOnDevError(field.builder.bindObjects(field))

	field.toggle.SetActive(value.Boolean())

	return field
}

// refreshContent MUST NOT be called from the UI thread
func (f *roomConfigFormFieldBoolean) refreshContent() {
	doInUIThread(func() {
		f.grid.SetVAlign(gtki.ALIGN_CENTER)
	})
}

// updateFieldValue MUST be called from the UI thread
func (f *roomConfigFormFieldBoolean) updateFieldValue() {
	f.value.SetBoolean(f.toggle.GetActive())
}

// focusWidget implements the hasRoomConfigFormField interface
func (f *roomConfigFormFieldBoolean) focusWidget() gtki.Widget {
	return f.toggle
}
