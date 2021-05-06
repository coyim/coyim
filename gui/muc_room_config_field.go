package gui

import (
	"errors"

	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/gtki"
)

type hasRoomConfigFormField interface {
	fieldName() string
	fieldLabel() string
	fieldDescription() string
	fieldValue() interface{}
	fieldWidget() gtki.Widget
	refreshContent()
	collectFieldValue()
}

type roomConfigFormField struct {
	field   *muc.RoomConfigFormField
	builder *builder

	widget      gtki.Box   `gtk-widget:"room-config-field-box"`
	label       gtki.Label `gtk-widget:"room-config-field-label"`
	description gtki.Label `gtk-widget:"room-config-field-description"`
}

func newRoomConfigFormField(f *muc.RoomConfigFormField, template string) *roomConfigFormField {
	field := &roomConfigFormField{
		field: f,
	}

	field.builder = newBuilder(template)
	panicOnDevError(field.builder.bindObjects(field))

	field.label.SetText(field.fieldLabel())

	description := field.fieldDescription()
	field.description.SetText(description)
	if description == "" {
		field.description.Hide()
	}

	return field
}

// fieldName implements the hasRoomConfigFormField interface
func (f *roomConfigFormField) fieldName() string {
	return f.field.Name
}

// fieldLabel implements the hasRoomConfigFormField interface
func (f *roomConfigFormField) fieldLabel() string {
	return f.field.Label
}

// fieldDescription implements the hasRoomConfigFormField interface
func (f *roomConfigFormField) fieldDescription() string {
	return f.field.Description
}

// fieldValue implements the hasRoomConfigFormField interface
func (f *roomConfigFormField) fieldValue() interface{} {
	return nil
}

// fieldWidget implements the hasRoomConfigFormField interface
func (f *roomConfigFormField) fieldWidget() gtki.Widget {
	return f.widget
}

// refreshContent implements the hasRoomConfigFormField interface
func (f *roomConfigFormField) refreshContent() {}

// collectFieldValue implements the hasRoomConfigFormField interface
func (f *roomConfigFormField) collectFieldValue() {}

var (
	errRoomConfigFieldNotSupported = errors.New("room configuration form field not supported")
)

func roomConfigFormFieldFactory(field *muc.RoomConfigFormField) (hasRoomConfigFormField, error) {
	switch valueHandler := field.ValueType().(type) {
	case *muc.RoomConfigFieldTextValue:
		return newRoomConfigFormTextField(field, valueHandler), nil
	case *muc.RoomConfigFieldBooleanValue:
		return newRoomConfigFormFieldBoolean(field, valueHandler), nil
	case *muc.RoomConfigFieldListValue:
		return newRoomConfigFormFieldList(field, valueHandler), nil
	}

	return nil, errRoomConfigFieldNotSupported
}
