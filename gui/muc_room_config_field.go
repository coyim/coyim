package gui

import (
	"errors"

	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/gtki"
)

type hasRoomConfigFormField interface {
	fieldWidget() gtki.Widget
	fieldName() string
	fieldLabel() string
	fieldValue() interface{}
}

var (
	errRoomConfigFieldNotSupported = errors.New("room configuration form field not supported")
)

func roomConfigFormFieldFactory(field *muc.RoomConfigFormField) (hasRoomConfigFormField, error) {
	switch field.Type {
	case muc.RoomConfigFieldText:
		return newRoomConfigFormTextField(field), nil
	}

	return nil, errRoomConfigFieldNotSupported
}

type roomConfigFormField struct {
	field *muc.RoomConfigFormField

	widget      gtki.Box   `gtk-widget:"room-config-field-box"`
	label       gtki.Label `gtk-widget:"room-config-field-label"`
	description gtki.Label `gtk-widget:"room-config-field-description"`

	builder *builder
}

func newRoomConfigFormField(field *muc.RoomConfigFormField, template string) *roomConfigFormField {
	f := &roomConfigFormField{
		field: field,
	}

	f.initBuilder(template)
	f.initDefaults()

	return f
}

func (f *roomConfigFormField) initBuilder(template string) {
	f.builder = newBuilder(template)
	panicOnDevError(f.builder.bindObjects(f))
}

func (f *roomConfigFormField) initDefaults() {
	f.label.SetText(f.field.Label)
}

func (f *roomConfigFormField) fieldWidget() gtki.Widget {
	return f.widget
}

func (f *roomConfigFormField) fieldName() string {
	return f.field.Name
}

func (f *roomConfigFormField) fieldLabel() string {
	return f.field.Label
}

func (f *roomConfigFormField) fieldValue() interface{} {
	return f.field.Value
}
