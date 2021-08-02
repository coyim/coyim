package gtka

import (
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type sizeGroup struct {
	*gliba.Object
	internal *gtk.SizeGroup
}

func WrapSizeGroupSimple(v *gtk.SizeGroup) gtki.SizeGroup {
	if v == nil {
		return nil
	}
	return &sizeGroup{gliba.WrapObjectSimple(v.Object), v}
}

func WrapSizeGroup(v *gtk.SizeGroup, e error) (gtki.SizeGroup, error) {
	return WrapSizeGroupSimple(v), e
}

func UnwrapSizeGroup(v gtki.SizeGroup) *gtk.SizeGroup {
	if v == nil {
		return nil
	}
	return v.(*sizeGroup).internal
}

func (v *sizeGroup) SetMode(v1 gtki.SizeGroupMode) {
	v.internal.SetMode(gtk.SizeGroupMode(v1))
}
