package muc

// RoomConfigFieldTextMultiValue contains information of the value of the boolean field
type RoomConfigFieldTextMultiValue struct {
	value []string
}

func newRoomConfigFieldTextMultiValue(values []string) *RoomConfigFieldTextMultiValue {
	return &RoomConfigFieldTextMultiValue{values}
}

// Value implements the "HasRoomConfigFormFieldValue" interface
func (v *RoomConfigFieldTextMultiValue) Value() []string {
	return v.value
}

// SetValue implements the "HasRoomConfigFormFieldValue" interface
func (v *RoomConfigFieldTextMultiValue) SetValue(value interface{}) {
	if val, ok := value.([]string); ok {
		v.value = val
	}
}

// Text returns the current text (multi line) value
func (v *RoomConfigFieldTextMultiValue) Text() []string {
	return v.Value()
}
