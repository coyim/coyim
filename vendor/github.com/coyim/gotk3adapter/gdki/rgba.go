package gdki

type Rgba interface {
	String() string
}

func AssertRgba(_ Rgba) {}
