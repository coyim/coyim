package gtki

type Container interface {
	Widget

	Add(Widget)
	Remove(Widget)
	SetBorderWidth(uint)
}
