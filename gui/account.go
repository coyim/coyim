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

	onConnect                         chan<- *account
	onDisconnect                      chan<- *account
	onEdit                            chan<- *account
	onRemove                          chan<- *account
	toggleConnectAutomaticallyRequest chan<- *account

	sessionObserver chan interface{}
}

type byAccountNameAlphabetic []*account

func (s byAccountNameAlphabetic) Len() int { return len(s) }
func (s byAccountNameAlphabetic) Less(i, j int) bool {
	return s[i].session.CurrentAccount.Account < s[j].session.CurrentAccount.Account
}
func (s byAccountNameAlphabetic) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

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
		u.addAndSaveAccountConfig(c)
	})

	return nil
}

func (account *account) destroyMenu() {
	close(account.sessionObserver)

	//I dont know how to remove it from its current parent without breaking the menu
	//unparent does not seem to work as expected
	account.menu.Destroy()
}

func (account *account) appendMenuTo(submenu *gtk.Menu) {
	if account.menu != nil {
		account.destroyMenu()
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

	editItem, _ := gtk.MenuItemNewWithMnemonic(i18n.Local("_Edit..."))
	accountSubMenu.Append(editItem)

	removeItem, _ := gtk.MenuItemNewWithMnemonic(i18n.Local("_Remove"))
	accountSubMenu.Append(removeItem)

	connectAutomaticallyItem, _ := gtk.CheckMenuItemNewWithMnemonic(i18n.Local("Connect _Automatically"))
	accountSubMenu.Append(connectAutomaticallyItem)
	connectAutomaticallyItem.SetActive(account.session.CurrentAccount.ConnectAutomatically)

	toggleConnectAndDisconnectMenuItems(account.session, connectItem, disconnectItem)

	connectItem.Connect("activate", func() {
		account.onConnect <- account
	})

	disconnectItem.Connect("activate", func() {
		account.onDisconnect <- account
	})

	connectAutomaticallyItem.Connect("activate", func() {
		account.toggleConnectAutomaticallyRequest <- account
	})

	editItem.Connect("activate", func() {
		account.onEdit <- account
	})

	removeItem.Connect("activate", func() {
		account.onRemove <- account
	})

	go account.watchAndToggleMenuItems(connectItem, disconnectItem)
	account.menu = menuitem
}

func (account *account) watchAndToggleMenuItems(connectItem, disconnectItem *gtk.MenuItem) {
	account.sessionObserver = make(chan interface{})
	account.session.Subscribe(account.sessionObserver)

	for ev := range account.sessionObserver {
		switch t := ev.(type) {
		case session.Event:
			switch t.Type {
			case session.Connected, session.Disconnected, session.Connecting:
				toggleConnectAndDisconnectMenuItems(t.Session, connectItem, disconnectItem)
			}
		}
	}
}
