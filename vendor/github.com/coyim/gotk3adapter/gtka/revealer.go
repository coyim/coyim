package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type revealer struct {
	*bin
	internal *gtk.Revealer
}

func WrapRevealerSimple(v *gtk.Revealer) gtki.Revealer {
	if v == nil {
		return nil
	}
	return &revealer{WrapBinSimple(&v.Bin).(*bin), v}
}

func WrapRevealer(v *gtk.Revealer, e error) (gtki.Revealer, error) {
	return WrapRevealerSimple(v), e
}

func UnwrapRevealer(v gtki.Revealer) *gtk.Revealer {
	if v == nil {
		return nil
	}
	return v.(*revealer).internal
}

func (v *revealer) SetRevealChild(revealChild bool) {
	v.internal.SetRevealChild(revealChild)
}
