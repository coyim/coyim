package gtk_mock

import "github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtki"

type MockContainer struct {
	MockWidget
}

func (*MockContainer) Add(v2 gtki.Widget) {
}

func (*MockContainer) Remove(v2 gtki.Widget) {
}
