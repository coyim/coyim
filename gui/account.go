package gui

import (
	"log"
	"strings"
	"sync"

	"github.com/coyim/coyim/config"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/access"
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/xmpp/interfaces"
	"github.com/coyim/coyim/xmpp/utils"
	"github.com/coyim/gotk3adapter/gtki"
)

// account wraps a Session with GUI functionality
type account struct {
	menu                gtki.MenuItem
	currentNotification gtki.InfoBar
	xmlConsole          gtki.Dialog

	//TODO: Should this be a map of roster.Peer and conversationView?
	conversations map[string]conversationView

	session access.Session

	sessionObserver         chan interface{}
	connectionEventHandlers []func()
	sessionObserverLock     sync.RWMutex

	delayedConversations     map[string][]func(conversationView)
	delayedConversationsLock sync.Mutex

	askingForPassword bool
	cachedPassword    string

	sync.RWMutex
}

func (account *account) getConversationView(fullJID string) (conversationView, bool) {
	account.RLock()
	defer account.RUnlock()

	c, ok := account.conversations[fullJID]
	return c, ok
}

func (account *account) executeDelayed(fullJID string) {
	account.delayedConversationsLock.Lock()
	defer account.delayedConversationsLock.Unlock()

	cv, ok := account.getConversationView(fullJID)
	if !ok {
		panic("race condition")
	}

	for _, f := range account.delayedConversations[fullJID] {
		f(cv)
	}

	delete(account.delayedConversations, fullJID)
}

func (account *account) setConversationView(fullJID string, c conversationView) {
	defer account.executeDelayed(fullJID)

	account.Lock()
	if c, ok := account.conversations[fullJID]; ok {
		c.destroy() // We dont want the previous conversation to be visible after this moment
	}

	account.conversations[fullJID] = c
	account.Unlock()
}

type byAccountNameAlphabetic []*account

func (s byAccountNameAlphabetic) Len() int { return len(s) }
func (s byAccountNameAlphabetic) Less(i, j int) bool {
	return s[i].session.GetConfig().Account < s[j].session.GetConfig().Account
}
func (s byAccountNameAlphabetic) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func newAccount(conf *config.ApplicationConfig, currentConf *config.Account, sf access.Factory, df interfaces.DialerFactory) *account {
	return &account{
		session:              sf(conf, currentConf, df),
		conversations:        make(map[string]conversationView),
		delayedConversations: make(map[string][]func(conversationView)),
	}
}

func (account *account) ID() string {
	return account.session.GetConfig().ID()
}

func (account *account) getConversationWith(to, resource string, ui *gtkUI) (c conversationView, ok bool) {
	fullJid := utils.ComposeFullJid(to, resource)
	c, ok = account.getConversationView(fullJid)
	if !ok {
		return
	}

	_, unifiedType := c.(*conversationStackItem)
	if ui.settings.GetSingleWindow() == unifiedType {
		//This is only true when c is nil
		return
	}

	// Why not a simple type switch?
	if !unifiedType {
		cv1 := c.(*conversationWindow)
		c = ui.unified.createConversation(account, to, resource, cv1.conversationPane)
	} else {
		cv1 := c.(*conversationStackItem)
		c = newConversationWindow(account, fullJid, ui, cv1.conversationPane)
	}

	account.setConversationView(fullJid, c)
	return
}
func (account *account) afterConversationWindowCreated(to string, f func(conversationView)) {
	account.delayedConversationsLock.Lock()
	defer account.delayedConversationsLock.Unlock()

	account.delayedConversations[to] = append(account.delayedConversations[to], f)
}

func (account *account) enableExistingConversationWindows(enable bool) {
	if account == nil {
		return
	}

	account.RLock()
	defer account.RUnlock()

	for _, cv := range account.conversations {
		cv.setEnabled(enable)
	}
}

func (account *account) executeCmd(c interface{}) {
	account.session.CommandManager().ExecuteCmd(c)
}

func (account *account) connected() bool {
	return account.session.IsConnected()
}

func (u *gtkUI) showAddAccountWindow() {
	c, _ := config.NewAccount()

	u.accountDialog(nil, c, func() {
		_, exists := u.config.GetAccount(c.Account)
		if exists {
			log.Printf("[add account] Can't add account %s since you already have an account "+
				"configured with the same name. Remove that account and add it again if you "+
				"really want to overwrite it.", c.Account)
			u.notify(i18n.Local("Unable to Add Account"), i18n.Localf("Can't add account:\n\n"+
				"You already have an account with this name."))
			return
		}

		u.addAndSaveAccountConfig(c)
		log.Printf("[add account] Account sucessfully added.")
		u.notify(i18n.Local("Account added"), i18n.Localf("%s successfully added.", c.Account))
	})
}

