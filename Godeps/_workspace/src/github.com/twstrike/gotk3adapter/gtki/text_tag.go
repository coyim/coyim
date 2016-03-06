package gtki

import "github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glibi"

type TextTag interface {
	glibi.Object
}

func AssertTextTag(_ TextTag) {}
