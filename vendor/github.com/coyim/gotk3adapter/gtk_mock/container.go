package gtk_mock

import "github.com/coyim/gotk3adapter/gtki"

type MockContainer struct {
	MockWidget
}

func (*MockContainer) Add(v2 gtki.Widget) {
}

func (*MockContainer) Remove(v2 gtki.Widget) {
}

func (*MockContainer) SetBorderWidth(v1 uint) {
}
