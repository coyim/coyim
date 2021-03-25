package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type listBox struct {
	*container
	internal *gtk.ListBox
}

func WrapListBoxSimple(v *gtk.ListBox) gtki.ListBox {
	if v == nil {
		return nil
	}
	return &listBox{WrapContainerSimple(&v.Container).(*container), v}
}

func WrapListBox(v *gtk.ListBox, e error) (gtki.ListBox, error) {
	return WrapListBoxSimple(v), e
}

func UnwrapListBox(v gtki.ListBox) *gtk.ListBox {
	if v == nil {
		return nil
	}
	return v.(*listBox).internal
}

func (v *listBox) SelectRow(v1 gtki.ListBoxRow) {
	v.internal.SelectRow(v1.(*listBoxRow).internal)
}

func (v *listBox) GetRowAtIndex(index int) gtki.ListBoxRow {
	return WrapListBoxRowSimple(v.internal.GetRowAtIndex(index))
}
