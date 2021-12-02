package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type expander struct {
	*bin
	internal *gtk.Expander
}

func WrapExpanderSimple(v *gtk.Expander) gtki.Expander {
	if v == nil {
		return nil
	}
	return &expander{WrapBinSimple(&v.Bin).(*bin), v}
}

func WrapExpander(v *gtk.Expander, e error) (gtki.Expander, error) {
	return WrapExpanderSimple(v), e
}

func UnwrapExpander(v gtki.Expander) *gtk.Expander {
	if v == nil {
		return nil
	}
	return v.(*expander).internal
}
