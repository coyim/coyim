package gui

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/i18n"
	xroster "github.com/twstrike/coyim/roster"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type gtkUI struct {
	roster           *roster
	window           *gtk.Window
	accountsMenu     *gtk.MenuItem
	notificationArea *gtk.Box
	viewMenu         *viewMenu

	config *config.ApplicationConfig

	*accountManager

	displaySettings *displaySettings

	keySupplier config.KeySupplier

	tags *tags

	toggleConnectAllAutomaticallyRequest chan bool

	commands chan interface{}
}

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
func NewGTK() UI {
	//*.mo files should be in ./i18n/locale_code.utf8/LC_MESSAGES/
	glib.InitI18n("coy", "./i18n")
	gtk.Init(argsWithApplicationName())

	ret := &gtkUI{
		commands: make(chan interface{}, 5),
		toggleConnectAllAutomaticallyRequest: make(chan bool, 100),
	}

	ret.applyStyle()
	ret.keySupplier = config.CachingKeySupplier(ret.getMasterPassword)

	ret.accountManager = newAccountManager(ret)

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
		err := u.showAddAccountWindow()
		if err != nil {
			log.Println("Failed to add account:", err.Error())
		}
	})
}

func (u *gtkUI) loadConfig(configFile string) {
	config, ok, err := config.LoadOrCreate(configFile, u.keySupplier)

	if !ok {
		u.keySupplier.Invalidate()
		u.loadConfig(configFile)
		// TODO: tell the user we couldn't open the encrypted file
		log.Printf("couldn't open encrypted file - either the user didn't supply a password, or the password was incorrect")
		return
	}

	// We assign config here, AFTER the return - so that a nil config means we are in a state of incorrectness and shouldn't do stuff.
	u.config = config

	if err != nil {
		log.Printf(err.Error())
		glib.IdleAdd(u.initialSetupWindow)
		return
	}

	u.configLoaded()
}

func (u *gtkUI) configLoaded() {
	u.buildAccounts(u.config)

	//TODO: replace me by session observer
	for _, acc := range u.accounts {
		acc.session.SessionEventHandler = u
		u.roster.update(acc, xroster.New())
	}

	glib.IdleAdd(func() bool {
		if u.viewMenu != nil {
			u.viewMenu.setFromConfig(u.config)
		}

		if u.window != nil {
			u.window.Emit(accountChangedSignal.String())
		}

		return false
	})

	if u.config.ConnectAutomatically {
		u.connectAllAutomatics(false)
	}

	go u.listenToToggleConnectAllAutomatically()
}

func (u *gtkUI) saveConfigInternal() error {
	err := u.saveConfigOnlyInternal()
	if err != nil {
		return err
	}

	u.addNewAccountsFromConfig(u.config)

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
	u.config.Remove(acc)
	u.SaveConfig()
	u.accountManager.buildAccounts(u.config)
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

func (u *gtkUI) Loop() {
	go u.watchCommands()
	go u.observeAccountEvents()
	go u.loadConfig(*config.ConfigFile)

	u.mainWindow()
	gtk.Main()
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
		"on_toggled_check_Item_Merge_signal":        u.toggleMergeAccounts,
		"on_toggled_check_Item_Show_Offline_signal": u.toggleShowOffline,
	})

	win, err := builder.GetObject("mainWindow")
	if err != nil {
		panic(err)
	}

	u.window = win.(*gtk.Window)

	u.displaySettings = detectCurrentDisplaySettingsFrom(&u.window.Bin.Container.Widget)
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
	u.addFeedbackInfoBar()

	u.connectShortcutsMainWindow(u.window)

	u.window.SetIcon(coyimIcon.getPixbuf())
	gtk.WindowSetDefaultIcon(coyimIcon.getPixbuf())

	u.window.ShowAll()
}

func (u *gtkUI) addFeedbackInfoBar() {
	builder := builderForDefinition("FeedbackInfo")

	obj, _ := builder.GetObject("feedbackInfo")
	infobar := obj.(*gtk.InfoBar)
	u.notificationArea.PackEnd(infobar, true, true, 0)

	obj, _ = builder.GetObject("feedbackButton")
	button := obj.(*gtk.Button)
	button.Connect("clicked", func() {
		glib.IdleAdd(u.feedbackDialog)
	})
}

func (u *gtkUI) quit() {
	// TODO: we should probably disconnect before quitting, if any account is connected
	gtk.MainQuit()
}

func (u *gtkUI) askForPassword(accountName string, connect func(string) error) {
	dialogTemplate := "AskForPassword"

	builder := builderForDefinition(dialogTemplate)

	obj, _ := builder.GetObject(dialogTemplate)
	dialog := obj.(*gtk.Dialog)

	obj, _ = builder.GetObject("accountName")
	label := obj.(*gtk.Label)
	label.SetText(accountName)

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

	response := dialog.Run()
	if gtk.ResponseType(response) == gtk.RESPONSE_CLOSE {
		dialog.Destroy()
	}
}

func (u *gtkUI) shouldViewAccounts() bool {
	return !u.config.MergeAccounts
}

