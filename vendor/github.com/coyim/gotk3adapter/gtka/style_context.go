package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/coyim/gotk3adapter/gtki"
)

type styleContext struct {
	*gliba.Object
	internal *gtk.StyleContext
}

func wrapStyleContextSimple(v *gtk.StyleContext) *styleContext {
	if v == nil {
		return nil
	}
	return &styleContext{gliba.WrapObjectSimple(v.Object), v}
}

func wrapStyleContext(v *gtk.StyleContext, e error) (*styleContext, error) {
	return wrapStyleContextSimple(v), e
}

func unwrapStyleContext(v gtki.StyleContext) *gtk.StyleContext {
	if v == nil {
		return nil
	}
	return v.(*styleContext).internal
}

func (v *styleContext) AddClass(v1 string) {
	v.internal.AddClass(v1)
}

func (v *styleContext) AddProvider(v1 gtki.StyleProvider, v2 uint) {
	v.internal.AddProvider(unwrapStyleProvider(v1), v2)
}

func (v *styleContext) GetProperty2(v1 string, v2 gtki.StateFlags) (interface{}, error) {
	vx1, vx2 := v.internal.GetProperty(v1, gtk.StateFlags(v2))
	return gliba.WrapAllGuard(vx1), vx2
}
