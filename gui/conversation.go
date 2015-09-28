package gui

import (
	"fmt"
	"unsafe"

	"github.com/twstrike/coyim/session"
	"github.com/twstrike/coyim/ui"
	"github.com/twstrike/go-gtk/gdk"
	"github.com/twstrike/go-gtk/glib"
	"github.com/twstrike/go-gtk/gtk"
	"github.com/twstrike/otr3"
)

type conversationWindow struct {
	to            string
	session       *session.Session
	win           *gtk.Window
	history       *gtk.TextView
	scrollHistory *gtk.ScrolledWindow
}

func newConversationWindow(s *session.Session, uid string) *conversationWindow {
	conv := &conversationWindow{
		to:            uid,
		session:       s,
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
		if s.GetConversationWith(conv.to).IsEncrypted() {
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
	//TODO: this should not be in both GUI and roster
	conversation := conv.session.GetConversationWith(conv.to)

	toSend, err := conversation.Send(otr3.ValidMessage(message))
	if err != nil {
		fmt.Println("Failed to generate OTR message")
		return
	}

	encrypted := conversation.IsEncrypted()
	glib.IdleAdd(func() bool {
		conv.appendMessage("ME", "NOW", encrypted, ui.StripHTML([]byte(message)))
		return false
	})

	for _, m := range toSend {
		//TODO: this should be session.Send(to, message)
		conv.session.Conn.Send(conv.to, string(m))
	}
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
