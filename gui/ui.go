package gui

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"

	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/i18n"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const debugEnabled = true

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
}

// UI is the user interface functionality exposed to main
type UI interface {
	Loop()
}

// NewGTK returns a new client for a GTK ui
func NewGTK() UI {
	gtk.Init(&os.Args)

	connect := make(chan *account, 0)
	disconnect := make(chan *account, 0)
	edit := make(chan *account, 0)

	res := &gtkUI{
		accountManager: newAccountManager(connect, disconnect, edit),
	}

	res.accountManager.saveConfiguration = res.SaveConfig
	res.keySupplier = config.CachingKeySupplier(res.getMasterPassword)

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
			}
		}
	}()

	return res
}

func (u *gtkUI) loadConfigInternal(configFile string) {
	accounts, ok, err := config.LoadOrCreate(configFile, u.keySupplier)
	u.config = accounts
	if ok {
		if u.viewMenu != nil {
			u.viewMenu.setFromConfig(accounts)
		}
		if err != nil {
			//TODO error
			log.Printf(err.Error())

			glib.IdleAdd(func() bool {
				u.wouldYouLikeToEncryptYourFile(func(res bool) {
					u.config.ShouldEncrypt = res
					u.showAddAccountWindow()
				})
				return false
			})
		}

		u.buildAccounts(u.config)

		//TODO: replace me by observer
		for _, acc := range u.accounts {
			acc.session.SessionEventHandler = u
		}

		if u.window != nil {
			u.window.Emit(accountChangedSignal.String())
		}
	} else {
		log.Printf("couldn't open encrypted file - either the user didn't supply a password, or the password was incorrect")
		// TODO: tell the user we couldn't open the encrypted file
	}
}

func (u *gtkUI) loadConfig(configFile string) {
	go u.loadConfigInternal(configFile)
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

func (u *gtkUI) Debug(m string) {
	if debugEnabled {
		fmt.Println(">>> DEBUG", m)
	}
}

func (u *gtkUI) Loop() {
	defer u.close()
	go u.observeAccountEvents()

	u.loadConfig(*config.ConfigFile)
	u.applyStyle()
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
	u.window, _ = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	u.displaySettings = detectCurrentDisplaySettingsFrom(&u.window.Bin.Container.Widget)
	u.initRoster()

	menubar := initMenuBar(u)
	vbox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 1)
	vbox.SetHomogeneous(false)
	vbox.PackStart(menubar, false, false, 0)
	vbox.PackStart(u.roster.widget, true, true, 0)
	u.window.Add(vbox)

	u.window.SetTitle(i18n.Local("Coy"))
	u.window.Connect("destroy", u.quit)
	u.window.SetSizeRequest(200, 600)

	u.connectShortcutsMainWindow(u.window)

	u.window.ShowAll()
}

func (u *gtkUI) quit() {
	// TODO: we should probably disconnect before quitting, if any account is connected
	gtk.MainQuit()
}

func (*gtkUI) askForPassword(connect func(string)) {
	reg := createWidgetRegistry()
	buttonID := "connect"
	dialog := dialog{
		title:    i18n.Local("Password"),
		position: gtk.WIN_POS_CENTER,
		id:       "dialog",
		content: []creatable{
			label{text: i18n.Local("Password")},
			entry{
				editable:   true,
				visibility: false,
				focused:    true,
				id:         "password",
				onActivate: onPasswordDialogClicked(reg, connect),
			},
			button{
				id:        buttonID,
				text:      i18n.Local("Connect"),
				onClicked: onPasswordDialogClicked(reg, connect),
			},
		},
	}
	dialog.createWithDefault(reg, buttonID)
	reg.dialogShowAll("dialog")
}

func onPasswordDialogClicked(reg *widgetRegistry, connect func(string)) func() {
	return func() {
		password := reg.getText("password")
		go connect(password)
		reg.dialogDestroy("dialog")
	}
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

func aboutDialog() {
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

	dialog := presenceSubscriptionDialog(accounts)
	dialog.ShowAll()
}

func (u *gtkUI) buildContactsMenu() *gtk.MenuItem {
	contactsMenu, _ := gtk.MenuItemNewWithMnemonic(i18n.Local("_Contacts"))

	submenu, _ := gtk.MenuNew()
	contactsMenu.SetSubmenu(submenu)

	menuitem, _ := gtk.MenuItemNewWithMnemonic(i18n.Local("_Add..."))
	submenu.Append(menuitem)

	menuitem.Connect("activate", u.addContactWindow)

	return contactsMenu
}

func initMenuBar(u *gtkUI) *gtk.MenuBar {
	menubar, _ := gtk.MenuBarNew()

	menubar.Append(u.buildContactsMenu())

	u.accountsMenu, _ = gtk.MenuItemNewWithMnemonic(i18n.Local("_Accounts"))
	menubar.Append(u.accountsMenu)

	//TODO: replace this by emiting the signal at startup
	u.buildAccountsMenu()
	u.window.Connect(accountChangedSignal.String(), func() {
		//TODO: should it destroy the current submenu? HOW?
		u.accountsMenu.SetSubmenu((*gtk.Widget)(nil))
		u.buildAccountsMenu()
	})

	u.createViewMenu(menubar)

	//Help -> About
	cascademenu, _ := gtk.MenuItemNewWithMnemonic(i18n.Local("_Help"))
	menubar.Append(cascademenu)
	submenu, _ := gtk.MenuNew()
	cascademenu.SetSubmenu(submenu)
	menuitem, _ := gtk.MenuItemNewWithMnemonic(i18n.Local("_About"))
	menuitem.Connect("activate", aboutDialog)
	submenu.Append(menuitem)
	return menubar
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

func (u *gtkUI) connect(account *account) {
	u.roster.connecting()
	connectFn := func(password string) {
		err := account.session.Connect(password, nil)
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
