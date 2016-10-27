package gdka

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/twstrike/gotk3adapter/gdki"
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
