package gui

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigFormFieldText struct {
	field *muc.RoomConfigFormField

	widget      gtki.Box   `gtk-widget:"room-config-field-box"`
	label       gtki.Label `gtk-widget:"room-config-field-label"`
	entry       gtki.Entry `gtk-widget:"room-config-text-field-entry"`
	description gtki.Label `gtk-widget:"room-config-field-description"`
}

func newRoomConfigFormTextField(field *muc.RoomConfigFormField) *roomConfigFormFieldText {
	f := &roomConfigFormFieldText{field: field}

	f.initBuilder()
	f.initDefaults()
	return f
}

func (f *roomConfigFormFieldText) initBuilder() {
	builder := newBuilder("MUCRoomConfigFormFieldText")
	panicOnDevError(builder.bindObjects(f))
}

func (f *roomConfigFormFieldText) initDefaults() {
	f.label.SetText(f.field.Label)
	f.entry.SetText(f.field.Value.(string))
	if f.field.Description != "" {
		f.description.SetText(f.field.Description)
	}
}

func (f *roomConfigFormFieldText) fieldWidget() gtki.Widget {
	return f.widget
}

func (f *roomConfigFormFieldText) fieldName() string {
	return f.field.Name
}

func (f *roomConfigFormFieldText) fieldLabel() string {
	return f.field.Label
}

func (f *roomConfigFormFieldText) fieldValue() interface{} {
	return getEntryText(f.entry)
}

func (f *roomConfigFormFieldText) refreshContent() {}
