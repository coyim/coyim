package gliba

import (
	"reflect"

	"github.com/gotk3/gotk3/glib"
	"github.com/coyim/gotk3adapter/glibi"
)

type Object struct {
	*glib.Object
}

func WrapObjectSimple(v *glib.Object) *Object {
	if v == nil {
		return nil
	}
	return &Object{v}
}

func unwrapObject(v glibi.Object) *glib.Object {
	if v == nil {
		return nil
	}
	return v.(*Object).Object
}

func FixupArray(v []interface{}) []interface{} {
	nv := make([]interface{}, len(v))
	for ix, vv := range v {
		nv[ix] = UnwrapAllGuard(vv)
	}
	return nv
}

func fixupReturnValue(v []reflect.Value) interface{} {
	return UnwrapAllGuard(v[0].Interface())
}

func fixupArg(tv reflect.Type, v interface{}) reflect.Value {
	vvt := reflect.TypeOf(v)

	switch vvt.Kind() {
	case reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32,
		reflect.Float64,
		reflect.Complex64,
		reflect.Complex128:
		if vvt != tv {
			return reflect.ValueOf(v).Convert(tv)
		} else {
			return reflect.ValueOf(v)
		}
	default:
		return reflect.ValueOf(WrapAllGuard(v))
	}
}

func fixupArgs(t reflect.Type, v ...interface{}) []reflect.Value {
	res := make([]reflect.Value, len(v))
	for ix, vv := range v {
		res[ix] = fixupArg(t.In(ix), vv)
	}
	return res
}

func FixupFunction(v interface{}) interface{} {
	rf := reflect.ValueOf(v)
	if rf.Type().Kind() != reflect.Func {
		panic("We can't fix up something that is not a function")
	}

	ni := rf.Type().NumIn()
	no := rf.Type().NumOut()

	if ni > 4 {
		panic("We can't handle more than 4 arguments to a closure")
	}

	if no > 1 {
		panic("We can't handle more than 1 output arguments for a closure")
	}

	switch ni {
	case 0:
		if no == 0 {
			return v
		} else {
			return func() interface{} {
				return fixupReturnValue(rf.Call([]reflect.Value{}))
			}
		}
	case 1:
		if no == 0 {
			return func(v1 interface{}) {
				rf.Call(fixupArgs(rf.Type(), v1))
			}
		} else {
			return func(v1 interface{}) interface{} {
				return fixupReturnValue(rf.Call(fixupArgs(rf.Type(), v1)))
			}
		}
	case 2:
		if no == 0 {
			return func(v1, v2 interface{}) {
				rf.Call(fixupArgs(rf.Type(), v1, v2))
			}
		} else {
			return func(v1, v2 interface{}) interface{} {
				return fixupReturnValue(rf.Call(fixupArgs(rf.Type(), v1, v2)))
			}
		}
	case 3:
		if no == 0 {
			return func(v1, v2, v3 interface{}) {
				rf.Call(fixupArgs(rf.Type(), v1, v2, v3))
			}
		} else {
			return func(v1, v2, v3 interface{}) interface{} {
				return fixupReturnValue(rf.Call(fixupArgs(rf.Type(), v1, v2, v3)))
			}
		}
	case 4:
		if no == 0 {
			return func(v1, v2, v3, v4 interface{}) {
				rf.Call(fixupArgs(rf.Type(), v1, v2, v3, v4))
			}
		} else {
			return func(v1, v2, v3, v4 interface{}) interface{} {
				return fixupReturnValue(rf.Call(fixupArgs(rf.Type(), v1, v2, v3, v4)))
			}
		}
	}

	panic("Shouldn't happen")
}

func (v *Object) Connect(v1 string, v2 interface{}, v3 ...interface{}) (glibi.SignalHandle, error) {
	nv2 := FixupFunction(v2)
	vx1, vx2 := v.Object.Connect(v1, nv2, FixupArray(v3)...)
	return glibi.SignalHandle(vx1), vx2
}

func (v *Object) ConnectAfter(v1 string, v2 interface{}, v3 ...interface{}) (glibi.SignalHandle, error) {
	nv2 := FixupFunction(v2)
	vx1, vx2 := v.Object.ConnectAfter(v1, nv2, FixupArray(v3)...)
	return glibi.SignalHandle(vx1), vx2
}

func (v *Object) Emit(v1 string, v2 ...interface{}) (interface{}, error) {
	vx1, vx2 := v.Object.Emit(v1, FixupArray(v2)...)
	return WrapAllGuard(vx1), vx2
}

func (v *Object) GetProperty(v1 string) (interface{}, error) {
	vx1, vx2 := v.Object.GetProperty(v1)
	return WrapAllGuard(vx1), vx2
}

func (v *Object) SetProperty(v1 string, v2 interface{}) error {
	return v.Object.SetProperty(v1, WrapAllGuard(v2))
}

func (v *Object) Ref() {
	v.Object.Ref()
}

func (v *Object) Unref() {
	v.Object.Unref()
}
