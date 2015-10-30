package gui

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/coyim/session"
)

type account struct {
	connectedSignal    *glib.Signal
	disconnectedSignal *glib.Signal
	menu               *gtk.MenuItem

	session *session.Session

	onConnect    chan<- *account
	onDisconnect chan<- *account
	onEdit       chan<- *account
}

func newAccount(conf *config.Accounts, currentConf *config.Account) (acc *account, err error) {
	acc = &account{}

	id := currentConf.ID()
	acc.connectedSignal, err = glib.SignalNew(signalName(id, "connected"))
	if err != nil {
		return
	}

	acc.disconnectedSignal, err = glib.SignalNew(signalName(id, "disconnected"))
	if err != nil {
		return
	}

	acc.session = session.NewSession(conf, currentConf)

	return
}

func signalName(id, signal string) string {
	return "coyim-account-" + signal + "-" + id
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
	if account.menu == nil {
		account.buildAccountSubmenu()
	}

	account.menu.SetLabel(account.session.CurrentAccount.Account)

	account.menu.Unparent()
	submenu.Append(account.menu)
	submenu.ShowAll()
}

func (account *account) buildAccountSubmenu() {
	menuitem, _ := gtk.MenuItemNew()

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
