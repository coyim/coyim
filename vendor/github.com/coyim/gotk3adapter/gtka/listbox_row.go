package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type listBoxRow struct {
	*bin
	internal *gtk.ListBoxRow
}

func WrapListBoxRowSimple(v *gtk.ListBoxRow) gtki.ListBoxRow {
	if v == nil {
		return nil
	}
	return &listBoxRow{WrapBinSimple(&v.Bin).(*bin), v}
}

func WrapListBoxRow(v *gtk.ListBoxRow, e error) (gtki.ListBoxRow, error) {
	return WrapListBoxRowSimple(v), e
}

func UnwrapListBoxRow(v gtki.ListBoxRow) *gtk.ListBoxRow {
	if v == nil {
		return nil
	}
	return v.(*listBoxRow).internal
}

func (v *listBoxRow) GetIndex() int {
	return v.internal.GetIndex()
}
