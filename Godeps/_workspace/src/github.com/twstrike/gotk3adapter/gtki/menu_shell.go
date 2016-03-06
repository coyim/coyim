package gtki

type MenuShell interface {
	Container

	Append(MenuItem)
}

func AssertMenuShell(_ MenuShell) {}
