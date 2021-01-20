package gotk3extra

// #cgo pkg-config: gdk-3.0 gio-2.0 glib-2.0 gobject-2.0 gtk+-3.0
// #include <stdlib.h>
// #include <gtk/gtk.h>
// #include "casting.go.h"
// #include "widget.go.h"
// #cgo CFLAGS: -Wno-deprecated-declarations
import "C"
import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

// Code in this file is lightly adapted from gotk3 to make it available for use in this package

// CastInternal casts the given object to the appropriate Go struct, but returns it as interface for later type assertions.
// The className is the results of C.object_get_class_name(c) called on the native object.
// The obj is the result of glib.Take(unsafe.Pointer(c)), used as a parameter for the wrapper functions.
func CastInternal(className string, obj *glib.Object) (interface{}, error) {
	fn, ok := gtk.WrapMap[className]
	if !ok {
		return nil, errors.New("unrecognized class name '" + className + "'")
	}

	// Check that the wrapper function is actually a function
	rf := reflect.ValueOf(fn)
	if rf.Type().Kind() != reflect.Func {
		return nil, errors.New("wraper is not a function")
	}

	// Call the wraper function with the *glib.Object as first parameter
	// e.g. "wrapWindow(obj)"
	v := reflect.ValueOf(obj)
	rv := rf.Call([]reflect.Value{v})

	// At most/max 1 return value
	if len(rv) != 1 {
		return nil, errors.New("wrapper did not return")
	}

	// Needs to be a pointer of some sort
	if k := rv[0].Kind(); k != reflect.Ptr {
		return nil, fmt.Errorf("wrong return type %s", k)
	}

	// Only get an interface value, type check will be done in more specific functions
	return rv[0].Interface(), nil
}

// CastX takes a native GObject and casts it to the appropriate Go struct.
func CastX(c unsafe.Pointer) (glib.IObject, error) {
	return Cast(C.toGObject(c))
}

// Cast takes a native GObject and casts it to the appropriate Go struct.
func Cast(c *C.GObject) (glib.IObject, error) {
	var (
		className = goString(C.object_get_class_name(c))
		obj       = glib.Take(unsafe.Pointer(c))
	)

	intf, err := CastInternal(className, obj)
	if err != nil {
		return nil, err
	}

	ret, ok := intf.(glib.IObject)
	if !ok {
		return nil, errors.New("did not return an IObject")
	}

	return ret, nil
}

func goString(cstr *C.gchar) string {
	return C.GoString((*C.char)(cstr))
}

// CastWidgetX takes a native GtkWidget and casts it to the appropriate Go struct.
func CastWidgetX(c unsafe.Pointer) (gtk.IWidget, error) {
	return CastWidget(C.toGtkWidget(c))
}

// CastWidget takes a native GtkWidget and casts it to the appropriate Go struct.
func CastWidget(c *C.GtkWidget) (gtk.IWidget, error) {
	var (
		className = goString((C.object_get_class_name(C.toGObject(unsafe.Pointer(c)))))
		obj       = glib.Take(unsafe.Pointer(c))
	)

	intf, err := CastInternal(className, obj)
	if err != nil {
		return nil, err
	}

	ret, ok := intf.(gtk.IWidget)
	if !ok {
		return nil, errors.New("did not return an IWidget")
	}

	return ret, nil
}
