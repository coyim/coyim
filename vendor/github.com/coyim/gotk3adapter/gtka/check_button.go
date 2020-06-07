package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type checkButton struct {
	*toggleButton
	internal *gtk.CheckButton
}

func WrapCheckButtonSimple(v *gtk.CheckButton) gtki.CheckButton {
	if v == nil {
		return nil
	}
	return &checkButton{WrapToggleButtonSimple(&v.ToggleButton).(*toggleButton), v}
}

func WrapCheckButton(v *gtk.CheckButton, e error) (gtki.CheckButton, error) {
	return WrapCheckButtonSimple(v), e
}

func UnwrapCheckButton(v gtki.CheckButton) *gtk.CheckButton {
	if v == nil {
		return nil
	}
	return v.(*checkButton).internal
}
