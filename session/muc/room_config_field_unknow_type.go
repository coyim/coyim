package muc

// RoomConfigFieldUnknowValue contains information of the value of an unknow field type
type RoomConfigFieldUnknowValue struct {
	value []string
}

func newRoomConfigFieldUnknowValue(values []string) *RoomConfigFieldUnknowValue {
	return &RoomConfigFieldUnknowValue{values}
}

// Value implements the "HasRoomConfigFormFieldValue" interface
func (v *RoomConfigFieldUnknowValue) Value() []string {
	return v.value
}

// SetValue sets the current value
func (v *RoomConfigFieldUnknowValue) SetValue(val []string) {
	v.value = val
}
