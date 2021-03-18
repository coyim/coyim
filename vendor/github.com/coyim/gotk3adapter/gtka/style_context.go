package gtka

import (
	"github.com/coyim/gotk3adapter/gdka"
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type styleContext struct {
	*gliba.Object
	internal *gtk.StyleContext
}

func WrapStyleContextSimple(v *gtk.StyleContext) gtki.StyleContext {
	if v == nil {
		return nil
	}
	return &styleContext{gliba.WrapObjectSimple(v.Object), v}
}

func WrapStyleContext(v *gtk.StyleContext, e error) (gtki.StyleContext, error) {
	return WrapStyleContextSimple(v), e
}

func UnwrapStyleContext(v gtki.StyleContext) *gtk.StyleContext {
	if v == nil {
		return nil
	}
	return v.(*styleContext).internal
}

func (v *styleContext) AddClass(v1 string) {
	v.internal.AddClass(v1)
}

func (v *styleContext) RemoveClass(v1 string) {
	v.internal.RemoveClass(v1)
}

func (v *styleContext) AddProvider(v1 gtki.StyleProvider, v2 uint) {
	v.internal.AddProvider(UnwrapStyleProvider(v1), v2)
}

func (v *styleContext) GetScreen() (gdki.Screen, error) {
	return gdka.WrapScreen(v.internal.GetScreen())
}

func (v *styleContext) GetProperty2(v1 string, v2 gtki.StateFlags) (interface{}, error) {
	vx1, vx2 := v.internal.GetProperty(v1, gtk.StateFlags(v2))
	return gliba.WrapAllGuard(vx1), vx2
}
