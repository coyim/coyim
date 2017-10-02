package gtk_mock

import "github.com/coyim/gotk3adapter/gtki"

type MockBin struct {
	MockContainer
}

func (*MockBin) GetChild() gtki.Widget {
	return nil
}
