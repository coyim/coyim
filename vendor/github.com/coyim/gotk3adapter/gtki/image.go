package gtki

import "github.com/coyim/gotk3adapter/gdki"

// Image is an interface of Gtk.Image
type Image interface {
	Widget

	SetFromIconName(string, IconSize)
	SetFromPixbuf(gdki.Pixbuf)
	Clear()
}

// AssertImage asserts the Image
func AssertImage(_ Image) {}
