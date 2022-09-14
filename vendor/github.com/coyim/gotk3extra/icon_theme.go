package gotk3extra

// #cgo pkg-config: gdk-3.0 gio-2.0 glib-2.0 gobject-2.0 gtk+-3.0
// #include <stdlib.h>
// #include <gtk/gtk.h>
// #include "icon_theme.go.h"
// #cgo CFLAGS: -Wno-deprecated-declarations
import "C"
import (
	"unsafe"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func init() {
	tm := []glib.TypeMarshaler{
		{glib.Type(C.gtk_icon_theme_get_type()), marshalIconTheme},
	}
	glib.RegisterGValueMarshalers(tm)

	gtk.WrapMap["GtkIconTheme"] = wrapIconTheme
}

// IconTheme is a representation of GTK's GtkIconTheme.
type IconTheme struct {
	*glib.Object
}

func marshalIconTheme(p uintptr) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(unsafe.Pointer(p)))
	obj := glib.Take(unsafe.Pointer(c))
	return wrapIconTheme(obj), nil
}

func wrapIconTheme(obj *glib.Object) *IconTheme {
	return &IconTheme{obj}
}

func (v *IconTheme) native() *C.GtkIconTheme {
	if v == nil || v.GObject == nil {
		return nil
	}
	p := unsafe.Pointer(v.GObject)
	return C.toGtkIconTheme(p)
}

// IconThemeNew is a wrapper around gtk_icon_theme_new()
func IconThemeNew() (*IconTheme, error) {
	c := C.gtk_icon_theme_new()
	if c == nil {
		return nil, nilPtrErr
	}
	return wrapIconTheme(glib.Take(unsafe.Pointer(c))), nil
}

// IconThemeGetDefault is a wrapper around gtk_icon_theme_get_default()
func IconThemeGetDefault() *IconTheme {
	c := C.gtk_icon_theme_get_default()
	if c == nil {
		return nil
	}
	return wrapIconTheme(glib.Take(unsafe.Pointer(c)))
}

// IconThemeGetForScreen is a wrapper around gtk_icon_theme_get_for_screen()
func IconThemeGetForScreen(screen *gdk.Screen) *IconTheme {
	c := C.gtk_icon_theme_get_for_screen(C.toGdkScreen(unsafe.Pointer(screen.Native())))
	if c == nil {
		return nil
	}
	return wrapIconTheme(glib.Take(unsafe.Pointer(c)))
}

// AddResourcePath is a wrapper around gtk_icon_theme_add_resource_path()
func (v *IconTheme) AddResourcePath(s string) {
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_icon_theme_add_resource_path(v.native(), (*C.gchar)(cstr))
}

// AppendSearchPath is a wrapper around gtk_icon_theme_append_search_path()
func (v *IconTheme) AppendSearchPath(s string) {
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_icon_theme_append_search_path(v.native(), (*C.gchar)(cstr))
}

// GetExampleIconName is a wrapper around gtk_icon_theme_get_example_icon_name()
func (v *IconTheme) GetExampleIconName() string {
	c := C.gtk_icon_theme_get_example_icon_name(v.native())
	if c == nil {
		return ""
	}
	return C.GoString((*C.char)(c))
}

// HasIcon is a wrapper around gtk_icon_theme_has_icon()
func (v *IconTheme) HasIcon(s string) bool {
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))
	return gobool(C.gtk_icon_theme_has_icon(v.native(), (*C.gchar)(cstr)))
}

// PrependSearchPath is a wrapper around gtk_icon_theme_prepend_search_path()
func (v *IconTheme) PrependSearchPath(s string) {
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_icon_theme_prepend_search_path(v.native(), (*C.gchar)(cstr))
}
