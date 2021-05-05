package muc

import (
	"fmt"
	"reflect"
	"strconv"
)

// HasRoomConfigFormFieldValue description
type HasRoomConfigFormFieldValue interface {
	Value() []string
	SetValue(interface{})
	Raw() interface{}
}

// RoomConfigFormField contains information of the field from the configuration form
type RoomConfigFormField struct {
	Name        string
	Type        string
	Label       string
	Description string
	Value       interface{}
}

func newRoomConfigFormField(name, typ, label, description string) *RoomConfigFormField {
	return &RoomConfigFormField{
		Name:        name,
		Type:        typ,
		Label:       label,
		Description: description,
	}
}

// ValueX implements the HasRoomConfigFormField interface
func (f *RoomConfigFormField) ValueX() []string {
	v := reflect.ValueOf(f.Value)

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
func (f *RoomConfigFormField) SetValue(v interface{}) {
	f.Value = v
}
