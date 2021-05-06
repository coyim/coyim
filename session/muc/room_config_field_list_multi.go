package muc

// RoomConfigFieldListMultiValue contains information of the value of the text field
type RoomConfigFieldListMultiValue struct {
	value []string
}

func newRoomConfigFieldListMultiValue(values []string) HasRoomConfigFormFieldValue {
	return &RoomConfigFieldListMultiValue{values}
}

// Value implements the "HasRoomConfigFormFieldValue" interface
func (v *RoomConfigFieldListMultiValue) Value() []string {
	return v.value
}

// SetValue implements the "HasRoomConfigFormFieldValue" interface
func (v *RoomConfigFieldListMultiValue) SetValue(value interface{}) {
	if val, ok := value.([]string); ok {
		v.value = val
	}
}

// Raw implements the "HasRoomConfigFormFieldValue" interface
func (v *RoomConfigFieldListMultiValue) Raw() interface{} {
	return v.value
}
