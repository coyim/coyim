package pangoa

import (
	"github.com/coyim/gotk3adapter/pangoi"
	"github.com/coyim/gotk3extra"
)

type pangoAttribute struct {
	internal *gotk3extra.PangoAttribute
}

func wrapPangoAttributeSimple(p *gotk3extra.PangoAttribute) pangoi.PangoAttribute {
	if p == nil {
		return nil
	}

	return &pangoAttribute{p}
}

func wrapPangoAttribute(p *gotk3extra.PangoAttribute, e error) (pangoi.PangoAttribute, error) {
	return wrapPangoAttributeSimple(p), e
}

func unwrapPangoAttribute(v pangoi.PangoAttribute) *gotk3extra.PangoAttribute {
	if v == nil {
		return nil
	}

	return v.(*pangoAttribute).internal
}

func (v *pangoAttribute) SetStartIndex(v1 int) {
	v.internal.SetStartIndex(uint(v1))
}

func (v *pangoAttribute) SetEndIndex(v1 int) {
	v.internal.SetEndIndex(uint(v1))
}
