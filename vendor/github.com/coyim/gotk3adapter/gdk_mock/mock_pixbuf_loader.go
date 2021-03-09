package gdk_mock

import (
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/glib_mock"
)

type MockPixbufLoader struct {
	glib_mock.MockObject
}

func (*MockPixbufLoader) Close() error {
	return nil
}

func (*MockPixbufLoader) GetPixbuf() (gdki.Pixbuf, error) {
	return nil, nil
}

func (*MockPixbufLoader) SetSize(int, int) {
}

func (*MockPixbufLoader) Write(b []byte) (int, error) {
	return 0, nil
}
