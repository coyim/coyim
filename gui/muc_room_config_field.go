package gui

import (
	"errors"

	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/gtki"
)

type hasRoomConfigFormField interface {
	fieldWidget() gtki.Widget
	refreshContent()
	collectFieldValue()
}

type roomConfigFormField struct {
	builder *builder

	widget      gtki.Box   `gtk-widget:"room-config-field-box"`
	label       gtki.Label `gtk-widget:"room-config-field-label"`
	description gtki.Label `gtk-widget:"room-config-field-description"`
}

func newRoomConfigFormField(fieldInfo roomConfigFieldTextInfo, template string) *roomConfigFormField {
	field := &roomConfigFormField{}

	field.builder = newBuilder(template)
	panicOnDevError(field.builder.bindObjects(field))

	field.label.SetText(fieldInfo.displayLabel)

	description := fieldInfo.displayDescription
	field.description.SetText(description)
	if description == "" {
		field.description.Hide()
	}

	return field
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

func roomConfigFormFieldFactory(fieldInfo roomConfigFieldTextInfo, fieldType interface{}) (hasRoomConfigFormField, error) {
	switch valueHandler := fieldType.(type) {
	case *muc.RoomConfigFieldTextValue:
		return newRoomConfigFormTextField(fieldInfo, valueHandler), nil
	case *muc.RoomConfigFieldTextMultiValue:
		return newRoomConfigFormTextMulti(fieldInfo, valueHandler), nil
	case *muc.RoomConfigFieldBooleanValue:
		return newRoomConfigFormFieldBoolean(fieldInfo, valueHandler), nil
	case *muc.RoomConfigFieldListValue:
		return newRoomConfigFormFieldList(fieldInfo, valueHandler), nil
	case *muc.RoomConfigFieldListMultiValue:
		return newRoomConfigFieldListMulti(fieldInfo, valueHandler), nil
	}

	return nil, errRoomConfigFieldNotSupported
}
