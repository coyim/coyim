package gdka

import "github.com/gotk3/gotk3/gdk"
import "github.com/coyim/gotk3adapter/gdki"

func WrapEventAsEventButton(v *event) gdki.EventButton {
	wrapped, _ := WrapEventButton(&gdk.EventButton{v.Event}, nil)
	return wrapped
}
