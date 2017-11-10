package gtki

type Button interface {
	Bin

	SetImage(Widget)
}

func AssertButton(_ Button) {}
