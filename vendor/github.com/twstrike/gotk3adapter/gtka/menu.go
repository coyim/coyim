package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/gotk3adapter/gdka"
	"github.com/twstrike/gotk3adapter/gdki"
	"github.com/twstrike/gotk3adapter/gtki"
)

type menu struct {
	*menuShell
	internal *gtk.Menu
}

func wrapMenuSimple(v *gtk.Menu) *menu {
	if v == nil {
		return nil
	}
	return &menu{wrapMenuShellSimple(&v.MenuShell), v}
}

func wrapMenu(v *gtk.Menu, e error) (*menu, error) {
	return wrapMenuSimple(v), e
}

func unwrapMenu(v gtki.Menu) *gtk.Menu {
	if v == nil {
		return nil
	}
	return v.(*menu).internal
}

func unwrapMenuToIMenu(v gtki.Menu) gtk.IMenu {
	if v == nil {
		return nil
	}
	return v.(*menu).internal
}

func (v *menu) PopupAtMouseCursor(v1 gtki.Menu, v2 gtki.MenuItem, v3 int, v4 uint32) {
	v.internal.PopupAtMouseCursor(unwrapMenuToIMenu(v1), unwrapMenuItemToIMenuItem(v2), v3, v4)
}

func (v *menu) PopupAtPointer(v1 gdki.Event) {
	v.internal.PopupAtPointer(gdka.UnwrapEvent(v1))
}
