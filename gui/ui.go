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
	"strings"
	"time"

	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/i18n"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type gtkUI struct {
	roster       *roster
	window       *gtk.Window
	accountsMenu *gtk.MenuItem
	viewMenu     *viewMenu

	config *config.Accounts

	*accountManager

	displaySettings *displaySettings

	keySupplier config.KeySupplier

	tags *tags

	toggleConnectAllAutomaticallyRequest chan bool
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
	gtk.Init(argsWithApplicationName())

	connect := make(chan *account, 0)
	disconnect := make(chan *account, 0)
	edit := make(chan *account, 0)
	toggleConnect := make(chan *account, 0)

	res := &gtkUI{
		accountManager: newAccountManager(connect, disconnect, edit, toggleConnect),
	}

	res.applyStyle()
	res.accountManager.saveConfiguration = res.SaveConfig
	res.keySupplier = config.CachingKeySupplier(res.getMasterPassword)
	res.toggleConnectAllAutomaticallyRequest = make(chan bool, 100)

	go func() {
		for {
			select {
			case acc := <-connect:
				glib.IdleAdd(func() bool {
					res.connect(acc)
					return false
				})
			case acc := <-disconnect:
				glib.IdleAdd(func() bool {
					res.disconnect(acc)
					return false
				})
			case acc := <-edit:
				glib.IdleAdd(func() bool {
					accountDialog(acc.session.CurrentAccount, res.SaveConfig)
					return false
				})
			case acc := <-toggleConnect:
				go func() {
					acc.session.CurrentAccount.ConnectAutomatically = !acc.session.CurrentAccount.ConnectAutomatically
					res.saveConfigOnly()
				}()
			}
		}
	}()

	return res
}

func (u *gtkUI) loadConfigInternal(configFile string) {
	config, ok, err := config.LoadOrCreate(configFile, u.keySupplier)

	if !ok {
		u.keySupplier.Invalidate()
		// TODO: tell the user we couldn't open the encrypted file
		log.Printf("couldn't open encrypted file - either the user didn't supply a password, or the password was incorrect")
		return
	}

	// We assign config here, AFTER the return - so that a nil config means we are in a state of incorrectness and shouldn't do stuff.
	u.config = config

	if err != nil {
		log.Printf(err.Error())

		glib.IdleAdd(func() bool {
			u.wouldYouLikeToEncryptYourFile(func(res bool) {
				u.config.ShouldEncrypt = res
				err := u.showConfigAssistant()
				if err != nil {
					log.Println(err.Error())
				}
			})
			return false
		})
	} else {
		u.configLoaded()
	}
}

func (u *gtkUI) loadConfig(configFile string) {
	//IO would block the UI loop
	go u.loadConfigInternal(configFile)
}

