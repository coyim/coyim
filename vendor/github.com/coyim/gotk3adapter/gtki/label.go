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
	SetAttributes(pangoi.AttrList)
	GetAttributes() (pangoi.AttrList, error)
}

func AssertLabel(_ Label) {}
