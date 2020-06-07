package gdka

import (
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/gotk3/gotk3/gdk"
)

type eventKey struct {
	*gdk.EventKey
}

func WrapEventKeySimple(v *gdk.EventKey) gdki.EventKey {
	if v == nil {
		return nil
	}
	return &eventKey{v}
}

func WrapEventKey(v *gdk.EventKey, e error) (gdki.EventKey, error) {
	return WrapEventKeySimple(v), e
}

func UnwrapEventKey(v gdki.EventKey) *gdk.EventKey {
	if v == nil {
		return nil
	}
	return v.(*eventKey).EventKey
}
