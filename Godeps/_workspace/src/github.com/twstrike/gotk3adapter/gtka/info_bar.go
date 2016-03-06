package gtka

import (
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtki"
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
