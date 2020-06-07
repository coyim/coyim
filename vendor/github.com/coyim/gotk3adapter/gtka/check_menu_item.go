package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type checkMenuItem struct {
	*menuItem
	internal *gtk.CheckMenuItem
}

func WrapCheckMenuItemSimple(v *gtk.CheckMenuItem) gtki.CheckMenuItem {
	if v == nil {
		return nil
	}
	return &checkMenuItem{WrapMenuItemSimple(&v.MenuItem).(*menuItem), v}
}

func WrapCheckMenuItem(v *gtk.CheckMenuItem, e error) (gtki.CheckMenuItem, error) {
	return WrapCheckMenuItemSimple(v), e
}

func UnwrapCheckMenuItem(v gtki.CheckMenuItem) *gtk.CheckMenuItem {
	if v == nil {
		return nil
	}
	return v.(*checkMenuItem).internal
}

func (v *checkMenuItem) GetActive() bool {
	return v.internal.GetActive()
}

func (v *checkMenuItem) SetActive(v1 bool) {
	v.internal.SetActive(v1)
}
