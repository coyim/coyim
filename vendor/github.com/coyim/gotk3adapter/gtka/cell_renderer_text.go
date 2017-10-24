package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gtki"
)

type cellRendererText struct {
	*cellRenderer
	internal *gtk.CellRendererText
}

func wrapCellRendererTextSimple(v *gtk.CellRendererText) *cellRendererText {
	if v == nil {
		return nil
	}
	return &cellRendererText{wrapCellRendererSimple(&v.CellRenderer), v}
}

func wrapCellRendererText(v *gtk.CellRendererText, e error) (*cellRendererText, error) {
	return wrapCellRendererTextSimple(v), e
}

func unwrapCellRendererText(v gtki.CellRendererText) *gtk.CellRendererText {
	if v == nil {
		return nil
	}
	return v.(*cellRendererText).internal
}
