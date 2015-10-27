package gui

import (
	"fmt"
	"strings"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
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

func onAccountDialogClicked(account *account, saveFunction func() error, reg *widgetRegistry) func() {
	return func() {
		account.session.CurrentAccount.Account = reg.getText("account")
		account.session.CurrentAccount.Password = reg.getText("password")

		parts := strings.SplitN(account.session.CurrentAccount.Account, "@", 2)
		if len(parts) != 2 {
			fmt.Println("invalid username (want user@domain): " + account.session.CurrentAccount.Account)
			return
		}

		go func() {
			if err := saveFunction(); err != nil {
				//TODO: handle errors
				fmt.Println(err.Error())
			}
		}()

		reg.dialogDestroy("dialog")
	}
}

func accountDialog(account *account, saveFunction func() error) {
	reg := createWidgetRegistry()
	d := dialog{
		title:    i18n.Local("Account Details"),
		position: gtk.WIN_POS_CENTER,
		id:       "dialog",
		content: []createable{
			label{i18n.Local("Account")},
			entry{
				text:       account.session.CurrentAccount.Account,
				editable:   true,
				visibility: true,
				id:         "account",
			},

			label{i18n.Local("Password\nAlert!! Your password is going to be stored as plaintext")},
			entry{
				text:       account.session.CurrentAccount.Password,
				editable:   true,
				visibility: false,
				id:         "password",
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

func buildAccountSubmenu(u *gtkUI, account *account) *gtk.MenuItem {
	menuitem, _ := gtk.MenuItemNewWithMnemonic(account.session.CurrentAccount.Account)

	accountSubMenu, _ := gtk.MenuNew()
	menuitem.SetSubmenu(accountSubMenu)

	connectItem, _ := gtk.MenuItemNewWithMnemonic(i18n.Local("_Connect"))
	accountSubMenu.Append(connectItem)

	disconnectItem, _ := gtk.MenuItemNewWithMnemonic(i18n.Local("_Disconnect"))
	accountSubMenu.Append(disconnectItem)

	toggleConnectAndDisconnectMenuItems(account.session, connectItem, disconnectItem)

	connectItem.Connect("activate", func() {
		connectItem.SetSensitive(false)
		u.connect(account)
	})

	disconnectItem.Connect("activate", func() {
		u.disconnect(account)
	})

	connToggle := func() {
		toggleConnectAndDisconnectMenuItems(account.session, connectItem, disconnectItem)
	}

	u.window.Connect(account.connectedSignal.String(), connToggle)
	u.window.Connect(account.disconnectedSignal.String(), connToggle)

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
