package gdka

import (
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/gotk3/gotk3/gdk"
)

type rectangle struct {
	internal *gdk.Rectangle
}

func WrapRectangleSimple(v *gdk.Rectangle) gdki.Rectangle {
	if v == nil {
		return nil
	}
	return &rectangle{v}
}

func WrapRectangle(v *gdk.Rectangle, e error) (gdki.Rectangle, error) {
	return WrapRectangleSimple(v), e
}

func UnwrapRectangle(v gdki.Rectangle) *gdk.Rectangle {
	if v == nil {
		return nil
	}
	return v.(*rectangle).internal
}

func (v *rectangle) GetY() int {
	return v.internal.GetY()
}
