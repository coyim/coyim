package gui

import (
	"fmt"
	"sync"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/coyim/ui"
)

type conversationWindow struct {
	to            string
	account       *account
	win           *gtk.Window
	history       *gtk.TextView
	scrollHistory *gtk.ScrolledWindow
	sync.Mutex
}

type tags struct {
	table *gtk.TextTagTable
}

func (u *gtkUI) getTags() *tags {
	if u.tags == nil {
		u.tags = newTags()
	}
	return u.tags
}

func newTags() *tags {
	t := new(tags)

	t.table, _ = gtk.TextTagTableNew()

	outgoingUser, _ := gtk.TextTagNew("outgoingUser")
	outgoingUser.SetProperty("foreground", "#3465a4")

	incomingUser, _ := gtk.TextTagNew("incomingUser")
	incomingUser.SetProperty("foreground", "#a40000")

	outgoingText, _ := gtk.TextTagNew("outgoingText")
	outgoingText.SetProperty("foreground", "#555753")

	incomingText, _ := gtk.TextTagNew("incomingText")

	statusText, _ := gtk.TextTagNew("statusText")
	statusText.SetProperty("foreground", "#4e9a06")

	t.table.Add(outgoingUser)
	t.table.Add(incomingUser)
	t.table.Add(outgoingText)
	t.table.Add(incomingText)
	t.table.Add(statusText)

	return t
}

func (t *tags) createTextBuffer() *gtk.TextBuffer {
	buf, _ := gtk.TextBufferNew(t.table)
	return buf
}

func newConversationWindow(account *account, uid string, u *gtkUI) (*conversationWindow, error) {
	vars := make(map[string]string)
	vars["$uid"] = uid
	vars["$DevOptions"] = i18n.Local("Developer options")
	vars["$StartOTR"] = i18n.Local("Start encrypted chat")
	vars["$EndOTR"] = i18n.Local("End encrypted chat")
	vars["$VerifyFP"] = i18n.Local("_Verify fingerprint...")

	var win, history, scrollHistory, messageEntry glib.IObject

	builder, err := loadBuilderWith("ConversationDefinition", vars)
	if err != nil {
		return nil, err
	}

	win, err = builder.GetObject("conversation")
	if err != nil {
		return nil, err
	}

	history, err = builder.GetObject("history")
	if err != nil {
		return nil, err
	}

	scrollHistory, err = builder.GetObject("historyScroll")
	if err != nil {
		return nil, err
	}

	messageEntry, err = builder.GetObject("message")
	if err != nil {
		return nil, err
	}

	conv := &conversationWindow{
		to:            uid,
		account:       account,
		win:           win.(*gtk.Window),
		history:       history.(*gtk.TextView),
		scrollHistory: scrollHistory.(*gtk.ScrolledWindow),
	}

	builder.ConnectSignals(map[string]interface{}{
		"on_send_message_signal": func() {
			entry := messageEntry.(*gtk.Entry)
			entry.SetEditable(false)
			text, _ := entry.GetText()
			entry.SetText("")
			entry.SetEditable(true)
			if text != "" {
				sendError := conv.sendMessage(text)
				if sendError != nil {
					fmt.Printf(i18n.Local("Failed to generate OTR message: %s\n"), sendError.Error())
				}
			}
			entry.GrabFocus()
		},
		// TODO: basically I think this whole menu should be rethought. It's useful for us to have during development
		"on_start_otr_signal": func() {
			//TODO: enable/disable depending on the conversation's encryption state
			//TODO: errors
			err := conv.account.session.StartEncryptedChatWith(conv.to)
			if err != nil {
				fmt.Printf(i18n.Local("Failed to start the encrypted chat: %s\n"), err.Error())
			}
		},
		"on_end_otr_signal": func() {
			//TODO: errors
			//TODO: enable/disable depending on the conversation's encryption state
			err := conv.account.session.TerminateConversationWith(conv.to)
			if err != nil {
				fmt.Printf(i18n.Local("Failed to terminate the encrypted chat: %s\n"), err.Error())
			}
		},
		"on_verify_fp_signal": func() {
			verifyFingerprintDialog(conv.account, conv.to, conv.win)
		},
	})

	// Unlike the GTK version, this is not supposed to be used as a callback but
	// it attaches the callback to the widget
	conv.win.HideOnDelete()

	conv.history.SetBuffer(u.getTags().createTextBuffer())

	u.displaySettings.control(&conv.history.Container.Widget)
	u.displaySettings.control(&messageEntry.(*gtk.Entry).Widget)

	return conv, nil
}

func (conv *conversationWindow) Hide() {
	conv.win.Hide()
}

