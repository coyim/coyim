package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type infoBar struct {
	*box
	internal *gtk.InfoBar
}

func WrapInfoBarSimple(v *gtk.InfoBar) gtki.InfoBar {
	if v == nil {
		return nil
	}
	return &infoBar{WrapBoxSimple(&v.Box).(*box), v}
}

func WrapInfoBar(v *gtk.InfoBar, e error) (gtki.InfoBar, error) {
	return WrapInfoBarSimple(v), e
}

func UnwrapInfoBar(v gtki.InfoBar) *gtk.InfoBar {
	if v == nil {
		return nil
	}
	return v.(*infoBar).internal
}
