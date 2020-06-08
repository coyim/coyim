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

func UnwrapMenuShellOnly(v gtki.MenuShell) *gtk.MenuShell {
	if v == nil {
		return nil
	}
	return v.(*menuShell).internal
}

func UnwrapMenuShell(v gtki.MenuShell) *gtk.MenuShell {
	switch oo := v.(type) {
	case *menuBar:
		val := UnwrapMenuBar(oo)
		if val == nil {
			return nil
		}
		return &val.MenuShell
	case *menu:
		val := UnwrapMenu(oo)
		if val == nil {
			return nil
		}
		return &val.MenuShell
	case *menuShell:
		return UnwrapMenuShellOnly(oo)
	default:
		return nil
	}
}

func (v *menuShell) Append(v1 gtki.MenuItem) {
	v.internal.Append(unwrapMenuItemToIMenuItem(v1))
}
