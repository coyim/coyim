package gtka

import (
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/coyim/gotk3extra"
)

type iconTheme struct {
	*gliba.Object
	internal *gotk3extra.IconTheme
}

func WrapIconThemeSimple(v *gotk3extra.IconTheme) gtki.IconTheme {
	if v == nil {
		return nil
	}
	return &iconTheme{gliba.WrapObjectSimple(v.Object), v}
}

func WrapIconTheme(v *gotk3extra.IconTheme, e error) (gtki.IconTheme, error) {
	return WrapIconThemeSimple(v), e
}

func UnwrapIconTheme(v gtki.IconTheme) *gotk3extra.IconTheme {
	if v == nil {
		return nil
	}
	return v.(*iconTheme).internal
}

func (v *iconTheme) AddResourcePath(path string) {
	v.internal.AddResourcePath(path)
}

func (v *iconTheme) AppendSearchPath(path string) {
	v.internal.AppendSearchPath(path)
}

func (v *iconTheme) GetExampleIconName() string {
	return v.internal.GetExampleIconName()
}

func (v *iconTheme) HasIcon(name string) bool {
	return v.internal.HasIcon(name)
}

func (v *iconTheme) PrependSearchPath(path string) {
	v.internal.PrependSearchPath(path)
}
