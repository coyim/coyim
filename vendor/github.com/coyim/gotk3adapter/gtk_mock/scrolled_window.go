package gtk_mock

import "github.com/coyim/gotk3adapter/gtki"

type MockScrolledWindow struct {
	MockBin
}

func (*MockScrolledWindow) GetVAdjustment() gtki.Adjustment {
	return nil
}
