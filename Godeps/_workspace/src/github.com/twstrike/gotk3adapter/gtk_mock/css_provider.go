package gtk_mock

import "github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glib_mock"

type MockCssProvider struct {
	glib_mock.MockObject
}

func (*MockCssProvider) LoadFromData(v1 string) error {
	return nil
}
