package gui

import (
	"os"
	"path/filepath"

	"github.com/coyim/coyim/config"
	"github.com/coyim/gotk3adapter/gtki"
)

var emacsKeyConf string

func init() {
	path := "Emacs/gtk-3.0/gtk-keys.css"

	toLook := []string{
		filepath.Join(config.SystemDataDir(), "themes", path),
		filepath.Join(config.WithHome(".themes"), path),
	}

	for _, dd := range config.XdgDataDirs() {
		toLook = append(toLook, filepath.Join(dd, "themes", path))
	}

	toLook = append(toLook, filepath.Join("/usr/share/themes", path))

	ek, ok := config.FindFile(toLook)
	if ok {
		content, _ := os.ReadFile(filepath.Clean(ek))
		emacsKeyConf = string(content)
	}
}

type keyboardSettings struct {
	emacs    bool
	provider *cssProvider
}

func (ks *keyboardSettings) control(w gtki.Widget) {
	doInUIThread(func() {
		styleContext, _ := w.GetStyleContext()
		styleContext.AddProvider(ks.provider.provider, 9999)
	})
}

func (ks *keyboardSettings) update() {
	doInUIThread(func() {
		if ks.emacs {
			ks.provider.load("emacs key config", emacsKeyConf)
		} else {
			ks.provider.load("empty key config", "")
		}
	})
}

func newKeyboardSettings(wl withLog) *keyboardSettings {
	ks := &keyboardSettings{}
	ks.provider = newCSSProvider(wl)
	return ks
}

func (u *gtkUI) increaseFontSize(w gtki.Window) {
	u.displaySettings.increaseFontSize()
}

func (u *gtkUI) decreaseFontSize(w gtki.Window) {
	u.displaySettings.decreaseFontSize()
}

func (u *gtkUI) showSearchBar(w gtki.Window) {
	u.search.SetSearchMode(true)
	u.searchEntry.GrabFocus()
}

func (u *gtkUI) closeApplication(w gtki.Window) {
	u.quit()
}

func (u *gtkUI) closeWindow(w gtki.Window) {
	w.Hide()
}

func (u *gtkUI) closeApplicationOrConversation(w gtki.Window) {
	if u.settings.GetSingleWindow() {
		page := u.unified.notebook.GetCurrentPage()
		if page < 0 {
			u.quit()
		} else {
			u.unified.onCloseClicked()
		}
	} else {
		u.quit()
	}
}

func (u *gtkUI) closeWindowOrConversation(w gtki.Window) {
	if u.settings.GetSingleWindow() {
		page := u.unified.notebook.GetCurrentPage()
		if page < 0 {
			w.Hide()
		} else {
			u.unified.onCloseClicked()
		}
	} else {
		w.Hide()
	}
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
	connectShortcut("<Primary>F", w, u.showSearchBar)
	connectShortcut("<Primary>f", w, u.showSearchBar)
	connectShortcut("<Primary>q", w, u.closeApplication)
	connectShortcut("<Primary>w", w, u.closeApplicationOrConversation)
	connectShortcut("<Alt>F4", w, u.closeApplication)
}

func (u *gtkUI) connectShortcutsChildWindow(w gtki.Window) {
	// <Primary> maps to Command on OS X, but Control on other platforms
	connectShortcut("<Primary>q", w, u.closeApplication)
	connectShortcut("<Primary>w", w, u.closeWindowOrConversation)
	connectShortcut("<Primary>F4", w, u.closeWindow)
	connectShortcut("<Alt>F4", w, u.closeApplication)
	connectShortcut("Escape", w, u.closeWindow)
}

func (u *gtkUI) connectShortcutsConversationWindow(c *conversationWindow) {
	// <Primary> maps to Command on OS X, but Control on other platforms
	connectShortcut("<Primary>plus", c.win, u.increaseFontSize)
	connectShortcut("<Primary>minus", c.win, u.decreaseFontSize)
}

func (u *gtkUI) connectShortcutsMucRoomWindow(window gtki.Window, closeWindow func(_ gtki.Window)) {
	// <Primary> maps to Command on OS X, but Control on other platforms
	connectShortcut("<Primary>q", window, u.closeApplication)
	connectShortcut("<Primary>w", window, closeWindow)
	connectShortcut("<Primary>F4", window, closeWindow)
	connectShortcut("Escape", window, closeWindow)
}

func (u *gtkUI) connectShortcutsMucConfigRoomWindow(window gtki.Window, closeWindow func(_ gtki.Window)) {
	// <Primary> maps to Command on OS X, but Control on other platforms
	connectShortcut("<Primary>q", window, u.closeApplication)
	connectShortcut("<Primary>w", window, closeWindow)
	connectShortcut("<Primary>F4", window, closeWindow)
}
