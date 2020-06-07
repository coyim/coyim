package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type textMark struct {
	*gtk.TextMark
}

func WrapTextMarkSimple(v *gtk.TextMark) gtki.TextMark {
	if v == nil {
		return nil
	}
	return &textMark{v}
}

func WrapTextMark(v *gtk.TextMark, e error) (gtki.TextMark, error) {
	return WrapTextMarkSimple(v), e
}

func UnwrapTextMark(v gtki.TextMark) *gtk.TextMark {
	if v == nil {
		return nil
	}
	return v.(*textMark).TextMark
}
