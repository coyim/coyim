package gtk_mock

import (
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glib_mock"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtki"
)

type MockSettings interface {
	glib_mock.MockObject
}

func (*MockSettings) GetProperty(string) (interface{}, error) {
	return nil, nil
}
