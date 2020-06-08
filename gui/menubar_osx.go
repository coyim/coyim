// +build darwin

package gui

import "github.com/coyim/gotk3adapter/gtka"
import "github.com/coyim/gotk3adapter/gtki"

func (u *gtkUI) removeMenuItemsBestPlacedElsewhere(mb gtki.MenuBar) {
	// TODO: we should remove the "settings" and "about" parts here
	// They will show up in other places anyway
}

func (u *gtkUI) initializeMenus() {
	mb := u.mainBuilder.getObj("menubar").(gtki.MenuBar)
	u.removeMenuItemsBestPlacedElsewhere(mb)

	mb.ShowAll()

	app := u.hooks.(*osxHooks).app
	app.SetMenuBar(gtka.UnwrapMenuShell(mb))
	app.SetHelpMenu(nil)
	app.SetWindowMenu(nil)

	mb.SetNoShowAll(true)
	mb.Hide()
}
