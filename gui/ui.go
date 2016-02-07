package gui

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/i18n"
	rosters "github.com/twstrike/coyim/roster"
	sessions "github.com/twstrike/coyim/session/access"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type gtkUI struct {
	roster           *roster
	app              *gtk.Application
	window           *gtk.ApplicationWindow
	accountsMenu     *gtk.MenuItem
	notificationArea *gtk.Box
	viewMenu         *viewMenu

	config *config.ApplicationConfig

	*accountManager

	displaySettings *displaySettings

	keySupplier config.KeySupplier

	tags *tags

	toggleConnectAllAutomaticallyRequest chan bool
	setShowAdvancedSettingsRequest       chan bool

	commands chan interface{}

	sessionFactory sessions.Factory
}

var coyimVersion string

// UI is the user interface functionality exposed to main
type UI interface {
	Loop()
}

func argsWithApplicationName() *[]string {
	newSlice := make([]string, len(os.Args))
	copy(newSlice, os.Args)
	newSlice[0] = "CoyIM"
	return &newSlice
}

// NewGTK returns a new client for a GTK ui
func NewGTK(version string, sf sessions.Factory) UI {
	coyimVersion = version
	//*.mo files should be in ./i18n/locale_code.utf8/LC_MESSAGES/
	glib.InitI18n("coy", "./i18n")
	gtk.Init(argsWithApplicationName())

	ret := &gtkUI{
		commands: make(chan interface{}, 5),
		toggleConnectAllAutomaticallyRequest: make(chan bool, 100),
		setShowAdvancedSettingsRequest:       make(chan bool, 100),
	}

	var err error
	flags := glib.APPLICATION_FLAGS_NONE
	if *config.MultiFlag {
		flags = glib.APPLICATION_NON_UNIQUE
	}
	ret.app, err = gtk.ApplicationNew("im.coy.CoyIM", flags)
	if err != nil {
		panic(err)
	}

	ret.keySupplier = config.CachingKeySupplier(ret.getMasterPassword)

	ret.accountManager = newAccountManager(ret)

	ret.sessionFactory = sf

	return ret
}

func (u *gtkUI) confirmAccountRemoval(acc *config.Account, removeAccountFunc func(*config.Account)) {
	builder := builderForDefinition("ConfirmAccountRemoval")

	obj, _ := builder.GetObject("RemoveAccount")
	dialog := obj.(*gtk.MessageDialog)
	dialog.SetTransientFor(u.window)
	dialog.SetProperty("text", acc.Account)

	response := dialog.Run()
	if gtk.ResponseType(response) == gtk.RESPONSE_YES {
		removeAccountFunc(acc)
	}

	dialog.Destroy()
}

