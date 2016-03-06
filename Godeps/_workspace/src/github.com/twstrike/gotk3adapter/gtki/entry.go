package gtki

type Entry interface {
	Widget
	Editable

	GetText() (string, error)
	SetHasFrame(bool)
	SetText(string)
	SetVisibility(bool)
}

func AssertEntry(_ Entry) {}