func (conv *conversationWindow) Show() {
	conv.win.ShowAll()
}

func (conv *conversationWindow) sendMessage(message string) error {
	err := conv.account.session.EncryptAndSendTo(conv.to, message)
	if err != nil {
		return err
	}

	//TODO: this should not be in both GUI and roster
	conversation := conv.account.session.GetConversationWith(conv.to)
	conv.appendMessage(conv.account.session.CurrentAccount.Account, time.Now(), conversation.IsEncrypted(), ui.StripHTML([]byte(message)), true)

	return nil
}

const timeDisplay = "15:04:05"

// Expects to be called from the GUI thread.
// Expects to be called when conv is already locked
func insertAtEnd(buff *gtk.TextBuffer, text string) {
	buff.Insert(buff.GetEndIter(), text)
}

// Expects to be called from the GUI thread.
// Expects to be called when conv is already locked
func insertWithTag(buff *gtk.TextBuffer, tagName, text string) {
	charCount := buff.GetCharCount()
	insertAtEnd(buff, text)
	oldEnd := buff.GetIterAtOffset(charCount)
	newEnd := buff.GetEndIter()
	buff.ApplyTagByName(tagName, oldEnd, newEnd)
}

func is(v bool, left, right string) string {
	if v {
		return left
	}
	return right
}

func showForDisplay(show string, gone bool) string {
	switch show {
	case "", "available", "online":
		if gone {
			return ""
		}
		return i18n.Local("Available")
	case "xa":
		return i18n.Local("Not Available")
	case "away":
		return i18n.Local("Away")
	case "dnd":
		return i18n.Local("Busy")
	case "chat":
		return i18n.Local("Free for Chat")
	case "invisible":
		return i18n.Local("Invisible")
	}
	return show
}

func onlineStatus(show, showStatus string) string {
	sshow := showForDisplay(show, false)
	if sshow != "" {
		return sshow + showStatusForDisplay(showStatus)
	}
	return ""
}

func showStatusForDisplay(showStatus string) string {
	if showStatus != "" {
		return " (" + showStatus + ")"
	}
	return ""
}

func extraOfflineStatus(show, showStatus string) string {
	sshow := showForDisplay(show, true)
	if sshow == "" {
		return showStatusForDisplay(showStatus)
	}

	if showStatus != "" {
		return " (" + sshow + ": " + showStatus + ")"
	}
	return " (" + sshow + ")"
}

func createStatusMessage(from string, show, showStatus string, gone bool) string {
	tail := ""
	if gone {
		tail = i18n.Local("Offline") + extraOfflineStatus(show, showStatus)
	} else {
		tail = onlineStatus(show, showStatus)
	}

	if tail != "" {
		return from + i18n.Local(" is now ") + tail
	}
	return ""
}

func (conv *conversationWindow) scrollToBottom() {
	adj := conv.scrollHistory.GetVAdjustment()
	adj.SetValue(adj.GetUpper())
}

// ADDS_TO_GUI_THREAD
// LOCKS_CONV
func (conv *conversationWindow) appendStatusString(text string, timestamp time.Time) {
	glib.IdleAdd(func() bool {
		conv.Lock()
		defer conv.Unlock()

		buff, _ := conv.history.GetBuffer()
		insertAtEnd(buff, "[")
		insertAtEnd(buff, timestamp.Format(timeDisplay))
		insertAtEnd(buff, "]")
		insertWithTag(buff, "statusText", text)
		insertAtEnd(buff, "\n")

		conv.scrollToBottom()

		return false
	})
}

func (conv *conversationWindow) appendStatus(from string, timestamp time.Time, show, showStatus string, gone bool) {
	conv.appendStatusString(createStatusMessage(from, show, showStatus, gone), timestamp)
}

// ADDS_TO_GUI_THREAD
// LOCKS_CONV
func (conv *conversationWindow) appendMessage(from string, timestamp time.Time, encrypted bool, message []byte, outgoing bool) {
	glib.IdleAdd(func() bool {
		conv.Lock()
		defer conv.Unlock()

		buff, _ := conv.history.GetBuffer()
		insertAtEnd(buff, "[")
		insertAtEnd(buff, timestamp.Format(timeDisplay))
		insertAtEnd(buff, "] ")
		insertWithTag(buff, is(outgoing, "outgoingUser", "incomingUser"), from)
		insertAtEnd(buff, ":  ")
		insertWithTag(buff, is(outgoing, "outgoingText", "incomingText"), string(message))
		insertAtEnd(buff, "\n")

		conv.scrollToBottom()

		return false
	})
}
