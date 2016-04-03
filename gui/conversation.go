package gui

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gdki"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glibi"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtki"
	"github.com/twstrike/coyim/client"
	"github.com/twstrike/coyim/i18n"
	rosters "github.com/twstrike/coyim/roster"
	"github.com/twstrike/coyim/session/access"
	"github.com/twstrike/coyim/ui"
)

var (
	enableWindow  glibi.Signal
	disableWindow glibi.Signal
)

type conversationView interface {
	showIdentityVerificationWarning(u *gtkUI)
	updateSecurityWarning()
	show(userInitiated bool)
	appendStatus(from string, timestamp time.Time, show, showStatus string, gone bool)
	appendMessage(from string, timestamp time.Time, encrypted bool, message []byte, outgoing bool)
	displayNotification(notification string)
	displayNotificationVerifiedOrNot(notificationV, notificationNV string)
	setEnabled(enabled bool)
	isVisible() bool
}

type conversationWindow struct {
	*conversationPane
	win       gtki.Window
	parentWin gtki.Window
}

type conversationPane struct {
	to                 string
	account            *account
	widget             gtki.Box
	menubar            gtki.MenuBar
	entry              gtki.TextView
	entryScroll        gtki.ScrolledWindow
	history            gtki.TextView
	scrollHistory      gtki.ScrolledWindow
	notificationArea   gtki.Box
	securityWarning    gtki.InfoBar
	fingerprintWarning gtki.InfoBar
	// The window to set dialogs transient for
	transientParent gtki.Window
	sync.Mutex
	marks           []*timedMark
	hidden          bool
	afterNewMessage func()
	currentPeer     func() (*rosters.Peer, bool)
}

type tags struct {
	table gtki.TextTagTable
}

func (u *gtkUI) getTags() *tags {
	if u.tags == nil {
		u.tags = u.newTags()
	}
	return u.tags
}

func (u *gtkUI) newTags() *tags {
	cs := u.currentColorSet()
	t := new(tags)

	t.table, _ = g.gtk.TextTagTableNew()

	outgoingUser, _ := g.gtk.TextTagNew("outgoingUser")
	outgoingUser.SetProperty("foreground", cs.conversationOutgoingUserForeground)

	incomingUser, _ := g.gtk.TextTagNew("incomingUser")
	incomingUser.SetProperty("foreground", cs.conversationIncomingUserForeground)

	outgoingText, _ := g.gtk.TextTagNew("outgoingText")
	outgoingText.SetProperty("foreground", cs.conversationOutgoingTextForeground)

	incomingText, _ := g.gtk.TextTagNew("incomingText")
	outgoingText.SetProperty("foreground", cs.conversationIncomingTextForeground)

	statusText, _ := g.gtk.TextTagNew("statusText")
	statusText.SetProperty("foreground", cs.conversationStatusTextForeground)

	t.table.Add(outgoingUser)
	t.table.Add(incomingUser)
	t.table.Add(outgoingText)
	t.table.Add(incomingText)
	t.table.Add(statusText)

	return t
}

func (t *tags) createTextBuffer() gtki.TextBuffer {
	buf, _ := g.gtk.TextBufferNew(t.table)
	return buf
}

func getTextBufferFrom(e gtki.TextView) gtki.TextBuffer {
	tb, _ := e.GetBuffer()
	return tb
}

func getTextFrom(e gtki.TextView) string {
	tb := getTextBufferFrom(e)
	return tb.GetText(tb.GetStartIter(), tb.GetEndIter(), false)
}

func insertEnter(e gtki.TextView) {
	tb := getTextBufferFrom(e)
	tb.Insert(tb.GetEndIter(), "\n")
}

func clearIn(e gtki.TextView) {
	tb := getTextBufferFrom(e)
	tb.Delete(tb.GetStartIter(), tb.GetEndIter())
}

func (conv *conversationWindow) isVisible() bool {
	return conv.win.HasToplevelFocus()

}

func (conv *conversationPane) onSendMessageSignal() {
	conv.entry.SetEditable(false)
	text := getTextFrom(conv.entry)
	clearIn(conv.entry)
	conv.entry.SetEditable(true)
	if text != "" {
		sendError := conv.sendMessage(text)

		if sendError != nil {
			log.Printf(i18n.Local("Failed to generate OTR message: %s\n"), sendError.Error())
		}
	}
	conv.entry.GrabFocus()
}

func (conv *conversationPane) currentResource() string {
	resource := ""
	conv.withCurrentPeer(func(p *rosters.Peer) {
		resource = p.ResourceToUse()
	})
	return resource
}