func (u *gtkUI) initialSetupWindow() {
	u.wouldYouLikeToEncryptYourFile(func(res bool) {
		u.config.ShouldEncrypt = res
		k := func() {
			err := u.showAddAccountWindow()
			if err != nil {
				log.Println("Failed to add account:", err.Error())
			}
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
	u.buildAccounts(c, u.sessionFactory)

	doInUIThread(func() {
		if u.viewMenu != nil {
			u.viewMenu.setFromConfig(c)
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

	u.addNewAccountsFromConfig(u.config, u.sessionFactory)

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

func (*gtkUI) RegisterCallback(title, instructions string, fields []interface{}) error {
	//TODO: should open a registration window
	fmt.Println("TODO")
	return nil
}

func init() {
	runtime.LockOSThread()
}

func (u *gtkUI) Loop() {
	u.app.Connect("activate", func() {
		go u.watchCommands()
		go u.observeAccountEvents()

		applyHacks()
		u.mainWindow()
		go u.loadConfig(*config.ConfigFile)
	})

	u.app.Run([]string{})
}

func (u *gtkUI) initRoster() {
	u.roster = u.newRoster()
}

func (u *gtkUI) mainWindow() {
	builder := builderForDefinition("Main")

	builder.ConnectSignals(map[string]interface{}{
		"on_close_window_signal":                    u.quit,
		"on_add_contact_window_signal":              u.addContactWindow,
		"on_about_dialog_signal":                    u.aboutDialog,
		"on_feedback_dialog_signal":                 u.feedbackDialog,
		"on_toggled_check_Item_Merge_signal":        u.toggleMergeAccounts,
		"on_toggled_check_Item_Show_Offline_signal": u.toggleShowOffline,
	})

	win, err := builder.GetObject("mainWindow")
	if err != nil {
		panic(err)
	}

	u.window = win.(*gtk.ApplicationWindow)
	u.window.SetApplication(u.app)

	u.displaySettings = detectCurrentDisplaySettingsFrom(&u.window.Bin.Container.Widget)

	// This must happen after u.displaySettings is initialized
	// So now, roster depends on displaySettings which depends on mainWindow
	u.initRoster()

	// AccountsMenu
	am, _ := builder.GetObject("AccountsMenu")
	u.accountsMenu = am.(*gtk.MenuItem)
	// ViewMenu
	u.viewMenu = new(viewMenu)
	checkItemMerge, _ := builder.GetObject("CheckItemMerge")
	u.viewMenu.merge = checkItemMerge.(*gtk.CheckMenuItem)
	u.displaySettings.defaultSettingsOn(&u.viewMenu.merge.MenuItem.Bin.Container.Widget)

	checkItemShowOffline, _ := builder.GetObject("CheckItemShowOffline")
	u.viewMenu.offline = checkItemShowOffline.(*gtk.CheckMenuItem)
	u.displaySettings.defaultSettingsOn(&u.viewMenu.offline.MenuItem.Bin.Container.Widget)

	u.initMenuBar()
	vbox, _ := builder.GetObject("Vbox")
	vbox.(*gtk.Box).PackStart(u.roster.widget, true, true, 0)

	obj, _ := builder.GetObject("notification-area")
	u.notificationArea = obj.(*gtk.Box)

	u.config.WhenLoaded(func(a *config.ApplicationConfig) {
		if a.Display.HideFeedbackBar {
			return
		}

		doInUIThread(u.addFeedbackInfoBar)
	})

	u.connectShortcutsMainWindow(&u.window.Window)

	u.window.SetIcon(coyimIcon.getPixbuf())
	gtk.WindowSetDefaultIcon(coyimIcon.getPixbuf())

	u.window.ShowAll()
}

func (u *gtkUI) addInitialAccountsToRoster() {
	for _, account := range u.accounts {
		u.roster.update(account, rosters.New())
	}
}

func (u *gtkUI) addFeedbackInfoBar() {
	builder := builderForDefinition("FeedbackInfo")

	obj, _ := builder.GetObject("feedbackInfo")
	infobar := obj.(*gtk.InfoBar)

	u.notificationArea.PackEnd(infobar, true, true, 0)
	infobar.ShowAll()

	builder.ConnectSignals(map[string]interface{}{
		"handleResponse": func(info *gtk.InfoBar, response gtk.ResponseType) {
			if response != gtk.RESPONSE_CLOSE {
				return
			}

			infobar.Hide()
			infobar.Destroy()

			u.config.Display.HideFeedbackBar = true
			u.saveConfigOnly()
		},
	})

	obj, _ = builder.GetObject("feedbackButton")
	button := obj.(*gtk.Button)
	button.Connect("clicked", func() {
		doInUIThread(u.feedbackDialog)
	})
}

func (u *gtkUI) quit() {
	u.accountManager.disconnectAll()
	u.app.Quit()
}

func (u *gtkUI) askForPassword(accountName string, connect func(string) error) {
	dialogTemplate := "AskForPassword"

	builder := builderForDefinition(dialogTemplate)

	obj, _ := builder.GetObject(dialogTemplate)
	dialog := obj.(*gtk.Dialog)

	obj, _ = builder.GetObject("accountName")
	label := obj.(*gtk.Label)
	label.SetText(accountName)
	label.SetSelectable(true)

	builder.ConnectSignals(map[string]interface{}{
		"on_save_signal": func() {
			passwordObj, _ := builder.GetObject("password")
			passwordEntry := passwordObj.(*gtk.Entry)
			password, _ := passwordEntry.GetText()

			if len(password) > 0 {
				go connect(password)
				dialog.Destroy()
			}
		},
		"on_cancel_signal": dialog.Destroy,
	})

	dialog.SetTransientFor(u.window)
	dialog.ShowAll()
}

func (u *gtkUI) feedbackDialog() {
	builder := builderForDefinition("Feedback")

	obj, _ := builder.GetObject("dialog")
	dialog := obj.(*gtk.MessageDialog)
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
	}
}

func (u *gtkUI) aboutDialog() {
	dialog, _ := gtk.AboutDialogNew()
	dialog.SetName(i18n.Local("Coy IM!"))
	dialog.SetProgramName("CoyIM")
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

	dialog := presenceSubscriptionDialog(accounts, func(accountID, peer string) error {
		//TODO errors
		account, ok := u.roster.getAccount(accountID)
		if !ok {
			return fmt.Errorf(i18n.Local("There is no account with the id %q"), accountID)
		}

		if !account.connected() {
			return errors.New(i18n.Local("Can't send a contact request from an offline account"))
		}

		return account.session.RequestPresenceSubscription(peer)
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
		u.buildAccountsMenu()
		u.accountsMenu.ShowAll()
		u.rosterUpdated()
	})

	u.buildAccountsMenu()
	u.accountsMenu.ShowAll()
}

func (u *gtkUI) rosterUpdated() {
	doInUIThread(u.roster.redraw)
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
