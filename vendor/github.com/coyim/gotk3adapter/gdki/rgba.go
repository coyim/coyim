package gdki

type Rgba interface {
	String() string
	GetRed() float64
	GetGreen() float64
	GetBlue() float64
}

func AssertRgba(_ Rgba) {}
