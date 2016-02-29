package gui

import "github.com/twstrike/gotk3adapter/gtki"

func (u *gtkUI) increaseFontSize(w gtki.Window) {
	u.displaySettings.increaseFontSize()
}

func (u *gtkUI) decreaseFontSize(w gtki.Window) {
	u.displaySettings.decreaseFontSize()
}

func (u *gtkUI) closeApplication(w gtki.Window) {
	u.quit()
}

func (u *gtkUI) closeWindow(w gtki.Window) {
	w.Hide()
}

func connectShortcut(accel string, w gtki.Window, action func(gtki.Window)) {
	gr, _ := g.gtk.AccelGroupNew()
	key, mod := g.gtk.AcceleratorParse(accel)

	// Do not remove the closure here - there is a limitation
	// in gtk that makes it necessary to have different functions for different accelerator groups
	gr.Connect2(key, mod, gtki.ACCEL_VISIBLE, func() {
		action(w)
	})

	w.AddAccelGroup(gr)
}

func (u *gtkUI) connectShortcutsMainWindow(w gtki.Window) {
	// <Primary> maps to Command on OS X, but Control on other platforms
	connectShortcut("<Primary>q", w, u.closeApplication)
	connectShortcut("<Primary>w", w, u.closeApplication)
	connectShortcut("<Alt>F4", w, u.closeApplication)
}

func (u *gtkUI) connectShortcutsChildWindow(w gtki.Window) {
	// <Primary> maps to Command on OS X, but Control on other platforms
	connectShortcut("<Primary>q", w, u.closeApplication)
	connectShortcut("<Primary>w", w, u.closeWindow)
	connectShortcut("<Primary>F4", w, u.closeWindow)
	connectShortcut("<Alt>F4", w, u.closeApplication)
	connectShortcut("Escape", w, u.closeWindow)
}

func (u *gtkUI) connectShortcutsConversationWindow(c *conversationWindow) {
	// <Primary> maps to Command on OS X, but Control on other platforms
	connectShortcut("<Primary>plus", c.win, u.increaseFontSize)
	connectShortcut("<Primary>minus", c.win, u.decreaseFontSize)
}
