package gtka

import (
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type accelGroup struct {
	*gliba.Object
	internal *gtk.AccelGroup
}

func WrapAccelGroupSimple(v *gtk.AccelGroup) gtki.AccelGroup {
	if v == nil {
		return nil
	}
	return &accelGroup{gliba.WrapObjectSimple(v.Object), v}
}

func WrapAccelGroup(v *gtk.AccelGroup, e error) (gtki.AccelGroup, error) {
	return WrapAccelGroupSimple(v), e
}

func UnwrapAccelGroup(v gtki.AccelGroup) *gtk.AccelGroup {
	if v == nil {
		return nil
	}
	return v.(*accelGroup).internal
}

func (v *accelGroup) Connect2(v2 uint, v3 gdki.ModifierType, v4 gtki.AccelFlags, v5 interface{}) {
	v.internal.Connect(v2, gdk.ModifierType(v3), gtk.AccelFlags(v4), gliba.FixupFunction(v5))
}
