// +build darwin

package gotk3osx

// #include <gtk/gtk.h>
// #include <gdk/gdk.h>
// #include <gtkosxapplication.h>
// #include "gtkosx.go.h"
import "C"
import "unsafe"
import "github.com/gotk3/gotk3/gdk"

func nativePixbuf(v *gdk.Pixbuf) *C.GdkPixbuf {
	if v == nil || v.GObject == nil {
		return nil
	}
	p := unsafe.Pointer(v.GObject)
	return C.toGdkPixbuf(p)
}
