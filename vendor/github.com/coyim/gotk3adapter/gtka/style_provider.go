package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gtki"
)

func unwrapStyleProvider(s gtki.StyleProvider) gtk.IStyleProvider {
	if s == nil {
		return nil
	}
	return unwrapCssProvider(s.(gtki.CssProvider))
}