func (u *gtkUI) addAndSaveAccountConfig(c *config.Account) error {
	accountsLock.Lock()

	u.config.Add(c)
	accountsLock.Unlock()

	err := u.saveConfigInternal()
	if err != nil {
		log.Println("Failed to save config:", err)
	}

	doInUIThread(func() {
		if u.window != nil {
			u.window.Emit(accountChangedSignal.String())
		}
	})

	return err
}

func (account *account) destroyMenu() {
	account.sessionObserverLock.Lock()
	defer account.sessionObserverLock.Unlock()

	close(account.sessionObserver)
	account.sessionObserver = nil
	account.connectionEventHandlers = nil

	account.menu.Destroy()
	account.menu = nil
}

func (account *account) appendMenuTo(submenu gtki.Menu) {
	if account.menu != nil {
		account.destroyMenu()
	}

	account.buildAccountSubmenu()
	account.menu.Show()
	submenu.Append(account.menu)
}

func (account *account) runSessionObserver() {
	for ev := range account.sessionObserver {
		switch t := ev.(type) {
		case events.Event:
			switch t.Type {
			case events.Connected, events.Disconnected, events.Connecting:
				doInUIThread(func() {
					account.sessionObserverLock.RLock()
					for _, ff := range account.connectionEventHandlers {
						ff()
					}
					account.sessionObserverLock.RUnlock()
				})
			}
		}
	}
}

func (account *account) observeConnectionEvents(f func()) {
	account.sessionObserverLock.Lock()
	defer account.sessionObserverLock.Unlock()

	if account.sessionObserver == nil {
		account.sessionObserver = make(chan interface{})
		account.session.Subscribe(account.sessionObserver)
		go account.runSessionObserver()
	}
	account.connectionEventHandlers = append(account.connectionEventHandlers, f)
}

func (account *account) createCheckConnectionItem() gtki.MenuItem {
	checkConnectionItem, _ := g.gtk.MenuItemNewWithMnemonic(i18n.Local("_Check Connection"))
	checkConnectionItem.Connect("activate", func() {
		account.session.SendPing()
	})
	checkConnectionItem.SetSensitive(account.session.IsConnected())

	account.observeConnectionEvents(func() {
		checkConnectionItem.SetSensitive(account.session.IsConnected())
	})
	return checkConnectionItem
}

func (account *account) createConnectItem() gtki.MenuItem {
	connectItem, _ := g.gtk.MenuItemNewWithMnemonic(i18n.Local("_Connect"))
	connectItem.Connect("activate", func() {
		account.session.SetWantToBeOnline(true)
		account.Connect()
	})
	connectItem.SetSensitive(account.session.IsDisconnected())
	account.observeConnectionEvents(func() {
		connectItem.SetSensitive(account.session.IsDisconnected())
	})
	return connectItem
}

func (account *account) createDisconnectItem() gtki.MenuItem {
	disconnectItem, _ := g.gtk.MenuItemNewWithMnemonic(i18n.Local("_Disconnect"))
	disconnectItem.Connect("activate", func() {
		account.session.SetWantToBeOnline(false)
		account.disconnect()
	})
	disconnectItem.SetSensitive(!account.session.IsDisconnected())
	account.observeConnectionEvents(func() {
		disconnectItem.SetSensitive(!account.session.IsDisconnected())
	})
	return disconnectItem
}

func (account *account) createSeparatorItem() gtki.MenuItem {
	sep, _ := g.gtk.SeparatorMenuItemNew()
	return sep
}

func (account *account) createConnectionItem() gtki.MenuItem {
	connInfoItem, _ := g.gtk.MenuItemNewWithMnemonic(i18n.Local("Connection _information..."))
	connInfoItem.Connect("activate", account.connectionInfo)
	connInfoItem.SetSensitive(account.session.IsConnected())
	account.observeConnectionEvents(func() {
		connInfoItem.SetSensitive(account.session.IsConnected())
	})
	return connInfoItem
}

func (account *account) createEditItem() gtki.MenuItem {
	editItem, _ := g.gtk.MenuItemNewWithMnemonic(i18n.Local("_Edit..."))
	editItem.Connect("activate", account.edit)
	return editItem
}

func (account *account) createRemoveItem() gtki.MenuItem {
	removeItem, _ := g.gtk.MenuItemNewWithMnemonic(i18n.Local("_Remove"))
	removeItem.Connect("activate", account.remove)
	return removeItem
}

func (account *account) createConnectAutomaticallyItem() gtki.MenuItem {
	connectAutomaticallyItem, _ := g.gtk.CheckMenuItemNewWithMnemonic(i18n.Local("Connect _Automatically"))
	connectAutomaticallyItem.SetActive(account.session.GetConfig().ConnectAutomatically)
	connectAutomaticallyItem.Connect("activate", account.toggleAutoConnect)
	return connectAutomaticallyItem
}

func (account *account) createAlwaysEncryptItem() gtki.MenuItem {
	alwaysEncryptItem, _ := g.gtk.CheckMenuItemNewWithMnemonic(i18n.Local("Always Encrypt Conversation"))
	alwaysEncryptItem.SetActive(account.session.GetConfig().AlwaysEncrypt)
	alwaysEncryptItem.Connect("activate", account.toggleAlwaysEncrypt)
	return alwaysEncryptItem
}

