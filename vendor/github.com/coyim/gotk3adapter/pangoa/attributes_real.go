package pangoa

import (
	"github.com/coyim/gotk3adapter/pangoi"
	"github.com/gotk3/gotk3/pango"
)

type attrList struct {
	*pango.AttrList
}

func wrapAttrListSimple(p *pango.AttrList) pangoi.AttrList {
	if p == nil {
		return nil
	}

	return &attrList{p}
}

func wrapAttrList(p *pango.AttrList, e error) (pangoi.AttrList, error) {
	return wrapAttrListSimple(p), e
}

func unwrapAttrList(v pangoi.AttrList) *pango.AttrList {
	if v == nil {
		return nil
	}
	return v.(*attrList).AttrList
}
