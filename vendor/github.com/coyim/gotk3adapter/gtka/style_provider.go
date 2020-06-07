package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

func UnwrapStyleProvider(s gtki.StyleProvider) gtk.IStyleProvider {
	if s == nil {
		return nil
	}
	return UnwrapCssProvider(s.(gtki.CssProvider))
}
