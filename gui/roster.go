package gui

import (
	"fmt"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/i18n"
	rosters "github.com/twstrike/coyim/roster"
	"github.com/twstrike/coyim/ui"
	"github.com/twstrike/coyim/xmpp"
)

type roster struct {
	window *gtk.ScrolledWindow
	model  *gtk.ListStore
	view   *gtk.TreeView

	contacts map[*account]*rosters.List

	checkEncrypted func(to string) bool
	sendMessage    func(to, message string)
	conversations  map[string]*conversationWindow
}

func newRoster() *roster {
	w, _ := gtk.ScrolledWindowNew(nil, nil)
	m, _ := gtk.ListStoreNew(
		glib.TYPE_STRING, // jid
		glib.TYPE_STRING, // display name
		glib.TYPE_STRING, // account id
	)
	v, _ := gtk.TreeViewNew()

	r := &roster{
		window: w,
		model:  m,
		view:   v,

		conversations: make(map[string]*conversationWindow),
		contacts:      make(map[*account]*rosters.List),
	}

	r.window.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_AUTOMATIC)

	r.view.SetHeadersVisible(false)
	if s, err := r.view.GetSelection(); err != nil {
		s.SetMode(gtk.SELECTION_NONE)
	}

	cr, _ := gtk.CellRendererTextNew()
	c, _ := gtk.TreeViewColumnNewWithAttribute("name", cr, "text", 1)
	r.view.AppendColumn(c)

	r.view.SetModel(r.model)
	r.view.Connect("row-activated", r.onActivateBuddy)
	r.window.Add(r.view)

	//initialize the model
	r.clear()

	return r
}

func (r *roster) clear() {
	glib.IdleAdd(func() {
		r.model.TreeModel.Ref()

		r.view.SetModel((*gtk.TreeModel)(nil))
		r.model.Clear()

		iter := r.model.Append()
		r.model.SetValue(iter,
			0, i18n.Local("Disconnected.\nPlease connect from pref. menu"),
		)

		r.view.SetModel(r.model)
		r.model.TreeModel.Unref()
	})
}

//TODO: move somewhere else
func (r *roster) getAccount(id string) (*account, bool) {
	for account := range r.contacts {
		if account.ID() == id {
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

func (r *roster) messageReceived(account *account, from string, timestamp time.Time, encrypted bool, message []byte) {
	glib.IdleAdd(func() bool {
		conv := r.openConversationWindow(account, from)
		conv.appendMessage(r.displayNameFor(account, from), timestamp, encrypted, ui.StripHTML(message))
		return false
	})
}

//TODO: It should have a mutex
func (r *roster) update(account *account, entries *rosters.List) {
	r.contacts[account] = entries
}

func (r *roster) debugPrintRosterFor(nm string) {
	nnm := xmpp.RemoveResourceFromJid(nm)
	fmt.Printf(" ^^^ Roster for: %s ^^^ \n", nnm)

	for account, rs := range r.contacts {
		if account.Config.Account == nnm {
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
	return p.Subscription != "none" && p.Subscription != ""
}

func (r *roster) redraw() {
	r.model.TreeModel.Ref()
	r.view.SetModel((*gtk.TreeModel)(nil))

	r.model.Clear()
	for account, contacts := range r.contacts {
		contacts.Iter(func(_ int, item *rosters.Peer) {
			if shouldDisplay(item) {
				iter := r.model.Append()
				r.model.Set(iter, []int{0, 1, 2}, []interface{}{
					item.Jid,
					item.NameForPresentation(),
					account.ID(),
				})
			}
		})
	}

	r.view.SetModel(r.model)
	r.model.TreeModel.Unref()
}
