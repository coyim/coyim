package gtka

import (
	"github.com/coyim/gotk3adapter/gdka"
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type menu struct {
	*menuShell
	internal *gtk.Menu
}

func WrapMenuSimple(v *gtk.Menu) gtki.Menu {
	if v == nil {
		return nil
	}
	return &menu{WrapMenuShellSimple(&v.MenuShell).(*menuShell), v}
}

func WrapMenu(v *gtk.Menu, e error) (gtki.Menu, error) {
	return WrapMenuSimple(v), e
}

func UnwrapMenu(v gtki.Menu) *gtk.Menu {
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

func (v *menu) PopupAtPointer(v1 gdki.Event) {
	v.internal.PopupAtPointer(gdka.UnwrapEvent(v1))
}
