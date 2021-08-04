package gotk3extra

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

// WrapButton can be used to create a GTK button instance from the given object.
func WrapButton(obj *glib.Object) *gtk.Button {
	if obj == nil {
		return nil
	}

	actionable := &gtk.Actionable{obj}
	return &gtk.Button{gtk.Bin{gtk.Container{gtk.Widget{glib.InitiallyUnowned{obj}}}}, actionable}
}
