package gotk3extra

// #cgo pkg-config: gdk-3.0 gio-2.0 glib-2.0 gobject-2.0 gtk+-3.0
// #include <stdlib.h>
// #include <gtk/gtk.h>
// #include "status_icon.go.h"
// #cgo CFLAGS: -Wno-deprecated-declarations
import "C"
import (
	"errors"
	"unsafe"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func init() {
	tm := []glib.TypeMarshaler{
		{glib.Type(C.gtk_status_icon_get_type()), marshalStatusIcon},
	}
	glib.RegisterGValueMarshalers(tm)

	gtk.WrapMap["GtkStatusIcon"] = wrapStatusIcon
}

// StatusIcon is a representation of GTK's GtkStatusIcon.
// Deprecated since 3.14 in favor of notifications
// (no replacement, see https://stackoverflow.com/questions/41917903/gtk-3-statusicon-replacement)
type StatusIcon struct {
	*glib.Object
}

func marshalStatusIcon(p uintptr) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(unsafe.Pointer(p)))
	obj := glib.Take(unsafe.Pointer(c))
	return wrapStatusIcon(obj), nil
}

func wrapStatusIcon(obj *glib.Object) *StatusIcon {
	return &StatusIcon{obj}
}

func (v *StatusIcon) native() *C.GtkStatusIcon {
	if v == nil || v.GObject == nil {
		return nil
	}
	p := unsafe.Pointer(v.GObject)
	return C.toGtkStatusIcon(p)
}

// StatusIconNew is a wrapper around gtk_status_icon_new()
func StatusIconNew() (*StatusIcon, error) {
	c := C.gtk_status_icon_new()
	if c == nil {
		return nil, nilPtrErr
	}
	return wrapStatusIcon(glib.Take(unsafe.Pointer(c))), nil
}

// StatusIconNewFromFile is a wrapper around gtk_status_icon_new_from_file()
func StatusIconNewFromFile(filename string) (*StatusIcon, error) {
	cstr := C.CString(filename)
	defer C.free(unsafe.Pointer(cstr))
	c := C.gtk_status_icon_new_from_file((*C.gchar)(cstr))
	if c == nil {
		return nil, nilPtrErr
	}
	return wrapStatusIcon(glib.Take(unsafe.Pointer(c))), nil
}

// StatusIconNewFromIconName is a wrapper around gtk_status_icon_new_from_icon_name()
func StatusIconNewFromIconName(iconName string) (*StatusIcon, error) {
	cstr := C.CString(iconName)
	defer C.free(unsafe.Pointer(cstr))
	c := C.gtk_status_icon_new_from_icon_name((*C.gchar)(cstr))
	if c == nil {
		return nil, nilPtrErr
	}
	return wrapStatusIcon(glib.Take(unsafe.Pointer(c))), nil
}

// StatusIconNewFromPixbuf is a wrapper around gtk_status_icon_new_from_pixbuf().
func StatusIconNewFromPixbuf(pixbuf *gdk.Pixbuf) (*StatusIcon, error) {
	c := C.gtk_status_icon_new_from_pixbuf(C.toGdkPixbuf(unsafe.Pointer(pixbuf.Native())))
	if c == nil {
		return nil, nilPtrErr
	}
	obj := glib.Take(unsafe.Pointer(c))
	return wrapStatusIcon(obj), nil
}

// SetFromFile is a wrapper around gtk_status_icon_set_from_file()
func (v *StatusIcon) SetFromFile(filename string) {
	cstr := C.CString(filename)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_status_icon_set_from_file(v.native(), (*C.gchar)(cstr))
}

// SetFromIconName is a wrapper around gtk_status_icon_set_from_icon_name()
func (v *StatusIcon) SetFromIconName(iconName string) {
	cstr := C.CString(iconName)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_status_icon_set_from_icon_name(v.native(), (*C.gchar)(cstr))
}

// SetFromPixbuf is a wrapper around gtk_status_icon_set_from_pixbuf()
func (v *StatusIcon) SetFromPixbuf(pixbuf *gdk.Pixbuf) {
	C.gtk_status_icon_set_from_pixbuf(v.native(), C.toGdkPixbuf(unsafe.Pointer(pixbuf.Native())))
}

