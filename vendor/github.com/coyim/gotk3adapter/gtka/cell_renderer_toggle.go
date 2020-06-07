package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type cellRendererToggle struct {
	*cellRenderer
	internal *gtk.CellRendererToggle
}

func WrapCellRendererToggleSimple(v *gtk.CellRendererToggle) gtki.CellRendererToggle {
	if v == nil {
		return nil
	}
	return &cellRendererToggle{WrapCellRendererSimple(&v.CellRenderer).(*cellRenderer), v}
}

func WrapCellRendererToggle(v *gtk.CellRendererToggle, e error) (gtki.CellRendererToggle, error) {
	return WrapCellRendererToggleSimple(v), e
}

func UnwrapCellRendererToggle(v gtki.CellRendererToggle) *gtk.CellRendererToggle {
	if v == nil {
		return nil
	}
	return v.(*cellRendererToggle).internal
}
