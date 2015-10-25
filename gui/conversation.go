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

func newConversationWindow(account *account, uid string) *conversationWindow {
	des := `
	<interface>
	  <object class="GtkWindow" id="conversation">
	    <property name="window-position">0</property>
	    <property name="default-height">500</property>
	    <property name="default-width">400</property>
	    <property name="destroy-with-parent">true</property>
	    <property name="title">` + uid + `</property>

	    <child>
	      <object class="GtkVBox">
                <property name="homogeneous">false</property>
	      	<child>
	          <object class="GtkMenuBar" id="menubar">
	            <child>
	              <object class="GtkMenuItem" id="conversationMenu">
	                <property name="label">` + i18n.Local("Developer options") + `</property>
	                <child type="submenu">
	                  <object class="GtkMenu">
	                    <child>
	                      <object class="GtkMenuItem" id="startOTRMenu">
	                        <property name="label">` + i18n.Local("Start encrypted chat") + `</property>
		                <signal name="activate" handler="on_start_otr_signal" />
	            	      </object>
	            	    </child>
	                    <child>
	                      <object class="GtkMenuItem" id="endOTRMenu">
	                        <property name="label">` + i18n.Local("End encrypted chat") + `</property>
		                <signal name="activate" handler="on_end_otr_signal" />
	            	      </object>
	            	    </child>
	                    <child>
	                      <object class="GtkMenuItem" id="verifyFingerMenu">
	                        <property name="label">` + i18n.Local("_Verify fingerprint...") + `</property>
		                <signal name="activate" handler="on_verify_fp_signal" />
	            	      </object>
	            	    </child>
	                  </object>
	                </child>
	              </object>
	            </child>
	          </object>

	          <packing>
	            <property name="expand">false</property>
	            <property name="fill">true</property>
	            <property name="position">0</property>
	          </packing>
		</child>

		<child>
	            <object class="GtkScrolledWindow" id="historyScroll">
	              <property name="vscrollbar-policy">GTK_POLICY_AUTOMATIC</property>
	              <property name="hscrollbar-policy">GTK_POLICY_AUTOMATIC</property>
	              <child>
	                <object class="GtkTextView" id="history">
		          <property name="visible">true</property>
		          <property name="wrap-mode">1</property>
		          <property name="editable">false</property>
		          <property name="cursor-visible">false</property>
	                </object>

	              </child>

	            </object>

	          <packing>
	            <property name="expand">true</property>
	            <property name="fill">true</property>
	            <property name="position">1</property>
	          </packing>
		</child>

	        <child>
	          <object class="GtkEntry" id="message">
	            <property name="has-focus">true</property>
		    <signal name="activate" handler="on_send_message_signal" />
	          </object>
	          <packing>
	            <property name="expand">false</property>
	            <property name="fill">true</property>
	            <property name="position">2</property>
	          </packing>
	        </child>
	      </object>
	    </child>

	  </object>
	</interface>
	`
	builder, _ := gtk.BuilderNew()
	addError := builder.AddFromString(des)
	if addError != nil {
		fmt.Printf("Failed to add the UI XML: %s", addError.Error())
	}

	win, _ := builder.GetObject("conversation")
	history, _ := builder.GetObject("history")
	scrollHistory, _ := builder.GetObject("historyScroll")

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
			verifyFingerprintDialog(conv.account, conv.to)
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
