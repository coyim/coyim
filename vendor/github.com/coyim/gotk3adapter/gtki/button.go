package gtki

type Button interface {
	Bin

	SetImage(Widget)
	GetLabel() (string, error)
	SetLabel(string)
}

func AssertButton(_ Button) {}
