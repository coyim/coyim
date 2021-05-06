package muc

// RoomConfigFieldTextValue contains information of the value of the text field
type RoomConfigFieldTextValue struct {
	value string
}

func newRoomConfigFieldTextValue(values []string) *RoomConfigFieldTextValue {
	return &RoomConfigFieldTextValue{formFieldSingleString(values)}
}

// Value implements the "HasRoomConfigFormFieldValue" interface
func (v *RoomConfigFieldTextValue) Value() []string {
	return []string{v.value}
}

// SetText sets the current text value
func (v *RoomConfigFieldTextValue) SetText(t string) {
	v.value = t
}

// Text returns the current text value
func (v *RoomConfigFieldTextValue) Text() string {
	return v.value
}
