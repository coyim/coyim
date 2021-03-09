package gdk_mock

import "github.com/coyim/gotk3adapter/gdki"

type Mock struct{}

func (*Mock) EventButtonFrom(ev gdki.Event) gdki.EventButton {
	return nil
}

func (*Mock) EventKeyFrom(ev gdki.Event) gdki.EventKey {
	return nil
}

func (*Mock) PixbufLoaderNew() (gdki.PixbufLoader, error) {
	return &MockPixbufLoader{}, nil
}

func (*Mock) ScreenGetDefault() (gdki.Screen, error) {
	return &MockScreen{}, nil
}

func (*Mock) WorkspaceControlSupported() bool {
	return false
}
