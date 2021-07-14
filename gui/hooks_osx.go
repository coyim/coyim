// +build darwin

package gui

import "github.com/coyim/gotk3osx/access"

// CreateOSX will return os hooks for OS X
func CreateOSX() OSHooks {
	return &osxHooks{}
}

type osxHooks struct {
	app access.Application
	ui  *gtkUI
}

// BeforeMainWindow implements the OSHooks interface
func (h *osxHooks) BeforeMainWindow(ui *gtkUI) {
	h.ui = ui
}

// AfterInit implements the OSHooks interface
func (h *osxHooks) AfterInit() {
	h.app, _ = g.extra.(access.GTKOSX).GetApplication()
	h.app.Ready()

	p := coyimIcon.GetPixbuf()
	h.app.SetDockIconPixbuf(p)
}
