package gtk_mock

import "github.com/coyim/gotk3adapter/gtki"

type MockMenuItem struct {
	MockBin
}

func (*MockMenuItem) GetLabel() string {
	return ""
}

func (*MockMenuItem) SetLabel(v1 string) {
}

func (*MockMenuItem) SetSubmenu(v1 gtki.Widget) {
}
