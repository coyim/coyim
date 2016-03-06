package gui

import "github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtki"

func applyHacks() {
	fixPopupMenusWithoutFocus()
}

// See #189
func fixPopupMenusWithoutFocus() {
	prov, _ := g.gtk.CssProviderNew()
	prov.LoadFromData("GtkMenu { margin: 0; }")

	// It must be added to the screen.
	// Adding to the main window has the same effect as putting the CSS in
	// gtk-keys.css (it is overwritten by the theme)
	screen, _ := g.gdk.ScreenGetDefault()
	g.gtk.AddProviderForScreen(screen, prov, uint(gtki.STYLE_PROVIDER_PRIORITY_USER))
}
