package gtki

type Label interface {
	Widget

	GetLabel() string
	SetLabel(string)
	SetSelectable(bool)
	SetText(string)
	SetMarkup(string)
	GetMnemonicKeyval() uint
}

func AssertLabel(_ Label) {}
