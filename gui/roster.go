package gui

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/coyim/ui"
	"github.com/twstrike/coyim/xmpp"
)

type Roster struct {
	Window *gtk.ScrolledWindow
	model  *gtk.ListStore
	view   *gtk.TreeView

	contacts map[*Account][]xmpp.RosterEntry

	CheckEncrypted func(to string) bool
	SendMessage    func(to, message string)
	conversations  map[string]*conversationWindow
}

func NewRoster() *Roster {
	w, _ := gtk.ScrolledWindowNew(nil, nil)
	m, _ := gtk.ListStoreNew(
		glib.TYPE_STRING,  // user
		glib.TYPE_POINTER, // *Account
	)
	v, _ := gtk.TreeViewNew()

	r := &Roster{
		Window: w,
		model:  m,
		view:   v,

		conversations: make(map[string]*conversationWindow),
		contacts:      make(map[*Account][]xmpp.RosterEntry),
	}

	r.Window.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_AUTOMATIC)

	r.view.SetHeadersVisible(false)
	if s, err := r.view.GetSelection(); err != nil {
		s.SetMode(gtk.SELECTION_NONE)
	}

	cr, _ := gtk.CellRendererTextNew()
	c, _ := gtk.TreeViewColumnNewWithAttribute("user", cr, "text", 0)
	r.view.AppendColumn(c)

	r.view.SetModel(r.model)
	r.view.Connect("row-activated", r.onActivateBuddy)
	r.Window.Add(r.view)

	//initialize the model
	r.Clear()

	return r
}

func (r *Roster) Clear() {
	glib.IdleAdd(func() {
		//gobj := glib.ToGObject(unsafe.Pointer(r.model.GListStore))
		//gobj.Ref()
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

func (r *Roster) onActivateBuddy(_ *gtk.TreeView, path *gtk.TreePath) {
	iter, err := r.model.GetIter(path)
	if err != nil {
		return
	}

	val, _ := r.model.GetValue(iter, 0)
	to, _ := val.GetString()

	val2, _ := r.model.GetValue(iter, 1)
	account := (*Account)(val2.GetPointer())

	//TODO: change to IDS and fix me
	r.openConversationWindow(account, to)
}

func (r *Roster) openConversationWindow(account *Account, to string) *conversationWindow {
	//TODO: handle same account on multiple sessions
	c, ok := r.conversations[to]

	if !ok {
		c = newConversationWindow(account, to)
		r.conversations[to] = c
	}

	c.Show()
	return c
}

func (r *Roster) MessageReceived(account *Account, from, timestamp string, encrypted bool, message []byte) {
	glib.IdleAdd(func() bool {
		conv := r.openConversationWindow(account, from)
		conv.appendMessage(from, timestamp, encrypted, ui.StripHTML(message))
		return false
	})
}

func (r *Roster) AppendMessageToHistory(to, from, timestamp string, encrypted bool, message []byte) {
	conv, ok := r.conversations[to]
	if !ok {
		return
	}

	glib.IdleAdd(func() bool {
		conv.appendMessage(from, timestamp, encrypted, ui.StripHTML(message))
		return false
	})
}

//TODO: It should have a mutex
func (r *Roster) Update(account *Account, entries []xmpp.RosterEntry) {
	r.contacts[account] = entries
}

func (r *Roster) Redraw() {
	//gobj := glib.ObjectFromNative(unsafe.Pointer(r.model.GListStore))
	//gobj.Ref()

	r.model.TreeModel.Ref()
	r.view.SetModel((*gtk.TreeModel)(nil))

	r.model.Clear()
	for account, contacts := range r.contacts {
		for _, item := range contacts {
			iter := r.model.Append()
			r.model.Set(iter, []int{0, 1}, []interface{}{item.Jid, account})
		}
	}

	r.view.SetModel(r.model)
	r.model.TreeModel.Unref()
	//gobj.Unref()
}
