package gtki

type HeaderBar interface {
	Container

	SetSubtitle(string)
}

func AssertHeaderBar(_ HeaderBar) {}
