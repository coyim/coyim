package gtki

import "github.com/coyim/gotk3adapter/pangoi"

type Label interface {
	Widget

	GetLabel() string
	SetLabel(string)
	SetSelectable(bool)
	SetText(string)
	SetMarkup(string)
	GetMnemonicKeyval() uint
	SetPangoAttributes(pangoi.PangoAttrList)
	GetPangoAttributes() pangoi.PangoAttrList
}

func AssertLabel(_ Label) {}
