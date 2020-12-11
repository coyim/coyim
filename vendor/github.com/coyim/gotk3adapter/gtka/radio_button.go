package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type radioButton struct {
	*checkButton
	internal *gtk.RadioButton
}

func WrapRadioButtonSimple(v *gtk.RadioButton) gtki.RadioButton {
	if v == nil {
		return nil
	}
	return &radioButton{WrapCheckButtonSimple(&v.CheckButton).(*checkButton), v}
}

func WrapRadioButton(v *gtk.RadioButton, e error) (gtki.RadioButton, error) {
	return WrapRadioButtonSimple(v), e
}

func UnwrapRadioButton(v gtki.RadioButton) *gtk.RadioButton {
	if v == nil {
		return nil
	}
	return v.(*radioButton).internal
}
