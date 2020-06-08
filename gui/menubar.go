package gui

import "github.com/coyim/coyim/i18n"
import "github.com/coyim/gotk3adapter/glibi"

/*
  OK, we have a few different parameters:

  - PrefersAppMenu
  - gtk-shell-shows-app-menu
  - gtk-shell-shows-menubar

  We have three different types of menus:
  - Our internal menu
  - The menubar
  - The app-menu

  Our internal menu is Gtk
  The menubar is GMenuModel
  The app-menu is GMenuModel

  However, on OS X, we are using the gtk-mac-integration library, and if we do that we can set these things properly.
  So, on OS X we will ignore the parameters and we will:
    - Set menubar by taking the object "menubar" which is our internal menu, and remove it from the outside
      - We will remove some options from it, to make them show up on the app-menu
    - We will NOT use the app-menu directly, only indirectly

  We will ignore PrefersAppMenu, because it doesn't really suit our application.

  On all other platforms we will:
    - if gtk-shell-shows-app-menu:
      - we will create a few simple options - quit, preferences and about, primarily
    if gtk-shell-shows-menubar:
      - we will remove our internal menu and create the outside menu
    if !gtk-shell-shows-menubar
      - we will use our own internal menu, and do nothing

  TODO: We need to create a GMenuModel version of our internal menu
  TODO: So we need actions for all things you can do
  TODO: and then we need to bind things properly
*/

func (u *gtkUI) createSimpleAppMenu() glibi.MenuModel {
	top := g.glib.MenuNew()

	aboutSection := g.glib.MenuNew()
	aboutMenuItem := g.glib.MenuItemNew(i18n.Local("About CoyIM"), "app.about")
	aboutSection.AppendItem(aboutMenuItem)

	prefsSection := g.glib.MenuNew()
	prefsMenuItem := g.glib.MenuItemNew(i18n.Local("Preferences..."), "app.preferences")
	prefsSection.AppendItem(prefsMenuItem)

	quitSection := g.glib.MenuNew()
	quitMenuItem := g.glib.MenuItemNew(i18n.Local("Quit CoyIM"), "app.quit")
	quitSection.AppendItem(quitMenuItem)

	top.AppendSection("about", aboutSection)
	top.AppendSection("prefs", prefsSection)
	top.AppendSection("quit", quitSection)

	return top
}
