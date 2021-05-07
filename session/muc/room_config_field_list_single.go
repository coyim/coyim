package muc

// RoomConfigFieldListValue contains information of the value of the list single field
type RoomConfigFieldListValue struct {
	value   string
	options []*RoomConfigFieldOption
}

func newRoomConfigFieldListValue(values []string, options []*RoomConfigFieldOption) *RoomConfigFieldListValue {
	return &RoomConfigFieldListValue{formFieldSingleString(values), options}
}

// Value implements the "HasRoomConfigFormFieldValue" interface
func (v *RoomConfigFieldListValue) Value() []string {
	return []string{v.value}
}

// SetSelected sets the current selected value
func (v *RoomConfigFieldListValue) SetSelected(s string) {
	v.value = s
}

// Selected returns the current selected value
func (v *RoomConfigFieldListValue) Selected() string {
	return v.value
}

// SelectedOption returns the current option based on the selected value
func (v *RoomConfigFieldListValue) SelectedOption() *RoomConfigFieldOption {
	selected := v.value
	for _, op := range v.options {
		if op.Value == selected {
			return op
		}
	}
	return nil
}

// Options returns the available options for the field
func (v *RoomConfigFieldListValue) Options() []*RoomConfigFieldOption {
	return v.options
}

// SetOptions updates the options for the field
func (v *RoomConfigFieldListValue) SetOptions(options []*RoomConfigFieldOption) {
	v.options = options
}
