package gotk3extra

// #cgo pkg-config: gdk-3.0 gio-2.0 glib-2.0 gobject-2.0 gtk+-3.0
// #include <stdlib.h>
// #include <gtk/gtk.h>
// #include "css_class_cell_renderer.go.h"
// #cgo CFLAGS: -Wno-deprecated-declarations
import "C"
import (
	"unsafe"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func init() {
	tm := []glib.TypeMarshaler{
		{glib.Type(C.css_class_cell_renderer_get_type()), marshalCSSClassCellRenderer},
	}
	glib.RegisterGValueMarshalers(tm)

	gtk.WrapMap["CSSClassCellRenderer"] = wrapCSSClassCellRenderer
}

/*
 * CSSClassCellRenderer
 */

// CSSClassCellRenderer is a representation of the native implementation.
type CSSClassCellRenderer struct {
	gtk.CellRenderer
}

// native returns a pointer to the underlying GtkCellRendererText.
func (v *CSSClassCellRenderer) native() *C.CSSClassCellRenderer {
	if v == nil || v.GObject == nil {
		return nil
	}
	p := unsafe.Pointer(v.GObject)
	return C.toCSSClassCellRenderer(p)
}

func marshalCSSClassCellRenderer(p uintptr) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(unsafe.Pointer(p)))
	obj := glib.Take(unsafe.Pointer(c))
	return wrapCSSClassCellRenderer(obj), nil
}

func wrapCSSClassCellRenderer(obj *glib.Object) *CSSClassCellRenderer {
	if obj == nil {
		return nil
	}

	return &CSSClassCellRenderer{gtk.CellRenderer{glib.InitiallyUnowned{obj}}}
}

// CSSClassCellRendererNew is a wrapper around css_class_cell_renderer_new().
func CSSClassCellRendererNew() (*CSSClassCellRenderer, error) {
	c := C.css_class_cell_renderer_new()
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.Take(unsafe.Pointer(c))
	return wrapCSSClassCellRenderer(obj), nil
}

func nativeCellRenderer(v *gtk.CellRenderer) *C.GtkCellRenderer {
	if v == nil || v.GObject == nil {
		return nil
	}
	p := unsafe.Pointer(v.GObject)
	return C.toGtkCellRenderer(p)
}

// SetReal sets the real renderer for the CSS class renderer
func (v *CSSClassCellRenderer) SetReal(real *gtk.CellRenderer) {
	n := v.native()
	C.css_class_cell_renderer_set_real(n, nativeCellRenderer(real))
}
