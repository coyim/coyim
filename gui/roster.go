package gui

import (
	"fmt"
	"html"
	"log"
	"sort"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/config"
	rosters "github.com/twstrike/coyim/roster"
	"github.com/twstrike/coyim/ui"
)

type roster struct {
	widget *gtk.ScrolledWindow
	model  *gtk.TreeStore
	view   *gtk.TreeView

	checkEncrypted func(to string) bool
	sendMessage    func(to, message string)

	isCollapsed map[string]bool
	toCollapse  []*gtk.TreePath

	ui *gtkUI
}

const (
	indexJid               = 0
	indexDisplayName       = 1
	indexAccountID         = 2
	indexColor             = 3
	indexBackgroundColor   = 4
	indexWeight            = 5
	indexParentJid         = 0
	indexParentDisplayName = 1
	indexTooltip           = 6
	indexStatusIcon        = 7
	indexRowType           = 8
)

func (u *gtkUI) newRoster() *roster {
	builder := builderForDefinition("Roster")

	r := &roster{
		isCollapsed: make(map[string]bool),

		ui: u,
	}

	builder.ConnectSignals(map[string]interface{}{
		"on_activate_buddy": r.onActivateBuddy,
	})

	obj, _ := builder.GetObject("roster")
	r.widget = obj.(*gtk.ScrolledWindow)

	obj, _ = builder.GetObject("roster-view")
	r.view = obj.(*gtk.TreeView)

	obj, _ = builder.GetObject("roster-model")
	r.model = obj.(*gtk.TreeStore)

	u.displaySettings.update()

	r.view.Connect("button-press-event", r.onButtonPress)

	return r
}

func (r *roster) getAccount(id string) (*account, bool) {
	return r.ui.accountManager.getAccountByID(id)
}

func getFromModelIter(m *gtk.TreeStore, iter *gtk.TreeIter, index int) string {
	val, _ := m.GetValue(iter, index)
	v, _ := val.GetString()
	return v
}

func (r *roster) getAccountAndJidFromEvent(bt *gdk.EventButton) (jid string, account *account, rowType string, ok bool) {
	x := bt.X()
	y := bt.Y()
	path := new(gtk.TreePath)
	found := r.view.GetPathAtPos(int(x), int(y), path, nil, nil, nil)
	if !found {
		return "", nil, "", false
	}
	iter, err := r.model.GetIter(path)
	if err != nil {
		return "", nil, "", false
	}
	jid = getFromModelIter(r.model, iter, indexJid)
	accountID := getFromModelIter(r.model, iter, indexAccountID)
	rowType = getFromModelIter(r.model, iter, indexRowType)
	account, ok = r.getAccount(accountID)
	return jid, account, rowType, ok
}

func (r *roster) createAccountPeerPopup(jid string, account *account, bt *gdk.EventButton) {
	builder := builderForDefinition("ContactPopupMenu")
	mn := getObjIgnoringErrors(builder, "contactMenu").(*gtk.Menu)

	builder.ConnectSignals(map[string]interface{}{
		"on_remove_contact": func() {
			account.session.RemoveContact(jid)
			r.ui.removePeer(account, jid)
			r.redraw()
		},
		"on_allow_contact_to_see_status": func() {
			account.session.ApprovePresenceSubscription(jid, "" /* generate id */)
		},
		"on_forbid_contact_to_see_status": func() {
			account.session.DenyPresenceSubscription(jid, "" /* generate id */)
		},
		"on_ask_contact_to_see_status": func() {
			account.session.RequestPresenceSubscription(jid)
		},
		"on_peer_fingerprints": func() {
			r.ui.showFingerprintsForPeer(jid, account)
		},
		"on_dump_info": func() {
			r.debugPrintRosterFor(account.session.GetConfig().Account)
		},
		"on_rename_signal": func() {
			r.renameContactPopup(account.session.GetConfig(), jid)
		},
	})

	mn.ShowAll()
	mn.PopupAtMouseCursor(nil, nil, int(bt.Button()), bt.Time())
}

