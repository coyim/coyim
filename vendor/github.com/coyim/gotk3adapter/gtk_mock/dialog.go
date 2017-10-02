package gtk_mock

import "github.com/coyim/gotk3adapter/gtki"

type MockDialog struct {
	MockWindow
}

func (*MockDialog) Run() int {
	return 0
}

func (*MockDialog) SetDefaultResponse(v1 gtki.ResponseType) {
}
