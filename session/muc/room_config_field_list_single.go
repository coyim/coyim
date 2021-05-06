package muc

// RoomConfigFieldListValue contains information of the value of the list single field
type RoomConfigFieldListValue struct {
	value   string
	options []string
}

func newRoomConfigFieldListValue(values, options []string) *RoomConfigFieldListValue {
	return &RoomConfigFieldListValue{formFieldSingleString(values), options}
}

// Value implements the "HasRoomConfigFormFieldValue" interface
func (v *RoomConfigFieldListValue) Value() []string {
	return []string{v.value}
}

// SetValue implements the "HasRoomConfigFormFieldValue" interface
func (v *RoomConfigFieldListValue) SetValue(value interface{}) {
	if val, ok := value.(string); ok {
		v.value = val
	}
}

// Raw implements the "HasRoomConfigFormFieldValue" interface
func (v *RoomConfigFieldListValue) Raw() interface{} {
	return v.value
}

// Options returns the available options for the field
func (v *RoomConfigFieldListValue) Options() []string {
	return v.options
}