func (r *roster) renameContactPopup(conf *config.Account, jid string) {
	builder := builderForDefinition("RenameContact")
	obj, _ := builder.GetObject("RenameContactPopup")
	popup := obj.(*gtk.Dialog)
	builder.ConnectSignals(map[string]interface{}{
		"on_rename_signal": func() {
			obj, _ = builder.GetObject("rename")
			renameTxt := obj.(*gtk.Entry)
			newName, _ := renameTxt.GetText()
			conf.RenamePeer(jid, newName)
			r.ui.SaveConfig()
			r.redraw()
			popup.Destroy()
		},
	})
	popup.SetTransientFor(r.ui.window)
	popup.ShowAll()
}

func (r *roster) createAccountPopup(jid string, account *account, bt *gdk.EventButton) {
	builder := builderForDefinition("AccountPopupMenu")
	obj, _ := builder.GetObject("accountMenu")
	mn := obj.(*gtk.Menu)

	builder.ConnectSignals(map[string]interface{}{
		"on_connect": func() {
			account.session.WantToBeOnline = true
			account.Connect()
		},
		"on_disconnect": func() {
			account.session.WantToBeOnline = false
			account.disconnect()
		},
		"on_edit": account.edit,
		"on_dump_info": func() {
			r.debugPrintRosterFor(account.session.GetConfig().Account)
		},
	})

	connx, _ := builder.GetObject("connectMenuItem")
	connect := connx.(*gtk.MenuItem)

	dconnx, _ := builder.GetObject("disconnectMenuItem")
	disconnect := dconnx.(*gtk.MenuItem)

	connect.SetSensitive(account.session.IsDisconnected())
	disconnect.SetSensitive(account.session.IsConnected())

	mn.ShowAll()
	mn.PopupAtMouseCursor(nil, nil, int(bt.Button()), bt.Time())
}

func (r *roster) onButtonPress(view *gtk.TreeView, ev *gdk.Event) bool {
	bt := &gdk.EventButton{ev}
	if bt.Button() == 0x03 {
		jid, account, rowType, ok := r.getAccountAndJidFromEvent(bt)
		if ok {
			switch rowType {
			case "peer":
				r.createAccountPeerPopup(jid, account, bt)
			case "account":
				r.createAccountPopup(jid, account, bt)
			}
		}
	}

	return false
}

func (r *roster) onActivateBuddy(v *gtk.TreeView, path *gtk.TreePath) {
	selection, _ := v.GetSelection()
	defer selection.UnselectPath(path)

	iter, err := r.model.GetIter(path)
	if err != nil {
		return
	}

	jid := getFromModelIter(r.model, iter, indexJid)
	accountID := getFromModelIter(r.model, iter, indexAccountID)
	rowType := getFromModelIter(r.model, iter, indexRowType)

	if rowType != "peer" {
		r.isCollapsed[jid] = !r.isCollapsed[jid]
		r.redraw()
		return
	}

	account, ok := r.getAccount(accountID)
	if !ok {
		return
	}

	r.openConversationWindow(account, jid)
}

func (r *roster) openConversationWindow(account *account, to string) (*conversationWindow, error) {
	c, ok := account.getConversationWith(to)

	if !ok {
		textBuffer := r.ui.getTags().createTextBuffer()
		c = account.createConversationWindow(to, r.ui.displaySettings, textBuffer)

		r.ui.connectShortcutsChildWindow(c.win)
		r.ui.connectShortcutsConversationWindow(c)
		c.parentWin = &r.ui.window.Window
	}

	c.Show()
	return c, nil
}

func (r *roster) displayNameFor(account *account, from string) string {
	p, ok := r.ui.getPeer(account, from)
	if !ok {
		return from
	}

	return p.NameForPresentation()
}

func (r *roster) presenceUpdated(account *account, from, show, showStatus string, gone bool) {
	c, ok := account.getConversationWith(from)
	if !ok {
		return
	}

	doInUIThread(func() {
		c.appendStatus(r.displayNameFor(account, from), time.Now(), show, showStatus, gone)
	})
}

func (r *roster) messageReceived(account *account, from string, timestamp time.Time, encrypted bool, message []byte) {
	doInUIThread(func() {
		conv, err := r.openConversationWindow(account, from)
		if err != nil {
			return
		}

		conv.appendMessage(r.displayNameFor(account, from), timestamp, encrypted, ui.StripHTML(message), false)
	})
}

