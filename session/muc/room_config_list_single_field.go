package muc

// ConfigListSingleField represents a "list-single" configuration form field
type ConfigListSingleField interface {
	// UpdateField will update the field with the given "value" and "options"
	UpdateField(value string, options []string)
	// UpdateValue updates the field value
	UpdateValue(value string)
	// Value returns the field value
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
