package gtka

import (
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type textTag struct {
	*gliba.Object
	internal *gtk.TextTag
}

func WrapTextTagSimple(v *gtk.TextTag) gtki.TextTag {
	if v == nil {
		return nil
	}
	return &textTag{gliba.WrapObjectSimple(v.Object), v}
}

func WrapTextTag(v *gtk.TextTag, e error) (gtki.TextTag, error) {
	return WrapTextTagSimple(v), e
}

func UnwrapTextTag(v gtki.TextTag) *gtk.TextTag {
	if v == nil {
		return nil
	}
	return v.(*textTag).internal
}
