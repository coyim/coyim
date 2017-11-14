package xmpp

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/coyim/coyim/xmpp/data"
)

//Unlike xmpp.processForm() we are interested in getting values from forms
//Not rendering them a UI widgets

func parseForm(dstStructElementValue reflect.Value, form data.Form) error {
	elem := dstStructElementValue.Elem()
	dstType := elem.Type()

	for i := 0; i < dstType.NumField(); i++ {
		field := dstType.Field(i)
		tag, ok := field.Tag.Lookup("form-field")
		if !ok {
			continue
		}

		dstValue := elem.Field(i)
		if !dstValue.CanSet() {
			panic("struct element must be settable")
		}

		//Find value in the form
		for _, f := range form.Fields {
			if f.Var != tag {
				continue
			}

			if 0 == len(f.Values) {
				continue
			}

			switch f.Type {
			case "text-single":
				switch dstValue.Kind() {
				case reflect.String:
					dstValue.SetString(f.Values[0])
				case reflect.Int:
					v, _ := strconv.ParseInt(f.Values[0], 0, 64)
					dstValue.SetInt(v)

				}
			default:
				//TODO
				fmt.Printf("I dont know how to set %s fields\n", f.Type)
			}

		}
	}

	return nil
}

func parseForms(dst interface{}, forms []data.Form) error {
	if reflect.TypeOf(dst).Kind() != reflect.Ptr {
		panic("dst must be a pointer")
	}

	dstValue := reflect.ValueOf(dst)
	dstType := dstValue.Elem().Type()
	if dstType.Kind() != reflect.Struct {
		panic("dst must be a struct value")
	}

	for _, f := range forms {
		if f.XMLName.Space != "jabber:x:data" {
			continue
		}

		if 0 == len(f.Fields) || f.Fields[0].Var != "FORM_TYPE" {
			continue
		}

		if 0 == len(f.Fields[0].Values) {
			continue
		}

		formType := f.Fields[0].Values[0]

		for i := 0; i < dstType.NumField(); i++ {
			tag, ok := dstType.Field(i).Tag.Lookup("form-type")
			if !ok {
				continue
			}

			if tag != formType {
				continue
			}
		}

		return parseForm(dstValue, f)
	}

	return nil
}
