package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type cellRendererToggle struct {
	*cellRenderer
	internal *gtk.CellRendererToggle
}

func wrapCellRendererToggleSimple(v *gtk.CellRendererToggle) *cellRendererToggle {
	if v == nil {
		return nil
	}
	return &cellRendererToggle{wrapCellRendererSimple(&v.CellRenderer), v}
}

func wrapCellRendererToggle(v *gtk.CellRendererToggle, e error) (*cellRendererToggle, error) {
	return wrapCellRendererToggleSimple(v), e
}

func unwrapCellRendererToggle(v gtki.CellRendererToggle) *gtk.CellRendererToggle {
	if v == nil {
		return nil
	}
	return v.(*cellRendererToggle).internal
}
