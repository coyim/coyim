package gtki

// Image is an interface of Gtk.Image
type Image interface {
	Widget

	SetFromIconName(string, IconSize)
	Clear()
}

// AssertImage asserts the Image
func AssertImage(_ Image) {}
