package gtki

import "github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glibi"

type TextTagTable interface {
	glibi.Object

	Add(TextTag)
}

func AssertTextTagTable(_ TextTagTable) {}
