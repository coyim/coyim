package gui

import (
	"crypto/rand"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"

	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/coyim/session"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/otr3"
)

const debugEnabled = true

type gtkUI struct {
	roster       *Roster
	window       *gtk.Window
	accountsMenu *gtk.MenuItem

	configFileManager *config.ConfigFileManager
	multiConfig       *config.MultiAccountConfig

	accounts []*Account
}

func NewGTK() *gtkUI {
	return &gtkUI{}
}

func (u *gtkUI) LoadConfig(configFile string) error {
	u.configFileManager = config.NewConfigFileManager(configFile)

	err := u.configFileManager.ParseConfigFile()
	if err != nil {
		u.Alert(err.Error())

		u.configFileManager.MultiAccountConfig = &config.MultiAccountConfig{}

		//TODO: Remove this
		u.multiConfig = u.configFileManager.MultiAccountConfig

		glib.IdleAdd(func() bool {
			u.showAddAccountWindow()
			return false
		})

		return nil
	}

	//TODO: REMOVE this
	u.multiConfig = u.configFileManager.MultiAccountConfig

	u.accounts = BuildAccountsFrom(u.multiConfig, u.configFileManager, u)

	return nil
}

func (u *gtkUI) addNewAccountsFromConfig() {
	for _, configAccount := range u.configFileManager.Accounts {
		var found bool
		for _, acc := range u.accounts {
			if acc.Id() == configAccount.Id() {
				found = true
				break
			}
		}

		if found {
			continue
		}

		acc := newAccount(&configAccount, u.configFileManager)
		u.accounts = append(u.accounts, acc)
	}
}

func (u *gtkUI) SaveConfig() error {
	err := u.configFileManager.Save()
	if err != nil {
		return err
	}

	u.addNewAccountsFromConfig()

	if u.window != nil {
		u.window.Emit(AccountChangedSignal.String())
	}

	return nil
}

//TODO: Should it be per session?
func (u *gtkUI) Disconnected() {
	//TODO: Is it necessary?
}

func (*gtkUI) RegisterCallback(title, instructions string, fields []interface{}) error {
	//TODO: should open a registration window
	fmt.Println("TODO")
	return nil
}

func (u *gtkUI) findAccountForSession(s *session.Session) *Account {
	a, ok := s.Account.(*Account)
	if ok {
		return a
	}
	return nil
}

func (u *gtkUI) MessageReceived(s *session.Session, from, timestamp string, encrypted bool, message []byte) {
	account := u.findAccountForSession(s)
	if account == nil {
		//TODO error
		return
	}

	u.roster.MessageReceived(account, from, timestamp, encrypted, message)
}

func (u *gtkUI) NewOTRKeys(uid string, conversation *otr3.Conversation) {
	u.Info(fmt.Sprintf("TODO: notify new keys from %s", uid))
}

func (u *gtkUI) OTREnded(uid string) {
	//TODO: conversation ended
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
	gtk.Init(&os.Args)
	u.mainWindow()
	gtk.Main()
}

func (u *gtkUI) Close() {}

//func (u *gtkUI) onReceiveSignal(s *glib.Signal, f func()) {
//	u.window.Connect(s.String(), f)
//}

func (u *gtkUI) initRoster() {
	u.roster = NewRoster()
}

func (u *gtkUI) mainWindow() {
	u.window, _ = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	u.applyStyle()
	u.initRoster()

	menubar := initMenuBar(u)
	vbox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 1)
	vbox.SetHomogeneous(false)
	vbox.PackStart(menubar, false, false, 0)
	vbox.PackStart(u.roster.Window, true, true, 0)
	u.window.Add(vbox)

	u.window.SetTitle(i18n.Local("Coy"))
	u.window.Connect("destroy", u.quit)
	u.window.SetSizeRequest(200, 600)

	u.connectShortcuts()

	u.window.ShowAll()
}

func (u *gtkUI) quit() {
	// TODO: we should probably disconnect before quitting, if any account is connected
	gtk.MainQuit()
}

func (u *gtkUI) connectShortcuts() {
	// TODO: once we accelerator support
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
	accounts := make([]*Account, 0, len(u.accounts))

	for i := range u.accounts {
		acc := u.accounts[i]
		if acc.Connected() {
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
	u.window.Connect(AccountChangedSignal.String(), func() {
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

func (u *gtkUI) SubscriptionRequest(s *session.Session, from string) {
	confirmDialog := authorizePresenceSubscriptionDialog(u.window, from)

	glib.IdleAdd(func() bool {
		confirm := gtk.ResponseType(confirmDialog.Run()) == gtk.RESPONSE_YES
		confirmDialog.Destroy()

		s.HandleConfirmOrDeny(from, confirm)

		return false
	})
}

func (u *gtkUI) rosterUpdated() {
	glib.IdleAdd(func() bool {
		u.roster.Redraw()
		return false
	})
}

func (u *gtkUI) ProcessPresence(from, to, show, showStatus string, gone bool) {
	u.Debug(fmt.Sprintf("[%s] Presence from %s: show: %s status: %s gone: %v\n", to, from, show, showStatus, gone))
	//	u.roster.debugPrintRosterFor(to)
	u.rosterUpdated()
}

func (u *gtkUI) Subscribed(account, peer string) {
	u.Debug(fmt.Sprintf("[%s] Subscribed to %s\n", account, peer))
	u.rosterUpdated()
}

func (u *gtkUI) Unsubscribe(account, peer string) {
	u.Debug(fmt.Sprintf("[%s] Unsubscribed from %s\n", account, peer))
	u.rosterUpdated()
}

func (u *gtkUI) IQReceived(iq string) {
	u.Debug(fmt.Sprintf("received iq: %v\n", iq))
	//TODO
}

func (u *gtkUI) RosterReceived(s *session.Session) {
	account := u.findAccountForSession(s)
	if account == nil {
		//TODO error
		return
	}

	u.roster.Update(account, s.R)

	glib.IdleAdd(func() bool {
		u.roster.Redraw()
		return false
	})
}

func (u *gtkUI) disconnect(account *Account) {
	account.Session.Close()
	u.window.Emit(account.DisconnectedSignal.String())
}

func (u *gtkUI) ensureConfigHasKey(c *config.Config) {
	u.Debug(fmt.Sprintf("[%s] ensureConfigHasKey()\n", c.Account))

	if len(c.PrivateKey) == 0 {
		u.Debug(fmt.Sprintf("[%s] - No private key available. Generating...\n", c.Account))
		var priv otr3.PrivateKey
		priv.Generate(rand.Reader)
		c.PrivateKey = priv.Serialize()
		u.SaveConfig()
		u.Debug(fmt.Sprintf("[%s] - Saved\n", c.Account))
	}
}

func (u *gtkUI) connect(account *Account) {
	connectFn := func(password string) {
		err := account.Session.Connect(password, nil)
		if err != nil {
			u.window.Emit(account.DisconnectedSignal.String())
			return
		}

		u.window.Emit(account.ConnectedSignal.String())
	}

	if len(account.Password) == 0 {
		u.askForPassword(connectFn)
		return
	}

	go connectFn(account.Password)
}
