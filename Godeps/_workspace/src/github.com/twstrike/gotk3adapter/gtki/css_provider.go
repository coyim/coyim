package gtki

import "github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glibi"

type CssProvider interface {
	glibi.Object

	LoadFromData(string) error
}

func AssertCssProvider(_ CssProvider) {}
