package gui

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/i18n"
	rosters "github.com/twstrike/coyim/roster"
	"github.com/twstrike/coyim/ui"
	"github.com/twstrike/gotk3adapter/gdki"
	"github.com/twstrike/gotk3adapter/gtki"
)

type roster struct {
	widget gtki.ScrolledWindow
	model  gtki.TreeStore
	view   gtki.TreeView

	checkEncrypted func(to string) bool
	sendMessage    func(to, message string)

	isCollapsed map[string]bool
	toCollapse  []gtki.TreePath

	ui          *gtkUI
	deNotify    *desktopNotifications
	actionTimes map[string]time.Time
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
	builder := newBuilder("Roster")

	r := &roster{
		isCollapsed: make(map[string]bool),
		actionTimes: make(map[string]time.Time),
		deNotify:    newDesktopNotifications(),

		ui: u,
	}

	builder.ConnectSignals(map[string]interface{}{
		"on_activate_buddy": r.onActivateBuddy,
	})

	obj := builder.getObj("roster")
	r.widget = obj.(gtki.ScrolledWindow)

	obj = builder.getObj("roster-view")
	r.view = obj.(gtki.TreeView)

	obj = builder.getObj("roster-model")
	r.model = obj.(gtki.TreeStore)

	u.displaySettings.update()

	r.view.Connect("button-press-event", r.onButtonPress)

	return r
}

func (r *roster) getAccount(id string) (*account, bool) {
	return r.ui.accountManager.getAccountByID(id)
}

func getFromModelIter(m gtki.TreeStore, iter gtki.TreeIter, index int) string {
	val, _ := m.GetValue(iter, index)
	v, _ := val.GetString()
	return v
}

