package gui

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gdki"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glibi"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtki"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/pangoi"
	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/gui/settings"
	"github.com/twstrike/coyim/i18n"
	rosters "github.com/twstrike/coyim/roster"
	sessions "github.com/twstrike/coyim/session/access"
	"github.com/twstrike/coyim/xmpp/interfaces"
)

const (
	programName        = "CoyIM"
	applicationID      = "im.coy.CoyIM"
	localizationDomain = "coy"
)

type gtkUI struct {
	roster           *roster
	app              gtki.Application
	window           gtki.ApplicationWindow
	accountsMenu     gtki.MenuItem
	notificationArea gtki.Box
	viewMenu         *viewMenu
	optionsMenu      *optionsMenu

	unified *unifiedLayout

	config *config.ApplicationConfig

	*accountManager

	displaySettings *displaySettings

	keySupplier config.KeySupplier

	tags *tags

	toggleConnectAllAutomaticallyRequest chan bool
	setShowAdvancedSettingsRequest       chan bool

	commands chan interface{}

	sessionFactory sessions.Factory

	dialerFactory interfaces.DialerFactory

	settings *settings.Settings
}

// Graphics represent the graphic configuration
type Graphics struct {
	gtk   gtki.Gtk
	glib  glibi.Glib
	gdk   gdki.Gdk
	pango pangoi.Pango
}

// CreateGraphics creates a Graphic represention from the given arguments
func CreateGraphics(gtkVal gtki.Gtk, glibVal glibi.Glib, gdkVal gdki.Gdk, pangoVal pangoi.Pango) Graphics {
	return Graphics{
		gtk:   gtkVal,
		glib:  glibVal,
		gdk:   gdkVal,
		pango: pangoVal,
	}
}

var g Graphics

var coyimVersion string

// UI is the user interface functionality exposed to main
type UI interface {
	Loop()
}

func argsWithApplicationName() *[]string {
	newSlice := make([]string, len(os.Args))
	copy(newSlice, os.Args)
	newSlice[0] = programName
	return &newSlice
}

// NewGTK returns a new client for a GTK ui
func NewGTK(version string, sf sessions.Factory, df interfaces.DialerFactory, gx Graphics) UI {
	runtime.LockOSThread()
	coyimVersion = version
	g = gx
	initSignals()

	//*.mo files should be in ./i18n/locale_code.utf8/LC_MESSAGES/
	g.glib.InitI18n(localizationDomain, "./i18n")
	g.gtk.Init(argsWithApplicationName())

	ret := &gtkUI{
		commands: make(chan interface{}, 5),
		toggleConnectAllAutomaticallyRequest: make(chan bool, 100),
		setShowAdvancedSettingsRequest:       make(chan bool, 100),
		dialerFactory:                        df,
	}

	var err error
	flags := glibi.APPLICATION_FLAGS_NONE
	if *config.MultiFlag {
		flags = glibi.APPLICATION_NON_UNIQUE
	}
	ret.app, err = g.gtk.ApplicationNew(applicationID, flags)
	if err != nil {
		panic(err)
	}

	ret.keySupplier = config.CachingKeySupplier(ret.getMasterPassword)

	ret.accountManager = newAccountManager(ret)

	ret.sessionFactory = sf

	ret.settings = settings.For("")

	return ret
}

func (u *gtkUI) confirmAccountRemoval(acc *config.Account, removeAccountFunc func(*config.Account)) {
	builder := newBuilder("ConfirmAccountRemoval")

	obj := builder.getObj("RemoveAccount")
	dialog := obj.(gtki.MessageDialog)
	dialog.SetTransientFor(u.window)
	dialog.SetProperty("text", acc.Account)

	response := dialog.Run()
	if gtki.ResponseType(response) == gtki.RESPONSE_YES {
		removeAccountFunc(acc)
	}

	dialog.Destroy()
}

func (u *gtkUI) initialSetupWindow() {
	u.wouldYouLikeToEncryptYourFile(func(res bool) {
		u.config.SetShouldSaveFileEncrypted(res)
		k := func() {
			go u.showFirstAccountWindow()
		}
		if res {
			u.captureInitialMasterPassword(k)
		} else {
			k()
		}
	})
}

