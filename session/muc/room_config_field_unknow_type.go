package muc

// RoomConfigFieldUnknownValue contains information of the value of an unknow field type
type RoomConfigFieldUnknownValue struct {
	value []string
}

func newRoomConfigFieldUnknownValue(values []string) *RoomConfigFieldUnknownValue {
	return &RoomConfigFieldUnknownValue{values}
}

// Value implements the "HasRoomConfigFormFieldValue" interface
func (v *RoomConfigFieldUnknownValue) Value() []string {
	return v.value
}

// SetValue sets the current value
func (v *RoomConfigFieldUnknownValue) SetValue(val []string) {
	v.value = val
}
