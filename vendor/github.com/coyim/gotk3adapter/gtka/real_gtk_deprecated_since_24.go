//+build gtk_3_6 gtk_3_8 gtk_3_10 gtk_3_12 gtk_3_14 gtk_3_16 gtk_3_18 gtk_3_20 gtk_3_22

package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

func (*RealGtk) CssProviderGetDefault() (gtki.CssProvider, error) {
	return wrapCssProvider(gtk.CssProviderGetDefault())
}
