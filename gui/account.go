package gui

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/coyim/session"
)

type account struct {
	menu *gtk.MenuItem

	session *session.Session

	onConnect    chan<- *account
	onDisconnect chan<- *account
	onEdit       chan<- *account
}

func newAccount(conf *config.Accounts, currentConf *config.Account) (acc *account, err error) {
	acc = &account{}
	acc.session = session.NewSession(conf, currentConf)

	return
}

func (account *account) connected() bool {
	return account.session.ConnStatus == session.CONNECTED
}

func (u *gtkUI) showAddAccountWindow() error {
	c, err := config.NewAccount()
	if err != nil {
		return err
	}
	accountDialog(c, func() {
		u.config.Add(c)
		u.SaveConfig()
	})
	return nil
}

func (account *account) appendMenuTo(submenu *gtk.Menu) {
	if account.menu != nil {
		//I dont know how to remove it from its current parent without breaking the menu
		//unparent does not seem to work as expected
		account.menu.Destroy()
	}

	account.buildAccountSubmenu()
	account.menu.Show()
	submenu.Append(account.menu)
}

func (account *account) buildAccountSubmenu() {
	menuitem, _ := gtk.MenuItemNew()

	menuitem.SetLabel(account.session.CurrentAccount.Account)

	accountSubMenu, _ := gtk.MenuNew()
	menuitem.SetSubmenu(accountSubMenu)

	connectItem, _ := gtk.MenuItemNewWithMnemonic(i18n.Local("_Connect"))
	accountSubMenu.Append(connectItem)

	disconnectItem, _ := gtk.MenuItemNewWithMnemonic(i18n.Local("_Disconnect"))
	accountSubMenu.Append(disconnectItem)

	toggleConnectAndDisconnectMenuItems(account.session, connectItem, disconnectItem)

	connectItem.Connect("activate", func() {
		account.onConnect <- account
	})

	disconnectItem.Connect("activate", func() {
		account.onDisconnect <- account
	})

	editItem, _ := gtk.MenuItemNewWithMnemonic(i18n.Local("_Edit..."))
	accountSubMenu.Append(editItem)

	editItem.Connect("activate", func() {
		account.onEdit <- account
	})

	//TODO: add "Remove" menu item

	c := make(chan interface{})
	account.session.Subscribe(c)

	go func() {
		for ev := range c {
			switch t := ev.(type) {
			case session.Event:
				switch t.Type {
				case session.Connected, session.Disconnected:
					toggleConnectAndDisconnectMenuItems(t.Session, connectItem, disconnectItem)
				}
			}
		}
	}()

	account.menu = menuitem
}