func (r *roster) getAccountAndJidFromEvent(bt gdki.EventButton) (jid string, account *account, rowType string, ok bool) {
	x := bt.X()
	y := bt.Y()
	path := g.gtk.TreePathNew()
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

func sortedGroupNames(groups map[string]bool) []string {
	sortedNames := make([]string, 0, len(groups))
	for k := range groups {
		sortedNames = append(sortedNames, k)
	}

	sort.Strings(sortedNames)

	return sortedNames
}

func (r *roster) allGroupNames() []string {
	groups := map[string]bool{}
	for _, contacts := range r.ui.accountManager.getAllContacts() {
		for name := range contacts.GetGroupNames() {
			if groups[name] {
				continue
			}

			groups[name] = true
		}
	}

	return sortedGroupNames(groups)
}

func (r *roster) getGroupNamesFor(a *account) []string {
	groups := map[string]bool{}
	contacts := r.ui.accountManager.getAllContacts()[a]
	for name := range contacts.GetGroupNames() {
		if groups[name] {
			continue
		}

		groups[name] = true
	}

	return sortedGroupNames(groups)
}

func (r *roster) updatePeer(acc *account, jid, nickname string, groups []string, updateRequireEncryption, requireEncryption bool) error {
	peer, ok := r.ui.getPeer(acc, jid)
	if !ok {
		return fmt.Errorf("Could not find peer %s", jid)
	}

	// This updates what is displayed in the roster
	peer.Nickname = nickname
	peer.SetGroups(groups)

	// NOTE: This requires the account to be connected in order to rename peers,
	// which should not be the case. This is one example of why gui.account should
	// own the account config -  and not the session.
	conf := acc.session.GetConfig()
	conf.SavePeerDetails(jid, nickname, groups)
	if updateRequireEncryption {
		conf.UpdateEncryptionRequired(jid, requireEncryption)
	}

	r.ui.SaveConfig()
	doInUIThread(r.redraw)

	return nil
}

func (r *roster) renamePeer(acc *account, jid, nickname string) {
	peer, ok := r.ui.getPeer(acc, jid)
	if !ok {
		return
	}

	// This updates what is displayed in the roster
	peer.Nickname = nickname

	// This saves the nickname to the config file
	// NOTE: This requires the account to be connected in order to rename peers,
	// which should not be the case. This is one example of why gui.account should
	// own the account config -  and not the session.
	//acc.session.GetConfig().RenamePeer(jid, nickname)

	doInUIThread(r.redraw)
	r.ui.SaveConfig()
}

func toArray(groupList gtki.ListStore) []string {
	groups := []string{}

	iter, ok := groupList.GetIterFirst()
	for ok {
		gValue, _ := groupList.GetValue(iter, 0)
		if group, err := gValue.GetString(); err == nil {
			groups = append(groups, group)
		}

		ok = groupList.IterNext(iter)
	}

	return groups
}

func (r *roster) createAccountPeerPopup(jid string, account *account, bt gdki.EventButton) {
	builder := newBuilder("ContactPopupMenu")
	mn := builder.getObj("contactMenu").(gtki.Menu)

	builder.ConnectSignals(map[string]interface{}{
		"on_remove_contact": func() {
			account.session.RemoveContact(jid)
			r.ui.removePeer(account, jid)
			r.redraw()
		},
		"on_edit_contact": func() {
			doInUIThread(func() { r.openEditContactDialog(jid, account) })
		},
		"on_allow_contact_to_see_status": func() {
			account.session.ApprovePresenceSubscription(jid, "" /* generate id */)
		},
		"on_forbid_contact_to_see_status": func() {
			account.session.DenyPresenceSubscription(jid, "" /* generate id */)
		},
		"on_ask_contact_to_see_status": func() {
			account.session.RequestPresenceSubscription(jid, "")
		},
		"on_dump_info": func() {
			r.debugPrintRosterFor(account.session.GetConfig().Account)
		},
	})

	mn.ShowAll()
	mn.PopupAtPointer(bt)
}

func (r *roster) createAccountPopup(jid string, account *account, bt gdki.EventButton) {
	mn := account.createSubmenu()
	if *config.DebugFlag {
		mn.Append(account.createSeparatorItem())
		mn.Append(account.createDumpInfoItem(r))
		mn.Append(account.createXMLConsoleItem(r))
	}
	mn.ShowAll()
	mn.PopupAtPointer(bt)
}

func (r *roster) onButtonPress(view gtki.TreeView, ev gdki.Event) bool {
	bt := g.gdk.EventButtonFrom(ev)
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

func collapseTransform(s string) string {
	res := sha256.Sum256([]byte(s))
	return hex.EncodeToString(res[:])
}

func (r *roster) restoreCollapseStatus() {
	pieces := strings.Split(r.ui.settings.GetCollapsed(), ":")
	for _, p := range pieces {
		if p != "" {
			r.isCollapsed[p] = true
		}
	}
}

func (r *roster) saveCollapseStatus() {
	var vals []string
	for e, v := range r.isCollapsed {
		if v {
			vals = append(vals, e)
		}
	}
	r.ui.settings.SetCollapsed(strings.Join(vals, ":"))
}

func (r *roster) onActivateBuddy(v gtki.TreeView, path gtki.TreePath) {
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
		ix := collapseTransform(jid)
		r.isCollapsed[ix] = !r.isCollapsed[ix]
		r.saveCollapseStatus()
		r.redraw()
		return
	}

	account, ok := r.getAccount(accountID)
	if !ok {
		return
	}

	r.openConversationView(account, jid, true)
}

func (r *roster) openConversationView(account *account, to string, userInitiated bool) (conversationView, error) {
	c, ok := account.getConversationWith(to, r.ui)

	if !ok {
		c = account.createConversationView(to, r.ui)
	}

	c.show(userInitiated)
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
	c, ok := account.getConversationWith(from, r.ui)
	if !ok {
		return
	}

	doInUIThread(func() {
		c.appendStatus(r.displayNameFor(account, from), time.Now(), show, showStatus, gone)
	})
}

const mergeNotificationsThreshold = 7

func (r *roster) lastActionTimeFor(f string) time.Time {
	return r.actionTimes[f]
}

func (r *roster) registerLastActionTimeFor(f string, t time.Time) {
	r.actionTimes[f] = t
}

func (r *roster) maybeNotify(timestamp time.Time, account *account, from, message string) {
	dname := r.displayNameFor(account, from)

	if timestamp.Before(r.lastActionTimeFor(from).Add(time.Duration(mergeNotificationsThreshold) * time.Second)) {
		fmt.Println("Decided to not show notification, since the time is not ready")
		return
	}

	r.registerLastActionTimeFor(from, timestamp)

	err := r.deNotify.show(from, dname, message)
	if err != nil {
		log.Println(err)
	}
}

func (r *roster) messageReceived(account *account, from, resource string, timestamp time.Time, encrypted bool, message []byte) {
	p, ok := r.ui.getPeer(account, from)
	if ok {
		p.LastResource(resource)
	}

	doInUIThread(func() {
		conv, err := r.openConversationView(account, from, false)
		if err != nil {
			return
		}

		sent := sentMessage{
			from:            r.displayNameFor(account, from),
			timestamp:       timestamp,
			isEncrypted:     encrypted,
			isOutgoing:      false,
			strippedMessage: ui.StripSomeHTML(message),
		}
		conv.appendMessage(sent)

		if !conv.isVisible() && r.deNotify != nil {
			r.maybeNotify(timestamp, account, from, string(ui.StripSomeHTML(message)))
		}
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

func isNominallyVisible(p *rosters.Peer, showWaiting bool) bool {
	return (p.Subscription != "none" && p.Subscription != "") || (showWaiting && (p.PendingSubscribeID != "" || p.Asked))
}

func shouldDisplay(p *rosters.Peer, showOffline, showWaiting bool) bool {
	return isNominallyVisible(p, showWaiting) && (showOffline || p.Online || p.Asked)
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
	if p.PendingSubscribeID != "" || p.Asked {
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

func decideColorFor(cs colorSet, p *rosters.Peer) string {
	if !p.Online {
		return cs.rosterPeerOfflineForeground
	}
	return cs.rosterPeerOnlineForeground
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

func (r *roster) addItem(item *rosters.Peer, parentIter gtki.TreeIter, indent string) {
	cs := r.ui.currentColorSet()
	iter := r.model.Append(parentIter)
	potentialExtra := ""
	if item.Asked {
		potentialExtra = i18n.Local(" (waiting for approval)")
	}
	setAll(r.model, iter,
		item.Jid,
		fmt.Sprintf("%s %s%s", indent, item.NameForPresentation(), potentialExtra),
		item.BelongsTo,
		decideColorFor(cs, item),
		cs.rosterPeerBackground,
		nil,
		createTooltipFor(item),
	)

	r.model.SetValue(iter, indexRowType, "peer")
	r.model.SetValue(iter, indexStatusIcon, statusIcons[decideStatusFor(item)].getPixbuf())
}

func (r *roster) redrawMerged() {
	showOffline := !r.ui.config.Display.ShowOnlyOnline
	showWaiting := !r.ui.config.Display.ShowOnlyConfirmed

	r.ui.accountManager.RLock()
	defer r.ui.accountManager.RUnlock()

	r.toCollapse = nil

	grp := rosters.TopLevelGroup()
	for account, contacts := range r.ui.accountManager.getAllContacts() {
		contacts.AddTo(grp, account.session.GroupDelimiter())
	}

	accountCounter := &counter{}
	r.displayGroup(grp, nil, accountCounter, showOffline, showWaiting, "")

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

func (r *roster) sortedPeers(ps []*rosters.Peer) []*rosters.Peer {
	if r.ui.config.Display.SortByStatus {
		sort.Sort(byStatus(ps))
	} else {
		sort.Sort(byNameForPresentation(ps))
	}
	return ps
}

type byNameForPresentation []*rosters.Peer

func (s byNameForPresentation) Len() int { return len(s) }
func (s byNameForPresentation) Less(i, j int) bool {
	return s[i].NameForPresentation() < s[j].NameForPresentation()
}
func (s byNameForPresentation) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func statusValueFor(s string) int {
	switch s {
	case "available":
		return 0
	case "away":
		return 1
	case "extended-away":
		return 2
	case "busy":
		return 3
	case "offline":
		return 4
	case "unknown":
		return 5
	}
	return -1
}

type byStatus []*rosters.Peer

func (s byStatus) Len() int { return len(s) }
func (s byStatus) Less(i, j int) bool {
	return statusValueFor(decideStatusFor(s[i])) < statusValueFor(decideStatusFor(s[j]))
}
func (s byStatus) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (r *roster) displayGroup(g *rosters.Group, parentIter gtki.TreeIter, accountCounter *counter, showOffline, showWaiting bool, accountName string) {
	pi := parentIter
	groupCounter := &counter{}
	groupID := accountName + "//" + g.FullGroupName()

	isEmpty := true
	for _, item := range g.UnsortedPeers() {
		if shouldDisplay(item, showOffline, showWaiting) {
			isEmpty = false
		}
	}

	if g.GroupName != "" && (!isEmpty || r.showEmptyGroups()) {
		pi = r.model.Append(parentIter)
		r.model.SetValue(pi, indexParentJid, groupID)
		r.model.SetValue(pi, indexRowType, "group")
		r.model.SetValue(pi, indexWeight, 500)
		r.model.SetValue(pi, indexBackgroundColor, r.ui.currentColorSet().rosterGroupBackground)
	}

	for _, item := range r.sortedPeers(g.UnsortedPeers()) {
		vs := isNominallyVisible(item, showWaiting)
		o := isOnline(item)
		accountCounter.inc(vs, vs && o)
		groupCounter.inc(vs, vs && o)

		if shouldDisplay(item, showOffline, showWaiting) {
			r.addItem(item, pi, "")
		}
	}

	for _, gr := range g.Groups() {
		r.displayGroup(gr, pi, accountCounter, showOffline, showWaiting, accountName)
	}

	if g.GroupName != "" && (!isEmpty || r.showEmptyGroups()) {
		parentPath, _ := r.model.GetPath(pi)
		shouldCollapse, ok := r.isCollapsed[collapseTransform(groupID)]
		isExpanded := true
		if ok && shouldCollapse {
			isExpanded = false
			r.toCollapse = append(r.toCollapse, parentPath)
		}

		r.model.SetValue(pi, indexParentDisplayName, createGroupDisplayName(g.FullGroupName(), groupCounter, isExpanded))
	}
}

func (r *roster) redrawSeparateAccount(account *account, contacts *rosters.List, showOffline, showWaiting bool) {
	cs := r.ui.currentColorSet()
	parentIter := r.model.Append(nil)

	accountCounter := &counter{}

	grp := contacts.Grouped(account.session.GroupDelimiter())
	parentName := account.session.GetConfig().Account
	r.displayGroup(grp, parentIter, accountCounter, showOffline, showWaiting, parentName)

	r.model.SetValue(parentIter, indexParentJid, parentName)
	r.model.SetValue(parentIter, indexAccountID, account.session.GetConfig().ID())
	r.model.SetValue(parentIter, indexRowType, "account")
	r.model.SetValue(parentIter, indexWeight, 700)

	bgcolor := cs.rosterAccountOnlineBackground
	if account.session.IsDisconnected() {
		bgcolor = cs.rosterAccountOfflineBackground
	}
	r.model.SetValue(parentIter, indexBackgroundColor, bgcolor)

	parentPath, _ := r.model.GetPath(parentIter)
	shouldCollapse, ok := r.isCollapsed[collapseTransform(parentName)]
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

func (r *roster) showEmptyGroups() bool {
	return r.ui.settings.GetShowEmptyGroups()
}

func (r *roster) redrawSeparate() {
	showOffline := !r.ui.config.Display.ShowOnlyOnline
	showWaiting := !r.ui.config.Display.ShowOnlyConfirmed

	r.ui.accountManager.RLock()
	defer r.ui.accountManager.RUnlock()

	r.toCollapse = nil

	for _, account := range r.sortedAccounts() {
		r.redrawSeparateAccount(account, r.ui.accountManager.getContacts(account), showOffline, showWaiting)
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

func setAll(v gtki.TreeStore, iter gtki.TreeIter, values ...interface{}) {
	for i, val := range values {
		if val != nil {
			v.SetValue(iter, i, val)
		}
	}
}
