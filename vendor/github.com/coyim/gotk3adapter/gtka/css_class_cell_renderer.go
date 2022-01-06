package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/coyim/gotk3extra"
)

type cssClassCellRenderer struct {
	*cellRenderer
	internal *gotk3extra.CSSClassCellRenderer
}

func WrapCSSClassCellRendererSimple(v *gotk3extra.CSSClassCellRenderer) gtki.CSSClassCellRenderer {
	if v == nil {
		return nil
	}
	return &cssClassCellRenderer{WrapCellRendererSimple(&v.CellRenderer).(*cellRenderer), v}
}

func WrapCSSClassCellRenderer(v *gotk3extra.CSSClassCellRenderer, e error) (gtki.CSSClassCellRenderer, error) {
	return WrapCSSClassCellRendererSimple(v), e
}

func UnwrapCSSClassCellRenderer(v gtki.CSSClassCellRenderer) *gotk3extra.CSSClassCellRenderer {
	if v == nil {
		return nil
	}
	return v.(*cssClassCellRenderer).internal
}

func (v *cssClassCellRenderer) SetReal(real gtki.CellRenderer) {
	v.internal.SetReal(UnwrapCellRenderer(real))
}