func (u *gtkUI) configLoaded() {
	u.buildAccounts(u.config)

	//TODO: replace me by session observer
	for _, acc := range u.accounts {
		acc.session.SessionEventHandler = u
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

func (u *gtkUI) saveConfigInternal() {
	err := u.saveConfigOnlyInternal()
	if err != nil {
		return
	}

	u.addNewAccountsFromConfig(u.config)

	if u.window != nil {
		u.window.Emit(accountChangedSignal.String())
	}
}

func (u *gtkUI) saveConfigOnlyInternal() error {
	return u.config.Save(u.keySupplier)
}

func (u *gtkUI) SaveConfig() {
	go u.saveConfigInternal()
}

func (u *gtkUI) saveConfigOnly() {
	go u.saveConfigOnlyInternal()
}

func (*gtkUI) RegisterCallback(title, instructions string, fields []interface{}) error {
	//TODO: should open a registration window
	fmt.Println("TODO")
	return nil
}

func (u *gtkUI) Loop() {
	defer u.close()
	go u.observeAccountEvents()

	u.loadConfig(*config.ConfigFile)
	u.mainWindow()
	gtk.Main()
}

func (u *gtkUI) close() {}

//func (u *gtkUI) onReceiveSignal(s *glib.Signal, f func()) {
//	u.window.Connect(s.String(), f)
//}

func (u *gtkUI) initRoster() {
	u.roster = u.newRoster()
}

func (u *gtkUI) mainWindow() {
	vars := make(map[string]string)
	vars["$title"] = i18n.Local("Coy")
	vars["$contactsMenu"] = i18n.Local("Contacts")
	vars["$addMenu"] = i18n.Local("Add...")
	vars["$accountsMenu"] = i18n.Local("Accounts")
	vars["$helpMenu"] = i18n.Local("Help")
	vars["$aboutMenu"] = i18n.Local("About")
	vars["$viewMenu"] = i18n.Local("View")
	vars["$checkItemMerge"] = i18n.Local("Merge Accounts")
	vars["$checkItemShowOffline"] = i18n.Local("Show Offline Contacts")

	builder, _ := loadBuilderWith("MainDefinition", vars)
	builder.ConnectSignals(map[string]interface{}{
		"on_close_window_signal":                    u.quit,
		"on_add_contact_window_signal":              u.addContactWindow,
		"on_about_dialog_signal":                    u.aboutDialog,
		"on_toggled_check_Item_Merge_signal":        u.toggleMergeAccounts,
		"on_toggled_check_Item_Show_Offline_signal": u.toggleShowOffline,
	})
	win, _ := builder.GetObject("mainWindow")
	u.window, _ = win.(*gtk.Window)
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

	u.connectShortcutsMainWindow(u.window)

	pl, _ := gdk.PixbufLoaderNew()
	pl.Write(decodedIcon256x256)
	pl.Close()
	pixbuf, _ := pl.GetPixbuf()
	u.window.SetIcon(pixbuf)
	gtk.WindowSetDefaultIcon(pixbuf)

	u.window.ShowAll()
}

func (u *gtkUI) quit() {
	// TODO: we should probably disconnect before quitting, if any account is connected
	gtk.MainQuit()
}

func (*gtkUI) askForPassword(connect func(string)) {
	vars := make(map[string]string)
	vars["$title"] = i18n.Local("Password")
	vars["$passwordLabel"] = i18n.Local("Password")
	vars["$saveLabel"] = i18n.Local("Connect")

	builder, _ := loadBuilderWith("AskForPasswordDefinition", vars)

	dialogObj, _ := builder.GetObject("AskForPassword")
	dialog := dialogObj.(*gtk.Dialog)

	builder.ConnectSignals(map[string]interface{}{
		"on_save_signal": func() {
			passwordObj, _ := builder.GetObject("password")
			passwordEntry := passwordObj.(*gtk.Entry)
			password, _ := passwordEntry.GetText()
			go connect(password)
			dialog.Destroy()
		},
	})

	dialog.ShowAll()
}

func (u *gtkUI) shouldViewAccounts() bool {
	return !u.config.MergeAccounts
}

func authors() []string {
	if b, err := exec.Command("git", "log").Output(); err == nil {
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
		lines = append(lines, "STRIKE Team <strike-public(AT)thoughtworks.com>")
		return lines
	}
	return []string{"STRIKE Team <strike-public@thoughtworks.com>"}
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

		return account.session.Conn.SendPresence(peer, "subscribe", "" /* generate id */)
	})

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

func (u *gtkUI) disconnect(account *account) {
	account.session.Close()
}

func (u *gtkUI) alertTorIsNotRunning() {
	builder, err := loadBuilderWith("TorNotRunningDef", nil)
	if err != nil {
		return
	}

	obj, _ := builder.GetObject("TorNotRunningDialog")
	dialog := obj.(*gtk.Dialog)

	dialog.SetTransientFor(u.window)
	dialog.ShowAll()
}

func (u *gtkUI) connect(account *account) {
	u.roster.connecting()
	connectFn := func(password string) {
		err := account.session.Connect(password, nil)

		if err == config.ErrTorNotRunning {
			glib.IdleAdd(u.alertTorIsNotRunning)
		}

		if err != nil {
			return
		}
	}

	if len(account.session.CurrentAccount.Password) == 0 {
		u.askForPassword(connectFn)
		return
	}

	go connectFn(account.session.CurrentAccount.Password)
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
	a.onConnect <- a
}

func (u *gtkUI) connectAllAutomatics(all bool) {
	log.Printf("connectAllAutomatics(%v)\n", all)
	var acc []*account
	for _, a := range u.accounts {
		if (all || a.session.CurrentAccount.ConnectAutomatically) && a.session.IsDisconnected() {
			acc = append(acc, a)
		}
	}

	glib.IdleAdd(func() {
		if len(acc) > 0 && u.roster != nil {
			u.roster.connecting()
		}
	})

	for _, a := range acc {
		go u.connectWithRandomDelay(a)
	}
}
