package gui

import (
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"../config"
	"../i18n"
	"../session"
)

var (
	// TODO: shouldn't this be specific to the account ID in question?
	accountChangedSignal, _ = glib.SignalNew("coyim-account-changed")
)

func firstProxy(account *account) string {
	if len(account.session.CurrentAccount.Proxies) > 0 {
		return account.session.CurrentAccount.Proxies[0]
	}

	return ""
}

func (u *gtkUI) accountDialog(account *config.Account, saveFunction func()) {
	dialogID := "AccountDetails"
	builder := builderForDefinition(dialogID)

	obj, _ := builder.GetObject(dialogID)
	dialog := obj.(*gtk.Dialog)

	obj, _ = builder.GetObject("account")
	accEntry := obj.(*gtk.Entry)
	accEntry.SetText(account.Account)

	obj, _ = builder.GetObject("password")
	passEntry := obj.(*gtk.Entry)

	obj, _ = builder.GetObject("server")
	serverEntry := obj.(*gtk.Entry)
	serverEntry.SetText(account.Server)

	obj, _ = builder.GetObject("port")
	portEntry := obj.(*gtk.Entry)
	if account.Port == 0 {
		account.Port = 5222
	}
	portEntry.SetText(strconv.Itoa(account.Port))

	builder.ConnectSignals(map[string]interface{}{
		"on_save_signal": func() {
			accTxt, _ := accEntry.GetText()
			passTxt, _ := passEntry.GetText()
			servTxt, _ := serverEntry.GetText()
			portTxt, _ := portEntry.GetText()

			account.Account = accTxt
			account.Server = servTxt

			if passTxt != "" {
				account.Password = passTxt
			}

			convertedPort, e := strconv.Atoi(portTxt)
			if len(strings.TrimSpace(portTxt)) == 0 || e != nil {
				convertedPort = 5222
			}

			account.Port = convertedPort

			parts := strings.SplitN(account.Account, "@", 2)
			if len(parts) != 2 {
				log.Println("invalid username (want user@domain): " + account.Account)
				return
			}

			go saveFunction()
			dialog.Destroy()
		},
		"on_cancel_signal": func() {
			u.buildAccountsMenu()
			dialog.Destroy()
		},
	})

	dialog.SetTransientFor(u.window)
	dialog.ShowAll()
}

func toggleConnectAndDisconnectMenuItems(s *session.Session, connect, disconnect *gtk.MenuItem) {
	glib.IdleAdd(func() {
		connect.SetSensitive(s.ConnStatus == session.DISCONNECTED)
		disconnect.SetSensitive(s.ConnStatus == session.CONNECTED)
	})
}

var accountsLock sync.Mutex

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
	addAccMenu.Connect("activate", func() { u.showAddAccountWindow() })
	submenu.Append(addAccMenu)

	importMenu, _ := gtk.MenuItemNewWithMnemonic(i18n.Local("_Import..."))
	importMenu.Connect("activate", func() { u.runImporter() })
	submenu.Append(importMenu)

	submenu.ShowAll()

	u.accountsMenu.SetSubmenu(submenu)
}
