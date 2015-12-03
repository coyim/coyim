package gui

import (
	"fmt"
	"html"
	"sort"
	"sync"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	rosters "github.com/twstrike/coyim/roster"
	"github.com/twstrike/coyim/ui"
	"github.com/twstrike/coyim/xmpp"
)

type contacts struct {
	sync.RWMutex
	m map[*account]*rosters.List
}

type roster struct {
	widget *gtk.Notebook

	model *gtk.TreeStore
	view  *gtk.TreeView

	contacts       contacts
	checkEncrypted func(to string) bool
	sendMessage    func(to, message string)
	conversations  map[string]*conversationWindow

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
)

func (u *gtkUI) newRoster() *roster {
	builder, err := loadBuilderWith("Roster")
	if err != nil {
		panic(err)
	}

	r := &roster{
		conversations: make(map[string]*conversationWindow),
		contacts: contacts{
			m: make(map[*account]*rosters.List),
		},

		isCollapsed: make(map[string]bool),

		ui: u,
	}

	builder.ConnectSignals(map[string]interface{}{
		"on_activate_buddy": r.onActivateBuddy,
	})

	obj, _ := builder.GetObject("notebook")
	r.widget = obj.(*gtk.Notebook)

	obj, _ = builder.GetObject("roster-view")
	r.view = obj.(*gtk.TreeView)

	obj, _ = builder.GetObject("roster-model")
	r.model = obj.(*gtk.TreeStore)

	//TODO: cant this be achieved by using CSS classes / IDs?
	//For example: #notebook GtkBox { .... }
	obj, _ = builder.GetObject("disconnected-page")
	vbox := obj.(*gtk.Box)
	u.displaySettings.unifiedBackgroundColor(&vbox.Container.Widget)

	obj, _ = builder.GetObject("spinner-page")
	vboxSpinner := obj.(*gtk.Box)
	u.displaySettings.unifiedBackgroundColor(&vboxSpinner.Container.Widget)

	u.displaySettings.update()

	r.view.Connect("button-press-event", r.onButtonPress)

	return r
}

func (r *roster) connected() {
	glib.IdleAdd(func() {
		r.widget.SetCurrentPage(rosterPageIndex)
	})
}

func (r *roster) connecting() {
	//only if there is nothing connected
	if len(r.contacts.m) != 0 {
		return
	}

	glib.IdleAdd(func() {
		r.widget.SetCurrentPage(spinnerPageIndex)
	})
}

func (r *roster) disconnected() {
	glib.IdleAdd(func() {
		r.widget.SetCurrentPage(disconnectedPageIndex)
	})
}

//TODO: move somewhere else
func (r *roster) getAccount(id string) (*account, bool) {
	r.contacts.RLock()
	defer r.contacts.RUnlock()

	for account := range r.contacts.m {
		if account.session.CurrentAccount.ID() == id {
			return account, true
		}
	}

	return nil, false
}

func getFromModelIter(m *gtk.TreeStore, iter *gtk.TreeIter, index int) string {
	val, _ := m.GetValue(iter, index)
	v, _ := val.GetString()
	return v
}

func (r *roster) getAccountAndJidFromEvent(bt *gdk.EventButton) (jid string, account *account, ok bool) {
	x := bt.X()
	y := bt.Y()
	path := new(gtk.TreePath)
	found := r.view.GetPathAtPos(int(x), int(y), path, nil, nil, nil)
	if !found {
		return "", nil, false
	}
	iter, err := r.model.GetIter(path)
	if err != nil {
		return "", nil, false
	}
	jid = getFromModelIter(r.model, iter, indexJid)
	accountID := getFromModelIter(r.model, iter, indexAccountID)
	if accountID != accountIDTopLevelMarker {
		account, ok = r.getAccount(accountID)
		return jid, account, ok
	}

	return "", nil, false
}

