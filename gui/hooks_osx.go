// +build darwin

package gui

import "github.com/coyim/gotk3osx"
import "github.com/coyim/gotk3adapter/gdka"

// CreateOSX will return os hooks for OS X
func CreateOSX() OSHooks {
	return &osxHooks{}
}

type osxHooks struct {
	app *gotk3osx.GtkosxApplication
	ui  *gtkUI
}

// BeforeMainWindow implements the OSHooks interface
func (h *osxHooks) BeforeMainWindow(ui *gtkUI) {
	h.ui = ui
}

// AfterInit implements the OSHooks interface
func (h *osxHooks) AfterInit() {
	h.app, _ = gotk3osx.GetGtkosxApplication()
	h.app.Ready()

	p := coyimIcon.GetPixbuf()
	h.app.SetDockIconPixbuf(gdka.UnwrapPixbuf(p))
}
