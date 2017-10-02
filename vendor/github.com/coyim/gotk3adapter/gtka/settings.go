package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/coyim/gotk3adapter/gtki"
)

type settings struct {
	*gliba.Object
	internal *gtk.Settings
}

func wrapSettingsSimple(v *gtk.Settings) *settings {
	if v == nil {
		return nil
	}
	return &settings{gliba.WrapObjectSimple(v.Object), v}
}

func wrapSettings(v *gtk.Settings, e error) (*settings, error) {
	return wrapSettingsSimple(v), e
}

func unwrapSettings(v gtki.Settings) *gtk.Settings {
	if v == nil {
		return nil
	}
	return v.(*settings).internal
}
