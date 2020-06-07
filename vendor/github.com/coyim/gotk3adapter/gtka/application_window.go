package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type applicationWindow struct {
	*window
	internal *gtk.ApplicationWindow
}

func WrapApplicationWindowSimple(v *gtk.ApplicationWindow) gtki.ApplicationWindow {
	if v == nil {
		return nil
	}
	return &applicationWindow{WrapWindowSimple(&v.Window).(*window), v}
}

func WrapApplicationWindow(v *gtk.ApplicationWindow, e error) (gtki.ApplicationWindow, error) {
	return WrapApplicationWindowSimple(v), e
}

func UnwrapApplicationWindow(v gtki.ApplicationWindow) *gtk.ApplicationWindow {
	if v == nil {
		return nil
	}
	return v.(*applicationWindow).internal
}

func (v *applicationWindow) SetShowMenubar(val bool) {
	v.internal.SetShowMenubar(val)
}

func (v *applicationWindow) GetShowMenubar() bool {
	return v.internal.GetShowMenubar()
}

func (v *applicationWindow) GetID() uint {
	return v.internal.GetID()
}
