package gui

import (
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/coyim/coyim/config"
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/access"
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/xmpp/interfaces"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

// account wraps a Session with GUI functionality
type account struct {
	menu                gtki.MenuItem
	currentNotification gtki.InfoBar
	xmlConsole          gtki.Dialog

	// c contains all conversations. the ones indexed with a "resourced" JID will be locked to that view
	// everything else will be indexed with a bare jid
	c map[string]conversationView

	session access.Session

	sessionObserver         chan interface{}
	connectionEventHandlers []func()
	sessionObserverLock     sync.RWMutex

	delayedConversations     map[string][]func(conversationView)
	delayedConversationsLock sync.Mutex

	askingForPassword bool
	cachedPassword    string

	log coylog.Logger

	events chan interface{}

	sync.RWMutex

	multiUserChatRooms     map[string]*roomView
	multiUserChatRoomsLock sync.RWMutex
}

func (account *account) executeOneDelayed(ui *gtkUI, p string, cv conversationView) {
	for _, f := range account.delayedConversations[p] {
		f(cv)
	}

	delete(account.delayedConversations, p)
}

func samePeer(peer jid.Any, s string) bool {
	sp := jid.Parse(s)
	return peer.NoResource().String() == sp.NoResource().String()
}

func (account *account) executeDelayed(ui *gtkUI, peer jid.Any, targeted bool) {
	account.delayedConversationsLock.Lock()
	defer account.delayedConversationsLock.Unlock()

	ui.NewConversationViewFactory(account, peer, targeted).IfConversationView(func(cv conversationView) {
		if targeted {
			account.executeOneDelayed(ui, peer.String(), cv)
		} else {
			for s := range account.delayedConversations {
				if samePeer(peer, s) {
					account.executeOneDelayed(ui, s, cv)
				}
			}
		}

	}, func() {
		account.log.WithFields(log.Fields{
			"peer":     peer,
			"targeted": targeted,
		}).Warn("Race condition in executeDelayed() - this shouldn't happen")
	})
}

type byAccountNameAlphabetic []*account

func (s byAccountNameAlphabetic) Len() int { return len(s) }
func (s byAccountNameAlphabetic) Less(i, j int) bool {
	return s[i].Account() < s[j].Account()
}
func (s byAccountNameAlphabetic) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func newAccount(conf *config.ApplicationConfig, currentConf *config.Account, sf access.Factory, df interfaces.DialerFactory) *account {
	return &account{
		session:              sf(conf, currentConf, df),
		c:                    make(map[string]conversationView),
		delayedConversations: make(map[string][]func(conversationView)),
		events:               make(chan interface{}),
		multiUserChatRooms:   make(map[string]*roomView),
	}
}

// ID returns the id of the account
func (account *account) ID() string {
	return account.session.GetConfig().ID()
}

// Account returns the JID of the account
func (account *account) Account() string {
	return account.session.GetConfig().Account
}

func (account *account) afterConversationWindowCreated(peer jid.Any, f func(conversationView)) {
	account.delayedConversationsLock.Lock()
	defer account.delayedConversationsLock.Unlock()

	account.delayedConversations[peer.String()] = append(account.delayedConversations[peer.String()], f)
}

func (account *account) enableExistingConversationWindows(enable bool) {
	if account == nil {
		return
	}
	account.RLock()
	defer account.RUnlock()

	for _, cv := range account.c {
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
			u.log.WithFields(log.Fields{
				"feature": "addAccount",
				"account": c.Account,
			}).Warn("Can't add account since you already have an account " +
				"configured with the same name. Remove that account and add it again if you " +
				"really want to overwrite it.")
			u.notify(i18n.Local("Unable to Add Account"), i18n.Localf("Can't add account:\n\n"+
				"You already have an account with this name."))
			return
		}

		_ = u.addAndSaveAccountConfig(c)
		u.log.WithFields(log.Fields{
			"feature": "addAccount",
			"account": c.Account,
		}).Info("Account sucessfully added")
		u.notify(i18n.Local("Account added"), i18n.Localf("%s successfully added.", c.Account))
	})
}

