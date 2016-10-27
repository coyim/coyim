package gdki

type EventButton interface {
	Event

	Button() uint
	Time() uint32
	X() float64
	Y() float64
}

func AssertEventButton(_ EventButton) {}
