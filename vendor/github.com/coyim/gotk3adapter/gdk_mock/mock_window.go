package gdk_mock

import "github.com/coyim/gotk3adapter/glib_mock"

type MockWindow struct {
	glib_mock.MockObject
}

func (*MockWindow) GetDesktop() uint32 {
	return 0
}

func (*MockWindow) MoveToDesktop(uint32) {
}
