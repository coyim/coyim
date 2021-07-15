package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/coyim/gotk3adapter/pangoi"
)

type infoBarHighlightType int

const (
	infoBarHighlightNickname infoBarHighlightType = iota
	infoBarHighlightAffiliation
	infoBarHighlightRole
)

type infoBarHighlightAttributes struct {
	labelNickname    gtki.Label `gtk-widget:"labelNickname"`
	labelAffiliation gtki.Label `gtk-widget:"labelAffiliation"`
	labelRole        gtki.Label `gtk-widget:"labelRole"`
}

func newInfoBarHighlightAttributes(tp infoBarHighlightType) pangoi.PangoAttrList {
	ibh := &infoBarHighlightAttributes{}

	builder := newBuilder("InfoBarHighlightAttributes")
	panicOnDevError(builder.bindObjects(ibh))

	var highlightLabel gtki.Label
	switch tp {
	case infoBarHighlightNickname:
		highlightLabel = ibh.labelNickname
	case infoBarHighlightAffiliation:
		highlightLabel = ibh.labelAffiliation
	case infoBarHighlightRole:
		highlightLabel = ibh.labelRole
	}

	if highlightLabel != nil {
		return highlightLabel.GetPangoAttributes()
	}

	return nil
}
