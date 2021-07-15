package gui

import (
	"github.com/coyim/coyim/text"
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

type infobarHighlightFormatter struct {
	text string
}

func newInfobarHighlightFormatter(text string) *infobarHighlightFormatter {
	return &infobarHighlightFormatter{text}
}

func (f *infobarHighlightFormatter) formatLabel(label gtki.Label) {
	if formatted, ok := text.ParseWithFormat(f.text); ok {
		text, formats := formatted.Join()
		label.SetText(text)

		pangoAttrList := g.pango.PangoAttrListNew()

		for _, format := range formats {
			if highlightType, ok := presentationFormatsHighlight[format.Format]; ok {
				copy := newInfoBarHighlightAttributes(highlightType)
				copyAttributesTo(pangoAttrList, copy, format.Start, format.Start+format.Length)
			}
		}

		label.SetPangoAttributes(pangoAttrList)
	} else {
		label.SetText(f.text)
	}
}

func copyAttributesTo(toAttrList, fromAttrList pangoi.PangoAttrList, startIndex, endIndex int) {
	for _, attr := range fromAttrList.GetAttributes() {
		attr.SetStartIndex(startIndex)
		attr.SetEndIndex(endIndex)

		toAttrList.InsertPangoAttribute(attr)
	}
}
