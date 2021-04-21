package gui

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigFormFieldBoolean struct {
	field *muc.RoomConfigFormField

	widget      gtki.Box    `gtk-widget:"room-config-field-box"`
	contentGrid gtki.Grid   `gtk-widget:"room-config-field-boolean-grid"`
	label       gtki.Label  `gtk-widget:"room-config-field-label"`
	entry       gtki.Switch `gtk-widget:"room-config-field-boolean"`
	description gtki.Label  `gtk-widget:"room-config-field-description"`
}

func newRoomConfigFormFieldBoolean(field *muc.RoomConfigFormField) *roomConfigFormFieldBoolean {
	f := &roomConfigFormFieldBoolean{field: field}

	f.initBuilder()
	f.initDefaults()
	return f
}

func (f *roomConfigFormFieldBoolean) initBuilder() {
	builder := newBuilder("MUCRoomConfigFormFieldBoolean")
	panicOnDevError(builder.bindObjects(f))
}

func (f *roomConfigFormFieldBoolean) initDefaults() {
	f.label.SetText(f.field.Label)
	f.entry.SetActive(f.field.Value.(bool))
	if f.field.Description != "" {
		f.description.SetText(f.field.Description)
	}
}

func (f *roomConfigFormFieldBoolean) fieldWidget() gtki.Widget {
	return f.widget
}

func (f *roomConfigFormFieldBoolean) fieldName() string {
	return f.field.Name
}

func (f *roomConfigFormFieldBoolean) fieldLabel() string {
	return f.field.Label
}

func (f *roomConfigFormFieldBoolean) fieldValue() interface{} {
	return f.entry.GetActive()
}

func (f *roomConfigFormFieldBoolean) refreshFields() {
	doInUIThread(func() {
		f.contentGrid.SetVAlign(gtki.ALIGN_CENTER)
		f.description.SetVisible(f.field.Description != "")
	})
}
