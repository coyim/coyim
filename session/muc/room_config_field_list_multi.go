package muc

// ConfigListMultiField represents a "list-multi" configuration form field
type ConfigListMultiField interface {
	// UpdateField will update the field with the given "values" and "options"
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
