package gdka

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/twstrike/gotk3adapter/gdki"
)

type eventKey struct {
	*gdk.EventKey
}

func wrapEventKeySimple(v *gdk.EventKey) *eventKey {
	if v == nil {
		return nil
	}
	return &eventKey{v}
}

func wrapEventKey(v *gdk.EventKey, e error) (*eventKey, error) {
	return wrapEventKeySimple(v), e
}

func unwrapEventKey(v gdki.EventKey) *gdk.EventKey {
	if v == nil {
		return nil
	}
	return v.(*eventKey).EventKey
}
