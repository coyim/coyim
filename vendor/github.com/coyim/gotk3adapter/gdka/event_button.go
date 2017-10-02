package gdka

import (
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/gotk3/gotk3/gdk"
)

type eventButton struct {
	*gdk.EventButton
}

func wrapEventButtonSimple(v *gdk.EventButton) *eventButton {
	if v == nil {
		return nil
	}
	return &eventButton{v}
}

func wrapEventButton(v *gdk.EventButton, e error) (*eventButton, error) {
	return wrapEventButtonSimple(v), e
}

func unwrapEventButton(v gdki.EventButton) *gdk.EventButton {
	if v == nil {
		return nil
	}
	return v.(*eventButton).EventButton
}
