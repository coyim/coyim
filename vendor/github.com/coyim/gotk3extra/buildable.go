package gotk3extra

// #include <gtk/gtk.h>
// #include "buildable.go.h"
import "C"
import (
	"unsafe"

	"github.com/gotk3/gotk3/glib"
)

func nativeBuildable(v *glib.Object) *C.GtkBuildable {
	if v == nil || v.GObject == nil {
		return nil
	}
	p := unsafe.Pointer(v.GObject)
	return C.toGtkBuildable(p)
}

// GetBuildableName is a wrapper around gtk_buildable_get_name().
func GetBuildableName(obj *glib.Object) (string, error) {
	c := C.gtk_buildable_get_name(nativeBuildable(obj))
	if c == nil {
		return "", nilPtrErr
	}
	return C.GoString((*C.char)(c)), nil
}
