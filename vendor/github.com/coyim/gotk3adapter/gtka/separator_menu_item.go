package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type separatorMenuItem struct {
	*menuItem
	internal *gtk.SeparatorMenuItem
}

func WrapSeparatorMenuItemSimple(v *gtk.SeparatorMenuItem) gtki.SeparatorMenuItem {
	if v == nil {
		return nil
	}
	return &separatorMenuItem{WrapMenuItemSimple(&v.MenuItem).(*menuItem), v}
}

func WrapSeparatorMenuItem(v *gtk.SeparatorMenuItem, e error) (gtki.SeparatorMenuItem, error) {
	return WrapSeparatorMenuItemSimple(v), e
}

func UnwrapSeparatorMenuItem(v gtki.SeparatorMenuItem) *gtk.SeparatorMenuItem {
	if v == nil {
		return nil
	}
	return v.(*separatorMenuItem).internal
}
