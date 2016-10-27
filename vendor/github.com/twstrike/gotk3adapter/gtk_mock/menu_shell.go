package gtk_mock

import "github.com/twstrike/gotk3adapter/gtki"

type MockMenuShell struct {
	MockContainer
}

func (*MockMenuShell) Append(v1 gtki.MenuItem) {
}
