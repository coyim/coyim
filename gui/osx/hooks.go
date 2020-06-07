// +build darwin

package osx

import "github.com/coyim/coyim/gui"
import "github.com/coyim/gotk3osx"

// Create will return os hooks for OS X
func Create() gui.OSHooks {
	return &osxHooks{}
}

type osxHooks struct {
	app *gotk3osx.GtkosxApplication
}

func (h *osxHooks) AfterInit() {
	h.app, _ = gotk3osx.GetGtkosxApplication()
}
