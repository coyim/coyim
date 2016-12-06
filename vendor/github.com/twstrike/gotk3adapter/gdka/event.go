package gdka

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/twstrike/gotk3adapter/gdki"
)

type event struct {
	*gdk.Event
}

func wrapEventSimple(v *gdk.Event) *event {
	if v == nil {
		return nil
	}
	return &event{v}
}

func wrapEvent(v *gdk.Event, e error) (*event, error) {
	return wrapEventSimple(v), e
}

func unwrapEvent(v gdki.Event) *gdk.Event {
	if v == nil {
		return nil
	}
	return v.(*event).Event
}

func UnwrapEvent(v gdki.Event) *gdk.Event {
	switch oo := v.(type) {
	case *eventButton:
		val := unwrapEventButton(oo)
		if val == nil {
			return nil
		}
		return val.Event
	case *eventKey:
		val := unwrapEventKey(oo)
		if val == nil {
			return nil
		}
		return val.Event
	case *event:
		val := unwrapEvent(oo)
		if val == nil {
			return nil
		}
		return val
	default:
		return nil
	}
}
