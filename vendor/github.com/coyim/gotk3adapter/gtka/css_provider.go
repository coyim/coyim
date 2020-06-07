package gtka

import (
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type cssProvider struct {
	*gliba.Object
	internal *gtk.CssProvider
}

func WrapCssProviderSimple(v *gtk.CssProvider) gtki.CssProvider {
	if v == nil {
		return nil
	}
	return &cssProvider{gliba.WrapObjectSimple(v.Object), v}
}

func WrapCssProvider(v *gtk.CssProvider, e error) (gtki.CssProvider, error) {
	return WrapCssProviderSimple(v), e
}

func UnwrapCssProvider(v gtki.CssProvider) *gtk.CssProvider {
	if v == nil {
		return nil
	}
	return v.(*cssProvider).internal
}

func (v *cssProvider) LoadFromData(v1 string) error {
	return v.internal.LoadFromData(v1)
}
