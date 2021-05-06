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

// SetSelected sets the current selected values
func (v *RoomConfigFieldListMultiValue) SetSelected(s []string) {
	v.value = s
}

// Selected returns the current selected values
func (v *RoomConfigFieldListMultiValue) Selected() []string {
	return v.Value()
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
