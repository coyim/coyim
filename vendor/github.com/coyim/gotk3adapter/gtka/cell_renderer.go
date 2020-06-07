package gtka

import (
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type cellRenderer struct {
	*gliba.Object
	*gtk.CellRenderer
}

type asCellRenderer interface {
	toCellRenderer() *cellRenderer
}

func (v *cellRenderer) toCellRenderer() *cellRenderer {
	return v
}

func WrapCellRendererSimple(v *gtk.CellRenderer) gtki.CellRenderer {
	if v == nil {
		return nil
	}
	return &cellRenderer{gliba.WrapObjectSimple(v.Object), v}
}

func WrapCellRenderer(v *gtk.CellRenderer, e error) (gtki.CellRenderer, error) {
	return WrapCellRendererSimple(v), e
}

func UnwrapCellRenderer(v gtki.CellRenderer) *gtk.CellRenderer {
	if v == nil {
		return nil
	}
	return v.(asCellRenderer).toCellRenderer().CellRenderer
}
