package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type entry struct {
	*widget
	internal *gtk.Entry
}

func wrapEntrySimple(v *gtk.Entry) *entry {
	if v == nil {
		return nil
	}
	return &entry{wrapWidgetSimple(&v.Widget), v}
}

func wrapEntry(v *gtk.Entry, e error) (*entry, error) {
	return wrapEntrySimple(v), e
}

func unwrapEntry(v gtki.Entry) *gtk.Entry {
	if v == nil {
		return nil
	}
	return v.(*entry).internal
}

func (v *entry) GetText() (string, error) {
	return v.internal.GetText()
}

func (v *entry) SetHasFrame(v1 bool) {
	v.internal.SetHasFrame(v1)
}

func (v *entry) SetVisibility(v1 bool) {
	v.internal.SetVisibility(v1)
}

func (v *entry) SetText(v1 string) {
	v.internal.SetText(v1)
}

func (v *entry) SetEditable(v1 bool) {
	v.internal.SetEditable(v1)
}

func (v *entry) SetWidthChars(v1 int) {
	v.internal.SetWidthChars(v1)
}

func (v *entry) GetAlignment() float32 {
	return v.internal.GetAlignment()
}
func (v *entry) SetAlignment(v1 float32) {
	v.internal.SetAlignment(v1)
}

func (v *entry) SetPosition(p int) {
	v.internal.SetPosition(p)
}

func (v *entry) GetPosition() int {
	return v.internal.GetPosition()
}
