package pangoa

import (
	"github.com/coyim/gotk3adapter/pangoi"
	"github.com/gotk3/gotk3/pango"
)

type RealPango struct{}

var Real = &RealPango{}

func (*RealPango) AsFontDescription(v interface{}) pangoi.FontDescription {
	if v == nil {
		return nil
	}

	return wrapFontDescriptionSimple(v.(*pango.FontDescription))
}
