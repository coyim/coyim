package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type container struct {
	*widget
	*gtk.Container
}

func WrapContainerSimple(v *gtk.Container) gtki.Container {
	if v == nil {
		return nil
	}
	return &container{WrapWidgetSimple(&v.Widget).(*widget), v}
}

func WrapContainer(v *gtk.Container, e error) (gtki.Container, error) {
	return WrapContainerSimple(v), e
}

func UnwrapContainer(v gtki.Container) *gtk.Container {
	if v == nil {
		return nil
	}
	return v.(*container).Container
}

func (v *container) Add(v2 gtki.Widget) {
	v.Container.Add(UnwrapWidget(v2))
}

func (v *container) Remove(v2 gtki.Widget) {
	v.Container.Remove(UnwrapWidget(v2))
}

func (v *container) SetBorderWidth(v1 uint) {
	v.Container.SetBorderWidth(v1)
}
