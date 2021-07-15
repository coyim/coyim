package pangoa

import (
	"github.com/coyim/gotk3adapter/pangoi"
	"github.com/gotk3/gotk3/pango"
)

type attribute struct {
	internal *pango.Attribute
}

func wrapAttributeSimple(p *pango.Attribute) pangoi.Attribute {
	if p == nil {
		return nil
	}

	return &attribute{p}
}

func wrapAttribute(p *pango.Attribute, e error) (pangoi.Attribute, error) {
	return wrapAttributeSimple(p), e
}

func unwrapAttribute(v pangoi.Attribute) *pango.Attribute {
	if v == nil {
		return nil
	}

	return v.(*attribute).internal
}
