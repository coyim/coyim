package gui

import (
	"unsafe"

	"github.com/twstrike/go-gtk/gdk"
	"github.com/twstrike/go-gtk/glib"
	"github.com/twstrike/go-gtk/gtk"
)

type conversationWindow struct {
	to            string
	win           *gtk.Window
	history       *gtk.TextView
	scrollHistory *gtk.ScrolledWindow
	roster        *Roster
}

func newConversationWindow(r *Roster, uid string) *conversationWindow {
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

	//TODO: provide function to trigger AKE
	encryptedFlag := gtk.NewButton()
	vbox.Add(encryptedFlag)
	glib.IdleAdd(func() bool {
		if conv.roster.CheckEncrypted(conv.to) {
			encryptedFlag.SetLabel("encrypted")
		} else {
			encryptedFlag.SetLabel("unencrypted")
		}
		return true
	})

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
	conv.roster.SendMessage(conv.to, message)
}

func (conv *conversationWindow) appendMessage(from, timestamp string, encrypted bool, message []byte) {
	glib.IdleAdd(func() bool {
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
