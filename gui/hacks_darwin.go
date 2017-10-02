package gui

import "github.com/coyim/gotk3adapter/gtki"

func applyHacks() {
	fixPopupMenusWithoutFocus()
}

// See #189
func fixPopupMenusWithoutFocus() {
	prov, err := g.gtk.CssProviderNew()
	if err != nil {
		return
	}
	prov.LoadFromData("GtkMenu { margin: 0; }")

	// It must be added to the screen.
	// Adding to the main window has the same effect as putting the CSS in
	// gtk-keys.css (it is overwritten by the theme)
	screen, err := g.gdk.ScreenGetDefault()
	if err != nil {
		return
	}
	g.gtk.AddProviderForScreen(screen, prov, uint(gtki.STYLE_PROVIDER_PRIORITY_USER))
}
