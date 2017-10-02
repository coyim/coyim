package gdka

import (
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/gotk3/gotk3/gdk"
)

type window struct {
	*gliba.Object
	internal *gdk.Window
}

func WrapWindowSimple(v *gdk.Window) *window {
	if v == nil {
		return nil
	}
	return &window{gliba.WrapObjectSimple(v.Object), v}
}

func WrapWindow(v *gdk.Window, e error) (*window, error) {
	return WrapWindowSimple(v), e
}

func UnwrapWindow(v gdki.Window) *gdk.Window {
	if v == nil {
		return nil
	}
	return v.(*window).internal
}

func (v *window) GetDesktop() uint32 {
	return v.internal.GetDesktop()
}

func (v *window) MoveToDesktop(v1 uint32) {
	v.internal.MoveToDesktop(v1)
}
