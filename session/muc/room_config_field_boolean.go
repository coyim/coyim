package muc

import "strconv"

// RoomConfigFieldBooleanValue contains information of the value of the boolean field
type RoomConfigFieldBooleanValue struct {
	value bool
}

func newRoomConfigFieldBooleanValue(values []string) *RoomConfigFieldBooleanValue {
	return &RoomConfigFieldBooleanValue{formFieldBool(values)}
}

// Value implements the "HasRoomConfigFormFieldValue" interface
func (v *RoomConfigFieldBooleanValue) Value() []string {
	return []string{strconv.FormatBool(v.value)}
}

// SetValue implements the "HasRoomConfigFormFieldValue" interface
func (v *RoomConfigFieldBooleanValue) SetValue(value interface{}) {
	if val, ok := value.(bool); ok {
		v.value = val
	}
}

// Boolean returns the current boolean value
func (v *RoomConfigFieldBooleanValue) Boolean() bool {
	return v.value
}
