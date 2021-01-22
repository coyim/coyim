package gotk3extra

// #cgo pkg-config: gdk-3.0 gio-2.0 glib-2.0 gobject-2.0 gtk+-3.0
// #include <stdlib.h>
// #include <gtk/gtk.h>
// #include "widget.go.h"
// #cgo CFLAGS: -Wno-deprecated-declarations
import "C"
import (
	"unsafe"

	"github.com/gotk3/gotk3/gtk"
)

func nativeWidget(v *gtk.Widget) *C.GtkWidget {
	if v == nil || v.GObject == nil {
		return nil
	}
	p := unsafe.Pointer(v.GObject)
	return C.toGtkWidget(p)
}

// GetParent is a wrapper around gtk_widget_get_parent().
func GetParent(v *gtk.Widget) (gtk.IWidget, error) {
	c := C.gtk_widget_get_parent(nativeWidget(v))
	if c == nil {
		return nil, nilPtrErr
	}
	return CastWidget(c)
}

// GetBuildableName is a wrapper around BuildableGetName().
// It can be used to get the "buildable id" of any widget.
func WidgetGetName(v *gtk.Widget) (string, error) {
	if v == nil || v.Object == nil {
		return "", nilPtrErr
	}
	return BuildableGetName(v.Object)
}
