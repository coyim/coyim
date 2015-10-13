package gui

import (
	"fmt"
	"strconv"

	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/coyim/session"
	"github.com/twstrike/gotk3/glib"
	"github.com/twstrike/gotk3/gtk"
)

var (
	// TODO: shouldn't this be specific to the account ID in question?
	AccountChangedSignal, _ = glib.NewSignal("coyim-account-changed")
)

func firstProxy(account Account) string {
	if len(account.Proxies) > 0 {
		return account.Proxies[0]
	}
	return ""
}

func onAccountDialogClicked(account Account, saveFunction func() error, reg *widgetRegistry) func() {
	return func() {
		account.Account = reg.getText("account")
		account.Password = reg.getText("password")
		account.Server = reg.getText("server")

		if v, err := strconv.Atoi(reg.getText("port")); err == nil {
			account.Port = v
		}

		if len(account.Proxies) == 0 {
			account.Proxies = append(account.Proxies, "")
		}

		account.Proxies[0] = reg.getText("tor proxy")
		account.AlwaysEncrypt = reg.getActive("always encrypt")

		if err := saveFunction(); err != nil {
			//TODO: handle errors
			fmt.Println(err.Error())
		}

		reg.dialogDestroy("dialog")
	}
}

func accountDialog(account Account, saveFunction func() error) {
	reg := createWidgetRegistry()
	d := dialog{
		title:    i18n.Local("Account Details"),
		position: gtk.WIN_POS_CENTER,
		id:       "dialog",
		content: []createable{
			label{i18n.Local("Account")},
			entry{
				text:       account.Account,
				editable:   true,
				visibility: true,
				id:         "account",
			},
			label{i18n.Local("Password")},
			entry{
				text:       account.Password,
				editable:   true,
				visibility: false,
				id:         "password",
			},
			label{i18n.Local("Server")},
			entry{
				text:       account.Server,
				editable:   true,
				visibility: true,
				id:         "server",
			},
			label{i18n.Local("Port")},
			entry{
				text:       strconv.Itoa(account.Port),
				editable:   true,
				visibility: true,
				id:         "port",
			},

			//TODO: I'm considering to not show this, since it will autodetect Tor running
			//and autoconfigure the account to use separate circuits.
			//Should the user be able to disable this even if Tor is found?
			label{i18n.Local("Tor Proxy")},
			entry{
				text:       firstProxy(account),
				editable:   true,
				visibility: true,
				id:         "tor proxy",
			},

			//TODO: I'm also considering to hide this
			checkButton{
				text:   i18n.Local("Always Encrypt"),
				active: account.AlwaysEncrypt,
				id:     "always encrypt",
			},

			button{
				text:      i18n.Local("Save"),
				onClicked: onAccountDialogClicked(account, saveFunction, reg),
			},
		},
	}

	d.create(reg)
	reg.dialogShowAll("dialog")
}

func toggleConnectAndDisconnectMenuItems(s *session.Session, connect, disconnect *gtk.MenuItem) {
	connected := s.ConnStatus == session.CONNECTED
	connect.SetSensitive(!connected)
	disconnect.SetSensitive(connected)
}

func buildAccountSubmenu(u *gtkUI, account Account) *gtk.MenuItem {
	menuitem, _ := gtk.MenuItemNewWithMnemonic(account.Account)

	accountSubMenu, _ := gtk.MenuNew()
	menuitem.SetSubmenu(accountSubMenu)

	connectItem, _ := gtk.MenuItemNewWithMnemonic(i18n.Local("_Connect"))
	accountSubMenu.Append(connectItem)

	disconnectItem, _ := gtk.MenuItemNewWithMnemonic(i18n.Local("_Disconnect"))
	accountSubMenu.Append(disconnectItem)

	toggleConnectAndDisconnectMenuItems(account.Session, connectItem, disconnectItem)

	connectItem.Connect("activate", func() {
		connectItem.SetSensitive(false)
		u.connect(account)
	})

	disconnectItem.Connect("activate", func() {
		u.disconnect(account)
	})

	connToggle := func() {
		toggleConnectAndDisconnectMenuItems(account.Session, connectItem, disconnectItem)
	}

	u.window.Connect(account.ConnectedSignal.String(), connToggle)
	u.window.Connect(account.DisconnectedSignal.String(), connToggle)

	editItem, _ := gtk.MenuItemNewWithMnemonic(i18n.Local("_Edit..."))
	accountSubMenu.Append(editItem)

	editItem.Connect("activate", func() {
		accountDialog(account, u.SaveConfig)
	})

	//TODO: add "Remove" menu item

	return menuitem
}

func (u *gtkUI) buildAccountsMenu() {
	submenu, _ := gtk.MenuNew()

	for _, account := range u.accounts {
		submenu.Append(buildAccountSubmenu(u, account))
	}

	if len(u.accounts) > 0 {
		sep, _ := gtk.SeparatorMenuItemNew()
		submenu.Append(sep)
	}

	addAccMenu, _ := gtk.MenuItemNewWithMnemonic(i18n.Local("_Add..."))
	addAccMenu.Connect("activate", u.showAddAccountWindow)

	submenu.Append(addAccMenu)
	submenu.ShowAll()

	u.accountsMenu.SetSubmenu(submenu)
}