func (u *gtkUI) loadConfig(configFile string) {
	u.config.WhenLoaded(func(c *config.ApplicationConfig) {
		u.configLoaded(c)
	})

	ok := false
	var conf *config.ApplicationConfig
	var err error
	for !ok {
		conf, ok, err = config.LoadOrCreate(configFile, u.keySupplier)
		if !ok {
			log.Printf("couldn't open encrypted file - either the user didn't supply a password, or the password was incorrect: %v", err)
			u.keySupplier.Invalidate()
			u.keySupplier.LastAttemptFailed()
		}
	}

	// We assign config here, AFTER the return - so that a nil config means we are in a state of incorrectness and shouldn't do stuff.
	// We never check, since a panic here is a serious programming error
	u.config = conf

	if err != nil {
		log.Printf(err.Error())
		doInUIThread(u.initialSetupWindow)
		return
	}
	if u.config.UpdateToLatestVersion() {
		u.saveConfigOnlyInternal()
	}
}

func (u *gtkUI) configLoaded(c *config.ApplicationConfig) {
	u.settings = settings.For(c.GetUniqueID())
	u.roster.deNotify.updateWith(u.settings)
	if !u.settings.GetSingleWindow() {
		u.unified = nil
	}

	u.buildAccounts(c, u.sessionFactory, u.dialerFactory)

	doInUIThread(func() {
		if u.viewMenu != nil {
			u.viewMenu.setFromConfig(c)
		}

		if u.optionsMenu != nil {
			u.optionsMenu.setFromConfig(c)
		}

		if u.window != nil {
			u.window.Emit(accountChangedSignal.String())
		}
	})

	u.addInitialAccountsToRoster()

	if c.ConnectAutomatically {
		u.connectAllAutomatics(false)
	}

	go u.listenToToggleConnectAllAutomatically()
	go u.listenToSetShowAdvancedSettings()
}

func (u *gtkUI) saveConfigInternal() error {
	err := u.saveConfigOnlyInternal()
	if err != nil {
		return err
	}

	u.addNewAccountsFromConfig(u.config, u.sessionFactory, u.dialerFactory)

	if u.window != nil {
		u.window.Emit(accountChangedSignal.String())
	}

	return nil
}

func (u *gtkUI) saveConfigOnlyInternal() error {
	return u.config.Save(u.keySupplier)
}

func (u *gtkUI) SaveConfig() {
	go func() {
		err := u.saveConfigInternal()
		if err != nil {
			log.Println("Failed to save config file:", err.Error())
		}
	}()
}

func (u *gtkUI) removeSaveReload(acc *config.Account) {
	//TODO: the account configs should be managed by the account manager
	u.accountManager.removeAccount(acc, func() {
		u.config.Remove(acc)
		u.SaveConfig()
	})
}

func (u *gtkUI) saveConfigOnly() {
	go func() {
		err := u.saveConfigOnlyInternal()
		if err != nil {
			log.Println("Failed to save config file:", err.Error())
		}
	}()
}

func (u *gtkUI) onActivate() {
	if activeWindow := u.app.GetActiveWindow(); activeWindow != nil {
		activeWindow.Present()
		return
	}

	applyHacks()
	u.mainWindow()

	go u.watchCommands()
	go u.observeAccountEvents()
	go u.loadConfig(*config.ConfigFile)
}

func (u *gtkUI) Loop() {
	u.app.Connect("activate", u.onActivate)
	u.app.Run([]string{})
}

func (u *gtkUI) initRoster() {
	u.roster = u.newRoster()
}

