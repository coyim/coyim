package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type bin struct {
	*container
	*gtk.Bin
}

func WrapBinSimple(v *gtk.Bin) gtki.Bin {
	if v == nil {
		return nil
	}
	return &bin{WrapContainerSimple(&v.Container).(*container), v}
}

func WrapBin(v *gtk.Bin, e error) (gtki.Bin, error) {
	return WrapBinSimple(v), e
}

func UnwrapBin(v gtki.Bin) *gtk.Bin {
	if v == nil {
		return nil
	}
	return v.(*bin).Bin
}

func (v *bin) GetChild() gtki.Widget {
	vx, _ := v.Bin.GetChild()
	return Wrap(vx).(gtki.Widget)
}
