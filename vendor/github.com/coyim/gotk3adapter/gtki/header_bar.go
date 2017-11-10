package gtki

type HeaderBar interface {
	Container

	SetSubtitle(string)
	SetShowCloseButton(bool)
	GetShowCloseButton() bool
}

func AssertHeaderBar(_ HeaderBar) {}