func (r *roster) update(account *account, entries *rosters.List) {
	r.ui.accountManager.Lock()
	defer r.ui.accountManager.Unlock()

	r.ui.accountManager.setContacts(account, entries)
}

func (r *roster) debugPrintRosterFor(nm string) {
	r.ui.accountManager.RLock()
	defer r.ui.accountManager.RUnlock()

	for account, rs := range r.ui.accountManager.getAllContacts() {
		if account.session.GetConfig().Is(nm) {
			rs.Iter(func(_ int, item *rosters.Peer) {
				fmt.Printf("->   %s\n", item.Dump())
			})
		}
	}

	fmt.Printf(" ************************************** \n")
	fmt.Println()
}

func isNominallyVisible(p *rosters.Peer) bool {
	return (p.Subscription != "none" && p.Subscription != "") || p.PendingSubscribeID != ""
}

func shouldDisplay(p *rosters.Peer, showOffline bool) bool {
	return isNominallyVisible(p) && (showOffline || p.Online)
}

func isAway(p *rosters.Peer) bool {
	switch p.Status {
	case "dnd", "xa", "away":
		return true
	}
	return false
}

func isOnline(p *rosters.Peer) bool {
	return p.PendingSubscribeID == "" && p.Online
}

func decideStatusFor(p *rosters.Peer) string {
	if p.PendingSubscribeID != "" {
		return "unknown"
	}

	if !p.Online {
		return "offline"
	}

	switch p.Status {
	case "dnd":
		return "busy"
	case "xa":
		return "extended-away"
	case "away":
		return "away"
	}

	return "available"
}

func decideColorFor(p *rosters.Peer) string {
	if !p.Online {
		return "#aaaaaa"
	}
	return "#000000"
}

func createGroupDisplayName(parentName string, counter *counter, isExpanded bool) string {
	name := parentName
	if !isExpanded {
		name = fmt.Sprintf("[%s]", name)
	}
	return fmt.Sprintf("%s (%d/%d)", name, counter.online, counter.total)
}

func createTooltipFor(item *rosters.Peer) string {
	pname := html.EscapeString(item.NameForPresentation())
	jid := html.EscapeString(item.Jid)
	if pname != jid {
		return fmt.Sprintf("%s (%s)", pname, jid)
	}
	return jid
}

func (r *roster) addItem(item *rosters.Peer, parentIter *gtk.TreeIter, indent string) {
	iter := r.model.Append(parentIter)
	setAll(r.model, iter,
		item.Jid,
		fmt.Sprintf("%s %s", indent, item.NameForPresentation()),
		item.BelongsTo,
		decideColorFor(item),
		"#ffffff",
		nil,
		createTooltipFor(item),
	)

	r.model.SetValue(iter, indexRowType, "peer")
	r.model.SetValue(iter, indexStatusIcon, statusIcons[decideStatusFor(item)].getPixbuf())
}

func (r *roster) redrawMerged() {
	showOffline := !r.ui.config.Display.ShowOnlyOnline

	r.ui.accountManager.RLock()
	defer r.ui.accountManager.RUnlock()

	r.toCollapse = nil

	grp := rosters.TopLevelGroup()
	for account, contacts := range r.ui.accountManager.getAllContacts() {
		contacts.AddTo(grp, account.session.GroupDelimiter)
	}

	accountCounter := &counter{}
	r.displayGroup(grp, nil, accountCounter, showOffline, "")

	r.view.ExpandAll()
	for _, path := range r.toCollapse {
		r.view.CollapseRow(path)
	}
}

type counter struct {
	total  int
	online int
}

func (c *counter) inc(total, online bool) {
	if total {
		c.total++
	}
	if online {
		c.online++
	}
}

