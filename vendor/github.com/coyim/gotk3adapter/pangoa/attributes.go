package pangoa

import (
	"github.com/coyim/gotk3adapter/pangoi"
	"github.com/coyim/gotk3extra"
)

type pangoAttrList struct {
	internal *gotk3extra.PangoAttrList
}

func WrapPangoAttrListSimple(p *gotk3extra.PangoAttrList) pangoi.PangoAttrList {
	if p == nil {
		return nil
	}

	return &pangoAttrList{p}
}

func WrapPangoAttrList(p *gotk3extra.PangoAttrList, e error) (pangoi.PangoAttrList, error) {
	return WrapPangoAttrListSimple(p), e
}

func UnwrapPangoAttrList(v pangoi.PangoAttrList) *gotk3extra.PangoAttrList {
	if v == nil {
		return nil
	}
	return v.(*pangoAttrList).internal
}

func (v *pangoAttrList) GetAttributes() []pangoi.PangoAttribute {
	attributes := []pangoi.PangoAttribute{}
	for _, attr := range v.internal.GetAttributes() {
		attributes = append(attributes, wrapPangoAttributeSimple(attr))
	}
	return attributes
}

func (v *pangoAttrList) InsertPangoAttribute(v1 pangoi.PangoAttribute) {
	v.internal.Insert(unwrapPangoAttribute(v1))
}

func (v *pangoAttrList) Insert(v1 pangoi.Attribute) {
	v.internal.Insert(gotk3extra.PangoAttributeFromReal(unwrapAttribute(v1)))
}
