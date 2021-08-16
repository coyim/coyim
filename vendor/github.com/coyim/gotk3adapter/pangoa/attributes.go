package pangoa

import (
	"github.com/coyim/gotk3adapter/pangoi"
	"github.com/gotk3/gotk3/pango"
)

type attrList struct {
	internal *pango.AttrList
}

// WrapAttrListSimple wraps a *pango.AttrList to a pangoi.AttrList
func WrapAttrListSimple(p *pango.AttrList) pangoi.AttrList {
	if p == nil {
		return nil
	}

	return &attrList{p}
}

// WrapAttrList wraps a *pango.AttrList to a pangoi.AttrList
func WrapAttrList(p *pango.AttrList, e error) (pangoi.AttrList, error) {
	return WrapAttrListSimple(p), e
}

// UnwrapAttrList unwraps a pangoi.AttrList to a *pango.AttrList
func UnwrapAttrList(v pangoi.AttrList) *pango.AttrList {
	if v == nil {
		return nil
	}

	return v.(*attrList).internal
}

func (v *attrList) Insert(v1 pangoi.Attribute) {
	v.internal.Insert(UnwrapAttribute(v1))
}

func (v *attrList) GetAttributes() []pangoi.Attribute {
	attributes := []pangoi.Attribute{}

	if list, err := v.internal.GetAttributes(); err == nil {
		list.Foreach(func(v interface{}) {
			if attr, ok := v.(*pango.Attribute); ok {
				attributes = append(attributes, WrapAttributeSimple(attr))
			}
		})
	}

	return attributes
}
