package gotk3extra

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

// WrapNotebook can be used to create a GTK notebook instance from the given object.
func WrapNotebook(obj *glib.Object) *gtk.Notebook {
	if obj == nil {
		return nil
	}

	return &gtk.Notebook{gtk.Container{gtk.Widget{glib.InitiallyUnowned{obj}}}}
}
