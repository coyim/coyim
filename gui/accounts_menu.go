package gui

import (
	"fmt"
	"strconv"

	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/session"
	"github.com/twstrike/go-gtk/glib"
	"github.com/twstrike/go-gtk/gtk"
)

var (
	AccountChangedSignal = glib.NewSignal("coyim-account-changed")
)

//TODO: Why does it receive Account?
func accountDialog(account Account, saveFunction func() error) {
	dialog := gtk.NewDialog()
	dialog.SetTitle("Account Details")
	dialog.SetPosition(gtk.WIN_POS_CENTER)
	vbox := dialog.GetVBox()

	accountLabel := gtk.NewLabel("Account")
	vbox.Add(accountLabel)

	accountInput := gtk.NewEntry()
	accountInput.SetText(account.Account)
	accountInput.SetEditable(true)
	vbox.Add(accountInput)

	vbox.Add(gtk.NewLabel("Password"))
	passwordInput := gtk.NewEntry()
	passwordInput.SetText(account.Password)
	passwordInput.SetEditable(true)
	passwordInput.SetVisibility(false)
	vbox.Add(passwordInput)

	vbox.Add(gtk.NewLabel("Server"))
	serverInput := gtk.NewEntry()
	serverInput.SetText(account.Server)
	serverInput.SetEditable(true)
	vbox.Add(serverInput)

	vbox.Add(gtk.NewLabel("Port"))
	portInput := gtk.NewEntry()
	portInput.SetText(strconv.Itoa(account.Port))
	portInput.SetEditable(true)
	vbox.Add(portInput)

	vbox.Add(gtk.NewLabel("Tor Proxy"))
	proxyInput := gtk.NewEntry()
	if len(account.Proxies) > 0 {
		proxyInput.SetText(account.Proxies[0])
	}
	proxyInput.SetEditable(true)
	vbox.Add(proxyInput)

	alwaysEncrypt := gtk.NewCheckButtonWithLabel("Always Encrypt")
	alwaysEncrypt.SetActive(account.AlwaysEncrypt)
	vbox.Add(alwaysEncrypt)

	button := gtk.NewButtonWithLabel("Save")
	button.Connect("clicked", func() {
		account.Account = accountInput.GetText()
		account.Password = passwordInput.GetText()
		account.Server = serverInput.GetText()

		v, err := strconv.Atoi(portInput.GetText())
		if err == nil {
			account.Port = v
		}

		if len(account.Proxies) == 0 {
			account.Proxies = append(account.Proxies, "")
		}
		account.Proxies[0] = proxyInput.GetText()

		account.AlwaysEncrypt = alwaysEncrypt.GetActive()

		if err := saveFunction(); err != nil {
			//TODO: handle errors
			fmt.Println(err.Error())
		}

		dialog.Destroy()
	})
	vbox.Add(button)

	dialog.ShowAll()
}

func buildAccountSubmenu(u *gtkUI, account Account) *gtk.MenuItem {
	menuitem := gtk.NewMenuItemWithMnemonic(account.Account)

	accountSubMenu := gtk.NewMenu()
	menuitem.SetSubmenu(accountSubMenu)

	connectItem := gtk.NewMenuItemWithMnemonic("_Connect")
	accountSubMenu.Append(connectItem)

	disconnectItem := gtk.NewMenuItemWithMnemonic("_Disconnect")
	disconnectItem.SetSensitive(false)
	accountSubMenu.Append(disconnectItem)

	connectItem.Connect("activate", func() {
		connectItem.SetSensitive(false)
		u.connect(account)
	})

	disconnectItem.Connect("activate", func() {
		u.disconnect(account)
	})

	connToggle := func() {
		s := account.Session
		connected := s.ConnStatus == session.CONNECTED
		connectItem.SetSensitive(!connected)
		disconnectItem.SetSensitive(connected)
	}

	u.window.Connect(account.Connected.Name(), connToggle)
	u.window.Connect(account.Disconnected.Name(), connToggle)

	editItem := gtk.NewMenuItemWithMnemonic("_Edit...")
	editItem.Connect("activate", func() {
		accountDialog(account, func() error {
			defer u.window.Emit(AccountChangedSignal.Name())
			return u.SaveConfig()
		})
	})
	accountSubMenu.Append(editItem)

	return menuitem
}

func populateAccountsSubMenu(menu *gtk.MenuItem, accounts []config.Config) {
}

func (u *gtkUI) buildAccountsMenu() {
	submenu := gtk.NewMenu()

	for _, account := range u.accounts {
		submenu.Append(buildAccountSubmenu(u, account))
	}

	if len(u.accounts) > 0 {
		submenu.Append(gtk.NewSeparatorMenuItem())
	}

	addAccMenu := gtk.NewMenuItemWithMnemonic("_Add...")
	addAccMenu.Connect("activate", u.showAddAccountWindow)

	submenu.Append(addAccMenu)
	submenu.ShowAll()

	u.accountsMenu.SetSubmenu(submenu)
}
