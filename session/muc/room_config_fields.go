package muc

import "strconv"

// RoomConfigFormField contains information about a room config form field
type RoomConfigFormField struct {
	Name, Type, Label, Description string
	Value                          interface{}
	Options                        []string
}

// SetValue sets the room configuration form field value from the given parameter
func (f *RoomConfigFormField) SetValue(v interface{}) {
	f.Value = v
}

// GetValue returns value based on the field type
func (f *RoomConfigFormField) GetValue() (values []string) {
	switch f.Type {
	case RoomConfigFieldText, RoomConfigFieldTextPrivate, RoomConfigFieldList:
		values = append(values, f.Value.(string))
	case RoomConfigFieldBoolean:
		values = append(values, strconv.FormatBool(f.Value.(bool)))
	}
	return
}

// ConfigListSingleField description
type ConfigListSingleField interface {
	// UpdateField will update the list fields with the given "value" and "options"
	UpdateField(value string, options []string)
	// UpdateValue updates the current field value
	UpdateValue(value string)
	// Value returns the current list value
	CurrentValue() string
	// Options returns the field options
	Options() []string
}

type configListSingleField struct {
	value   string
	options []string
}

func newConfigListSingleField(o []string) ConfigListSingleField {
	return &configListSingleField{
		options: o,
	}
}

func (cf *configListSingleField) UpdateField(v string, o []string) {
	cf.UpdateValue(v)
	if len(o) != 0 {
		cf.options = o
	}
}

func (cf *configListSingleField) UpdateValue(v string) {
	cf.value = v
}

func (cf *configListSingleField) CurrentValue() string {
	return cf.value
}

func (cf *configListSingleField) Options() []string {
	return cf.options
}

// ConfigListMultiField description
type ConfigListMultiField interface {
	// UpdateField will update the list fields with the given "values" and "options"
	UpdateField(values, options []string)
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
