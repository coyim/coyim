package gtki

type TextView interface {
	Container

	BackwardDisplayLine(TextIter) bool
	BackwardDisplayLineStart(TextIter) bool
	ForwardDisplayLine(TextIter) bool
	ForwardDisplayLineEnd(TextIter) bool
	GetBuffer() (TextBuffer, error)
	MoveVisually(TextIter, int) bool
	SetBuffer(TextBuffer)
	SetCursorVisible(bool)
	SetEditable(bool)
	StartsDisplayLine(TextIter) bool
}

func AssertTextView(_ TextView) {}
