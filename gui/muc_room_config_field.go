package gui

import (
	"errors"
	"strconv"

	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/gtki"
)

type focusable interface {
	HasFocus() bool
	GrabFocus()
}

type hasRoomConfigFormField interface {
	fieldKey() muc.RoomConfigFieldType
	fieldWidget() gtki.Widget
	refreshContent()
	updateFieldValue()
	isValid() bool
	showValidationErrors()
	focusWidget() focusable
}

const fieldErrorIconName = "field_error"

type roomConfigFormField struct {
	field                  muc.RoomConfigFieldType
	builder                *builder
	onShowValidationErrors func()
	onHideValidationErrors func()

	widget      gtki.Box   `gtk-widget:"room-config-field-box"`
	icon        gtki.Image `gtk-widget:"icon-image"`
	label       gtki.Label `gtk-widget:"room-config-field-label"`
	description gtki.Label `gtk-widget:"room-config-field-description"`
}

func newRoomConfigFormField(ft muc.RoomConfigFieldType, fieldInfo roomConfigFieldTextInfo, template string, onShowValidationErrors func(), onHideValidationErrors func()) *roomConfigFormField {
	field := &roomConfigFormField{
		field:                  ft,
		onShowValidationErrors: onShowValidationErrors,
		onHideValidationErrors: onHideValidationErrors,
	}

	field.builder = newBuilder(template)
	panicOnDevError(field.builder.bindObjects(field))

	field.icon.SetFromPixbuf(getMUCIconPixbuf(fieldErrorIconName))

	field.label.SetText(fieldInfo.displayLabel)
	mucStyles.setErrorLabelClass(field.label)

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

// updateFieldValue implements the hasRoomConfigFormField interface
func (f *roomConfigFormField) updateFieldValue() {}

// isValid implements the hasRoomConfigFormField interface
func (f *roomConfigFormField) isValid() bool {
	return true
}

// showValidationErrors implements the hasRoomConfigFormField interface
func (f *roomConfigFormField) showValidationErrors() {
	f.onShowValidationErrors()

	sc, _ := f.label.GetStyleContext()
	sc.AddClass("label-error")

	f.icon.Show()
}

// fieldKey implements the hasRoomConfigFormField interface
func (f *roomConfigFormField) fieldKey() muc.RoomConfigFieldType {
	return f.field
}

// focusWidget implements the hasRoomConfigFormField interface
func (f *roomConfigFormField) focusWidget() focusable {
	return f.widget
}

// hideValidationErrors MUST be called from the UI thread
func (f *roomConfigFormField) hideValidationErrors() {
	f.onHideValidationErrors()

	sc, _ := f.label.GetStyleContext()
	sc.RemoveClass("label-error")

	f.icon.Hide()
}

var (
	errRoomConfigFieldNotSupported = errors.New("room configuration form field not supported")
)

func roomConfigFormUnknownFieldByType(fieldInfo roomConfigFieldTextInfo, fieldTypeValue interface{}, onShowValidationErrors func(), onHideValidationErrors func()) (hasRoomConfigFormField, error) {
	switch valueHandler := fieldTypeValue.(type) {
	case *muc.RoomConfigFieldTextValue:
		return newRoomConfigFormTextField(muc.RoomConfigFieldUnexpected, fieldInfo, valueHandler, onShowValidationErrors, onHideValidationErrors), nil
	case *muc.RoomConfigFieldTextMultiValue:
		return newRoomConfigFormTextMulti(muc.RoomConfigFieldUnexpected, fieldInfo, valueHandler, onShowValidationErrors, onHideValidationErrors), nil
	case *muc.RoomConfigFieldBooleanValue:
		return newRoomConfigFormFieldBoolean(muc.RoomConfigFieldUnexpected, fieldInfo, valueHandler, onShowValidationErrors, onHideValidationErrors), nil
	case *muc.RoomConfigFieldListValue:
		return newRoomConfigFormFieldList(muc.RoomConfigFieldUnexpected, fieldInfo, valueHandler, onShowValidationErrors, onHideValidationErrors), nil
	case *muc.RoomConfigFieldListMultiValue:
		return newRoomConfigFieldListMulti(muc.RoomConfigFieldUnexpected, fieldInfo, valueHandler, onShowValidationErrors, onHideValidationErrors), nil
	}

	return nil, errRoomConfigFieldNotSupported
}

func roomConfigFormFieldByType(fieldType muc.RoomConfigFieldType, fieldInfo roomConfigFieldTextInfo, fieldTypeValue interface{}, onShowValidationErrors func(), onHideValidationErrors func()) (hasRoomConfigFormField, error) {
	switch valueHandler := fieldTypeValue.(type) {
	case *muc.RoomConfigFieldTextValue:
		return roomConfigFormTextFieldByType(fieldType, valueHandler, onShowValidationErrors, onHideValidationErrors)
	case *muc.RoomConfigFieldTextMultiValue:
		return newRoomConfigFormTextMulti(fieldType, fieldInfo, valueHandler, onShowValidationErrors, onHideValidationErrors), nil
	case *muc.RoomConfigFieldBooleanValue:
		return newRoomConfigFormFieldBoolean(fieldType, fieldInfo, valueHandler, onShowValidationErrors, onHideValidationErrors), nil
	case *muc.RoomConfigFieldListValue:
		return roomConfigFormListFieldByType(fieldType, valueHandler, onShowValidationErrors, onHideValidationErrors)
	case *muc.RoomConfigFieldListMultiValue:
		return newRoomConfigFieldListMulti(fieldType, fieldInfo, valueHandler, onShowValidationErrors, onHideValidationErrors), nil
	}

	return nil, errRoomConfigFieldNotSupported
}

func roomConfigFormTextFieldByType(ft muc.RoomConfigFieldType, value *muc.RoomConfigFieldTextValue, onShowValidationErrors func(), onHideValidationErrors func()) (hasRoomConfigFormField, error) {
	switch ft {
	case muc.RoomConfigFieldDescription:
		return newRoomConfigFieldDescription(ft, roomConfigFieldsTexts[ft], value, onShowValidationErrors, onHideValidationErrors), nil
	case muc.RoomConfigFieldLanguage:
		return newRoomConfigFormFieldLanguage(ft, roomConfigFieldsTexts[ft], value, onShowValidationErrors, onHideValidationErrors), nil
	case muc.RoomConfigFieldPassword:
		return newRoomConfigFormFieldPassword(ft, roomConfigFieldsTexts[ft], value, onShowValidationErrors, onHideValidationErrors), nil
	}
	return newRoomConfigFormTextField(ft, roomConfigFieldsTexts[ft], value, onShowValidationErrors, onHideValidationErrors), nil
}

func roomConfigFormListFieldByType(ft muc.RoomConfigFieldType, value *muc.RoomConfigFieldListValue, onShowValidationErrors func(), onHideValidationErrors func()) (hasRoomConfigFormField, error) {
	if ft == muc.RoomConfigFieldMaxOccupantsNumber || ft == muc.RoomConfigFieldMaxHistoryFetch {
		return newRoomConfigFormFieldListEntry(ft, roomConfigFieldsTexts[ft], value, onShowValidationErrors, onHideValidationErrors), nil
	}
	return newRoomConfigFormFieldList(ft, roomConfigFieldsTexts[ft], value, onShowValidationErrors, onHideValidationErrors), nil
}

var numberValidator = func(v string) bool {
	_, err := strconv.Atoi(v)
	return err == nil
}

var roomConfigFieldValidator = map[muc.RoomConfigFieldType]func(string) bool{
	muc.RoomConfigFieldMaxOccupantsNumber:     numberValidator,
	muc.RoomConfigFieldVoiceRequestMinInteval: numberValidator,
}
