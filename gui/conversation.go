package gui

import (
	"fmt"

	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/coyim/ui"
	"github.com/twstrike/gotk3/gdk"
	"github.com/twstrike/gotk3/glib"
	"github.com/twstrike/gotk3/gtk"
)

type conversationWindow struct {
	to            string
	account       *Account
	win           *gtk.Window
	history       *gtk.TextView
	scrollHistory *gtk.ScrolledWindow
}

func (conv *conversationWindow) conversationMenu() *gtk.MenuBar {
	menubar, _ := gtk.MenuBarNew()
	conversationMenu, _ := gtk.MenuItemNewWithMnemonic(i18n.Local("Conversation"))
	menubar.Append(conversationMenu)

	submenu, _ := gtk.MenuNew()
	conversationMenu.SetSubmenu(submenu)

	startAKE, _ := gtk.MenuItemNewWithMnemonic(i18n.Local("Start encrypted chat"))
	submenu.Append(startAKE)

	//TODO: enable/disable depending on the conversation's encryption state
	startAKE.Connect("activate", func() {
		//TODO: errors
		err := conv.account.StartEncryptedChatWith(conv.to)
		if err != nil {
			fmt.Printf(i18n.Local("Failed to start the encrypted chat: %s\n"), err.Error())
		}
	})

	//TODO: enable/disable depending on the conversation's encryption state
	endOTR, _ := gtk.MenuItemNewWithMnemonic(i18n.Local("End encrypted chat"))
	submenu.Append(endOTR)

	endOTR.Connect("activate", func() {
		//TODO: errors
		err := conv.account.TerminateConversationWith(conv.to)
		if err != nil {
			fmt.Printf(i18n.Local("Failed to terminate the encrypted chat: %s\n"), err.Error())
		}
	})

	verify, _ := gtk.MenuItemNewWithMnemonic(i18n.Local("_Verify fingerprint..."))
	submenu.Append(verify)

	verify.Connect("activate", func() {
		verifyFingerprintDialog(conv.account, conv.to)
	})

	return menubar
}

func newConversationWindow(account *Account, uid string) *conversationWindow {
	win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	history, _ := gtk.TextViewNew()
	scrollHistory, _ := gtk.ScrolledWindowNew(nil, nil)

	conv := &conversationWindow{
		to:            uid,
		account:       account,
		win:           win,
		history:       history,
		scrollHistory: scrollHistory,
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

	vbox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 1)

	vbox.SetHomogeneous(false)
	vbox.SetSpacing(5)
	vbox.SetBorderWidth(5)

	text, _ := gtk.TextViewNew()
	text.SetWrapMode(gtk.WRAP_WORD)
	text.Connect("key-press-event", func(_ *gtk.TextView, ev *gdk.Event) bool {
		//Send message on ENTER press (without modifier key)
		evKey := gdk.EventKey{ev}
		if (evKey.State()&gdk.GDK_MODIFIER_MASK) == 0 && evKey.KeyVal() == 0xff0d {
			text.SetEditable(false)

			b, _ := text.GetBuffer()
			s := b.GetStartIter()
			e := b.GetEndIter()
			msg, _ := b.GetText(s, e, true)
			b.SetText("")

			text.SetEditable(true)

			err := conv.sendMessage(msg)
			if err != nil {
				fmt.Printf(i18n.Local("Failed to generate OTR message: %s\n"), err.Error())
			}

			return true
		}

		return false
	})

	scroll, _ := gtk.ScrolledWindowNew(nil, nil)
	scroll.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	scroll.Add(text)

	conv.scrollHistory.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	conv.scrollHistory.Add(conv.history)

	vbox.PackStart(conv.conversationMenu(), false, false, 0)
	vbox.PackStart(conv.scrollHistory, true, true, 0)
	vbox.PackStart(scroll, true, true, 0)

	//TODO: provide function to trigger AKE
	encryptedFlag, _ := gtk.ButtonNew()
	vbox.Add(encryptedFlag)

	//TODO this will run undefinitely
	//It should stop at least when the session disconnects
	//It would be better to connect this button to a signal that is emitted when
	//the conversation encrypted state changes
	//This way it would not keep updating the button when the window is not visible
	glib.IdleAdd(func() bool {
		if conv.account.GetConversationWith(conv.to).IsEncrypted() {
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

func (conv *conversationWindow) sendMessage(message string) error {
	err := conv.account.EncryptAndSendTo(conv.to, message)
	if err != nil {
		return err
	}

	//TODO: this should not be in both GUI and roster
	conversation := conv.account.GetConversationWith(conv.to)
	encrypted := conversation.IsEncrypted()
	glib.IdleAdd(func() bool {
		conv.appendMessage("ME", "NOW", encrypted, ui.StripHTML([]byte(message)))
		return false
	})

	return nil
}

func (conv *conversationWindow) appendMessage(from, timestamp string, encrypted bool, message []byte) {
	glib.IdleAdd(func() bool {
		buff, _ := conv.history.GetBuffer()
		buff.InsertAtCursor(timestamp)
		buff.InsertAtCursor(" - ")
		buff.InsertAtCursor(string(message))
		buff.InsertAtCursor("\n")

		adj := conv.scrollHistory.GetVAdjustment()
		adj.SetValue(adj.GetUpper())

		return false
	})
}
