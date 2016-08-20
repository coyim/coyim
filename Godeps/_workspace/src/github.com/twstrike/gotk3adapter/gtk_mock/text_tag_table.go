package gtk_mock

import (
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glib_mock"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtki"
)

type MockTextTagTable struct {
	glib_mock.MockObject
}

func (*MockTextTagTable) Add(v1 gtki.TextTag) {
}
