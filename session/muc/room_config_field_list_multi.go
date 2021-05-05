package muc

// ConfigListMultiField represents a "list-multi" configuration form field
type ConfigListMultiField interface {
	// UpdateField will update the field with the given "values" and "options"
	UpdateField(values, options []string)
	Values() []string
}

type configListMultiField struct {
	values  []string
	options []string
}

func newConfigListMultiField(o []string) ConfigListMultiField {
	return &configListMultiField{
		options: o,
	}
}

func (cf *configListMultiField) UpdateField(v, o []string) {
	cf.values = v
	cf.options = o
}

func (cf *configListMultiField) Values() []string {
	return cf.values
}

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
