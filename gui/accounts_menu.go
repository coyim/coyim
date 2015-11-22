package gui

import (
	"log"
	"strings"

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

func firstProxy(account *account) string {
	if len(account.session.CurrentAccount.Proxies) > 0 {
		return account.session.CurrentAccount.Proxies[0]
	}
	return ""
}

func accountDialog(account *config.Account, saveFunction func()) {
	vars := make(map[string]string)
	vars["$title"] = i18n.Local("Account Details")
	vars["$accountMessage"] = i18n.Local("Your account (for example: kim42@dukgo.com)")
	vars["$pswMessage"] = i18n.Local("Password")
	vars["$saveLabel"] = i18n.Local("Save")
	builder, buildError := loadBuilderWith("AccountDetails", vars)
	if buildError != nil {
		panic(buildError.Error())
	}
	obj, _ := builder.GetObject("AccountDetailsDialog")
	dialog := obj.(*gtk.Dialog)
	accObj, _ := builder.GetObject("account")
	accEntry := accObj.(*gtk.Entry)
	accEntry.SetText(account.Account)

	builder.ConnectSignals(map[string]interface{}{
		"on_save_signal": func() {
			passObj, _ := builder.GetObject("password")
			accTxt, _ := accEntry.GetText()
			passTxt, _ := passObj.(*gtk.Entry).GetText()
			account.Account = accTxt
			account.Password = passTxt

			parts := strings.SplitN(account.Account, "@", 2)
			if len(parts) != 2 {
				log.Println("invalid username (want user@domain): " + account.Account)
				return
			}

			go saveFunction()
			dialog.Destroy()
		},
	})
	dialog.ShowAll()
}

func toggleConnectAndDisconnectMenuItems(s *session.Session, connect, disconnect *gtk.MenuItem) {
	connected := s.ConnStatus == session.CONNECTED
	connect.SetSensitive(!connected)
	disconnect.SetSensitive(connected)
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

	connectAutomaticallyItem, _ := gtk.CheckMenuItemNewWithMnemonic(i18n.Local("Connect On _Startup"))
	u.config.WhenLoaded(func(a *config.Accounts) {
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
