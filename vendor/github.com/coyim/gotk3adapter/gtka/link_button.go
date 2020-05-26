package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type linkButton struct {
	*bin
	internal *gtk.LinkButton
}

func wrapLinkButtonSimple(v *gtk.LinkButton) *linkButton {
	if v == nil {
		return nil
	}
	return &linkButton{wrapBinSimple(&v.Bin), v}
}

func wrapLinkButton(v *gtk.LinkButton, e error) (*linkButton, error) {
	return wrapLinkButtonSimple(v), e
}

func unwrapLinkButton(v gtki.LinkButton) *gtk.LinkButton {
	if v == nil {
		return nil
	}
	return v.(*linkButton).internal
}

func (v *linkButton) GetUri() string {
	return v.internal.GetUri()
}

func (v *linkButton) SetUri(uri string) {
	v.internal.SetUri(uri)
}

func (v *linkButton) SetImage(v1 gtki.Widget) {
	v.internal.SetImage(unwrapWidget(v1))
}
