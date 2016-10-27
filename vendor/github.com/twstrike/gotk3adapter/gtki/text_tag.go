package gtki

import "github.com/twstrike/gotk3adapter/glibi"

type TextTag interface {
	glibi.Object
}

func AssertTextTag(_ TextTag) {}
