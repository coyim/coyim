package muc

// RoomConfigFieldTextMultiValue contains information of the value of the boolean field
type RoomConfigFieldTextMultiValue struct {
	value []string
}

func newRoomConfigFieldTextMultiValue(values []string) HasRoomConfigFormFieldValue {
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

// Raw implements the "HasRoomConfigFormFieldValue" interface
func (v *RoomConfigFieldTextMultiValue) Raw() interface{} {
	return v.value
}
