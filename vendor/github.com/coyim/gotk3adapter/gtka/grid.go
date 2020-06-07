package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type grid struct {
	*container
	internal *gtk.Grid
}

func WrapGridSimple(v *gtk.Grid) gtki.Grid {
	if v == nil {
		return nil
	}
	return &grid{WrapContainerSimple(&v.Container).(*container), v}
}

func WrapGrid(v *gtk.Grid, e error) (gtki.Grid, error) {
	return WrapGridSimple(v), e
}

func UnwrapGrid(v gtki.Grid) *gtk.Grid {
	if v == nil {
		return nil
	}
	return v.(*grid).internal
}

func (v *grid) Attach(v1 gtki.Widget, v2, v3, v4, v5 int) {
	v.internal.Attach(UnwrapWidget(v1), v2, v3, v4, v5)
}
