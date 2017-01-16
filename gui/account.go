package gui

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/coyim/servers"
	"github.com/twstrike/coyim/session/access"
	"github.com/twstrike/coyim/session/events"
	"github.com/twstrike/coyim/xmpp/interfaces"
	"github.com/twstrike/gotk3adapter/gtki"
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
	cachedPassword    string

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

func (account *account) getConversationWith(to string, ui *gtkUI) (conversationView, bool) {
	c, ok := account.conversations[to]

	if ok {
		_, unifiedType := c.(*conversationStackItem)

		if ui.settings.GetSingleWindow() && !unifiedType {
			cv1 := c.(*conversationWindow)
			c = ui.unified.createConversation(account, to, cv1.conversationPane)
			account.conversations[to] = c
		} else if !ui.settings.GetSingleWindow() && unifiedType {
			cv1 := c.(*conversationStackItem)
			c = newConversationWindow(account, to, ui, cv1.conversationPane)
			account.conversations[to] = c
		}
	}

	return c, ok
}

func (account *account) createConversationView(to string, ui *gtkUI) conversationView {
	var cv conversationView
	if ui.settings.GetSingleWindow() {
		cv = ui.unified.createConversation(account, to, nil)
	} else {
		cv = newConversationWindow(account, to, ui, nil)
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

func (u *gtkUI) showServerSelectionWindow() {
	builder := newBuilder("AccountRegistration")
	d := builder.getObj("dialog").(gtki.Dialog)
	d.SetTransientFor(u.window)
	message := builder.getObj("message").(gtki.Label)
	serverBox := builder.getObj("server").(gtki.ComboBoxText)
	applyBtn := builder.getObj("btn-apply").(gtki.Button)

	for _, s := range servers.GetServersForRegistration() {
		serverBox.AppendText(s.Name)
	}

	serverBox.SetActive(0)
	destroyed := false

	builder.ConnectSignals(map[string]interface{}{
		"on_save_signal": func() {
			message.SetLabel(i18n.Local("Connecting to server for registration..."))
			applyBtn.SetSensitive(false)
			serverBox.SetSensitive(false)

			server := serverBox.GetActiveText()

			form := &registrationForm{
				parent: u.window,
				server: server,
			}

			saveFn := func() {
				doInUIThread(func() {
					if !destroyed {
						d.Destroy()
						destroyed = true
					}
				})
				u.addAndSaveAccountConfig(form.conf)
				if acc, ok := u.getAccountByID(form.conf.ID()); ok {
					acc.session.SetWantToBeOnline(true)
					acc.Connect()
				}
			}

			renderFn := func(title, instructions string, fields []interface{}) error {
				doInUIThread(func() {
					if !destroyed {
						d.Destroy()
						destroyed = true
					}
				})
				return form.renderForm(title, instructions, fields)
			}

			errorFn := func(err error) {
				log.Printf("Error when trying to get registration form: %v", err)
				doInUIThread(func() {
					message.SetLabel(i18n.Local("We had an error when trying to contact the server. Please correct your server choice and try again. The registration process currently requires Tor in order to ensure your safety. If you don't have Tor turned on, make sure to do so."))
					applyBtn.SetSensitive(true)
					serverBox.SetSensitive(true)
				})
			}

			go requestAndRenderRegistrationForm(form.server, renderFn, saveFn, errorFn, u.dialerFactory, u.unassociatedVerifier())
		},

		"on_cancel_signal": func() {
			d.Destroy()
		},
	})

	d.ShowAll()
}

func (u *gtkUI) showAddAccountWindow() {
	c, _ := config.NewAccount()

	u.accountDialog(nil, c, func() {
		u.addAndSaveAccountConfig(c)
		u.notify(i18n.Local("Account added"), fmt.Sprintf(i18n.Local("The account %s was added successfully."), c.Account))
	})
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
		r.debugPrintRosterFor(account.session.GetConfig().Account)
	})
	return dumpInfoItem
}

func (account *account) createXMLConsoleItem(r *roster) gtki.MenuItem {
	consoleItem, _ := g.gtk.MenuItemNewWithMnemonic(i18n.Local("XML Console"))
	consoleItem.Connect("activate", func() {
		builder := newBuilder("XMLConsole")
		console := builder.getObj("XMLConsole").(gtki.Dialog)
		buf := builder.getObj("consoleContent").(gtki.TextBuffer)
		console.SetTransientFor(r.ui.window)
		console.SetTitle(strings.Replace(console.GetTitle(), "ACCOUNT_NAME", account.session.GetConfig().Account, -1))
		log := account.session.GetInMemoryLog()

		buf.Delete(buf.GetStartIter(), buf.GetEndIter())
		if log != nil {
			buf.Insert(buf.GetEndIter(), log.String())
		}

		builder.ConnectSignals(map[string]interface{}{
			"on_refresh_signal": func() {
				buf.Delete(buf.GetStartIter(), buf.GetEndIter())
				if log != nil {
					buf.Insert(buf.GetEndIter(), log.String())
				}
			},
			"on_close_signal": func() {
				console.Destroy()
			},
		})

		console.ShowAll()
	})
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

func (account *account) buildNotification(template, msg string, u *gtkUI, moreInfo func()) gtki.InfoBar {
	builder := newBuilder(template)

	infoBar := builder.getObj("infobar").(gtki.InfoBar)

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

func (account *account) buildConnectionNotification(u *gtkUI) gtki.InfoBar {
	return account.buildNotification("ConnectingAccountInfo", fmt.Sprintf(i18n.Local("Connecting account\n%s"), account.session.GetConfig().Account), u, nil)
}

func (account *account) buildConnectionFailureNotification(u *gtkUI, moreInfo func()) gtki.InfoBar {
	return account.buildNotification("ConnectionFailureNotification", fmt.Sprintf(i18n.Local("Connection failure\n%s"), account.session.GetConfig().Account), u, moreInfo)
}

func (account *account) buildTorNotRunningNotification(u *gtkUI) gtki.InfoBar {
	return account.buildNotification("TorNotRunningNotification", i18n.Local("Tor is not currently running"), u, nil)
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