func (conv *conversationPane) onStartOtrSignal() {
	//TODO: enable/disable depending on the conversation's encryption state
	session := conv.account.session
	c, _ := session.ConversationManager().EnsureConversationWith(conv.to)
	err := c.StartEncryptedChat(session, conv.currentResource())
	if err != nil {
		//TODO: notify failure
	}
}

func (conv *conversationPane) onEndOtrSignal() {
	//TODO: errors
	//TODO: enable/disable depending on the conversation's encryption state
	session := conv.account.session
	c, ok := session.ConversationManager().GetConversationWith(conv.to)
	if !ok {
		return
	}

	err := c.EndEncryptedChat(session, conv.currentResource())
	if err != nil {
		log.Printf(i18n.Local("Failed to terminate the encrypted chat: %s\n"), err.Error())
	}
}

func (conv *conversationPane) onVerifyFpSignal() {
	switch verifyFingerprintDialog(conv.account, conv.to, conv.transientParent) {
	case gtki.RESPONSE_YES:
		conv.removeIdentityVerificationWarning()
	}
}

func (conv *conversationPane) onConnect() {
	conv.entry.SetEditable(true)
	conv.entry.SetSensitive(true)
}

func (conv *conversationPane) onDisconnect() {
	conv.entry.SetEditable(false)
	conv.entry.SetSensitive(false)
}

func countVisibleLines(v gtki.TextView) uint {
	lines := uint(1)
	iter := getTextBufferFrom(v).GetStartIter()
	for v.ForwardDisplayLine(iter) {
		lines++
	}

	return lines
}

func (conv *conversationPane) calculateHeight(lines uint) uint {
	return lines * 2 * getFontSizeFrom(conv.entry)
}

func (conv *conversationPane) doPotentialEntryResize() {
	lines := countVisibleLines(conv.entry)
	scroll := true
	if lines > 3 {
		scroll = false
		lines = 3
	}
	conv.entryScroll.SetProperty("height-request", conv.calculateHeight(lines))
	if scroll {
		scrollToTop(conv.entryScroll)
	}
}

func createConversationPane(account *account, uid string, ui *gtkUI, transientParent gtki.Window) *conversationPane {
	builder := newBuilder("ConversationPane")

	cp := &conversationPane{
		to:              uid,
		account:         account,
		transientParent: transientParent,
		afterNewMessage: func() {},
		currentPeer: func() (*rosters.Peer, bool) {
			return ui.getPeer(account, uid)
		},
	}

	builder.getItems(
		"box", &cp.widget,
		"history", &cp.history,
		"historyScroll", &cp.scrollHistory,
		"message", &cp.entry,
		"notification-area", &cp.notificationArea,
		"security-warning", &cp.securityWarning,
		"menubar", &cp.menubar,
		"messageScroll", &cp.entryScroll,
	)

	builder.ConnectSignals(map[string]interface{}{
		"on_start_otr_signal": cp.onStartOtrSignal,
		"on_end_otr_signal":   cp.onEndOtrSignal,
		"on_verify_fp_signal": cp.onVerifyFpSignal,
		"on_connect":          cp.onConnect,
		"on_disconnect":       cp.onDisconnect,
	})

	cp.entryScroll.SetProperty("height-request", cp.calculateHeight(1))
	cp.history.SetBuffer(ui.getTags().createTextBuffer())
	cp.history.Connect("size-allocate", func() {
		scrollToBottom(cp.scrollHistory)
	})

	cp.entry.Connect("key-release-event", cp.doPotentialEntryResize)

	ui.displaySettings.control(cp.history)
	ui.displaySettings.control(cp.entry)

	return cp
}

func (conv *conversationPane) connectEnterHandler(target gtki.Widget) {
	if target == nil {
		target = conv.entry
	}

	target.Connect("key-press-event", func(_ gtki.Widget, ev gdki.Event) bool {
		evk := g.gdk.EventKeyFrom(ev)
		ret := false

		if conv.account.isInsertEnter(evk) {
			insertEnter(conv.entry)
			ret = true
		} else if conv.account.isSend(evk) {
			conv.onSendMessageSignal()
			ret = true
		}

		return ret
	})
}

func (a *account) isInsertEnter(evk gdki.EventKey) bool {
	return hasShift(evk) && hasEnter(evk)
}

func (a *account) isSend(evk gdki.EventKey) bool {
	return !hasControlingModifier(evk) && hasEnter(evk)
}

