package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gtki"
)

type menuShell struct {
	*container
	internal *gtk.MenuShell
}

func wrapMenuShellSimple(v *gtk.MenuShell) *menuShell {
	if v == nil {
		return nil
	}
	return &menuShell{wrapContainerSimple(&v.Container), v}
}

func wrapMenuShell(v *gtk.MenuShell, e error) (*menuShell, error) {
	return wrapMenuShellSimple(v), e
}

func unwrapMenuShell(v gtki.MenuShell) *gtk.MenuShell {
	if v == nil {
		return nil
	}
	return v.(*menuShell).internal
}

func (v *menuShell) Append(v1 gtki.MenuItem) {
	v.internal.Append(unwrapMenuItemToIMenuItem(v1))
}
