package gdka

import (
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/gotk3/gotk3/gdk"
)

func init() {
	gliba.AddWrapper(WrapLocal)

	gliba.AddUnwrapper(UnwrapLocal)
}

func WrapLocal(o interface{}) (interface{}, bool) {
	switch oo := o.(type) {
	case *gdk.EventButton:
		val := wrapEventButtonSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gdk.Event:
		val := wrapEventSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gdk.Pixbuf:
		val := wrapPixbufSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gdk.PixbufLoader:
		val := wrapPixbufLoaderSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gdk.Screen:
		val := wrapScreenSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gdk.Window:
		val := WrapWindowSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	default:
		return nil, false
	}
}

func UnwrapLocal(o interface{}) (interface{}, bool) {
	switch oo := o.(type) {
	case *eventButton:
		val := unwrapEventButton(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *event:
		val := unwrapEvent(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *pixbuf:
		val := UnwrapPixbuf(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *pixbufLoader:
		val := unwrapPixbufLoader(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *screen:
		val := UnwrapScreen(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *window:
		val := UnwrapWindow(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	default:
		return nil, false
	}
}
