package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type menuShell struct {
	*container
	internal *gtk.MenuShell
}

func WrapMenuShellSimple(v *gtk.MenuShell) gtki.MenuShell {
	if v == nil {
		return nil
	}
	return &menuShell{WrapContainerSimple(&v.Container).(*container), v}
}

func WrapMenuShell(v *gtk.MenuShell, e error) (gtki.MenuShell, error) {
	return WrapMenuShellSimple(v), e
}

func UnwrapMenuShell(v gtki.MenuShell) *gtk.MenuShell {
	if v == nil {
		return nil
	}
	return v.(*menuShell).internal
}

func (v *menuShell) Append(v1 gtki.MenuItem) {
	v.internal.Append(unwrapMenuItemToIMenuItem(v1))
}
