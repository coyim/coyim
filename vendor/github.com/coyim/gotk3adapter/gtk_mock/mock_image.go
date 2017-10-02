package gtk_mock

import (
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/gtki"
)

type MockImage struct {
	MockWidget
}

func (v *MockImage) SetFromIconName(v1 string, v2 gtki.IconSize) {
}

func (v *MockImage) Clear() {
}

func (v *MockImage) SetFromPixbuf(pb gdki.Pixbuf) {
}
