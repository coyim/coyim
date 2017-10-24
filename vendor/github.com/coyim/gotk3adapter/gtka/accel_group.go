package gtka

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/coyim/gotk3adapter/gtki"
)

type accelGroup struct {
	*gliba.Object
	internal *gtk.AccelGroup
}

func wrapAccelGroupSimple(v *gtk.AccelGroup) *accelGroup {
	if v == nil {
		return nil
	}
	return &accelGroup{gliba.WrapObjectSimple(v.Object), v}
}

func wrapAccelGroup(v *gtk.AccelGroup, e error) (*accelGroup, error) {
	return wrapAccelGroupSimple(v), e
}

func unwrapAccelGroup(v gtki.AccelGroup) *gtk.AccelGroup {
	if v == nil {
		return nil
	}
	return v.(*accelGroup).internal
}

func (v *accelGroup) Connect2(v2 uint, v3 gdki.ModifierType, v4 gtki.AccelFlags, v5 interface{}) {
	v.internal.Connect(v2, gdk.ModifierType(v3), gtk.AccelFlags(v4), gliba.FixupFunction(v5))
}
