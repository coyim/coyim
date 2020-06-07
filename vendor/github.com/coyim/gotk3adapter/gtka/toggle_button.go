package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type toggleButton struct {
	*button
	internal *gtk.ToggleButton
}

func WrapToggleButtonSimple(v *gtk.ToggleButton) gtki.ToggleButton {
	if v == nil {
		return nil
	}
	return &toggleButton{WrapButtonSimple(&v.Button).(*button), v}
}

func WrapToggleButton(v *gtk.ToggleButton, e error) (gtki.ToggleButton, error) {
	return WrapToggleButtonSimple(v), e
}

func UnwrapToggleButton(v gtki.ToggleButton) *gtk.ToggleButton {
	if v == nil {
		return nil
	}
	return v.(*toggleButton).internal
}

func (v *toggleButton) SetActive(v1 bool) {
	v.internal.SetActive(v1)
}

func (v *toggleButton) GetActive() bool {
	return v.internal.GetActive()
}
