package gtka

import (
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtki"
)

type container struct {
	*widget
	*gtk.Container
}

func wrapContainerSimple(v *gtk.Container) *container {
	if v == nil {
		return nil
	}
	return &container{wrapWidgetSimple(&v.Widget), v}
}

func wrapContainer(v *gtk.Container, e error) (*container, error) {
	return wrapContainerSimple(v), e
}

func unwrapContainer(v gtki.Container) *gtk.Container {
	if v == nil {
		return nil
	}
	return v.(*container).Container
}

func (v *container) Add(v2 gtki.Widget) {
	v.Container.Add(unwrapWidget(v2))
}

func (v *container) Remove(v2 gtki.Widget) {
	v.Container.Remove(unwrapWidget(v2))
}
