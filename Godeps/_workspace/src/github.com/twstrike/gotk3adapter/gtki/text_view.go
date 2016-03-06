package gtki

type TextView interface {
	Container

	GetBuffer() (TextBuffer, error)
	SetBuffer(TextBuffer)
	SetCursorVisible(bool)
	SetEditable(bool)
}

func AssertTextView(_ TextView) {}
