package muc

// RoomConfigFieldListMultiValue contains information of the value of the text field
type RoomConfigFieldListMultiValue struct {
	value   []string
	options []string
}

func newRoomConfigFieldListMultiValue(values, options []string) *RoomConfigFieldListMultiValue {
	return &RoomConfigFieldListMultiValue{values, options}
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

// Options return the list of options from where the value can be taken
func (v *RoomConfigFieldListMultiValue) Options() []string {
	return v.options
}

// SetOptions update the list of options from where the value can be taken, only if
// the given list is not empty
func (v *RoomConfigFieldListMultiValue) SetOptions(options []string) {
	if len(options) > 0 {
		v.options = options
	}
}
