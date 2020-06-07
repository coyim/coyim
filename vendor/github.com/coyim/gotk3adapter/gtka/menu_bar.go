package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type menuBar struct {
	*menuShell
	internal *gtk.MenuBar
}

func WrapMenuBarSimple(v *gtk.MenuBar) gtki.MenuBar {
	if v == nil {
		return nil
	}
	return &menuBar{WrapMenuShellSimple(&v.MenuShell).(*menuShell), v}
}

func WrapMenuBar(v *gtk.MenuBar, e error) (gtki.MenuBar, error) {
	return WrapMenuBarSimple(v), e
}

func UnwrapMenuBar(v gtki.MenuBar) *gtk.MenuBar {
	if v == nil {
		return nil
	}
	return v.(*menuBar).internal
}
