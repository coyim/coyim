package gui

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"

	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/coyim/session"
	"github.com/twstrike/coyim/xmpp"

	"github.com/twstrike/go-gtk/glib"
	"github.com/twstrike/go-gtk/gtk"
	"github.com/twstrike/otr3"
)

type gtkUI struct {
	roster       *Roster
	window       *gtk.Window
	accountsMenu *gtk.MenuItem

	configFileManager *config.ConfigFileManager
	multiConfig       *config.MultiAccountConfig

	accounts []Account
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

	u.accounts = BuildAccountsFrom(u.multiConfig, u.configFileManager)

	return nil
}

func (u *gtkUI) addNewAccountsFromConfig() {
	for i := range u.configFileManager.Accounts {
		conf := &u.configFileManager.Accounts[i]

		var found bool
		for _, acc := range u.accounts {
			if acc.Config == conf {
				found = true
				break
			}
		}

		if found {
			continue
		}

		u.accounts = append(u.accounts, newAccount(conf))
	}
}

func (u *gtkUI) SaveConfig() error {
	err := u.configFileManager.Save()
	if err != nil {
		return err
	}

	u.addNewAccountsFromConfig()
	u.window.Emit(AccountChangedSignal.Name())

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
	for i := range u.accounts {
		account := &u.accounts[i]
		if account.Session == s {
			return account
		}
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

func (u *gtkUI) onReceiveSignal(s *glib.Signal, f func()) {
	u.window.Connect(s.Name(), f)
}

func (u *gtkUI) initRoster() {
	u.roster = NewRoster()
}

func (u *gtkUI) mainWindow() {
	u.window = gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	u.applyStyle()
	u.initRoster()

	menubar := initMenuBar(u)
	vbox := gtk.NewVBox(false, 1)
	vbox.PackStart(menubar, false, false, 0)
	vbox.Add(u.roster.Window)
	u.window.Add(vbox)

	u.window.SetTitle(i18n.Local("Coy"))
	u.window.Connect("destroy", gtk.MainQuit)
	u.window.SetSizeRequest(200, 600)

	u.window.ShowAll()
}

func (*gtkUI) askForPassword(connect func(string)) {
	glib.IdleAdd(func() bool {
		dialog := gtk.NewDialog()
		dialog.SetTitle(i18n.Local("Password"))
		dialog.SetPosition(gtk.WIN_POS_CENTER)
		vbox := dialog.GetVBox()

		vbox.Add(gtk.NewLabel(i18n.Local("Password")))
		passwordInput := gtk.NewEntry()
		passwordInput.SetEditable(true)
		passwordInput.SetVisibility(false)
		vbox.Add(passwordInput)

		button := gtk.NewButtonWithLabel(i18n.Local("Send"))
		button.Connect("clicked", func() {
			go connect(passwordInput.GetText())
			dialog.Destroy()
		})
		vbox.Add(button)

		dialog.ShowAll()
		return false
	})
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
	dialog := gtk.NewAboutDialog()
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
		acc := &u.accounts[i]
		if acc.Connected() {
			accounts = append(accounts, acc)
		}
	}

	dialog := presenceSubscriptionDialog(accounts)
	dialog.ShowAll()
}

func (u *gtkUI) buildContactsMenu() *gtk.MenuItem {
	contactsMenu := gtk.NewMenuItemWithMnemonic(i18n.Local("_Contacts"))

	submenu := gtk.NewMenu()
	contactsMenu.SetSubmenu(submenu)

	menuitem := gtk.NewMenuItemWithMnemonic(i18n.Local("_Add..."))
	submenu.Append(menuitem)

	menuitem.Connect("activate", u.addContactWindow)

	return contactsMenu
}

func initMenuBar(u *gtkUI) *gtk.MenuBar {
	menubar := gtk.NewMenuBar()

	menubar.Append(u.buildContactsMenu())

	u.accountsMenu = gtk.NewMenuItemWithMnemonic(i18n.Local("_Accounts"))
	menubar.Append(u.accountsMenu)

	//TODO: replace this by emiting the signal at startup
	u.buildAccountsMenu()
	u.window.Connect(AccountChangedSignal.Name(), func() {
		//TODO: should it destroy the current submenu? HOW?
		u.accountsMenu.SetSubmenu(nil)

		u.buildAccountsMenu()
	})

	//Help -> About
	cascademenu := gtk.NewMenuItemWithMnemonic(i18n.Local("_Help"))
	menubar.Append(cascademenu)
	submenu := gtk.NewMenu()
	cascademenu.SetSubmenu(submenu)
	menuitem := gtk.NewMenuItemWithMnemonic(i18n.Local("_About"))
	menuitem.Connect("activate", aboutDialog)
	submenu.Append(menuitem)
	return menubar
}

func (u *gtkUI) SubscriptionRequest(s *session.Session, from string) {
	confirmDialog := authorizePresenceSubscriptionDialog(u.window, from)

	glib.IdleAdd(func() bool {
		confirm := confirmDialog.Run() == gtk.RESPONSE_YES
		confirmDialog.Destroy()

		s.HandleConfirmOrDeny(from, confirm)

		return false
	})
}

func (u *gtkUI) ProcessPresence(stanza *xmpp.ClientPresence, gone bool) {
	//TODO: Notify via UI
	jid := xmpp.RemoveResourceFromJid(stanza.From)
	status := "available"
	if stanza.Show != "" {
		status = stanza.Show
		if stanza.Status != "" {
			status = status + " (" + stanza.Status + ")"
		}
	}
	fmt.Printf("%s is %s\n", jid, status)
}

func (u *gtkUI) IQReceived(string) {
	//TODO
}

//TODO: we should update periodically (like Pidgin does) if we include the status (online/offline/away) on the label
func (u *gtkUI) RosterReceived(s *session.Session, roster []xmpp.RosterEntry) {
	account := u.findAccountForSession(s)
	if account == nil {
		//TODO error
		return
	}

	u.roster.Update(account, roster)

	glib.IdleAdd(func() bool {
		u.roster.Redraw()
		return false
	})
}

func (u *gtkUI) disconnect(account Account) error {
	account.Session.Close()
	u.window.Emit(account.DisconnectedSignal.Name())
	return nil
}

func (u *gtkUI) connect(account Account) {
	//TODO find a better place to initialize the eventHandler
	s := account.Session
	s.SessionEventHandler = guiSessionEventHandler{u}

	var registerCallback xmpp.FormCallback
	if *config.CreateAccount {
		registerCallback = u.RegisterCallback
	}

	connectFn := func(password string) {
		err := s.Connect(password, registerCallback)
		if err != nil {
			u.window.Emit(account.DisconnectedSignal.Name())
			return
		}

		u.window.Emit(account.ConnectedSignal.Name())
	}

	//TODO We do not support saved empty passwords
	if len(account.Password) == 0 {
		u.askForPassword(connectFn)
		return
	}

	go connectFn(account.Password)
}
