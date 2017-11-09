package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type mucMockupView struct {
	gtki.Window
}

func (u *gtkUI) openMUCMockup() {
	builder := newBuilder("MUCMockup")

	mockup := &mucMockupView{
		Window: builder.get("muc-window").(gtki.Window),
	}

	mockup.SetTransientFor(u.window)

	mockup.Show()
}
