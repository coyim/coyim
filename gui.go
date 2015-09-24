// +build nocli

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"time"
	"unsafe"

	coyconf "github.com/twstrike/coyim/config"
	coyui "github.com/twstrike/coyim/ui"
	"github.com/twstrike/coyim/xmpp"
	"github.com/twstrike/go-gtk/gdk"
	"github.com/twstrike/go-gtk/glib"
	"github.com/twstrike/go-gtk/gtk"
	"github.com/twstrike/otr3"
)

type roster struct {
	window *gtk.ScrolledWindow
	model  *gtk.ListStore
	view   *gtk.TreeView

	sendMessage   func(to, message string)
	conversations map[string]*conversationWindow
}

func newRoster() *roster {
	r := &roster{
		window: gtk.NewScrolledWindow(nil, nil),

		model: gtk.NewListStore(
			gtk.TYPE_STRING, // user
			gtk.TYPE_INT,    // id
		),
		view: gtk.NewTreeView(),

		conversations: make(map[string]*conversationWindow),
	}

	r.window.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_AUTOMATIC)

	r.view.SetHeadersVisible(false)
	r.view.GetSelection().SetMode(gtk.SELECTION_NONE)

	r.view.AppendColumn(
		gtk.NewTreeViewColumnWithAttributes("user",
			gtk.NewCellRendererText(), "text", 0),
	)

	//TODO: Replace by something better
	iter := &gtk.TreeIter{}
	r.model.Append(iter)
	r.model.Set(iter,
		0, "CONNECTING...",
	)

	r.view.SetModel(r.model)
	r.view.Connect("row-activated", r.onActivateBuddy)

	r.window.Add(r.view)

	return r
}

func (r *roster) onActivateBuddy(ctx *glib.CallbackContext) {
	var path *gtk.TreePath
	var column *gtk.TreeViewColumn
	r.view.GetCursor(&path, &column)

	if path == nil {
		return
	}

	iter := &gtk.TreeIter{}
	if !r.model.GetIter(iter, path) {
		return
	}

	val := &glib.GValue{}
	r.model.GetValue(iter, 0, val)
	to := val.GetString()

	r.openConversationWindow(to)

	//TODO call g_value_unset() but the lib doesnt provide
}

func (r *roster) openConversationWindow(to string) *conversationWindow {
	//TODO: There is no validation if this person is on our roster
	c, ok := r.conversations[to]

	if !ok {
		c = newConversationWindow(r, to)
		r.conversations[to] = c
	}

	c.Show()
	return c
}

type conversationWindow struct {
	to            string
	win           *gtk.Window
	history       *gtk.TextView
	scrollHistory *gtk.ScrolledWindow
	roster        *roster
}

func newConversationWindow(r *roster, uid string) *conversationWindow {
	conv := &conversationWindow{
		to:            uid,
		roster:        r,
		win:           gtk.NewWindow(gtk.WINDOW_TOPLEVEL),
		history:       gtk.NewTextView(),
		scrollHistory: gtk.NewScrolledWindow(nil, nil),
	}
	// Unlike the GTK version, this is not supposed to be used as a callback but
	// it attaches the callback to the widget
	conv.win.HideOnDelete()

	conv.win.SetPosition(gtk.WIN_POS_CENTER)
	conv.win.SetDefaultSize(300, 300)
	conv.win.SetDestroyWithParent(true)
	conv.win.SetTitle(uid)

	//TODO: Load recent messages
	conv.history.SetWrapMode(gtk.WRAP_WORD)
	conv.history.SetEditable(false)
	conv.history.SetCursorVisible(false)

	vbox := gtk.NewVBox(false, 1)
	vbox.SetHomogeneous(false)
	vbox.SetSpacing(5)
	vbox.SetBorderWidth(5)

	text := gtk.NewTextView()
	text.SetWrapMode(gtk.WRAP_WORD)
	text.Connect("key-press-event", func(ctx *glib.CallbackContext) bool {
		arg := ctx.Args(0)
		evKey := *(**gdk.EventKey)(unsafe.Pointer(&arg))

		//Send message on ENTER press (without modifier key)
		if evKey.State == 0 && evKey.Keyval == 0xff0d {
			text.SetEditable(false)

			b := text.GetBuffer()
			s := &gtk.TextIter{}
			e := &gtk.TextIter{}
			b.GetStartIter(s)
			b.GetEndIter(e)
			msg := b.GetText(s, e, true)
			b.SetText("")

			text.SetEditable(true)

			conv.sendMessage(msg)
			return true
		}

		return false
	})

	scroll := gtk.NewScrolledWindow(nil, nil)
	scroll.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	scroll.Add(text)

	conv.scrollHistory.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	conv.scrollHistory.Add(conv.history)

	vbox.PackStart(conv.scrollHistory, true, true, 0)
	vbox.Add(scroll)

	conv.win.Add(vbox)

	return conv
}

