package gdka

import (
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/gotk3/gotk3/gdk"
)

type eventButton struct {
	*gdk.EventButton
}

func WrapEventButtonSimple(v *gdk.EventButton) gdki.EventButton {
	if v == nil {
		return nil
	}
	return &eventButton{v}
}

func WrapEventButton(v *gdk.EventButton, e error) (gdki.EventButton, error) {
	return WrapEventButtonSimple(v), e
}

func UnwrapEventButton(v gdki.EventButton) *gdk.EventButton {
	if v == nil {
		return nil
	}
	return v.(*eventButton).EventButton
}
