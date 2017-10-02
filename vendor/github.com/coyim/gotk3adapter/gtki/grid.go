package gtki

type Grid interface {
	Container

	Attach(Widget, int, int, int, int)
}

func AssertGrid(_ Grid) {}
