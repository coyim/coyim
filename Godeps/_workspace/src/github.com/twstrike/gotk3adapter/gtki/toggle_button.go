package gtki

type ToggleButton interface {
	Button

	GetActive() bool
	SetActive(bool)
}

func AssertToggleButton(_ ToggleButton) {}
