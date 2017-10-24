package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gtki"
)

type checkButton struct {
	*toggleButton
	internal *gtk.CheckButton
}

func wrapCheckButtonSimple(v *gtk.CheckButton) *checkButton {
	if v == nil {
		return nil
	}
	return &checkButton{wrapToggleButtonSimple(&v.ToggleButton), v}
}

func wrapCheckButton(v *gtk.CheckButton, e error) (*checkButton, error) {
	return wrapCheckButtonSimple(v), e
}

func unwrapCheckButton(v gtki.CheckButton) *gtk.CheckButton {
	if v == nil {
		return nil
	}
	return v.(*checkButton).internal
}