func (r *roster) createAccountPopup(jid string, account *account, bt *gdk.EventButton) {
	builder, err := loadBuilderWith("ContactPopupMenu")
	if err != nil {
		// TODO: print error
		return
	}

	obj, err := builder.GetObject("contactMenu")
	if err != nil {
		// TODO: print error
		return
	}

	builder.ConnectSignals(map[string]interface{}{
		"on_remove_contact": func() {
			account.session.Conn.SendIQ("" /* to the server */, "set", xmpp.RosterRequest{
				Item: xmpp.RosterRequestItem{
					Jid:          jid,
					Subscription: "remove",
				},
			})
			r.contacts.Lock()
			r.contacts.m[account].Remove(jid)
			r.contacts.Unlock()
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
		"on_dump_info": func() {
			r.debugPrintRosterFor(account.session.CurrentAccount.Account)
		},
	})

	mn := obj.(*gtk.Menu)
	mn.ShowAll()
	mn.PopupAtMouseCursor(nil, nil, int(bt.Button()), bt.Time())
}

func (r *roster) onButtonPress(view *gtk.TreeView, ev *gdk.Event) bool {
	bt := &gdk.EventButton{ev}
	if bt.Button() == 0x03 {
		jid, account, ok := r.getAccountAndJidFromEvent(bt)
		if ok {
			r.createAccountPopup(jid, account, bt)
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

	if accountID == accountIDTopLevelMarker {
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
	//TODO: handle same account on multiple sessions
	c, ok := r.conversations[to]

	if !ok {
		var err error
		c, err = newConversationWindow(account, to, r.ui)
		if err != nil {
			return nil, err
		}

		r.ui.connectShortcutsChildWindow(c.win)
		r.ui.connectShortcutsConversationWindow(c)
		r.conversations[to] = c
	}

	c.Show()
	return c, nil
}

func (r *roster) displayNameFor(account *account, from string) string {
	r.contacts.RLock()
	l, ok := r.contacts.m[account]
	r.contacts.RUnlock()
	if !ok {
		return from
	}

	p, ok := l.Get(from)
	if ok {
		return p.NameForPresentation()
	}
	return from
}

func (r *roster) presenceUpdated(account *account, from, show, showStatus string, gone bool) {
	c, ok := r.conversations[from]
	if !ok {
		return
	}

	glib.IdleAdd(func() bool {
		c.appendStatus(r.displayNameFor(account, from), time.Now(), show, showStatus, gone)
		return false
	})
}

func (r *roster) messageReceived(account *account, from string, timestamp time.Time, encrypted bool, message []byte) {
	glib.IdleAdd(func() bool {
		conv, err := r.openConversationWindow(account, from)
		if err != nil {
			return false
		}

		conv.appendMessage(r.displayNameFor(account, from), timestamp, encrypted, ui.StripHTML(message), false)
		return false
	})
}

func (r *roster) enableExistingConversationWindows(account *account, enable bool) {
	//TODO: shold lock r.conversations
	for _, convWindow := range r.conversations {
		if enable {
			convWindow.win.Emit("enable")
		} else {
			convWindow.win.Emit("disable")
		}
	}
}

func (r *roster) update(account *account, entries *rosters.List) {
	r.contacts.Lock()
	defer r.contacts.Unlock()
	r.contacts.m[account] = entries
}

func (r *roster) debugPrintRosterFor(nm string) {
	r.contacts.RLock()
	defer r.contacts.RUnlock()

	for account, rs := range r.contacts.m {
		if account.session.CurrentAccount.Is(nm) {
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

func decideStatusGlyphFor(p *rosters.Peer) string {
	if p.PendingSubscribeID != "" {
		return "?"
	}

	if !p.Online {
		return "✘"
	}

	if isAway(p) {
		return "⛔"
	}

	return "✔"
}

func decideColorFor(p *rosters.Peer) string {
	if !p.Online {
		return "#aaaaaa"
	}
	return "#000000"
}

const accountIDTopLevelMarker = "--is-account-top-level--"

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
		fmt.Sprintf("%s%s %s", indent, decideStatusGlyphFor(item), item.NameForPresentation()),
		item.BelongsTo,
		decideColorFor(item),
		"#ffffff",
		nil,
		createTooltipFor(item),
	)
}

func (r *roster) redrawMerged() {
	showOffline := !r.ui.config.ShowOnlyOnline

	r.contacts.RLock()
	defer r.contacts.RUnlock()

	r.toCollapse = nil

	grp := rosters.TopLevelGroup()
	for account, contacts := range r.contacts.m {
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
		r.model.SetValue(pi, indexAccountID, accountIDTopLevelMarker)
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
	parentName := account.session.CurrentAccount.Account
	r.displayGroup(grp, parentIter, accountCounter, showOffline, parentName)

	r.model.SetValue(parentIter, indexParentJid, parentName)
	r.model.SetValue(parentIter, indexAccountID, accountIDTopLevelMarker)
	r.model.SetValue(parentIter, indexWeight, 500)
	r.model.SetValue(parentIter, indexBackgroundColor, "#918caa")
	parentPath, _ := r.model.GetPath(parentIter)
	shouldCollapse, ok := r.isCollapsed[parentName]
	isExpanded := true
	if ok && shouldCollapse {
		isExpanded = false
		r.toCollapse = append(r.toCollapse, parentPath)
	}
	r.model.SetValue(parentIter, indexParentDisplayName, createGroupDisplayName(parentName, accountCounter, isExpanded))
}

func (r *roster) sortedAccounts() []*account {
	var as []*account
	for account := range r.contacts.m {
		as = append(as, account)
	}
	sort.Sort(byAccountNameAlphabetic(as))
	return as
}

func (r *roster) redrawSeparate() {
	showOffline := !r.ui.config.ShowOnlyOnline

	r.contacts.RLock()

	r.toCollapse = nil

	defer r.contacts.RUnlock()

	for _, account := range r.sortedAccounts() {
		r.redrawSeparateAccount(account, r.contacts.m[account], showOffline)
	}

	r.view.ExpandAll()
	for _, path := range r.toCollapse {
		r.view.CollapseRow(path)
	}
}

const disconnectedPageIndex = 0
const spinnerPageIndex = 1
const rosterPageIndex = 2

func (r *roster) redrawIfRosterVisible() {
	if r.widget.GetCurrentPage() == rosterPageIndex {
		r.redraw()
	}
}

func (r *roster) redraw() {
	r.model.Clear()

	if r.ui.shouldViewAccounts() {
		r.redrawSeparate()
	} else {
		r.redrawMerged()
	}

	// We call connected here to make sure we don't display the roster until we have some roster data
	r.connected()
}

func setAll(v *gtk.TreeStore, iter *gtk.TreeIter, values ...interface{}) {
	for i, val := range values {
		if val != nil {
			v.SetValue(iter, i, val)
		}
	}
}
