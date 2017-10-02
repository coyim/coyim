package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gtki"
)

type textMark struct {
	*gtk.TextMark
}

func wrapTextMarkSimple(v *gtk.TextMark) *textMark {
	if v == nil {
		return nil
	}
	return &textMark{v}
}

func wrapTextMark(v *gtk.TextMark, e error) (*textMark, error) {
	return wrapTextMarkSimple(v), e
}

func unwrapTextMark(v gtki.TextMark) *gtk.TextMark {
	if v == nil {
		return nil
	}
	return v.(*textMark).TextMark
}