func (r *roster) displayGroup(g *rosters.Group, parentIter *gtk.TreeIter, accountCounter *counter, showOffline bool, accountName string) {
	pi := parentIter
	groupCounter := &counter{}
	groupID := accountName + "//" + g.FullGroupName()
	if g.GroupName != "" {
		pi = r.model.Append(parentIter)
		r.model.SetValue(pi, indexParentJid, groupID)
		r.model.SetValue(pi, indexRowType, "group")
		r.model.SetValue(pi, indexWeight, 500)
		r.model.SetValue(pi, indexBackgroundColor, "#e9e7f3")
	}

	for _, item := range g.Peers() {
		vs := isNominallyVisible(item)
		o := isOnline(item)
		accountCounter.inc(vs, vs && o)
		groupCounter.inc(vs, vs && o)

		if shouldDisplay(item, showOffline) {
			r.addItem(item, pi, "")
		}
	}

	for _, gr := range g.Groups() {
		r.displayGroup(gr, pi, accountCounter, showOffline, accountName)
	}

	if g.GroupName != "" {
		parentPath, _ := r.model.GetPath(pi)
		shouldCollapse, ok := r.isCollapsed[groupID]
		isExpanded := true
		if ok && shouldCollapse {
			isExpanded = false
			r.toCollapse = append(r.toCollapse, parentPath)
		}

		r.model.SetValue(pi, indexParentDisplayName, createGroupDisplayName(g.FullGroupName(), groupCounter, isExpanded))
	}
}

func (r *roster) redrawSeparateAccount(account *account, contacts *rosters.List, showOffline bool) {
	parentIter := r.model.Append(nil)

	accountCounter := &counter{}

	grp := contacts.Grouped(account.session.GroupDelimiter)
	parentName := account.session.GetConfig().Account
	r.displayGroup(grp, parentIter, accountCounter, showOffline, parentName)

	r.model.SetValue(parentIter, indexParentJid, parentName)
	r.model.SetValue(parentIter, indexAccountID, account.session.GetConfig().ID())
	r.model.SetValue(parentIter, indexRowType, "account")
	r.model.SetValue(parentIter, indexWeight, 700)

	bgcolor := "#918caa"
	if account.session.IsDisconnected() {
		bgcolor = "#d5d3de"
	}
	r.model.SetValue(parentIter, indexBackgroundColor, bgcolor)

	parentPath, _ := r.model.GetPath(parentIter)
	shouldCollapse, ok := r.isCollapsed[parentName]
	isExpanded := true
	if ok && shouldCollapse {
		isExpanded = false
		r.toCollapse = append(r.toCollapse, parentPath)
	}
	var stat string
	if account.session.IsDisconnected() {
		stat = "offline"
	} else if account.session.IsConnected() {
		stat = "available"
	} else {
		stat = "connecting"
	}

	r.model.SetValue(parentIter, indexStatusIcon, statusIcons[stat].getPixbuf())
	r.model.SetValue(parentIter, indexParentDisplayName, createGroupDisplayName(parentName, accountCounter, isExpanded))
}

func (r *roster) sortedAccounts() []*account {
	var as []*account
	for account := range r.ui.accountManager.getAllContacts() {
		if account == nil {
			log.Printf("adding an account that is nil...\n")
		}
		as = append(as, account)
	}
	//TODO sort by nickname if available
	sort.Sort(byAccountNameAlphabetic(as))
	return as
}

func (r *roster) redrawSeparate() {
	showOffline := !r.ui.config.Display.ShowOnlyOnline

	r.ui.accountManager.RLock()
	defer r.ui.accountManager.RUnlock()

	r.toCollapse = nil

	for _, account := range r.sortedAccounts() {
		r.redrawSeparateAccount(account, r.ui.accountManager.getContacts(account), showOffline)
	}

	r.view.ExpandAll()
	for _, path := range r.toCollapse {
		r.view.CollapseRow(path)
	}
}

const disconnectedPageIndex = 0
const spinnerPageIndex = 1
const rosterPageIndex = 2

func (r *roster) redraw() {
	//TODO: this should be behind a mutex
	r.model.Clear()

	if r.ui.shouldViewAccounts() {
		r.redrawSeparate()
	} else {
		r.redrawMerged()
	}
}

func setAll(v *gtk.TreeStore, iter *gtk.TreeIter, values ...interface{}) {
	for i, val := range values {
		if val != nil {
			v.SetValue(iter, i, val)
		}
	}
}
