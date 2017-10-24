package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/coyim/gotk3adapter/gtki"
)

type textTag struct {
	*gliba.Object
	internal *gtk.TextTag
}

func wrapTextTagSimple(v *gtk.TextTag) *textTag {
	if v == nil {
		return nil
	}
	return &textTag{gliba.WrapObjectSimple(v.Object), v}
}

func wrapTextTag(v *gtk.TextTag, e error) (*textTag, error) {
	return wrapTextTagSimple(v), e
}

func unwrapTextTag(v gtki.TextTag) *gtk.TextTag {
	if v == nil {
		return nil
	}
	return v.(*textTag).internal
}
