package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gtki"
)

type infoBar struct {
	*box
	internal *gtk.InfoBar
}

func wrapInfoBarSimple(v *gtk.InfoBar) *infoBar {
	if v == nil {
		return nil
	}
	return &infoBar{wrapBoxSimple(&v.Box), v}
}

func wrapInfoBar(v *gtk.InfoBar, e error) (*infoBar, error) {
	return wrapInfoBarSimple(v), e
}

func unwrapInfoBar(v gtki.InfoBar) *gtk.InfoBar {
	if v == nil {
		return nil
	}
	return v.(*infoBar).internal
}
