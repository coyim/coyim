package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gtki"
)

type spinButton struct {
	*entry
	internal *gtk.SpinButton
}

func wrapSpinButtonSimple(v *gtk.SpinButton) *spinButton {
	if v == nil {
		return nil
	}
	return &spinButton{wrapEntrySimple(&v.Entry), v}
}

func wrapSpinButton(v *gtk.SpinButton, e error) (*spinButton, error) {
	return wrapSpinButtonSimple(v), e
}

func unwrapSpinButton(v gtki.SpinButton) *gtk.SpinButton {
	if v == nil {
		return nil
	}
	return v.(*spinButton).internal
}

func (v *spinButton) GetValueAsInt() int {
	return v.internal.GetValueAsInt()
}

func (v *spinButton) SetValue(v1 float64) {
	v.internal.SetValue(v1)
}

func (v *spinButton) GetValue() float64 {
	return v.internal.GetValue()
}

func (v *spinButton) GetAdjustment() gtki.Adjustment {
	return wrapAdjustmentSimple(v.internal.GetAdjustment())
}

func (v *spinButton) SetRange(v1 float64, v2 float64) {
	v.internal.SetRange(v1, v2)
}

func (v *spinButton) SetIncrements(v1 float64, v2 float64) {
	v.internal.SetIncrements(v1, v2)
}
