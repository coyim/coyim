package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type linkButton struct {
	*bin
	internal *gtk.LinkButton
}

func WrapLinkButtonSimple(v *gtk.LinkButton) gtki.LinkButton {
	if v == nil {
		return nil
	}
	return &linkButton{WrapBinSimple(&v.Bin).(*bin), v}
}

func WrapLinkButton(v *gtk.LinkButton, e error) (gtki.LinkButton, error) {
	return WrapLinkButtonSimple(v), e
}

func UnwrapLinkButton(v gtki.LinkButton) *gtk.LinkButton {
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
	v.internal.SetImage(UnwrapWidget(v1))
}
