package gtk_mock

import "github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtki"

type MockBin struct {
	MockContainer
}

func (*MockBin) GetChild() gtki.Widget {
	return nil
}
