package muc

// RoomConfigFieldListMultiOption represents the list multi type for the
// config form
type RoomConfigFieldListMultiOption struct {
	Value string
	Label string
}

func newRoomConfigFieldListMultiOption(v, l string) *RoomConfigFieldListMultiOption {
	return &RoomConfigFieldListMultiOption{v, l}
}

// RoomConfigFieldListMultiValue contains information of the value of the text field
type RoomConfigFieldListMultiValue struct {
	value   []string
	options []*RoomConfigFieldOption
}

func newRoomConfigFieldListMultiValue(values []string, options []*RoomConfigFieldOption) *RoomConfigFieldListMultiValue {
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
func (v *RoomConfigFieldListMultiValue) Options() []*RoomConfigFieldOption {
	return v.options
}

// SetOptions update the list of options from where the value can be taken, only if
// the given list is not empty
func (v *RoomConfigFieldListMultiValue) SetOptions(options []*RoomConfigFieldOption) {
	if len(options) > 0 {
		v.options = options
	}
}

// IsSelected returns a boolean that indicates if the given option is selected
func (v *RoomConfigFieldListMultiValue) IsSelected(option *RoomConfigFieldOption) bool {
	for _, o := range v.Selected() {
		if o == option.Value {
			return true
		}
	}
	return false
}
