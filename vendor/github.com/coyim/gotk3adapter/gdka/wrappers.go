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
		val := WrapEventButtonSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gdk.Event:
		val := WrapEventSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gdk.Pixbuf:
		val := WrapPixbufSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gdk.PixbufLoader:
		val := WrapPixbufLoaderSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *gdk.Screen:
		val := WrapScreenSimple(oo)
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
		val := UnwrapEventButton(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	case *event:
		val := UnwrapEvent(oo)
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
		val := UnwrapPixbufLoader(oo)
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
