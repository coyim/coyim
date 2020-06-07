package gtka

import (
	"github.com/coyim/gotk3adapter/gdka"
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type image struct {
	*widget
	internal *gtk.Image
}

func WrapImageSimple(v *gtk.Image) gtki.Image {
	if v == nil {
		return nil
	}

	return &image{WrapWidgetSimple(&v.Widget).(*widget), v}
}

func WrapImage(v *gtk.Image, e error) (gtki.Image, error) {
	return WrapImageSimple(v), e
}

func UnwrapImage(v gtki.Image) *gtk.Image {
	if v == nil {
		return nil
	}
	return v.(*image).internal
}

func (v *image) SetFromIconName(v1 string, v2 gtki.IconSize) {
	v.internal.SetFromIconName(v1, gtk.IconSize(v2))
}

func (v *image) Clear() {
	v.internal.Clear()
}

func (v *image) SetFromPixbuf(pb gdki.Pixbuf) {
	v.internal.SetFromPixbuf(gdka.UnwrapPixbuf(pb))
}
