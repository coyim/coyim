package muc

// RoomConfigFieldTextValue contains information of the value of the text field
type RoomConfigFieldTextValue struct {
	value string
}

func newRoomConfigFieldTextValue(values []string) *RoomConfigFieldTextValue {
	return &RoomConfigFieldTextValue{formFieldSingleString(values)}
}

// Value implements the "HasRoomConfigFormFieldValue" interface
func (v *RoomConfigFieldTextValue) Value() []string {
	return []string{v.value}
}

// SetValue implements the "HasRoomConfigFormFieldValue" interface
func (v *RoomConfigFieldTextValue) SetValue(value interface{}) {
	if val, ok := value.(string); ok {
		v.value = val
	}
}

// Text returns the current text value
func (v *RoomConfigFieldTextValue) Text() string {
	return v.value
}
