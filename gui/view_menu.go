package gui

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/i18n"
)

type viewMenu struct {
	merge   *gtk.CheckMenuItem
	offline *gtk.CheckMenuItem
}

func (u *gtkUI) createViewMenu(bar *gtk.MenuBar) {
	u.viewMenu = new(viewMenu)

	viewMenu, _ := gtk.MenuItemNewWithMnemonic(i18n.Local("_View"))
	bar.Append(viewMenu)
	viewSubmenu, _ := gtk.MenuNew()
	viewMenu.SetSubmenu(viewSubmenu)

	checkItemMerge, _ := gtk.CheckMenuItemNewWithMnemonic(i18n.Local("_Merge Accounts"))
	u.viewMenu.merge = checkItemMerge
	viewSubmenu.Append(checkItemMerge)
	checkItemMerge.Connect("toggled", u.toggleMergeAccounts)

	checkItemShowOffline, _ := gtk.CheckMenuItemNewWithMnemonic(i18n.Local("Show _Offline Contacts"))
	u.viewMenu.offline = checkItemShowOffline
	viewSubmenu.Append(checkItemShowOffline)
	checkItemShowOffline.Connect("toggled", u.toggleShowOffline)

	u.displaySettings.defaultSettingsOn(&checkItemMerge.MenuItem.Bin.Container.Widget)
	u.displaySettings.defaultSettingsOn(&checkItemShowOffline.MenuItem.Bin.Container.Widget)
}

func (v *viewMenu) setFromConfig(c *config.Accounts) {
	glib.IdleAdd(func() bool {
		v.merge.SetActive(c.MergeAccounts)
		v.offline.SetActive(!c.ShowOnlyOnline)
		return false
	})
}

func (u *gtkUI) toggleMergeAccounts() {
	if u.config != nil {
		u.config.MergeAccounts = u.viewMenu.merge.GetActive()
		u.saveConfigOnly()
	}

	u.roster.redrawIfRosterVisible()
}

func (u *gtkUI) toggleShowOffline() {
	if u.config != nil {
		u.config.ShowOnlyOnline = !u.viewMenu.offline.GetActive()
		u.saveConfigOnly()
	}

	u.roster.redrawIfRosterVisible()
}
