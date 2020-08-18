package gui

import (
	"fmt"

	rosters "github.com/coyim/coyim/roster"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type conversationViewFactory interface {
	OpenConversationView(userInitiated bool) conversationView
	IfConversationView(whenExists func(conversationView), whenDoesntExist func())
}

type ourConversationViewFactory struct {
	account  *account
	peer     jid.Any
	ui       *gtkUI
	ul       *unifiedLayout
	targeted bool
}

func (u *gtkUI) NewConversationViewFactory(account *account, peer jid.Any, targeted bool) conversationViewFactory {
	return &ourConversationViewFactory{
		ui:       u,
		ul:       u.unified,
		account:  account,
		peer:     peer,
		targeted: targeted,
	}
}

func (cvf *ourConversationViewFactory) OpenConversationView(userInitiated bool) conversationView {
	// fmt.Printf("OpenConversationView(peer=%s, user=%v, targeted=%v)\n", cvf.peer, userInitiated, cvf.targeted)
	c, ok := cvf.getConversationViewSafely()
	if !ok {
		c = cvf.createConversationView(nil)
	}
	c.show(userInitiated)
	return c
}

func (cvf *ourConversationViewFactory) IfConversationView(whenExists func(conversationView), whenDoesntExist func()) {
	// fmt.Printf("IfConversationView(peer=%s)\n", cvf.peer)
	c, ok := cvf.getConversationViewSafely()
	// fmt.Printf("    IfConversationView(peer=%s) ok=%v\n", cvf.peer, ok)
	if ok {
		whenExists(c)
	} else {
		whenDoesntExist()
	}
}

func (cvf *ourConversationViewFactory) createConversationView(existing *conversationPane) conversationView {
	// fmt.Printf("createConversationView(peer=%s, targeted=%v)\n", cvf.peer, cvf.targeted)
	var cv conversationView

	if cvf.ui.settings.GetSingleWindow() {
		cv = cvf.createUnifiedConversationView(existing)
	} else {
		cv = cvf.createWindowedConversationView(existing)
	}
	cvf.setConversationView(cv)

	return cv
}

func (cvf *ourConversationViewFactory) potentialTarget() string {
	p := cvf.peer.PotentialResource().String()
	if cvf.targeted && p != "" {
		return fmt.Sprintf(" (%s)", p)
	}
	return ""
}

func windowConversationTitle(ui *gtkUI, peer jid.Any, account *account, potentialTarget string) string {
	// TODO: Can we put the security rating here, maybe?

	otherName := peer.String()
	if p, ok := ui.accountManager.getPeer(account, peer.NoResource()); ok {
		otherName = p.NameForPresentation()
	}

	return fmt.Sprintf("%s%s (%s)", otherName, potentialTarget, account.session.DisplayName())
}

func (cvf *ourConversationViewFactory) recreateWindowOn(conv *conversationWindow) {
	cp := conv.conversationPane

	builder := newBuilder("Conversation")
	win := builder.getObj("conversation").(gtki.Window)
	win.SetApplication(cvf.ui.app)

	win.SetTitle(windowConversationTitle(cvf.ui, cvf.peer, cvf.account, cvf.potentialTarget()))
	winBox := builder.getObj("box").(gtki.Box)
	cp.menubar.Show()
	winBox.PackStart(cp.widget, true, true, 0)

	conv.win = win
	cp.connectEnterHandler(conv.win)

	cvf.setWindowEvents(conv, winBox)

	cvf.ui.connectShortcutsChildWindow(conv.win)
	cvf.ui.connectShortcutsConversationWindow(conv)

	cp.transientParent = win

	// This 115 is apparently for the letter "s"
	win.AddMnemonic(uint(115), cp.encryptedLabel)
}

func (cvf *ourConversationViewFactory) setWindowEvents(conv *conversationWindow, winBox gtki.Box) {
	cp := conv.conversationPane

	inEventHandler := false
	_, _ = conv.win.Connect("set-focus", func() {
		if !inEventHandler {
			inEventHandler = true
			conv.entry.GrabFocus()
			inEventHandler = false
		}
	})

	_, _ = conv.win.Connect("focus-in-event", func() {
		conv.unsetUrgent()
	})

	_, _ = conv.win.Connect("delete-event", func() {
		winBox.Remove(cp.widget)
		conv.win = nil
	})

	_, _ = conv.win.Connect("notify::is-active", func() {
		if conv.win.IsActive() {
			inEventHandler = true
			conv.entry.GrabFocus()
			inEventHandler = false
		}
	})

	_, _ = conv.win.Connect("hide", func() {
		conv.onHide()
	})

	_, _ = conv.win.Connect("show", func() {
		conv.onShow()
	})

}

func (cvf *ourConversationViewFactory) createWindowedConversationView(existing *conversationPane) *conversationWindow {
	// fmt.Printf("createWindowedConversationView(peer=%s, targeted=%v)\n", cvf.peer, cvf.targeted)
	builder := newBuilder("Conversation")
	win := builder.getObj("conversation").(gtki.Window)
	win.SetApplication(cvf.ui.app)

	win.SetTitle(windowConversationTitle(cvf.ui, cvf.peer, cvf.account, cvf.potentialTarget()))
	winBox := builder.getObj("box").(gtki.Box)

	cp := cvf.createConversationPane(win)
	if existing != nil {
		b, _ := existing.history.GetBuffer()
		cp.history.SetBuffer(b)
	}

	cp.menubar.Show()
	winBox.PackStart(cp.widget, true, true, 0)

	conv := &conversationWindow{
		conversationPane: cp,
		win:              win,
	}

	cp.connectEnterHandler(conv.win)
	cp.afterNewMessage = conv.potentiallySetUrgent

	cvf.setWindowEvents(conv, winBox)

	cvf.ui.connectShortcutsChildWindow(conv.win)
	cvf.ui.connectShortcutsConversationWindow(conv)
	conv.parentWin = cvf.ui.window
	return conv
}

func (cvf *ourConversationViewFactory) createUnifiedConversationView(existing *conversationPane) conversationView {
	// fmt.Printf("createUnifiedConversationView(peer=%s, targeted=%v)\n", cvf.peer, cvf.targeted)
	cp := cvf.createConversationPane(cvf.ui.window)

	if existing != nil {
		b, _ := existing.history.GetBuffer()
		cp.history.SetBuffer(b)
	}

	cp.connectEnterHandler(nil)

	idx := cvf.ul.notebook.AppendPage(cp.widget, nil)
	if idx < 0 {
		panic("Failed to append page to notebook")
	}

	csi := &conversationStackItem{
		conversationPane: cp,
		pageIndex:        idx,
		layout:           cvf.ul,
	}

	//	csi.entry.SetHasFrame(true)
	csi.entryScroll.SetMarginTop(5)
	csi.entryScroll.SetMarginBottom(5)

	tabLabel := csi.shortName()
	resource := cvf.peer.PotentialResource().String()
	if resource != "" {
		tabLabel = tabLabel + " [at] " + resource
	}
	cvf.ul.notebook.SetTabLabelText(cp.widget, tabLabel)
	cvf.ul.itemMap[idx] = csi
	buffer, _ := csi.history.GetBuffer()
	_, _ = buffer.Connect("changed", func() {
		cvf.ul.onConversationChanged(csi)
	})
	return csi
}

func (cvf *ourConversationViewFactory) createConversationPane(win gtki.Window) *conversationPane {
	// fmt.Printf("createConversationPane(peer=%s, targeted=%v)\n", cvf.peer, cvf.targeted)
	builder := newBuilder("ConversationPane")

	var target jid.Any = cvf.peer.NoResource()
	if cvf.targeted {
		target = cvf.peer.(jid.WithResource)
	}

	cp := &conversationPane{
		isTargeted: cvf.targeted,
		target:     target,
		otrLock:    nil,

		account:              cvf.account,
		fileTransferNotif:    builder.fileTransferNotifInit(),
		securityWarningNotif: builder.securityWarningNotifInit(),
		transientParent:      win,
		shiftEnterSends:      cvf.ui.settings.GetShiftEnterForSend(),
		afterNewMessage:      func() {},
		delayed:              make(map[int]sentMessage),
		currentPeer: func() (*rosters.Peer, bool) {
			p, ok := cvf.ui.getPeer(cvf.account, cvf.peer.NoResource())
			if !ok {
				cvf.ui.log.WithField("peer", cvf.peer.NoResource().String()).Warn("Failure to look up peer from account")
			}
			return p, ok
		},
	}

	panicOnDevError(builder.bindObjects(cp))

	builder.ConnectSignals(map[string]interface{}{
		"on_start_otr":             cp.onStartOtrSignal,
		"on_end_otr":               cp.onEndOtrSignal,
		"on_verify_fp":             cp.onVerifyFpSignal,
		"on_connect":               cp.onConnect,
		"on_disconnect":            cp.onDisconnect,
		"on_destroy_file_transfer": cp.onDestroyFileTransferNotif,
		"on_send_file_to_contact": func() {
			cvf.account.sendFileTo(cp.currentPeerForSending(), cvf.ui, cp)
		},
		"on_send_dir_to_contact": func() {
			cvf.account.sendDirectoryTo(cp.currentPeerForSending(), cvf.ui, cp)
		},
	})

	// This 115 is apparently for the letter "s"
	win.AddMnemonic(uint(115), cp.encryptedLabel)

	_ = cp.entryScroll.SetProperty("height-request", cp.calculateHeight(1))

	prov := providerWithCSS("scrolledwindow { border-top: 2px solid #d3d3d3; } ")
	updateWithStyle(cp.entryScroll, prov)

	cp.history.SetBuffer(cvf.ui.getTags().createTextBuffer())
	_, _ = cp.history.Connect("size-allocate", func() {
		scrollToBottom(cp.scrollHistory)
	})

	cp.pending.SetBuffer(cvf.ui.getTags().createTextBuffer())

	_, _ = cp.entry.Connect("key-release-event", cp.doPotentialEntryResize)

	cvf.ui.displaySettings.control(cp.history)
	cvf.ui.displaySettings.shadeBackground(cp.pending)
	cvf.ui.displaySettings.control(cp.entry)
	cvf.ui.keyboardSettings.control(cp.entry)
	cvf.ui.keyboardSettings.update()

	cp.verifier = newVerifier(cvf.ui, cp)
	cp.encryptionStatus = &encryptionStatus{}

	return cp
}

func (cvf *ourConversationViewFactory) setConversationView(c conversationView) {
	// fmt.Printf("setConversationView(peer=%s)\n", cvf.peer)
	defer cvf.account.executeDelayed(cvf.ui, cvf.peer, cvf.targeted)

	cvf.account.Lock()
	defer cvf.account.Unlock()

	if cold, ok := cvf.account.c[c.getTarget().String()]; ok {
		cold.destroy()
	}

	// fmt.Printf("setConversationView(target=%s)\n", c.getTarget())
	cvf.account.c[c.getTarget().String()] = c
}

func (cvf *ourConversationViewFactory) isWindowingStyleConsistent(c conversationView) bool {
	unifiedLayout := cvf.ui.settings.GetSingleWindow()
	_, windowUnifiedLayout := c.(*conversationStackItem)
	return unifiedLayout == windowUnifiedLayout
}

func (cvf *ourConversationViewFactory) getConversationViewSafely() (conversationView, bool) {
	// fmt.Printf("getConversationViewSafely(peer=%s)\n", cvf.peer)
	c, ok := cvf.basicGetConversationView()
	// fmt.Printf("    getConversationViewSafely(peer=%s) ok=%v\n", cvf.peer, ok)
	if !ok {
		return nil, false
	}

	if val, okWin := c.(*conversationWindow); okWin && val.win == nil {
		cvf.recreateWindowOn(val)
	}

	if cvf.isWindowingStyleConsistent(c) {
		return c, true
	}

	defer c.destroy()

	var pane *conversationPane
	switch v := c.(type) {
	case *conversationWindow:
		pane = v.conversationPane
	case *conversationStackItem:
		pane = v.conversationPane
	}

	return cvf.createConversationView(pane), true
}

func (cvf *ourConversationViewFactory) countPeerWindows(peer jid.Any) int {
	c := 0
	for k := range cvf.account.c {
		if samePeer(peer, k) {
			c++
		}
	}
	return c
}

func (cvf *ourConversationViewFactory) basicGetConversationView() (conversationView, bool) {
	// fmt.Printf("basicGetConversationView(peer=%s)\n", cvf.peer)
	cvf.account.RLock()
	defer cvf.account.RUnlock()

	pw, pwo := jid.WithAndWithout(cvf.peer)
	// fmt.Printf("    basicGetConversationView(peer=%s) with=%v without=%v\n", cvf.peer, pw, pwo)

	if pw != nil {
		if c, ok := cvf.account.c[pw.String()]; ok {
			// This check is not strictly necessary - something should go VERY wrong if it triggers
			if !c.isOtrLocked() || c.isOtrLockedTo(cvf.peer) {
				return c, true
			}
		}
	}

	if c, ok := cvf.account.c[pwo.String()]; ok && !cvf.targeted && (!c.isOtrLocked() || c.isOtrLockedTo(cvf.peer) || cvf.countPeerWindows(pwo) == 1) {
		return c, true
	}

	return nil, false
}
