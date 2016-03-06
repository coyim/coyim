package gtka

import (
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtki"
)

type textIter struct {
	*gtk.TextIter
}

func wrapTextIterSimple(v *gtk.TextIter) *textIter {
	if v == nil {
		return nil
	}
	return &textIter{v}
}

func wrapTextIter(v *gtk.TextIter, e error) (*textIter, error) {
	return wrapTextIterSimple(v), e
}

func unwrapTextIter(v gtki.TextIter) *gtk.TextIter {
	if v == nil {
		return nil
	}
	return v.(*textIter).TextIter
}