func (conv *conversationWindow) Hide() {
	conv.win.Hide()
}

func (conv *conversationWindow) Show() {
	conv.win.ShowAll()
}

func (conv *conversationWindow) sendMessage(message string) {
	conv.roster.sendMessage(conv.to, message)
}

func (conv *conversationWindow) appendMessage(from, timestamp string, encrypted bool, message []byte) {
	glib.IdleAdd(func() bool {
		fmt.Println("Appending message", string(message))
		buff := conv.history.GetBuffer()
		buff.InsertAtCursor(timestamp)
		buff.InsertAtCursor(" - ")
		buff.InsertAtCursor(string(message))
		buff.InsertAtCursor("\n")
		adj := conv.scrollHistory.GetVAdjustment()
		adj.SetValue(adj.GetUpper())
		return false
	})
}

func (r *roster) update(entries []xmpp.RosterEntry) {
	gobj := glib.ObjectFromNative(unsafe.Pointer(r.model.GListStore))

	gobj.Ref()
	r.view.SetModel(nil)

	r.model.Clear()
	iter := &gtk.TreeIter{}
	for _, item := range entries {
		r.model.Append(iter)

		//state, ok := s.knownStates[item.Jid]
		// Subscription, knownState
		r.model.Set(iter,
			0, item.Jid,
			// 1, item.Name,
		)
	}

	r.view.SetModel(r.model)
	gobj.Unref()
}

type gtkUI struct {
	roster  *roster
	session *Session
	window  *gtk.Window
}

func (*gtkUI) RegisterCallback() xmpp.FormCallback {
	if *createAccount {
		return func(title, instructions string, fields []interface{}) error {
			//TODO: should open a registration window
			fmt.Println("TODO")
			return nil
		}
	}

	return nil
}

