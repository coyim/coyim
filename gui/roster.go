package gui

import (
	"fmt"
	"sync"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/i18n"
	rosters "github.com/twstrike/coyim/roster"
	"github.com/twstrike/coyim/ui"
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

	isExpanded map[string]bool

	ui *gtkUI
}

func (u *gtkUI) newNotebook() *gtk.Notebook {
	notebook, err := gtk.NotebookNew()
	if err != nil {
		panic("failed")
	}

	notebook.SetShowTabs(false)
	notebook.SetShowBorder(false)
	notebook.PopupDisable()

	vbox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 1)
	vbox.SetHomogeneous(false)
	vbox.SetBorderWidth(3)
	u.displaySettings.unifiedBackgroundColor(&vbox.Container.Widget)

	welcome, _ := gtk.LabelNew(i18n.Local("You are not connected to any account.\nPlease connect to view your online contacts."))

	welcome.SetProperty("margin-start", 5)
	welcome.SetProperty("margin-end", 5)
	welcome.SetMarginTop(7)
	welcome.Show()

	vbox.PackStart(welcome, false, false, 0)

	notebook.AppendPage(vbox, nil)

	vboxSpinner, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 1)
	vboxSpinner.SetHomogeneous(false)
	vboxSpinner.SetBorderWidth(3)
	u.displaySettings.unifiedBackgroundColor(&vboxSpinner.Container.Widget)
	spinner, _ := gtk.SpinnerNew()
	spinner.Start()
	vboxSpinner.PackStart(spinner, true, true, 0)

	notebook.AppendPage(vboxSpinner, nil)

	u.displaySettings.update()

	return notebook
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
)

func (u *gtkUI) newRoster() *roster {
	w, _ := gtk.ScrolledWindowNew(nil, nil)
	m, _ := gtk.TreeStoreNew(
		glib.TYPE_STRING, // jid
		glib.TYPE_STRING, // display name
		glib.TYPE_STRING, // account id
		glib.TYPE_STRING, // color (used to indicate status)
		glib.TYPE_STRING, // background color (used for background of all cell renderers
		glib.TYPE_INT,    // weight of font
	)

	v, _ := gtk.TreeViewNew()

	r := &roster{
		widget: u.newNotebook(),
		model:  m,
		view:   v,

		conversations: make(map[string]*conversationWindow),
		contacts: contacts{
			m: make(map[*account]*rosters.List),
		},
		isExpanded: make(map[string]bool),

		ui: u,
	}

	w.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_AUTOMATIC)

	r.view.SetHeadersVisible(false)
	if s, err := r.view.GetSelection(); err != nil {
		s.SetMode(gtk.SELECTION_NONE)
	}

	cr, _ := gtk.CellRendererTextNew()
	c, _ := gtk.TreeViewColumnNewWithAttribute("name", cr, "text", indexDisplayName)
	c.AddAttribute(cr, "foreground", indexColor)
	c.AddAttribute(cr, "background", indexBackgroundColor)
	c.AddAttribute(cr, "weight", indexWeight)

	r.view.AppendColumn(c)

	r.view.SetShowExpanders(false)
	r.view.SetLevelIndentation(3)

	r.view.SetModel(r.model)
	r.view.Connect("row-activated", r.onActivateBuddy)
	w.Add(r.view)
	w.ShowAll()

	r.widget.AppendPage(w, nil)
	r.disconnected()

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

func (r *roster) onActivateBuddy(_ *gtk.TreeView, path *gtk.TreePath) {
	iter, err := r.model.GetIter(path)
	if err != nil {
		return
	}

	jid := getFromModelIter(r.model, iter, indexJid)
	accountID := getFromModelIter(r.model, iter, indexAccountID)

	if accountID == accountIDTopLevelMarker {
		r.isExpanded[jid] = !r.isExpanded[jid]
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
		return ""
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

//TODO: Does it need the account? Why not having the session?
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
				fmt.Printf("->   #%v\n", item)
			})
		}
	}

	fmt.Printf(" ************************************** \n")
	fmt.Println()
}

//TODO: I believe we can achieve the same with a GtkTreeModelFilter
//See: gtk_tree_model_filter_set_visible_func()
func shouldDisplay(p *rosters.Peer, showOffline bool) bool {
	if showOffline {
		return (p.Subscription != "none" && p.Subscription != "") || p.PendingSubscribeID != ""
	}
	return p.Online && ((p.Subscription != "none" && p.Subscription != "") || p.PendingSubscribeID != "")
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

func createAccountDisplayName(parentName string, count, countOnline int, isExpanded bool) string {
	name := parentName
	if !isExpanded {
		name = fmt.Sprintf("[%s]", name)
	}
	return fmt.Sprintf("%s (%d/%d)", name, countOnline, count)
}

func (r *roster) addItem(item *rosters.Peer, parentIter *gtk.TreeIter, indent string) {
	iter := r.model.Append(parentIter)
	setAll(r.model, iter,
		item.Jid,
		fmt.Sprintf("%s%s %s", indent, decideStatusGlyphFor(item), item.NameForPresentation()),
		item.BelongsTo,
		decideColorFor(item),
		"#ffffff",
	)
}

func (r *roster) redrawMerged() {
	showOffline := !r.ui.config.ShowOnlyOnline

	r.contacts.RLock()
	allContacts := make([]*rosters.List, 0, len(r.contacts.m))
	for _, contacts := range r.contacts.m {
		allContacts = append(allContacts, contacts)
	}
	r.contacts.RUnlock()

	rosters.IterAll(func(_ int, item *rosters.Peer) {
		if shouldDisplay(item, showOffline) {
			r.addItem(item, nil, "")
		}
	}, allContacts...)
}

func (r *roster) redrawSeparate() {
	var parentIter *gtk.TreeIter
	showOffline := !r.ui.config.ShowOnlyOnline

	r.contacts.RLock()
	defer r.contacts.RUnlock()
	for account, contacts := range r.contacts.m {
		parentIter = r.model.Append(nil)

		count := 0
		countOnline := 0

		contacts.Iter(func(_ int, item *rosters.Peer) {
			if shouldDisplay(item, showOffline) {
				count++
				if isOnline(item) {
					countOnline++
				}
				r.addItem(item, parentIter, "   ")
			}
		})

		parentName := account.session.CurrentAccount.Account
		r.model.SetValue(parentIter, indexParentJid, parentName)
		r.model.SetValue(parentIter, indexAccountID, accountIDTopLevelMarker)
		r.model.SetValue(parentIter, indexWeight, 500)
		r.model.SetValue(parentIter, indexBackgroundColor, "#918caa")
		parentPath, _ := r.model.GetPath(parentIter)
		isExpanded, ok := r.isExpanded[parentName]
		if !ok {
			isExpanded = true
			r.isExpanded[parentName] = isExpanded
		}
		r.model.SetValue(parentIter, indexParentDisplayName, createAccountDisplayName(parentName, count, countOnline, isExpanded))
		if isExpanded {
			r.view.ExpandRow(parentPath, false)
		} else {
			r.view.CollapseRow(parentPath)
		}
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
		v.SetValue(iter, i, val)
	}
}
