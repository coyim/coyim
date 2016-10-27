package gtk_mock

import (
	"github.com/twstrike/gotk3adapter/gdki"
	"github.com/twstrike/gotk3adapter/gtki"
)

type MockWindow struct {
	MockBin
}

func (*MockWindow) AddAccelGroup(v2 gtki.AccelGroup) {
}

func (*MockWindow) GetTitle() string {
	return ""
}

func (*MockWindow) HasToplevelFocus() bool {
	return false
}

func (*MockWindow) Fullscreen() {
}

func (*MockWindow) Unfullscreen() {
}

func (*MockWindow) IsActive() bool {
	return false
}

func (*MockWindow) Resize(v1, v2 int) {
}

func (*MockWindow) SetApplication(v2 gtki.Application) {
}

func (*MockWindow) SetIcon(v2 gdki.Pixbuf) {
}

func (*MockWindow) SetTitle(v1 string) {
}

func (*MockWindow) SetTitlebar(v2 gtki.Widget) {
}

func (*MockWindow) SetTransientFor(v2 gtki.Window) {
}

func (*MockWindow) SetUrgencyHint(v2 bool) {
}

func (*MockWindow) Present() {
}
