package gtka

import (
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtki"
)

type applicationWindow struct {
	*window
	internal *gtk.ApplicationWindow
}

func wrapApplicationWindowSimple(v *gtk.ApplicationWindow) *applicationWindow {
	if v == nil {
		return nil
	}
	return &applicationWindow{wrapWindowSimple(&v.Window), v}
}

func wrapApplicationWindow(v *gtk.ApplicationWindow, e error) (*applicationWindow, error) {
	return wrapApplicationWindowSimple(v), e
}

func unwrapApplicationWindow(v gtki.ApplicationWindow) *gtk.ApplicationWindow {
	if v == nil {
		return nil
	}
	return v.(*applicationWindow).internal
}
