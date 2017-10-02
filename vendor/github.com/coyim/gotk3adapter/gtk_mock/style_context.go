package gtk_mock

import (
	"github.com/coyim/gotk3adapter/glib_mock"
	"github.com/coyim/gotk3adapter/gtki"
)

type MockStyleContext struct {
	glib_mock.MockObject
}

func (v *MockStyleContext) AddClass(v1 string) {
}

func (v *MockStyleContext) AddProvider(v1 gtki.StyleProvider, v2 uint) {
}

func (v *MockStyleContext) GetProperty2(v1 string, v2 gtki.StateFlags) (interface{}, error) {
	return nil, nil
}
