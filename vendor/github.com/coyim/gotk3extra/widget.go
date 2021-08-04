package gotk3extra

// #cgo pkg-config: gdk-3.0 gio-2.0 glib-2.0 gobject-2.0 gtk+-3.0
// #include <stdlib.h>
// #include <gtk/gtk.h>
// #include "widget.go.h"
// #cgo CFLAGS: -Wno-deprecated-declarations
import "C"
import (
	"unsafe"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

// WrapWidget can be used to create a GTK widget instance from the given object.
func WrapWidget(obj *glib.Object) *gtk.Widget {
	if obj == nil {
		return nil
	}

	return &gtk.Widget{glib.InitiallyUnowned{obj}}
}

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

// GetWidgetBuildableName is a wrapper around GetBuildableName().
// It can be used to get the "buildable id" of any widget.
func GetWidgetBuildableName(v *gtk.Widget) (string, error) {
	if v == nil || v.Object == nil {
		return "", nilPtrErr
	}
	return GetBuildableName(v.Object)
}

// GetWidgetTemplateChild is a wrapper around gtk_widget_get_template_child().
// This will only report children which were previously declared with
// gtk_widget_class_bind_template_child_full() or one of its variants.
// In other case, it will return an error.
func GetWidgetTemplateChild(vv gtk.IWidget, name string) (*glib.Object, error) {
	if vv == nil {
		return nil, nilPtrErr
	}

	v := vv.ToWidget()
	if v.Object == nil {
		return nil, nilPtrErr
	}

	widget := nativeWidget(v)
	gtype := C.GType(v.Object.TypeFromInstance())
	c := C.gtk_widget_get_template_child(widget, gtype, (*C.gchar)(C.CString(name)))

	obj := glib.Take(unsafe.Pointer(c))
	if obj == nil {
		return nil, nilPtrErr
	}

	return obj, nil
}
