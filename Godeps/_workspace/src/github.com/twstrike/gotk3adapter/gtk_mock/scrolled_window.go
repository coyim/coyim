package gtk_mock

import "github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtki"

type MockScrolledWindow struct {
	MockBin
}

func (*MockScrolledWindow) GetVAdjustment() gtki.Adjustment {
	return nil
}
