package gtki

type CheckMenuItem interface {
	MenuItem

	GetActive() bool
	SetActive(bool)
}

func AssertCheckMenuItem(_ CheckMenuItem) {}
