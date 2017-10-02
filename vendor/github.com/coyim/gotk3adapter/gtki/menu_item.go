package gtki

type MenuItem interface {
	Bin

	GetLabel() string
	SetLabel(string)
	SetSubmenu(Widget)
}

func AssertMenuItem(_ MenuItem) {}
