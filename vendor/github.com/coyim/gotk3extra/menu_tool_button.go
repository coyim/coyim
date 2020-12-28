package gotk3extra

// #cgo pkg-config: gdk-3.0 gio-2.0 glib-2.0 gobject-2.0 gtk+-3.0
// #include <stdlib.h>
// #include <gtk/gtk.h>
// #include "menu_tool_button.go.h"
// #cgo CFLAGS: -Wno-deprecated-declarations
import "C"
import (
	"unsafe"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func init() {
	tm := []glib.TypeMarshaler{
		{glib.Type(C.gtk_menu_tool_button_get_type()), marshalMenuToolButton},
	}
	glib.RegisterGValueMarshalers(tm)

	gtk.WrapMap["GtkMenuToolButton"] = wrapMenuToolButton
}

// MenuToolButton is a representation of GTK's GtkMenuToolButton.
type MenuToolButton struct {
	gtk.ToolButton
}

func marshalMenuToolButton(p uintptr) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(unsafe.Pointer(p)))
	obj := glib.Take(unsafe.Pointer(c))
	return wrapMenuToolButton(obj), nil
}

func wrapMenuToolButton(obj *glib.Object) *MenuToolButton {
	return &MenuToolButton{gtk.ToolButton{gtk.ToolItem{gtk.Bin{gtk.Container{gtk.Widget{
		glib.InitiallyUnowned{obj}}}}}}}
}

func (v *MenuToolButton) native() *C.GtkMenuToolButton {
	if v == nil || v.GObject == nil {
		return nil
	}
	p := unsafe.Pointer(v.GObject)
	return C.toGtkMenuToolButton(p)
}
