package muc

import (
	"fmt"
	"reflect"
	"strconv"
)

// HasRoomConfigFormField represents a configuration form field
type HasRoomConfigFormField interface {
	Name() string
	Type() string
	Label() string
	Description() string
	Value() interface{}
	SetValue(interface{})
	ValueX() []string
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
func (f *roomConfigFormField) Value() interface{} {
	return f.value
}

// ValueX implements the HasRoomConfigFormField interface
func (f *roomConfigFormField) ValueX() []string {
	v := reflect.ValueOf(f.Value())

	switch t := v.Kind(); t {
	case reflect.String:
		return []string{v.String()}
	case reflect.Bool:
		return []string{strconv.FormatBool(v.Bool())}
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
		return []string{strconv.FormatInt(v.Int(), 10)}
	case reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64:
		return []string{strconv.FormatUint(v.Uint(), 10)}
	case reflect.Slice:
		values := []string{}
		if list, ok := v.Interface().([]string); ok {
			for _, itm := range list {
				values = append(values, string(itm))
			}
		}
		return values
	default:
		fmt.Printf("DON'T KNOW ABOUT TYPE: %d\n", t)

	}

	return []string{}
}

// SetValue implements the HasRoomConfigFormField interface
func (f *roomConfigFormField) SetValue(v interface{}) {
	f.value = v
}
