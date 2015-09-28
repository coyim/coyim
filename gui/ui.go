package gui

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/event"
	"github.com/twstrike/coyim/session"
	"github.com/twstrike/coyim/ui"
	"github.com/twstrike/coyim/xmpp"

	"github.com/twstrike/go-gtk/gdk"
	"github.com/twstrike/go-gtk/glib"
	"github.com/twstrike/go-gtk/gtk"
	"github.com/twstrike/otr3"
)

//TODO: this signal should be per Session
var (
	CONNECTED_SIG    = glib.NewSignal("coyim-account-connected")
	DISCONNECTED_SIG = glib.NewSignal("coyim-account-disconnected")
)

type gtkUI struct {
	roster *Roster
	window *gtk.Window

	multiConfig *config.MultiAccountConfig
	sessions    map[*config.Config]*session.Session

	connected bool
}

func NewGTK() *gtkUI {
	return &gtkUI{
		sessions: make(map[*config.Config]*session.Session),
	}
}

func (ui *gtkUI) LoadConfig(configFile string) {
	var err error
	if ui.multiConfig, err = config.Load(configFile); err != nil {
		ui.Alert(err.Error())
		ui.enroll()
		return
	}
}

//TODO: Should be per session
func (*gtkUI) Disconnected() {
	//TODO: remove everybody from the roster
	fmt.Println("TODO: Should disconnect the account")
}

func (*gtkUI) RegisterCallback() xmpp.FormCallback {
	//if !*createAccount {
	//  return nil
	//}
	return nil

	return func(title, instructions string, fields []interface{}) error {
		//TODO: should open a registration window
		fmt.Println("TODO")
		return nil
	}
}

func (u *gtkUI) MessageReceived(from, timestamp string, encrypted bool, message []byte) {
	u.roster.MessageReceived(from, timestamp, encrypted, message)
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
	gdk.ThreadsInit()

	gdk.ThreadsEnter()
	u.mainWindow()
	gtk.Main()
	gdk.ThreadsLeave()
}

func (u *gtkUI) onReceiveSignal(s *glib.Signal, f func()) {
	u.window.Connect(s.Name(), f)
}

func (u *gtkUI) initRoster() {
	u.roster = NewRoster()
	u.onReceiveSignal(DISCONNECTED_SIG, u.roster.Clear)
}

func (u *gtkUI) mainWindow() {
	u.window = gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	u.initRoster()

	menubar := initMenuBar(u)
	vbox := gtk.NewVBox(false, 1)
	vbox.PackStart(menubar, false, false, 0)
	vbox.Add(u.roster.Window)
	u.window.Add(vbox)

	u.window.SetTitle("Coy")
	u.window.Connect("destroy", gtk.MainQuit)
	u.window.SetSizeRequest(200, 600)
	u.window.ShowAll()
}

//TODO: REMOVE ME
func (*gtkUI) AskForPassword(c *config.Config) (string, error) {
	return "", nil
}

func (*gtkUI) askForPassword(connect func(string)) {
	glib.IdleAdd(func() bool {
		dialog := gtk.NewDialog()
		dialog.SetTitle("Password")
		dialog.SetPosition(gtk.WIN_POS_CENTER)
		vbox := dialog.GetVBox()

		vbox.Add(gtk.NewLabel("Password"))
		passwordInput := gtk.NewEntry()
		passwordInput.SetEditable(true)
		passwordInput.SetVisibility(false)
		vbox.Add(passwordInput)

		button := gtk.NewButtonWithLabel("Send")
		button.Connect("clicked", func() {
			go connect(passwordInput.GetText())
			dialog.Destroy()
		})
		vbox.Add(button)

		dialog.ShowAll()
		return false
	})
}

