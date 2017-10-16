package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

func unwrapStyleProvider(s gtki.StyleProvider) gtk.IStyleProvider {
	if s == nil {
		return nil
	}
	return unwrapCssProvider(s.(gtki.CssProvider))
}
