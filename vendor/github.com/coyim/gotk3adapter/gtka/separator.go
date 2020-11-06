package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type separator struct {
	*widget
	internal *gtk.Separator
}

type asSeparator interface {
	toSeparator() *separator
}

func (v *separator) toSeparator() *separator {
	return v
}

func WrapSeparatorSimple(v *gtk.Separator) gtki.Separator {
	if v == nil {
		return nil
	}
	return &separator{WrapWidgetSimple(&v.Widget).(*widget), v}
}

func WrapSeparator(v *gtk.Separator, e error) (gtki.Separator, error) {
	return WrapSeparatorSimple(v), e
}

func UnwrapSeparator(v gtki.Separator) *gtk.Separator {
	if v == nil {
		return nil
	}
	return v.(asSeparator).toSeparator().internal
}

func (*RealGtk) SeparatorNew(v1 gtki.Orientation) (gtki.Separator, error) {
	return WrapSeparator(gtk.SeparatorNew(gtk.Orientation(v1)))
}
