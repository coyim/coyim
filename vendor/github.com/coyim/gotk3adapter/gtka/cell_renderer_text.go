package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type cellRendererText struct {
	*cellRenderer
	internal *gtk.CellRendererText
}

func WrapCellRendererTextSimple(v *gtk.CellRendererText) gtki.CellRendererText {
	if v == nil {
		return nil
	}
	return &cellRendererText{WrapCellRendererSimple(&v.CellRenderer).(*cellRenderer), v}
}

func WrapCellRendererText(v *gtk.CellRendererText, e error) (gtki.CellRendererText, error) {
	return WrapCellRendererTextSimple(v), e
}

func UnwrapCellRendererText(v gtki.CellRendererText) *gtk.CellRendererText {
	if v == nil {
		return nil
	}
	return v.(*cellRendererText).internal
}
