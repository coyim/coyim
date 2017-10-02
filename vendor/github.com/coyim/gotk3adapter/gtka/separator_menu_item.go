package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gtki"
)

type separatorMenuItem struct {
	*menuItem
	internal *gtk.SeparatorMenuItem
}

func wrapSeparatorMenuItemSimple(v *gtk.SeparatorMenuItem) *separatorMenuItem {
	if v == nil {
		return nil
	}
	return &separatorMenuItem{wrapMenuItemSimple(&v.MenuItem), v}
}

func wrapSeparatorMenuItem(v *gtk.SeparatorMenuItem, e error) (*separatorMenuItem, error) {
	return wrapSeparatorMenuItemSimple(v), e
}

func unwrapSeparatorMenuItem(v gtki.SeparatorMenuItem) *gtk.SeparatorMenuItem {
	if v == nil {
		return nil
	}
	return v.(*separatorMenuItem).internal
}
