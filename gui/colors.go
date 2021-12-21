package gui

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type colorValue float64

type rgb struct {
	red   colorValue
	green colorValue
	blue  colorValue
}

type rgba struct {
	*rgb
	alpha colorValue
}

type rgbaGetters interface {
	GetRed() float64
	GetGreen() float64
	GetBlue() float64
}

var colorThemeInsensitiveForeground = rgbFrom(131, 119, 119)

func createColorValueFrom(v uint8) colorValue {
	return colorValue(float64(v) / 255)
}

func (v colorValue) toScaledValue() uint8 {
	return uint8(v * 255)
}

func colorValueFromHex(s string) colorValue {
	value, err := strconv.ParseUint(s, 16, 8)
	if err != nil {
		return colorValue(0)
	}
	return createColorValueFrom(uint8(value))
}

func doubleString(s string) string {
	return s + s
}

// rgbFromHex will return an rgb object from either #xxxxxx or #xxx representation
// it returns nil if parsing fails
func rgbFromHex(spec string) *rgb {
	s := strings.TrimPrefix(spec, "#")
	switch len(s) {
	case 3:
		return &rgb{
			red:   colorValueFromHex(doubleString(s[0:1])),
			green: colorValueFromHex(doubleString(s[1:2])),
			blue:  colorValueFromHex(doubleString(s[2:3])),
		}
	case 6:
		return &rgb{
			red:   colorValueFromHex(s[0:2]),
			green: colorValueFromHex(s[2:4]),
			blue:  colorValueFromHex(s[4:6]),
		}
	}
	return nil
}

func rgbFrom(r, g, b uint8) *rgb {
	return &rgb{
		red:   createColorValueFrom(r),
		green: createColorValueFrom(g),
		blue:  createColorValueFrom(b),
	}
}

func rgbaFrom(r, g, b uint8, a float64) *rgba {
	return &rgba{
		rgb:   rgbFrom(r, g, b),
		alpha: colorValue(a),
	}
}

func rgbFromPercent(r, g, b float64) *rgb {
	return &rgb{
		red:   colorValue(r),
		green: colorValue(g),
		blue:  colorValue(b),
	}
}

func rgbFromGetters(v rgbaGetters) *rgb {
	return rgbFromPercent(v.GetRed(), v.GetGreen(), v.GetBlue())
}

const lightnessThreshold = 0.8

func (r *rgb) isDark() bool {
	return r.lightness() < lightnessThreshold
}

func (r *rgb) lightness() float64 {
	// We are using the formula found in https://en.wikipedia.org/wiki/HSL_and_HSV#From_RGB
	max := math.Max(math.Max(float64(r.red), float64(r.green)), float64(r.blue))
	min := math.Min(math.Min(float64(r.red), float64(r.green)), float64(r.blue))

	return (max + min) / 2
}

type hexColor interface {
	toHex() string
}

type color interface {
	cssColor
	hexColor
}

func (r *rgba) String() string {
	return r.toCSS()
}

func (r *rgb) toHex() string {
	return fmt.Sprintf("#%02x%02x%02x", r.red.toScaledValue(), r.green.toScaledValue(), r.blue.toScaledValue())
}

func (r *rgb) String() string {
	return r.toCSS()
}
