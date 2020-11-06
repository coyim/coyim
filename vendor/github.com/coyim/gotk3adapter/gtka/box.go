package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type box struct {
	*container
	internal *gtk.Box
}

func WrapBoxSimple(v *gtk.Box) gtki.Box {
	if v == nil {
		return nil
	}
	return &box{WrapContainerSimple(&v.Container).(*container), v}
}

func WrapBox(v *gtk.Box, e error) (gtki.Box, error) {
	return WrapBoxSimple(v), e
}

func UnwrapBox(v gtki.Box) *gtk.Box {
	if v == nil {
		return nil
	}
	return v.(*box).internal
}

func (*RealGtk) BoxNew(v1 gtki.Orientation, v2 int) (gtki.Box, error) {
	return WrapBox(gtk.BoxNew(gtk.Orientation(v1), v2))
}

func (v *box) PackEnd(v1 gtki.Widget, v2, v3 bool, v4 uint) {
	v.internal.PackEnd(UnwrapWidget(v1), v2, v3, v4)
}

func (v *box) PackStart(v1 gtki.Widget, v2, v3 bool, v4 uint) {
	v.internal.PackStart(UnwrapWidget(v1), v2, v3, v4)
}

func (v *box) SetChildPacking(v1 gtki.Widget, v2, v3 bool, v4 uint, v5 gtki.PackType) {
	v.internal.SetChildPacking(UnwrapWidget(v1), v2, v3, v4, gtk.PackType(v5))
}

func (v *box) GetOrientation() gtki.Orientation {
	return WrapOrientation(v.internal.GetOrientation())
}

func (v *box) SetOrientation(o gtki.Orientation) {
	v.internal.SetOrientation(UnwrapOrientation(o))
}