func authors() []string {
	strikeMessage := "STRIKE Team <coyim@thoughtworks.com>"

	b, err := exec.Command("git", "log").Output()
	if err != nil {
		return []string{strikeMessage}
	}

	lines := strings.Split(string(b), "\n")

	var a []string
	r := regexp.MustCompile(`^Author:\s*([^ <]+).*$`)
	for _, e := range lines {
		ms := r.FindStringSubmatch(e)
		if ms == nil {
			continue
		}
		a = append(a, ms[1])
	}
	sort.Strings(a)
	var p string
	lines = []string{}
	for _, e := range a {
		if p == e {
			continue
		}
		lines = append(lines, e)
		p = e
	}
	lines = append(lines, strikeMessage)
	return lines
}

func (u gtkUI) aboutDialog() {
	dialog, _ := gtk.AboutDialogNew()
	dialog.SetName(i18n.Local("Coy IM!"))
	dialog.SetProgramName("Coyim")
	dialog.SetAuthors(authors())
	// dir, _ := path.Split(os.Args[0])
	// imagefile := path.Join(dir, "../../data/coyim-logo.png")
	// pixbuf, _ := gdkpixbuf.NewFromFile(imagefile)
	// dialog.SetLogo(pixbuf)
	dialog.SetLicense(`Copyright (c) 2012 The Go Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.`)
	dialog.SetWrapLicense(true)
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
			return fmt.Errorf("There is no account with the id %q", accountID)
		}

		if !account.connected() {
			return errors.New("Cant send a contact request from an offline account")
		}

		return account.session.RequestPresenceSubscription(peer)
	})

	dialog.SetTransientFor(u.window)
	dialog.ShowAll()
}

func (u *gtkUI) listenToToggleConnectAllAutomatically() {
	for {
		<-u.toggleConnectAllAutomaticallyRequest
		u.config.ConnectAutomatically = !u.config.ConnectAutomatically
		u.saveConfigOnly()
	}
}

func (u *gtkUI) toggleConnectAllAutomatically() {
	u.toggleConnectAllAutomaticallyRequest <- true
}

func (u *gtkUI) initMenuBar() {
	u.window.Connect(accountChangedSignal.String(), func() {
		u.buildAccountsMenu()
		u.accountsMenu.ShowAll()
	})
}

func (u *gtkUI) rosterUpdated() {
	glib.IdleAdd(func() bool {
		u.roster.redraw()
		return false
	})
}

func (u *gtkUI) alertTorIsNotRunning() {
	//TODO: should it notify instead of alert?

	builder := builderForDefinition("TorNotRunning")

	obj, _ := builder.GetObject("TorNotRunningDialog")
	dialog := obj.(*gtk.Dialog)

	dialog.SetTransientFor(u.window)
	dialog.ShowAll()
}

func (u *gtkUI) askForServerDetails(conf *config.Account, connectFn func() error) {
	builder := builderForDefinition("ConnectionSettings")

	obj, _ := builder.GetObject("ConnectionSettingsDialog")
	dialog := obj.(*gtk.Dialog)

	obj, _ = builder.GetObject("server")
	serverEntry := obj.(*gtk.Entry)

	obj, _ = builder.GetObject("port")
	portEntry := obj.(*gtk.Entry)

	if conf.Port == 0 {
		conf.Port = 5222
	}

	serverEntry.SetText(conf.Server)
	portEntry.SetText(strconv.Itoa(conf.Port))

	builder.ConnectSignals(map[string]interface{}{
		"reconnect": func() {
			defer dialog.Destroy()

			//TODO: validate
			conf.Server, _ = serverEntry.GetText()

			p, _ := portEntry.GetText()
			conf.Port, _ = strconv.Atoi(p)

			go func() {
				if connectFn() != nil {
					return
				}

				u.saveConfigOnly()
			}()
		},
	})

	dialog.SetTransientFor(u.window)
	dialog.ShowAll()
}

func (u *gtkUI) disconnectAccount(account *account) {
	go account.session.Close()
}

func (u *gtkUI) editAccount(account *account) {
	u.accountDialog(account.session.CurrentAccount, u.SaveConfig)
}

func (u *gtkUI) removeAccount(account *account) {
	u.confirmAccountRemoval(account.session.CurrentAccount, func(c *config.Account) {
		u.disconnectAccount(account)
		u.removeSaveReload(c)
	})
}

func (u *gtkUI) toggleAutoConnectAccount(account *account) {
	account.session.CurrentAccount.ToggleConnectAutomatically()
	u.saveConfigOnly()
}

// implemented using Sattolo’s variant of the Fisher–Yates shuffle
func shuffleAccounts(a []*account) {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}

func (u *gtkUI) connectWithRandomDelay(a *account) {
	sleepDelay := time.Duration(rand.Int31n(7643)) * time.Millisecond
	log.Printf("connectWithRandomDelay(%v, %vms)\n", a.session.CurrentAccount.Account, sleepDelay)
	time.Sleep(sleepDelay)
	a.connect()
}

func (u *gtkUI) connectAllAutomatics(all bool) {
	log.Printf("connectAllAutomatics(%v)\n", all)
	var acc []*account
	for _, a := range u.accounts {
		if (all || a.session.CurrentAccount.ConnectAutomatically) && a.session.IsDisconnected() {
			acc = append(acc, a)
		}
	}

	//TODO: add notification?

	for _, a := range acc {
		go u.connectWithRandomDelay(a)
	}
}
