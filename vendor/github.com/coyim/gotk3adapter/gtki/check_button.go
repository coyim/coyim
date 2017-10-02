package gtki

type CheckButton interface {
	ToggleButton
}

func AssertCheckButton(_ CheckButton) {}
