// +build nocli

package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"time"
	"unsafe"

	. "github.com/twstrike/coyim/config"
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

func (r *roster) openConversationWindow(to string) {
	if c, ok := r.conversations[to]; ok {
		c.Show()
		return
	}

	r.conversations[to] = newConversationWindow(to)

}

type conversationWindow struct {
	to  string
	win *gtk.Window
}

func newConversationWindow(uid string) *conversationWindow {
	conv := &conversationWindow{
		to:  uid,
		win: gtk.NewWindow(gtk.WINDOW_TOPLEVEL),
	}

	conv.win.SetPosition(gtk.WIN_POS_CENTER)
	conv.win.SetDefaultSize(300, 300)
	conv.win.SetDestroyWithParent(true)

	// Unlike the GTK version, this is not supposed to be used as a callback but
	// it attaches the callback to the widget
	conv.win.HideOnDelete()

	textview := gtk.NewTextView()
	textview.SetEditable(false)
	textview.SetCursorVisible(false)
	buffer := textview.GetBuffer()
	//TODO: Load recent messages
	buffer.InsertAtCursor("** History here **")

	vbox := gtk.NewVBox(false, 1)
	vbox.SetHomogeneous(false)
	vbox.SetSpacing(5)
	vbox.SetBorderWidth(5)

	text := gtk.NewTextView()
	text.Connect("key_press_event", func(ctx *glib.CallbackContext) {
		arg := ctx.Args(0)
		evKey := *(**gdk.EventKey)(unsafe.Pointer(&arg))

		if evKey.Keyval == 0xff0d {
			msg := ""
			conv.sendMessage(msg)

			//target := ctx.Target()
			//if t, ok := target.(*gtk.TextView); ok {
			//	t.GetBuffer().SetText("")
			//}
		}
	})

	scroll := gtk.NewScrolledWindow(nil, nil)
	scroll.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	scroll.Add(text)

	vbox.PackStart(textview, true, true, 0)
	vbox.PackStart(scroll, true, true, 0)

	conv.win.Add(vbox)
	conv.win.ShowAll()

	return conv
}

func (conv *conversationWindow) Hide() {
	conv.win.Hide()
}

func (conv *conversationWindow) Show() {
	conv.win.Show()
}

func (*conversationWindow) sendMessage(message string) {
	fmt.Println("SEND MESSAGE")

	/*
		conversation, ok := s.conversations[cmd.to]
		isEncrypted := ok && conversation.IsEncrypted()
		if cmd.setPromptIsEncrypted != nil {
			cmd.setPromptIsEncrypted <- isEncrypted
		}
		if !isEncrypted && config.ShouldEncryptTo(cmd.to) {
			warn(ui.term, fmt.Sprintf("Did not send: no encryption established with %s", cmd.to))
			continue
		}
		var msgs [][]byte
		message := []byte(cmd.msg)
		// Automatically tag all outgoing plaintext
		// messages with a whitespace tag that
		// indicates that we support OTR.
		if config.OTRAutoAppendTag &&
			!bytes.Contains(message, []byte("?OTR")) &&
			(!ok || !conversation.IsEncrypted()) {
			message = append(message, OTRWhitespaceTag...)
		}
		if ok {
			var err error
			validMsgs, err := conversation.Send(message)
			msgs = otr3.Bytes(validMsgs)
			if err != nil {
				alert(ui.term, err.Error())
				break
			}
		} else {
			msgs = [][]byte{[]byte(message)}
		}
		for _, message := range msgs {
			s.conn.Send(cmd.to, string(message))
		}
	*/
}

