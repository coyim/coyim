package gotk3extra

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

// WrapHeaderBar can be used to create a GTK header bar instance from the given object.
func WrapHeaderBar(obj *glib.Object) *gtk.HeaderBar {
	if obj == nil {
		return nil
	}

	return &gtk.HeaderBar{gtk.Container{gtk.Widget{glib.InitiallyUnowned{obj}}}}
}