func (u *gtkUI) MessageReceived(from, timestamp string, encrypted bool, message []byte) {
	u.Info("TODO: message received " + from)
	u.Info(string(message))

	//TODO show the message on the history
	glib.IdleAdd(func() bool {
		conv := u.roster.openConversationWindow(from)
		conv.appendMessage(from, timestamp, encrypted, coyui.StripHTML(message))
		return false
	})
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

func (u *gtkUI) Disconnected() {
	fmt.Println("TODO: Should disconnect the account")
}

func (u *gtkUI) Loop() {
	gtk.Init(&os.Args)
	gdk.ThreadsInit()

	gdk.ThreadsEnter()
	u.mainWindow()
	gtk.Main()
	gdk.ThreadsLeave()
}

func NewGTK() *gtkUI {
	return &gtkUI{}
}

func (u *gtkUI) mainWindow() {
	u.window = gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	u.roster = newRoster()
	u.roster.sendMessage = u.sendMessage

	menubar := initMenuBar()
	vbox := gtk.NewVBox(false, 1)
	vbox.PackStart(menubar, false, false, 0)
	vbox.Add(u.roster.window)
	u.window.Add(vbox)

	u.window.SetTitle("Coy")
	u.window.Connect("destroy", gtk.MainQuit)
	u.window.SetSizeRequest(200, 600)
	u.window.ShowAll()
}

func (u *gtkUI) sendMessage(to, message string) {
	conversation := u.session.getConversationWith(to)

	toSend, err := conversation.Send(otr3.ValidMessage(message))
	if err != nil {
		fmt.Println("Failed to generate OTR message")
		return
	}

	glib.IdleAdd(func() bool {
		//TODO: refactor
		conv, _ := u.roster.conversations[to]

		encrypted := conversation.IsEncrypted()
		conv.appendMessage("ME", "NOW", encrypted, coyui.StripHTML([]byte(message)))
		return false
	})

	for _, m := range toSend {
		//TODO: this should be session.Send(to, message)
		fmt.Printf("[send] %q\n", string(m))
		u.session.conn.Send(to, string(m))
	}
}

func (*gtkUI) AskForPassword(*coyconf.Config) (string, error) {
	//TODO
	return "", nil
}

func (*gtkUI) Enroll(*coyconf.Config) bool {
	//TODO
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

func accountDialog() {
	//TODO It should not load config here
	config := &coyconf.Config{}
	dialog := gtk.NewDialog()
	dialog.SetTitle("Account Details")
	dialog.SetPosition(gtk.WIN_POS_CENTER)
	vbox := dialog.GetVBox()

	accountLabel := gtk.NewLabel("Account:")
	vbox.Add(accountLabel)

	accountInput := gtk.NewEntry()
	accountInput.SetText(config.Account)
	accountInput.SetEditable(true)
	vbox.Add(accountInput)

	button := gtk.NewButtonWithLabel("OK")
	button.Connect("clicked", func() {
		fmt.Println(accountInput.GetText())
		dialog.Destroy()
	})
	vbox.Add(button)

	dialog.ShowAll()
}

func initMenuBar() *gtk.MenuBar {
	menubar := gtk.NewMenuBar()

	//Config -> Account
	cascademenu := gtk.NewMenuItemWithMnemonic("_Preference")
	menubar.Append(cascademenu)
	submenu := gtk.NewMenu()
	cascademenu.SetSubmenu(submenu)
	menuitem := gtk.NewMenuItemWithMnemonic("_Account")
	menuitem.Connect("activate", accountDialog)
	submenu.Append(menuitem)

	//Help -> About
	cascademenu = gtk.NewMenuItemWithMnemonic("_Help")
	menubar.Append(cascademenu)
	submenu = gtk.NewMenu()
	cascademenu.SetSubmenu(submenu)
	menuitem = gtk.NewMenuItemWithMnemonic("_About")
	menuitem.Connect("activate", aboutDialog)
	submenu.Append(menuitem)
	return menubar
}

func (u *gtkUI) ProcessPresence(stanza *xmpp.ClientPresence, ignore, gone bool) {

	jid := xmpp.RemoveResourceFromJid(stanza.From)
	state, ok := u.session.knownStates[jid]
	if !ok || len(state) == 0 {
		state = "unknown"
	}

	//TODO: Notify via UI
	fmt.Println(jid, "is", state)
}

func (u *gtkUI) IQReceived(string) {
	//TODO
}

//TODO: we should update periodically (like Pidgin does) if we include the status (online/offline/away) on the label
func (u *gtkUI) RosterReceived(roster []xmpp.RosterEntry) {
	glib.IdleAdd(func() bool {
		u.roster.update(roster)
		return false
	})
}

func main() {
	flag.Parse()

	ui := NewGTK()

	if err := ui.Connect(); err != nil {
		//TODO: Handle error?
		return
	}

	//ticker := time.NewTicker(1 * time.Second)
	//quit := make(chan bool)
	//go timeoutLoop(&s, ticker.C)

	ui.Loop()
	os.Stdout.Write([]byte("\n"))
}

func (u *gtkUI) Connect() error {
	config, password, err := loadConfig(u)
	if err != nil {
		return err
	}

	//TODO support one session per account
	u.session = &Session{
		ui: u,

		account:           config.Account,
		conversations:     make(map[string]*otr3.Conversation),
		eh:                make(map[string]*eventHandler),
		knownStates:       make(map[string]string),
		privateKey:        new(otr3.PrivateKey),
		config:            config,
		pendingRosterChan: make(chan *coyui.RosterEdit),
		pendingSubscribes: make(map[string]string),
		lastActionTime:    time.Now(),
		sessionHandler:    u,
	}

	// TODO: GTK main loop freezes unless this is run on a Go routine
	// and I have no idea why
	go func() {
		conn, err := NewXMPPConn(u, config, password, u.RegisterCallback(), os.Stdout)
		if err != nil {
			u.Alert(err.Error())
			return
		}

		u.session.conn = conn
		u.onConnect()
	}()

	u.session.privateKey.Parse(config.PrivateKey)
	u.session.timeouts = make(map[xmpp.Cookie]time.Time)

	fmt.Printf("Your fingerprint is %x", u.session.privateKey.DefaultFingerprint())

	return nil
}

func (ui *gtkUI) onConnect() {
	go ui.session.WatchTimeout()
	go ui.session.WatchRosterEvents()
	go ui.session.WatchStanzas()
}
