package gliba

import (
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/gotk3/gotk3/glib"
)

type settingsBackend struct {
	*Object
	*glib.SettingsBackend
}

func WrapSettingsBackendSimple(v *glib.SettingsBackend) glibi.SettingsBackend {
	if v == nil {
		return nil
	}
	return &settingsBackend{WrapObjectSimple(v.Object), v}
}

func UnwrapSettingsBackend(v glibi.SettingsBackend) *glib.SettingsBackend {
	if v == nil {
		return nil
	}
	return v.(*settingsBackend).SettingsBackend
}
