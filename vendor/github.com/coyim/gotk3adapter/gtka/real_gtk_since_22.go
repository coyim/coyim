// +build !gtk_3_6,!gtk_3_8,!gtk_3_10,!gtk_3_12,!gtk_3_14,!gtk_3_16,!gtk_3_18,!gtk_3_20

package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/coyim/gotk3extra"
)

func (*RealGtk) InfoBarSetRevealed(infobar gtki.InfoBar, setting bool) {
	gotk3extra.InfoBarSetRevealed(UnwrapInfoBar(infobar), setting)
}

func (*RealGtk) InfoBarGetRevealed(infobar gtki.InfoBar) bool {
	return gotk3extra.InfoBarGetRevealed(UnwrapInfoBar(infobar))
}
