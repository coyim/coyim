package gtki

import "github.com/coyim/gotk3adapter/glibi"

type TextTagTable interface {
	glibi.Object

	Add(TextTag)
}

func AssertTextTagTable(_ TextTagTable) {}
