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
	isValid() bool
	showValidationErrors()
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

// isValid implements the hasRoomConfigFormField interface
func (f *roomConfigFormField) isValid() bool {
	return true
}

// showValidationErrors implements the hasRoomConfigFormField interface
func (f *roomConfigFormField) showValidationErrors() {}

var (
	errRoomConfigFieldNotSupported = errors.New("room configuration form field not supported")
)

func roomConfigFormUnknownFieldFactory(fieldInfo roomConfigFieldTextInfo, fieldTypeValue interface{}) (hasRoomConfigFormField, error) {
	switch valueHandler := fieldTypeValue.(type) {
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

func roomConfigFormFieldFactory(fieldType muc.RoomConfigFieldType, fieldInfo roomConfigFieldTextInfo, fieldTypeValue interface{}) (hasRoomConfigFormField, error) {
	switch valueHandler := fieldTypeValue.(type) {
	case *muc.RoomConfigFieldTextValue:
		return roomConfigFormTextFieldFactory(fieldType, valueHandler)
	case *muc.RoomConfigFieldTextMultiValue:
		return newRoomConfigFormTextMulti(fieldInfo, valueHandler), nil
	case *muc.RoomConfigFieldBooleanValue:
		return newRoomConfigFormFieldBoolean(fieldInfo, valueHandler), nil
	case *muc.RoomConfigFieldListValue:
		return roomConfigFormListFieldFactory(fieldType, valueHandler)
	case *muc.RoomConfigFieldListMultiValue:
		return newRoomConfigFieldListMulti(fieldInfo, valueHandler), nil
	}

	return nil, errRoomConfigFieldNotSupported
}

func roomConfigFormTextFieldFactory(ft muc.RoomConfigFieldType, value *muc.RoomConfigFieldTextValue) (hasRoomConfigFormField, error) {
	switch ft {
	case muc.RoomConfigFieldLanguage:
		return newRoomConfigFormFieldLanguage(roomConfigFieldsTexts[ft], value), nil
	case muc.RoomConfigFieldPassword:
		return newRoomConfigFormFieldPassword(roomConfigFieldsTexts[ft], value), nil
	}
	return newRoomConfigFormTextField(roomConfigFieldsTexts[ft], value), nil
}

func roomConfigFormListFieldFactory(ft muc.RoomConfigFieldType, value *muc.RoomConfigFieldListValue) (hasRoomConfigFormField, error) {
	if ft == muc.RoomConfigFieldMaxOccupantsNumber || ft == muc.RoomConfigFieldMaxHistoryFetch {
		return newRoomConfigFormFieldListEntry(roomConfigFieldsTexts[ft], value), nil
	}
	return newRoomConfigFormFieldList(roomConfigFieldsTexts[ft], value), nil
}
