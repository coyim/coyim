package gui

import (
	"io/ioutil"
	"path/filepath"

	"github.com/coyim/coyim/config"
	"github.com/coyim/gotk3adapter/gtki"
)

var emacsKeyConf string

func init() {
	path := "Emacs/gtk-3.0/gtk-keys.css"

	toLook := []string{
		filepath.Join(config.XdgDataHome(), "themes", path),
		filepath.Join(config.WithHome(".themes"), path),
	}

	for _, dd := range config.XdgDataDirs() {
		toLook = append(toLook, filepath.Join(dd, "themes", path))
	}

	toLook = append(toLook, filepath.Join("/usr/share/themes", path))

	ek, ok := config.FindFile(toLook)
	if ok {
		content, _ := ioutil.ReadFile(ek)
		emacsKeyConf = string(content)
	}
}

type keyboardSettings struct {
	emacs    bool
	provider gtki.CssProvider
}

func (ks *keyboardSettings) control(w gtki.Widget) {
	doInUIThread(func() {
		styleContext, _ := w.GetStyleContext()
		styleContext.AddProvider(ks.provider, 9999)
	})
}

func (ks *keyboardSettings) update() {
	doInUIThread(func() {
		if ks.emacs {
			ks.provider.LoadFromData(emacsKeyConf)
		} else {
			ks.provider.LoadFromData("")
		}
	})
}

func newKeyboardSettings() *keyboardSettings {
	ks := &keyboardSettings{}
	prov, _ := g.gtk.CssProviderNew()
	ks.provider = prov
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
