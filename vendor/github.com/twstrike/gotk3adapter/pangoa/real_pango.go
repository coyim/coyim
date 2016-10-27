package pangoa

import (
	"github.com/gotk3/gotk3/pango"
	"github.com/twstrike/gotk3adapter/pangoi"
)

type RealPango struct{}

var Real = &RealPango{}

func (*RealPango) AsFontDescription(v interface{}) pangoi.FontDescription {
	if v == nil {
		return nil
	}

	return wrapFontDescriptionSimple(v.(*pango.FontDescription))
}
