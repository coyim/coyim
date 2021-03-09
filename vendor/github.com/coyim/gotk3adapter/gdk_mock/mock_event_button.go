package gdk_mock

type MockEventButton struct {
	MockEvent
}

func (*MockEventButton) Button() uint {
	return 0
}

func (*MockEventButton) Time() uint32 {
	return 0
}

func (*MockEventButton) X() float64 {
	return 0
}

func (*MockEventButton) Y() float64 {
	return 0
}