func (u *gtkUI) addAndSaveAccountConfig(c *config.Account) error {
	accountsLock.Lock()

	u.config.Add(c)
	accountsLock.Unlock()

	err := u.saveConfigInternal()
	if err != nil {
		u.log.WithField("account", c.Account).WithError(err).Warn("Failed to save config")
	}

	doInUIThread(func() {
		if u.window != nil {
			_, _ = u.window.Emit(accountChangedSignal.String())
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

func (account *account) appendMenuTo(u *gtkUI, submenu gtki.Menu) {
	if account.menu != nil {
		account.destroyMenu()
	}

	account.buildAccountSubmenu(u)
	account.menu.Show()
	submenu.Append(account.menu)
}

func (account *account) runSessionObserver(u *gtkUI) {
	for ev := range account.sessionObserver {
		switch t := ev.(type) {
		case events.Event:
			switch t.Type {
			case events.Connected, events.Disconnected, events.Connecting:
				doInUIThread(func() {
					account.sessionObserverLock.RLock()
					defer account.sessionObserverLock.RUnlock()
					u.updateGlobalMenuStatus()
					for _, ff := range account.connectionEventHandlers {
						ff()
					}
				})
			}
		}
	}
}

func (account *account) observeConnectionEvents(u *gtkUI, f func()) {
	account.sessionObserverLock.Lock()
	defer account.sessionObserverLock.Unlock()

	if account.sessionObserver == nil {
		account.sessionObserver = make(chan interface{})
		account.session.Subscribe(account.sessionObserver)
		go account.runSessionObserver(u)
	}
	account.connectionEventHandlers = append(account.connectionEventHandlers, f)
}

func (account *account) createCheckConnectionItem(u *gtkUI) gtki.MenuItem {
	checkConnectionItem, _ := g.gtk.MenuItemNewWithMnemonic(i18n.Local("_Check Connection"))
	_, _ = checkConnectionItem.Connect("activate", func() {
		account.session.SendPing()
	})
	checkConnectionItem.SetSensitive(account.session.IsConnected())

	account.observeConnectionEvents(u, func() {
		checkConnectionItem.SetSensitive(account.session.IsConnected())
	})
	return checkConnectionItem
}

func (account *account) createConnectItem(u *gtkUI) gtki.MenuItem {
	connectItem, _ := g.gtk.MenuItemNewWithMnemonic(i18n.Local("_Connect"))
	_, _ = connectItem.Connect("activate", account.Connect)
	connectItem.SetSensitive(account.session.IsDisconnected())
	account.observeConnectionEvents(u, func() {
		connectItem.SetSensitive(account.session.IsDisconnected())
	})
	return connectItem
}

func (account *account) createDisconnectItem(u *gtkUI) gtki.MenuItem {
	disconnectItem, _ := g.gtk.MenuItemNewWithMnemonic(i18n.Local("_Disconnect"))
	_, _ = disconnectItem.Connect("activate", func() {
		account.session.SetWantToBeOnline(false)
		account.disconnect()
	})
	disconnectItem.SetSensitive(!account.session.IsDisconnected())
	account.observeConnectionEvents(u, func() {
		disconnectItem.SetSensitive(!account.session.IsDisconnected())
	})
	return disconnectItem
}

func (account *account) createSeparatorItem() gtki.MenuItem {
	sep, _ := g.gtk.SeparatorMenuItemNew()
	return sep
}

func (account *account) createConnectionItem(u *gtkUI) gtki.MenuItem {
	connInfoItem, _ := g.gtk.MenuItemNewWithMnemonic(i18n.Local("Connection _information..."))
	_, _ = connInfoItem.Connect("activate", account.connectionInfo)
	connInfoItem.SetSensitive(account.session.IsConnected())
	account.observeConnectionEvents(u, func() {
		connInfoItem.SetSensitive(account.session.IsConnected())
	})
	return connInfoItem
}

func (account *account) createEditItem() gtki.MenuItem {
	editItem, _ := g.gtk.MenuItemNewWithMnemonic(i18n.Local("_Edit..."))
	_, _ = editItem.Connect("activate", account.edit)
	return editItem
}

func (account *account) createChangePasswordItem(u *gtkUI) gtki.MenuItem {
	changePasswordItem, _ := g.gtk.MenuItemNewWithMnemonic(i18n.Local("_Change Password..."))
	_, _ = changePasswordItem.Connect("activate", account.changePassword)
	changePasswordItem.SetSensitive(account.session.IsConnected())
	return changePasswordItem
}

func (account *account) createRemoveItem() gtki.MenuItem {
	removeItem, _ := g.gtk.MenuItemNewWithMnemonic(i18n.Local("_Remove"))
	_, _ = removeItem.Connect("activate", account.remove)
	return removeItem
}

func (account *account) createConnectAutomaticallyItem() gtki.MenuItem {
	connectAutomaticallyItem, _ := g.gtk.CheckMenuItemNewWithMnemonic(i18n.Local("Connect _Automatically"))
	connectAutomaticallyItem.SetActive(account.session.GetConfig().ConnectAutomatically)
	_, _ = connectAutomaticallyItem.Connect("activate", account.toggleAutoConnect)
	return connectAutomaticallyItem
}

func (account *account) createAlwaysEncryptItem() gtki.MenuItem {
	alwaysEncryptItem, _ := g.gtk.CheckMenuItemNewWithMnemonic(i18n.Local("Always Encrypt Conversation"))
	alwaysEncryptItem.SetActive(account.session.GetConfig().AlwaysEncrypt)
	_, _ = alwaysEncryptItem.Connect("activate", account.toggleAlwaysEncrypt)
	return alwaysEncryptItem
}

func (account *account) createDumpInfoItem(r *roster) gtki.MenuItem {
	dumpInfoItem, _ := g.gtk.MenuItemNewWithMnemonic(i18n.Local("Dump info"))
	_, _ = dumpInfoItem.Connect("activate", func() {
		r.ui.accountManager.debugPeersFor(account)
	})
	return dumpInfoItem
}

func (account *account) createXMLConsoleItem(parent gtki.Window) gtki.MenuItem {
	if account.xmlConsole == nil {
		account.xmlConsole = newXMLConsoleView(account.session.GetInMemoryLog())
		account.xmlConsole.SetTransientFor(parent)
		account.xmlConsole.SetTitle(strings.Replace(account.xmlConsole.GetTitle(), "ACCOUNT_NAME", account.Account(), -1))
	}

	consoleItem, _ := g.gtk.MenuItemNewWithMnemonic(i18n.Local("XML Console"))
	_, _ = consoleItem.Connect("activate", account.xmlConsole.ShowAll)

	return consoleItem
}

func (account *account) createSubmenu(u *gtkUI) gtki.Menu {
	m, _ := g.gtk.MenuNew()

	m.Append(account.createConnectItem(u))
	m.Append(account.createDisconnectItem(u))
	m.Append(account.createCheckConnectionItem(u))
	m.Append(account.createSeparatorItem())
	m.Append(account.createConnectionItem(u))
	m.Append(account.createEditItem())
	m.Append(account.createChangePasswordItem(u))
	m.Append(account.createRemoveItem())
	m.Append(account.createSeparatorItem())
	m.Append(account.createConnectAutomaticallyItem())
	m.Append(account.createAlwaysEncryptItem())

	return m
}

func (account *account) buildAccountSubmenu(u *gtkUI) {
	menuitem, _ := g.gtk.MenuItemNew()
	menuitem.SetLabel(account.Account())
	menuitem.SetSubmenu(account.createSubmenu(u))
	account.menu = menuitem
}

func (account *account) Connect() {
	account.session.SetWantToBeOnline(true)
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

func (account *account) changePassword() {
	account.executeCmd(changePasswordAccountCmd{account})
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
		"on_more_info": func() {
			if moreInfo != nil {
				moreInfo()
			}
		},
		"on_close": func() {
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
	return account.buildNotification("ConnectingAccountInfo", i18n.Localf("Connecting account\n%s", account.Account()), nil)
}

func (account *account) buildConnectionFailureNotification(moreInfo func()) gtki.InfoBar {
	return account.buildNotification("ConnectionFailureNotification", i18n.Localf("Connection failure\n%s", account.Account()), moreInfo)
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