// GetStorageType is a wrapper around gtk_status_icon_get_storage_type()
func (v *StatusIcon) GetStorageType() gtk.ImageType {
	return (gtk.ImageType)(C.gtk_status_icon_get_storage_type(v.native()))
}

// SetTooltipText is a wrapper around gtk_status_icon_set_tooltip_text()
func (v *StatusIcon) SetTooltipText(text string) {
	cstr := C.CString(text)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_status_icon_set_tooltip_text(v.native(), (*C.gchar)(cstr))
}

// GetTooltipText is a wrapper around gtk_status_icon_get_tooltip_text()
func (v *StatusIcon) GetTooltipText() string {
	c := C.gtk_status_icon_get_tooltip_text(v.native())
	if c == nil {
		return ""
	}
	return C.GoString((*C.char)(c))
}

// SetTooltipMarkup is a wrapper around gtk_status_icon_set_tooltip_markup()
func (v *StatusIcon) SetTooltipMarkup(markup string) {
	cstr := C.CString(markup)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_status_icon_set_tooltip_markup(v.native(), (*C.gchar)(cstr))
}

// GetTooltipMarkup is a wrapper around gtk_status_icon_get_tooltip_markup()
func (v *StatusIcon) GetTooltipMarkup() string {
	c := C.gtk_status_icon_get_tooltip_markup(v.native())
	if c == nil {
		return ""
	}
	return C.GoString((*C.char)(c))
}

// SetHasTooltip is a wrapper around gtk_status_icon_set_has_tooltip()
func (v *StatusIcon) SetHasTooltip(hasTooltip bool) {
	C.gtk_status_icon_set_has_tooltip(v.native(), gbool(hasTooltip))
}

// GetTitle is a wrapper around gtk_status_icon_get_title()
func (v *StatusIcon) GetTitle() string {
	c := C.gtk_status_icon_get_title(v.native())
	if c == nil {
		return ""
	}
	return C.GoString((*C.char)(c))
}

// SetName is a wrapper around gtk_status_icon_set_name()
func (v *StatusIcon) SetName(name string) {
	cstr := C.CString(name)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_status_icon_set_name(v.native(), (*C.gchar)(cstr))
}

// SetVisible is a wrapper around gtk_status_icon_set_visible()
func (v *StatusIcon) SetVisible(visible bool) {
	C.gtk_status_icon_set_visible(v.native(), gbool(visible))
}

// GetVisible is a wrapper around gtk_status_icon_get_visible()
func (v *StatusIcon) GetVisible() bool {
	return gobool(C.gtk_status_icon_get_visible(v.native()))
}

// IsEmbedded is a wrapper around gtk_status_icon_is_embedded()
func (v *StatusIcon) IsEmbedded() bool {
	return gobool(C.gtk_status_icon_is_embedded(v.native()))
}

// GetX11WindowID is a wrapper around gtk_status_icon_get_x11_window_id()
func (v *StatusIcon) GetX11WindowID() uint32 {
	return uint32(C.gtk_status_icon_get_x11_window_id(v.native()))
}

// GetHasTooltip is a wrapper around gtk_status_icon_get_has_tooltip()
func (v *StatusIcon) GetHasTooltip() bool {
	return gobool(C.gtk_status_icon_get_has_tooltip(v.native()))
}

// SetTitle is a wrapper around gtk_status_icon_set_title()
func (v *StatusIcon) SetTitle(title string) {
	cstr := C.CString(title)
	defer C.free(unsafe.Pointer(cstr))
	C.gtk_status_icon_set_title(v.native(), (*C.gchar)(cstr))
}

// GetIconName is a wrapper around gtk_status_icon_get_icon_name()
func (v *StatusIcon) GetIconName() string {
	c := C.gtk_status_icon_get_icon_name(v.native())
	if c == nil {
		return ""
	}
	return C.GoString((*C.char)(c))
}

// GetSize is a wrapper around gtk_status_icon_get_size()
func (v *StatusIcon) GetSize() int {
	return int(C.gtk_status_icon_get_size(v.native()))
}

var nilPtrErr = errors.New("cgo returned unexpected nil pointer")

func gbool(b bool) C.gboolean {
	if b {
		return C.gboolean(1)
	}
	return C.gboolean(0)
}

func gobool(b C.gboolean) bool {
	return b != C.FALSE
}