func (u *gtkUI) mainWindow() {
	builder := newBuilder("Main")

	builder.ConnectSignals(map[string]interface{}{
		"on_close_window_signal":                       u.quit,
		"on_add_contact_window_signal":                 u.addContactWindow,
		"on_about_dialog_signal":                       u.aboutDialog,
		"on_feedback_dialog_signal":                    u.feedbackDialog,
		"on_toggled_check_Item_Merge_signal":           u.toggleMergeAccounts,
		"on_toggled_check_Item_Show_Offline_signal":    u.toggleShowOffline,
		"on_toggled_encrypt_configuration_file_signal": u.toggleEncryptedConfig,
		"on_preferences_signal":                        u.showGlobalPreferences,
	})

	win, err := builder.GetObject("mainWindow")
	if err != nil {
		panic(err)
	}

	u.window = win.(gtki.ApplicationWindow)
	u.window.SetApplication(u.app)

	u.displaySettings = detectCurrentDisplaySettingsFrom(u.window)

	// This must happen after u.displaySettings is initialized
	// So now, roster depends on displaySettings which depends on mainWindow
	u.initRoster()

	g.gtk.LabelNew("hello")

	// AccountsMenu
	u.accountsMenu = builder.getObj("AccountsMenu").(gtki.MenuItem)

	// ViewMenu
	u.viewMenu = new(viewMenu)
	u.viewMenu.merge = builder.getObj("CheckItemMerge").(gtki.CheckMenuItem)
	u.displaySettings.defaultSettingsOn(u.viewMenu.merge)

	u.viewMenu.offline = builder.getObj("CheckItemShowOffline").(gtki.CheckMenuItem)
	u.displaySettings.defaultSettingsOn(u.viewMenu.offline)

	// OptionsMenu
	u.optionsMenu = new(optionsMenu)
	u.optionsMenu.encryptConfig = builder.getObj("EncryptConfigurationFileCheckMenuItem").(gtki.CheckMenuItem)
	u.displaySettings.defaultSettingsOn(u.optionsMenu.encryptConfig)

	u.initMenuBar()
	obj := builder.getObj("Vbox")
	vbox := obj.(gtki.Box)
	vbox.PackStart(u.roster.widget, true, true, 0)

	obj = builder.getObj("Hbox")
	hbox := obj.(gtki.Box)
	u.unified = newUnifiedLayout(u, vbox, hbox)

	u.notificationArea = builder.getObj("notification-area").(gtki.Box)

	u.config.WhenLoaded(func(a *config.ApplicationConfig) {
		if a.Display.HideFeedbackBar {
			return
		}

		doInUIThread(u.addFeedbackInfoBar)
	})

	u.connectShortcutsMainWindow(u.window)

	u.window.SetIcon(coyimIcon.getPixbuf())
	g.gtk.WindowSetDefaultIcon(coyimIcon.getPixbuf())

	u.window.ShowAll()
}

func (u *gtkUI) addInitialAccountsToRoster() {
	for _, account := range u.accounts {
		u.roster.update(account, rosters.New())
	}
}

func (u *gtkUI) addFeedbackInfoBar() {
	builder := newBuilder("FeedbackInfo")

	obj := builder.getObj("feedbackInfo")
	infobar := obj.(gtki.InfoBar)

	u.notificationArea.PackEnd(infobar, true, true, 0)
	infobar.ShowAll()

	builder.ConnectSignals(map[string]interface{}{
		"handleResponse": func(info gtki.InfoBar, response gtki.ResponseType) {
			if response != gtki.RESPONSE_CLOSE {
				return
			}

			infobar.Hide()
			infobar.Destroy()

			u.config.Display.HideFeedbackBar = true
			u.saveConfigOnly()
		},
	})

	obj = builder.getObj("feedbackButton")
	button := obj.(gtki.Button)
	button.Connect("clicked", func() {
		doInUIThread(u.feedbackDialog)
	})
}

func (u *gtkUI) quit() {
	u.accountManager.disconnectAll()
	u.app.Quit()
}

func (u *gtkUI) askForPassword(accountName string, cancel func(), connect func(string) error) {
	dialogTemplate := "AskForPassword"

	builder := newBuilder(dialogTemplate)

	obj := builder.getObj(dialogTemplate)
	dialog := obj.(gtki.Dialog)

	obj = builder.getObj("accountName")
	label := obj.(gtki.Label)
	label.SetText(accountName)
	label.SetSelectable(true)

	builder.ConnectSignals(map[string]interface{}{
		"on_entered_password_signal": func() {
			passwordObj := builder.getObj("password")
			passwordEntry := passwordObj.(gtki.Entry)
			password, _ := passwordEntry.GetText()

			if len(password) > 0 {
				go connect(password)
				dialog.Destroy()
			}
		},
		"on_cancel_password_signal": func() {
			cancel()
			dialog.Destroy()
		},
	})

	dialog.SetTransientFor(u.window)
	dialog.ShowAll()
}

func (u *gtkUI) feedbackDialog() {
	builder := newBuilder("Feedback")

	obj := builder.getObj("dialog")
	dialog := obj.(gtki.MessageDialog)
	dialog.SetTransientFor(u.window)

	dialog.Run()
	dialog.Destroy()
}

func (u *gtkUI) shouldViewAccounts() bool {
	return !u.config.Display.MergeAccounts
}

