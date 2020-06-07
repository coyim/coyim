package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type menuItem struct {
	*bin
	internal *gtk.MenuItem
}

type asMenuItem interface {
	toMenuItem() *menuItem
}

func (v *menuItem) toMenuItem() *menuItem {
	return v
}

func WrapMenuItemSimple(v *gtk.MenuItem) gtki.MenuItem {
	if v == nil {
		return nil
	}
	return &menuItem{WrapBinSimple(&v.Bin).(*bin), v}
}

func WrapMenuItem(v *gtk.MenuItem, e error) (gtki.MenuItem, error) {
	return WrapMenuItemSimple(v), e
}

func UnwrapMenuItem(v gtki.MenuItem) *gtk.MenuItem {
	if v == nil {
		return nil
	}
	return v.(asMenuItem).toMenuItem().internal
}

func unwrapMenuItemToIMenuItem(v gtki.MenuItem) gtk.IMenuItem {
	if v == nil {
		return nil
	}
	return v.(asMenuItem).toMenuItem().internal
}

func (v *menuItem) GetLabel() string {
	return v.internal.GetLabel()
}

func (v *menuItem) SetLabel(v1 string) {
	v.internal.SetLabel(v1)
}

func (v *menuItem) SetSubmenu(v1 gtki.Widget) {
	v.internal.SetSubmenu(UnwrapWidget(v1))
}