//TODO: Remove?
func (*gtkUI) Enroll(c *config.Config) bool {
	return false
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
	dialog.SetName("Coy IM!")
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

func accountDialog(c *config.Config) {
	dialog := gtk.NewDialog()
	dialog.SetTitle("Account Details")
	dialog.SetPosition(gtk.WIN_POS_CENTER)
	vbox := dialog.GetVBox()

	accountLabel := gtk.NewLabel("Account")
	vbox.Add(accountLabel)

	accountInput := gtk.NewEntry()
	accountInput.SetText(c.Account)
	accountInput.SetEditable(true)
	vbox.Add(accountInput)

	vbox.Add(gtk.NewLabel("Password"))
	passwordInput := gtk.NewEntry()
	passwordInput.SetText(c.Password)
	passwordInput.SetEditable(true)
	passwordInput.SetVisibility(false)
	vbox.Add(passwordInput)

	vbox.Add(gtk.NewLabel("Server"))
	serverInput := gtk.NewEntry()
	serverInput.SetText(c.Server)
	serverInput.SetEditable(true)
	vbox.Add(serverInput)

	vbox.Add(gtk.NewLabel("Port"))
	portInput := gtk.NewEntry()
	portInput.SetText(strconv.Itoa(c.Port))
	portInput.SetEditable(true)
	vbox.Add(portInput)

	vbox.Add(gtk.NewLabel("Tor Proxy"))
	proxyInput := gtk.NewEntry()
	if len(c.Proxies) > 0 {
		proxyInput.SetText(c.Proxies[0])
	}
	proxyInput.SetEditable(true)
	vbox.Add(proxyInput)

	alwaysEncrypt := gtk.NewCheckButtonWithLabel("Always Encrypt")
	alwaysEncrypt.SetActive(c.AlwaysEncrypt)
	vbox.Add(alwaysEncrypt)

	button := gtk.NewButtonWithLabel("Save")
	button.Connect("clicked", func() {
		c.Account = accountInput.GetText()
		c.Password = passwordInput.GetText()
		c.Server = serverInput.GetText()

		v, err := strconv.Atoi(portInput.GetText())
		if err == nil {
			c.Port = v
		}

		if len(c.Proxies) == 0 {
			c.Proxies = append(c.Proxies, "")
		}
		c.Proxies[0] = proxyInput.GetText()

		c.AlwaysEncrypt = alwaysEncrypt.GetActive()

		c.Save()
		dialog.Destroy()
	})
	vbox.Add(button)

	dialog.ShowAll()
}

func buildAccountSubmenu(u *gtkUI, c *config.Config) *gtk.MenuItem {
	menuitem := gtk.NewMenuItemWithMnemonic(c.Account)

	accountSubMenu := gtk.NewMenu()
	menuitem.SetSubmenu(accountSubMenu)

	connectItem := gtk.NewMenuItemWithMnemonic("_Connect")
	accountSubMenu.Append(connectItem)

	disconnectItem := gtk.NewMenuItemWithMnemonic("_Disconnect")
	//disconnectItem.SetSensitive(false)
	accountSubMenu.Append(disconnectItem)

	connectItem.Connect("activate", func() {
		//connectItem.SetSensitive(false)
		u.connect(c)
	})

	disconnectItem.Connect("activate", func() {
		u.disconnect(c)
	})

	//TODO: this should be per session
	//connToggle := func() {
	//	s, ok := u.sessions[c]
	//	if !ok {
	//		return
	//	}

	//	connected := s.ConnStatus == session.CONNECTED
	//	connectItem.SetSensitive(!connected)
	//	disconnectItem.SetSensitive(connected)
	//}

	//u.onReceiveSignal(CONNECTED_SIG, connToggle)
	//u.onReceiveSignal(DISCONNECTED_SIG, connToggle)

	editItem := gtk.NewMenuItemWithMnemonic("_Edit")
	editItem.Connect("activate", func() {
		accountDialog(c)
	})
	accountSubMenu.Append(editItem)

	return menuitem
}

func initMenuBar(u *gtkUI) *gtk.MenuBar {
	menubar := gtk.NewMenuBar()

	//Config -> Account
	cascademenu := gtk.NewMenuItemWithMnemonic("_Accounts")
	menubar.Append(cascademenu)
	submenu := gtk.NewMenu()
	cascademenu.SetSubmenu(submenu)

	for i := range u.multiConfig.Accounts {
		p := &u.multiConfig.Accounts[i]
		submenu.Append(buildAccountSubmenu(u, p))
	}

	//Help -> About
	cascademenu = gtk.NewMenuItemWithMnemonic("_Help")
	menubar.Append(cascademenu)
	submenu = gtk.NewMenu()
	cascademenu.SetSubmenu(submenu)
	menuitem := gtk.NewMenuItemWithMnemonic("_About")
	menuitem.Connect("activate", aboutDialog)
	submenu.Append(menuitem)
	return menubar
}

func (u *gtkUI) SubscriptionRequest(from string) {
}

func (u *gtkUI) ProcessPresence(stanza *xmpp.ClientPresence, gone bool) {

	//TODO: Notify via UI
	jid := xmpp.RemoveResourceFromJid(stanza.From)
	fmt.Println(jid, "is", stanza.Show)
}

func (u *gtkUI) IQReceived(string) {
	//TODO
}

//TODO: we should update periodically (like Pidgin does) if we include the status (online/offline/away) on the label
func (u *gtkUI) RosterReceived(s *session.Session, roster []xmpp.RosterEntry) {
	glib.IdleAdd(func() bool {
		u.roster.Update(s, roster)
		u.roster.Redraw()
		return false
	})
}

func (u *gtkUI) enroll() {
	//TODO: import private key?

	filename, err := config.FindConfigFile(os.Getenv("HOME"))
	if err != nil {
		//TODO cant write config file. Should it be a problem?
		return
	}

	//TODO: extract to function when implementing "add account"
	u.multiConfig = &config.MultiAccountConfig{
		Filename: *filename,
		Accounts: []config.Config{
			config.Config{},
		},
	}

	glib.IdleAdd(func() bool {
		accountDialog(&u.multiConfig.Accounts[0])
		return false
	})
}

func (u *gtkUI) disconnect(c *config.Config) error {
	//if !u.connected {
	//	return nil
	//}

	s, ok := u.sessions[c]
	if !ok {
		return errors.New("tried to disconnect an unexisting session")
	}

	s.Terminate()
	delete(u.sessions, c)

	u.window.EmitSignal(DISCONNECTED_SIG)
	//u.connected = false

	return nil
}

func newSession(u *gtkUI, c *config.Config) *session.Session {
	s := &session.Session{
		Config: c,

		Conversations:     make(map[string]*otr3.Conversation),
		OtrEventHandler:   make(map[string]*event.OtrEventHandler),
		KnownStates:       make(map[string]string),
		PrivateKey:        new(otr3.PrivateKey),
		PendingRosterChan: make(chan *ui.RosterEdit),
		PendingSubscribes: make(map[string]string),
		LastActionTime:    time.Now(),

		SessionEventHandler: u,
	}

	s.PrivateKey.Parse(c.PrivateKey)

	//TODO: This should happen regardless of connecting
	fmt.Printf("Your fingerprint is %x\n", s.PrivateKey.DefaultFingerprint())

	return s
}

func (u *gtkUI) findOrCreateSession(c *config.Config) *session.Session {
	s, ok := u.sessions[c]
	if !ok {
		s = newSession(u, c)
		u.sessions[c] = s
	}

	return s
}

func (u *gtkUI) connect(c *config.Config) {
	//if u.connected {
	//	return
	//}

	s := u.findOrCreateSession(c)

	connectFn := func(password string) {
		err := s.Connect(password)
		if err != nil {
			return
		}

		//TODO: change this signal to be session specific
		u.window.EmitSignal(CONNECTED_SIG)
		//u.connected = true
		u.onConnect(s)
	}

	//TODO We do not support empty passwords
	if len(c.Password) == 0 {
		u.askForPassword(connectFn)
		return
	}

	go connectFn(c.Password)
}

func (ui *gtkUI) onConnect(s *session.Session) {
	go s.WatchTimeout()
	go s.WatchRosterEvents()
	go s.WatchStanzas()
}
