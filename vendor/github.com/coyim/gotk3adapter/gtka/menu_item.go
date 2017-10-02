package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gtki"
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

func wrapMenuItemSimple(v *gtk.MenuItem) *menuItem {
	if v == nil {
		return nil
	}
	return &menuItem{wrapBinSimple(&v.Bin), v}
}

func wrapMenuItem(v *gtk.MenuItem, e error) (*menuItem, error) {
	return wrapMenuItemSimple(v), e
}

func unwrapMenuItem(v gtki.MenuItem) *gtk.MenuItem {
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
	v.internal.SetSubmenu(unwrapWidget(v1))
}
