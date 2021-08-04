package gotk3extra

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

// WrapBox can be used to create a GTK box instance from the given object.
func WrapBox(obj *glib.Object) *gtk.Box {
	if obj == nil {
		return nil
	}

	return &gtk.Box{gtk.Container{gtk.Widget{glib.InitiallyUnowned{obj}}}}
}
