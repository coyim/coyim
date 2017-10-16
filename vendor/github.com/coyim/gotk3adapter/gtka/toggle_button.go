package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type toggleButton struct {
	*button
	*gtk.ToggleButton
}

func wrapToggleButtonSimple(v *gtk.ToggleButton) *toggleButton {
	if v == nil {
		return nil
	}
	return &toggleButton{wrapButtonSimple(&v.Button), v}
}

func wrapToggleButton(v *gtk.ToggleButton, e error) (*toggleButton, error) {
	return wrapToggleButtonSimple(v), e
}

func unwrapToggleButton(v gtki.ToggleButton) *gtk.ToggleButton {
	if v == nil {
		return nil
	}
	return v.(*toggleButton).ToggleButton
}