func (r *roster) update(entries []xmpp.RosterEntry) {
	gobj := glib.ObjectFromNative(unsafe.Pointer(r.model.GListStore))

	gdk.ThreadsEnter()
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
	gdk.ThreadsLeave()
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
	u.Info("TODO: message received")
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

func NewGTK() *gtkUI {
	return &gtkUI{}
}

func (ui *gtkUI) mainWindow() {
	ui.window = gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	ui.roster = newRoster()

	menubar := initMenuBar()
	vbox := gtk.NewVBox(false, 1)
	vbox.PackStart(menubar, false, false, 0)
	vbox.Add(ui.roster.window)
	ui.window.Add(vbox)

	ui.window.SetTitle("Coy")
	ui.window.Connect("destroy", gtk.MainQuit)
	ui.window.SetSizeRequest(200, 600)
	ui.window.ShowAll()
}

func (*gtkUI) AskForPassword(*Config) (string, error) {
	//TODO
	return "", nil
}

func (*gtkUI) Enroll(*Config) bool {
	//TODO
	return false
}

//TODO: we should update periodically (like Pidgin does) if we include the status (online/offline/away) on the label
func (ui *gtkUI) updateRoster(roster []xmpp.RosterEntry) {
	ui.roster.update(roster)
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
	config := &Config{}
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

func (ui *gtkUI) ProcessPresence(stanza *xmpp.ClientPresence) {
	jid := xmpp.RemoveResourceFromJid(stanza.From)
	state, ok := ui.session.knownStates[jid]
	if !ok || len(state) == 0 {
		state = "unknown"
	}

	//TODO: Notify via UI
	fmt.Println(jid, "is", state)
}

func (ui *gtkUI) IQReceived(string) {
	//TODO
}

func (ui *gtkUI) RosterReceived(roster []xmpp.RosterEntry) {
	//TODO
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

func (ui *gtkUI) Connect() error {
	config, password, err := loadConfig(ui)
	if err != nil {
		return err
	}

	//TODO support one session per account
	ui.session = &Session{
		ui: ui,

		account:           config.Account,
		conversations:     make(map[string]*otr3.Conversation),
		eh:                make(map[string]*eventHandler),
		knownStates:       make(map[string]string),
		privateKey:        new(otr3.PrivateKey),
		config:            config,
		pendingRosterChan: make(chan *rosterEdit),
		pendingSubscribes: make(map[string]string),
		lastActionTime:    time.Now(),
	}

	// TODO: GTK main loop freezes unless this is run on a Go routine
	// and I have no idea why
	go func() {
		logger := bytes.NewBuffer(nil)
		conn, err := NewXMPPConn(ui, config, password, ui.RegisterCallback(), logger)
		if err != nil {
			ui.Alert(err.Error())
			//gtk.MainQuit()
			return
		}

		ui.session.conn = conn
		ui.session.conn.SignalPresence("")
		ui.onConnect()
	}()

	ui.session.privateKey.Parse(config.PrivateKey)
	ui.session.timeouts = make(map[xmpp.Cookie]time.Time)

	fmt.Printf("Your fingerprint is %x", ui.session.privateKey.DefaultFingerprint())

	return nil
}

func (ui *gtkUI) onConnect() {
	go ui.handleRosterEvents()
	go ui.handleStanzaEvents()
}

//TODO: use the rosterLoop() from main_loop.go
func (ui *gtkUI) handleRosterEvents() {
	s := ui.session

	fmt.Println("Fetching roster")
	rosterReply, _, err := s.conn.RequestRoster()
	if err != nil {
		fmt.Println("Failed to request roster: " + err.Error())
		return
	}

RosterLoop:
	for {
		select {
		case rosterStanza, ok := <-rosterReply:
			if !ok {
				ui.Alert("Failed to read roster: " + err.Error())
				break RosterLoop
			}

			if s.roster, err = xmpp.ParseRoster(rosterStanza); err != nil {
				ui.Alert("Failed to parse roster: " + err.Error())
				break RosterLoop
			}

			ui.updateRoster(s.roster)

			fmt.Println("Roster received")

		case edit := <-s.pendingRosterChan:
			if !edit.isComplete {
				ui.Info("Please edit " + edit.fileName + " and run /rostereditdone when complete")
				s.pendingRosterEdit = edit
				continue
			}
			if s.processEditedRoster(edit) {
				s.pendingRosterEdit = nil
			} else {
				ui.Alert("Please reedit file and run /rostereditdone again")
			}
		}
	}

	//TODO should it quit?
	gdk.ThreadsEnter()
	gtk.MainQuit()
	gdk.ThreadsLeave()
}

func (ui *gtkUI) handleStanzaEvents() {
	stanzaChan := make(chan xmpp.Stanza)
	go ui.session.readMessages(stanzaChan)

StanzaLoop:
	for {
		select {
		case rawStanza, ok := <-stanzaChan:
			if !ok {
				fmt.Println("Stanza channel closed")
				break StanzaLoop
			}

			switch stanza := rawStanza.Value.(type) {
			case *xmpp.StreamError:
				var text string
				if len(stanza.Text) > 0 {
					text = stanza.Text
				} else {
					text = fmt.Sprintf("%s", stanza.Any)
				}
				fmt.Println("Exiting in response to fatal error from server: " + text)
				break StanzaLoop
			case *xmpp.ClientMessage:
				ui.session.processClientMessage(stanza)
			case *xmpp.ClientPresence:
				//OKish
				ui.session.processPresence(stanza)
				ui.ProcessPresence(stanza)
			case *xmpp.ClientIQ:
				if stanza.Type != "get" && stanza.Type != "set" {
					continue
				}
				reply := ui.session.processIQ(stanza)
				if reply == nil {
					reply = xmpp.ErrorReply{
						Type:  "cancel",
						Error: xmpp.ErrorBadRequest{},
					}
				}

				if err := ui.session.conn.SendIQReply(stanza.From, "result", stanza.Id, reply); err != nil {
					fmt.Println("Failed to send IQ message: " + err.Error())
				}
			default:
				fmt.Println(fmt.Sprintf("%s %s", rawStanza.Name, rawStanza.Value))
			}
		}
	}

	//TODO should it quit?
	gdk.ThreadsEnter()
	gtk.MainQuit()
	gdk.ThreadsLeave()
}
