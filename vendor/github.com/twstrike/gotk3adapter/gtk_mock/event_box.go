package gtk_mock

type MockEventBox struct {
	MockBin
}

func (*MockEventBox) SetAboveChild(v1 bool) {
}

func (*MockEventBox) GetAboveChild() bool {
	return false
}

func (*MockEventBox) SetVisibleWindow(v1 bool) {
}

func (*MockEventBox) GetVisibleWindow() bool {
	return false
}