func authors() []string {
	return []string{
		"Fan Jiang  -  fan.torchz@gmail.com",
		"Iván Pazmiño  -  iapazmino@gmail.com",
		"Ola Bini  -  ola@olabini.se",
		"Reinaldo de Souza Jr  -  juniorz@gmail.com",
		"Tania Silva  -  tsilva@thoughtworks.com",
		"Adam Langley",
		"Gray Leonard - gl7039a@american.edu",
		"Bruce Leidl - bruce@subgraph.com",
		"xSmurf - matth@subgraph.com",
	}
}

func (u *gtkUI) aboutDialog() {
	dialog, _ := g.gtk.AboutDialogNew()
	dialog.SetName(i18n.Local("Coy IM!"))
	dialog.SetProgramName(programName)
	dialog.SetAuthors(authors())
	dialog.SetVersion(coyimVersion)
	dialog.SetLicense(`GNU GENERAL PUBLIC LICENSE, Version 3`)
	dialog.SetWrapLicense(true)

	dialog.SetTransientFor(u.window)
	dialog.Run()
	dialog.Destroy()
}

func (u *gtkUI) addContactWindow() {
	accounts := make([]*account, 0, len(u.accounts))

	for i := range u.accounts {
		acc := u.accounts[i]
		if acc.connected() {
			accounts = append(accounts, acc)
		}
	}

	dialog := presenceSubscriptionDialog(accounts, func(accountID, peer, msg, nick string, autoAuth bool) error {
		account, ok := u.roster.getAccount(accountID)
		if !ok {
			return fmt.Errorf(i18n.Local("There is no account with the id %q"), accountID)
		}

		if !account.connected() {
			return errors.New(i18n.Local("Can't send a contact request from an offline account"))
		}

		err := account.session.RequestPresenceSubscription(peer, msg)

		if nick != "" {
			account.session.GetConfig().SavePeerDetails(peer, nick, []string{})
			u.SaveConfig()
		}

		if autoAuth {
			account.session.AutoApprove(peer)
		}

		return err
	})

	dialog.SetTransientFor(u.window)
	dialog.ShowAll()
}

func (u *gtkUI) listenToToggleConnectAllAutomatically() {
	for {
		val := <-u.toggleConnectAllAutomaticallyRequest
		u.config.ConnectAutomatically = val
		u.saveConfigOnly()
	}
}

func (u *gtkUI) setConnectAllAutomatically(val bool) {
	u.toggleConnectAllAutomaticallyRequest <- val
}

func (u *gtkUI) setShowAdvancedSettings(val bool) {
	u.setShowAdvancedSettingsRequest <- val
}

func (u *gtkUI) listenToSetShowAdvancedSettings() {
	for {
		val := <-u.setShowAdvancedSettingsRequest
		u.config.AdvancedOptions = val
		u.saveConfigOnly()
	}
}

func (u *gtkUI) initMenuBar() {
	u.window.Connect(accountChangedSignal.String(), func() {
		doInUIThread(func() {
			u.buildAccountsMenu()
			u.accountsMenu.ShowAll()
			u.rosterUpdated()
		})
	})

	u.buildAccountsMenu()
	u.accountsMenu.ShowAll()
}

func (u *gtkUI) rosterUpdated() {
	doInUIThread(u.roster.redraw)
	if u.unified != nil {
		doInUIThread(u.unified.update)
	}
}

func (u *gtkUI) connectionInfo(account *account) {
	u.connectionInfoDialog(account)
}

func (u *gtkUI) editAccount(account *account) {
	u.accountDialog(account.session, account.session.GetConfig(), func() {
		u.SaveConfig()
		account.session.ReloadKeys()
	})
}

func (u *gtkUI) removeAccount(account *account) {
	u.confirmAccountRemoval(account.session.GetConfig(), func(c *config.Account) {
		account.session.SetWantToBeOnline(false)
		account.disconnect()
		u.removeSaveReload(c)
	})
}

func (u *gtkUI) toggleAutoConnectAccount(account *account) {
	account.session.GetConfig().ToggleConnectAutomatically()
	u.saveConfigOnly()
}

func (u *gtkUI) toggleAlwaysEncryptAccount(account *account) {
	account.session.GetConfig().ToggleAlwaysEncrypt()
	u.saveConfigOnly()
}
