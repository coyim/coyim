package gtk_mock

import "github.com/coyim/gotk3adapter/gtki"

type MockTextView struct {
	MockContainer
}

func (*MockTextView) SetEditable(v1 bool) {
}

func (*MockTextView) SetCursorVisible(v1 bool) {
}

func (*MockTextView) SetBuffer(v1 gtki.TextBuffer) {
}

func (*MockTextView) GetBuffer() (gtki.TextBuffer, error) {
	return nil, nil
}

func (*MockTextView) ForwardDisplayLine(gtki.TextIter) bool {
	return false
}

func (*MockTextView) BackwardDisplayLine(gtki.TextIter) bool {
	return false
}

func (*MockTextView) ForwardDisplayLineEnd(gtki.TextIter) bool {
	return false
}

func (*MockTextView) BackwardDisplayLineStart(gtki.TextIter) bool {
	return false
}

func (*MockTextView) StartsDisplayLine(gtki.TextIter) bool {
	return false
}

func (*MockTextView) MoveVisually(gtki.TextIter, int) bool {
	return false
}
