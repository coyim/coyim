package gui

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"../config"
	"../i18n"
	"../session"
)

// account wraps a Session with GUI functionality
type account struct {
	menu                   *gtk.MenuItem
	connectionNotification *gtk.InfoBar

	session         *session.Session
	sessionObserver chan interface{}
}

type byAccountNameAlphabetic []*account

func (s byAccountNameAlphabetic) Len() int { return len(s) }
func (s byAccountNameAlphabetic) Less(i, j int) bool {
	return s[i].session.CurrentAccount.Account < s[j].session.CurrentAccount.Account
}
func (s byAccountNameAlphabetic) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func newAccount(conf *config.ApplicationConfig, currentConf *config.Account) (acc *account, err error) {
	return &account{
		session: session.NewSession(conf, currentConf),
	}, nil
}

func (account *account) executeCmd(c interface{}) {
	account.session.CommandManager.ExecuteCmd(c)
}

func (account *account) connected() bool {
	return account.session.ConnStatus == session.CONNECTED
}

func (u *gtkUI) showAddAccountWindow() error {
	c, err := config.NewAccount()
	if err != nil {
		return err
	}

	u.accountDialog(c, func() {
		u.addAndSaveAccountConfig(c)
	})

	return nil
}

func (account *account) destroyMenu() {
	close(account.sessionObserver)

	account.menu.Destroy()
	account.menu = nil
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

	alwaysEncryptItem, _ := gtk.CheckMenuItemNewWithMnemonic(i18n.Local("Always Encrypt Conversation"))
	accountSubMenu.Append(alwaysEncryptItem)
	alwaysEncryptItem.SetActive(account.session.CurrentAccount.AlwaysEncrypt)

	toggleConnectAndDisconnectMenuItems(account.session, connectItem, disconnectItem)

	connectItem.Connect("activate", account.connect)
	disconnectItem.Connect("activate", account.disconnect)
	editItem.Connect("activate", account.edit)
	removeItem.Connect("activate", account.remove)
	connectAutomaticallyItem.Connect("activate", account.toggleAutoConnect)
	alwaysEncryptItem.Connect("activate", account.toggleAlwaysEncrypt)

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

func (account *account) connect() {
	account.executeCmd(connectAccountCmd(account))
}

func (account *account) disconnect() {
	account.executeCmd(disconnectAccountCmd(account))
}

func (account *account) toggleAutoConnect() {
	account.executeCmd(toggleAutoConnectCmd(account))
}

func (account *account) toggleAlwaysEncrypt() {
	account.executeCmd(toggleAlwaysEncryptCmd(account))
}

func (account *account) edit() {
	account.executeCmd(editAccountCmd(account))
}

func (account *account) remove() {
	account.executeCmd(removeAccountCmd(account))
}

func (account *account) buildConnectionNotification() {
	builder := builderForDefinition("ConnectingAccountInfo")

	obj, _ := builder.GetObject("infobar")
	infoBar := obj.(*gtk.InfoBar)

	obj, _ = builder.GetObject("message")
	msg := obj.(*gtk.Label)

	text := fmt.Sprintf(i18n.Local("Connecting account\n%s"),
		account.session.CurrentAccount.Account)

	msg.SetText(text)

	account.connectionNotification = infoBar
}

func (account *account) removeConnectionNotification() {
	account.connectionNotification.Hide()
	account.connectionNotification.Destroy()
	account.connectionNotification = nil
}
