package gui

import "github.com/coyim/gotk3adapter/gtki"

func applyHacks(wl withLog) {
	fixPopupMenusWithoutFocus(wl)
}

// See #189
func fixPopupMenusWithoutFocus(wl withLog) {
	prov := newCSSProvider(wl)
	prov.load("popup menu without margin", "GtkMenu { margin: 0; }")

	// It must be added to the screen.
	// Adding to the main window has the same effect as putting the CSS in
	// gtk-keys.css (it is overwritten by the theme)
	screen, err := g.gdk.ScreenGetDefault()
	if err != nil {
		return
	}
	g.gtk.AddProviderForScreen(screen, prov.provider, uint(gtki.STYLE_PROVIDER_PRIORITY_USER))
}