func newConversationWindow(account *account, uid string, ui *gtkUI) *conversationWindow {
	builder := newBuilder("Conversation")
	win := builder.getObj("conversation").(gtki.Window)

	peer, ok := ui.accountManager.contacts[account].Get(uid)
	otherName := uid
	if ok {
		otherName = peer.NameForPresentation()
	}

	// TODO: Can we put the security rating here, maybe?
	title := fmt.Sprintf("%s <-> %s", account.session.GetConfig().Account, otherName)
	win.SetTitle(title)

	winBox := builder.getObj("box").(gtki.Box)

	cp := createConversationPane(account, uid, ui, win)
	winBox.PackStart(cp.widget, true, true, 0)

	conv := &conversationWindow{
		conversationPane: cp,
		win:              win,
	}

	cp.connectEnterHandler(conv.win)
	cp.afterNewMessage = conv.potentiallySetUrgent

	// Unlike the GTK version, this is not supposed to be used as a callback but
	// it attaches the callback to the widget
	conv.win.HideOnDelete()

	inEventHandler := false
	conv.win.Connect("set-focus", func() {
		if !inEventHandler {
			inEventHandler = true
			conv.entry.GrabFocus()
			inEventHandler = false
		}
	})

	conv.win.Connect("focus-in-event", func() {
		conv.unsetUrgent()
	})

	conv.win.Connect("notify::is-active", func() {
		if conv.win.IsActive() {
			inEventHandler = true
			conv.entry.GrabFocus()
			inEventHandler = false
		}
	})

	conv.win.Connect("hide", func() {
		conv.onHide()
	})

	conv.win.Connect("show", func() {
		conv.onShow()
	})

	ui.connectShortcutsChildWindow(conv.win)
	ui.connectShortcutsConversationWindow(conv)
	conv.parentWin = ui.window

	return conv
}

func (conv *conversationPane) addNotification(notification gtki.InfoBar) {
	conv.notificationArea.Add(notification)
}

func (conv *conversationWindow) Hide() {
	conv.win.Hide()
}

func (conv *conversationWindow) tryEnsureCorrectWorkspace() {
	if g.gdk.WorkspaceControlSupported() {
		wi, _ := conv.parentWin.GetWindow()
		parentPlace := wi.GetDesktop()
		cwi, _ := conv.win.GetWindow()
		cwi.MoveToDesktop(parentPlace)
	}
}

func (conv *conversationPane) getConversation() (client.Conversation, bool) {
	return conv.account.session.ConversationManager().GetConversationWith(conv.to)
}

func (conv *conversationPane) withCurrentPeer(f func(*rosters.Peer)) {
	p, ok := conv.currentPeer()
	if ok {
		f(p)
	}
}

func (conv *conversationPane) isVerified() bool {
	conversation, exists := conv.getConversation()
	if !exists {
		log.Println("Conversation does not exist - this shouldn't happen")
		return false
	}

	fingerprint := conversation.TheirFingerprint()
	conf := conv.account.session.GetConfig()

	p, hasPeer := conf.GetPeer(conv.to)

	if hasPeer {
		p.EnsureHasFingerprint(fingerprint)
	}

	return hasPeer && p.HasTrustedFingerprint(fingerprint)
}

func (conv *conversationPane) showIdentityVerificationWarning(u *gtkUI) {
	conv.Lock()
	defer conv.Unlock()

	if conv.fingerprintWarning != nil {
		log.Println("we are already showing a fingerprint warning, so not doing it again")
		return
	}

	if conv.isVerified() {
		log.Println("We have a peer and a trusted fingerprint already, so no reason to warn")
		return
	}

	conv.fingerprintWarning = buildVerifyIdentityNotification(conv.account, conv.to, conv.transientParent)
	conv.addNotification(conv.fingerprintWarning)
}

func (conv *conversationPane) removeIdentityVerificationWarning() {
	conv.Lock()
	defer conv.Unlock()

	if conv.fingerprintWarning != nil {
		conv.fingerprintWarning.Hide()
		conv.fingerprintWarning.Destroy()
		conv.fingerprintWarning = nil
	}
}

func (conv *conversationPane) updateSecurityWarning() {
	conversation, ok := conv.getConversation()
	if !ok {
		return
	}

	conv.securityWarning.SetVisible(!conversation.IsEncrypted())
}

func (conv *conversationWindow) show(userInitiated bool) {
	conv.updateSecurityWarning()
	conv.win.Show()
	conv.tryEnsureCorrectWorkspace()
}

func (conv *conversationPane) sendMessage(message string) error {
	err := conv.account.session.EncryptAndSendTo(conv.to, conv.currentResource(), message)

	if err != nil {
		oerr, isoff := err.(*access.OfflineError)
		if !isoff {
			return err
		}

		conv.displayNotification(oerr.Error())
	} else {
		//TODO: review whether it should create a conversation
		//TODO: this should be whether the message was encrypted or not, rather than
		//whether the conversation is encrypted or not
		conversation, _ := conv.account.session.ConversationManager().EnsureConversationWith(conv.to)
		conv.appendMessage(conv.account.session.GetConfig().Account, time.Now(), conversation.IsEncrypted(), ui.StripSomeHTML([]byte(message)), true)
	}
	return nil
}

