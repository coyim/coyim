package gui

import (
	"github.com/coyim/gotk3adapter/gtka"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/coyim/gotk3extra"
)

type infoBarHighlightType int

const (
	infoBarHighlightFontWeight infoBarHighlightType = iota
)

type infoBarHighlightAttributes struct {
	labelWithFontWeight gtki.Label `gtk-widget:"labelWithFontWeight"`
}

func newInfoBarHighlightAttributes(tp infoBarHighlightType) *gotk3extra.PangoAttrList {
	ibh := &infoBarHighlightAttributes{}

	builder := newBuilder("InfoBarHighlightAttributes")
	panicOnDevError(builder.bindObjects(ibh))

	var highlightLabel gtki.Label
	switch tp {
	case infoBarHighlightFontWeight:
		highlightLabel = ibh.labelWithFontWeight
	}

	if highlightLabel != nil {
		return gotk3extra.LabelGetAttributes(gtka.UnwrapLabel(highlightLabel))
	}

	return nil
}
