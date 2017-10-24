package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/coyim/gotk3adapter/gtki"
)

type cssProvider struct {
	*gliba.Object
	internal *gtk.CssProvider
}

func wrapCssProviderSimple(v *gtk.CssProvider) *cssProvider {
	if v == nil {
		return nil
	}
	return &cssProvider{gliba.WrapObjectSimple(v.Object), v}
}

func wrapCssProvider(v *gtk.CssProvider, e error) (*cssProvider, error) {
	return wrapCssProviderSimple(v), e
}

func unwrapCssProvider(v gtki.CssProvider) *gtk.CssProvider {
	if v == nil {
		return nil
	}
	return v.(*cssProvider).internal
}

func (v *cssProvider) LoadFromData(v1 string) error {
	return v.internal.LoadFromData(v1)
}
