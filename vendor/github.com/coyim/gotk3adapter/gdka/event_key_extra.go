package gdka

import "github.com/gotk3/gotk3/gdk"

func wrapEventAsEventKey(v *event) *eventKey {
	wrapped, _ := wrapEventKey(&gdk.EventKey{v.Event}, nil)
	return wrapped
}
