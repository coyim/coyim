package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gtki"
)

type label struct {
	*widget
	internal *gtk.Label
}

func wrapLabelSimple(v *gtk.Label) *label {
	if v == nil {
		return nil
	}
	return &label{wrapWidgetSimple(&v.Widget), v}
}

func wrapLabel(v *gtk.Label, e error) (*label, error) {
	return wrapLabelSimple(v), e
}

func unwrapLabel(v gtki.Label) *gtk.Label {
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
