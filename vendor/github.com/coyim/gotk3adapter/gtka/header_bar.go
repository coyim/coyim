package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type headerBar struct {
	*container
	internal *gtk.HeaderBar
}

func WrapHeaderBarSimple(v *gtk.HeaderBar) gtki.HeaderBar {
	if v == nil {
		return nil
	}
	return &headerBar{WrapContainerSimple(&v.Container).(*container), v}
}

func WrapHeaderBar(v *gtk.HeaderBar, e error) (gtki.HeaderBar, error) {
	return WrapHeaderBarSimple(v), e
}

func UnwrapHeaderBar(v gtki.HeaderBar) *gtk.HeaderBar {
	if v == nil {
		return nil
	}
	return v.(*headerBar).internal
}

func (v *headerBar) SetSubtitle(v1 string) {
	v.internal.SetSubtitle(v1)
}

func (v *headerBar) SetShowCloseButton(v1 bool) {
	v.internal.SetShowCloseButton(v1)
}

func (v *headerBar) GetShowCloseButton() bool {
	return v.internal.GetShowCloseButton()
}
