package gtki

type Entry interface {
	Widget
	Editable

	GetText() (string, error)
	SetHasFrame(bool)
	SetText(string)
	GetVisibility() bool
	SetVisibility(bool)
	SetWidthChars(int)
	GetAlignment() float32
	SetAlignment(float32)
	SetCompletion(EntryCompletion)
	SetPlaceholderText(string)
}

func AssertEntry(_ Entry) {}
