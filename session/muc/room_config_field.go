package muc

// HasRoomConfigFormField represents a configuration form field
type HasRoomConfigFormField interface {
	Name() string
	Type() string
	Label() string
	Description() string
	Value() []string
	SetValue(interface{})
}

// HasRoomConfigFormFieldOptions represents a configuration form field that has options
type HasRoomConfigFormFieldOptions interface {
	Options() []string
	SetOptions([]string)
}

type roomConfigFormField struct {
	name        string
	typ         string
	label       string
	description string
	value       interface{}
}

func newRoomConfigFormField(name, typ, label, description string) *roomConfigFormField {
	return &roomConfigFormField{
		name:        name,
		typ:         typ,
		label:       label,
		description: description,
	}
}

// Name implements the HasRoomConfigFormField interface
func (f *roomConfigFormField) Name() string {
	return f.name
}

// Type implements the HasRoomConfigFormField interface
func (f *roomConfigFormField) Type() string {
	return f.typ
}

// Label implements the HasRoomConfigFormField interface
func (f *roomConfigFormField) Label() string {
	return f.label
}

// Description implements the HasRoomConfigFormField interface
func (f *roomConfigFormField) Description() string {
	return f.description
}

// Value implements the HasRoomConfigFormField interface
func (f *roomConfigFormField) Value() []string {
	return []string{f.value.(string)}
}

// SetValue implements the HasRoomConfigFormField interface
func (f *roomConfigFormField) SetValue(v interface{}) {
	f.value = v
}
