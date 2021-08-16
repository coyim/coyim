package pangoa

import (
	"github.com/coyim/gotk3adapter/pangoi"
	"github.com/gotk3/gotk3/pango"
)

type attribute struct {
	internal *pango.Attribute
}

// WrapAttributeSimple wraps a *pango.Attribute to a pangoi.Attribute
func WrapAttributeSimple(p *pango.Attribute) pangoi.Attribute {
	if p == nil {
		return nil
	}

	return &attribute{p}
}

// WrapAttribute wraps a *pango.Attribute to a pangoi.Attribute
func WrapAttribute(p *pango.Attribute, e error) (pangoi.Attribute, error) {
	return WrapAttributeSimple(p), e
}

// UnwrapAttribute unwraps a pangoi.Attribute to a *pango.Attribute
func UnwrapAttribute(v pangoi.Attribute) *pango.Attribute {
	if v == nil {
		return nil
	}

	return v.(*attribute).internal
}

func (v *attribute) SetStartIndex(v1 int) {
	v.internal.SetStartIndex(uint(v1))
}

func (v *attribute) SetEndIndex(v1 int) {
	v.internal.SetEndIndex(uint(v1))
}
