package gotk3extra

// #cgo pkg-config: gdk-3.0 gio-2.0 glib-2.0 gobject-2.0 gtk+-3.0
// #include <stdlib.h>
// #include <gtk/gtk.h>
// #include "size_group.go.h"
// #cgo CFLAGS: -Wno-deprecated-declarations

import "C"
import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func WrapSizeGroupSimple(obj *glib.Object) *gtk.SizeGroup {
	if obj == nil {
		return nil
	}

	return &gtk.SizeGroup{obj}
}

func WrapSizeGroup(obj *glib.Object) (*gtk.SizeGroup, error) {
	if obj == nil {
		return nil, nilPtrErr
	}

	return WrapSizeGroupSimple(obj), nil
}
