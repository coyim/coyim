package gtk_mock

import "github.com/coyim/gotk3adapter/gtki"

type MockInfoBar struct {
	MockBox
}

func (*MockInfoBar) GetOrientation() gtki.Orientation {
	return gtki.HorizontalOrientation
}

func (*MockInfoBar) SetOrientation(o gtki.Orientation) {
}
