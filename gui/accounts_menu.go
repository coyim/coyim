package gui

import (
	"sync"

	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glibi"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtki"
	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/coyim/session/access"
)

var (
	// TODO: shouldn't this be specific to the account ID in question?
	accountChangedSignal glibi.Signal
)

func toggleConnectAndDisconnectMenuItems(s access.Session, connect, disconnect gtki.MenuItem) {
	doInUIThread(func() {
		connect.SetSensitive(s.IsDisconnected())
		disconnect.SetSensitive(!s.IsDisconnected())
	})
}

var accountsLock sync.Mutex

func (u *gtkUI) buildStaticAccountsMenu(submenu gtki.Menu) {
	connectAutomaticallyItem, _ := g.gtk.CheckMenuItemNewWithMnemonic(i18n.Local("Connect On _Startup"))
	u.config.WhenLoaded(func(a *config.ApplicationConfig) {
		connectAutomaticallyItem.SetActive(a.ConnectAutomatically)
	})

	connectAutomaticallyItem.Connect("activate", func() {
		u.setConnectAllAutomatically(connectAutomaticallyItem.GetActive())
	})
	submenu.Append(connectAutomaticallyItem)

	connectAllMenu, _ := g.gtk.MenuItemNewWithMnemonic(i18n.Local("_Connect All"))
	connectAllMenu.Connect("activate", func() { u.connectAllAutomatics(true) })
	submenu.Append(connectAllMenu)

	sep2, _ := g.gtk.SeparatorMenuItemNew()
	submenu.Append(sep2)

	addAccMenu, _ := g.gtk.MenuItemNewWithMnemonic(i18n.Local("_Add..."))
	addAccMenu.Connect("activate", u.showAddAccountWindow)
	submenu.Append(addAccMenu)

	importMenu, _ := g.gtk.MenuItemNewWithMnemonic(i18n.Local("_Import..."))
	importMenu.Connect("activate", u.runImporter)
	submenu.Append(importMenu)

	registerAccMenu, _ := g.gtk.MenuItemNewWithMnemonic(i18n.Local("_Register..."))
	registerAccMenu.Connect("activate", u.showServerSelectionWindow)
	submenu.Append(registerAccMenu)
}

func (u *gtkUI) buildAccountsMenu() {
	accountsLock.Lock()
	defer accountsLock.Unlock()

	submenu, _ := g.gtk.MenuNew()

	for _, account := range u.accounts {
		account.appendMenuTo(submenu)
	}

	if len(u.accounts) > 0 {
		sep, _ := g.gtk.SeparatorMenuItemNew()
		submenu.Append(sep)
	}

	u.buildStaticAccountsMenu(submenu)

	submenu.ShowAll()

	u.accountsMenu.SetSubmenu(submenu)
}
