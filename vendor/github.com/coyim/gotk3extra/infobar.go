package gotk3extra

// #include <gtk/gtk.h>
// #include "infobar.go.h"
import "C"
import (
	"unsafe"

	"github.com/gotk3/gotk3/gtk"
)

func nativeInfoBar(v *gtk.InfoBar) *C.GtkInfoBar {
	if v == nil || v.GObject == nil {
		return nil
	}

	p := unsafe.Pointer(v.GObject)
	return C.toGtkInfoBar(p)
}

// InfoBarSetRevealed sets the revealed property to the infobar's revealer.
// This will cause infobar to show up with a slide-in transition.
// Note that this property does not automatically show infobar and
// thus won’t have any effect if it is invisible.
func InfoBarSetRevealed(infobar *gtk.InfoBar, setting bool) {
	C.gtk_info_bar_set_revealed(nativeInfoBar(infobar), gbool(setting))
}

// InfoBarGetRevealed returns the current value of the infobar's revealed property.
func InfoBarGetRevealed(infobar *gtk.InfoBar) bool {
	b := C.gtk_info_bar_get_revealed(nativeInfoBar(infobar))
	return gobool(b)
}
