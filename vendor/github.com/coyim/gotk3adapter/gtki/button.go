package gtki

type Button interface {
	Bin

	SetImage(Widget)
	GetLabel() (string, error)
}

func AssertButton(_ Button) {}
