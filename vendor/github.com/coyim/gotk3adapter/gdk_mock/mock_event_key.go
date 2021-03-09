package gdk_mock

type MockEventKey struct {
	MockEvent
}

func (*MockEventKey) KeyVal() uint {
	return 0
}

func (*MockEventKey) State() uint {
	return 0
}
