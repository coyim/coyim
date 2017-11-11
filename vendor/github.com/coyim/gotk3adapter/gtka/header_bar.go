package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type headerBar struct {
	*container
	internal *gtk.HeaderBar
}

func wrapHeaderBarSimple(v *gtk.HeaderBar) *headerBar {
	if v == nil {
		return nil
	}
	return &headerBar{wrapContainerSimple(&v.Container), v}
}

func wrapHeaderBar(v *gtk.HeaderBar, e error) (*headerBar, error) {
	return wrapHeaderBarSimple(v), e
}

func unwrapHeaderBar(v gtki.HeaderBar) *gtk.HeaderBar {
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
