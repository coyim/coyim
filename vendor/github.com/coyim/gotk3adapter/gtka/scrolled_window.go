package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gtki"
)

type scrolledWindow struct {
	*bin
	internal *gtk.ScrolledWindow
}

func wrapScrolledWindowSimple(v *gtk.ScrolledWindow) *scrolledWindow {
	if v == nil {
		return nil
	}
	return &scrolledWindow{wrapBinSimple(&v.Bin), v}
}

func wrapScrolledWindow(v *gtk.ScrolledWindow, e error) (*scrolledWindow, error) {
	return wrapScrolledWindowSimple(v), e
}

func unwrapScrolledWindow(v gtki.ScrolledWindow) *gtk.ScrolledWindow {
	if v == nil {
		return nil
	}
	return v.(*scrolledWindow).internal
}

func (v *scrolledWindow) GetVAdjustment() gtki.Adjustment {
	return wrapAdjustmentSimple(v.internal.GetVAdjustment())
}
