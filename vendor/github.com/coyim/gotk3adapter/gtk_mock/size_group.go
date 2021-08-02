package gtk_mock

import (
	"github.com/coyim/gotk3adapter/glib_mock"
	"github.com/coyim/gotk3adapter/gtki"
)

type MockSizeGroup struct {
	glib_mock.MockObject
}

func (*MockSizeGroup) SetMode(gtki.SizeGroupMode) {
}
