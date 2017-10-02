package gliba

import (
	"fmt"

	"github.com/gotk3/gotk3/glib"
)

type Wrapper func(interface{}) (interface{}, bool)
type Unwrapper func(interface{}) (interface{}, bool)

var wrappers []Wrapper
var unwrappers []Unwrapper

func AddWrapper(f Wrapper) {
	wrappers = append(wrappers, f)
}

func AddUnwrapper(f Unwrapper) {
	unwrappers = append(unwrappers, f)
}

func WrapAllGuard(v interface{}) interface{} {
	vv, ok := WrapAll(v)
	if !ok {
		panic(fmt.Sprintf("Unrecognized type of object: %#v", v))
	}
	return vv
}

func UnwrapAllGuard(v interface{}) interface{} {
	vv, ok := UnwrapAll(v)
	if !ok {
		panic(fmt.Sprintf("Unrecognized type of object: %#v", v))
	}
	return vv
}

func WrapAll(v interface{}) (interface{}, bool) {
	for _, w := range wrappers {
		v1, ok := w(v)
		if ok {
			return v1, ok
		}
	}
	return nil, false
}

func UnwrapAll(v interface{}) (interface{}, bool) {
	for _, u := range unwrappers {
		v1, ok := u(v)
		if ok {
			return v1, ok
		}
	}
	return nil, false
}

func init() {
	AddWrapper(WrapPrimitive)
	AddWrapper(WrapLocal)

	AddUnwrapper(UnwrapPrimitive)
	AddUnwrapper(UnwrapLocal)
}

func UnwrapPrimitive(v interface{}) (interface{}, bool) {
	if v == nil {
		return nil, true
	}

	switch e := v.(type) {
	case bool:
		return e, true
	case int8:
		return e, true
	case int64:
		return e, true
	case int:
		return e, true
	case uint8:
		return e, true
	case uint64:
		return e, true
	case uint:
		return e, true
	case float32:
		return e, true
	case float64:
		return e, true
	case string:
		return e, true
	}
	return nil, false
}

func WrapPrimitive(v interface{}) (interface{}, bool) {
	return UnwrapPrimitive(v)
}

func Wrap(o interface{}) interface{} {
	v1, ok := WrapLocal(o)
	if !ok {
		panic(fmt.Sprintf("Unrecognized type of object: %#v", o))
	}
	return v1
}

func Unwrap(o interface{}) interface{} {
	v1, ok := UnwrapLocal(o)
	if !ok {
		panic(fmt.Sprintf("Unrecognized type of object: %#v", o))
	}
	return v1
}

func WrapLocal(o interface{}) (interface{}, bool) {
	switch oo := o.(type) {
	case *glib.Application:
		val := WrapApplicationSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *glib.Object:
		val := WrapObjectSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *glib.Signal:
		val := wrapSignalSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *glib.Value:
		val := wrapValueSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	}
	return nil, false
}

func UnwrapLocal(o interface{}) (interface{}, bool) {
	switch oo := o.(type) {
	case *Application:
		val := unwrapApplication(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *Object:
		val := unwrapObject(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *signal:
		val := unwrapSignal(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *value:
		val := unwrapValue(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	}
	return nil, false
}
