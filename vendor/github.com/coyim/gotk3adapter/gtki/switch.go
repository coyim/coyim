package gtki

type Switch interface {
	Widget

	GetActive() bool
	SetActive(bool)
}

func AssertSwitch(_ Switch) {}
