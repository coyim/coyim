package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gtki"
)

type revealer struct {
	*bin
	internal *gtk.Revealer
}

func wrapRevealerSimple(v *gtk.Revealer) *revealer {
	if v == nil {
		return nil
	}
	return &revealer{wrapBinSimple(&v.Bin), v}
}

func wrapRevealer(v *gtk.Revealer, e error) (*revealer, error) {
	return wrapRevealerSimple(v), e
}

func unwrapRevealer(v gtki.Revealer) *gtk.Revealer {
	if v == nil {
		return nil
	}
	return v.(*revealer).internal
}
