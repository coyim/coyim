package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

type builder struct {
	*gliba.Object
	internal *gtk.Builder
}

func wrapBuilderSimple(v *gtk.Builder) *builder {
	if v == nil {
		return nil
	}
	return &builder{gliba.WrapObjectSimple(v.Object), v}
}

func wrapBuilder(v *gtk.Builder, e error) (*builder, error) {
	return wrapBuilderSimple(v), e
}

func unwrapBuilder(v gtki.Builder) *gtk.Builder {
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
