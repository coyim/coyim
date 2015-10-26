package gui

import (
	"crypto/rand"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"

	"github.com/twstrike/coyim/client"
	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/i18n"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/otr3"
)

const debugEnabled = true

type gtkUI struct {
	roster       *roster
	window       *gtk.Window
	accountsMenu *gtk.MenuItem

	config *config.Accounts

	*accountManager
}

// NewGTK returns a new client for a GTK ui
func NewGTK() client.Client {
	return &gtkUI{
		accountManager: newAccountManager(),
	}
}

func (u *gtkUI) LoadConfig(configFile string) error {
	accounts, err := config.LoadOrCreate(configFile)
	u.config = accounts

	if err != nil {
		u.Alert(err.Error())

		glib.IdleAdd(func() bool {
			u.showAddAccountWindow()
			return false
		})
	}

	u.buildAccounts(u.config)

	//TODO: replace me by observer
	for _, acc := range u.accounts {
		acc.session.SessionEventHandler = u
	}

	return nil
}

func (u *gtkUI) SaveConfig() error {
	err := u.config.Save()
	if err != nil {
		return err
	}

	u.addNewAccountsFromConfig(u.config)

	if u.window != nil {
		u.window.Emit(accountChangedSignal.String())
	}

	return nil
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

func (u *gtkUI) Info(m string) {
	fmt.Println(">>> INFO", m)
}

func (u *gtkUI) Warn(m string) {
	fmt.Println(">>> WARN", m)
}

func (u *gtkUI) Alert(m string) {
	fmt.Println(">>> ALERT", m)
}

func (u *gtkUI) Loop() {
	go u.observeAccountEvents()

	gtk.Init(&os.Args)
	u.applyStyle()
	u.mainWindow()
	gtk.Main()
}

func (u *gtkUI) Close() {}

//func (u *gtkUI) onReceiveSignal(s *glib.Signal, f func()) {
//	u.window.Connect(s.String(), f)
//}

func (u *gtkUI) initRoster() {
	u.roster = u.newRoster()
}

func (u *gtkUI) mainWindow() {
	u.window, _ = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
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
	dialog := dialog{
		title:    i18n.Local("Password"),
		position: gtk.WIN_POS_CENTER,
		id:       "dialog",
		content: []createable{
			label{i18n.Local("Password")},
			entry{
				editable:   true,
				visibility: false,
				id:         "password",
			},
			button{
				text:      i18n.Local("Connect"),
				onClicked: onPasswordDialogClicked(reg, connect),
			},
		},
	}
	dialog.create(reg)
	reg.dialogShowAll("dialog")
}

func onPasswordDialogClicked(reg *widgetRegistry, connect func(string)) func() {
	return func() {
		password := reg.getText("password")
		go connect(password)
		reg.dialogDestroy("dialog")
	}
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

func (u *gtkUI) ensureConfigHasKey(c *config.Account) {
	u.Debug(fmt.Sprintf("[%s] ensureConfigHasKey()\n", c.Account))

	if len(c.PrivateKey) == 0 {
		u.Debug(fmt.Sprintf("[%s] - No private key available. Generating...\n", c.Account))
		var priv otr3.PrivateKey

		//TODO: error
		priv.Generate(rand.Reader)
		c.PrivateKey = priv.Serialize()
		u.SaveConfig()
		u.Debug(fmt.Sprintf("[%s] - Saved\n", c.Account))
	}
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
