package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gtki"
)

type grid struct {
	*container
	internal *gtk.Grid
}

func wrapGridSimple(v *gtk.Grid) *grid {
	if v == nil {
		return nil
	}
	return &grid{wrapContainerSimple(&v.Container), v}
}

func wrapGrid(v *gtk.Grid, e error) (*grid, error) {
	return wrapGridSimple(v), e
}

func unwrapGrid(v gtki.Grid) *gtk.Grid {
	if v == nil {
		return nil
	}
	return v.(*grid).internal
}

func (v *grid) Attach(v1 gtki.Widget, v2, v3, v4, v5 int) {
	v.internal.Attach(unwrapWidget(v1), v2, v3, v4, v5)
}
