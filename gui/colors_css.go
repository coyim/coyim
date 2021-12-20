package gui

import "fmt"

type cssColor interface {
	toCSS() string
}

var (
	colorNone                       = cssColorReferenceFrom("none")
	colorTransparent                = cssColorReferenceFrom("transparent")
	colorThemeBase                  = cssColorReferenceFrom("@theme_base_color")
	colorThemeBackground            = cssColorReferenceFrom("@theme_bg_color")
	colorThemeForeground            = cssColorReferenceFrom("@theme_fg_color")
	colorThemeInsensitiveBackground = cssColorReferenceFrom("@insensitive_bg_color")
)

type cssColorReference struct {
	ref string
}

func cssColorReferenceFrom(ref string) *cssColorReference {
	return &cssColorReference{
		ref: ref,
	}
}

func (r *rgba) toCSS() string {
	return fmt.Sprintf("rgba(%d, %d, %d, %f)",
		r.red.toScaledValue(),
		r.green.toScaledValue(),
		r.blue.toScaledValue(),
		float64(r.alpha),
	)
}

func (r *rgb) toCSS() string {
	return fmt.Sprintf("rgb(%d, %d, %d)", r.red.toScaledValue(), r.green.toScaledValue(), r.blue.toScaledValue())
}

func (r *cssColorReference) String() string {
	return r.toCSS()
}

func (r *cssColorReference) toCSS() string {
	return r.ref
}
