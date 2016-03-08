package gui

import (
	"fmt"
	"log"
	"sync"

	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtki"
	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/coyim/session/access"
	"github.com/twstrike/coyim/session/events"
	"github.com/twstrike/coyim/xmpp/interfaces"
)

// account wraps a Session with GUI functionality
type account struct {
	menu                gtki.MenuItem
	currentNotification gtki.InfoBar

	//TODO: Should this be a map of roster.Peer and conversationView?
	conversations map[string]conversationView

	session access.Session

	sessionObserver         chan interface{}
	connectionEventHandlers []func()
	sessionObserverLock     sync.RWMutex

	delayedConversations     map[string][]func(conversationView)
	delayedConversationsLock sync.Mutex

	askingForPassword bool

	sync.Mutex
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

func (account *account) getConversationWith(to string) (conversationView, bool) {
	c, ok := account.conversations[to]
	return c, ok
}

func (account *account) createConversationView(to string, ui *gtkUI) conversationView {
	var cv conversationView
	if *config.SingleWindowFlag {
		cv = ui.unified.createConversation(account, to)
	} else {
		cv = newConversationWindow(account, to, ui)
	}
	account.conversations[to] = cv

	account.delayedConversationsLock.Lock()
	defer account.delayedConversationsLock.Unlock()
	for _, f := range account.delayedConversations[to] {
		f(cv)
	}
	delete(account.delayedConversations, to)

	return cv
}

func (account *account) afterConversationWindowCreated(to string, f func(conversationView)) {
	account.delayedConversationsLock.Lock()
	defer account.delayedConversationsLock.Unlock()

	account.delayedConversations[to] = append(account.delayedConversations[to], f)
}

func (account *account) enableExistingConversationWindows(enable bool) {
	if account != nil {
		account.Lock()
		defer account.Unlock()

		for _, cv := range account.conversations {
			cv.setEnabled(enable)
		}
	}
}

func (account *account) executeCmd(c interface{}) {
	account.session.CommandManager().ExecuteCmd(c)
}

func (account *account) connected() bool {
	return account.session.IsConnected()
}

func (u *gtkUI) showServerSelectionWindow() error {
	builder := builderForDefinition("AccountRegistration")
	obj, _ := builder.GetObject("dialog")

	d := obj.(gtki.Dialog)
	defer d.Destroy()

	d.SetTransientFor(u.window)
	d.ShowAll()

	resp := d.Run()
	if gtki.ResponseType(resp) != gtki.RESPONSE_APPLY {
		return nil
	}

	obj, _ = builder.GetObject("server")
	iter, _ := obj.(gtki.ComboBox).GetActiveIter()

	obj, _ = builder.GetObject("servers-model")
	val, _ := obj.(gtki.ListStore).GetValue(iter, 0)
	server, _ := val.GetString()

	form := &registrationForm{
		parent: u.window,
		server: server,
	}

	saveFn := func() {
		u.addAndSaveAccountConfig(form.conf)
		if acc, ok := u.getAccountByID(form.conf.ID()); ok {
			acc.session.SetWantToBeOnline(true)
			acc.Connect()
		}
	}

	go requestAndRenderRegistrationForm(form.server, form.renderForm, saveFn, u.dialerFactory, u.unassociatedVerifier())

	return nil
}

func (u *gtkUI) showAddAccountWindow() error {
	c, err := config.NewAccount()
	if err != nil {
		return err
	}

	u.accountDialog(nil, c, func() {
		u.addAndSaveAccountConfig(c)
		u.notify(i18n.Local("Account added"), fmt.Sprintf(i18n.Local("The account %s was added successfully."), c.Account))
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
	doInUIThread(func() {
		if u.window != nil {
			u.window.Emit(accountChangedSignal.String())
		}
	})
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
				account.sessionObserverLock.RLock()
				for _, ff := range account.connectionEventHandlers {
					ff()
				}
				account.sessionObserverLock.RUnlock()
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
		r.debugPrintRosterFor(account.session.GetConfig().Account)
	})
	return dumpInfoItem
}

func (account *account) createSubmenu() gtki.Menu {
	m, _ := g.gtk.MenuNew()

	m.Append(account.createConnectItem())
	m.Append(account.createDisconnectItem())
	m.Append(account.createSeparatorItem())
	m.Append(account.createEditItem())
	m.Append(account.createRemoveItem())
	m.Append(account.createSeparatorItem())
	m.Append(account.createConnectAutomaticallyItem())
	m.Append(account.createAlwaysEncryptItem())
	m.Append(account.createSeparatorItem())

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

func (account *account) remove() {
	account.executeCmd(removeAccountCmd{account})
}

func (account *account) buildNotification(template, msg string) gtki.InfoBar {
	builder := builderForDefinition(template)

	builder.ConnectSignals(map[string]interface{}{
		"handleResponse": func(info gtki.InfoBar, response gtki.ResponseType) {
			if response != gtki.RESPONSE_CLOSE {
				return
			}

			info.Hide()
			info.Destroy()
		},
	})

	obj, _ := builder.GetObject("infobar")
	infoBar := obj.(gtki.InfoBar)

	obj, _ = builder.GetObject("message")
	msgLabel := obj.(gtki.Label)
	msgLabel.SetSelectable(true)
	msgLabel.SetText(msg)

	return infoBar
}

func (account *account) buildConnectionNotification() gtki.InfoBar {
	return account.buildNotification("ConnectingAccountInfo", fmt.Sprintf(i18n.Local("Connecting account\n%s"), account.session.GetConfig().Account))
}

func (account *account) buildConnectionFailureNotification() gtki.InfoBar {
	return account.buildNotification("ConnectionFailureNotification", fmt.Sprintf(i18n.Local("Connection failure\n%s"), account.session.GetConfig().Account))
}

func (account *account) buildTorNotRunningNotification() gtki.InfoBar {
	return account.buildNotification("TorNotRunningNotification", i18n.Local("Tor is not currently running"))
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
