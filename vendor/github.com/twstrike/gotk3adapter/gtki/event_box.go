package gtki

type EventBox interface {
	Bin

	SetAboveChild(bool)
	GetAboveChild() bool
	SetVisibleWindow(bool)
	GetVisibleWindow() bool
}

func AssertEventBox(_ EventBox) {}
