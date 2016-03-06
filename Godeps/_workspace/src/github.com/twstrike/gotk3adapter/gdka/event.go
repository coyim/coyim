package gdka

import (
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/gotk3/gotk3/gdk"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gdki"
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
