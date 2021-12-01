package gdki

type Rectangle interface {
	GetY() int
}

func AssertRectangle(_ Rectangle) {}
