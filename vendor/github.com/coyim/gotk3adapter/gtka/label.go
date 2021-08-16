package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/coyim/gotk3adapter/pangoa"
	"github.com/coyim/gotk3adapter/pangoi"
	"github.com/gotk3/gotk3/gtk"
)

type label struct {
	*widget
	internal *gtk.Label
}

func WrapLabelSimple(v *gtk.Label) gtki.Label {
	if v == nil {
		return nil
	}
	return &label{WrapWidgetSimple(&v.Widget).(*widget), v}
}

func WrapLabel(v *gtk.Label, e error) (gtki.Label, error) {
	return WrapLabelSimple(v), e
}

func UnwrapLabel(v gtki.Label) *gtk.Label {
	if v == nil {
		return nil
	}
	return v.(*label).internal
}

func (v *label) GetLabel() string {
	return v.internal.GetLabel()
}

func (v *label) SetLabel(v1 string) {
	v.internal.SetLabel(v1)
}

func (v *label) SetText(v1 string) {
	v.internal.SetText(v1)
}

func (v *label) SetMarkup(v1 string) {
	v.internal.SetMarkup(v1)
}

func (v *label) SetSelectable(v1 bool) {
	v.internal.SetSelectable(v1)
}

func (v *label) GetMnemonicKeyval() uint {
	return v.internal.GetMnemonicKeyval()
}

func (v *label) SetAttributes(v1 pangoi.AttrList) {
	v.internal.SetAttributes(pangoa.UnwrapAttrList(v1))
}

func (v *label) GetAttributes() (pangoi.AttrList, error) {
	atts, err := v.internal.GetAttributes()
	return pangoa.WrapAttrList(atts, err)
}
