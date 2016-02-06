package gui

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func applyHacks() {
	fixPopupMenusWithoutFocus()
}

// See #189
func fixPopupMenusWithoutFocus() {
	prov, _ := gtk.CssProviderNew()
	prov.LoadFromData("GtkMenu { margin: 0; }")

	// It must be added to the screen.
	// Adding to the main window has the same effect as putting the CSS in
	// gtk-keys.css (it is overwritten by the theme)
	screen, _ := gdk.ScreenGetDefault()
	gtk.AddProviderForScreen(screen, prov, gtk.STYLE_PROVIDER_PRIORITY_USER)
}
