package gui

import (
	"fmt"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/i18n"
	rosters "github.com/twstrike/coyim/roster"
	"github.com/twstrike/coyim/ui"
)

type roster struct {
	widget *gtk.Notebook

	model *gtk.ListStore
	view  *gtk.TreeView

	contacts map[*account]*rosters.List

	checkEncrypted func(to string) bool
	sendMessage    func(to, message string)
	conversations  map[string]*conversationWindow

	ui *gtkUI
}

func newNotebook() *gtk.Notebook {
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

	welcome, _ := gtk.LabelNew(i18n.Local("You are not connected to any account.\nPlease connect to view your online contacts."))

	welcome.SetProperty("margin-start", 5)
	welcome.SetProperty("margin-end", 5)
	welcome.SetMarginTop(7)
	welcome.Show()

	vbox.PackStart(welcome, false, false, 0)

	notebook.AppendPage(vbox, nil)

	spinner, _ := gtk.SpinnerNew()
	spinner.Start()
	notebook.AppendPage(spinner, nil)

	return notebook
}

func (u *gtkUI) newRoster() *roster {
	w, _ := gtk.ScrolledWindowNew(nil, nil)
	m, _ := gtk.ListStoreNew(
		glib.TYPE_STRING, // jid
		glib.TYPE_STRING, // display name
		glib.TYPE_STRING, // account id
		glib.TYPE_STRING, // account status
		glib.TYPE_STRING, // color (used to indicate status)
	)
	v, _ := gtk.TreeViewNew()

	r := &roster{
		widget: newNotebook(),
		model:  m,
		view:   v,

		conversations: make(map[string]*conversationWindow),
		contacts:      make(map[*account]*rosters.List),

		ui: u,
	}

	w.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_AUTOMATIC)

	r.view.SetHeadersVisible(false)
	if s, err := r.view.GetSelection(); err != nil {
		s.SetMode(gtk.SELECTION_NONE)
	}

	cr, _ := gtk.CellRendererTextNew()
	c, _ := gtk.TreeViewColumnNewWithAttribute("name", cr, "text", 1)
	c.AddAttribute(cr, "foreground", 4)

	cr2, _ := gtk.CellRendererTextNew()
	c2, _ := gtk.TreeViewColumnNewWithAttribute("status", cr2, "text", 3)

	r.view.AppendColumn(c2)
	r.view.AppendColumn(c)

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
		r.widget.SetCurrentPage(2)
	})
}

func (r *roster) connecting() {
	glib.IdleAdd(func() {
		r.widget.SetCurrentPage(1)
	})
}

func (r *roster) disconnected() {
	//TODO: should it destroy all conversations?
	glib.IdleAdd(func() {
		r.widget.SetCurrentPage(0)
	})
}

//TODO: move somewhere else
func (r *roster) getAccount(id string) (*account, bool) {
	for account := range r.contacts {
		if account.session.CurrentAccount.ID() == id {
			return account, true
		}
	}

	return nil, false
}

func (r *roster) onActivateBuddy(_ *gtk.TreeView, path *gtk.TreePath) {
	iter, err := r.model.GetIter(path)
	if err != nil {
		return
	}

	val, _ := r.model.GetValue(iter, 0)
	to, _ := val.GetString()

	val2, _ := r.model.GetValue(iter, 2)
	accountName, _ := val2.GetString()
	account, ok := r.getAccount(accountName)

	if !ok {
		return
	}

	r.openConversationWindow(account, to)
}

func (r *roster) openConversationWindow(account *account, to string) *conversationWindow {
	//TODO: handle same account on multiple sessions
	c, ok := r.conversations[to]

	if !ok {
		c = newConversationWindow(account, to)
		r.ui.connectShortcutsChildWindow(c.win)
		r.conversations[to] = c
	}

	c.Show()
	return c
}

func (r *roster) displayNameFor(account *account, from string) string {
	p, ok := r.contacts[account].Get(from)
	if ok {
		return p.NameForPresentation()
	}
	return from
}

func (r *roster) presenceUpdated(account *account, from, show, showStatus string, gone bool) {
	_, ok := r.conversations[from]
	if ok {
		glib.IdleAdd(func() bool {
			conv := r.openConversationWindow(account, from)
			conv.appendStatus(r.displayNameFor(account, from), time.Now(), show, showStatus, gone)
			return false
		})
	}
}

func (r *roster) messageReceived(account *account, from string, timestamp time.Time, encrypted bool, message []byte) {
	glib.IdleAdd(func() bool {
		conv := r.openConversationWindow(account, from)
		conv.appendMessage(r.displayNameFor(account, from), timestamp, encrypted, ui.StripHTML(message), false)
		return false
	})
}

//TODO: It should have a mutex
//Does it need the account? Why not having the session?
func (r *roster) update(account *account, entries *rosters.List) {
	r.contacts[account] = entries
}

func (r *roster) debugPrintRosterFor(nm string) {
	for account, rs := range r.contacts {
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
func shouldDisplay(p *rosters.Peer) bool {
	return (p.Subscription != "none" && p.Subscription != "") || p.PendingSubscribeID != ""
}

func isAway(p *rosters.Peer) bool {
	switch p.Status {
	case "dnd", "xa", "away":
		return true
	}
	return false
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

func (r *roster) redraw() {
	r.model.TreeModel.Ref()
	r.view.SetModel((*gtk.TreeModel)(nil))

	r.model.Clear()
	for account, contacts := range r.contacts {
		contacts.Iter(func(_ int, item *rosters.Peer) {
			if shouldDisplay(item) {
				iter := r.model.Append()
				r.model.Set(iter, []int{0, 1, 2, 3, 4}, []interface{}{
					item.Jid,
					item.NameForPresentation(),
					account.session.CurrentAccount.ID(),
					decideStatusGlyphFor(item),
					decideColorFor(item),
				})
			}
		})
	}

	r.view.SetModel(r.model)
	r.model.TreeModel.Unref()

	// We call connected here to make sure we don't display the roster until we have some roster data
	r.connected()
}
