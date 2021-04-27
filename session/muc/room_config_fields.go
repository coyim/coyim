package muc

import "strconv"

// RoomConfigFormField contains information about a room config form field
type RoomConfigFormField struct {
	Name, Type, Label, Description string
	Value                          interface{}
	Options                        []string
}

// SetValue sets the room configuration form field value from the given parameter
func (f *RoomConfigFormField) SetValue(v interface{}) {
	f.Value = v
}

// GetValue returns value based on the field type
func (f *RoomConfigFormField) GetValue() (values []string) {
	switch f.Type {
	case RoomConfigFieldText, RoomConfigFieldTextPrivate, RoomConfigFieldList:
		values = append(values, f.Value.(string))
	case RoomConfigFieldBoolean:
		values = append(values, strconv.FormatBool(f.Value.(bool)))
	}
	return
}
