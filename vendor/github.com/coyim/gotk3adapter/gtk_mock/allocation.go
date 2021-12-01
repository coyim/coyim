package gtk_mock

import (
	"github.com/coyim/gotk3adapter/gdk_mock"
)

type MockAllocation struct {
	gdk_mock.MockRectangle
}

func (*MockAllocation) GetY() int {
	return 0
}
