package gui

import (
	"log"
	"strings"
	"sync"

	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/coyim/servers"
	"github.com/twstrike/coyim/session/access"
	"github.com/twstrike/coyim/session/events"
	"github.com/twstrike/coyim/xmpp"
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

var (
	torErrorMessage = "The registration process currently requires Tor in order to ensure your safety.\n\n" +
		"You don't have Tor running. Please, start it.\n\n"
	torLogMessage         = "We had an error when trying to register your account: Tor is not running. %v"
	storeAccountInfoError = "We had an error when trying to store your account information."
	storeAccountInfoLog   = "We had an error when trying to store your account information. %v"
	contactServerError    = "Could not contact the server.\n\n Please correct your server choice and try again."
	contactServerLog      = "Error when trying to get registration form: %v"
	requiredFieldsError   = "We had an error:\n\nSome required fields are missing."
	requiredFieldsLog     = "Error when trying to get registration form: %v"
)

// TODO: check rendering of images
func renderError(doneMessage gtki.Label, errorMessage, logMessage string, err error) {
	log.Printf(logMessage, err)
	//doneImage.SetFromIconName("software-update-urgent", gtki.ICON_SIZE_DIALOG)
	doneMessage.SetLabel(i18n.Local(errorMessage))
}

func renderTorError(assistant gtki.Assistant, pg gtki.Widget, formMessage gtki.Label, err error) {
	log.Printf(torLogMessage, err)
	assistant.SetPageType(pg, gtki.ASSISTANT_PAGE_SUMMARY)
	formMessage.SetLabel(i18n.Local(torErrorMessage))
	//formImage.Clear()
	//formImage.SetFromIconName("software-update-urgent", gtki.ICON_SIZE_DIALOG)
}

type serverSelectionWindow struct {
	b           *builder
	assistant   gtki.Assistant
	formMessage gtki.Label
	doneMessage gtki.Label
	serverBox   gtki.ComboBoxText
	spinner     gtki.Spinner
	grid        gtki.Grid
	// formImage := builder.getObj("formImage").(gtki.Image)
	// doneImage := builder.getObj("doneImage").(gtki.Image)

	formSubmitted chan error
	done          chan error

	form *registrationForm
}

func createServerSelectionWindow() *serverSelectionWindow {
	w := &serverSelectionWindow{b: newBuilder("AccountRegistration")}

	w.b.getItems(
		"assistant", &w.assistant,
		"formMessage", &w.formMessage,
		"doneMessage", &w.doneMessage,
		"server", &w.serverBox,
		"spinner", &w.spinner,
		"formGrid", &w.grid,
	)

	w.formSubmitted = make(chan error)
	w.done = make(chan error)

	w.form = &registrationForm{grid: w.grid}

	return w
}

func (w *serverSelectionWindow) initializeServers() {
	for _, s := range servers.GetServersForRegistration() {
		w.serverBox.AppendText(s.Name)
	}
	w.serverBox.SetActive(0)
}

func (u *gtkUI) showServerSelectionWindow() {
	w := createServerSelectionWindow()
	w.assistant.SetTransientFor(u.window)
	w.initializeServers()

	w.b.ConnectSignals(map[string]interface{}{
		"on_prepare": func(_ gtki.Assistant, pg gtki.Widget) {
			switch w.assistant.GetCurrentPage() {
			case 0:
				w.serverBox.SetSensitive(true)
				w.form.server = ""

				//TODO: Destroy everything in the grid on page 1?
			case 1:
				w.serverBox.SetSensitive(false)
				w.form.server = w.serverBox.GetActiveText()

				renderFn := func(title, instructions string, fields []interface{}) error {
					w.spinner.Stop()
					w.formMessage.SetLabel("")
					w.doneMessage.SetLabel("")

					w.form.renderForm(title, fields)
					w.assistant.SetPageComplete(pg, true)

					return <-w.formSubmitted
				}

				w.spinner.Start()
				w.formMessage.SetLabel(i18n.Local("Connecting to server for registration... \n\n " +
					"This might take a while."))

				go func() {
					err := requestAndRenderRegistrationForm(w.form.server, renderFn, u.dialerFactory, u.unassociatedVerifier(), u.config)
					if err != nil && w.assistant.GetCurrentPage() != 2 {
						if err != config.ErrTorNotRunning {
							go w.assistant.SetCurrentPage(2)
						}
						w.spinner.Stop()
						renderTorError(w.assistant, pg, w.formMessage, err)
						return
					}

					w.done <- err
				}()
			case 2:
				w.formSubmitted <- w.form.accepted()
				err := <-w.done
				w.spinner.Stop()

				if err != nil {
					if err != xmpp.ErrMissingRequiredRegistrationInfo {
						renderError(w.doneMessage, contactServerError, contactServerLog, err)
					} else {
						renderError(w.doneMessage, requiredFieldsError, requiredFieldsLog, err)
					}

					return
				}

				//Save the account
				err = u.addAndSaveAccountConfig(w.form.conf)

				if err != nil {
					renderError(w.doneMessage, storeAccountInfoError, storeAccountInfoLog, err)
					return
				}

				if acc, ok := u.getAccountByID(w.form.conf.ID()); ok {
					acc.session.SetWantToBeOnline(true)
					acc.Connect()
				}

				// doneImage.SetFromIconName("emblem-default", gtki.ICON_SIZE_DIALOG)
				w.doneMessage.SetLabel(i18n.Localf("%s successfully created.", w.form.conf.Account))
			}
		},
		"on_cancel_signal": w.assistant.Destroy,
	})

	w.assistant.ShowAll()
}

func (u *gtkUI) showAddAccountWindow() {
	c, _ := config.NewAccount()

	u.accountDialog(nil, c, func() {
		u.addAndSaveAccountConfig(c)
		u.notify(i18n.Local("Account added"), i18n.Localf("%s was added successfully.", c.Account))
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
	return account.buildNotification("ConnectingAccountInfo", i18n.Localf("Connecting account\n%s", account.session.GetConfig().Account), u, nil)
}

func (account *account) buildConnectionFailureNotification(u *gtkUI, moreInfo func()) gtki.InfoBar {
	return account.buildNotification("ConnectionFailureNotification", i18n.Localf("Connection failure\n%s", account.session.GetConfig().Account), u, moreInfo)
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
