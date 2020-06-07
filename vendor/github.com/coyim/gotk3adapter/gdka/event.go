package gdka

import (
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/gotk3/gotk3/gdk"
)

type event struct {
	*gdk.Event
}

func WrapEventSimple(v *gdk.Event) gdki.Event {
	if v == nil {
		return nil
	}
	return &event{v}
}

func WrapEvent(v *gdk.Event, e error) (gdki.Event, error) {
	return WrapEventSimple(v), e
}

func UnwrapEventOnly(v gdki.Event) *gdk.Event {
	if v == nil {
		return nil
	}
	return v.(*event).Event
}

func UnwrapEvent(v gdki.Event) *gdk.Event {
	switch oo := v.(type) {
	case *eventButton:
		val := UnwrapEventButton(oo)
		if val == nil {
			return nil
		}
		return val.Event
	case *eventKey:
		val := UnwrapEventKey(oo)
		if val == nil {
			return nil
		}
		return val.Event
	case *event:
		val := UnwrapEventOnly(oo)
		if val == nil {
			return nil
		}
		return val
	default:
		return nil
	}
}