func (account *account) createDumpInfoItem(r *roster) gtki.MenuItem {
	dumpInfoItem, _ := g.gtk.MenuItemNewWithMnemonic(i18n.Local("Dump info"))
	dumpInfoItem.Connect("activate", func() {
		r.ui.accountManager.debugPeersFor(account)
	})
	return dumpInfoItem
}

func (account *account) createXMLConsoleItem(parent gtki.Window) gtki.MenuItem {
	if account.xmlConsole == nil {
		account.xmlConsole = newXMLConsoleView(account.session.GetInMemoryLog())
		account.xmlConsole.SetTransientFor(parent)
		account.xmlConsole.SetTitle(strings.Replace(account.xmlConsole.GetTitle(), "ACCOUNT_NAME", account.session.GetConfig().Account, -1))
	}

	consoleItem, _ := g.gtk.MenuItemNewWithMnemonic(i18n.Local("XML Console"))
	consoleItem.Connect("activate", account.xmlConsole.ShowAll)

	return consoleItem
}

func (account *account) createSubmenu() gtki.Menu {
	m, _ := g.gtk.MenuNew()

	m.Append(account.createConnectItem())
	m.Append(account.createDisconnectItem())
	m.Append(account.createCheckConnectionItem())
	m.Append(account.createSeparatorItem())
	m.Append(account.createConnectionItem())
	m.Append(account.createEditItem())
	m.Append(account.createRemoveItem())
	m.Append(account.createSeparatorItem())
	m.Append(account.createConnectAutomaticallyItem())
	m.Append(account.createAlwaysEncryptItem())

	return m
}

func (account *account) buildAccountSubmenu() {
	menuitem, _ := g.gtk.MenuItemNew()
	menuitem.SetLabel(account.session.GetConfig().Account)
	menuitem.SetSubmenu(account.createSubmenu())
	account.menu = menuitem
}

func (account *account) Connect() {
	account.executeCmd(connectAccountCmd{account})
}

func (account *account) disconnect() {
	account.session.SetWantToBeOnline(false)
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

func (account *account) connectionInfo() {
	account.executeCmd(connectionInfoCmd{account})
}

func (account *account) remove() {
	account.executeCmd(removeAccountCmd{account})
}

func (account *account) buildNotification(template, msg string, moreInfo func()) gtki.InfoBar {
	builder := newBuilder(template)

	infoBar := builder.getObj("infobar").(gtki.InfoBar)
	box := builder.getObj("box").(gtki.InfoBar)

	prov := providerWithCSS("box { border: none; }")
	updateWithStyle(box, prov)

	builder.ConnectSignals(map[string]interface{}{
		"on_more_info_signal": func() {
			if moreInfo != nil {
				moreInfo()
			}
		},
		"on_close_signal": func() {
			infoBar.Hide()
			infoBar.Destroy()
		},
		"handleResponse": func(info gtki.InfoBar, response gtki.ResponseType) {
			if response != gtki.RESPONSE_CLOSE {
				return
			}

			info.Hide()
			info.Destroy()
		},
	})

	msgLabel := builder.getObj("message").(gtki.Label)
	msgLabel.SetText(msg)

	return infoBar
}

func (account *account) buildConnectionNotification() gtki.InfoBar {
	return account.buildNotification("ConnectingAccountInfo", i18n.Localf("Connecting account\n%s", account.session.GetConfig().Account), nil)
}

func (account *account) buildConnectionFailureNotification(moreInfo func()) gtki.InfoBar {
	return account.buildNotification("ConnectionFailureNotification", i18n.Localf("Connection failure\n%s", account.session.GetConfig().Account), moreInfo)
}

func (account *account) buildTorNotRunningNotification(moreInfo func()) gtki.InfoBar {
	return account.buildNotification("TorNotRunningNotification", i18n.Local("Tor is not currently running"), moreInfo)
}

func (account *account) removeCurrentNotification() {
	if account.currentNotification != nil {
		account.currentNotification.Hide()
		account.currentNotification.Destroy()
		account.currentNotification = nil
	}
}

func (account *account) removeCurrentNotificationIf(ib gtki.InfoBar) {
	if account.currentNotification == ib {
		account.currentNotification.Hide()
		account.currentNotification.Destroy()
		account.currentNotification = nil
	}
}

func (account *account) setCurrentNotification(ib gtki.InfoBar, notificationArea gtki.Box) {
	account.Lock()
	defer account.Unlock()

	account.removeCurrentNotification()
	account.currentNotification = ib
	notificationArea.Add(ib)
	ib.ShowAll()
}

func (account *account) IsAskingForPassword() bool {
	return account.askingForPassword
}

func (account *account) AskForPassword() {
	account.askingForPassword = true
}

func (account *account) AskedForPassword() {
	account.askingForPassword = false
}
