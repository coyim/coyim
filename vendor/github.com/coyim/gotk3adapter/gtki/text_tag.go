package gtki

import "github.com/coyim/gotk3adapter/glibi"

type TextTag interface {
	glibi.Object
}

func AssertTextTag(_ TextTag) {}
