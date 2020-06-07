// +build darwin

package gotk3osx

// #include <gtk/gtk.h>
// #include <gtkosxapplication.h>
// #include "gtkosx.go.h"
import "C"
import "unsafe"
import "github.com/gotk3/gotk3/gtk"

func nativeMenuItem(v *gtk.MenuItem) *C.GtkMenuItem {
	if v == nil || v.GObject == nil {
		return nil
	}
	p := unsafe.Pointer(v.GObject)
	return C.toGtkMenuItem(p)
}
