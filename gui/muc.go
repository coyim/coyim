package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type mucMockupView struct {
	gtki.Dialog
}

func (u *gtkUI) openMUCMockup() {
	builder := newBuilder("MUCMockup")

	mockup := &mucMockupView{
		Dialog: builder.get("muc-dialog").(gtki.Dialog),
	}

	mockup.SetTransientFor(u.window)

	mockup.Show()
}
