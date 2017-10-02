package gtk_mock

import "github.com/coyim/gotk3adapter/glib_mock"

type MockCssProvider struct {
	glib_mock.MockObject
}

func (*MockCssProvider) LoadFromData(v1 string) error {
	return nil
}
