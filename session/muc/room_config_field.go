package muc

import xmppData "github.com/coyim/coyim/xmpp/data"

// HasRoomConfigFormFieldValue description
type HasRoomConfigFormFieldValue interface {
	Value() []string
}

// RoomConfigFormField contains information of the field from the configuration form
type RoomConfigFormField struct {
	Name        string
	Type        string
	Label       string
	Description string
	value       HasRoomConfigFormFieldValue
}

func newRoomConfigFormField(field xmppData.FormFieldX) *RoomConfigFormField {
	return &RoomConfigFormField{
		Name:        field.Var,
		Type:        field.Type,
		Label:       field.Label,
		Description: field.Desc,
		value:       roomConfigFormFieldValueFactory(field),
	}
}

// Value returns the current field value
func (f *RoomConfigFormField) Value() []string {
	return f.value.Value()
}

// RawValue returns the raw value of the field
func (f *RoomConfigFormField) RawValue() interface{} {
	switch v := f.value.(type) {
	case *RoomConfigFieldTextValue:
		return v.Text()
	case *RoomConfigFieldTextMultiValue:
		return v.Text()
	case *RoomConfigFieldBooleanValue:
		return v.Boolean()
	case *RoomConfigFieldListValue:
		return v.Selected()
	case *RoomConfigFieldListMultiValue:
		return v.Selected()
	case *RoomConfigFieldJidMultiValue:
		return v.List()
	}

	return f.value.Value()
}

// ValueType returns the value type handler of the field
func (f *RoomConfigFormField) ValueType() HasRoomConfigFormFieldValue {
	return f.value
}

func (f *RoomConfigFormField) updateBooleanValue(v bool) {
	if field, ok := f.value.(*RoomConfigFieldBooleanValue); ok {
		field.SetBoolean(v)
	}
}

func roomConfigFormFieldValueFactory(field xmppData.FormFieldX) HasRoomConfigFormFieldValue {
	values := field.Values
	options := formFieldOptionsValues(field.Options)
	fieldType := field.Type
	if field.Var == configFieldMaxHistoryFetch || field.Var == configFieldMaxHistoryLength {
		fieldType = RoomConfigFieldList
		options = maxHistoryFetchDefaultOptions
	}

	switch fieldType {
	case RoomConfigFieldText, RoomConfigFieldTextPrivate:
		return newRoomConfigFieldTextValue(values)
	case RoomConfigFieldTextMulti:
		return newRoomConfigFieldTextMultiValue(values)
	case RoomConfigFieldBoolean:
		return newRoomConfigFieldBooleanValue(values)
	case RoomConfigFieldList:
		return newRoomConfigFieldListValue(values, options)
	case RoomConfigFieldListMulti:
		return newRoomConfigFieldListMultiValue(values, options)
	case RoomConfigFieldJidMulti:
		return newRoomConfigFieldJidMultiValue(values)
	}

	return newRoomConfigFieldUnknownValue(values)
}
