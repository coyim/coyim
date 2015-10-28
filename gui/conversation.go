package gui

import (
	"fmt"
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
}

func newConversationWindow(account *account, uid string, u *gtkUI) *conversationWindow {
	vars := make(map[string]string)
	vars["$uid"] = uid
	vars["$DevOptions"] = i18n.Local("Developer options")
	vars["$StartOTR"] = i18n.Local("Start encrypted chat")
	vars["$EndOTR"] = i18n.Local("End encrypted chat")
	vars["$VerifyFP"] = i18n.Local("_Verify fingerprint...")
	builder, _ := loadBuilderWith("ConversationDefinition", vars)
	win, _ := builder.GetObject("conversation")
	history, _ := builder.GetObject("history")
	scrollHistory, _ := builder.GetObject("historyScroll")
	messageEntry, _ := builder.GetObject("message")

	conv := &conversationWindow{
		to:            uid,
		account:       account,
		win:           win.(*gtk.Window),
		history:       history.(*gtk.TextView),
		scrollHistory: scrollHistory.(*gtk.ScrolledWindow),
	}

	builder.ConnectSignals(map[string]interface{}{
		"on_send_message_signal": func() {
			entryNode, _ := builder.GetObject("message")
			entry := entryNode.(*gtk.Entry)
			entry.SetEditable(false)
			text, _ := entry.GetText()
			entry.SetText("")
			entry.SetEditable(true)
			sendError := conv.sendMessage(text)
			if sendError != nil {
				fmt.Printf(i18n.Local("Failed to generate OTR message: %s\n"), sendError.Error())
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

	buff, _ := conv.history.GetBuffer()
	ttable, _ := buff.GetTagTable()

	outgoingUser, _ := gtk.TextTagNew("outgoingUser")
	outgoingUser.SetProperty("foreground", "#3465a4")
	ttable.Add(outgoingUser)

	incomingUser, _ := gtk.TextTagNew("incomingUser")
	incomingUser.SetProperty("foreground", "#a40000")
	ttable.Add(incomingUser)

	outgoingText, _ := gtk.TextTagNew("outgoingText")
	outgoingText.SetProperty("foreground", "#555753")
	ttable.Add(outgoingText)

	incomingText, _ := gtk.TextTagNew("incomingText")
	ttable.Add(incomingText)

	statusText, _ := gtk.TextTagNew("statusText")
	statusText.SetProperty("foreground", "#4e9a06")
	ttable.Add(statusText)

	u.displaySettings.control(&conv.history.Container.Widget)
	u.displaySettings.control(&messageEntry.(*gtk.Entry).Widget)

	return conv
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
	encrypted := conversation.IsEncrypted()
	glib.IdleAdd(func() bool {
		conv.appendMessage(conv.account.session.CurrentAccount.Account, time.Now(), encrypted, ui.StripHTML([]byte(message)), true)
		return false
	})

	return nil
}

const timeDisplay = "15:04:05"

func insertWithTag(buff *gtk.TextBuffer, tagName, text string) {
	charCount := buff.GetCharCount()
	buff.InsertAtCursor(text)
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

func (conv *conversationWindow) appendStatusString(text string, timestamp time.Time) {
	glib.IdleAdd(func() bool {
		buff, _ := conv.history.GetBuffer()
		buff.InsertAtCursor("[")
		buff.InsertAtCursor(timestamp.Format(timeDisplay))
		buff.InsertAtCursor("] ")
		insertWithTag(buff, "statusText", text)
		buff.InsertAtCursor("\n")

		return false
	})
}

func (conv *conversationWindow) appendStatus(from string, timestamp time.Time, show, showStatus string, gone bool) {
	conv.appendStatusString(createStatusMessage(from, show, showStatus, gone), timestamp)
}

func (conv *conversationWindow) appendMessage(from string, timestamp time.Time, encrypted bool, message []byte, outgoing bool) {
	glib.IdleAdd(func() bool {
		buff, _ := conv.history.GetBuffer()
		buff.InsertAtCursor("[")
		buff.InsertAtCursor(timestamp.Format(timeDisplay))
		buff.InsertAtCursor("] ")
		insertWithTag(buff, is(outgoing, "outgoingUser", "incomingUser"), from)
		buff.InsertAtCursor(":  ")
		insertWithTag(buff, is(outgoing, "outgoingText", "incomingText"), string(message))
		buff.InsertAtCursor("\n")

		adj := conv.scrollHistory.GetVAdjustment()
		adj.SetValue(adj.GetUpper())

		return false
	})
}
