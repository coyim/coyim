package gdka

import (
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/gotk3/gotk3/gdk"
)

type rgba struct {
	internal *gdk.RGBA
}

func WrapRgbaSimple(v *gdk.RGBA) gdki.Rgba {
	if v == nil {
		return nil
	}
	return &rgba{v}
}

func WrapRgba(v *gdk.RGBA, e error) (gdki.Rgba, error) {
	return WrapRgbaSimple(v), e
}

func UnwrapRgba(v gdki.Rgba) *gdk.RGBA {
	if v == nil {
		return nil
	}
	return v.(*rgba).internal
}

func (v *rgba) String() string {
	return UnwrapRgba(v).String()
}
