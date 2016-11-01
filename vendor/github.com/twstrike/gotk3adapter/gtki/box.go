package gtki

type Box interface {
	Container
	Orientable

	PackEnd(Widget, bool, bool, uint)
	PackStart(Widget, bool, bool, uint)
	SetChildPacking(Widget, bool, bool, uint, PackType)
}

func AssertBox(_ Box) {}
