// +build darwin

package osx

import "github.com/coyim/coyim/gui"
import "github.com/coyim/gotk3osx"
import "github.com/coyim/gotk3adapter/gdka"

// Create will return os hooks for OS X
func Create() gui.OSHooks {
	return &osxHooks{}
}

type osxHooks struct {
	app *gotk3osx.GtkosxApplication
}

func (h *osxHooks) AfterInit() {
	h.app, _ = gotk3osx.GetGtkosxApplication()
	h.app.Ready()

	p := gui.CoyimIcon.GetPixbuf()
	h.app.SetDockIconPixbuf(gdka.UnwrapPixbuf(p))
}
