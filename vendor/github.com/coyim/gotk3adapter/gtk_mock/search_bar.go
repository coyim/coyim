package gtk_mock

import (
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/gtki"
)

type MockSearchBar struct {
	MockBin
}

func (*MockSearchBar) ConnectEntry(v1 gtki.Entry) {
}

func (*MockSearchBar) GetSearchMode() bool {
	return false
}

func (*MockSearchBar) SetSearchMode(v1 bool) {
}

func (*MockSearchBar) GetShowCloseButton() bool {
	return false
}

func (*MockSearchBar) SetShowCloseButton(v1 bool) {
}

func (*MockSearchBar) HandleEvent(v1 gdki.Event) {
}
