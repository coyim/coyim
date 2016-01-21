package gui

import (
	"fmt"
	"log"
	"sync"

	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/coyim/session"
)

// account wraps a Session with GUI functionality
type account struct {
	menu                *gtk.MenuItem
	currentNotification *gtk.InfoBar

	//TODO: Should this be a map of roster.Peer and conversationWindow?
	conversations map[string]*conversationWindow

	session         *session.Session
	sessionObserver chan interface{}

	sync.Mutex
}

type byAccountNameAlphabetic []*account

func (s byAccountNameAlphabetic) Len() int { return len(s) }
func (s byAccountNameAlphabetic) Less(i, j int) bool {
	return s[i].session.GetConfig().Account < s[j].session.GetConfig().Account
}
func (s byAccountNameAlphabetic) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func newAccount(conf *config.ApplicationConfig, currentConf *config.Account) (acc *account, err error) {
	return &account{
		session:       session.NewSession(conf, currentConf),
		conversations: make(map[string]*conversationWindow),
	}, nil
}

func (account *account) ID() string {
	return account.session.GetConfig().ID()
}

func (account *account) getConversationWith(to string) (*conversationWindow, bool) {
	c, ok := account.conversations[to]
	return c, ok
}

func (account *account) createConversationWindow(to string, displaySettings *displaySettings, textBuffer *gtk.TextBuffer) *conversationWindow {
	c := newConversationWindow(account, to, displaySettings, textBuffer)
	account.conversations[to] = c
	return c
}

func (account *account) enableExistingConversationWindows(enable bool) {
	//TODO: should lock account.conversations
	for _, convWindow := range account.conversations {
		if enable {
			convWindow.win.Emit("enable")
		} else {
			convWindow.win.Emit("disable")
		}
	}
}

func (account *account) executeCmd(c interface{}) {
	account.session.CommandManager.ExecuteCmd(c)
}

func (account *account) connected() bool {
	return account.session.IsConnected()
}

func (u *gtkUI) showServerSelectionWindow() error {
	builder := builderForDefinition("AccountRegistration")
	builder.ConnectSignals(map[string]interface{}{
		"response-handler": func(d *gtk.Dialog, resp gtk.ResponseType) {
			defer d.Destroy()

			if resp != gtk.RESPONSE_APPLY {
				return
			}

			obj, _ := builder.GetObject("server")
			iter, _ := obj.(*gtk.ComboBox).GetActiveIter()

			obj, _ = builder.GetObject("servers-model")
			val, _ := obj.(*gtk.ListStore).GetValue(iter, 0)
			server, _ := val.GetString()

			form := &registrationForm{
				parent: u.window,
				server: server,
			}

			saveFn := func() {
				u.addAndSaveAccountConfig(form.conf)
				if acc, ok := u.getAccountByID(form.conf.ID()); ok {
					acc.connect()
				}
			}

			go requestAndRenderRegistrationForm(form.server, form.renderForm, saveFn)
		},
	})

	obj, _ := builder.GetObject("dialog")

	dialog := obj.(*gtk.Dialog)
	dialog.SetTransientFor(u.window)
	dialog.ShowAll()

	return nil
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

func (u *gtkUI) addAndSaveAccountConfig(c *config.Account) {
	accountsLock.Lock()
	u.config.Add(c)
	accountsLock.Unlock()

	err := u.saveConfigInternal()
	if err != nil {
		log.Println("Failed to save config:", err)
	}
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

	menuitem.SetLabel(account.session.GetConfig().Account)

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
	connectAutomaticallyItem.SetActive(account.session.GetConfig().ConnectAutomatically)

	alwaysEncryptItem, _ := gtk.CheckMenuItemNewWithMnemonic(i18n.Local("Always Encrypt Conversation"))
	accountSubMenu.Append(alwaysEncryptItem)
	alwaysEncryptItem.SetActive(account.session.GetConfig().AlwaysEncrypt)

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
	account.executeCmd(connectAccountCmd{account})
}

func (account *account) disconnect() {
	account.executeCmd(disconnectAccountCmd{account})
}

func (account *account) toggleAutoConnect() {
	account.executeCmd(toggleAutoConnectCmd{account})
}

func (account *account) toggleAlwaysEncrypt() {
	account.executeCmd(toggleAlwaysEncryptCmd{account})
}

func (account *account) edit() {
	account.executeCmd(editAccountCmd{account})
}

func (account *account) remove() {
	account.executeCmd(removeAccountCmd{account})
}

func (account *account) buildNotification(template, msg string) *gtk.InfoBar {
	builder := builderForDefinition(template)

	builder.ConnectSignals(map[string]interface{}{
		"handleResponse": func(info *gtk.InfoBar, response gtk.ResponseType) {
			if response != gtk.RESPONSE_CLOSE {
				return
			}

			info.Hide()
			info.Destroy()
		},
	})

	obj, _ := builder.GetObject("infobar")
	infoBar := obj.(*gtk.InfoBar)

	obj, _ = builder.GetObject("message")
	msgLabel := obj.(*gtk.Label)
	msgLabel.SetSelectable(true)

	text := fmt.Sprintf(i18n.Local(msg), account.session.GetConfig().Account)
	msgLabel.SetText(text)

	return infoBar
}

func (account *account) buildConnectionNotification() *gtk.InfoBar {
	return account.buildNotification("ConnectingAccountInfo", "Connecting account\n%s")
}

func (account *account) buildConnectionFailureNotification() *gtk.InfoBar {
	return account.buildNotification("ConnectionFailureNotification", "Connection failure\n%s")
}

func (account *account) removeCurrentNotification() {
	if account.currentNotification != nil {
		account.currentNotification.Hide()
		account.currentNotification.Destroy()
		account.currentNotification = nil
	}
}

func (account *account) removeCurrentNotificationIf(ib *gtk.InfoBar) {
	if account.currentNotification == ib {
		account.currentNotification.Hide()
		account.currentNotification.Destroy()
		account.currentNotification = nil
	}
}

func (account *account) setCurrentNotification(ib *gtk.InfoBar, notificationArea *gtk.Box) {
	account.Lock()
	defer account.Unlock()

	account.removeCurrentNotification()
	account.currentNotification = ib
	notificationArea.Add(ib)
	ib.ShowAll()
}
