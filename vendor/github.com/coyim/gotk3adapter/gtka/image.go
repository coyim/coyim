package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gdka"
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/gtki"
)

type image struct {
	*widget
	internal *gtk.Image
}

func wrapImageSimple(v *gtk.Image) *image {
	if v == nil {
		return nil
	}

	return &image{wrapWidgetSimple(&v.Widget), v}
}

func wrapImage(v *gtk.Image, e error) (*image, error) {
	return wrapImageSimple(v), e
}

func unwrapImage(v gtki.Image) *gtk.Image {
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
