package gtki

type Label interface {
	Widget

	SetLabel(string)
	SetSelectable(bool)
	SetText(string)
}

func AssertLabel(_ Label) {}
