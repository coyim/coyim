package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gtki"
)

type box struct {
	*container
	internal *gtk.Box
}

func wrapBoxSimple(v *gtk.Box) *box {
	if v == nil {
		return nil
	}
	return &box{wrapContainerSimple(&v.Container), v}
}

func wrapBox(v *gtk.Box, e error) (*box, error) {
	return wrapBoxSimple(v), e
}

func unwrapBox(v gtki.Box) *gtk.Box {
	if v == nil {
		return nil
	}
	return v.(*box).internal
}

func (v *box) PackEnd(v1 gtki.Widget, v2, v3 bool, v4 uint) {
	v.internal.PackEnd(unwrapWidget(v1), v2, v3, v4)
}

func (v *box) PackStart(v1 gtki.Widget, v2, v3 bool, v4 uint) {
	v.internal.PackStart(unwrapWidget(v1), v2, v3, v4)
}

func (v *box) SetChildPacking(v1 gtki.Widget, v2, v3 bool, v4 uint, v5 gtki.PackType) {
	v.internal.SetChildPacking(unwrapWidget(v1), v2, v3, v4, gtk.PackType(v5))
}

func (v *box) GetOrientation() gtki.Orientation {
	return wrapOrientation(v.internal.GetOrientation())
}

func (v *box) SetOrientation(o gtki.Orientation) {
	v.internal.SetOrientation(unwrapOrientation(o))
}
