package gtk_mock

import (
	"github.com/twstrike/gotk3adapter/gdki"
	"github.com/twstrike/gotk3adapter/gtki"
)

type MockMenu struct {
	MockMenuShell
}

func (*MockMenu) PopupAtMouseCursor(v1 gtki.Menu, v2 gtki.MenuItem, v3 int, v4 uint32) {
}

func (*MockMenu) PopupAtPointer(_ gdki.Event) {
}
