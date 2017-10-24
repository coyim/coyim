package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gtki"
)

type checkMenuItem struct {
	*menuItem
	internal *gtk.CheckMenuItem
}

func wrapCheckMenuItemSimple(v *gtk.CheckMenuItem) *checkMenuItem {
	if v == nil {
		return nil
	}
	return &checkMenuItem{wrapMenuItemSimple(&v.MenuItem), v}
}

func wrapCheckMenuItem(v *gtk.CheckMenuItem, e error) (*checkMenuItem, error) {
	return wrapCheckMenuItemSimple(v), e
}

func unwrapCheckMenuItem(v gtki.CheckMenuItem) *gtk.CheckMenuItem {
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
