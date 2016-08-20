package gtk_mock

type MockEntry struct {
	MockWidget
}

func (*MockEntry) GetText() (string, error) {
	return "", nil
}

func (*MockEntry) SetHasFrame(v1 bool) {
}

func (*MockEntry) SetVisibility(v1 bool) {
}

func (*MockEntry) SetText(v1 string) {
}

func (*MockEntry) SetEditable(v1 bool) {
}
