package gui

import (
	"unsafe"

	"github.com/twstrike/coyim/session"
	"github.com/twstrike/coyim/ui"
	"github.com/twstrike/coyim/xmpp"
	"github.com/twstrike/go-gtk/glib"
	"github.com/twstrike/go-gtk/gtk"
)

type Roster struct {
	Window *gtk.ScrolledWindow
	model  *gtk.ListStore
	view   *gtk.TreeView

	contacts map[*session.Session][]xmpp.RosterEntry

	CheckEncrypted func(to string) bool
	SendMessage    func(to, message string)
	conversations  map[string]*conversationWindow
}

func NewRoster() *Roster {
	r := &Roster{
		Window: gtk.NewScrolledWindow(nil, nil),

		model: gtk.NewListStore(
			gtk.TYPE_STRING,  // user
			gtk.TYPE_POINTER, // *Session
		),
		view: gtk.NewTreeView(),

		conversations: make(map[string]*conversationWindow),
		contacts:      make(map[*session.Session][]xmpp.RosterEntry),
	}

	r.Window.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_AUTOMATIC)

	r.view.SetHeadersVisible(false)
	r.view.GetSelection().SetMode(gtk.SELECTION_NONE)

	r.view.AppendColumn(
		gtk.NewTreeViewColumnWithAttributes("user",
			gtk.NewCellRendererText(), "text", 0),
	)

	r.view.SetModel(r.model)
	r.view.Connect("row-activated", r.onActivateBuddy)
	r.Window.Add(r.view)

	//initialize the model
	r.Clear()

	return r
}

func (r *Roster) Clear() {
	glib.IdleAdd(func() bool {
		gobj := glib.ObjectFromNative(unsafe.Pointer(r.model.GListStore))

		gobj.Ref()
		r.view.SetModel(nil)
		r.model.Clear()

		//TODO: Replace by something better
		iter := &gtk.TreeIter{}
		r.model.Append(iter)
		r.model.Set(iter,
			0, "Disconnected.\nPlease connect from pref. menu",
		)

		r.view.SetModel(r.model)
		gobj.Unref()
		return false
	})
}

func (r *Roster) onActivateBuddy(ctx *glib.CallbackContext) {
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

	val2 := &glib.GValue{}
	r.model.GetValue(iter, 1, val2)
	s := (*session.Session)(val2.GetPointer())

	r.openConversationWindow(s, to)

	//TODO call g_value_unset() but the lib doesnt provide
}

func (r *Roster) openConversationWindow(s *session.Session, to string) *conversationWindow {
	//TODO: handle same account on multiple sessions
	c, ok := r.conversations[to]

	if !ok {
		c = newConversationWindow(s, to)
		r.conversations[to] = c
	}

	c.Show()
	return c
}

func (r *Roster) MessageReceived(from, timestamp string, encrypted bool, message []byte) {
	glib.IdleAdd(func() bool {
		//TODO: fix me
		conv := r.openConversationWindow(nil, from)
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

func (r *Roster) Update(s *session.Session, entries []xmpp.RosterEntry) {
	r.contacts[s] = entries
}

func (r *Roster) Redraw() {
	gobj := glib.ObjectFromNative(unsafe.Pointer(r.model.GListStore))

	gobj.Ref()
	r.view.SetModel(nil)

	r.model.Clear()
	iter := &gtk.TreeIter{}
	for s, contacts := range r.contacts {
		for _, item := range contacts {
			r.model.Append(iter)

			//state, ok := s.knownStates[item.Jid]
			// Subscription, knownState
			r.model.Set(iter,
				0, item.Jid,
				1, s,
			)
		}
	}

	r.view.SetModel(r.model)
	gobj.Unref()
}
