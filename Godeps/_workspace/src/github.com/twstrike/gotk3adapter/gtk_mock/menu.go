package gtk_mock

import "github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtki"

type MockMenu struct {
	MockMenuShell
}

func (*MockMenu) PopupAtMouseCursor(v1 gtki.Menu, v2 gtki.MenuItem, v3 int, v4 uint32) {
}
