package gtki

type Label interface {
	Widget

	GetLabel() string
	SetLabel(string)
	SetSelectable(bool)
	SetText(string)
}

func AssertLabel(_ Label) {}
