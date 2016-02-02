package gui

import (
	"sync"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/coyim/session"
)

var (
	// TODO: shouldn't this be specific to the account ID in question?
	accountChangedSignal, _ = glib.SignalNew("coyim-account-changed")
)

func toggleConnectAndDisconnectMenuItems(s *session.Session, connect, disconnect *gtk.MenuItem) {
	doInUIThread(func() {
		connect.SetSensitive(s.IsDisconnected())
		disconnect.SetSensitive(s.IsConnected())
	})
}

var accountsLock sync.Mutex

func (u *gtkUI) buildStaticAccountsMenu(submenu *gtk.Menu) {
	connectAutomaticallyItem, _ := gtk.CheckMenuItemNewWithMnemonic(i18n.Local("Connect On _Startup"))
	u.config.WhenLoaded(func(a *config.ApplicationConfig) {
		connectAutomaticallyItem.SetActive(a.ConnectAutomatically)
	})

	connectAutomaticallyItem.Connect("activate", func() {
		u.toggleConnectAllAutomatically()
	})
	submenu.Append(connectAutomaticallyItem)

	connectAllMenu, _ := gtk.MenuItemNewWithMnemonic(i18n.Local("_Connect All"))
	connectAllMenu.Connect("activate", func() { u.connectAllAutomatics(true) })
	submenu.Append(connectAllMenu)

	sep2, _ := gtk.SeparatorMenuItemNew()
	submenu.Append(sep2)

	addAccMenu, _ := gtk.MenuItemNewWithMnemonic(i18n.Local("_Add..."))
	addAccMenu.Connect("activate", u.showAddAccountWindow)
	submenu.Append(addAccMenu)

	importMenu, _ := gtk.MenuItemNewWithMnemonic(i18n.Local("_Import..."))
	importMenu.Connect("activate", u.runImporter)
	submenu.Append(importMenu)

	registerAccMenu, _ := gtk.MenuItemNewWithMnemonic(i18n.Local("_Register..."))
	registerAccMenu.Connect("activate", u.showServerSelectionWindow)
	submenu.Append(registerAccMenu)
}

func (u *gtkUI) buildAccountsMenu() {
	accountsLock.Lock()
	defer accountsLock.Unlock()

	submenu, _ := gtk.MenuNew()

	for _, account := range u.accounts {
		account.appendMenuTo(submenu)
	}

	if len(u.accounts) > 0 {
		sep, _ := gtk.SeparatorMenuItemNew()
		submenu.Append(sep)
	}

	u.buildStaticAccountsMenu(submenu)

	submenu.ShowAll()

	u.accountsMenu.SetSubmenu(submenu)
}
