package muc

import xmppData "github.com/coyim/coyim/xmpp/data"

// HasRoomConfigFormFieldValue description
type HasRoomConfigFormFieldValue interface {
	Value() []string
	SetValue(interface{})
	Raw() interface{}
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
		value:       roomConfigFormFieldValueFactory(field.Type, field.Values),
	}
}

// SetValue sets the field value with the given "v" parameter
func (f *RoomConfigFormField) SetValue(v interface{}) {
	f.value.SetValue(v)
}

// Value returns the current field value
func (f *RoomConfigFormField) Value() []string {
	return f.value.Value()
}

// RawValue returns the raw value of the field
func (f *RoomConfigFormField) RawValue() interface{} {
	return f.value.Raw()
}

func roomConfigFormFieldValueFactory(typ string, values []string) HasRoomConfigFormFieldValue {
	switch typ {
	case RoomConfigFieldText, RoomConfigFieldTextPrivate:
		return newRoomConfigFieldTextValue(values)
	case RoomConfigFieldTextMulti:
		return newRoomConfigFieldTextMultiValue(values)
	case RoomConfigFieldBoolean:
		return newRoomConfigFieldBooleanValue(values)
	case RoomConfigFieldList:
		return newRoomConfigFieldListValue(values)
	case RoomConfigFieldListMulti:
		return newRoomConfigFieldListMultiValue(values)
	case RoomConfigFieldJidMulti:
		return newRoomConfigFieldJidMultiValue(values)
	}
	return newRoomConfigFieldUnknowValue(values)
}
