package gtk_mock

import "github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtki"

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
