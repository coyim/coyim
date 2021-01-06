package gtk_mock

import "github.com/coyim/gotk3adapter/gtki"

type MockEntry struct {
	MockWidget
}

func (*MockEntry) GetText() (string, error) {
	return "", nil
}

func (*MockEntry) SetHasFrame(v1 bool) {
}

func (*MockEntry) GetVisibility() bool {
	return false
}

func (*MockEntry) SetVisibility(v1 bool) {
}

func (*MockEntry) SetText(v1 string) {
}

func (*MockEntry) SetEditable(v1 bool) {
}

func (*MockEntry) SetWidthChars(v1 int) {
}

func (*MockEntry) GetAlignment() float32 {
	return 0.0
}

func (*MockEntry) SetAlignment(v1 float32) {
}

func (*MockEntry) SetPosition(p int) {
}

func (*MockEntry) GetPosition() int {
	return 0
}

func (*MockEntry) SetCompletion(v1 gtki.EntryCompletion) {
}

func (*MockEntry) SetPlaceholderText(v1 string) {
}
