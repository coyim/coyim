package gui

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigFormTextField struct {
	field *muc.RoomConfigFormField

	widget      gtki.Box   `gtk-widget:"room-config-field-box"`
	label       gtki.Label `gtk-widget:"room-config-field-label"`
	entry       gtki.Entry `gtk-widget:"room-config-text-field-entry"`
	description gtki.Label `gtk-widget:"room-config-field-description"`
}

func newRoomConfigFormTextField(field *muc.RoomConfigFormField) *roomConfigFormTextField {
	f := &roomConfigFormTextField{field: field}

	f.initBuilder()
	f.initDefaults()
	return f
}

func (f *roomConfigFormTextField) initBuilder() {
	builder := newBuilder("MUCRoomConfigFormTextField")
	panicOnDevError(builder.bindObjects(f))
}

func (f *roomConfigFormTextField) initDefaults() {
	f.label.SetText(f.field.Label)
	f.entry.SetText(f.field.Value.(string))
	if f.field.Description != "" {
		f.description.SetText(f.field.Description)
	}
}

func (f *roomConfigFormTextField) fieldWidget() gtki.Widget {
	return f.widget
}

func (f *roomConfigFormTextField) fieldName() string {
	return f.field.Name
}

func (f *roomConfigFormTextField) fieldLabel() string {
	return f.field.Label
}

func (f *roomConfigFormTextField) fieldValue() interface{} {
	return f.field.Value
}
