package gdki

type Rgba interface {
	String() string

	GetRed() float64
	GetGreen() float64
	GetBlue() float64
	GetAlpha() float64

	SetRed(float64)
	SetGreen(float64)
	SetBlue(float64)
	SetAlpha(float64)

	Colors() (r, g, b, a float64)
	SetColors(r, g, b, a float64)

	Parse(string) bool
}

func AssertRgba(_ Rgba) {}
