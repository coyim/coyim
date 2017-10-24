package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gtki"
)

type menuBar struct {
	*menuShell
	internal *gtk.MenuBar
}

func wrapMenuBarSimple(v *gtk.MenuBar) *menuBar {
	if v == nil {
		return nil
	}
	return &menuBar{wrapMenuShellSimple(&v.MenuShell), v}
}

func wrapMenuBar(v *gtk.MenuBar, e error) (*menuBar, error) {
	return wrapMenuBarSimple(v), e
}

func unwrapMenuBar(v gtki.MenuBar) *gtk.MenuBar {
	if v == nil {
		return nil
	}
	return v.(*menuBar).internal
}
