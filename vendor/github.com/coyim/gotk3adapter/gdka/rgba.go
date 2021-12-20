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
	return v.internal.String()
}

func (v *rgba) GetRed() float64 {
	return v.internal.GetRed()
}

func (v *rgba) GetGreen() float64 {
	return v.internal.GetGreen()
}

func (v *rgba) GetBlue() float64 {
	return v.internal.GetBlue()
}

func (v *rgba) GetAlpha() float64 {
	return v.internal.GetAlpha()
}

func (v *rgba) SetRed(c float64) {
	v.internal.SetRed(c)
}

func (v *rgba) SetGreen(c float64) {
	v.internal.SetGreen(c)
}

func (v *rgba) SetBlue(c float64) {
	v.internal.SetBlue(c)
}

func (v *rgba) SetAlpha(c float64) {
	v.internal.SetAlpha(c)
}

func (v *rgba) Colors() (r, g, b, a float64) {
	f := v.internal.Floats()
	return f[0], f[1], f[2], f[3]
}

func (v *rgba) SetColors(r, g, b, a float64) {
	v.internal.SetColors(r, g, b, a)
}

func (v *rgba) Parse(spec string) bool {
	return v.internal.Parse(spec)
}
