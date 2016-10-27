package pangoa

import (
	"github.com/gotk3/gotk3/pango"
	"github.com/twstrike/gotk3adapter/pangoi"
)

type fontDescription struct {
	*pango.FontDescription
}

func wrapFontDescriptionSimple(p *pango.FontDescription) pangoi.FontDescription {
	if p == nil {
		return nil
	}

	return &fontDescription{p}
}

func wrapFontDescription(p *pango.FontDescription, e error) (pangoi.FontDescription, error) {
	return wrapFontDescriptionSimple(p), e
}

func unwrapFontDescription(v pangoi.FontDescription) *pango.FontDescription {
	if v == nil {
		return nil
	}
	return v.(*fontDescription).FontDescription
}
