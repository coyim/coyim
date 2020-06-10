// +build darwin

package gui

import "github.com/coyim/gotk3adapter/gtka"
import "github.com/coyim/gotk3adapter/gtki"

func (u *gtkUI) initializeMenus() {
	mb := u.mainBuilder.getObj("menubar").(gtki.MenuBar)
	mb.ShowAll()

	app := u.hooks.(*osxHooks).app
	app.SetMenuBar(gtka.UnwrapMenuShell(mb))
	app.SetHelpMenu(nil)
	app.SetWindowMenu(nil)

	aboutMenuItem := u.mainBuilder.getObj("aboutMenu").(gtki.MenuItem)
	prefsMenuItem := u.mainBuilder.getObj("preferencesMenuItem").(gtki.MenuItem)
	sepMenuItem, _ := g.gtk.SeparatorMenuItemNew()

	app.InsertAppMenuItem(gtka.UnwrapWidget(aboutMenuItem), 0)
	app.InsertAppMenuItem(gtka.UnwrapWidget(sepMenuItem), 1)
	app.InsertAppMenuItem(gtka.UnwrapWidget(prefsMenuItem), 2)

	mb.SetNoShowAll(true)
	mb.Hide()
}
