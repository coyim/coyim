package gtk_mock

import "github.com/coyim/gotk3adapter/glib_mock"

type MockIconTheme struct {
	glib_mock.MockObject
}

func (*MockIconTheme) AddResourcePath(path string) {
}

func (*MockIconTheme) AppendSearchPath(path string) {
}

func (*MockIconTheme) GetExampleIconName() string {
	return ""
}

func (*MockIconTheme) HasIcon(name string) bool {
	return false
}

func (*MockIconTheme) PrependSearchPath(path string) {
}
