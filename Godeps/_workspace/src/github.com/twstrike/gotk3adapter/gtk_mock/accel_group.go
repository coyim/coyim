package gtk_mock

import (
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gdki"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glib_mock"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtki"
)

type MockAccelGroup struct {
	glib_mock.MockObject
}

func (*MockAccelGroup) Connect2(v2 uint, v3 gdki.ModifierType, v4 gtki.AccelFlags, v5 interface{}) {
}
