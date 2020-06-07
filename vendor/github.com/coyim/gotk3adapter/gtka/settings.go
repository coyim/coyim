package gtka

import (
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type settings struct {
	*gliba.Object
	internal *gtk.Settings
}

func WrapSettingsSimple(v *gtk.Settings) gtki.Settings {
	if v == nil {
		return nil
	}
	return &settings{gliba.WrapObjectSimple(v.Object), v}
}

func WrapSettings(v *gtk.Settings, e error) (gtki.Settings, error) {
	return WrapSettingsSimple(v), e
}

func UnwrapSettings(v gtki.Settings) *gtk.Settings {
	if v == nil {
		return nil
	}
	return v.(*settings).internal
}