const timeDisplay = "15:04:05"

// Expects to be called from the GUI thread.
// Expects to be called when conv is already locked
func insertAtEnd(buff gtki.TextBuffer, text string) {
	buff.Insert(buff.GetEndIter(), text)
}

// Expects to be called from the GUI thread.
// Expects to be called when conv is already locked
func insertWithTag(buff gtki.TextBuffer, tagName, text string) {
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

func scrollToBottom(sw gtki.ScrolledWindow) {
	adj := sw.GetVAdjustment()
	adj.SetValue(adj.GetUpper() - adj.GetPageSize())
}

func scrollToTop(sw gtki.ScrolledWindow) {
	adj := sw.GetVAdjustment()
	adj.SetValue(adj.GetLower())
}

type taggableText struct {
	tag  string
	text string
}

func (conv *conversationPane) appendToHistory(timestamp time.Time, entries ...taggableText) {
	conv.markNow()
	doInUIThread(func() {
		conv.Lock()
		defer conv.Unlock()

		buff, _ := conv.history.GetBuffer()
		if buff.GetCharCount() != 0 {
			insertAtEnd(buff, "\n")
		}

		insertAtEnd(buff, "[")
		insertAtEnd(buff, timestamp.Format(timeDisplay))
		insertAtEnd(buff, "] ")

		for _, entry := range entries {
			if entry.tag != "" {
				insertWithTag(buff, entry.tag, entry.text)
			} else {
				insertAtEnd(buff, entry.text)
			}
		}
		conv.afterNewMessage()
	})
}

func (conv *conversationPane) appendStatus(from string, timestamp time.Time, show, showStatus string, gone bool) {
	conv.appendToHistory(timestamp, taggableText{"statusText", createStatusMessage(from, show, showStatus, gone)})
}

func (conv *conversationPane) appendMessage(from string, timestamp time.Time, encrypted bool, message []byte, outgoing bool) {
	conv.appendToHistory(timestamp,
		taggableText{
			is(outgoing, "outgoingUser", "incomingUser"),
			from,
		},
		taggableText{
			text: ":  ",
		},
		taggableText{
			is(outgoing, "outgoingText", "incomingText"),
			string(message),
		})
}

func (conv *conversationPane) displayNotification(notification string) {
	conv.appendToHistory(time.Now(), taggableText{"statusText", notification})
}

func (conv *conversationPane) displayNotificationVerifiedOrNot(notificationV, notificationNV string) {
	if conv.isVerified() {
		conv.displayNotification(notificationV)
	} else {
		conv.displayNotification(notificationNV)
	}
}

func (conv *conversationWindow) setEnabled(enabled bool) {
	if enabled {
		conv.win.Emit("enable")
	} else {
		conv.win.Emit("disable")
	}
}

type timedMark struct {
	at     time.Time
	offset int
}

func (conv *conversationPane) markNow() {
	conv.Lock()
	defer conv.Unlock()

	buf, _ := conv.history.GetBuffer()
	offset := buf.GetCharCount()

	conv.marks = append(conv.marks, &timedMark{
		at:     time.Now(),
		offset: offset,
	})
}

const reapInterval = time.Duration(1) * time.Hour

func (conv *conversationPane) reapOlderThan(t time.Time) {
	conv.Lock()
	defer conv.Unlock()

	newMarks := []*timedMark{}
	var lastMark *timedMark
	isEnd := false

	for ix, m := range conv.marks {
		if t.Before(m.at) {
			newMarks = conv.marks[ix:]
			break
		}
		lastMark = m
		isEnd = len(conv.marks)-1 == ix
	}

	if lastMark != nil {
		off := lastMark.offset + 1
		buf, _ := conv.history.GetBuffer()
		sit := buf.GetStartIter()
		eit := buf.GetIterAtOffset(off)
		if isEnd {
			eit = buf.GetEndIter()
			newMarks = []*timedMark{}
		}

		buf.Delete(sit, eit)

		for _, nm := range newMarks {
			nm.offset = nm.offset - off
		}

		conv.marks = newMarks
	}
}

func (conv *conversationPane) onHide() {
	conv.reapOlderThan(time.Now().Add(-reapInterval))
	conv.hidden = true
}

func (conv *conversationPane) onShow() {
	if conv.hidden {
		conv.reapOlderThan(time.Now().Add(-reapInterval))
		conv.hidden = false
	}
}

func (conv *conversationWindow) potentiallySetUrgent() {
	if !conv.win.HasToplevelFocus() {
		conv.win.SetUrgencyHint(true)
	}
}

func (conv *conversationWindow) unsetUrgent() {
	conv.win.SetUrgencyHint(false)
}
