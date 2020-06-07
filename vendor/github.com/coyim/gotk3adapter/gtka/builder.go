package gtka

import (
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type builder struct {
	*gliba.Object
	internal *gtk.Builder
}

func WrapBuilderSimple(v *gtk.Builder) gtki.Builder {
	if v == nil {
		return nil
	}
	return &builder{gliba.WrapObjectSimple(v.Object), v}
}

func WrapBuilder(v *gtk.Builder, e error) (gtki.Builder, error) {
	return WrapBuilderSimple(v), e
}

func UnwrapBuilder(v gtki.Builder) *gtk.Builder {
	if v == nil {
		return nil
	}
	return v.(*builder).internal
}

func (v *builder) AddFromString(v1 string) error {
	return v.internal.AddFromString(v1)
}

func (v *builder) AddFromResource(v1 string) error {
	return v.internal.AddFromResource(v1)
}

func (v *builder) GetObject(v1 string) (glibi.Object, error) {
	vx1, vx2 := v.internal.GetObject(v1)
	return Wrap(vx1).(glibi.Object), vx2
}

func (v *builder) ConnectSignals(v1 map[string]interface{}) {
	newSignals := make(map[string]interface{})
	for k, vv := range v1 {
		newSignals[k] = gliba.FixupFunction(vv)
	}
	v.internal.ConnectSignals(newSignals)
}
