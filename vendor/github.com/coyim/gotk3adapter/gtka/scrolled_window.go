package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type scrolledWindow struct {
	*bin
	internal *gtk.ScrolledWindow
}

func WrapScrolledWindowSimple(v *gtk.ScrolledWindow) gtki.ScrolledWindow {
	if v == nil {
		return nil
	}
	return &scrolledWindow{WrapBinSimple(&v.Bin).(*bin), v}
}

func WrapScrolledWindow(v *gtk.ScrolledWindow, e error) (gtki.ScrolledWindow, error) {
	return WrapScrolledWindowSimple(v), e
}

func UnwrapScrolledWindow(v gtki.ScrolledWindow) *gtk.ScrolledWindow {
	if v == nil {
		return nil
	}
	return v.(*scrolledWindow).internal
}

func (v *scrolledWindow) GetVAdjustment() gtki.Adjustment {
	return WrapAdjustmentSimple(v.internal.GetVAdjustment())
}
